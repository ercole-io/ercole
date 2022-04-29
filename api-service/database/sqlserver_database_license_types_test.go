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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/bson"
)

var sqlServerLicenseTypeExample = model.SqlServerDatabaseLicenseType{
	ID:              "Test",
	ItemDescription: "Sql Server Database Enterprise Edition",
	Edition:         "Enterprise",
	Version:         "Test",
}

func (m *MongodbSuite) TestGetSqlServerDatabaseLicenseTypes() {
	defer m.db.Client.Database(m.dbname).Collection("ms_sqlserver_database_license_types").DeleteMany(context.TODO(), bson.M{})

	m.T().Run("success with empty table", func(t *testing.T) {
		actual, err := m.db.GetSqlServerDatabaseLicenseTypes()
		m.Require().NoError(err)

		expected := make([]model.SqlServerDatabaseLicenseType, 0)
		assert.ElementsMatch(m.T(), expected, actual)
	})

	m.T().Run("success with some values", func(t *testing.T) {
		expected := []interface{}{
			model.SqlServerDatabaseLicenseType{
				ID:              "test01",
				ItemDescription: "desc001",
				Edition:         "Standard",
				Version:         "test",
			},
			model.SqlServerDatabaseLicenseType{
				ID:              "test02",
				ItemDescription: "desc002",
				Edition:         "Standard",
				Version:         "test",
			},
		}

		ctx := context.TODO()
		_, err := m.db.Client.Database(m.dbname).
			Collection("ms_sqlserver_database_license_types").
			InsertMany(ctx, expected)
		m.Require().NoError(err)

		actual, err := m.db.GetSqlServerDatabaseLicenseTypes()
		m.Require().NoError(err)

		assert.ElementsMatch(m.T(), expected, actual)
	})
}
