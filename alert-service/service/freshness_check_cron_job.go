package service

import (
	"time"

	"github.com/amreo/ercole-services/utils"
)

// FreshnessCheckJob is the job used to check the freshness of the current hosts
type FreshnessCheckJob struct {
	alertService *AlertService
}

// Run throws NO_DATA alert for each hosts that haven't sent a hostdata withing the FreshnessCheck.DaysThreshold
func (fcj *FreshnessCheckJob) Run() {
	//Find the current hosts older than FreshnessCheck.DaysThreshold days
	hosts, err := fcj.alertService.Database.FindOldCurrentHosts(time.Now().AddDate(0, 0, -fcj.alertService.Config.AlertService.FreshnessCheck.DaysThreshold))
	if err != nil {
		utils.LogErr(err)
		return
	}

	//For each host, throw a NO_DATA alert
	for _, host := range hosts {
		//Throw a NO_DATA alert if the host doesn't already have a new NO_DATA alert
		if exist, err := fcj.alertService.Database.ExistNoDataAlertByHost(host); err != nil {
			utils.LogErr(err)
			return
		} else if !exist {
			fcj.alertService.ThrowNoDataAlert(host, fcj.alertService.Config.AlertService.FreshnessCheck.DaysThreshold)
		}
	}
}
