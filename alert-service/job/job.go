package job

import (
	"time"

	"github.com/bamzi/jobrunner"
	"github.com/ercole-io/ercole/v2/alert-service/database"
	"github.com/ercole-io/ercole/v2/alert-service/emailer"
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
	Emailer  emailer.Emailer
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

	removeAlertJob := RemoveAlertJob{Database: j.Database, Config: j.Config, Log: j.Log}
	if err := jobrunner.Schedule(j.Config.AlertService.RemoveAlertJob.Crontab, &removeAlertJob); err != nil {
		j.Log.Errorf("something went wrong scheduling removeAlertJob: %v", err)
	}

	if j.Config.AlertService.RemoveAlertJob.RunAtStartup {
		jobrunner.Now(&removeAlertJob)
	}

	reportAlertJob := ReportAlertJob{Database: j.Database, Config: j.Config, Log: j.Log, Emailer: j.Emailer}
	if err := jobrunner.Schedule(j.Config.AlertService.ReportAlertJob.Crontab, &reportAlertJob); err != nil {
		j.Log.Errorf("something went wrong scheduling reportAlertJob: %v", err)
	}

	if j.Config.AlertService.ReportAlertJob.RunAtStartup {
		jobrunner.Now(&reportAlertJob)
	}

	simulatedHostAlertJob := SimulatedHostAlertJob{Database: j.Database, Config: j.Config, Log: j.Log, Emailer: j.Emailer}
	jobrunner.Every(time.Minute*5, &simulatedHostAlertJob)
}
