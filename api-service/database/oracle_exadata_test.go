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

	"github.com/ercole-io/ercole/v2/utils"
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
				"content": []interface{}{
					map[string]interface{}{
						"createdAt": utils.P("2020-05-12T10:39:44.831+02:00").Local(),
						"dbServers": []map[string]interface{}{
							{
								"runningCPUCount":    48,
								"totalCPUCount":      48,
								"swVersion":          "19.2.4.0.0.190709",
								"hostname":           "zombie-0d1347d47a10b673a4df7aeeecc24a8a",
								"memory":             376,
								"model":              "X7-2",
								"runningPowerSupply": 2,
								"totalPowerSupply":   2,
								"tempActual":         24.0,
							},
							{
								"runningCPUCount":    48,
								"totalCPUCount":      48,
								"swVersion":          "19.2.4.0.0.190709",
								"hostname":           "kantoor-43a6cdc54bb211eb127bca5c6651950c",
								"memory":             376,
								"model":              "X7-2",
								"runningPowerSupply": 2,
								"totalPowerSupply":   2,
								"tempActual":         24.0,
							},
						},
						"environment": "PROD",
						"hostname":    "test-exadata",
						"ibSwitches": []map[string]interface{}{
							{
								"swVersion": "2.2.13-2.190326",
								"hostname":  "off-df8b95a01746a464e69203c840a6a46a",
								"model":     "SUN_DCS_36p",
							},
							{
								"swVersion": "2.2.13-2.190326",
								"hostname":  "aspen-8d1d1b210625b1f1024b686135f889a1",
								"model":     "SUN_DCS_36p",
							},
						},
						"location": "Italy",
						"storageServers": []map[string]interface{}{
							{
								"runningCPUCount":    20,
								"totalCPUCount":      40,
								"swVersion":          "19.2.4.0.0.190709",
								"hostname":           "s75-c2449b0e89e5a0b38401636eaa07abd5",
								"memory":             188,
								"model":              "X7-2L_High_Capacity",
								"runningPowerSupply": 2,
								"totalPowerSupply":   2,
								"tempActual":         23.0,
							},
							{
								"runningCPUCount":    20,
								"totalCPUCount":      40,
								"swVersion":          "19.2.4.0.0.190709",
								"hostname":           "itl-b22fa37cad1326aba990cdec7facace2",
								"memory":             188,
								"model":              "X7-2L_High_Capacity",
								"runningPowerSupply": 2,
								"totalPowerSupply":   2,
								"tempActual":         24.0,
							},
						},
						"_id": utils.Str2oid("5eba60d00b606515fdc2c554"),
					},
				},
				"metadata": map[string]interface{}{
					"empty":         false,
					"first":         true,
					"last":          true,
					"number":        0,
					"size":          1,
					"totalElements": 1,
					"totalPages":    1,
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
				"createdAt": utils.P("2020-05-12T10:39:44.831+02:00").Local(),
				"dbServers": []map[string]interface{}{
					{
						"runningCPUCount":    48,
						"totalCPUCount":      48,
						"swVersion":          "19.2.4.0.0.190709",
						"hostname":           "zombie-0d1347d47a10b673a4df7aeeecc24a8a",
						"memory":             376,
						"model":              "X7-2",
						"runningPowerSupply": 2,
						"totalPowerSupply":   2,
						"tempActual":         24.0,
					},
					{
						"runningCPUCount":    48,
						"totalCPUCount":      48,
						"swVersion":          "19.2.4.0.0.190709",
						"hostname":           "kantoor-43a6cdc54bb211eb127bca5c6651950c",
						"memory":             376,
						"model":              "X7-2",
						"runningPowerSupply": 2,
						"totalPowerSupply":   2,
						"tempActual":         24.0,
					},
				},
				"environment": "PROD",
				"hostname":    "test-exadata",
				"ibSwitches": []map[string]interface{}{
					{
						"swVersion": "2.2.13-2.190326",
						"hostname":  "off-df8b95a01746a464e69203c840a6a46a",
						"model":     "SUN_DCS_36p",
					},
					{
						"swVersion": "2.2.13-2.190326",
						"hostname":  "aspen-8d1d1b210625b1f1024b686135f889a1",
						"model":     "SUN_DCS_36p",
					},
				},
				"location": "Italy",
				"storageServers": []map[string]interface{}{
					{
						"runningCPUCount":    20,
						"totalCPUCount":      40,
						"swVersion":          "19.2.4.0.0.190709",
						"hostname":           "s75-c2449b0e89e5a0b38401636eaa07abd5",
						"memory":             188,
						"model":              "X7-2L_High_Capacity",
						"runningPowerSupply": 2,
						"totalPowerSupply":   2,
						"tempActual":         23.0,
					},
					{
						"runningCPUCount":    20,
						"totalCPUCount":      40,
						"swVersion":          "19.2.4.0.0.190709",
						"hostname":           "itl-b22fa37cad1326aba990cdec7facace2",
						"memory":             188,
						"model":              "X7-2L_High_Capacity",
						"runningPowerSupply": 2,
						"totalPowerSupply":   2,
						"tempActual":         24.0,
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
				"createdAt": utils.P("2020-05-12T08:39:44.831Z").Local(),
				"dbServers": []map[string]interface{}{
					{
						"runningCPUCount":      48,
						"totalCPUCount":        48,
						"cellsrvServiceStatus": nil,
						"swVersion":            "19.2.4.0.0.190709",
						"runningFanCount":      16,
						"totalFanCount":        16,
						"fanStatus":            "normal",
						"hostname":             "zombie-0d1347d47a10b673a4df7aeeecc24a8a",
						"memory":               376,
						"model":                "X7-2",
						"msServiceStatus":      "running",
						"runningPowerSupply":   2,
						"totalPowerSupply":     2,
						"powerStatus":          "normal",
						"rsServiceStatus":      "running",
						"status":               "online",
						"tempActual":           24.0,
						"tempStatus":           "normal",
					},
					{
						"runningCPUCount":      48,
						"totalCPUCount":        48,
						"cellsrvServiceStatus": nil,
						"swVersion":            "19.2.4.0.0.190709",
						"runningFanCount":      16,
						"totalFanCount":        16,
						"fanStatus":            "normal",
						"hostname":             "kantoor-43a6cdc54bb211eb127bca5c6651950c",
						"memory":               376,
						"model":                "X7-2",
						"msServiceStatus":      "running",
						"runningPowerSupply":   2,
						"totalPowerSupply":     2,
						"powerStatus":          "normal",
						"rsServiceStatus":      "running",
						"status":               "online",
						"tempActual":           24.0,
						"tempStatus":           "normal",
					},
				},
				"environment": "PROD",
				"hostname":    "test-exadata",
				"ibSwitches": []map[string]interface{}{
					{
						"swVersion": "2.2.13-2.190326",
						"hostname":  "off-df8b95a01746a464e69203c840a6a46a",
						"model":     "SUN_DCS_36p",
					},
					{
						"swVersion": "2.2.13-2.190326",
						"hostname":  "aspen-8d1d1b210625b1f1024b686135f889a1",
						"model":     "SUN_DCS_36p",
					},
				},
				"location": "Italy",
				"storageServers": []map[string]interface{}{
					{
						"runningCPUCount": 20,
						"totalCPUCount":   40,
						"cellDisks": []map[string]interface{}{
							{
								"errCount": 0,
								"name":     "fanshop-5bde7badf2c9deceea5f615c3840c0b9",
								"status":   "normal",
								"usedPerc": 32,
							},
							{
								"errCount": 0,
								"name":     "globe-b31fa3756675a5c8ac052437ff7e439b",
								"status":   "normal",
								"usedPerc": 54,
							},
							{
								"errCount": 3,
								"name":     "srvc28-3807c977788c598cfa31fe21c0d3d5be",
								"status":   "normal",
								"usedPerc": 90,
							},
						},
						"cellsrvServiceStatus": "running",
						"swVersion":            "19.2.4.0.0.190709",
						"runningFanCount":      8,
						"totalFanCount":        8,
						"fanStatus":            "normal",
						"flashcacheMode":       "WriteBack",
						"hostname":             "s75-c2449b0e89e5a0b38401636eaa07abd5",
						"memory":               188,
						"model":                "X7-2L_High_Capacity",
						"msServiceStatus":      "running",
						"runningPowerSupply":   2,
						"totalPowerSupply":     2,
						"powerStatus":          "normal",
						"rsServiceStatus":      "running",
						"status":               "online",
						"tempActual":           23.0,
						"tempStatus":           "normal",
					},
					{
						"runningCPUCount": 20,
						"totalCPUCount":   40,
						"cellDisks": []map[string]interface{}{
							{
								"errCount": 0,
								"name":     "server52-390b5dac8c3b68e3471c657ef97a7ae6",
								"status":   "normal",
								"usedPerc": 32,
							},
							{
								"errCount": 1,
								"name":     "srvc07-5b888ab40dbd25e106309fd482f859a0",
								"status":   "normal",
								"usedPerc": 54,
							},
						},
						"cellsrvServiceStatus": "running",
						"swVersion":            "19.2.4.0.0.190709",
						"runningFanCount":      8,
						"totalFanCount":        8,
						"fanStatus":            "normal",
						"flashcacheMode":       "WriteBack",
						"hostname":             "itl-b22fa37cad1326aba990cdec7facace2",
						"memory":               188,
						"model":                "X7-2L_High_Capacity",
						"msServiceStatus":      "running",
						"runningPowerSupply":   2,
						"totalPowerSupply":     2,
						"powerStatus":          "normal",
						"rsServiceStatus":      "running",
						"status":               "online",
						"tempActual":           24.0,
						"tempStatus":           "normal",
					},
				},
				"_id": utils.Str2oid("5eba60d00b606515fdc2c554"),
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
