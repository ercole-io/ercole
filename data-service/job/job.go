// Package service is a package that provides methods for manipulating host informations
package job

import (
	"time"

	"github.com/bamzi/jobrunner"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/data-service/database"
	"github.com/sirupsen/logrus"
)

type JobInterface interface {
	Init()
}

type Job struct {
	Config        config.Configuration
	ServerVersion string
	Database      database.MongoDatabaseInterface
	TimeNow       func() time.Time
	Log           *logrus.Logger
}

func (job *Job) Init() {
	jobrunner.Start()

	currentHostCleaningJob := &CurrentHostCleaningJob{TimeNow: job.TimeNow, Database: job.Database, Config: job.Config, Log: job.Log}
	if err := jobrunner.Schedule(job.Config.DataService.CurrentHostCleaningJob.Crontab, currentHostCleaningJob); err != nil {
		job.Log.Errorf("Something went wrong scheduling CurrentHostCleaningJob: %v", err)
	}

	if job.Config.DataService.CurrentHostCleaningJob.RunAtStartup {
		jobrunner.Now(currentHostCleaningJob)
	}

	archivedHostCleaningJob := &ArchivedHostCleaningJob{TimeNow: job.TimeNow, Database: job.Database, Config: job.Config, Log: job.Log}
	if err := jobrunner.Schedule(job.Config.DataService.ArchivedHostCleaningJob.Crontab, archivedHostCleaningJob); err != nil {
		job.Log.Errorf("Something went wrong scheduling ArchivedHostCleaningJob: %v", err)
	}

	if job.Config.DataService.ArchivedHostCleaningJob.RunAtStartup {
		jobrunner.Now(archivedHostCleaningJob)
	}

	oracleDbsLicensesHistory := &OracleDbsLicensesHistory{
		Database: job.Database,
		TimeNow:  job.TimeNow,
		Config:   job.Config,
		Log:      job.Log,
	}

	jobrunner.Every(5*time.Minute, oracleDbsLicensesHistory)
}
