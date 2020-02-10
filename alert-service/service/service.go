// Copyright (c) 2019 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package service is a package that provides methods for manipulating host informations
package service

import (
	"log"
	"sync"
	"time"

	"github.com/amreo/ercole-services/alert-service/database"
	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	"github.com/bamzi/jobrunner"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/amreo/ercole-services/config"

	"github.com/leandro-lugaresi/hub"
)

// AlertServiceInterface is a interface that wrap methods used to insert and process alert messages
type AlertServiceInterface interface {
	// Init initializes the service
	Init(wg *sync.WaitGroup)
	// HostDataInsertion inserts the host data insertion in the queue
	HostDataInsertion(id primitive.ObjectID) utils.AdvancedErrorInterface
	// ProcessMsg processes the message msg
	ProcessMsg(msg hub.Message)
	// ThrowNewDatabaseAlert create and insert in the database a new NEW_DATABASE alert
	ThrowNewDatabaseAlert(dbname string, hostname string) utils.AdvancedErrorInterface
	// ThrowNewServerAlert create and insert in the database a new NEW_SERVER alert
	ThrowNewServerAlert(hostname string) utils.AdvancedErrorInterface
	// ThrowNewEnterpriseLicenseAlert create and insert in the database a new NEW_DATABASE alert
	ThrowNewEnterpriseLicenseAlert(hostname string) utils.AdvancedErrorInterface
	// ThrowActivatedFeaturesAlert create and insert in the database a new NEW_OPTION alert
	ThrowActivatedFeaturesAlert(dbname string, hostname string, activatedFeatures []string) utils.AdvancedErrorInterface
	// ThrowNoDataAlert create and insert in the database a new NO_DATA alert
	ThrowNoDataAlert(hostname string, freshnessThreshold int) utils.AdvancedErrorInterface
}

// AlertService is the concrete implementation of HostDataServiceInterface. It saves data to a MongoDB database
type AlertService struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Queue that contains all messages to be processed
	Queue *hub.Hub
	// Database contains the database layer
	Database database.MongoDatabaseInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
}

// Init initializes the service and database
func (as *AlertService) Init(wg *sync.WaitGroup) {
	//Create a new queue
	as.Queue = hub.New()

	//Subscribe the alert-service
	sub := as.Queue.Subscribe(0, "hostdata.insertion")
	wg.Add(1)
	go func(s hub.Subscription) {
		log.Println("Start alert-service/queue")
		for msg := range s.Receiver {
			as.ProcessMsg(msg)
		}
		log.Println("Stop alert-service/queue")
		wg.Done()
	}(sub)

	//Start cron jobs
	jobrunner.Start()

	jobrunner.Schedule(as.Config.AlertService.FreshnessCheckJob.Crontab, &FreshnessCheckJob{alertService: as, TimeNow: as.TimeNow, Database: as.Database})
	if as.Config.AlertService.FreshnessCheckJob.RunAtStartup {
		jobrunner.Now(&FreshnessCheckJob{alertService: as, TimeNow: as.TimeNow, Database: as.Database})
	}
}

// HostDataInsertion inserts the host data insertion in the queue
func (as *AlertService) HostDataInsertion(id primitive.ObjectID) utils.AdvancedErrorInterface {
	as.Queue.Publish(hub.Message{
		Name: "hostdata.insertion",
		Fields: hub.Fields{
			"id": id,
		},
	})
	return nil
}

// ProcessMsg processes the message msg
func (as *AlertService) ProcessMsg(msg hub.Message) {
	if as.Config.AlertService.LogMessages {
		log.Printf("RECEIVED EVENT %s: %s", msg.Topic(), utils.ToJSON(msg.Fields))
	}

	switch msg.Topic() {
	case model.TopicHostDataInsertion:
		as.ProcessHostDataInsertion(msg.Fields)
	default:
	}
}

// ProcessHostDataInsertion processes the host data insertion event
func (as *AlertService) ProcessHostDataInsertion(params hub.Fields) {
	id := params["id"].(primitive.ObjectID)

	//Get the original data
	newData, err := as.Database.FindHostData(id)
	if err != nil {
		utils.LogErr(err)
		return
	}

	//Get the previous data
	oldData, err := as.Database.FindMostRecentHostDataOlderThan(newData.Hostname, newData.CreatedAt)
	if err != nil {
		utils.LogErr(err)
		return
	}

	//Find the data difference and generate eventually alerts
	if err := as.DiffHostDataAndGenerateAlert(oldData, newData); err != nil {
		utils.LogErr(err)
		return
	}
}

// DiffHostDataAndGenerateAlert find the difference between the data and generate eventually alerts for such difference
func (as *AlertService) DiffHostDataAndGenerateAlert(oldData model.HostData, newData model.HostData) utils.AdvancedErrorInterface {
	newEnterpriseLicenseAlertThrown := false
	//Modify the data to make the comparison more easier
	if oldData.Extra.Databases == nil {
		oldData.Extra.Databases = []model.Database{}
	}
	if newData.Extra.Databases == nil {
		newData.Extra.Databases = []model.Database{}
	}

	//Convert databases array to map
	oldDatabases := model.DatabasesArrayAsMap(oldData.Extra.Databases)
	newDatabases := model.DatabasesArrayAsMap(newData.Extra.Databases)

	//If the oldData is empty, fire a new server
	if oldData.Hostname == "" {
		if err := as.ThrowNewServerAlert(newData.Hostname); err != nil {
			return err
		}
	}

	//For each new database
	for _, newDb := range newDatabases {
		//Get the old database if exist
		var oldDb model.Database
		if val, ok := oldDatabases[newDb.Name]; ok {
			oldDb = val
		} else {
			oldDb = model.Database{
				Features: []model.Feature{},
			}
			// fire NEW_DATABASE alert
			if err := as.ThrowNewDatabaseAlert(newDb.Name, newData.Hostname); err != nil {
				return err
			}
		}

		//Find new enterprises licenses
		if ((oldData.Info.CPUCores < newData.Info.CPUCores) || (!model.HasEnterpriseLicense(oldDb) && model.HasEnterpriseLicense(newDb))) && !newEnterpriseLicenseAlertThrown {
			if err := as.ThrowNewEnterpriseLicenseAlert(newData.Hostname); err != nil {
				return err
			}
			newEnterpriseLicenseAlertThrown = true
		}

		//Get the difference of features
		diff := model.DiffFeature(oldDb.Features, newDb.Features)

		//Extract from the diff the activated features
		activatedFeatures := []string{}
		for feature, val := range diff {
			if val == model.DiffFeatureActivated {
				activatedFeatures = append(activatedFeatures, feature)
			}
		}

		//Throw alert for activated features
		if len(activatedFeatures) > 0 {
			if err := as.ThrowActivatedFeaturesAlert(newDb.Name, newData.Hostname, activatedFeatures); err != nil {
				return err
			}
		}
	}

	return nil
}
