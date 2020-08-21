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

package service

import (
	"time"

	"github.com/ercole-io/ercole/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/alert-service/database"
	"github.com/ercole-io/ercole/utils"
)

// FreshnessCheckJob is the job used to check the freshness of the current hosts
type FreshnessCheckJob struct {
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Database contains the database layer
	Database database.MongoDatabaseInterface
	// Config contains the dataservice global configuration
	Config config.Configuration
	// alertService contains the underlyng alert service
	alertService AlertServiceInterface
	// Log contains logger formatted
	Log *logrus.Logger
}

// Run throws NO_DATA alert for each hosts that haven't sent a hostdata withing the FreshnessCheck.DaysThreshold
func (job *FreshnessCheckJob) Run() {
	if job.Config.AlertService.FreshnessCheckJob.DaysThreshold <= 0 {
		job.Log.Errorf("AlertService.FreshnessCheckJob.DaysThreshold must be higher than 0, but it's set to %v. Job failed.",
			job.Config.AlertService.FreshnessCheckJob.DaysThreshold)

		return
	}

	if err := job.Database.DeleteAllNoDataAlerts(); err != nil {
		utils.LogErr(job.Log, err)
		return
	}

	hosts, err := job.Database.FindOldCurrentHosts(
		job.TimeNow().AddDate(0, 0, -job.Config.AlertService.FreshnessCheckJob.DaysThreshold))
	if err != nil {
		utils.LogErr(job.Log, err)
		return
	}

	//TODO remove mongo-driver dependency from here...
	for _, host := range hosts {
		if t, ok := host["createdAt"].(primitive.DateTime); !ok {
			panic("Can't obtain createdAt value")
		} else {

			elapsed := job.TimeNow().Sub(t.Time())
			elapsedDays := int(elapsed.Truncate(time.Hour*24).Hours() / 24)

			hostname := host["hostname"].(string)
			err := job.alertService.ThrowNoDataAlert(hostname, elapsedDays)
			if err != nil {
				utils.LogErr(job.Log, err)
				return
			}
		}
	}
}
