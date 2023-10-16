// Copyright (c) 2023 Sorint.lab S.p.A.
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

package service

import (
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
)

func (as *APIService) ListExadataInstances(filter dto.GlobalFilter) ([]model.OracleExadataInstance, error) {
	return as.Database.ListExadataInstances(filter)
}

func (as *APIService) UpdateExadataVmClusterName(rackID, hostID, vmname, clustername string) error {
	instance, err := as.Database.GetExadataInstance(rackID)
	if err != nil {
		return err
	}

	for i := range instance.Components {
		if instance.Components[i].HostID == hostID {
			for j := range instance.Components[i].VMs {
				if instance.Components[i].VMs[j].Name == vmname {
					instance.Components[i].VMs[j].ClusterName = clustername
				}
			}
		}
	}

	return as.Database.UpdateExadataInstance(*instance)
}

func (as *APIService) UpdateExadataComponentClusterName(RackID, hostID, clustername string) error {
	instance, err := as.Database.GetExadataInstance(RackID)
	if err != nil {
		return err
	}

	for i := range instance.Components {
		if instance.Components[i].HostID == hostID {
			instance.Components[i].ClusterName = clustername
		}
	}

	return as.Database.UpdateExadataInstance(*instance)
}
