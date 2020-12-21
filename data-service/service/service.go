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
	"time"

	"github.com/bamzi/jobrunner"
	"github.com/ercole-io/ercole/v2/data-service/database"
	"github.com/sirupsen/logrus"

	"github.com/ercole-io/ercole/v2/utils"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
)

type HostDataServiceInterface interface {
	Init()

	InsertHostData(hostdata model.HostDataBE) (interface{}, utils.AdvancedErrorInterface)
}

type HostDataService struct {
	Config        config.Configuration
	ServerVersion string
	Database      database.MongoDatabaseInterface
	TimeNow       func() time.Time
	Log           *logrus.Logger
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
	jobrunner.Every(5*time.Minute, oracleDbsLicensesHistory)
}
