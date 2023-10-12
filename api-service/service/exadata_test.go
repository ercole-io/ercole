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
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListExadataInstances_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := []model.OracleExadataInstance{
		{
			Hostname: "exadataTest01",
			RackID:   "RACTEST01",
			Components: []model.OracleExadataComponent{
				{
					RackID:   "RACCMPTEST01",
					HostType: "DOM0",
				},
			},
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

	original := model.OracleExadataInstance{
		RackID: "rackID_test",
		Components: []model.OracleExadataComponent{
			{
				HostID: "hostID_test",
				VMs: []model.OracleExadataVM{
					{
						Name:        "vm_test",
						ClusterName: "",
					},
				},
			},
		},
	}

	db.EXPECT().GetExadataInstance("rackID_test").Return(&original, nil).AnyTimes()
	db.EXPECT().UpdateExadataInstance(original).Return(nil).AnyTimes()

	err := as.UpdateExadataVmClusterName("rackID_test", "hostID_test", "vm_test", "updated_cluster_name")
	require.NoError(t, err)
}
