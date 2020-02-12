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
	"bytes"
	"net/http"
	"time"

	"github.com/amreo/ercole-services/data-service/database"
	"github.com/bamzi/jobrunner"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/amreo/ercole-services/utils"

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/model"
)

// HostDataServiceInterface is a interface that wrap methods used to manipulate and save data
type HostDataServiceInterface interface {
	// Init initialize the service
	Init()

	// UpdateHostInfo update the host informations using the provided hostdata
	UpdateHostInfo(hostdata model.HostData) (interface{}, utils.AdvancedErrorInterface)

	// ArchiveHost archive the host
	// ArchiveHost(hostname string) utils.AdvancedError
}

// HostDataService is the concrete implementation of HostDataServiceInterface. It saves data to a MongoDB database
type HostDataService struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Version of the saved data
	Version string
	// Database contains the database layer
	Database database.MongoDatabaseInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Log contains logger formatted
	Log *logrus.Logger
}

// Init initializes the service and database
func (hds *HostDataService) Init() {
	//Start cron jobs
	jobrunner.Start()

	jobrunner.Schedule(hds.Config.DataService.CurrentHostCleaningJob.Crontab, &CurrentHostCleaningJob{hostDataService: hds, TimeNow: hds.TimeNow, Database: hds.Database, Config: hds.Config})
	jobrunner.Schedule(hds.Config.DataService.ArchivedHostCleaningJob.Crontab, &ArchivedHostCleaningJob{hostDataService: hds, TimeNow: hds.TimeNow, Database: hds.Database, Config: hds.Config})
	if hds.Config.DataService.CurrentHostCleaningJob.RunAtStartup {
		jobrunner.Now(&CurrentHostCleaningJob{hostDataService: hds, TimeNow: hds.TimeNow, Database: hds.Database, Config: hds.Config})
	}
	if hds.Config.DataService.ArchivedHostCleaningJob.RunAtStartup {
		jobrunner.Now(&ArchivedHostCleaningJob{hostDataService: hds, TimeNow: hds.TimeNow, Database: hds.Database, Config: hds.Config})
	}
}

// UpdateHostInfo saves the hostdata
func (hds *HostDataService) UpdateHostInfo(hostdata model.HostData) (interface{}, utils.AdvancedErrorInterface) {
	hostdata.ServerVersion = hds.Version
	hostdata.Archived = false
	hostdata.CreatedAt = hds.TimeNow()
	hostdata.SchemaVersion = model.SchemaVersion
	hostdata.ID = primitive.NewObjectIDFromTimestamp(hds.TimeNow())

	//Archive the host
	_, err := hds.Database.ArchiveHost(hostdata.Hostname)
	if err != nil {
		return nil, err
	}

	//Insert the host
	if hds.Config.DataService.LogInsertingHostdata {
		hds.Log.Info(utils.ToJSON(hostdata))
	}
	res, err := hds.Database.InsertHostData(hostdata)
	if err != nil {
		return nil, err
	}

	//Enqueue the insertion
	if resp, err := http.Post(
		utils.NewAPIUrlNoParams(
			hds.Config.AlertService.RemoteEndpoint,
			hds.Config.AlertService.PublisherUsername,
			hds.Config.AlertService.PublisherPassword,
			"/queue/host-data-insertion/"+res.InsertedID.(primitive.ObjectID).Hex(),
		).String(), "application/json", bytes.NewReader([]byte{})); err != nil {

		return nil, utils.NewAdvancedErrorPtr(err, "EVENT ENQUEUE")
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, utils.NewAdvancedErrorPtr(utils.ErrEventEnqueue, "EVENT ENQUEUE")
	}

	return res.InsertedID, nil
}
