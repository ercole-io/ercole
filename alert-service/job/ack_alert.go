package job

import (
	"github.com/ercole-io/ercole/v2/alert-service/database"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
)

type AckAlertJob struct {
	Database database.MongoDatabaseInterface
	Config   config.Configuration
	Log      logger.Logger
}

func (j *AckAlertJob) Run() {
	res, err := j.Database.AckOldAlerts()
	if err != nil {
		j.Log.Errorf("ack alert job", err)
		return
	}

	j.Log.Infof("matched %v documents and modified %v document\n", res.MatchedCount, res.ModifiedCount)
}
