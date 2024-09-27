package job

import (
	"github.com/ercole-io/ercole/v2/alert-service/database"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
)

type RemoveAlertJob struct {
	Database database.MongoDatabaseInterface
	Config   config.Configuration
	Log      logger.Logger
}

func (j *RemoveAlertJob) Run() {
	res, err := j.Database.RemoveOldAlerts(j.Config.AlertService.AckAlertJob.DueDays)
	if err != nil {
		j.Log.Errorf("remove alert job", err)
		return
	}

	j.Log.Infof("removed %d documents", res.DeletedCount)
}
