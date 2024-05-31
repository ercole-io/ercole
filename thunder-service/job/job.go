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

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/thunder-service/database"
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

	ociDataRetrieveJob := &OciDataRetrieveJob{TimeNow: j.TimeNow, Database: j.Database, Config: j.Config, Log: j.Log}
	if err := jobrunner.Schedule(j.Config.ThunderService.OciDataRetrieveJob.Crontab, ociDataRetrieveJob); err != nil {
		j.Log.Errorf("Something went wrong scheduling OciDataRetrieveJob: %v", err)
	}

	if j.Config.ThunderService.OciDataRetrieveJob.RunAtStartup {
		jobrunner.Now(ociDataRetrieveJob)
	}

	ociRemoveOldDataObjectsJob := &OciRemoveOldDataObjectsJob{TimeNow: j.TimeNow, Database: j.Database, Config: j.Config, Log: j.Log}
	if err := jobrunner.Schedule(j.Config.ThunderService.OciRemoveOldDataObjectsJob.Crontab, ociRemoveOldDataObjectsJob); err != nil {
		j.Log.Errorf("Something went wrong scheduling OciRemoveOldDataObjectsJob: %v", err)
	}

	if j.Config.ThunderService.OciRemoveOldDataObjectsJob.RunAtStartup {
		jobrunner.Now(ociRemoveOldDataObjectsJob)
	}

	awsDataRetrieveJob := AwsDataRetrieveJob{Database: j.Database, Config: j.Config, Log: j.Log}
	if err := jobrunner.Schedule(j.Config.ThunderService.AwsDataRetrieveJob.Crontab, &awsDataRetrieveJob); err != nil {
		j.Log.Errorf("Something went wrong scheduling AwsDataRetrieveJob: %v", err)
	}

	if j.Config.ThunderService.AwsDataRetrieveJob.RunAtStartup {
		jobrunner.Now(&awsDataRetrieveJob)
	}

	gcpDataRetrieverJob := GcpDataRetrieveJob{j.Database, j.Config, j.Log, nil}
	if err := jobrunner.Schedule(j.Config.ThunderService.GcpDataRetrieveJob.Crontab, &gcpDataRetrieverJob); err != nil {
		j.Log.Errorf("Something went wrong scheduling GcpDataRetrieveJob: %v", err)
	}

	if j.Config.ThunderService.GcpDataRetrieveJob.RunAtStartup {
		jobrunner.Now(&gcpDataRetrieverJob)
	}
}
