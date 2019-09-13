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

// Package service is a package that provides methods for manipulating host informations
package service

import (
	"log"
	"time"

	"github.com/amreo/ercole-services/utils"
)

// CurrentHostCleaningJob is the job used to clean and archive old current host
type CurrentHostCleaningJob struct {
	hostDataService *HostDataService
}

// Run archive every hostdata that is older than a amount
func (job *CurrentHostCleaningJob) Run() {
	//Find the current hosts older than CurrentHostCleaningJob.HourThreshold days
	hosts, err := job.hostDataService.Database.FindOldCurrentHosts(time.Now().Add(time.Duration(-job.hostDataService.Config.DataService.CurrentHostCleaningJob.HourThreshold) * time.Hour))
	if err != nil {
		utils.LogErr(err)
		return
	}

	//For each host, archive the host
	for _, host := range hosts {
		//Archive the host
		_, err := job.hostDataService.Database.ArchiveHost(host)
		if err != nil {
			utils.LogErr(err)
			return
		}
		log.Printf("%s has been moved because it have passed more than %d hours from last update", host, job.hostDataService.Config.DataService.CurrentHostCleaningJob.HourThreshold)
	}
}
