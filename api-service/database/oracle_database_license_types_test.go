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

package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/model"
)

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
