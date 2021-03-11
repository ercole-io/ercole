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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchMySQLInstances(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		filter := dto.GlobalFilter{
			Location:    "Greece",
			Environment: "TEST",
			OlderThan:   utils.P("2020-05-20T09:53:34+00:00"),
		}

		expected := []dto.MySQLInstance{
			{
				Hostname:      "pippo",
				Location:      "",
				Environment:   "",
				MySQLInstance: model.MySQLInstance{},
			},
		}

		db.EXPECT().SearchMySQLInstances(filter).
			Return(expected, nil).Times(1)

		actual, err := as.SearchMySQLInstances(filter)
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		filter := dto.GlobalFilter{
			Location:    "Greece",
			Environment: "TEST",
			OlderThan:   utils.P("2020-05-20T09:53:34+00:00"),
		}

		db.EXPECT().SearchMySQLInstances(filter).
			Return(nil, errMock).Times(1)

		actual, err := as.SearchMySQLInstances(filter)
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestSearchMySQLInstancesAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
	}

	returned := []dto.MySQLInstance{
		{
			Hostname:    "pippo",
			Location:    "Italy",
			Environment: "TST",
			MySQLInstance: model.MySQLInstance{
				Name:               "pippo",
				Version:            "1.1.1",
				Edition:            "Vanilla",
				Platform:           "Linux",
				Architecture:       "",
				Engine:             "",
				RedoLogEnabled:     "",
				CharsetServer:      "",
				CharsetSystem:      "",
				PageSize:           0,
				ThreadsConcurrency: 0,
				BufferPoolSize:     0,
				LogBufferSize:      0,
				SortBufferSize:     0,
				ReadOnly:           false,
				Databases: []model.MySQLDatabase{
					{Name: "pluto"}, {Name: "topolino"}, {Name: "minnie"},
				},
				TableSchemas: []model.MySQLTableSchema{
					{Name: "marte"}, {Name: "venere"}, {Name: "saturno"},
				},
				SegmentAdvisors: []model.MySQLSegmentAdvisor{},
			},
		},
		{
			Hostname:    "pluto",
			Location:    "",
			Environment: "TST",
			MySQLInstance: model.MySQLInstance{
				Name:               "Ash",
				Version:            "",
				Edition:            "Ketchup",
				Platform:           "",
				Architecture:       "",
				Engine:             "",
				RedoLogEnabled:     "",
				CharsetServer:      "",
				CharsetSystem:      "",
				PageSize:           0,
				ThreadsConcurrency: 0,
				BufferPoolSize:     0,
				LogBufferSize:      0,
				SortBufferSize:     0,
				ReadOnly:           false,
				Databases: []model.MySQLDatabase{
					{Name: "Picka"}, {Name: "Bulbasaur"},
				},
				TableSchemas:    []model.MySQLTableSchema{},
				SegmentAdvisors: []model.MySQLSegmentAdvisor{},
			},
		},
	}

	thisMoment := utils.P("2019-11-05T14:02:03+01:00")

	globalFilter := dto.GlobalFilter{
		Location:    "Dubai",
		Environment: "TEST",
		OlderThan:   thisMoment,
	}

	db.EXPECT().SearchMySQLInstances(globalFilter).
		Return(returned, nil)

	actual, err := as.SearchMySQLInstancesAsXLSX(globalFilter)
	require.NoError(t, err)

	assert.Equal(t, "Name", actual.GetCellValue("Instances", "A1"))
	assert.Equal(t, returned[0].MySQLInstance.Name, actual.GetCellValue("Instances", "A2"))
	assert.Equal(t, returned[1].MySQLInstance.Name, actual.GetCellValue("Instances", "A3"))

	assert.Equal(t, "Version", actual.GetCellValue("Instances", "B1"))
	assert.Equal(t, returned[0].Version, actual.GetCellValue("Instances", "B2"))
	assert.Equal(t, returned[1].Version, actual.GetCellValue("Instances", "B3"))

	assert.Equal(t, "Table Schemas", actual.GetCellValue("Instances", "Q1"))
	assert.Equal(t, "marte, venere, saturno", actual.GetCellValue("Instances", "Q2"))
	assert.Equal(t, "", actual.GetCellValue("Instances", "Q3"))
}
