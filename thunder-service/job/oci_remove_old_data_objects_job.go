// Copyright (c) 2021 Sorint.lab S.p.A.
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
	"github.com/ercole-io/ercole/v2/logger"
	db "github.com/ercole-io/ercole/v2/thunder-service/database"
)

// OciRemoveOldDataObjectsJob is the job used to retrieve data from Oracle Cloud installation
type OciRemoveOldDataObjectsJob struct {
	// Database contains the database layer
	Database db.MongoDatabaseInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Log contains logger formatted
	Log logger.Logger
}

// Run job that remove data object older than 5 days
func (job *OciRemoveOldDataObjectsJob) Run() {

	currentTime := time.Now().UTC()
	dateFrom := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()-5, 0, 0, 0, 0, currentTime.Location())

	err := job.Database.DeleteOldOciObjects(dateFrom)
	if err != nil {
		job.Log.Error(err)
	}
}
