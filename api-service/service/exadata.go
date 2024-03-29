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
	"github.com/ercole-io/ercole/v2/api-service/domain"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/mongo"
)

func (as *APIService) ListExadataInstances(filter dto.GlobalFilter, hidden bool) ([]dto.ExadataInstanceResponse, error) {
	return as.Database.ListExadataInstances(filter, hidden)
}

func (as *APIService) GetExadataInstance(rackid string, hidden bool) (*domain.OracleExadataInstance, error) {
	instance, err := as.Database.FindExadataInstance(rackid, hidden)
	if err != nil {
		return nil, err
	}

	dom, err := domain.ToOracleExadataInstance(*instance)
	if err != nil {
		return nil, err
	}

	return dom, nil
}

func (as *APIService) UpdateExadataVmClusterName(rackID, hostID, vmname, clustername string) error {
	if _, err := as.Database.FindExadataVmClustername(rackID, hostID, vmname); err != nil {
		if err == mongo.ErrNoDocuments {
			return as.Database.InsertExadataVmClustername(rackID, hostID, vmname, clustername)
		}

		return err
	}

	return as.Database.UpdateExadataVmClustername(rackID, hostID, vmname, clustername)
}

func (as *APIService) UpdateExadataComponentClusterName(RackID, hostID string, clusternames []string) error {
	instance, err := as.Database.FindExadataInstance(RackID, false)
	if err != nil {
		return err
	}

	for i := range instance.Components {
		if instance.Components[i].HostID == hostID {
			instance.Components[i].ClusterNames = clusternames
		}
	}

	return as.Database.UpdateExadataInstance(*instance)
}

func (as *APIService) UpdateExadataRdma(rackID string, rdma model.OracleExadataRdma) error {
	instance, err := as.Database.FindExadataInstance(rackID, false)
	if err != nil {
		return err
	}

	instance.RDMA = &rdma

	return as.Database.UpdateExadataInstance(*instance)
}

func (as *APIService) HideExadataInstance(rackID string) error {
	instance, err := as.Database.FindExadataInstance(rackID, false)
	if err != nil {
		return err
	}

	instance.Hidden = true

	return as.Database.UpdateExadataInstance(*instance)
}

func (as *APIService) ShowExadataInstance(rackID string) error {
	instance, err := as.Database.FindExadataInstance(rackID, true)
	if err != nil {
		return err
	}

	instance.Hidden = false

	return as.Database.UpdateExadataInstance(*instance)
}
