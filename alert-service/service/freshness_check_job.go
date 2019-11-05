package service

import (
	"time"

	"github.com/amreo/ercole-services/utils"
)

// FreshnessCheckJob is the job used to check the freshness of the current hosts
type FreshnessCheckJob struct {
	alertService *AlertService
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
}

// Run throws NO_DATA alert for each hosts that haven't sent a hostdata withing the FreshnessCheck.DaysThreshold
func (job *FreshnessCheckJob) Run() {
	//Find the current hosts older than FreshnessCheck.DaysThreshold days
	hosts, err := job.alertService.Database.FindOldCurrentHosts(job.TimeNow().AddDate(0, 0, -job.alertService.Config.AlertService.FreshnessCheckJob.DaysThreshold))
	if err != nil {
		utils.LogErr(err)
		return
	}

	//For each host, throw a NO_DATA alert
	for _, host := range hosts {
		//Throw a NO_DATA alert if the host doesn't already have a new NO_DATA alert
		if exist, err := job.alertService.Database.ExistNoDataAlertByHost(host); err != nil {
			utils.LogErr(err)
			return
		} else if !exist {
			job.alertService.ThrowNoDataAlert(host, job.alertService.Config.AlertService.FreshnessCheckJob.DaysThreshold)
		}
	}
}
