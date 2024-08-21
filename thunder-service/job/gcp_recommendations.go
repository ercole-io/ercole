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
	"google.golang.org/api/compute/v1"
)

func (job *GcpDataRetrieveJob) AuditInstancePoint(queryType string, points []*monitoringpb.Point) model.CountValue {
	switch queryType {
	case "avg_cpu":
		counter := 0

		for _, point := range points {
			if point.Value != nil &&
				(point.Value.GetDoubleValue()*100) > float64(job.Config.ThunderService.GcpDataRetrieveJob.AvgCpuPercentage) {
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
			if point.Value != nil &&
				(point.Value.GetDoubleValue()*100) > float64(job.Config.ThunderService.GcpDataRetrieveJob.MaxCpuPercentage) {
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
	}

	return model.CountValue{IsOptimizable: false}
}

func (job *GcpDataRetrieveJob) AuditDiskPoint(queryType string, disk model.GcpDisk, points []*monitoringpb.Point) model.OptimizableValue {
	switch queryType {
	case "max_read_iops":
		var maxMeasurement float64

		rIops := disk.ReadIopsPerGib()
		limit := rIops * (float64(job.Config.ThunderService.GcpDataRetrieveJob.IopsStoragePercentage) / 100)

		for _, point := range points {
			if point.Value == nil {
				continue
			}

			if float64(point.Value.GetInt64Value()) > maxMeasurement {
				maxMeasurement = float64(point.Value.GetInt64Value())
			}

			job.Log.Debugf(
				`--Received--
				disk name: %s
				disk type: %s
				vpcus: %d
				shared core: %t
				read iops: %v
				read iops limit: %v
				condition limit: %v
				----`,
				disk.GetName(),
				disk.Type(),
				disk.InstanceVcpus,
				disk.IsSharedCore,
				float64(point.Value.GetInt64Value()),
				rIops,
				limit,
			)
		}

		if maxMeasurement < limit {
			return model.OptimizableValue{
				IsOptimizable:  true,
				RetrievedValue: maxMeasurement,
				TargetValue:    rIops,
			}
		}

		return model.OptimizableValue{
			IsOptimizable:  false,
			RetrievedValue: maxMeasurement,
			TargetValue:    rIops,
		}

	case "max_write_iops":
		var maxMeasurement float64

		wIops := disk.WriteIopsPerGib()
		limit := wIops * (float64(job.Config.ThunderService.GcpDataRetrieveJob.IopsStoragePercentage) / 100)

		for _, point := range points {
			if point.Value == nil {
				continue
			}

			if float64(point.Value.GetInt64Value()) > maxMeasurement {
				maxMeasurement = float64(point.Value.GetInt64Value())
			}

			job.Log.Debugf(
				`--Received--
				disk name: %s
				disk type: %s
				vpcus: %d
				shared core: %t
				write iops: %v
				write iops limit: %v
				condition limit: %v
				----`,
				disk.GetName(),
				disk.Type(),
				disk.InstanceVcpus,
				disk.IsSharedCore,
				float64(point.Value.GetInt64Value()),
				wIops,
				limit,
			)
		}

		if maxMeasurement < limit {
			return model.OptimizableValue{
				IsOptimizable:  true,
				RetrievedValue: maxMeasurement,
				TargetValue:    wIops,
			}
		}

		return model.OptimizableValue{
			IsOptimizable:  false,
			RetrievedValue: maxMeasurement,
			TargetValue:    wIops,
		}

	case "max_read_throughput":
		var maxMeasurement, pointValue float64

		rThroughput := disk.ReadThroughputPerMib()
		limit := rThroughput * (float64(job.Config.ThunderService.GcpDataRetrieveJob.ThroughputStoragePercentage) / 100)

		for _, point := range points {
			if point.Value == nil {
				continue
			}

			pointValue = float64(point.Value.GetInt64Value()) / 1048576

			if pointValue > maxMeasurement {
				maxMeasurement = pointValue
			}

			job.Log.Debugf(
				`--Received--
				disk name: %s
				disk type: %s
				vpcus: %d
				shared core: %t
				read throughput: %v
				read throughput limit: %v
				condition limit: %v
				----`,
				disk.GetName(),
				disk.Type(),
				disk.InstanceVcpus,
				disk.IsSharedCore,
				pointValue,
				rThroughput,
				limit,
			)
		}

		if maxMeasurement < limit {
			return model.OptimizableValue{
				IsOptimizable:  true,
				RetrievedValue: maxMeasurement,
				TargetValue:    rThroughput,
			}
		}

		return model.OptimizableValue{
			IsOptimizable:  false,
			RetrievedValue: maxMeasurement,
			TargetValue:    rThroughput,
		}

	case "max_write_throughput":
		var maxMeasurement, pointValue float64

		wThroughput := disk.WriteThroughputPerMib()
		limit := wThroughput * (float64(job.Config.ThunderService.GcpDataRetrieveJob.ThroughputStoragePercentage) / 100)

		for _, point := range points {
			if point.Value == nil {
				continue
			}

			pointValue = float64(point.Value.GetInt64Value()) / 1048576

			if pointValue > maxMeasurement {
				maxMeasurement = pointValue
			}

			job.Log.Debugf(
				`--Received--
				disk name: %s
				disk type: %s
				vpcus: %d
				shared core: %t
				write throughput: %v
				write throughput limit: %v
				condition limit: %v
				----`,
				disk.GetName(),
				disk.Type(),
				disk.InstanceVcpus,
				disk.IsSharedCore,
				pointValue,
				wThroughput,
				limit,
			)
		}

		if maxMeasurement < limit {
			return model.OptimizableValue{
				IsOptimizable:  true,
				RetrievedValue: maxMeasurement,
				TargetValue:    wThroughput,
			}
		}

		return model.OptimizableValue{
			IsOptimizable:  false,
			RetrievedValue: maxMeasurement,
			TargetValue:    wThroughput,
		}
	}

	return model.OptimizableValue{}
}

func (job *GcpDataRetrieveJob) AuditMaxMemPoint(machineType *compute.MachineType, points []*monitoringpb.Point) model.CountValue {
	counter := 0

	for _, point := range points {
		if point.Value != nil {
			maxMem := point.Value.GetInt64Value()
			maxMemMib := float64(maxMem) / 1048576
			percentage := (maxMemMib / float64(machineType.MemoryMb)) * 100

			if percentage > float64(job.Config.ThunderService.GcpDataRetrieveJob.MaxMemPercentage) {
				counter++
			}
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

func (job *GcpDataRetrieveJob) AddRecommendation(recommendation model.GcpRecommendation) error {
	return job.Database.AddGcpRecommendation(recommendation)
}

func (job *GcpDataRetrieveJob) AddError(gcperror model.GcpError) error {
	return job.Database.AddGcpError(gcperror)
}
