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
	"strconv"
	"strings"

	"cloud.google.com/go/compute/apiv1/computepb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/cloudresourcemanager/v1"
)

type GcpDisk struct {
	InstanceID    uint64
	InstanceZone  string
	MachineType   string
	InstanceVcpus int
	IsSharedCore  bool
	ProfileID     primitive.ObjectID
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
	Limit map[string]RW
}

type ThroughputPerMib struct {
	Limit map[string]RW
}

const (
	DiskTypeStandard = "pd-standard"
	DiskTypeBalanced = "pd-balanced"
	DiskTypeSsd      = "pd-ssd"
)

var (
	IopsLimits = map[string]IopsPerGib{
		DiskTypeStandard: {
			Limit: map[string]RW{
				"shared-core": {Read: 1000, Write: 10_000},
				"2-7":         {Read: 3000, Write: 15_000},
				"8-15":        {Read: 5000, Write: 15_000},
				"16+":         {Read: 7500, Write: 15_000},
			},
		},

		DiskTypeBalanced: {
			Limit: map[string]RW{
				"shared-core": {Write: 10_000, Read: 12_000},
				"2-7":         {Write: 15_000, Read: 15_000},
				"8-15":        {Write: 15_000, Read: 15_000},
				"16-31":       {Write: 20_000, Read: 20_000},
				"32+":         {Write: 50_000, Read: 50_000},
			},
		},
	}

	ThroughputLimits = map[string]ThroughputPerMib{
		DiskTypeStandard: {
			Limit: map[string]RW{
				"shared-core": {Read: 200, Write: 200},
				"2-7":         {Read: 240, Write: 240},
				"8-15":        {Read: 800, Write: 400},
				"16+":         {Read: 1200, Write: 400},
			},
		},

		DiskTypeBalanced: {
			Limit: map[string]RW{
				"shared-core": {Write: 200, Read: 200},
				"2-7":         {Write: 240, Read: 240},
				"8-15":        {Write: 800, Read: 800},
				"16-31":       {Write: 1000, Read: 1200},
				"32+":         {Write: 1000, Read: 1200},
			},
		},
	}
)

func (d GcpDisk) ReadIopsPerGib() float64 {
	return IopsLimits[d.Type()].Limit[mustReturnRange(d.InstanceVcpus, d.IsSharedCore)].Read
}

func (d GcpDisk) WriteIopsPerGib() float64 {
	return IopsLimits[d.Type()].Limit[mustReturnRange(d.InstanceVcpus, d.IsSharedCore)].Write
}

func (d GcpDisk) ReadThroughputPerMib() float64 {
	return ThroughputLimits[d.Type()].Limit[mustReturnRange(d.InstanceVcpus, d.IsSharedCore)].Read
}

func (d GcpDisk) WriteThroughputPerMib() float64 {
	return ThroughputLimits[d.Type()].Limit[mustReturnRange(d.InstanceVcpus, d.IsSharedCore)].Write
}

func mustReturnRange(vcpu int, isSharedCore bool) string {
	if isSharedCore {
		return "shared-core"
	}

	cpuRanges := []string{
		"shared-core",
		"2-7",
		"8-15",
		"16+",
	}

	for _, cpuRange := range cpuRanges {
		if strings.Contains(cpuRange, "-") {
			bounds := strings.Split(cpuRange, "-")
			lowerBound, err1 := strconv.Atoi(bounds[0])
			upperBound, err2 := strconv.Atoi(bounds[1])

			if err1 != nil && err2 != nil {
				continue
			}

			if vcpu >= lowerBound && vcpu <= upperBound {
				return cpuRange
			}
		} else if strings.HasSuffix(cpuRange, "+") {
			lowerBound, err := strconv.Atoi(strings.TrimSuffix(cpuRange, "+"))
			if err != nil {
				continue
			}

			if vcpu >= lowerBound {
				return cpuRange
			}
		}
	}

	return ""
}

type OptimizableValue struct {
	IsOptimizable  bool
	RetrievedValue float64
	TargetValue    float64
}
