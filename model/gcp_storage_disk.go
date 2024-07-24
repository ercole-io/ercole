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
package model

import (
	"strings"

	"cloud.google.com/go/compute/apiv1/computepb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/cloudresourcemanager/v1"
)

type GcpDisk struct {
	InstanceID   uint64
	InstanceZone string
	MachineType  string
	ProfileID    primitive.ObjectID
	*cloudresourcemanager.Project
	*computepb.Disk
}

func (d GcpDisk) Type() string {
	parts := strings.Split(d.GetType(), "/")

	if len(parts) == 0 {
		return ""
	}

	return parts[len(parts)-1]
}

type RW struct {
	Read  float64
	Write float64
}

type IopsPerGib struct {
	RW
	Limit RW
}

type ThroughputPerGib struct {
	RW
	Limit RW
}

const (
	DiskTypeStandard = "pd-standard"
	DiskTypeBalanced = "pd-balanced"
	DiskTypeSsd      = "pd-ssd"
)

var (
	IopsLimits = map[string]IopsPerGib{
		DiskTypeStandard: {RW: RW{Read: 0.75, Write: 1.5}, Limit: RW{Read: 7500, Write: 15000}},
		DiskTypeBalanced: {RW: RW{Read: 6, Write: 6}, Limit: RW{Read: 80000, Write: 80000}},
		DiskTypeSsd:      {RW: RW{Read: 30, Write: 30}, Limit: RW{Read: 100000, Write: 100000}},
	}

	ThroughputLimits = map[string]ThroughputPerGib{
		DiskTypeStandard: {RW: RW{Read: 0.12, Write: 0.12}, Limit: RW{Read: 1200, Write: 400}},
		DiskTypeBalanced: {RW: RW{Read: 0.28, Write: 0.28}, Limit: RW{Read: 1200, Write: 1200}},
		DiskTypeSsd:      {RW: RW{Read: 0.48, Write: 0.48}, Limit: RW{Read: 1200, Write: 1200}},
	}
)

func (d GcpDisk) ReadIopsPerGib() float64 {
	limit := float64(d.GetSizeGb() * int64(IopsLimits[d.Type()].Read))

	if limit > IopsLimits[d.Type()].Limit.Read {
		return IopsLimits[d.Type()].Limit.Read
	}

	return limit
}

func (d GcpDisk) WriteIopsPerGib() float64 {
	limit := float64(d.GetSizeGb() * int64(IopsLimits[d.Type()].Write))

	if limit > IopsLimits[d.Type()].Limit.Write {
		return IopsLimits[d.Type()].Limit.Write
	}

	return limit
}

func (d GcpDisk) ReadThroughputPerGib() float64 {
	limit := float64(d.GetSizeGb() * int64(ThroughputLimits[d.Type()].Read))

	if limit > ThroughputLimits[d.Type()].Limit.Read {
		return ThroughputLimits[d.Type()].Limit.Read
	}

	return limit
}

func (d GcpDisk) WriteThroughputPerGib() float64 {
	limit := float64(d.GetSizeGb() * int64(ThroughputLimits[d.Type()].Write))

	if limit > IopsLimits[d.Type()].Limit.Write {
		return IopsLimits[d.Type()].Limit.Write
	}

	return limit
}

type OptimizableValue struct {
	IsOptimizable  bool
	RetrievedValue int64
	TargetValue    float64
}
