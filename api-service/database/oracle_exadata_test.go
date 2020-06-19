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

	"github.com/ercole-io/ercole/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestSearchOracleExadata() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_10.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchOracleExadata(false, []string{""}, "", false, -1, -1, "", "FOOBAR", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchOracleExadata(false, []string{""}, "", false, -1, -1, "France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchOracleExadata(false, []string{""}, "", false, -1, -1, "", "", utils.P("1999-05-04T16:09:46.608+02:00"))
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_anything", func(t *testing.T) {
		out, err := m.db.SearchOracleExadata(false, []string{"foobar"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_paging", func(t *testing.T) {
		out, err := m.db.SearchOracleExadata(false, []string{}, "", false, 0, 1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"Content": []interface{}{
					map[string]interface{}{
						"CreatedAt": utils.P("2020-05-12T10:39:44.831+02:00").Local(),
						"DBServers": []map[string]interface{}{
							{
								"CPUEnabled":   "48/48",
								"ExaSwVersion": "19.2.4.0.0.190709",
								"Hostname":     "zombie-0d1347d47a10b673a4df7aeeecc24a8a",
								"Memory":       "376GB",
								"Model":        "X7-2",
								"PowerCount":   "2/2",
								"TempActual":   "24.0",
							},
							{
								"CPUEnabled":   "48/48",
								"ExaSwVersion": "19.2.4.0.0.190709",
								"Hostname":     "kantoor-43a6cdc54bb211eb127bca5c6651950c",
								"Memory":       "376GB",
								"Model":        "X7-2",
								"PowerCount":   "2/2",
								"TempActual":   "24.0",
							},
						},
						"Environment": "PROD",
						"Hostname":    "test-exadata",
						"IBSwitches": []map[string]interface{}{
							{
								"ExaSwVersion": "2.2.13-2.190326",
								"Hostname":     "off-df8b95a01746a464e69203c840a6a46a",
								"Model":        "SUN_DCS_36p",
							},
							{
								"ExaSwVersion": "2.2.13-2.190326",
								"Hostname":     "aspen-8d1d1b210625b1f1024b686135f889a1",
								"Model":        "SUN_DCS_36p",
							},
						},
						"Location": "Italy",
						"StorageServers": []map[string]interface{}{
							{
								"CPUEnabled":   "20/40",
								"ExaSwVersion": "19.2.4.0.0.190709",
								"Hostname":     "s75-c2449b0e89e5a0b38401636eaa07abd5",
								"Memory":       "188GB",
								"Model":        "X7-2L_High_Capacity",
								"PowerCount":   "2/2",
								"TempActual":   "23.0",
							},
							{
								"CPUEnabled":   "20/40",
								"ExaSwVersion": "19.2.4.0.0.190709",
								"Hostname":     "itl-b22fa37cad1326aba990cdec7facace2",
								"Memory":       "188GB",
								"Model":        "X7-2L_High_Capacity",
								"PowerCount":   "2/2",
								"TempActual":   "24.0",
							},
						},
						"_id": utils.Str2oid("5eba60d00b606515fdc2c554"),
					},
				},
				"Metadata": map[string]interface{}{
					"Empty":         false,
					"First":         true,
					"Last":          true,
					"Number":        0,
					"Size":          1,
					"TotalElements": 1,
					"TotalPages":    1,
				},
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	//FIXME: The sorting is not tested...

	m.T().Run("should_search_test_exadata", func(t *testing.T) {
		out, err := m.db.SearchOracleExadata(false, []string{"test-exadata", "s75-c2449b0e89e5a0b38401636eaa07abd5"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"CreatedAt": utils.P("2020-05-12T10:39:44.831+02:00").Local(),
				"DBServers": []map[string]interface{}{
					{
						"CPUEnabled":   "48/48",
						"ExaSwVersion": "19.2.4.0.0.190709",
						"Hostname":     "zombie-0d1347d47a10b673a4df7aeeecc24a8a",
						"Memory":       "376GB",
						"Model":        "X7-2",
						"PowerCount":   "2/2",
						"TempActual":   "24.0",
					},
					{
						"CPUEnabled":   "48/48",
						"ExaSwVersion": "19.2.4.0.0.190709",
						"Hostname":     "kantoor-43a6cdc54bb211eb127bca5c6651950c",
						"Memory":       "376GB",
						"Model":        "X7-2",
						"PowerCount":   "2/2",
						"TempActual":   "24.0",
					},
				},
				"Environment": "PROD",
				"Hostname":    "test-exadata",
				"IBSwitches": []map[string]interface{}{
					{
						"ExaSwVersion": "2.2.13-2.190326",
						"Hostname":     "off-df8b95a01746a464e69203c840a6a46a",
						"Model":        "SUN_DCS_36p",
					},
					{
						"ExaSwVersion": "2.2.13-2.190326",
						"Hostname":     "aspen-8d1d1b210625b1f1024b686135f889a1",
						"Model":        "SUN_DCS_36p",
					},
				},
				"Location": "Italy",
				"StorageServers": []map[string]interface{}{
					{
						"CPUEnabled":   "20/40",
						"ExaSwVersion": "19.2.4.0.0.190709",
						"Hostname":     "s75-c2449b0e89e5a0b38401636eaa07abd5",
						"Memory":       "188GB",
						"Model":        "X7-2L_High_Capacity",
						"PowerCount":   "2/2",
						"TempActual":   "23.0",
					},
					{
						"CPUEnabled":   "20/40",
						"ExaSwVersion": "19.2.4.0.0.190709",
						"Hostname":     "itl-b22fa37cad1326aba990cdec7facace2",
						"Memory":       "188GB",
						"Model":        "X7-2L_High_Capacity",
						"PowerCount":   "2/2",
						"TempActual":   "24.0",
					},
				},
				"_id": utils.Str2oid("5eba60d00b606515fdc2c554"),
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("fullmode", func(t *testing.T) {
		out, err := m.db.SearchOracleExadata(true, []string{}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"CreatedAt": utils.P("2020-05-12T08:39:44.831Z").Local(),
				"DBServers": []map[string]interface{}{
					{
						"CPUEnabled":     "48/48",
						"CellsrvService": "-",
						"ExaSwVersion":   "19.2.4.0.0.190709",
						"FanCount":       "16/16",
						"FanStatus":      "normal",
						"Hostname":       "zombie-0d1347d47a10b673a4df7aeeecc24a8a",
						"Memory":         "376GB",
						"Model":          "X7-2",
						"MsService":      "running",
						"PowerCount":     "2/2",
						"PowerStatus":    "normal",
						"RsService":      "running",
						"Status":         "online",
						"TempActual":     "24.0",
						"TempStatus":     "normal",
					},
					{
						"CPUEnabled":     "48/48",
						"CellsrvService": "-",
						"ExaSwVersion":   "19.2.4.0.0.190709",
						"FanCount":       "16/16",
						"FanStatus":      "normal",
						"Hostname":       "kantoor-43a6cdc54bb211eb127bca5c6651950c",
						"Memory":         "376GB",
						"Model":          "X7-2",
						"MsService":      "running",
						"PowerCount":     "2/2",
						"PowerStatus":    "normal",
						"RsService":      "running",
						"Status":         "online",
						"TempActual":     "24.0",
						"TempStatus":     "normal",
					},
				},
				"Environment": "PROD",
				"Hostname":    "test-exadata",
				"IBSwitches": []map[string]interface{}{
					{
						"ExaSwVersion": "2.2.13-2.190326",
						"Hostname":     "off-df8b95a01746a464e69203c840a6a46a",
						"Model":        "SUN_DCS_36p",
					},
					{
						"ExaSwVersion": "2.2.13-2.190326",
						"Hostname":     "aspen-8d1d1b210625b1f1024b686135f889a1",
						"Model":        "SUN_DCS_36p",
					},
				},
				"Location": "Italy",
				"StorageServers": []map[string]interface{}{
					{
						"CPUEnabled": "20/40",
						"CellDisks": []map[string]interface{}{
							{
								"ErrCount": "0",
								"Name":     "fanshop-5bde7badf2c9deceea5f615c3840c0b9",
								"Status":   "normal",
								"UsedPerc": "32",
							},
							{
								"ErrCount": "0",
								"Name":     "globe-b31fa3756675a5c8ac052437ff7e439b",
								"Status":   "normal",
								"UsedPerc": "54",
							},
							{
								"ErrCount": "3",
								"Name":     "srvc28-3807c977788c598cfa31fe21c0d3d5be",
								"Status":   "normal",
								"UsedPerc": "90",
							},
						},
						"CellsrvService": "running",
						"ExaSwVersion":   "19.2.4.0.0.190709",
						"FanCount":       "8/8",
						"FanStatus":      "normal",
						"FlashcacheMode": "WriteBack",
						"Hostname":       "s75-c2449b0e89e5a0b38401636eaa07abd5",
						"Memory":         "188GB",
						"Model":          "X7-2L_High_Capacity",
						"MsService":      "running",
						"PowerCount":     "2/2",
						"PowerStatus":    "normal",
						"RsService":      "running",
						"Status":         "online",
						"TempActual":     "23.0",
						"TempStatus":     "normal",
					},
					{
						"CPUEnabled": "20/40",
						"CellDisks": []map[string]interface{}{
							{
								"ErrCount": "0",
								"Name":     "server52-390b5dac8c3b68e3471c657ef97a7ae6",
								"Status":   "normal",
								"UsedPerc": "32",
							},
							{
								"ErrCount": "1",
								"Name":     "srvc07-5b888ab40dbd25e106309fd482f859a0",
								"Status":   "normal",
								"UsedPerc": "54",
							},
						},
						"CellsrvService": "running",
						"ExaSwVersion":   "19.2.4.0.0.190709",
						"FanCount":       "8/8",
						"FanStatus":      "normal",
						"FlashcacheMode": "WriteBack",
						"Hostname":       "itl-b22fa37cad1326aba990cdec7facace2",
						"Memory":         "188GB",
						"Model":          "X7-2L_High_Capacity",
						"MsService":      "running",
						"PowerCount":     "2/2",
						"PowerStatus":    "normal",
						"RsService":      "running",
						"Status":         "online",
						"TempActual":     "24.0",
						"TempStatus":     "normal",
					},
				},
				"_id": utils.Str2oid("5eba60d00b606515fdc2c554"),
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
