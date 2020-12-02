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
	"sort"
	"testing"

	dto "github.com/ercole-io/ercole/api-service/dto"
	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadOracleDatabaseAgreementParts_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Log:      utils.NewLogger("TEST"),
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
	}
	as.loadOracleDatabaseAgreementParts()

	expected := []model.OracleDatabasePart{
		{PartID: "A11111", ItemDescription: "Database Enterprise Edition", Metric: model.AgreementPartMetricNamedUserPlusPerpetual, Cost: 42, Aliases: []string{"Db ENT"}},
		{PartID: "B22222", ItemDescription: "Database Standard Edition", Metric: model.AgreementPartMetricProcessorPerpetual, Cost: 43, Aliases: []string{"Db STD"}},
		{PartID: "C33333", ItemDescription: "Tuning", Metric: model.AgreementPartMetricStreamPerpetual, Cost: 44, Aliases: []string{"Tuning"}},
	}

	assert.ElementsMatch(t, expected, as.OracleDatabaseAgreementParts)
}

func TestGetOracleDatabaseAgreementPartsList_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabasePart{
			{
				PartID:          "Pippo",
				ItemDescription: "Pluto",
				Metric:          "Topolino",
				Cost:            12,
				Aliases:         []string{"Minny"},
			},
		},
	}
	res, err := as.GetOracleDatabaseAgreementPartsList()
	require.NoError(t, err)
	assert.Equal(t, []model.OracleDatabasePart{
		{
			PartID:          "Pippo",
			ItemDescription: "Pluto",
			Metric:          "Topolino",
			Cost:            12,
			Aliases:         []string{"Minny"},
		},
	}, res)
}

func TestGetLicensesCompliance(t *testing.T) {
	var sampleAgreements []dto.OracleDatabaseAgreementFE = []dto.OracleDatabaseAgreementFE{
		{
			ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			AgreementID:     "",
			CSI:             "",
			PartID:          "PID001",
			ItemDescription: "",
			Metric:          "",
			ReferenceNumber: "",
			Unlimited:       false,
			Count:           50,
			CatchAll:        false,
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{Hostname: "pippo"},
				{Hostname: "pluto"},
			},
			AvailableCount: 50,
			LicensesCount:  0,
			UsersCount:     0,
		},
		{
			ID:              utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
			AgreementID:     "",
			CSI:             "",
			PartID:          "PID002",
			ItemDescription: "",
			Metric:          "",
			ReferenceNumber: "",
			Unlimited:       false,
			Count:           75,
			CatchAll:        false,
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{Hostname: "topolino"},
				{Hostname: "minnie"},
			},
			AvailableCount: 75,
			LicensesCount:  0,
			UsersCount:     0,
		},
		{
			ID:              utils.Str2oid("cccccccccccccccccccccccc"),
			AgreementID:     "",
			CSI:             "",
			PartID:          "PID003",
			ItemDescription: "",
			Metric:          "",
			ReferenceNumber: "",
			Unlimited:       true,
			Count:           75,
			CatchAll:        false,
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{Hostname: "topolino"},
				{Hostname: "minnie"},
			},
			AvailableCount: 75,
			LicensesCount:  0,
			UsersCount:     0,
		},
	}

	var sampleHosts []dto.HostUsingOracleDatabaseLicenses = []dto.HostUsingOracleDatabaseLicenses{
		{LicenseName: "alias1", Name: "test1", Type: "host", LicenseCount: 3, OriginalCount: 3},
		{LicenseName: "alias1", Name: "pluto", Type: "host", LicenseCount: 1.5, OriginalCount: 1.5},
		{LicenseName: "alias1", Name: "pippo", Type: "host", LicenseCount: 5.5, OriginalCount: 5.5},

		{LicenseName: "alias2", Name: "topolino", Type: "cluster", LicenseCount: 7, OriginalCount: 7},
		{LicenseName: "alias2", Name: "minnie", Type: "host", LicenseCount: 4, OriginalCount: 4},
		{LicenseName: "alias2", Name: "minnie", Type: "host", LicenseCount: 8, OriginalCount: 8},

		{LicenseName: "alias3", Name: "minnie", Type: "host", LicenseCount: 0.5, OriginalCount: 0.5},
		{LicenseName: "alias3", Name: "pippo", Type: "host", LicenseCount: 0.5, OriginalCount: 0.5},
		{LicenseName: "alias3", Name: "test2", Type: "host", LicenseCount: 4, OriginalCount: 4},
		{LicenseName: "alias3", Name: "test3", Type: "cluster", LicenseCount: 6, OriginalCount: 6},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:                     db,
		Log:                          utils.NewLogger("TEST"),
		OracleDatabaseAgreementParts: sampleParts,
	}

	gomock.InOrder(
		db.EXPECT().
			ListOracleDatabaseAgreements().
			Return(sampleAgreements, nil),
		db.EXPECT().
			ListHostUsingOracleDatabaseLicenses().
			Return(sampleHosts, nil),
	)

	actual, err := as.GetOracleDatabaseLicensesCompliance()
	require.NoError(t, err)

	compliance := 75.0 / 275.0
	expected := []dto.OracleDatabaseLicenseUsage{
		{PartID: "PID001", ItemDescription: "itemDesc1", Metric: "Processor Perpetual", Consumed: 7, Covered: 7, Compliance: 1, Unlimited: false},
		{PartID: "PID002", ItemDescription: "itemDesc2", Metric: "Named User Plus Perpetual", Consumed: 275, Covered: 75, Compliance: compliance, Unlimited: false},
		{PartID: "PID003", ItemDescription: "itemDesc3", Metric: "Computer Perpetual", Consumed: 0.5, Covered: 0.5, Compliance: 1, Unlimited: true},
	}

	sort.Slice(actual, func(i, j int) bool {
		return actual[i].PartID < actual[j].PartID
	})

	require.Equal(t, expected, actual)
}
