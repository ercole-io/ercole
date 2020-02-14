// Copyright (c) 2019 Sorint.lab S.p.A.
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

	"github.com/amreo/ercole-services/config"
	"github.com/sirupsen/logrus"

	"github.com/amreo/ercole-services/alert-service/database"
	"github.com/amreo/ercole-services/utils"
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
	//Find the current hosts older than FreshnessCheck.DaysThreshold days
	hosts, err := job.Database.FindOldCurrentHosts(job.TimeNow().AddDate(0, 0, -job.Config.AlertService.FreshnessCheckJob.DaysThreshold))
	if err != nil {
		utils.LogErr(job.Log, err)
		return
	}

	//For each host, throw a NO_DATA alert
	for _, host := range hosts {
		//Throw a NO_DATA alert if the host doesn't already have a new NO_DATA alert
		if exist, err := job.Database.ExistNoDataAlertByHost(host); err != nil {
			utils.LogErr(job.Log, err)
			return
		} else if !exist {
			err = job.alertService.ThrowNoDataAlert(host, job.Config.AlertService.FreshnessCheckJob.DaysThreshold)
			if err != nil {
				utils.LogErr(job.Log, err)
				return
			}
		}
	}
}
