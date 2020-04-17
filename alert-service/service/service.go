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
	"fmt"
	"sync"
	"time"

	"github.com/amreo/ercole-services/alert-service/database"
	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	"github.com/bamzi/jobrunner"
	"github.com/sirupsen/logrus"
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
	// Log contains logger formatted
	Log *logrus.Logger
	// Emailer contains the emailer layer
	Emailer Emailer
}

// Init initializes the service and database
func (as *AlertService) Init(wg *sync.WaitGroup) {
	//Create a new queue
	as.Queue = hub.New()

	//Subscribe the alert-service
	sub := as.Queue.Subscribe(as.Config.AlertService.QueueBufferSize, model.TopicHostDataInsertion, model.TopicAlertInsertion)
	wg.Add(1)
	go func(s hub.Subscription) {
		as.Log.Info("Start alert-service/queue")
		for msg := range s.Receiver {
			as.ProcessMsg(msg)
		}
		as.Log.Info("Stop alert-service/queue")
		wg.Done()
	}(sub)

	//Start cron jobs
	jobrunner.Start()

	jobrunner.Schedule(as.Config.AlertService.FreshnessCheckJob.Crontab, &FreshnessCheckJob{alertService: as, TimeNow: as.TimeNow, Database: as.Database, Log: as.Log})
	if as.Config.AlertService.FreshnessCheckJob.RunAtStartup {
		jobrunner.Now(&FreshnessCheckJob{alertService: as, TimeNow: as.TimeNow, Database: as.Database, Log: as.Log})
	}
}

// HostDataInsertion inserts the host data insertion in the queue
func (as *AlertService) HostDataInsertion(id primitive.ObjectID) utils.AdvancedErrorInterface {
	as.Queue.Publish(hub.Message{
		Name: model.TopicHostDataInsertion,
		Fields: hub.Fields{
			"id": id,
		},
	})
	return nil
}

// AlertInsertion inserts a alert insertion in the queue
func (as *AlertService) AlertInsertion(alr model.Alert) utils.AdvancedErrorInterface {
	as.Queue.Publish(hub.Message{
		Name: model.TopicAlertInsertion,
		Fields: hub.Fields{
			"alert": alr,
		},
	})
	return nil
}

// ProcessMsg processes the message msg
func (as *AlertService) ProcessMsg(msg hub.Message) {
	if as.Config.AlertService.LogMessages {
		as.Log.Infof("RECEIVED EVENT %s: %s", msg.Topic(), utils.ToJSON(msg.Fields))
	}

	switch msg.Topic() {
	case model.TopicHostDataInsertion:
		as.ProcessHostDataInsertion(msg.Fields)
	case model.TopicAlertInsertion:
		as.ProcessAlertInsertion(msg.Fields)
	default:
		as.Log.Warnf("Received message with unknown topic: %s", msg.Topic())
	}
}

// ProcessHostDataInsertion processes the host data insertion event
func (as *AlertService) ProcessHostDataInsertion(params hub.Fields) {
	id := params["id"].(primitive.ObjectID)

	//Get the original data
	newData, err := as.Database.FindHostData(id)
	if err != nil {
		utils.LogErr(as.Log, err)
		return
	}

	//Get the previous data
	oldData, err := as.Database.FindMostRecentHostDataOlderThan(newData.Hostname(), newData.CreatedAt())
	if err != nil {
		utils.LogErr(as.Log, err)
		return
	}

	//Find the data difference and generate eventually alerts
	if err := as.DiffHostDataMapAndGenerateAlert(oldData, newData); err != nil {
		utils.LogErr(as.Log, err)
		return
	}
}

// ProcessAlertInsertion processes the alert insertion event
func (as *AlertService) ProcessAlertInsertion(params hub.Fields) {
	alert := params["alert"].(model.Alert)

	//Create the subject and message
	var subject string
	var message string
	if val, ok := alert.OtherInfo["Hostname"]; ok {
		subject = fmt.Sprintf("%s %s on %s", alert.AlertSeverity, alert.Description, val)
		message = fmt.Sprintf("Date: %s\nSeverity: %s\nHost: %s\nCode: %s\n%s", alert.Date, alert.AlertSeverity, val, alert.AlertCode, alert.Description)
	} else {
		subject = fmt.Sprintf("%s %s", alert.AlertSeverity, alert.Description)
		message = fmt.Sprintf("Date: %s\nSeverity: %s\nCode: %s\n%s", alert.Date, alert.AlertSeverity, alert.AlertCode, alert.Description)
	}

	// Send the email
	err := as.Emailer.SendEmail(subject, message, as.Config.AlertService.Emailer.To)
	if err != nil {
		utils.LogErr(as.Log, err)
		return
	}
}

// DiffHostDataMapAndGenerateAlert find the difference between the data and generate eventually alerts for such difference
func (as *AlertService) DiffHostDataMapAndGenerateAlert(oldData model.HostDataMap, newData model.HostDataMap) utils.AdvancedErrorInterface {
	newEnterpriseLicenseAlertThrown := false

	oldExtra := oldData.Extra()
	newExtra := newData.Extra()

	//Convert databases array to map
	oldDatabases := model.DatabaseMapArrayAsMap(oldExtra.Databases())
	newDatabases := model.DatabaseMapArrayAsMap(newExtra.Databases())

	//If the oldData is empty, fire a new server
	if oldData.Hostname() == "" {
		if err := as.ThrowNewServerAlert(newData.Hostname()); err != nil {
			return err
		}
	}

	//For each new database
	for _, newDb := range newDatabases {
		//Get the old database if exist
		var oldDb model.DatabaseMap
		if val, ok := oldDatabases[newDb.Name()]; ok {
			oldDb = val
		} else {
			oldDb = model.DatabaseMap{
				"Licenses": primitive.A{},
				"Features": primitive.A{},
			}
			// fire NEW_DATABASE alert
			if err := as.ThrowNewDatabaseAlert(newDb.Name(), newData.Hostname()); err != nil {
				return err
			}
		}

		oldDataInfo := oldData.Info()
		newDataInfo := oldData.Info()

		//Find new enterprises licenses
		if ((oldDataInfo.CPUCores() < newDataInfo.CPUCores()) || (!oldDb.HasEnterpriseLicense() && newDb.HasEnterpriseLicense())) && !newEnterpriseLicenseAlertThrown {
			if err := as.ThrowNewEnterpriseLicenseAlert(newData.Hostname()); err != nil {
				return err
			}
			newEnterpriseLicenseAlertThrown = true
		}

		fmt.Printf("%v", oldDb.Features())
		fmt.Printf("%v", newDb.Features())

		//Get the difference of features
		diff := model.DiffFeatureMap(oldDb.Features(), newDb.Features())

		//Extract from the diff the activated features
		activatedFeatures := []string{}
		for feature, val := range diff {
			if val == model.DiffFeatureActivated {
				activatedFeatures = append(activatedFeatures, feature)
			}
		}

		//Throw alert for activated features
		if len(activatedFeatures) > 0 {
			if err := as.ThrowActivatedFeaturesAlert(newDb.Name(), newData.Hostname(), activatedFeatures); err != nil {
				return err
			}
		}
	}

	return nil
}
