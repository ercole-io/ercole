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

	"github.com/ercole-io/ercole/v2/api-service/domain"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
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

	db.EXPECT().ListExadataInstances(f, false).Return(expected, nil)

	res, err := as.ListExadataInstances(f, false)
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

	modelInstance := mongoutils.LoadFixtureExadata(t, "../../fixture/test_apiservice_exadata_01.json")
	println(modelInstance.Hostname)

	db.EXPECT().FindExadataInstance("3M2ORPFI9Q", false).Return(&modelInstance, nil)

	dom, err := domain.ToOracleExadataInstance(modelInstance)
	require.NoError(t, err)

	res, err := as.GetExadataInstance("3M2ORPFI9Q", false)
	require.NoError(t, err)
	assert.Equal(t, dom, res)
}

func TestHideExadataInstance_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	exadata := mongoutils.LoadFixtureExadata(t, "../../fixture/test_apiservice_exadata_01.json")

	db.EXPECT().FindExadataInstance("3M2ORPFI9Q", false).Return(&exadata, nil)

	exadata.Hidden = true

	db.EXPECT().UpdateExadataInstance(exadata).Return(nil).AnyTimes()

	err := as.HideExadataInstance("3M2ORPFI9Q")

	require.NoError(t, err)

	assert.Equal(t, true, exadata.Hidden)
}

func TestShowExadataInstance(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	exadata := mongoutils.LoadFixtureExadata(t, "../../fixture/test_apiservice_exadata_02.json")

	db.EXPECT().FindExadataInstance("3M2ORPFI9W", true).Return(&exadata, nil)

	exadata.Hidden = false

	db.EXPECT().UpdateExadataInstance(exadata).Return(nil).AnyTimes()

	err := as.ShowExadataInstance("3M2ORPFI9W")

	require.NoError(t, err)

	assert.Equal(t, false, exadata.Hidden)
}

func TestGetExadataPatchAdvisors(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := []dto.OracleExadataPatchAdvisor{
		{
			Hostname:     "exa01db01",
			RackID:       "3M2ORPFI9W",
			ImageVersion: "22.1.19.0.0.240119.1",
			ReleaseDate:  utils.P("2024-01-19T00:00:00Z"),
			FourMonths:   false,
			SixMonths:    false,
			TwelveMonths: false,
		},
	}

	db.EXPECT().FindAllExadataPatchAdvisors().Return(expected, nil)

	actual, err := as.GetExadataPatchAdvisors()

	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}
