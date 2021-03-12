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
	"github.com/sirupsen/logrus"
)

// CurrentHostCleaningJob is the job used to clean and archive old current host
type CurrentHostCleaningJob struct {
	// Database contains the database layer
	Database database.MongoDatabaseInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Log contains logger formatted
	Log *logrus.Logger
}

// Run archive every hostdata that is older than a amount
func (job *CurrentHostCleaningJob) Run() {
	timeLimit := job.TimeNow().Add(time.Duration(-job.Config.DataService.CurrentHostCleaningJob.HourThreshold) * time.Hour)
	hosts, err := job.Database.FindOldCurrentHostnames(timeLimit)
	if err != nil {
		job.Log.Error(err)
		return
	}

	for _, host := range hosts {
		_, err := job.Database.ArchiveHost(host)
		if err != nil {
			job.Log.Error(err)
			return
		}
		job.Log.Infof("%s has been moved because it have passed more than %d hours from last update", host,
			job.Config.DataService.CurrentHostCleaningJob.HourThreshold)
	}
}
