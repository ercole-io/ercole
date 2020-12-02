// Copyright (c) 2020 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadManagedTechnologiesList_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
	}
	as.loadManagedTechnologiesList()

	assert.Equal(t, "Oracle/Database", as.TechnologyInfos[0].Product)
	assert.Equal(t, "Microsoft/SQLServer", as.TechnologyInfos[1].Product)
	assert.Equal(t, "iVBORw0K", as.TechnologyInfos[0].Logo[:8])
	assert.Equal(t, "iVBORw0K", as.TechnologyInfos[1].Logo[:8])
}

func TestGetDefaultDatabaseTags_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			APIService: config.APIService{
				DefaultDatabaseTags: []string{
					"foo",
					"bar",
				},
			},
		},
	}
	res, err := as.GetDefaultDatabaseTags()
	require.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar"}, res)
}

func TestGetErcoleFeatures_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := map[string]bool{
		"Oracle/Database": true,
		"Oracle/Exadata":  false,
	}

	getTechnologiesUsageRes := map[string]float64{
		"Oracle/Database": 8,
		"Oracle/Exadata":  0,
	}

	db.EXPECT().
		GetHostsCountUsingTechnologies("", "", utils.MAX_TIME).
		Return(getTechnologiesUsageRes, nil)

	res, err := as.GetErcoleFeatures()

	require.NoError(t, err)
	assert.JSONEq(t, utils.ToJSON(expectedRes), utils.ToJSON(res))
}

func TestGetErcoleFeatures_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().
		GetHostsCountUsingTechnologies("", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	_, err := as.GetErcoleFeatures()

	require.Equal(t, aerrMock, err)
}

func TestGetTechnologyList_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		TechnologyInfos: []model.TechnologyInfo{
			{},
		},
	}
	res, err := as.GetTechnologyList()
	require.NoError(t, err)
	assert.Equal(t, []model.TechnologyInfo{
		{},
	}, res)
}
