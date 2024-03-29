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

package job

import (
	"time"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/data-service/database"
	"github.com/ercole-io/ercole/v2/logger"
)

// ArchivedHostCleaningJob is the job used to clean and remove old archived host
type ArchivedHostCleaningJob struct {
	// Database contains the database layer
	Database database.MongoDatabaseInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Log contains logger formatted
	Log logger.Logger
}

// Run archive every archived hostdata that is older than a amount
func (job *ArchivedHostCleaningJob) Run() {
	//Find the archived hosts older than ArchivedHostCleaningJob.HourThreshold hours
	ids, err := job.Database.FindOldArchivedHosts(job.TimeNow().Add(time.Duration(-job.Config.DataService.ArchivedHostCleaningJob.HourThreshold) * time.Hour))
	if err != nil {
		job.Log.Error(err)
		return
	}

	//For each host, archive the host
	for _, id := range ids {
		//Delete the host
		err := job.Database.DeleteHostData(id)
		if err != nil {
			job.Log.Error(err)
			return
		}

		job.Log.Infof("%s has been deleted because it have passed more than %d hours from the host data insertion", id, job.Config.DataService.ArchivedHostCleaningJob.HourThreshold)
	}
}
