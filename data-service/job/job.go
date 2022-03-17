// Copyright (c) 2022 Sorint.lab S.p.A.
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
package job

import (
	"time"

	"github.com/bamzi/jobrunner"
	"go.mongodb.org/mongo-driver/bson/primitive"

	alert_service_client "github.com/ercole-io/ercole/v2/alert-service/client"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/data-service/database"
	"github.com/ercole-io/ercole/v2/logger"
)

type JobInterface interface {
	Init()
}

type Job struct {
	Config        config.Configuration
	ServerVersion string
	Database      database.MongoDatabaseInterface
	TimeNow       func() time.Time
	Log           logger.Logger
}

func (j *Job) Init() {
	jobrunner.Start()

	currentHostCleaningJob := &CurrentHostCleaningJob{TimeNow: j.TimeNow, Database: j.Database, Config: j.Config, Log: j.Log}
	if err := jobrunner.Schedule(j.Config.DataService.CurrentHostCleaningJob.Crontab, currentHostCleaningJob); err != nil {
		j.Log.Errorf("Something went wrong scheduling CurrentHostCleaningJob: %v", err)
	}

	if j.Config.DataService.CurrentHostCleaningJob.RunAtStartup {
		jobrunner.Now(currentHostCleaningJob)
	}

	archivedHostCleaningJob := &ArchivedHostCleaningJob{TimeNow: j.TimeNow, Database: j.Database, Config: j.Config, Log: j.Log}
	if err := jobrunner.Schedule(j.Config.DataService.ArchivedHostCleaningJob.Crontab, archivedHostCleaningJob); err != nil {
		j.Log.Errorf("Something went wrong scheduling ArchivedHostCleaningJob: %v", err)
	}

	if j.Config.DataService.ArchivedHostCleaningJob.RunAtStartup {
		jobrunner.Now(archivedHostCleaningJob)
	}

	freshnessJob := &FreshnessCheckJob{
		TimeNow:        j.TimeNow,
		Database:       j.Database,
		AlertSvcClient: alert_service_client.NewClient(j.Config.AlertService),
		Config:         j.Config,
		Log:            j.Log,
		NewObjectID: func() primitive.ObjectID {
			return primitive.NewObjectIDFromTimestamp(j.TimeNow())
		},
	}
	if err := jobrunner.Schedule(j.Config.DataService.FreshnessCheckJob.Crontab, freshnessJob); err != nil {
		j.Log.Errorf("Something went wrong scheduling FreshnessCheckJob: %v", err)
	}

	if j.Config.DataService.FreshnessCheckJob.RunAtStartup {
		jobrunner.Now(freshnessJob)
	}

	historicizeLicensesComplianceJob := &HistoricizeLicensesComplianceJob{
		Database: j.Database,
		TimeNow:  j.TimeNow,
		Config:   j.Config,
		Log:      j.Log,
	}
	jobrunner.Every(5*time.Minute, historicizeLicensesComplianceJob)
}
