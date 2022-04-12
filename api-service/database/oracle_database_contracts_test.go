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

package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

var licenseTypeSample = model.OracleDatabaseLicenseType{
	ID:              "ID00001",
	ItemDescription: "ItemDesc 1",
	Cost:            42,
	Metric:          model.LicenseTypeMetricProcessorPerpetual,
	Aliases:         []string{},
}

var contractSample = model.OracleDatabaseContract{
	ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
	ContractID:      "AID001",
	CSI:             "csi001",
	LicenseTypeID:   licenseTypeSample.ID,
	ReferenceNumber: "R00001",
	Unlimited:       true,
	Count:           345,
	Basket:          true,
	Restricted:      true, // there shouldn't be Basket==true && Restricted && true, this is only for tests
	Hosts:           []string{"foo", "bar"},
}

func (m *MongodbSuite) TestInsertOracleDatabaseContract_Success() {
	aerr := m.db.InsertOracleDatabaseContract(contractSample)
	require.NoError(m.T(), aerr)
	defer m.db.Client.Database(m.dbname).Collection(oracleDbContractsCollection).DeleteMany(context.TODO(), bson.M{})

	val := m.db.Client.Database(m.dbname).Collection(oracleDbContractsCollection).FindOne(context.TODO(), bson.M{
		"_id": contractSample.ID,
	})
	require.NoError(m.T(), val.Err())

	var out model.OracleDatabaseContract
	err := val.Decode(&out)
	assert.NoError(m.T(), err)

	assert.Equal(m.T(), contractSample, out)
}

func (m *MongodbSuite) TestInsertOracleDatabaseContract_DuplicateError() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbContractsCollection).DeleteMany(context.TODO(), bson.M{})

	err := m.db.InsertOracleDatabaseContract(contractSample)
	require.NoError(m.T(), err)

	err = m.db.InsertOracleDatabaseContract(contractSample)
	require.Error(m.T(), err, "Should not accept two contracts with same ID")
}

func (m *MongodbSuite) TestGetOracleDatabaseContract() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbContractsCollection).DeleteMany(context.TODO(), bson.M{})

	err := m.db.InsertOracleDatabaseContract(contractSample)
	require.NoError(m.T(), err)

	m.T().Run("id_exist", func(t *testing.T) {
		out, err := m.db.GetOracleDatabaseContract(contractSample.ID)
		require.NoError(t, err)
		assert.Equal(t, contractSample, *out)
	})

	m.T().Run("id_not_exist", func(t *testing.T) {
		out, err := m.db.GetOracleDatabaseContract(utils.Str2oid("xxxxxxxxxxxxxxxxxxxxxxxx"))
		require.Nil(t, out)
		require.Equal(t, utils.ErrOracleDatabaseContractNotFound, err)
	})
}

func (m *MongodbSuite) TestUpdateOracleDatabaseContract() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbContractsCollection).DeleteMany(context.TODO(), bson.M{})

	err := m.db.InsertOracleDatabaseContract(contractSample)
	require.NoError(m.T(), err)

	m.T().Run("id_exist", func(t *testing.T) {
		contractSampleUpdated := model.OracleDatabaseContract{
			ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			ContractID:      "AID001",
			CSI:             "000001",
			LicenseTypeID:   licenseTypeSample.ID,
			ReferenceNumber: "000002",
			Unlimited:       true,
			Count:           345,
			Basket:          true,
			Restricted:      false,
			Hosts:           []string{"foo", "bar"},
		}

		err := m.db.UpdateOracleDatabaseContract(contractSampleUpdated)
		require.NoError(t, err)

		val := m.db.Client.Database(m.dbname).Collection(oracleDbContractsCollection).FindOne(context.TODO(), bson.M{
			"_id": contractSampleUpdated.ID,
		})
		require.NoError(m.T(), val.Err())

		var out model.OracleDatabaseContract
		err2 := val.Decode(&out)
		assert.NoError(t, err2)

		assert.Equal(m.T(), contractSampleUpdated, out)
	})

	m.T().Run("id_not_exist", func(t *testing.T) {
		contractSampleUpdated := model.OracleDatabaseContract{
			ID: utils.Str2oid("doesn't exist"),
		}
		err := m.db.UpdateOracleDatabaseContract(contractSampleUpdated)

		require.Equal(t, utils.ErrOracleDatabaseContractNotFound, err)
	})
}

func (m *MongodbSuite) TestRemoveOracleDatabaseContract() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbContractsCollection).DeleteMany(context.TODO(), bson.M{})

	err := m.db.InsertOracleDatabaseContract(contractSample)
	require.NoError(m.T(), err)

	err = m.db.RemoveOracleDatabaseContract(contractSample.ID)
	require.NoError(m.T(), err)

	err = m.db.RemoveOracleDatabaseContract(utils.Str2oid("5dcad8933b243f80e2ed8538"))
	require.Equal(m.T(), utils.ErrOracleDatabaseContractNotFound, err)
}

func (m *MongodbSuite) TestListOracleDatabaseContracts() {
	defer m.db.Client.Database(m.dbname).Collection("oracle_database_license_types").DeleteMany(context.TODO(), bson.M{})
	licenseTypeSample1 := model.OracleDatabaseLicenseType{
		ID:              "ID00001",
		ItemDescription: "ItemDesc 1",
		Cost:            42,
		Metric:          model.LicenseTypeMetricProcessorPerpetual,
		Aliases:         []string{},
	}
	licenseTypeSample2 := model.OracleDatabaseLicenseType{
		ID:              "ID00002",
		ItemDescription: "ItemDesc 2",
		Cost:            24,
		Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
		Aliases:         []string{},
	}
	_, err2 := m.db.Client.Database(m.dbname).Collection(oracleDbLicenseTypesCollection).
		InsertMany(context.TODO(), []interface{}{licenseTypeSample1, licenseTypeSample2})
	require.NoError(m.T(), err2)

	contractSample := model.OracleDatabaseContract{
		ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		ContractID:      "agr001",
		CSI:             "csi001",
		LicenseTypeID:   licenseTypeSample1.ID,
		ReferenceNumber: "R00001",
		Basket:          true,
		Restricted:      false,
		Count:           345,
		Hosts:           []string{"foo", "bar"},
		Unlimited:       true,
	}
	contractSample2 := model.OracleDatabaseContract{
		ID:              utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
		ContractID:      "agr002",
		CSI:             "csi002",
		LicenseTypeID:   licenseTypeSample1.ID,
		ReferenceNumber: "R00002",
		Basket:          false,
		Restricted:      true,
		Count:           111,
		Hosts:           []string{"pippo", "clarabella"},
		Unlimited:       false,
	}

	contractSample3 := model.OracleDatabaseContract{
		ID:              utils.Str2oid("cccccccccccccccccccccccc"),
		ContractID:      "agr002",
		CSI:             "csi002",
		LicenseTypeID:   licenseTypeSample2.ID,
		ReferenceNumber: "R00003",
		Basket:          false,
		Restricted:      false,
		Count:           222,
		Hosts:           []string{"topolino", "minni"},
		Unlimited:       true,
	}

	m.T().Run("Empty collection", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(oracleDbContractsCollection).DeleteMany(context.TODO(), bson.M{})

		out, err := m.db.ListOracleDatabaseContracts()
		m.Require().NoError(err)

		assert.Equal(m.T(), []dto.OracleDatabaseContractFE{}, out)
	})

	m.T().Run("One contract", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(oracleDbContractsCollection).DeleteMany(context.TODO(), bson.M{})
		err := m.db.InsertOracleDatabaseContract(contractSample)
		require.NoError(m.T(), err)

		out, err := m.db.ListOracleDatabaseContracts()
		m.Require().NoError(err)

		assert.Equal(m.T(), []dto.OracleDatabaseContractFE{
			{
				ID:                       utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				ContractID:               "agr001",
				CSI:                      "csi001",
				LicenseTypeID:            "ID00001",
				ItemDescription:          "ItemDesc 1",
				Metric:                   model.LicenseTypeMetricProcessorPerpetual,
				ReferenceNumber:          "R00001",
				Unlimited:                true,
				Basket:                   true,
				Restricted:               false,
				Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{Hostname: "foo"}, {Hostname: "bar"}},
				LicensesPerCore:          345,
				LicensesPerUser:          0,
				AvailableLicensesPerCore: 345,
				AvailableLicensesPerUser: 0,
			},
		}, out)
	})

	m.T().Run("Multiple associations contract-licenseTypes", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(oracleDbContractsCollection).DeleteMany(context.TODO(), bson.M{})

		err := m.db.InsertOracleDatabaseContract(contractSample)
		require.NoError(m.T(), err)
		err = m.db.InsertOracleDatabaseContract(contractSample2)
		require.NoError(m.T(), err)
		err = m.db.InsertOracleDatabaseContract(contractSample3)
		require.NoError(m.T(), err)

		out, err := m.db.ListOracleDatabaseContracts()
		m.Require().NoError(err)

		assert.Equal(m.T(), []dto.OracleDatabaseContractFE{
			{
				ID:                       utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				ContractID:               "agr001",
				CSI:                      "csi001",
				LicenseTypeID:            "ID00001",
				ItemDescription:          "ItemDesc 1",
				Metric:                   "Processor Perpetual",
				ReferenceNumber:          "R00001",
				Unlimited:                true,
				Basket:                   true,
				Restricted:               false,
				Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{Hostname: "foo", CoveredLicensesCount: 0, TotalCoveredLicensesCount: 0, ConsumedLicensesCount: 0}, {Hostname: "bar", CoveredLicensesCount: 0, TotalCoveredLicensesCount: 0, ConsumedLicensesCount: 0}},
				LicensesPerCore:          345,
				LicensesPerUser:          0,
				AvailableLicensesPerCore: 345,
				AvailableLicensesPerUser: 0,
			},
			{
				ID:                       utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
				ContractID:               "agr002",
				CSI:                      "csi002",
				LicenseTypeID:            "ID00001",
				ItemDescription:          "ItemDesc 1",
				Metric:                   "Processor Perpetual",
				ReferenceNumber:          "R00002",
				Unlimited:                false,
				Basket:                   false,
				Restricted:               true,
				Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{Hostname: "pippo", CoveredLicensesCount: 0, TotalCoveredLicensesCount: 0, ConsumedLicensesCount: 0}, {Hostname: "clarabella", CoveredLicensesCount: 0, TotalCoveredLicensesCount: 0, ConsumedLicensesCount: 0}},
				LicensesPerCore:          111,
				LicensesPerUser:          0,
				AvailableLicensesPerCore: 111,
				AvailableLicensesPerUser: 0,
			},
			{
				ID:                       utils.Str2oid("cccccccccccccccccccccccc"),
				ContractID:               "agr002",
				CSI:                      "csi002",
				LicenseTypeID:            "ID00002",
				ItemDescription:          "ItemDesc 2",
				Metric:                   "Named User Plus Perpetual",
				ReferenceNumber:          "R00003",
				Unlimited:                true,
				Basket:                   false,
				Restricted:               false,
				Hosts:                    []dto.OracleDatabaseContractAssociatedHostFE{{Hostname: "topolino", CoveredLicensesCount: 0, TotalCoveredLicensesCount: 0, ConsumedLicensesCount: 0}, {Hostname: "minni", CoveredLicensesCount: 0, TotalCoveredLicensesCount: 0, ConsumedLicensesCount: 0}},
				LicensesPerCore:          0,
				LicensesPerUser:          222,
				AvailableLicensesPerCore: 0,
				AvailableLicensesPerUser: 222,
			}},

			out)
	})
}
