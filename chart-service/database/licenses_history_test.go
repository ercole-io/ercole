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

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/chart-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
)

func (m *MongodbSuite) TestGetLicenseComplianceHistory() {
	m.T().Run("should_get_license_compliance_history", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection("oracle_database_licenses_history").DeleteMany(context.TODO(), bson.M{})
		m.InsertHostDataHistory(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_chartservice_mongohostdata_03.json"))

		out, err := m.db.GetLicenseComplianceHistory()
		m.Require().NoError(err)

		expectedOut := []dto.LicenseComplianceHistory{
			{
				LicenseTypeID:   "A90649",
				ItemDescription: "TEST",
				Metric:          "TEST",
				History: []dto.LicenseComplianceHistoricValue{
					{
						Date:      utils.P("2019-06-24T17:34:20Z"),
						Consumed:  10,
						Covered:   10,
						Purchased: 20,
					},
				},
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
