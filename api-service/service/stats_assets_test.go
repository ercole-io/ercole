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

func TestGetTotalAssetsComplianceStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := map[string]interface{}{
		"Compliant":  false,
		"Count":      9,
		"Used":       12,
		"TotalCost":  130,
		"PaidCost":   85,
		"HostsCount": 20,
	}

	getAssetsUsageRes := map[string]float32{
		"Oracle/Exadata": 2,
	}

	db.EXPECT().
		GetAssetsUsage("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
		Return(getAssetsUsageRes, nil)
	db.EXPECT().
		GetHostsCountStats("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
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
		ListLicenses(false, "", false, -1, -1, "Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
		Return(listLicensesRes, nil).AnyTimes().MinTimes(1)

	res, err := as.GetTotalAssetsComplianceStats(
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.JSONEq(t, utils.ToJSON(expectedRes), utils.ToJSON(res))
}

func TestGetTotalAssetsComplianceStats_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().
		GetAssetsUsage("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
		Return(nil, aerrMock)

	_, err := as.GetTotalAssetsComplianceStats(
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)

	require.Equal(t, aerrMock, err)
}
