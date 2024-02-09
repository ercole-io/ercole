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
	"testing"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListExadataInstances_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := []dto.ExadataInstanceResponse{
		{
			Hostname: "exadataTest01",
			RackID:   "RACTEST01",
		},
	}

	f := dto.GlobalFilter{OlderThan: utils.MAX_TIME}

	db.EXPECT().ListExadataInstances(f).Return(expected, nil)

	res, err := as.ListExadataInstances(f)
	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestUpdateExadataVmClusterName_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	original := model.OracleExadataVmClustername{
		InstanceRackID: "rackID_test",
		HostID:         "hostID_test",
		VmName:         "vm_test",
		Clustername:    "original_clustername",
	}

	db.EXPECT().FindExadataVmClustername("rackID_test", "hostID_test", "vm_test").Return(&original, nil).AnyTimes()
	db.EXPECT().UpdateExadataVmClustername("rackID_test", "hostID_test", "vm_test", "updated_cluster_name").Return(nil).AnyTimes()

	err := as.UpdateExadataVmClusterName("rackID_test", "hostID_test", "vm_test", "updated_cluster_name")
	require.NoError(t, err)
}

func TestGetExadataInstance_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := &model.OracleExadataInstance{
		RackID: "rackid_test",
	}

	db.EXPECT().FindExadataInstance("rackid_test").Return(expected, nil)

	res, err := as.GetExadataInstance("rackid_test")
	require.NoError(t, err)
	assert.Equal(t, expected, res)
}
