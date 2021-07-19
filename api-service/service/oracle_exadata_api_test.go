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
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"testing"
	"time"

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchOracleExadata_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	var expectedRes = []dto.OracleExadataResponse{
		{
			Content: []dto.OracleExadata{
				{
					Id:        "5e8c234b24f648a08585bd3e",
					CreatedAt: time.Time{},
					DbServers: []dto.DbServers{
						{
							Hostname:           "zombie-0d1347d47a10b673a4df7aeeecc24a8a",
							Memory:             48,
							Model:              "19.2.4.0.0.190709",
							RunningCPUCount:    48,
							RunningPowerSupply: 376,
							SwVersion:          "X7-2",
							TempActual:         2,
							TotalCPUCount:      2,
							TotalPowerSupply:   24.0,
						},
					},
					Environment: "PROD",
					Hostname:    "engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804",
					IbSwitches: []dto.IbSwitches{
						{
							Hostname:  "2.2.13-2.190326",
							Model:     "off-df8b95a01746a464e69203c840a6a46a",
							SwVersion: "SUN_DCS_36p",
						},
					},
					Location: "Italy",
					StorageServers: []dto.StorageServers{
						{
							Hostname:           "zombie-0d1347d47a10b673a4df7aeeecc24a8a",
							Memory:             48,
							Model:              "19.2.4.0.0.190709",
							RunningCPUCount:    48,
							RunningPowerSupply: 376,
							SwVersion:          "X7-2",
							TempActual:         2,
							TotalCPUCount:      2,
							TotalPowerSupply:   24.0,
						},
					},
				},
			},
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          1,
				TotalElements: 1,
				TotalPages:    0,
			},
		},
	}

	db.EXPECT().SearchOracleExadata(
		false, []string{"foo", "bar", "foobarx"}, "CPU",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchOracleExadata(
		false, "foo bar foobarx", "CPU",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchOracleExadata_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchOracleExadata(
		false, []string{"foo", "bar", "foobarx"}, "CPU",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchOracleExadata(
		false, "foo bar foobarx", "CPU",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestSearchOracleExadataAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	var expectedRes = []dto.OracleExadataResponse{
		{
			Content: []dto.OracleExadata{
				{
					Id:        "5e8c234b24f648a08585bd3e",
					CreatedAt: time.Time{},
					DbServers: []dto.DbServers{
						{
							Hostname:           "zombie-0d1347d47a10b673a4df7aeeecc24a8a",
							Memory:             48,
							Model:              "19.2.4.0.0.190709",
							RunningCPUCount:    48,
							RunningPowerSupply: 376,
							SwVersion:          "X7-2",
							TempActual:         2,
							TotalCPUCount:      2,
							TotalPowerSupply:   24.0,
						},
					},
					Environment:    "PROD",
					Hostname:       "engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804",
					IbSwitches:     nil,
					Location:       "",
					StorageServers: nil,
				},
			},
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          1,
				TotalElements: 1,
				TotalPages:    0,
			},
		},
	}

	filter := dto.GlobalFilter{
		Location:    "Italy",
		Environment: "TST",
		OlderThan:   utils.P("2020-06-10T11:54:59Z"),
	}

	db.EXPECT().SearchOracleExadata(true, []string{}, "", false, -1, -1, filter.Location, filter.Environment, filter.OlderThan).
		Return(expectedRes, nil).Times(1)

	actual, err := as.SearchOracleExadataAsXLSX(filter)
	require.NoError(t, err)
	assert.Equal(t, "engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804", actual.GetCellValue("engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804", "A1"))

	assert.Equal(t, "zombie-0d1347d47a10b673a4df7aeeecc24a8a", actual.GetCellValue("engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804", "A4"))
	assert.Equal(t, "19.2.4.0.0.190709", actual.GetCellValue("engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804", "B4"))
	assert.Equal(t, "2", actual.GetCellValue("engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804", "C4"))
	assert.Equal(t, "48", actual.GetCellValue("engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804", "D4"))
	assert.Equal(t, "X7-2", actual.GetCellValue("engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804", "E4"))
	assert.Equal(t, "24", actual.GetCellValue("engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804", "F4"))
}
