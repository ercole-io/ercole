// Copyright (c) 2019 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInfoForFrontendDashboard_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := map[string]interface{}{
		"Alerts": []map[string]interface{}{
			{
				"AffectedHosts": 12,
				"Category":      "SYSTEM",
				"Code":          "NEW_SERVER",
				"Count":         12,
				"OldestAlert":   "2020-05-06T15:40:04.543+02:00",
				"Severity":      "NOTICE",
			},
			{
				"AffectedHosts": 12,
				"Category":      "SYSTEM",
				"Code":          "NEW_SERVER",
				"Count":         12,
				"OldestAlert":   "2020-05-06T15:40:04.543+02:00",
				"Severity":      "NOTICE",
			},
		},
		"Technologies": map[string]interface{}{
			"Technologies": []map[string]interface{}{
				{
					"Compliance": false,
					"Count":      7,
					"Name":       model.TechnologyOracleDatabase,
					"Used":       10,
					"TotalCost":  130,
					"PaidCost":   85,
					"HostsCount": 8,
				},
			},
			"Total": map[string]interface{}{
				"Compliant":  false,
				"TotalCost":  130,
				"PaidCost":   85,
				"Count":      7,
				"Used":       10,
				"HostsCount": 20,
			},
		},
		"Features": map[string]interface{}{
			"Oracle/Database": true,
			"Oracle/Exadata":  true,
		},
	}

	getTechnologiesUsageRes := map[string]float32{
		"Oracle/Database_HostsCount": 8,
		"Oracle/Exadata":             0,
	}
	db.EXPECT().
		GetTechnologiesUsage("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
		Return(getTechnologiesUsageRes, nil).AnyTimes().MinTimes(1)

	db.EXPECT().
		GetHostsCountStats("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
		Return(20, nil).AnyTimes().MinTimes(1)

	listLicensesRes := []interface{}{
		map[string]interface{}{
			"Compliance":       false,
			"Count":            4,
			"Used":             4,
			"_id":              "Partitioning",
			"TotalCost":        40,
			"PaidCost":         40,
			"CostPerProcessor": 10,
		},
		map[string]interface{}{
			"Compliance":       false,
			"Count":            3,
			"Used":             6,
			"_id":              "Diagnostics Pack",
			"TotalCost":        90,
			"PaidCost":         45,
			"CostPerProcessor": 15,
		},
		map[string]interface{}{
			"Compliance":       true,
			"Count":            5,
			"Used":             0,
			"_id":              "Advanced Analytics",
			"TotalCost":        0,
			"PaidCost":         5,
			"CostPerProcessor": 1,
		},
	}
	db.EXPECT().
		ListLicenses(false, "", false, -1, -1, "Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
		Return(listLicensesRes, nil).AnyTimes().MinTimes(1)

	searchAlertsRes := []interface{}{
		map[string]interface{}{
			"Category":      "SYSTEM",
			"AffectedHosts": 12,
			"Code":          "NEW_SERVER",
			"Count":         12,
			"OldestAlert":   "2020-05-06T15:40:04.543+02:00",
			"Severity":      "NOTICE",
		},
		map[string]interface{}{
			"Category":      "SYSTEM",
			"AffectedHosts": 12,
			"Code":          "NEW_SERVER",
			"Count":         12,
			"OldestAlert":   "2020-05-06T15:40:04.543+02:00",
			"Severity":      "NOTICE",
		},
	}
	db.EXPECT().SearchAlerts(
		"aggregated-category-severity",
		[]string{""}, "", false,
		-1, -1, "", "",
		utils.MIN_TIME, utils.P("2019-12-05T14:02:03Z"),
	).Return(
		searchAlertsRes,
		nil,
	)

	getTechnologiesUsageRes2 := map[string]float32{
		"Oracle/Database_HostsCount": 8,
		"Oracle/Exadata":             2,
	}
	db.EXPECT().
		GetTechnologiesUsage("", "", utils.MAX_TIME).
		Return(getTechnologiesUsageRes2, nil)

	res, err := as.GetInfoForFrontendDashboard("Italy", "PRD", utils.P("2019-12-05T14:02:03Z"))

	require.NoError(t, err)
	assert.JSONEq(t, utils.ToJSON(expectedRes), utils.ToJSON(res))
}

func TestGetInfoForFrontendDashboard_Fail1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().
		GetTechnologiesUsage("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
		Return(nil, aerrMock).AnyTimes().MinTimes(1)

	_, err := as.GetInfoForFrontendDashboard("Italy", "PRD", utils.P("2019-12-05T14:02:03Z"))

	require.Equal(t, aerrMock, err)
}

func TestGetInfoForFrontendDashboard_Fail2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	getTechnologiesUsageRes := map[string]float32{
		"Oracle/Database": 8,
		"Oracle/Exadata":  0,
	}

	db.EXPECT().
		GetTechnologiesUsage("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
		Return(getTechnologiesUsageRes, nil).AnyTimes().MinTimes(1)
	db.EXPECT().
		GetHostsCountStats("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
		Return(20, nil).AnyTimes().MinTimes(1)

	listLicensesRes := []interface{}{
		map[string]interface{}{
			"Compliance": false,
			"Count":      4,
			"Used":       4,
			"_id":        "Partitioning",
		},
		map[string]interface{}{
			"Compliance": false,
			"Count":      3,
			"Used":       6,
			"_id":        "Diagnostics Pack",
		},
		map[string]interface{}{
			"Compliance": true,
			"Count":      5,
			"Used":       0,
			"_id":        "Advanced Analytics",
		},
	}
	db.EXPECT().
		ListLicenses(false, "", false, -1, -1, "Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
		Return(listLicensesRes, nil).AnyTimes().MinTimes(1)

	searchAlertsRes := []interface{}{
		map[string]interface{}{
			"AffectedHosts": 12,
			"Code":          "NEW_SERVER",
			"Count":         12,
			"OldestAlert":   "2020-05-06T15:40:04.543+02:00",
			"Severity":      "NOTICE",
		},
		map[string]interface{}{
			"AffectedHosts": 12,
			"Code":          "NEW_SERVER",
			"Count":         12,
			"OldestAlert":   "2020-05-06T15:40:04.543+02:00",
			"Severity":      "NOTICE",
		},
	}
	db.EXPECT().SearchAlerts(
		"aggregated-category-severity",
		[]string{""}, "", false,
		-1, -1, "", "",
		utils.MIN_TIME, utils.P("2019-12-05T14:02:03Z"),
	).Return(
		searchAlertsRes,
		nil,
	)

	db.EXPECT().
		GetTechnologiesUsage("", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	_, err := as.GetInfoForFrontendDashboard("Italy", "PRD", utils.P("2019-12-05T14:02:03Z"))

	require.Equal(t, aerrMock, err)
}
