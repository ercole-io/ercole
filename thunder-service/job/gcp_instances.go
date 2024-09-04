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
	"sync"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/ercole-io/ercole/v2/model"
	"google.golang.org/api/iterator"
)

func (job *GcpDataRetrieveJob) GetInstances(ctx context.Context, projectID string) ([]*computepb.Instance, error) {
	c, err := compute.NewInstancesRESTClient(ctx, *job.Opt)
	if err != nil {
		return nil, err
	}

	defer c.Close()

	req := &computepb.AggregatedListInstancesRequest{
		Project: projectID,
	}

	it := c.AggregatedList(ctx, req)

	res := make([]*computepb.Instance, 0)

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		if resp.Value.Instances != nil {
			res = append(res, resp.Value.Instances...)
		}
	}

	return res, nil
}

func (job *GcpDataRetrieveJob) FetchGcpInstanceRightsizing(ctx context.Context, gcpInstance model.GcpInstance, seqValue uint64, wg *sync.WaitGroup, ch chan<- model.GcpRecommendation) {
	defer wg.Done()

	avgcpumetrics, err := job.IsAvgCpuUtilizationOptimizable(ctx, gcpInstance)
	if err != nil {
		job.Log.Error(err)
		return
	}

	maxcpumetrics, err := job.IsMaxCpuUtilizationOptimizable(ctx, gcpInstance)
	if err != nil {
		job.Log.Error(err)
		return
	}

	maxmemmetrics, err := job.IsMaxMemUtilizationOptimizable(ctx, gcpInstance)
	if err != nil {
		job.Log.Error(err)
		return
	}

	optimizable := avgcpumetrics.IsOptimizable && maxcpumetrics.IsOptimizable && maxmemmetrics.IsOptimizable

	if optimizable {
		job.Log.Debugf("avgcpumetrics percentage: %.2f", avgcpumetrics.GetPercentage())
		job.Log.Debugf("maxcpumetrics percentage: %.2f", maxcpumetrics.GetPercentage())
		job.Log.Debugf("maxmemmetrics percentage: %.2f", maxmemmetrics.GetPercentage())

		optimizationScore := job.GetOptimizationScore(avgcpumetrics.GetPercentage(), maxcpumetrics.GetPercentage(), maxmemmetrics.GetPercentage())

		ch <- model.GcpRecommendation{
			SeqValue:     seqValue,
			CreatedAt:    time.Now(),
			ProfileID:    gcpInstance.ProfileID,
			ResourceID:   gcpInstance.GetId(),
			ResourceName: gcpInstance.GetName(),
			Category:     "Compute Instance Rightsizing",
			Suggestion:   "Resize oversized compute instance",
			ProjectID:    gcpInstance.ProjectId,
			ProjectName:  gcpInstance.Project.Name,
			ObjectType:   "Compute Instance",
			Details: map[string]string{
				"Instance Name": gcpInstance.GetName(),
				"Cpu Average": fmt.Sprintf("%%Cpu Average 90dd - Number of Threshold Reached (>%d%%): %d/%d",
					job.Config.ThunderService.GcpDataRetrieveJob.AvgCpuPercentage,
					avgcpumetrics.Count,
					avgcpumetrics.TargetValue),
				"Cpu Max": fmt.Sprintf("%%Cpu Max 7dd - Number of Threshold Reached (>%d%%): %d/%d",
					job.Config.ThunderService.GcpDataRetrieveJob.MaxCpuPercentage,
					maxcpumetrics.Count,
					maxcpumetrics.TargetValue),
				"Mem Max": fmt.Sprintf("%%Memory Average 7dd - Number of Threshold Reached (>%d%%): %d/%d",
					job.Config.ThunderService.GcpDataRetrieveJob.MaxMemPercentage,
					maxmemmetrics.Count,
					maxmemmetrics.TargetValue),
			},
			OptimizationScore: optimizationScore,
		}
	}
}
