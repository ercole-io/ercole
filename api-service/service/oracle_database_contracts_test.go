// Copyright (c) 2022 Sorint.lab S.p.A.
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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

var lt1 = model.OracleDatabaseLicenseType{
	ID:              "PID001",
	ItemDescription: "itemDesc1",
	Aliases:         []string{"alias1"},
	Metric:          "metric1",
}

func TestAddOracleDatabaseContract_Success_InsertNew(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	contract := model.OracleDatabaseContract{
		ContractID:      "AID001",
		LicenseTypeID:   "PID001",
		CSI:             "CSI001",
		ReferenceNumber: "RF0001",
		Unlimited:       true,
		Count:           30,
		Basket:          true,
		Restricted:      false,
		Hosts:           []string{"test-db", "ercsoldbx"},
	}

	expectedAgr := contract
	expectedAgr.ID = utils.Str2oid("000000000000000000000001")
	commonFilters := dto.NewSearchHostsFilters()

	gomock.InOrder(
		db.EXPECT().SearchHosts("hostnames",
			commonFilters).Return([]map[string]interface{}{
			{"hostname": "test-db"},
			{"hostname": "foobar"},
			{"hostname": "ercsoldbx"},
		}, nil),
		db.EXPECT().GetOracleDatabaseLicenseType("PID001").
			Return(&lt1, nil),
		db.EXPECT().InsertOracleDatabaseContract(expectedAgr).
			Return(nil),
	)

	searchedContractItem := dto.OracleDatabaseContractFE{
		ID:                       expectedAgr.ID,
		ContractID:               contract.ContractID,
		CSI:                      contract.CSI,
		LicenseTypeID:            contract.LicenseTypeID,
		ItemDescription:          "",
		Metric:                   "",
		ReferenceNumber:          "",
		Unlimited:                false,
		Basket:                   false,
		Restricted:               false,
		Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{},
		LicensesPerCore:          0,
		LicensesPerUser:          0,
		AvailableLicensesPerCore: 0,
		AvailableLicensesPerUser: 0,
	}
	as.mockGetOracleDatabaseContracts = func(filters dto.GetOracleDatabaseContractsFilter) ([]dto.OracleDatabaseContractFE, error) {
		return []dto.OracleDatabaseContractFE{searchedContractItem}, nil
	}

	res, err := as.AddOracleDatabaseContract(contract)
	require.NoError(t, err)
	assert.Equal(t,
		searchedContractItem,
		*res)
}

func TestAddOracleDatabaseContracts_Fail(t *testing.T) {
	commonFilters := dto.NewSearchHostsFilters()
	t.Run("Fail: can't find host", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		db := NewMockMongoDatabaseInterface(mockCtrl)

		as := APIService{
			Database: db,
			Config: config.Configuration{
				ResourceFilePath: "../../resources",
			},
			TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
			NewObjectID: utils.NewObjectIDForTests(),
		}

		addRequest := model.OracleDatabaseContract{
			ContractID:      "AID001",
			LicenseTypeID:   "PID001",
			CSI:             "CSI001",
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			Basket:          true,
			Restricted:      false,
			Hosts: []string{
				"test-db",
				"ercsoldbx",
			},
		}

		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames",
				commonFilters).
				Return([]map[string]interface{}{
					{"hostname": "paperino"},
					{"hostname": "pippo"},
					{"hostname": "pluto"},
				}, nil),
		)

		res, err := as.AddOracleDatabaseContract(addRequest)
		assert.EqualError(t, err, utils.ErrHostNotFound.Error())
		assert.Nil(t, res)
	})

	t.Run("Fail: can't find licenseType", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		db := NewMockMongoDatabaseInterface(mockCtrl)

		as := APIService{
			Database: db,
			Config: config.Configuration{
				ResourceFilePath: "../../resources",
			},
			TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
			NewObjectID: utils.NewObjectIDForTests(),
		}

		contractWrongLicenseType := model.OracleDatabaseContract{
			ContractID:      "AID001",
			LicenseTypeID:   "xxxxxx",
			CSI:             "CSI001",
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			Basket:          true,
			Restricted:      false,
			Hosts: []string{
				"test-db",
				"ercsoldbx",
			},
		}
		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames",
				commonFilters).
				Return([]map[string]interface{}{
					{"hostname": "test-db"},
					{"hostname": "ercsoldbx"},
				}, nil),
			db.EXPECT().GetOracleDatabaseLicenseType("xxxxxx").
				Return(nil, nil),
		)

		res, err := as.AddOracleDatabaseContract(contractWrongLicenseType)

		assert.EqualError(t, err, utils.ErrOracleDatabaseLicenseTypeIDNotFound.Error())
		assert.Nil(t, res)
	})
}

func TestUpdateOracleDatabaseContract_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := NewMockMongoDatabaseInterface(mockCtrl)

	commonFilters := dto.NewSearchHostsFilters()

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	contract := model.OracleDatabaseContract{
		ContractID:      "AID001",
		CSI:             "CSI001",
		ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		LicenseTypeID:   lt1.ID,
		ReferenceNumber: "RF0001",
		Unlimited:       true,
		Count:           30,
		Basket:          true,
		Restricted:      false,
		Hosts:           []string{"test-db", "ercsoldbx"},
	}

	gomock.InOrder(
		db.EXPECT().SearchHosts("hostnames",
			commonFilters).
			Return([]map[string]interface{}{
				{"hostname": "test-db"},
				{"hostname": "foobar"},
				{"hostname": "ercsoldbx"},
			}, nil),
		db.EXPECT().GetOracleDatabaseLicenseType("PID001").
			Return(&lt1, nil),
		db.EXPECT().UpdateOracleDatabaseContract(contract).Return(nil),
	)

	searchedContractItem := dto.OracleDatabaseContractFE{
		ID:                       contract.ID,
		ContractID:               contract.ContractID,
		CSI:                      contract.CSI,
		LicenseTypeID:            contract.LicenseTypeID,
		ItemDescription:          "",
		Metric:                   "",
		ReferenceNumber:          "",
		Unlimited:                false,
		Basket:                   false,
		Restricted:               false,
		Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{},
		LicensesPerCore:          0,
		LicensesPerUser:          0,
		AvailableLicensesPerCore: 0,
		AvailableLicensesPerUser: 0,
	}
	as.mockGetOracleDatabaseContracts = func(filters dto.GetOracleDatabaseContractsFilter) ([]dto.OracleDatabaseContractFE, error) {
		return []dto.OracleDatabaseContractFE{searchedContractItem}, nil
	}

	actualContract, err := as.UpdateOracleDatabaseContract(contract)
	require.NoError(t, err)
	assert.Equal(t, searchedContractItem, *actualContract)
}

func TestUpdateOracleDatabaseContract_Fail_LicenseTypeIdNotValid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db := NewMockMongoDatabaseInterface(mockCtrl)

	commonFilters := dto.NewSearchHostsFilters()

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	contract := model.OracleDatabaseContract{
		ContractID:      "AID001",
		CSI:             "CSI001",
		ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		LicenseTypeID:   "invalidLicenseTypeID",
		ReferenceNumber: "RF0001",
		Unlimited:       true,
		Count:           30,
		Basket:          true,
		Restricted:      false,
		Hosts:           []string{"test-db", "ercsoldbx"},
	}

	gomock.InOrder(
		db.EXPECT().SearchHosts("hostnames",
			commonFilters).
			Return([]map[string]interface{}{
				{"hostname": "test-db"},
				{"hostname": "foobar"},
				{"hostname": "ercsoldbx"},
			}, nil),
		db.EXPECT().GetOracleDatabaseLicenseType("invalidLicenseTypeID").
			Return(nil, nil),
	)

	actual, err := as.UpdateOracleDatabaseContract(contract)

	assert.EqualError(t, err, utils.ErrOracleDatabaseLicenseTypeIDNotFound.Error())
	assert.Nil(t, actual)
}

func TestGetOracleDatabaseContracts_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}

	returnedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
		},
	}
	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "PID002",
				DbName:        "test-dbname",
				Hostname:      "test-db",
				UsedLicenses:  3,
			},
		},
	}
	clusters := []dto.Cluster{}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "test-db",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    false,
				VeritasClusterHostnames: []string{},
			},
			Info: model.Host{
				CPUCores: 42,
			},
		},
	}
	globalFilterAny := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	expectedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 3, Hostname: "test-db", TotalCoveredLicensesCount: 3, ConsumedLicensesCount: 3}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
			CoveredLicenses:          3,
		},
	}

	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return(hostdatas, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).
			Return(returnedContracts, nil),

		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
	)

	res, err := as.GetOracleDatabaseContracts(dto.NewGetOracleDatabaseContractsFilter())
	require.NoError(t, err)
	assert.Equal(t, expectedContracts, res)
}

func TestGetOracleDatabaseContractsCluster_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}

	returnedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 3, Hostname: "test-db", TotalCoveredLicensesCount: 3}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
		},
	}
	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "PID002",
				DbName:        "test-dbname",
				Hostname:      "test-db",
				UsedLicenses:  3,
			},
		},
	}
	clusters := []dto.Cluster{
		{
			ID:                          [12]byte{},
			CreatedAt:                   time.Time{},
			Hostname:                    "bart",
			HostnameAgentVirtualization: "",
			Name:                        "bart",
			Environment:                 "",
			Location:                    "",
			FetchEndpoint:               "",
			CPU:                         0,
			Sockets:                     0,
			Type:                        "vmware",
			VirtualizationNodes:         []string{},
			VirtualizationNodesCount:    0,
			VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
			VMs: []dto.VM{
				{
					CappedCPU:          false,
					Hostname:           "test-db",
					Name:               "test-db",
					VirtualizationNode: "",
					IsErcoleInstalled:  false,
				},
			},
			VMsCount:            0,
			VMsErcoleAgentCount: 0,
		},
	}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "test-db",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    false,
				VeritasClusterHostnames: []string{},
			},
			Info: model.Host{
				CPUCores: 42,
			},
		},
	}
	globalFilterAny := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	expectedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 3, Hostname: "test-db", TotalCoveredLicensesCount: 3, ConsumedLicensesCount: 0}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
			CoveredLicenses:          0,
		},
	}

	cluster := dto.Cluster{
		ID:                          [12]byte{},
		CreatedAt:                   time.Time{},
		Hostname:                    "bart",
		HostnameAgentVirtualization: "",
		Name:                        "bart",
		Environment:                 "",
		Location:                    "",
		FetchEndpoint:               "",
		CPU:                         0,
		Sockets:                     0,
		Type:                        "vmware",
		VirtualizationNodes:         []string{},
		VirtualizationNodesCount:    0,
		VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
		VMs: []dto.VM{
			{
				CappedCPU:          false,
				Hostname:           "test-db",
				Name:               "test-db",
				VirtualizationNode: "",
				IsErcoleInstalled:  false,
			},
		},
		VMsCount:            0,
		VMsErcoleAgentCount: 0,
	}

	db.EXPECT().ExistHostdata(gomock.Any()).Return(false, nil).AnyTimes()
	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return(hostdatas, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	db.EXPECT().GetCluster("bart", utils.MAX_TIME).Return(&cluster, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).
			Return(returnedContracts, nil),

		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
	)

	res, err := as.GetOracleDatabaseContracts(dto.NewGetOracleDatabaseContractsFilter())
	require.NoError(t, err)
	assert.Equal(t, expectedContracts, res)
}

func TestGetOracleDatabaseContractsClusterCappedCPU_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}

	returnedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 3, Hostname: "test-db", TotalCoveredLicensesCount: 3}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
		},
	}
	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "PID002",
				DbName:        "test-dbname",
				Hostname:      "test-db",
				UsedLicenses:  3,
			},
		},
	}
	clusters := []dto.Cluster{
		{
			ID:                          [12]byte{},
			CreatedAt:                   time.Time{},
			Hostname:                    "bart",
			HostnameAgentVirtualization: "",
			Name:                        "bart",
			Environment:                 "",
			Location:                    "",
			FetchEndpoint:               "",
			CPU:                         0,
			Sockets:                     0,
			Type:                        "vmware",
			VirtualizationNodes:         []string{},
			VirtualizationNodesCount:    0,
			VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
			VMs: []dto.VM{
				{
					CappedCPU:          true,
					Hostname:           "test-db",
					Name:               "test-db",
					VirtualizationNode: "",
					IsErcoleInstalled:  false,
				},
			},
			VMsCount:            0,
			VMsErcoleAgentCount: 0,
		},
	}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "test-db",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    false,
				VeritasClusterHostnames: []string{},
			},
			Info: model.Host{
				CPUCores: 42,
			},
		},
	}
	globalFilterAny := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	expectedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 3, Hostname: "test-db", TotalCoveredLicensesCount: 3, ConsumedLicensesCount: 0}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
			CoveredLicenses:          0,
		},
	}

	cluster := dto.Cluster{
		ID:                          [12]byte{},
		CreatedAt:                   time.Time{},
		Hostname:                    "bart",
		HostnameAgentVirtualization: "",
		Name:                        "bart",
		Environment:                 "",
		Location:                    "",
		FetchEndpoint:               "",
		CPU:                         0,
		Sockets:                     0,
		Type:                        "vmware",
		VirtualizationNodes:         []string{},
		VirtualizationNodesCount:    0,
		VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
		VMs: []dto.VM{
			{
				CappedCPU:          true,
				Hostname:           "test-db",
				Name:               "test-db",
				VirtualizationNode: "",
				IsErcoleInstalled:  false,
			},
		},
		VMsCount:            0,
		VMsErcoleAgentCount: 0,
	}

	db.EXPECT().ExistHostdata(gomock.Any()).Return(false, nil).AnyTimes()
	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return(hostdatas, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	db.EXPECT().GetCluster("bart", utils.MAX_TIME).Return(&cluster, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).
			Return(returnedContracts, nil),

		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
	)

	res, err := as.GetOracleDatabaseContracts(dto.NewGetOracleDatabaseContractsFilter())
	require.NoError(t, err)
	assert.Equal(t, expectedContracts, res)
}

func TestGetOracleDatabaseContractsClusterCappedCPU2_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}

	returnedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:      "AID001",
			CSI:             "CSI001",
			LicenseTypeID:   "PID002",
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Basket:          false,
			Restricted:      false,
			Hosts: []dto.OracleDatabaseContractAssociatedHostFE{
				{CoveredLicensesCount: 3, Hostname: "test-db", TotalCoveredLicensesCount: 3},
				{CoveredLicensesCount: 5, Hostname: "test-db2", TotalCoveredLicensesCount: 5},
			},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
		},
	}
	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "PID002",
				DbName:        "test-dbname",
				Hostname:      "test-db",
				UsedLicenses:  3,
			},
			{
				LicenseTypeID: "PID002",
				DbName:        "test-dbname",
				Hostname:      "test-db2",
				UsedLicenses:  5,
			},
		},
	}
	clusters := []dto.Cluster{
		{
			ID:                          [12]byte{},
			CreatedAt:                   time.Time{},
			Hostname:                    "bart",
			HostnameAgentVirtualization: "",
			Name:                        "bart",
			Environment:                 "",
			Location:                    "",
			FetchEndpoint:               "",
			CPU:                         0,
			Sockets:                     0,
			Type:                        "vmware",
			VirtualizationNodes:         []string{},
			VirtualizationNodesCount:    0,
			VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
			VMs: []dto.VM{
				{
					CappedCPU:          true,
					Hostname:           "test-db",
					Name:               "test-db",
					VirtualizationNode: "",
					IsErcoleInstalled:  false,
				},
				{
					CappedCPU:          false,
					Hostname:           "test-db2",
					Name:               "test-db2",
					VirtualizationNode: "",
					IsErcoleInstalled:  false,
				},
			},
			VMsCount:            0,
			VMsErcoleAgentCount: 0,
		},
	}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "test-db",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    false,
				VeritasClusterHostnames: []string{},
			},
			Info: model.Host{
				CPUCores: 42,
			},
		},
		{
			Hostname: "test-db2",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    false,
				VeritasClusterHostnames: []string{},
			},
			Info: model.Host{
				CPUCores: 42,
			},
		},
	}
	globalFilterAny := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	expectedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:              utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:      "AID001",
			CSI:             "CSI001",
			LicenseTypeID:   "PID002",
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Basket:          false,
			Restricted:      false,
			Hosts: []dto.OracleDatabaseContractAssociatedHostFE{
				{CoveredLicensesCount: 3, Hostname: "test-db", TotalCoveredLicensesCount: 3, ConsumedLicensesCount: 0},
				{CoveredLicensesCount: 5, Hostname: "test-db2", TotalCoveredLicensesCount: 5, ConsumedLicensesCount: 0},
			},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
			CoveredLicenses:          0,
		},
	}

	cluster := dto.Cluster{
		ID:                          [12]byte{},
		CreatedAt:                   time.Time{},
		Hostname:                    "bart",
		HostnameAgentVirtualization: "",
		Name:                        "bart",
		Environment:                 "",
		Location:                    "",
		FetchEndpoint:               "",
		CPU:                         0,
		Sockets:                     0,
		Type:                        "vmware",
		VirtualizationNodes:         []string{},
		VirtualizationNodesCount:    0,
		VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
		VMs: []dto.VM{
			{
				CappedCPU:          true,
				Hostname:           "test-db",
				Name:               "test-db",
				VirtualizationNode: "",
				IsErcoleInstalled:  false,
			},
			{
				CappedCPU:          false,
				Hostname:           "test-db2",
				Name:               "test-db2",
				VirtualizationNode: "",
				IsErcoleInstalled:  false,
			},
		},
		VMsCount:            0,
		VMsErcoleAgentCount: 0,
	}

	db.EXPECT().ExistHostdata(gomock.Any()).Return(false, nil).AnyTimes()
	db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).
		Return(returnedContracts, nil)

	db.EXPECT().
		SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(&oracleLics, nil)
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(licenseTypes, nil)
	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return(hostdatas, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	db.EXPECT().GetCluster("bart", utils.MAX_TIME).Return(&cluster, nil).AnyTimes()
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(licenseTypes, nil)

	res, err := as.GetOracleDatabaseContracts(dto.NewGetOracleDatabaseContractsFilter())
	require.NoError(t, err)
	assert.Equal(t, expectedContracts, res)
}

func TestGetOracleDatabaseContracts_SuccessFilter1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}

	returnedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 7,
			AvailableLicensesPerUser: 0,
		},
	}

	oracleLics := dto.OracleDatabaseUsedLicenseSearchResponse{
		Content: []dto.OracleDatabaseUsedLicense{
			{
				LicenseTypeID: "PID002",
				DbName:        "test-dbname",
				Hostname:      "test-db",
				UsedLicenses:  3,
			},
		},
	}
	clusters := []dto.Cluster{}
	hostdatas := []model.HostDataBE{
		{
			Hostname: "test-db",
			ClusterMembershipStatus: model.ClusterMembershipStatus{
				OracleClusterware:       false,
				SunCluster:              false,
				HACMP:                   false,
				VeritasClusterServer:    false,
				VeritasClusterHostnames: []string{},
			},
			Info: model.Host{
				CPUCores: 42,
			},
		},
	}
	globalFilterAny := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}
	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}

	db.EXPECT().GetHostDatas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).
		Return(hostdatas, nil).AnyTimes()
	db.EXPECT().GetClusters(globalFilterAny).
		Return(clusters, nil).AnyTimes()
	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).
			Return(returnedContracts, nil),
		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(parts, nil),
	)

	res, err := as.GetOracleDatabaseContracts(dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "asddfa",
		Unlimited:                   "",
		Basket:                      "",
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerCoreLTE: -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerCoreLTE:          -1,
		LicensesPerUserGTE:          -1,
		LicensesPerUserLTE:          -1,
	})
	require.NoError(t, err)
	assert.Empty(t, res)

	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).
			Return(returnedContracts, nil),

		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(parts, nil),
	)

	res, err = as.GetOracleDatabaseContracts(dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "asddfa",
		Unlimited:                   "",
		Basket:                      "",
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerCoreLTE: -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerCoreLTE:          -1,
		LicensesPerUserGTE:          -1,
		LicensesPerUserLTE:          -1,
	})
	require.NoError(t, err)
	assert.Empty(t, res)

	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).
			Return(returnedContracts, nil),

		db.EXPECT().
			SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(&oracleLics, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(licenseTypes, nil),
		db.EXPECT().GetOracleDatabaseLicenseTypes().
			Return(parts, nil),
	)

	res, err = as.GetOracleDatabaseContracts(dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "asddfa",
		Unlimited:                   "",
		Basket:                      "",
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerCoreLTE: -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerCoreLTE:          -1,
		LicensesPerUserGTE:          -1,
		LicensesPerUserLTE:          -1,
	})

	require.NoError(t, err)
	assert.Empty(t, res)
}

func TestGetOracleDatabaseContracts_Failed2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	returnedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 7,
			AvailableLicensesPerUser: 0,
		},
	}

	gomock.InOrder(
		db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).
			Return(returnedContracts, nil),
		db.EXPECT().SearchOracleDatabaseUsedLicenses("", "", false, -1, -1, "", "", utils.MAX_TIME).
			Return(nil, aerrMock),
	)

	_, err := as.GetOracleDatabaseContracts(dto.NewGetOracleDatabaseContractsFilter())
	require.Equal(t, aerrMock, err)
}

func TestCheckOracleDatabaseContractMatchFilter(t *testing.T) {
	agg1 := dto.OracleDatabaseContractFE{
		ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
		ContractID:               "5051863",
		CSI:                      "6871235",
		LicenseTypeID:            "A90620",
		ItemDescription:          "Oracle Partitioning",
		Metric:                   model.LicenseTypeMetricProcessorPerpetual,
		ReferenceNumber:          "10032246681",
		Unlimited:                false,
		Basket:                   true,
		Restricted:               false,
		Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: -1, Hostname: "test-db", TotalCoveredLicensesCount: -1}, {CoveredLicensesCount: -1, Hostname: "ercsoldbx", TotalCoveredLicensesCount: -1}},
		LicensesPerCore:          30,
		LicensesPerUser:          5,
		AvailableLicensesPerCore: 7,
		AvailableLicensesPerUser: 0,
	}

	assert.True(t, checkOracleDatabaseContractMatchFilter(agg1, dto.NewGetOracleDatabaseContractsFilter()))

	assert.True(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "5051",
		LicenseTypeID:               "A9062",
		ItemDescription:             "Partitioning",
		CSI:                         "6871",
		Metric:                      model.LicenseTypeMetricProcessorPerpetual,
		ReferenceNumber:             "100322",
		Unlimited:                   "false",
		Basket:                      "true",
		LicensesPerCoreLTE:          35,
		LicensesPerCoreGTE:          25,
		LicensesPerUserLTE:          10,
		LicensesPerUserGTE:          0,
		AvailableLicensesPerCoreLTE: 8,
		AvailableLicensesPerCoreGTE: 6,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.True(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          30,
		LicensesPerCoreGTE:          30,
		LicensesPerUserLTE:          5,
		LicensesPerUserGTE:          5,
		AvailableLicensesPerCoreLTE: 7,
		AvailableLicensesPerCoreGTE: 7,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))

	assert.False(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "fdgdfgsdsfg",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "fdgdfgsdsfg",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "fdgdfgsdsfg",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "fdgdfgsdsfg",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "fdgdfgsdsfg",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "fdgdfgsdsfg",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "true",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "false",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          35,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          25,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          0,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          10,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: 3,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
	assert.False(t, checkOracleDatabaseContractMatchFilter(agg1, dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: 8,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}))
}

func TestAssignOracleDatabaseContractsToHosts_SimpleUnlimitedCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	contracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{Hostname: "test-db", CoveredLicensesCount: 0, TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 7,
			AvailableLicensesPerUser: 0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseTypeID: "PID002",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 3, Hostname: "test-db", TotalCoveredLicensesCount: 3, ConsumedLicensesCount: 3}},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
			CoveredLicenses:          3,
		},
	}

	err := as.assignOracleDatabaseContractsToHosts(contracts, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedContracts, contracts)
}

func TestAssignOracleDatabaseContractsToHosts_SimpleProcessorPerpetualCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	contracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 5,
			AvailableLicensesPerUser: 0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseTypeID: "PID002",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "AID001",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 3, Hostname: "test-db", TotalCoveredLicensesCount: 3, ConsumedLicensesCount: 3}},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 2,
			AvailableLicensesPerUser: 0,
			CoveredLicenses:          3,
		},
	}

	err := as.assignOracleDatabaseContractsToHosts(contracts, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedContracts, contracts)
}

func TestAssignOracleDatabaseContractsToHosts_SimpleNamedUserPlusCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	contracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricNamedUserPlusPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          0,
			LicensesPerUser:          250,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 250,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  5,
			LicenseTypeID: "PID002",
			OriginalCount: 5,
			Type:          "host",
		},
	}

	expectedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricNamedUserPlusPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 5, Hostname: "test-db", TotalCoveredLicensesCount: 5, ConsumedLicensesCount: 5}},
			LicensesPerCore:          0,
			LicensesPerUser:          250,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 245,
			CoveredLicenses:          5,
		},
	}

	err := as.assignOracleDatabaseContractsToHosts(contracts, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedContracts, contracts)
}

func TestAssignOracleDatabaseContractsToHosts_SharedContract(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	contracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}, {CoveredLicensesCount: 0, Hostname: "test-db2", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 5,
			AvailableLicensesPerUser: 0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseTypeID: "PID002",
			OriginalCount: 3,
			Type:          "host",
		},
		{
			Name:          "test-db2",
			LicenseCount:  4,
			LicenseTypeID: "PID002",
			OriginalCount: 4,
			Type:          "host",
		},
	}

	expectedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 4, Hostname: "test-db2", TotalCoveredLicensesCount: 4, ConsumedLicensesCount: 4}, {CoveredLicensesCount: 1, Hostname: "test-db", TotalCoveredLicensesCount: 1, ConsumedLicensesCount: 3}},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
			CoveredLicenses:          5,
		},
	}

	err := as.assignOracleDatabaseContractsToHosts(contracts, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedContracts, contracts)
}

func TestAssignOracleDatabaseContractsToHosts_SharedHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	contracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 5,
			AvailableLicensesPerUser: 0,
		},
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 0, Hostname: "test-db", TotalCoveredLicensesCount: 0}},
			LicensesPerCore:          10,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 10,
			AvailableLicensesPerUser: 0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  20,
			LicenseTypeID: "PID002",
			OriginalCount: 20,
			Type:          "host",
		},
	}

	expectedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 10, Hostname: "test-db", TotalCoveredLicensesCount: 15, ConsumedLicensesCount: 20}},
			LicensesPerCore:          10,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
			CoveredLicenses:          10,
		},
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   false,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{CoveredLicensesCount: 5, Hostname: "test-db", TotalCoveredLicensesCount: 15, ConsumedLicensesCount: 20}},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
			CoveredLicenses:          5,
		},
	}

	err := as.assignOracleDatabaseContractsToHosts(contracts, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedContracts, contracts)
}

func TestAssignOracleDatabaseContractsToHosts_SimpleUnlimitedCaseNoAssociatedHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	contracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			Basket:                   true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseTypeID: "ID Partioning",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                true,
			Basket:                   true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{},
			LicensesPerCore:          0,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
		},
	}

	err := as.assignOracleDatabaseContractsToHosts(contracts, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedContracts, contracts)
}

func TestAssignOracleDatabaseContractsToHosts_SimpleProcessorPerpetualCaseNoAssociatedHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	contracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 5,
			AvailableLicensesPerUser: 0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseTypeID: "PID002",
			OriginalCount: 3,
			Type:          "host",
		},
	}

	expectedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{},
			LicensesPerCore:          5,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 2,
			AvailableLicensesPerUser: 0,
			CoveredLicenses:          3,
		},
	}

	err := as.assignOracleDatabaseContractsToHosts(contracts, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedContracts, contracts)
}

func TestAssignOracleDatabaseContractsToHosts_SimpleNamedUserPlusCaseNoAssociatedHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	contracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricNamedUserPlusPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{},
			LicensesPerCore:          0,
			LicensesPerUser:          200,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 200,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  5,
			LicenseTypeID: "PID002",
			OriginalCount: 5,
			Type:          "host",
		},
	}

	expectedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricNamedUserPlusPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{},
			LicensesPerCore:          0,
			LicensesPerUser:          200,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 200,
			CoveredLicenses:          0,
		},
	}

	err := as.assignOracleDatabaseContractsToHosts(contracts, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedContracts, contracts)
}

func TestAssignOracleDatabaseContractsToHosts_CompleCase1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		Config:      config.Configuration{},
		NewObjectID: utils.NewObjectIDForTests(),
	}

	parts := []model.OracleDatabaseLicenseType{
		{
			ID:              "PID002",
			Aliases:         []string{"Partitioning"},
			ItemDescription: "Oracle Partitioning",
			Metric:          model.LicenseTypeMetricProcessorPerpetual,
		},
	}
	db.EXPECT().GetOracleDatabaseLicenseTypes().
		Return(parts, nil)

	cluster := dto.Cluster{
		ID:                          [12]byte{},
		CreatedAt:                   time.Time{},
		Hostname:                    "dbclust",
		HostnameAgentVirtualization: "",
		Name:                        "dbclust",
		Environment:                 "",
		Location:                    "",
		FetchEndpoint:               "",
		CPU:                         0,
		Sockets:                     0,
		Type:                        "",
		VirtualizationNodes:         []string{},
		VirtualizationNodesCount:    0,
		VirtualizationNodesStats:    []dto.VirtualizationNodesStat{},
		VMs:                         []dto.VM{},
		VMsCount:                    0,
		VMsErcoleAgentCount:         0,
	}
	db.EXPECT().GetCluster("dbclust", utils.MAX_TIME).Return(&cluster, nil)

	contracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{Hostname: "test-db"}},
			LicensesPerCore:          10,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 10,
			AvailableLicensesPerUser: 0,
		},
	}
	hosts := []dto.HostUsingOracleDatabaseLicenses{
		{
			Name:          "test-db",
			LicenseCount:  3,
			LicenseTypeID: "PID002",
			OriginalCount: 3,
			Type:          "host",
		},
		{
			Name:          "dbclust",
			LicenseCount:  20,
			LicenseTypeID: "PID002",
			OriginalCount: 20,
			Type:          "cluster",
		},
	}

	expectedContracts := []dto.OracleDatabaseContractFE{
		{
			ID:                       utils.Str2oid("5f4d0ab1c6bc19e711bbcce6"),
			ContractID:               "5051863",
			CSI:                      "CSI001",
			LicenseTypeID:            "PID002",
			ItemDescription:          "Oracle Partitioning",
			Metric:                   model.LicenseTypeMetricProcessorPerpetual,
			ReferenceNumber:          "RF0001",
			Unlimited:                false,
			Basket:                   true,
			Restricted:               false,
			Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{Hostname: "test-db", CoveredLicensesCount: 3, TotalCoveredLicensesCount: 3, ConsumedLicensesCount: 3}},
			LicensesPerCore:          10,
			LicensesPerUser:          0,
			AvailableLicensesPerCore: 0,
			AvailableLicensesPerUser: 0,
			CoveredLicenses:          10,
		},
	}

	err := as.assignOracleDatabaseContractsToHosts(contracts, hosts)
	assert.NoError(t, err)

	assert.Equal(t, expectedContracts, contracts)
}

func TestSortHostsUsingLicenses(t *testing.T) {
	list := []dto.HostUsingOracleDatabaseLicenses{
		{
			LicenseTypeID: "Diagnostics Pack",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
		},
		{
			LicenseTypeID: "Real Application Clusters",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  1.5,
		},
		{
			LicenseTypeID: "Diagnostics Pack",
			Name:          "test-db4",
			Type:          "host",
			LicenseCount:  0.5,
		},
		{
			LicenseTypeID: "Oracle ENT",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  0.5,
		},
		{
			LicenseTypeID: "Oracle ENT",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
		},
	}

	expected := []dto.HostUsingOracleDatabaseLicenses{
		{
			LicenseTypeID: "Oracle ENT",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
		},
		{
			LicenseTypeID: "Diagnostics Pack",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
		},
		{
			LicenseTypeID: "Real Application Clusters",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  1.5,
		},
		{
			LicenseTypeID: "Diagnostics Pack",
			Name:          "test-db4",
			Type:          "host",
			LicenseCount:  0.5,
		},
		{
			LicenseTypeID: "Oracle ENT",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  0.5,
		},
	}

	sortHostsByLicenses(list)

	assert.Equal(t, expected, list)
}

func TestSortOracleDatabaseContracts(t *testing.T) {
	list := []dto.OracleDatabaseContractFE{
		{Basket: true, Unlimited: false, LicensesPerUser: 10},
		{Basket: true, Unlimited: false, LicensesPerCore: 10},
		{Basket: true, Unlimited: true, LicensesPerUser: 20},
		{Basket: false, Unlimited: false, LicensesPerCore: 20},
		{Basket: false, Unlimited: true, LicensesPerUser: 10},
		{Basket: false, Unlimited: true, LicensesPerCore: 20},
		{Basket: false, Unlimited: false, LicensesPerCore: 10},
		{Basket: true, Unlimited: true, LicensesPerCore: 10},
		{Basket: false, Unlimited: true, LicensesPerUser: 20},
		{Basket: false, Unlimited: false, LicensesPerUser: 10},
		{Basket: true, Unlimited: true, LicensesPerUser: 10},
		{Basket: true, Unlimited: true, LicensesPerCore: 20},
		{Basket: true, Unlimited: false, LicensesPerCore: 20},
		{Basket: false, Unlimited: false, LicensesPerUser: 20},
		{Basket: false, Unlimited: true, LicensesPerCore: 10},
		{Basket: true, Unlimited: false, LicensesPerUser: 20},
	}

	expected := []dto.OracleDatabaseContractFE{
		{Basket: false, Unlimited: false, LicensesPerUser: 20},
		{Basket: false, Unlimited: false, LicensesPerUser: 10},
		{Basket: false, Unlimited: false, LicensesPerCore: 20},
		{Basket: false, Unlimited: false, LicensesPerCore: 10},
		{Basket: false, Unlimited: true, LicensesPerUser: 20},
		{Basket: false, Unlimited: true, LicensesPerUser: 10},
		{Basket: false, Unlimited: true, LicensesPerCore: 20},
		{Basket: false, Unlimited: true, LicensesPerCore: 10},
		{Basket: true, Unlimited: false, LicensesPerUser: 20},
		{Basket: true, Unlimited: false, LicensesPerUser: 10},
		{Basket: true, Unlimited: false, LicensesPerCore: 20},
		{Basket: true, Unlimited: false, LicensesPerCore: 10},
		{Basket: true, Unlimited: true, LicensesPerUser: 20},
		{Basket: true, Unlimited: true, LicensesPerUser: 10},
		{Basket: true, Unlimited: true, LicensesPerCore: 20},
		{Basket: true, Unlimited: true, LicensesPerCore: 10},
	}

	sortOracleDatabaseContracts(list)

	assert.Equal(t, expected, list)
}

func TestSortAssociatedHostsInOracleDatabaseContract(t *testing.T) {
	hostsMap := map[string]map[string]*dto.HostUsingOracleDatabaseLicenses{
		"L10005": {
			"test-db1": {
				LicenseCount: 30,
			},
			"test-db2": {
				LicenseCount: 30,
			},
			"test-db3": {
				LicenseCount: 15,
			},
			"test-db4": {
				LicenseCount: 35,
			},
		},
	}

	agr := dto.OracleDatabaseContractFE{
		LicenseTypeID: "L10005",
		Hosts: []dto.OracleDatabaseContractAssociatedHostFE{
			{Hostname: "test-db2"},
			{Hostname: "test-db1"},
			{Hostname: "test-db4"},
			{Hostname: "test-db3"},
		},
	}

	expected := []dto.OracleDatabaseContractAssociatedHostFE{
		{Hostname: "test-db4"},
		{Hostname: "test-db2"},
		{Hostname: "test-db1"},
		{Hostname: "test-db3"},
	}

	sortHostsInContractByLicenseCount(&agr, hostsMap)

	assert.Equal(t, expected, agr.Hosts)
}

func TestBuildHostUsingLicensesMap(t *testing.T) {
	list := []dto.HostUsingOracleDatabaseLicenses{
		{
			LicenseTypeID: "LTID01",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
		},
		{
			LicenseTypeID: "LTID02",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
		},
		{
			LicenseTypeID: "LTID03",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  1.5,
		},
		{
			LicenseTypeID: "LTID02",
			Name:          "test-db4",
			Type:          "host",
			LicenseCount:  0.5,
		},
		{
			LicenseTypeID: "LTID01",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  0.5,
		},
	}

	expected := map[string]map[string]*dto.HostUsingOracleDatabaseLicenses{
		"LTID01": {
			"Puzzait":  &list[0],
			"test-db3": &list[4],
		},
		"LTID02": {
			"Puzzait":  &list[1],
			"test-db4": &list[3],
		},
		"LTID03": {
			"test-db3": &list[2],
		},
	}

	assert.Equal(t, expected, buildHostUsingLicensesMap(list))
}

func TestBuildContractPartMap(t *testing.T) {
	list := []model.OracleDatabaseLicenseType{
		{
			ItemDescription: "itemDesc1",
			Aliases:         []string{"alias1"},
			Metric:          "metric1",
			ID:              "PID001",
		},
		{
			ItemDescription: "itemDesc2",
			Aliases:         []string{"alias1"},
			Metric:          "metric2",
			ID:              "PID002",
		},
	}

	expected := map[string]*model.OracleDatabaseLicenseType{
		"PID001": &list[0],
		"PID002": &list[1],
	}

	assert.Equal(t, expected, buildLicenseTypesMap(list))
}

func TestDeleteOracleDatabaseContract(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	contractID := utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")

	t.Run("Fail: can't find contract", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().RemoveOracleDatabaseContract(contractID).
				Return(utils.ErrContractNotFound),
		)

		err := as.DeleteOracleDatabaseContract(contractID)
		require.EqualError(t, err, utils.ErrContractNotFound.Error())
	})

	t.Run("Success", func(t *testing.T) {
		contract := model.OracleDatabaseContract{
			ID:              contractID,
			ContractID:      "AID001",
			CSI:             "CSI001",
			LicenseTypeID:   lt1.ID,
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			Basket:          true,
			Hosts:           []string{"test-db", "ercsoldbx"},
		}

		gomock.InOrder(
			db.EXPECT().RemoveOracleDatabaseContract(contract.ID).
				Return(nil),
		)

		err := as.DeleteOracleDatabaseContract(contractID)
		assert.Nil(t, err)
	})

}

func TestAddHostToOracleDatabaseContract(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	commonFilters := dto.NewSearchHostsFilters()
	anotherAssociatedPartID := utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb")
	t.Run("Fail: can't find host", func(t *testing.T) {

		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames",
				commonFilters).Return([]map[string]interface{}{
				{"hostname": "test-db"},
				{"hostname": "foobar"},
				{"hostname": "ercsoldbx"},
			}, nil),
		)

		err := as.AddHostToOracleDatabaseContract(anotherAssociatedPartID, "pippo")
		assert.EqualError(t, err, utils.ErrHostNotFound.Error())
	})

	id := utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")

	t.Run("Fail: can't find contract", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames",
				commonFilters).
				Return([]map[string]interface{}{
					{"hostname": "test-db"},
					{"hostname": "foobar"},
					{"hostname": "ercsoldbx"},
				}, nil),
			db.EXPECT().GetOracleDatabaseContract(id).
				Return(nil, utils.ErrContractNotFound),
		)

		err := as.AddHostToOracleDatabaseContract(id, "foobar")
		assert.EqualError(t, err, utils.ErrContractNotFound.Error())
	})

	t.Run("Success", func(t *testing.T) {
		contract := model.OracleDatabaseContract{
			ContractID:      "AID001",
			CSI:             "CSI001",
			ID:              id,
			LicenseTypeID:   lt1.ID,
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			Basket:          true,
			Hosts:           []string{"test-db", "ercsoldbx"},
		}

		contractPostAdd := model.OracleDatabaseContract{
			ID:              id,
			ContractID:      "AID001",
			CSI:             "CSI001",
			LicenseTypeID:   lt1.ID,
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			Basket:          true,
			Hosts:           []string{"test-db", "ercsoldbx", "foobar"},
		}

		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames",
				commonFilters).
				Return([]map[string]interface{}{
					{"hostname": "test-db"},
					{"hostname": "foobar"},
					{"hostname": "ercsoldbx"},
				}, nil),
			db.EXPECT().GetOracleDatabaseContract(id).
				Return(&contract, nil),
			db.EXPECT().UpdateOracleDatabaseContract(contractPostAdd).
				Return(nil),
		)

		err := as.AddHostToOracleDatabaseContract(id, "foobar")
		assert.Nil(t, err)
	})
}

func TestDeleteHostFromOracleDatabaseContract(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	commonFilters := dto.NewSearchHostsFilters()
	id := utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb")

	t.Run("Fail: can't get contract", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames", commonFilters).
				Return([]map[string]interface{}{
					{"hostname": "pippo"},
				}, nil),
			db.EXPECT().GetOracleDatabaseContract(id).
				Return(nil, utils.ErrContractNotFound),
		)

		err := as.DeleteHostFromOracleDatabaseContract(id, "pippo")
		require.EqualError(t, err, utils.ErrContractNotFound.Error())
	})

	anotherId := utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")

	t.Run("Success", func(t *testing.T) {
		contract := model.OracleDatabaseContract{
			ContractID:      "AID001",
			CSI:             "CSI001",
			ID:              anotherId,
			LicenseTypeID:   lt1.ID,
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			Basket:          true,
			Hosts:           []string{"test-db", "ercsoldbx"},
		}

		contractPostAdd := model.OracleDatabaseContract{
			ContractID:      "AID001",
			CSI:             "CSI001",
			ID:              anotherId,
			LicenseTypeID:   lt1.ID,
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			Basket:          true,
			Hosts:           []string{"test-db"},
		}

		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames",
				commonFilters).
				Return([]map[string]interface{}{
					{"hostname": "ercsoldbx"},
				}, nil),
			db.EXPECT().GetOracleDatabaseContract(anotherId).
				Return(&contract, nil),
			db.EXPECT().UpdateOracleDatabaseContract(contractPostAdd).
				Return(nil),
		)

		err := as.DeleteHostFromOracleDatabaseContract(anotherId, "ercsoldbx")
		assert.Nil(t, err)
	})
}

func TestDeleteHostFromOracleDatabaseContracts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	anotherId := utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa")
	commonFilters := dto.NewSearchHostsFilters()

	t.Run("Success", func(t *testing.T) {
		listContract := []dto.OracleDatabaseContractFE{}

		contract := model.OracleDatabaseContract{
			ContractID:      "AID001",
			CSI:             "CSI001",
			ID:              anotherId,
			LicenseTypeID:   lt1.ID,
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			Basket:          true,
			Hosts:           []string{"test-db", "ercsoldbx"},
		}

		contractPostAdd := model.OracleDatabaseContract{
			ContractID:      "AID001",
			CSI:             "CSI001",
			ID:              anotherId,
			LicenseTypeID:   lt1.ID,
			ReferenceNumber: "RF0001",
			Unlimited:       true,
			Count:           30,
			Basket:          true,
			Hosts:           []string{"test-db"},
		}

		gomock.InOrder(
			db.EXPECT().SearchHosts("hostnames",
				commonFilters).
				Return([]map[string]interface{}{
					{"hostname": "ercsoldbx"},
				}, nil),
			db.EXPECT().ListOracleDatabaseContracts(gomock.Any()).
				Return(listContract, nil),
			db.EXPECT().GetOracleDatabaseContract(anotherId).
				Return(&contract, nil).AnyTimes(),
			db.EXPECT().UpdateOracleDatabaseContract(contractPostAdd).
				Return(nil).AnyTimes(),
		)

		err := as.DeleteHostFromOracleDatabaseContracts("ercsoldbx")
		assert.Nil(t, err)
	})
}

func TestGetOracleDatabaseContractsAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	contract := model.OracleDatabaseContract{
		ContractID:      "5051863",
		CSI:             "13902248",
		ID:              utils.Str2oid("609ce3072eff5d5540ec4a28"),
		LicenseTypeID:   lt1.ID,
		ReferenceNumber: "37255828",
		Unlimited:       false,
		Count:           30,
		Basket:          false,
		Restricted:      false,
		Hosts:           []string{"test-db", "ercsoldbx"},
	}

	searchedContractItem := dto.OracleDatabaseContractFE{
		ID:                       contract.ID,
		ContractID:               contract.ContractID,
		CSI:                      contract.CSI,
		LicenseTypeID:            contract.LicenseTypeID,
		ItemDescription:          "Oracle Database Enterprise Edition",
		Metric:                   "Named User Plus Perpetual",
		ReferenceNumber:          contract.ReferenceNumber,
		Unlimited:                false,
		Basket:                   false,
		Restricted:               false,
		Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{},
		LicensesPerCore:          0,
		LicensesPerUser:          350,
		AvailableLicensesPerCore: 0,
		AvailableLicensesPerUser: 0,
	}
	as.mockGetOracleDatabaseContracts = func(filters dto.GetOracleDatabaseContractsFilter) ([]dto.OracleDatabaseContractFE, error) {
		return []dto.OracleDatabaseContractFE{searchedContractItem}, nil
	}

	filter := dto.GetOracleDatabaseContractsFilter{
		ContractID:                  "",
		LicenseTypeID:               "",
		ItemDescription:             "",
		CSI:                         "",
		Metric:                      "",
		ReferenceNumber:             "",
		Unlimited:                   "true",
		Basket:                      "",
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}

	actual, err := as.GetOracleDatabaseContractsAsXLSX(filter)
	require.NoError(t, err)

	assert.Equal(t, "5051863", actual.GetCellValue("Contracts", "A2"))
	assert.Equal(t, "PID001", actual.GetCellValue("Contracts", "B2"))
	assert.Equal(t, "Oracle Database Enterprise Edition", actual.GetCellValue("Contracts", "C2"))
	assert.Equal(t, "Named User Plus Perpetual", actual.GetCellValue("Contracts", "D2"))
	assert.Equal(t, "13902248", actual.GetCellValue("Contracts", "E2"))
	assert.Equal(t, "37255828", actual.GetCellValue("Contracts", "F2"))

	assert.Equal(t, "", actual.GetCellValue("Contracts", "H2"))
	assert.Equal(t, "", actual.GetCellValue("Contracts", "I2"))
	assert.Equal(t, "", actual.GetCellValue("Contracts", "J2"))

	assert.Equal(t, "0", actual.GetCellValue("Contracts", "K2"))
	assert.Equal(t, "0", actual.GetCellValue("Contracts", "L2"))
	assert.Equal(t, "350", actual.GetCellValue("Contracts", "M2"))
	assert.Equal(t, "0", actual.GetCellValue("Contracts", "N2"))
	assert.Equal(t, "0", actual.GetCellValue("Contracts", "O2"))
}
