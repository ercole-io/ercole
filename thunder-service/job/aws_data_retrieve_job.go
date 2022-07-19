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
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	db "github.com/ercole-io/ercole/v2/thunder-service/database"
)

type AwsDataRetrieveJob struct {
	Database db.MongoDatabaseInterface
	Config   config.Configuration
	Log      logger.Logger
}

func (job *AwsDataRetrieveJob) Run() {
	awsProfiles, err := job.Database.GetAwsProfiles(false)
	if err != nil {
		job.Log.Error(err)
		return
	}

	seqValue, err := job.Database.GetLastAwsSeqValue()
	if err != nil {
		job.Log.Error(err)
		return
	}

	seqValue = seqValue + 1

	c := make(chan error)

	for _, profile := range awsProfiles {
		go func(profile model.AwsProfile, seq uint64) {
			if err := job.FetchObjectStorageOptimization(profile, seq); err != nil {
				c <- err
			}
		}(profile, seqValue)

		go func(profile model.AwsProfile, seq uint64) {
			if err := job.FetchAwsUnusedLoadBalancers(profile, seq); err != nil {
				c <- err
			}
		}(profile, seqValue)

		go func(profile model.AwsProfile, seq uint64) {
			if err := job.FetchAwsUnusedIPAddresses(profile, seq); err != nil {
				c <- err
			}
		}(profile, seqValue)

		go func(profile model.AwsProfile, seq uint64) {
			if err := job.FetchAwsVolumesNotUsed(profile, seq); err != nil {
				c <- err
			}
		}(profile, seqValue)

		go func(profile model.AwsProfile, seq uint64) {
			if err := job.FetchAwsNotActiveInstances(profile, seq); err != nil {
				c <- err
			}
		}(profile, seqValue)
	}

	job.Log.Error(<-c)
}
