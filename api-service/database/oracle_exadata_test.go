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

package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
)

func (m *MongodbSuite) TestSearchOracleExadata() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	expectedOut := dto.OracleExadataResponse{
		Content: []dto.OracleExadata{},
		Metadata: dto.PagingMetadata{
			Empty:         true,
			First:         true,
			Last:          true,
			Number:        0,
			Size:          0,
			TotalElements: 0,
			TotalPages:    0,
		},
	}
	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchOracleExadata(false, []string{""}, "", false, -1, -1, "", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchOracleExadata(false, []string{""}, "", false, -1, -1, "France", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchOracleExadata(false, []string{""}, "", false, -1, -1, "", "", utils.P("1999-05-04T16:09:46.608+02:00"))
		m.Require().NoError(err)

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_anything", func(t *testing.T) {
		out, err := m.db.SearchOracleExadata(false, []string{"foobar"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	//FIXME: The sorting is not tested...
	expected := dto.OracleExadataResponse{
		Content: []dto.OracleExadata{
			{
				Id:        "5eba60d00b606515fdc2c554",
				CreatedAt: utils.P("2020-05-12T08:39:44.831Z"),
				DbServers: []dto.DbServers{
					{
						Hostname:           "zombie-0d1347d47a10b673a4df7aeeecc24a8a",
						Memory:             376,
						Model:              "X7-2",
						RunningCPUCount:    48,
						RunningPowerSupply: 2,
						SwVersion:          "19.2.4.0.0.190709",
						TempActual:         24,
						TotalCPUCount:      48,
						TotalPowerSupply:   2,
					},
					{
						Hostname:           "kantoor-43a6cdc54bb211eb127bca5c6651950c",
						Memory:             376,
						Model:              "X7-2",
						RunningCPUCount:    48,
						RunningPowerSupply: 2,
						SwVersion:          "19.2.4.0.0.190709",
						TempActual:         24,
						TotalCPUCount:      48,
						TotalPowerSupply:   2,
					},
				},
				Environment: "PROD",
				Hostname:    "test-exadata",
				IbSwitches: []dto.IbSwitches{
					{
						Hostname:  "off-df8b95a01746a464e69203c840a6a46a",
						Model:     "SUN_DCS_36p",
						SwVersion: "2.2.13-2.190326",
					},
					{
						Hostname:  "aspen-8d1d1b210625b1f1024b686135f889a1",
						Model:     "SUN_DCS_36p",
						SwVersion: "2.2.13-2.190326",
					},
				},
				Location: "Italy",
				StorageServers: []dto.StorageServers{
					{
						Hostname:           "s75-c2449b0e89e5a0b38401636eaa07abd5",
						Memory:             188,
						Model:              "X7-2L_High_Capacity",
						RunningCPUCount:    20,
						RunningPowerSupply: 2,
						SwVersion:          "19.2.4.0.0.190709",
						TempActual:         23,
						TotalCPUCount:      40,
						TotalPowerSupply:   2,
					},
					{
						Hostname:           "itl-b22fa37cad1326aba990cdec7facace2",
						Memory:             188,
						Model:              "X7-2L_High_Capacity",
						RunningCPUCount:    20,
						RunningPowerSupply: 2,
						SwVersion:          "19.2.4.0.0.190709",
						TempActual:         24,
						TotalCPUCount:      40,
						TotalPowerSupply:   2,
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
	}

	m.T().Run("should_search_test_exadata", func(t *testing.T) {
		out, err := m.db.SearchOracleExadata(false, []string{"test-exadata", "s75-c2449b0e89e5a0b38401636eaa07abd5"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.JSONEq(t, utils.ToJSON(expected), utils.ToJSON(out))
	})

	m.T().Run("fullmode", func(t *testing.T) {
		out, err := m.db.SearchOracleExadata(false, []string{}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		assert.JSONEq(t, utils.ToJSON(expected), utils.ToJSON(out))
	})
}
