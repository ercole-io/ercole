// Copyright (c) 2025 Sorint.lab S.p.A.
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

// Package service is a package that provides methods for querying data
package service

import (
	"fmt"

	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (hds *HostDataService) createDR(hostdata model.HostDataBE) error {
	drname := fmt.Sprintf("%s-DR", hostdata.Hostname)

	if !hds.Database.ExistsDR(drname) {
		return nil
	}

	if err := hds.Database.DismissHost(drname); err != nil {
		return err
	}

	hostdata.ID = primitive.NewObjectID()
	hostdata.Hostname = drname
	hostdata.IsDR = true

	if hostdata.ClusterMembershipStatus.VeritasClusterServer {
		for _, hostname := range hostdata.ClusterMembershipStatus.VeritasClusterHostnames {
			if err := hds.Database.UpdateVeritasClusterHostnames(hostname, drname); err != nil {
				return err
			}
		}

		hostdata.ClusterMembershipStatus.VeritasClusterHostnames =
			append(hostdata.ClusterMembershipStatus.VeritasClusterHostnames, drname)
	}

	return hds.Database.InsertHostData(hostdata)
}
