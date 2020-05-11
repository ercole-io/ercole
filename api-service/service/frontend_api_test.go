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

	"github.com/amreo/ercole-services/utils"
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
				"Code":          "NEW_SERVER",
				"Count":         12,
				"OldestAlert":   "2020-05-06T15:40:04.543+02:00",
				"Severity":      "NOTICE",
			},
			{
				"AffectedHosts": 12,
				"Code":          "NEW_SERVER",
				"Count":         12,
				"OldestAlert":   "2020-05-06T15:40:04.543+02:00",
				"Severity":      "NOTICE",
			},
		},
		"Assets": map[string]interface{}{
			"Assets": []map[string]interface{}{
				{
					"Compliance": false,
					"Cost":       0,
					"Count":      0,
					"Name":       "Oracle/Database",
					"Used":       8,
				},
			},
			"Total": map[string]interface{}{
				"Compliant": false,
				"Cost":      0,
				"Count":     0,
				"Used":      8,
			},
		},
		"Features": map[string]interface{}{
			"Oracle/Database": true,
			"Oracle/Exadata":  true,
		},
	}

	getAssetsUsageRes := map[string]float32{
		"Oracle/Database": 8,
		"Oracle/Exadata":  0,
	}
	db.EXPECT().
		GetAssetsUsage("", false, "Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
		Return(getAssetsUsageRes, nil).AnyTimes().MinTimes(1)

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
		"aggregated-code-severity",
		[]string{""}, "", false,
		-1, -1, "Italy", "PRD",
		utils.MIN_TIME, utils.P("2019-12-05T14:02:03Z"),
	).Return(
		searchAlertsRes,
		nil,
	)

	getAssetsUsageRes2 := map[string]float32{
		"Oracle/Database": 8,
		"Oracle/Exadata":  2,
	}
	db.EXPECT().
		GetAssetsUsage("", false, "", "", utils.MAX_TIME).
		Return(getAssetsUsageRes2, nil)

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
		GetAssetsUsage("", false, "Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
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

	getAssetsUsageRes := map[string]float32{
		"Oracle/Database": 8,
		"Oracle/Exadata":  0,
	}
	db.EXPECT().
		GetAssetsUsage("", false, "Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
		Return(getAssetsUsageRes, nil).AnyTimes().MinTimes(1)

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
		"aggregated-code-severity",
		[]string{""}, "", false,
		-1, -1, "Italy", "PRD",
		utils.MIN_TIME, utils.P("2019-12-05T14:02:03Z"),
	).Return(
		searchAlertsRes,
		nil,
	)

	db.EXPECT().
		GetAssetsUsage("", false, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	_, err := as.GetInfoForFrontendDashboard("Italy", "PRD", utils.P("2019-12-05T14:02:03Z"))

	require.Equal(t, aerrMock, err)
}
