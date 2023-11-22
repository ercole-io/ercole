package job

import (
	"github.com/bamzi/jobrunner"
	"github.com/ercole-io/ercole/v2/alert-service/database"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
)

type JobInterface interface {
	Init()
}

type Job struct {
	Config   config.Configuration
	Database database.MongoDatabaseInterface
	Log      logger.Logger
}

func (j *Job) Init() {
	j.Log.Infof("init alert-service jobs")

	jobrunner.Start()

	ackAlertJob := AckAlertJob{Database: j.Database, Config: j.Config, Log: j.Log}
	if err := jobrunner.Schedule(j.Config.AlertService.AckAlertJob.Crontab, &ackAlertJob); err != nil {
		j.Log.Errorf("something went wrong scheduling ackAlertJob: %v", err)
	}

	if j.Config.AlertService.AckAlertJob.RunAtStartup {
		jobrunner.Now(&ackAlertJob)
	}
}
