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
	"github.com/stretchr/testify/require"
)

//TODO
//func TestGetInfoForFrontendDashboard_Success(t *testing.T) {
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	db := NewMockMongoDatabaseInterface(mockCtrl)
//	as := APIService{
//		Database: db,
//	}

//	expectedRes := map[string]interface{}{
//		"technologies": map[string]interface{}{
//			"technologies": []map[string]interface{}{
//				{
//					"compliance": 7.0 / 10.0,
//					"product":    model.TechnologyOracleDatabase,
//					"hostsCount": 8,
//					"unpaidDues": 45,
//				},
//				{
//					"compliance": 1.0,
//					"unpaidDues": 0,
//					"product":    model.TechnologyMariaDBFoundationMariaDB,
//					"hostsCount": 0,
//				},
//				{
//					"compliance": 1.0,
//					"unpaidDues": 0,
//					"product":    model.TechnologyPostgreSQLPostgreSQL,
//					"hostsCount": 0,
//				},
//				{
//					"compliance": 1.0,
//					"unpaidDues": 0,
//					"product":    model.TechnologyOracleMySQL,
//					"hostsCount": 0,
//				},
//				{
//					"compliance": 1.0,
//					"unpaidDues": 0,
//					"product":    model.TechnologyMicrosoftSQLServer,
//					"hostsCount": 0,
//				},
//			},
//			"total": map[string]interface{}{
//				"compliance": 7.0 / 10.0,
//				"unpaidDues": 45,
//				"hostsCount": 20,
//			},
//		},
//		"features": map[string]interface{}{
//			"Oracle/Database": true,
//			"Oracle/Exadata":  true,
//		},
//	}

//	getTechnologiesUsageRes := map[string]float64{
//		"Oracle/Database_hostsCount": 8,
//		"Oracle/Exadata":             0,
//	}

//	db.EXPECT().
//		GetTechnologiesUsage("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
//		Return(getTechnologiesUsageRes, nil).AnyTimes().MinTimes(1)

//	db.EXPECT().
//		GetHostsCountStats("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
//		Return(20, nil).AnyTimes().MinTimes(1)

//	listLicensesRes := []interface{}{
//		map[string]interface{}{
//			"compliance":       false,
//			"count":            4,
//			"used":             4,
//			"_id":              "Partitioning",
//			"totalCost":        40,
//			"paidCost":         40,
//			"costPerProcessor": 10,
//		},
//		map[string]interface{}{
//			"compliance":       false,
//			"count":            3,
//			"used":             6,
//			"_id":              "Diagnostics Pack",
//			"totalCost":        90,
//			"paidCost":         45,
//			"costPerProcessor": 15,
//		},
//		map[string]interface{}{
//			"compliance":       true,
//			"count":            5,
//			"used":             0,
//			"_id":              "Advanced Analytics",
//			"totalCost":        0,
//			"paidCost":         5,
//			"costPerProcessor": 1,
//		},
//	}
//	db.EXPECT().
//		SearchLicenses("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
//		Return(listLicensesRes, nil).AnyTimes().MinTimes(1)

//	getTechnologiesUsageRes2 := map[string]float64{
//		"Oracle/Database_hostsCount": 8,
//		"Oracle/Exadata":             2,
//	}
//	db.EXPECT().
//		GetTechnologiesUsage("", "", utils.MAX_TIME).
//		Return(getTechnologiesUsageRes2, nil)

//	res, err := as.GetInfoForFrontendDashboard("Italy", "PRD", utils.P("2019-12-05T14:02:03Z"))

//	require.NoError(t, err)
//	assert.JSONEq(t, utils.ToJSON(expectedRes), utils.ToJSON(res))
//}

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

//TODO
//func TestGetInfoForFrontendDashboard_Fail2(t *testing.T) {
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	db := NewMockMongoDatabaseInterface(mockCtrl)
//	as := APIService{
//		Database: db,
//	}

//	getTechnologiesUsageRes := map[string]float64{
//		"Oracle/Database": 8,
//		"Oracle/Exadata":  0,
//	}

//	db.EXPECT().
//		GetTechnologiesUsage("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
//		Return(getTechnologiesUsageRes, nil).AnyTimes().MinTimes(1)
//	db.EXPECT().
//		GetHostsCountStats("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
//		Return(20, nil).AnyTimes().MinTimes(1)

//	listLicensesRes := []interface{}{
//		map[string]interface{}{
//			"Compliance": false,
//			"Count":      4,
//			"Used":       4,
//			"_id":        "Partitioning",
//		},
//		map[string]interface{}{
//			"Compliance": false,
//			"Count":      3,
//			"Used":       6,
//			"_id":        "Diagnostics Pack",
//		},
//		map[string]interface{}{
//			"Compliance": true,
//			"Count":      5,
//			"Used":       0,
//			"_id":        "Advanced Analytics",
//		},
//	}
//	db.EXPECT().
//		SearchLicenses("Italy", "PRD", utils.P("2019-12-05T14:02:03Z")).
//		Return(listLicensesRes, nil).AnyTimes().MinTimes(1)

//	db.EXPECT().
//		GetTechnologiesUsage("", "", utils.MAX_TIME).
//		Return(nil, aerrMock)

//	_, err := as.GetInfoForFrontendDashboard("Italy", "PRD", utils.P("2019-12-05T14:02:03Z"))

//	require.Equal(t, aerrMock, err)
//}
