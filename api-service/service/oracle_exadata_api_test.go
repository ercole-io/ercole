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

	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
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

	expectedRes := []interface{}{
		map[string]interface{}{
			"CreatedAt": utils.P("2020-04-07T08:52:59.865+02:00"),
			"DBServers": []interface{}{
				map[string]interface{}{
					"RunningCPUCount":    48,
					"TotalCPUCount":      48,
					"SwVersion":          "19.2.4.0.0.190709",
					"Hostname":           "zombie-0d1347d47a10b673a4df7aeeecc24a8a",
					"Memory":             376,
					"Model":              "X7-2",
					"RunningPowerSupply": 2,
					"TotalPowerSupply":   2,
					"TempActual":         24.0,
				},
				map[string]interface{}{
					"RunningCPUCount":    48,
					"TotalCPUCount":      48,
					"SwVersion":          "19.2.4.0.0.190709",
					"Hostname":           "kantoor-43a6cdc54bb211eb127bca5c6651950c",
					"Memory":             376,
					"Model":              "X7-2",
					"RunningPowerSupply": 2,
					"TotalPowerSupply":   2,
					"TempActual":         24.0,
				},
			},
			"Environment": "PROD",
			"Hostname":    "engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804",
			"IBSwitches": []interface{}{
				map[string]interface{}{
					"SwVersion": "2.2.13-2.190326",
					"Hostname":  "off-df8b95a01746a464e69203c840a6a46a",
					"Model":     "SUN_DCS_36p",
				},
				map[string]interface{}{
					"SwVersion": "2.2.13-2.190326",
					"Hostname":  "aspen-8d1d1b210625b1f1024b686135f889a1",
					"Model":     "SUN_DCS_36p",
				},
			},
			"Location": "Italy",
			"StorageServers": []interface{}{
				map[string]interface{}{
					"RunningCPUCount":    20,
					"TotalCPUCount":      40,
					"SwVersion":          "19.2.4.0.0.190709",
					"Hostname":           "s75-c2449b0e89e5a0b38401636eaa07abd5",
					"Memory":             188,
					"Model":              "X7-2L_High_Capacity",
					"RunningPowerSupply": 2,
					"TotalPowerSupply":   2,
					"TempActual":         23.0,
				},
				map[string]interface{}{
					"RunningCPUCount":    20,
					"TotalCPUCount":      40,
					"SwVersion":          "19.2.4.0.0.190709",
					"Hostname":           "itl-b22fa37cad1326aba990cdec7facace2",
					"Memory":             188,
					"Model":              "X7-2L_High_Capacity",
					"RunningPowerSupply": 2,
					"TotalPowerSupply":   2,
					"TempActual":         24.0,
				},
			},
			"_id": utils.Str2oid("5e8c234b24f648a08585bd3e"),
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
