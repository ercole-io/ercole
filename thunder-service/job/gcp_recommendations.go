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
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/ercole-io/ercole/v2/model"
)

func (job *GcpDataRetrieveJob) AuditInstancePoint(queryType string, points []*monitoringpb.Point) model.CountValue {
	switch queryType {
	case "avg_cpu":
		counter := 0

		for _, point := range points {
			if point.Value != nil && point.Value.GetDoubleValue() > 0.5 {
				counter++
			}

			if counter >= int(job.Config.ThunderService.GcpDataRetrieveJob.AvgCpuUtilizationThreshold) {
				return model.CountValue{IsOptimizable: false}
			}
		}

		return model.CountValue{
			IsOptimizable: true,
			Count:         counter,
		}

	case "max_cpu":
		counter := 0

		for _, point := range points {
			if point.Value != nil && point.Value.GetDoubleValue() > 0.5 {
				counter++
			}

			if counter >= int(job.Config.ThunderService.GcpDataRetrieveJob.MaxCpuUtilizationThreshold) {
				return model.CountValue{IsOptimizable: false}
			}
		}

		return model.CountValue{
			IsOptimizable: true,
			Count:         counter,
		}

	case "max_mem":
		counter := 0

		for _, point := range points {
			if point.Value != nil && point.Value.GetDoubleValue() > 0.9 {
				counter++
			}

			if counter >= int(job.Config.ThunderService.GcpDataRetrieveJob.MaxMemUtilizationThreshold) {
				return model.CountValue{IsOptimizable: false}
			}
		}

		return model.CountValue{
			IsOptimizable: true,
			Count:         counter,
		}
	}

	return model.CountValue{IsOptimizable: false}
}

func (job *GcpDataRetrieveJob) AuditDiskPoint(queryType string, disk model.GcpDisk, points []*monitoringpb.Point) model.OptimizableValue {
	switch queryType {
	case "max_read_iops":
		var maxMeasurement float64

		for _, point := range points {
			if point.Value != nil && float64(point.Value.GetInt64Value()) > maxMeasurement {
				maxMeasurement = float64(point.Value.GetInt64Value())
			}

			if point.Value != nil && float64(point.Value.GetInt64Value()) < disk.ReadIopsPerGib()/2 {
				return model.OptimizableValue{
					IsOptimizable:  true,
					RetrievedValue: float64(point.Value.GetInt64Value()),
					TargetValue:    disk.ReadIopsPerGib(),
				}
			}
		}

		return model.OptimizableValue{
			IsOptimizable:  false,
			RetrievedValue: maxMeasurement,
			TargetValue:    disk.ReadIopsPerGib(),
		}

	case "max_write_iops":
		var maxMeasurement float64

		for _, point := range points {
			if point.Value != nil && float64(point.Value.GetInt64Value()) > maxMeasurement {
				maxMeasurement = float64(point.Value.GetInt64Value())
			}

			if point.Value != nil && float64(point.Value.GetInt64Value()) < disk.WriteIopsPerGib()/2 {
				return model.OptimizableValue{
					IsOptimizable:  true,
					RetrievedValue: float64(point.Value.GetInt64Value()),
					TargetValue:    disk.WriteIopsPerGib(),
				}
			}
		}

		return model.OptimizableValue{
			IsOptimizable:  false,
			RetrievedValue: maxMeasurement,
			TargetValue:    disk.WriteIopsPerGib(),
		}

	case "max_read_throughput":
		var maxMeasurement, pointValue float64

		for _, point := range points {
			if point.Value != nil {
				pointValue = float64(point.Value.GetInt64Value()) / 1048576

				if pointValue > maxMeasurement {
					maxMeasurement = pointValue
				}
			}

			if point.Value != nil && pointValue < disk.ReadThroughputPerMib()/2 {
				return model.OptimizableValue{
					IsOptimizable:  true,
					RetrievedValue: pointValue,
					TargetValue:    disk.ReadThroughputPerMib(),
				}
			}
		}

		return model.OptimizableValue{
			IsOptimizable:  false,
			RetrievedValue: maxMeasurement,
			TargetValue:    disk.ReadThroughputPerMib(),
		}

	case "max_write_throughput":
		var maxMeasurement, pointValue float64

		for _, point := range points {
			if point.Value != nil {
				pointValue = float64(point.Value.GetInt64Value()) / 1048576

				if pointValue > maxMeasurement {
					maxMeasurement = pointValue
				}
			}

			if point.Value != nil && pointValue < disk.WriteThroughputPerMib()/2 {
				return model.OptimizableValue{
					IsOptimizable:  true,
					RetrievedValue: pointValue,
					TargetValue:    disk.WriteThroughputPerMib(),
				}
			}
		}

		return model.OptimizableValue{
			IsOptimizable:  false,
			RetrievedValue: maxMeasurement,
			TargetValue:    disk.WriteThroughputPerMib(),
		}
	}

	return model.OptimizableValue{}
}

func (job *GcpDataRetrieveJob) AddRecommendation(recommendation model.GcpRecommendation) error {
	return job.Database.AddGcpRecommendation(recommendation)
}

func (job *GcpDataRetrieveJob) AddError(gcperror model.GcpError) error {
	return job.Database.AddGcpError(gcperror)
}
