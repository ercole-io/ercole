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
	"regexp"
	"strconv"
	"strings"

	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/cloudresourcemanager/v1"
)

type GcpInstance struct {
	*computepb.Instance
	*cloudresourcemanager.Project
	ProfileID primitive.ObjectID
}

func (i GcpInstance) Zone() string {
	parts := strings.Split(i.GetZone(), "/")

	if len(parts) == 0 {
		return ""
	}

	return parts[len(parts)-1]
}

func (i GcpInstance) MachineType() string {
	parts := strings.Split(i.GetMachineType(), "/")

	if len(parts) == 0 {
		return ""
	}

	return parts[len(parts)-1]
}

func (i GcpInstance) VCpu() int {
	machineType := i.MachineType()
	re := regexp.MustCompile(`\d+$`)
	match := re.FindString(machineType)

	if match != "" {
		vcpu, err := strconv.Atoi(match)
		if err != nil {
			return 0
		}

		return vcpu
	}

	return 2
}

func (i GcpInstance) IsSharedCore() bool {
	sharedMachineType := []string{"e2-micro", "e2-small", "e2-medium"}
	return utils.Contains(sharedMachineType, i.MachineType())
}

type CountValue struct {
	IsOptimizable bool
	Count         int
}
