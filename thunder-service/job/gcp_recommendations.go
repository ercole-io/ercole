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

func (job *GcpDataRetrieveJob) AuditInstancePoint(queryType string, points []*monitoringpb.Point) bool {
	switch queryType {
	case "avg_cpu":
		counter := 0

		for _, point := range points {
			if point.Value != nil && point.Value.GetDoubleValue() > 0.5 {
				counter++
			}

			if counter >= int(job.Config.ThunderService.GcpDataRetrieveJob.AvgCpuUtilizationThreshold) {
				return false
			}
		}

		return true

	case "max_cpu":
		counter := 0

		for _, point := range points {
			if point.Value != nil && point.Value.GetDoubleValue() > 0.5 {
				counter++
			}

			if counter >= int(job.Config.ThunderService.GcpDataRetrieveJob.MaxCpuUtilizationThreshold) {
				return false
			}
		}

		return true

	case "max_mem":
		counter := 0

		for _, point := range points {
			if point.Value != nil && point.Value.GetDoubleValue() > 0.9 {
				counter++
			}

			if counter >= int(job.Config.ThunderService.GcpDataRetrieveJob.MaxMemUtilizationThreshold) {
				return false
			}
		}

		return true
	}

	return false
}

func (job *GcpDataRetrieveJob) AuditDiskPoint(queryType string, disk model.GcpDisk, points []*monitoringpb.Point) bool {
	switch queryType {
	case "max_read_iops":
		for _, point := range points {
			if point.Value != nil && point.Value.GetInt64Value() < int64(disk.ReadIopsPerGib()/2) {
				return true
			}
		}

	case "max_write_iops":
		for _, point := range points {
			if point.Value != nil && point.Value.GetInt64Value() < int64(disk.WriteIopsPerGib()/2) {
				return true
			}
		}

	case "max_read_throughput":
		for _, point := range points {
			if point.Value != nil && point.Value.GetInt64Value() < int64(disk.ReadThroughputPerGib()/2) {
				return true
			}
		}

	case "max_write_throughput":
		for _, point := range points {
			if point.Value != nil && point.Value.GetInt64Value() < int64(disk.WriteThroughputPerGib()/2) {
				return true
			}
		}
	}

	return false
}

func (job *GcpDataRetrieveJob) AddRecommendation(recommendation model.GcpRecommendation) error {
	return job.Database.AddGcpRecommendation(recommendation)
}

func (job *GcpDataRetrieveJob) AddError(gcperror model.GcpError) error {
	return job.Database.AddGcpError(gcperror)
}
