// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"time"

	"github.com/ercole-io/ercole/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestSearchOracleDatabasePatchAdvisors() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_12.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_13.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{""}, "", false, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "PROD", utils.MAX_TIME, "")
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{""}, "", false, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "France", "", utils.MAX_TIME, "")
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{""}, "", false, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "", utils.P("1999-05-04T16:09:46.608+02:00"), "")
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_paging", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{""}, "", false, 0, 1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "", utils.MAX_TIME, "")
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"Content": []interface{}{
					map[string]interface{}{
						"_id":         utils.Str2oid("5ebbaaf747c3fcf9dc0a1f51"),
						"Hostname":    "test-db2",
						"Location":    "Germany",
						"Environment": "TST",
						"Date":        utils.P("2012-04-16T02:00:00+02:00").Local(),
						"CreatedAt":   utils.P("2020-05-13T10:08:23.885+02:00").Local(),
						"Dbname":      "foobar1",
						"Dbver":       "12.2.0.1.0 Enterprise Edition",
						"Description": "PSU 11.2.0.3.2",
						"Status":      "KO",
					},
				},
				"Metadata": map[string]interface{}{
					"Empty":         false,
					"First":         true,
					"Last":          false,
					"Number":        0,
					"Size":          1,
					"TotalElements": 4,
					"TotalPages":    4,
				},
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorting", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{}, "Dbname", true, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "", utils.MAX_TIME, "")
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"_id":         utils.Str2oid("5ec2518bbc4991e955e2cb3f"),
				"Hostname":    "test-db3",
				"Location":    "Germany",
				"Environment": "PRD",
				"Date":        time.Unix(0, 0),
				"CreatedAt":   utils.P("2020-05-13T10:08:23.885+02:00").Local(),
				"Dbname":      "foobar4",
				"Dbver":       "18.2.0.1.0 Enterprise Edition",
				"Description": "",
				"Status":      "KO",
			},
			map[string]interface{}{
				"_id":         utils.Str2oid("5ec2518bbc4991e955e2cb3f"),
				"Hostname":    "test-db3",
				"Location":    "Germany",
				"Environment": "PRD",
				"Date":        utils.P("2020-04-16T02:00:00+02:00").Local(),
				"CreatedAt":   utils.P("2020-05-13T10:08:23.885+02:00").Local(),
				"Dbname":      "foobar3",
				"Dbver":       "16.2.0.1.0 Enterprise Edition",
				"Description": "PSU 11.2.0.3.7",
				"Status":      "OK",
			},
			map[string]interface{}{
				"_id":         utils.Str2oid("5ebbaaf747c3fcf9dc0a1f51"),
				"Hostname":    "test-db2",
				"Location":    "Germany",
				"Environment": "TST",
				"Date":        utils.P("2012-04-16T02:00:00+02:00").Local(),
				"CreatedAt":   utils.P("2020-05-13T10:08:23.885+02:00").Local(),
				"Dbname":      "foobar2",
				"Dbver":       "12.2.0.1.0 Enterprise Edition",
				"Description": "PSU 11.2.0.3.2",
				"Status":      "KO",
			},
			map[string]interface{}{
				"_id":         utils.Str2oid("5ebbaaf747c3fcf9dc0a1f51"),
				"Hostname":    "test-db2",
				"Location":    "Germany",
				"Environment": "TST",
				"Date":        utils.P("2012-04-16T02:00:00+02:00").Local(),
				"CreatedAt":   utils.P("2020-05-13T10:08:23.885+02:00").Local(),
				"Dbname":      "foobar1",
				"Dbver":       "12.2.0.1.0 Enterprise Edition",
				"Description": "PSU 11.2.0.3.2",
				"Status":      "KO",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_anything", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{"barfoo"}, "", false, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "", utils.MAX_TIME, "")
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_found", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{"test-db2", "foobar1"}, "", false, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "", utils.MAX_TIME, "")
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"_id":         utils.Str2oid("5ebbaaf747c3fcf9dc0a1f51"),
				"Hostname":    "test-db2",
				"Location":    "Germany",
				"Environment": "TST",
				"Date":        utils.P("2012-04-16T02:00:00+02:00").Local(),
				"CreatedAt":   utils.P("2020-05-13T10:08:23.885+02:00").Local(),
				"Dbname":      "foobar1",
				"Dbver":       "12.2.0.1.0 Enterprise Edition",
				"Description": "PSU 11.2.0.3.2",
				"Status":      "KO",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_by_status", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{}, "", false, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "", utils.MAX_TIME, "OK")
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"_id":         utils.Str2oid("5ec2518bbc4991e955e2cb3f"),
				"Hostname":    "test-db3",
				"Location":    "Germany",
				"Environment": "PRD",
				"Date":        utils.P("2020-04-16T02:00:00+02:00").Local(),
				"CreatedAt":   utils.P("2020-05-13T10:08:23.885+02:00").Local(),
				"Dbname":      "foobar3",
				"Dbver":       "16.2.0.1.0 Enterprise Edition",
				"Description": "PSU 11.2.0.3.7",
				"Status":      "OK",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{}, "", false, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "", utils.MAX_TIME, "")
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"_id":         utils.Str2oid("5ebbaaf747c3fcf9dc0a1f51"),
				"Hostname":    "test-db2",
				"Location":    "Germany",
				"Environment": "TST",
				"Date":        utils.P("2012-04-16T02:00:00+02:00").Local(),
				"CreatedAt":   utils.P("2020-05-13T10:08:23.885+02:00").Local(),
				"Dbname":      "foobar1",
				"Dbver":       "12.2.0.1.0 Enterprise Edition",
				"Description": "PSU 11.2.0.3.2",
				"Status":      "KO",
			},
			map[string]interface{}{
				"_id":         utils.Str2oid("5ebbaaf747c3fcf9dc0a1f51"),
				"Hostname":    "test-db2",
				"Location":    "Germany",
				"Environment": "TST",
				"Date":        utils.P("2012-04-16T02:00:00+02:00").Local(),
				"CreatedAt":   utils.P("2020-05-13T10:08:23.885+02:00").Local(),
				"Dbname":      "foobar2",
				"Dbver":       "12.2.0.1.0 Enterprise Edition",
				"Description": "PSU 11.2.0.3.2",
				"Status":      "KO",
			},
			map[string]interface{}{
				"_id":         utils.Str2oid("5ec2518bbc4991e955e2cb3f"),
				"Hostname":    "test-db3",
				"Location":    "Germany",
				"Environment": "PRD",
				"Date":        utils.P("2020-04-16T02:00:00+02:00").Local(),
				"CreatedAt":   utils.P("2020-05-13T10:08:23.885+02:00").Local(),
				"Dbname":      "foobar3",
				"Dbver":       "16.2.0.1.0 Enterprise Edition",
				"Description": "PSU 11.2.0.3.7",
				"Status":      "OK",
			},
			map[string]interface{}{
				"_id":         utils.Str2oid("5ec2518bbc4991e955e2cb3f"),
				"Hostname":    "test-db3",
				"Location":    "Germany",
				"Environment": "PRD",
				"Date":        time.Unix(0, 0),
				"CreatedAt":   utils.P("2020-05-13T10:08:23.885+02:00").Local(),
				"Dbname":      "foobar4",
				"Dbver":       "18.2.0.1.0 Enterprise Edition",
				"Description": "",
				"Status":      "KO",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
