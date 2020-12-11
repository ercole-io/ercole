// Copyright (c) 2020 Sorint.lab S.p.A.
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

	"github.com/bamzi/jobrunner"
	"github.com/ercole-io/ercole/data-service/database"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/utils"

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/model"
)

// HostDataServiceInterface is a interface that wrap methods used to manipulate and save data
type HostDataServiceInterface interface {
	// Init initialize the service
	Init()

	// UpdateHostInfo update the host informations using the provided hostdata
	UpdateHostInfo(hostdata model.HostDataBE) (interface{}, utils.AdvancedErrorInterface)

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

func (hds *HostDataService) Init() {
	jobrunner.Start()

	currentHostCleaningJob := &CurrentHostCleaningJob{hostDataService: hds, TimeNow: hds.TimeNow, Database: hds.Database, Config: hds.Config, Log: hds.Log}
	if err := jobrunner.Schedule(hds.Config.DataService.CurrentHostCleaningJob.Crontab, currentHostCleaningJob); err != nil {
		hds.Log.Errorf("Something went wrong scheduling CurrentHostCleaningJob: %v", err)
	}

	if hds.Config.DataService.CurrentHostCleaningJob.RunAtStartup {
		jobrunner.Now(currentHostCleaningJob)
	}

	archivedHostCleaningJob := &ArchivedHostCleaningJob{hostDataService: hds, TimeNow: hds.TimeNow, Database: hds.Database, Config: hds.Config, Log: hds.Log}
	if err := jobrunner.Schedule(hds.Config.DataService.ArchivedHostCleaningJob.Crontab, archivedHostCleaningJob); err != nil {
		hds.Log.Errorf("Something went wrong scheduling ArchivedHostCleaningJob: %v", err)
	}

	if hds.Config.DataService.ArchivedHostCleaningJob.RunAtStartup {
		jobrunner.Now(archivedHostCleaningJob)
	}

	oracleDbsLicensesHistory := &OracleDbsLicensesHistory{
		Database:        hds.Database,
		TimeNow:         hds.TimeNow,
		Config:          hds.Config,
		hostDataService: hds,
		Log:             hds.Log,
	}
	if err := jobrunner.Schedule("@daily", oracleDbsLicensesHistory); err != nil {
		hds.Log.Errorf("Something went wrong scheduling OracleDbsLicensesHistory: %v", err)
	}

	go func() {
		time.Sleep(time.Second * 3)
		jobrunner.Now(oracleDbsLicensesHistory)
	}()
}

// UpdateHostInfo saves the hostdata
// TODO move in its hosts.go file
func (hds *HostDataService) UpdateHostInfo(hostdata model.HostDataBE) (interface{}, utils.AdvancedErrorInterface) {
	var aerr utils.AdvancedErrorInterface

	hostdata.ServerVersion = hds.Version
	hostdata.Archived = false
	hostdata.CreatedAt = hds.TimeNow()
	hostdata.ServerSchemaVersion = model.SchemaVersion
	hostdata.ID = primitive.NewObjectIDFromTimestamp(hds.TimeNow())

	if hds.Config.DataService.EnablePatching {
		hostdata, aerr = hds.PatchHostData(hostdata)
		if aerr != nil {
			return nil, aerr
		}
	}

	if hostdata.Features.Oracle != nil {
		hds.oracleDatabasesChecks(hostdata.Info, hostdata.Features.Oracle)
	}

	_, aerr = hds.Database.ArchiveHost(hostdata.Hostname)
	if aerr != nil {
		return nil, aerr
	}

	if hds.Config.DataService.LogInsertingHostdata {
		hds.Log.Info(utils.ToJSON(hostdata))
	}
	res, aerr := hds.Database.InsertHostData(hostdata)
	if aerr != nil {
		return nil, aerr
	}

	alertHostDataInsertionURL := utils.NewAPIUrlNoParams(
		hds.Config.AlertService.RemoteEndpoint,
		hds.Config.AlertService.PublisherUsername,
		hds.Config.AlertService.PublisherPassword,
		"/queue/host-data-insertion/"+res.InsertedID.(primitive.ObjectID).Hex()).String()

	if resp, err := http.Post(alertHostDataInsertionURL, "application/json", bytes.NewReader([]byte{})); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "EVENT ENQUEUE")

	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, utils.NewAdvancedErrorPtr(utils.ErrEventEnqueue, "EVENT ENQUEUE")
	}

	return res.InsertedID, nil
}

// PatchHostData patch the hostdata using the pf stored in the db
func (hds *HostDataService) PatchHostData(hostdata model.HostDataBE) (model.HostDataBE, utils.AdvancedErrorInterface) {
	//Find the patch
	patch, err := hds.Database.FindPatchingFunction(hostdata.Hostname)
	if err != nil {
		return model.HostDataBE{}, err
	}

	//If patch is valid, apply the path the data
	if patch.Hostname == hostdata.Hostname && patch.Code != "" {
		if hds.Config.DataService.LogDataPatching {
			hds.Log.Printf("Patching %s hostdata with the patch %s\n", patch.Hostname, patch.ID)
		}

		return utils.PatchHostdata(patch, hostdata)
	}

	return hostdata, nil
}
