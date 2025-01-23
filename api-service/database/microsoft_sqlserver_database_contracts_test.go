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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

var sqlServerLicenseTypeSample = model.SqlServerDatabaseLicenseType{
	ID:              "359-06320",
	ItemDescription: "ItemDesc 1",
	Edition:         "ED00001",
	Version:         "V 0.0.1",
}

var msContractSample = model.SqlServerDatabaseContract{
	ID:             utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
	ContractID:     "AID001",
	LicenseTypeID:  "359-06320",
	Type:           "TYPE001",
	LicensesNumber: 1,
	Clusters:       []string{"foo", "bar"},
	Hosts:          []string{"foo", "bar"},
}

func (m *MongodbSuite) TestInsertSqlServerDatabaseContract_Success() {
	_, aerr := m.db.InsertSqlServerDatabaseContract(msContractSample)
	require.NoError(m.T(), aerr)
	defer m.db.Client.Database(m.dbname).Collection(sqlServerDbContractsCollection).DeleteMany(context.TODO(), bson.M{})

	val := m.db.Client.Database(m.dbname).Collection(sqlServerDbContractsCollection).FindOne(context.TODO(), bson.M{
		"_id": msContractSample.ID,
	})
	require.NoError(m.T(), val.Err())

	var out model.SqlServerDatabaseContract
	err := val.Decode(&out)
	assert.NoError(m.T(), err)

	assert.Equal(m.T(), msContractSample, out)
}

func (m *MongodbSuite) TestInsertSqlServerDatabaseContract_DuplicateError() {
	defer m.db.Client.Database(m.dbname).Collection(sqlServerDbContractsCollection).DeleteMany(context.TODO(), bson.M{})

	_, err := m.db.InsertSqlServerDatabaseContract(msContractSample)
	require.NoError(m.T(), err)

	_, err = m.db.InsertSqlServerDatabaseContract(msContractSample)
	require.Error(m.T(), err, "Should not accept two contracts with same ID")
}

func (m *MongodbSuite) TestUpdateSqlServerDatabaseContract() {
	defer m.db.Client.Database(m.dbname).Collection(sqlServerDbContractsCollection).DeleteMany(context.TODO(), bson.M{})

	_, err := m.db.InsertSqlServerDatabaseContract(msContractSample)
	require.NoError(m.T(), err)

	m.T().Run("id_exist", func(t *testing.T) {
		contractSampleUpdated := model.SqlServerDatabaseContract{
			ID:            utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			ContractID:    "AID001",
			LicenseTypeID: "359-06320",
			Hosts:         []string{"foo", "bar"},
		}

		err := m.db.UpdateSqlServerDatabaseContract(contractSampleUpdated)
		require.NoError(t, err)

		val := m.db.Client.Database(m.dbname).Collection(sqlServerDbContractsCollection).FindOne(context.TODO(), bson.M{
			"_id": contractSampleUpdated.ID,
		})
		require.NoError(m.T(), val.Err())

		var out model.SqlServerDatabaseContract
		err2 := val.Decode(&out)
		assert.NoError(t, err2)

		assert.Equal(m.T(), contractSampleUpdated, out)
	})

	m.T().Run("id_not_exist", func(t *testing.T) {
		contractSampleUpdated := model.OracleDatabaseContract{
			ID: utils.Str2oid("doesn't exist"),
		}
		err := m.db.UpdateOracleDatabaseContract(contractSampleUpdated)

		require.Equal(t, utils.ErrContractNotFound, err)
	})
}

func (m *MongodbSuite) TestRemoveSqlServerDatabaseContract() {
	defer m.db.Client.Database(m.dbname).Collection(sqlServerDbContractsCollection).DeleteMany(context.TODO(), bson.M{})

	_, err := m.db.InsertSqlServerDatabaseContract(msContractSample)
	require.NoError(m.T(), err)

	err = m.db.RemoveSqlServerDatabaseContract(msContractSample.ID)
	require.NoError(m.T(), err)

	err = m.db.RemoveSqlServerDatabaseContract(utils.Str2oid("5dcad8933b243f80e2ed8538"))
	require.Equal(m.T(), utils.ErrContractNotFound, err)
}

func (m *MongodbSuite) TestListSqlServerDatabaseContracts() {
	defer m.db.Client.Database(m.dbname).Collection(sqlServerDbLicenseTypesCollection).DeleteMany(context.TODO(), bson.M{})
	licenseTypeSample1 := model.SqlServerDatabaseLicenseType{
		ID:              "359-06320",
		ItemDescription: "ItemDesc 1",
		Edition:         "ED00001",
		Version:         "V 0.0.1",
	}
	licenseTypeSample2 := model.SqlServerDatabaseLicenseType{
		ID:              "ID00002",
		ItemDescription: "ItemDesc 2",
		Edition:         "ED00002",
		Version:         "V 0.0.2",
	}
	_, err2 := m.db.Client.Database(m.dbname).Collection(sqlServerDbLicenseTypesCollection).
		InsertMany(context.TODO(), []interface{}{licenseTypeSample1, licenseTypeSample2})
	require.NoError(m.T(), err2)

	contractSample := model.SqlServerDatabaseContract{
		ID:             utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		ContractID:     "AID001",
		LicenseTypeID:  "359-06320",
		Type:           "TYPE001",
		LicensesNumber: 1,
		Clusters:       []string{"foo", "bar"},
		Hosts:          []string{"foo", "bar"},
	}

	m.T().Run("Empty collection", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(sqlServerDbContractsCollection).DeleteMany(context.TODO(), bson.M{})

		out, err := m.db.ListSqlServerDatabaseContracts([]string{})
		m.Require().NoError(err)

		assert.Equal(m.T(), []model.SqlServerDatabaseContract{}, out)
	})

	m.T().Run("One contract", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(sqlServerDbContractsCollection).DeleteMany(context.TODO(), bson.M{})
		_, err := m.db.InsertSqlServerDatabaseContract(contractSample)
		require.NoError(m.T(), err)

		out, err := m.db.ListSqlServerDatabaseContracts([]string{})
		m.Require().NoError(err)

		assert.Equal(m.T(), []model.SqlServerDatabaseContract{
			{
				ID:             utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				ContractID:     "AID001",
				LicenseTypeID:  "359-06320",
				Type:           "TYPE001",
				LicensesNumber: 1,
				Clusters:       []string{"foo", "bar"},
				Hosts:          []string{"foo", "bar"},
			},
		}, out)
	})
}
