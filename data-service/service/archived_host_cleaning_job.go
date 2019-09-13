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

// ArchivedHostCleaningJob is the job used to clean and remove old archived host
type ArchivedHostCleaningJob struct {
	hostDataService *HostDataService
}

// Run archive every archived hostdata that is older than a amount
func (job *ArchivedHostCleaningJob) Run() {
	//Find the archived hosts older than ArchivedHostCleaningJob.HourThreshold days
	ids, err := job.hostDataService.Database.FindOldArchivedHosts(time.Now().Add(time.Duration(-job.hostDataService.Config.DataService.ArchivedHostCleaningJob.HourThreshold) * time.Hour))
	if err != nil {
		utils.LogErr(err)
		return
	}

	//For each host, archive the host
	for _, id := range ids {
		//Delete the host
		err := job.hostDataService.Database.DeleteHostData(id)
		if err != nil {
			utils.LogErr(err)
			return
		}
		log.Printf("%s has been deleted because it have passed more than %d hours from the host data insertion", id, job.hostDataService.Config.DataService.ArchivedHostCleaningJob.HourThreshold)
	}
}
