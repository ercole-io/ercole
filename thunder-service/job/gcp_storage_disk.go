// Copyright (c) 2024 Sorint.lab S.p.A.
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
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/ercole-io/ercole/v2/model"
)

func (job *GcpDataRetrieveJob) GetDisk(ctx context.Context, projectID, diskname, zone string) (*computepb.Disk, error) {
	c, err := compute.NewDisksRESTClient(ctx, *job.Opt)
	if err != nil {
		return nil, err
	}

	defer c.Close()

	req := &computepb.GetDiskRequest{
		Project: projectID,
		Disk:    diskname,
		Zone:    zone,
	}

	return c.Get(ctx, req)
}

func (job *GcpDataRetrieveJob) FetchGcpStorageDisk(ctx context.Context, gcpDisk model.GcpDisk, seqValue uint64, wg *sync.WaitGroup, ch chan<- model.GcpRecommendation) {
	defer wg.Done()

	maxReadIops, err := job.IsMaxReadIopsStorageOptimizable(ctx, gcpDisk)
	if err != nil {
		job.Log.Error(err)
		return
	}

	maxWriteIops, err := job.IsMaxWriteIopsStorageOptimizable(ctx, gcpDisk)
	if err != nil {
		job.Log.Error(err)
		return
	}

	maxReadThroughput, err := job.IsMaxThroughputReadStorageOptimizable(ctx, gcpDisk)
	if err != nil {
		job.Log.Error(err)
		return
	}

	maxWriteThroughput, err := job.IsMaxThroughputWriteStorageOptimizable(ctx, gcpDisk)
	if err != nil {
		job.Log.Error(err)
		return
	}

	optimizable := maxReadIops.IsOptimizable && maxWriteIops.IsOptimizable &&
		maxReadThroughput.IsOptimizable && maxWriteThroughput.IsOptimizable

	if optimizable {
		sizeGbStr := strconv.Itoa(int(gcpDisk.GetSizeGb()))

		job.Log.Debugf("riops percentage: %.2f", maxReadIops.GetPercentage())
		job.Log.Debugf("wiops percentage: %.2f", maxWriteIops.GetPercentage())
		job.Log.Debugf("rthroughput percentage: %.2f", maxReadThroughput.GetPercentage())
		job.Log.Debugf("wthroughput percentage: %.2f", maxWriteThroughput.GetPercentage())

		optimizationScore := job.GetOptimizationScore(maxReadIops.GetPercentage(), maxWriteIops.GetPercentage(), maxReadThroughput.GetPercentage(), maxWriteThroughput.GetPercentage())

		ch <- model.GcpRecommendation{
			SeqValue:     seqValue,
			CreatedAt:    time.Now(),
			ProfileID:    gcpDisk.ProfileID,
			ResourceID:   gcpDisk.GetId(),
			ResourceName: gcpDisk.GetName(),
			Category:     "Block Storage Rightsizing",
			Suggestion:   "Resize Oversized Disk",
			ProjectID:    gcpDisk.ProjectId,
			ProjectName:  gcpDisk.Project.Name,
			ObjectType:   "Disk",
			Details: map[string]string{
				"Block Storage Name":           gcpDisk.Disk.GetName(),
				"Size GB":                      sizeGbStr,
				"THROUGHPUT W MAX 5DD (MiBps)": fmt.Sprintf("%.2f/%v", maxWriteThroughput.RetrievedValue, maxWriteThroughput.TargetValue),
				"THROUGHPUT R MAX 5DD (MiBps)": fmt.Sprintf("%.2f/%v", maxReadThroughput.RetrievedValue, maxReadThroughput.TargetValue),
				"IOPS W MAX 5DD":               fmt.Sprintf("%.0f/%v", maxWriteIops.RetrievedValue, maxWriteIops.TargetValue),
				"IOPS R MAX 5DD":               fmt.Sprintf("%.0f/%v", maxReadIops.RetrievedValue, maxReadIops.TargetValue),
				"storage type":                 gcpDisk.Type(),
			},
			OptimizationScore: optimizationScore,
		}
	}
}
