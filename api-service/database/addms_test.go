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

	"github.com/ercole-io/ercole/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestSearchAddms() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchAddms([]string{}, "", false, -1, -1, "Italy", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchAddms([]string{}, "", false, -1, -1, "", "PRD", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchAddms([]string{}, "", false, -1, -1, "", "", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_be_paging", func(t *testing.T) {
		out, err := m.db.SearchAddms([]string{}, "", false, 0, 1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"Content": []interface{}{
					map[string]interface{}{
						"Action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
						"Benefit":        83.34,
						"CreatedAt":      utils.P("2020-04-15T08:46:58.471+02:00").Local(),
						"Dbname":         "ERCOLE",
						"Environment":    "TST",
						"Finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
						"Hostname":       "test-db",
						"Location":       "Germany",
						"Recommendation": "SQL Tuning",
						"_id":            utils.Str2oid("5e96ade270c184faca93fe36"),
					},
				},
				"Metadata": map[string]interface{}{
					"Empty":         false,
					"First":         true,
					"Last":          false,
					"Number":        0,
					"Size":          1,
					"TotalElements": 2,
					"TotalPages":    2,
				},
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_be_sorting", func(t *testing.T) {
		out, err := m.db.SearchAddms([]string{}, "Benefit", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"Action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
				"Benefit":        68.24,
				"CreatedAt":      utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"Dbname":         "ERCOLE",
				"Environment":    "TST",
				"Finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
				"Hostname":       "test-db",
				"Location":       "Germany",
				"Recommendation": "Segment Tuning",
				"_id":            utils.Str2oid("5e96ade270c184faca93fe36"),
			},
			map[string]interface{}{
				"Action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
				"Benefit":        83.34,
				"CreatedAt":      utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"Dbname":         "ERCOLE",
				"Environment":    "TST",
				"Finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
				"Hostname":       "test-db",
				"Location":       "Germany",
				"Recommendation": "SQL Tuning",
				"_id":            utils.Str2oid("5e96ade270c184faca93fe36"),
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_search1", func(t *testing.T) {
		out, err := m.db.SearchAddms([]string{"foobar"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_search2", func(t *testing.T) {
		out, err := m.db.SearchAddms([]string{"test-db", "ERCOLE"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"Action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
				"Benefit":        83.34,
				"CreatedAt":      utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"Dbname":         "ERCOLE",
				"Environment":    "TST",
				"Finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
				"Hostname":       "test-db",
				"Location":       "Germany",
				"Recommendation": "SQL Tuning",
				"_id":            utils.Str2oid("5e96ade270c184faca93fe36"),
			},
			map[string]interface{}{
				"Action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
				"Benefit":        68.24,
				"CreatedAt":      utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"Dbname":         "ERCOLE",
				"Environment":    "TST",
				"Finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
				"Hostname":       "test-db",
				"Location":       "Germany",
				"Recommendation": "Segment Tuning",
				"_id":            utils.Str2oid("5e96ade270c184faca93fe36"),
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
