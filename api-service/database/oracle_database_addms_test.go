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

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestSearchOracleDatabaseAddms() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseAddms([]string{}, "", false, -1, -1, "Italy", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseAddms([]string{}, "", false, -1, -1, "", "PRD", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseAddms([]string{}, "", false, -1, -1, "", "", utils.MIN_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_be_paging", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseAddms([]string{}, "", false, 0, 1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"content": []interface{}{
					map[string]interface{}{
						"action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
						"benefit":        83.34,
						"createdAt":      utils.P("2020-04-15T08:46:58.471+02:00").Local(),
						"dbname":         "ERCOLE",
						"environment":    "TST",
						"finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
						"hostname":       "test-db",
						"location":       "Germany",
						"recommendation": "SQL Tuning",
						"_id":            utils.Str2oid("5e96ade270c184faca93fe36"),
					},
				},
				"metadata": map[string]interface{}{
					"empty":         false,
					"first":         true,
					"last":          false,
					"number":        0,
					"size":          1,
					"totalElements": 2,
					"totalPages":    2,
				},
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_be_sorting", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseAddms([]string{}, "benefit", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
				"benefit":        68.24,
				"createdAt":      utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"dbname":         "ERCOLE",
				"environment":    "TST",
				"finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
				"hostname":       "test-db",
				"location":       "Germany",
				"recommendation": "Segment Tuning",
				"_id":            utils.Str2oid("5e96ade270c184faca93fe36"),
			},
			map[string]interface{}{
				"action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
				"benefit":        83.34,
				"createdAt":      utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"dbname":         "ERCOLE",
				"environment":    "TST",
				"finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
				"hostname":       "test-db",
				"location":       "Germany",
				"recommendation": "SQL Tuning",
				"_id":            utils.Str2oid("5e96ade270c184faca93fe36"),
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_search1", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseAddms([]string{"foobar"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_search2", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseAddms([]string{"test-db", "ERCOLE"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
				"benefit":        83.34,
				"createdAt":      utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"dbname":         "ERCOLE",
				"environment":    "TST",
				"finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
				"hostname":       "test-db",
				"location":       "Germany",
				"recommendation": "SQL Tuning",
				"_id":            utils.Str2oid("5e96ade270c184faca93fe36"),
			},
			map[string]interface{}{
				"action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
				"benefit":        68.24,
				"createdAt":      utils.P("2020-04-15T08:46:58.471+02:00").Local(),
				"dbname":         "ERCOLE",
				"environment":    "TST",
				"finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
				"hostname":       "test-db",
				"location":       "Germany",
				"recommendation": "Segment Tuning",
				"_id":            utils.Str2oid("5e96ade270c184faca93fe36"),
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
