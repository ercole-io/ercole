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

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dto "github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestGetTotalTechnologiesComplianceStats_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	getTechnologiesUsageRes := map[string]float64{
		"Oracle/Exadata": 2,
	}
	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().
		GetHostsCountUsingTechnologies("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
		Return(getTechnologiesUsageRes, nil)
	db.EXPECT().
		GetHostsCountStats("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
		Return(20, nil).AnyTimes().MinTimes(1)
	db.EXPECT().
		GetOracleDatabaseLicenseTypes().
		Return(licenseTypes, nil)

	returnedAgreements := []dto.OracleDatabaseAgreementFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			AgreementID:              "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          55,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 55,
			AvailableLicensesPerUser: 0,
		},
	}
	db.EXPECT().ListOracleDatabaseAgreements().Return(returnedAgreements, nil)

	returnedHosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  100,
			LicenseTypeID: "PID002",
			OriginalCount: 100,
			Type:          "host",
		},
	}
	db.EXPECT().ListHostUsingOracleDatabaseLicenses().Return(returnedHosts, nil)

	res, err := as.GetTotalTechnologiesComplianceStats(
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)

	require.NoError(t, err)

	expectedRes := map[string]interface{}{
		"compliance": 0.55,
		"unpaidDues": 0,
		"hostsCount": 20,
	}
	assert.JSONEq(t, utils.ToJSON(expectedRes), utils.ToJSON(res))
}

func TestGetTotalTechnologiesComplianceStats_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().
		GetHostsCountUsingTechnologies("Italy", "PROD", utils.P("2020-12-05T14:02:03Z")).
		Return(nil, aerrMock)

	_, err := as.GetTotalTechnologiesComplianceStats(
		"Italy", "PROD", utils.P("2020-12-05T14:02:03Z"),
	)

	require.Equal(t, aerrMock, err)
}
