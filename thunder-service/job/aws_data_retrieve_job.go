// Copyright (c) 2022 Sorint.lab S.p.A.
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
	"github.com/ercole-io/ercole/v2/model"
	db "github.com/ercole-io/ercole/v2/thunder-service/database"
)

type AwsDataRetrieveJob struct {
	// Database contains the database layer
	Database db.MongoDatabaseInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Log contains logger formatted
	Log logger.Logger
}

func (job *AwsDataRetrieveJob) Run() {
	if err := job.RetrieveObjectStorageOptimization(); err != nil {
		job.Log.Error(err)
	}

	var profiles []model.AwsProfile

	var seqValue uint64

	var newSeqValue uint64

	awsProfiles, err := job.Database.GetAwsProfiles(false)
	if err != nil {
		job.Log.Error(err)
		return
	}

	for _, val := range awsProfiles {
		if val.Selected {
			profiles = append(profiles, val)
		}
	}

	seqValue, err = job.Database.GetLastAwsSeqValue()

	if err != nil {
		job.Log.Error(err)
		return
	}

	newSeqValue = seqValue + 1
	job.GetAwsUnusedLoadBalancers(profiles, newSeqValue)
	job.GetAwsUnusedIPAddresses(profiles, newSeqValue)
}
