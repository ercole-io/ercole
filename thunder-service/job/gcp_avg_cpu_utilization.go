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
	"time"

	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/ercole-io/ercole/v2/model"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (job *GcpDataRetrieveJob) IsAvgCpuUtilizationOptimizable(ctx context.Context, instance model.GcpInstance) (*model.CountValue, error) {
	now := time.Now()
	endMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startdate := now.AddDate(0, 0, -90)
	startdateMidnight := time.Date(startdate.Year(), startdate.Month(), startdate.Day(), 0, 0, 0, 0, startdate.Location())

	filter := fmt.Sprintf(`metric.type = "compute.googleapis.com/instance/cpu/utilization"
	AND resource.labels.instance_id = "%d"
	AND resource.labels.zone = "%s"`, instance.GetId(), instance.Zone())

	req := &monitoringpb.ListTimeSeriesRequest{
		Name:   fmt.Sprintf("projects/%s", instance.ProjectId),
		Filter: filter,
		Interval: &monitoringpb.TimeInterval{
			StartTime: timestamppb.New(startdateMidnight),
			EndTime:   timestamppb.New(endMidnight),
		},
		Aggregation: &monitoringpb.Aggregation{
			AlignmentPeriod:  &durationpb.Duration{Seconds: 86400},
			PerSeriesAligner: monitoringpb.Aggregation_ALIGN_MEAN,
		},
	}

	series, err := job.GetTimeSeries(ctx, *job.Opt, req)
	if err != nil {
		return nil, err
	}

	if series != nil && series.Points != nil {
		countVal := job.AuditInstancePoint("avg_cpu", series.Points)

		return &countVal, nil
	}

	return &model.CountValue{IsOptimizable: false}, nil
}
