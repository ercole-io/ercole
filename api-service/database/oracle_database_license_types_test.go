// Copyright (c) 2021 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.mongodb.org/mongo-driver/bson"
)

var licenseTypeExample = model.OracleDatabaseLicenseType{
	ID:              "Test",
	ItemDescription: "Oracle Database Enterprise Edition",
	Metric:          "Processor Perpetual",
	Cost:            500,
	Aliases:         []string{"Tuning Pack"},
	Option:          false,
}

func (m *MongodbSuite) TestGetOracleDatabaseLicenseTypes() {
	defer m.db.Client.Database(m.dbname).Collection("oracle_database_license_types").DeleteMany(context.TODO(), bson.M{})

	m.T().Run("success with empty table", func(t *testing.T) {
		actual, err := m.db.GetOracleDatabaseLicenseTypes()
		m.Require().NoError(err)

		expected := make([]model.OracleDatabaseLicenseType, 0)
		assert.ElementsMatch(m.T(), expected, actual)
	})

	m.T().Run("success with some values", func(t *testing.T) {
		expected := []interface{}{
			model.OracleDatabaseLicenseType{
				ID:              "PID001",
				ItemDescription: "desc001",
				Metric:          model.LicenseTypeMetricProcessorPerpetual,
				Cost:            42,
				Aliases:         []string{"pippo"},
			},
			model.OracleDatabaseLicenseType{
				ID:              "PID002",
				ItemDescription: "desc002",
				Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
				Cost:            7,
				Aliases:         []string{"qui", "quo", "qua"},
			},
		}

		ctx := context.TODO()
		_, err := m.db.Client.Database(m.dbname).
			Collection("oracle_database_license_types").
			InsertMany(ctx, expected)
		m.Require().NoError(err)

		actual, err := m.db.GetOracleDatabaseLicenseTypes()
		m.Require().NoError(err)

		assert.ElementsMatch(m.T(), expected, actual)
	})
}

func (m *MongodbSuite) TestInsertOracleDatabaseLicenseTypes_Success() {
	aerr := m.db.InsertOracleDatabaseLicenseType(licenseTypeExample)
	require.NoError(m.T(), aerr)
	defer m.db.Client.Database(m.dbname).Collection(oracleDbLicenseTypesCollection).DeleteMany(context.TODO(), bson.M{})

	val := m.db.Client.Database(m.dbname).Collection(oracleDbLicenseTypesCollection).FindOne(context.TODO(), bson.M{
		"_id": licenseTypeExample.ID,
	})
	require.NoError(m.T(), val.Err())

	var out model.OracleDatabaseLicenseType
	err := val.Decode(&out)
	assert.NoError(m.T(), err)

	assert.Equal(m.T(), licenseTypeExample, out)
}

func (m *MongodbSuite) TestUpdateOracleDatabaseLicenseTypes() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbLicenseTypesCollection).DeleteMany(context.TODO(), bson.M{})

	err := m.db.InsertOracleDatabaseLicenseType(licenseTypeExample)
	require.NoError(m.T(), err)

	m.T().Run("id_exist", func(t *testing.T) {
		licenseTypeSampleUpdated := model.OracleDatabaseLicenseType{
			ID:              "Test",
			ItemDescription: "Oracle Database Enterprise Edition",
			Metric:          "Processor Perpetual",
			Cost:            500,
			Aliases:         []string{"Tuning Pack"},
			Option:          false,
		}

		err := m.db.UpdateOracleDatabaseLicenseType(licenseTypeSampleUpdated)
		require.NoError(t, err)

		val := m.db.Client.Database(m.dbname).Collection(oracleDbLicenseTypesCollection).FindOne(context.TODO(), bson.M{
			"_id": licenseTypeSampleUpdated.ID,
		})
		require.NoError(m.T(), val.Err())

		var out model.OracleDatabaseLicenseType
		err2 := val.Decode(&out)
		assert.NoError(t, err2)

		assert.Equal(m.T(), licenseTypeSampleUpdated, out)
	})

	m.T().Run("id_not_exist", func(t *testing.T) {
		licenseTypeSampleUpdated := model.OracleDatabaseLicenseType{
			ID: "doesn't exist",
		}
		err := m.db.UpdateOracleDatabaseLicenseType(licenseTypeSampleUpdated)

		require.Equal(t, utils.ErrOracleDatabaseLicenseTypeIDNotFound, err)
	})
}

func (m *MongodbSuite) TestRemoveOracleDatabaseLicenseType() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbLicenseTypesCollection).DeleteMany(context.TODO(), bson.M{})

	err := m.db.InsertOracleDatabaseLicenseType(licenseTypeExample)
	require.NoError(m.T(), err)

	err = m.db.RemoveOracleDatabaseLicenseType(licenseTypeExample.ID)
	require.NoError(m.T(), err)

	err = m.db.RemoveOracleDatabaseLicenseType("Test")
	require.Equal(m.T(), utils.ErrOracleDatabaseLicenseTypeIDNotFound, err)
}

func (m *MongodbSuite) TestLicenseHostIgnoredField() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))

	m.T().Run("update_ignored", func(t *testing.T) {

		hostname, dbname, licenseTypeID := "test-db", "ERCOLE1", "A90611"
		ignored := false

		err := m.db.UpdateLicenseIgnoredField(hostname, dbname, licenseTypeID, ignored)
		require.NoError(t, err)
	})
}
