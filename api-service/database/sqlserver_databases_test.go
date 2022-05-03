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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestSearchSqlServerInstances() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_27.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_28.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchSqlServerInstances([]string{""}, "", false, -1, -1, "", "PROD", utils.MAX_TIME)
		m.Require().NoError(err)

		expectedOut := dto.SqlServerInstanceResponse{
			Content: []dto.SqlServerInstance{},
			Metadata: dto.PagingMetadata{
				Empty:         true,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          0,
				TotalElements: 0,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchSqlServerInstances([]string{""}, "", false, -1, -1, "France", "", utils.MAX_TIME)
		m.Require().NoError(err)

		expectedOut := dto.SqlServerInstanceResponse{
			Content: []dto.SqlServerInstance{},
			Metadata: dto.PagingMetadata{
				Empty:         true,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          0,
				TotalElements: 0,
				TotalPages:    0,
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchSqlServerInstances([]string{""}, "", false, -1, -1, "", "", utils.P("1999-05-04T16:09:46.608+02:00"))
		m.Require().NoError(err)

		expectedOut := dto.SqlServerInstanceResponse{
			Content: []dto.SqlServerInstance{},
			Metadata: dto.PagingMetadata{
				Empty:         true,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          0,
				TotalElements: 0,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_paging", func(t *testing.T) {
		out, err := m.db.SearchSqlServerInstances([]string{""}, "", false, 0, 1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		var expectedContent []dto.SqlServerInstance = []dto.SqlServerInstance{
			{
				Hostname:      "test-db",
				Environment:   "TST",
				Name:          "MSSQLSERVER",
				Status:        "ONLINE",
				Edition:       "ENT",
				CollationName: "Latin1_General_CI_AS",
				Version:       "15.0.2000.5",
			},
		}

		expectedOut := dto.SqlServerInstanceResponse{
			Content: expectedContent,
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          false,
				Number:        0,
				Size:          1,
				TotalElements: 2,
				TotalPages:    2,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorting", func(t *testing.T) {
		out, err := m.db.SearchSqlServerInstances([]string{""}, "hostname", true, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedContent []dto.SqlServerInstance = []dto.SqlServerInstance{
			{
				Hostname:      "test-db2",
				Environment:   "PRD",
				Name:          "MSSQLSERVER",
				Status:        "ONLINE",
				Edition:       "STD",
				CollationName: "Latin1_General_CI_AS",
				Version:       "15.0.2000.5",
			},
		}

		expectedOut := dto.SqlServerInstanceResponse{
			Content: expectedContent,
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          2,
				TotalElements: 2,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut.Content[0]), utils.ToJSON(out.Content[0]))
	})

	m.T().Run("should_search_return_anything", func(t *testing.T) {
		out, err := m.db.SearchSqlServerInstances([]string{"foobar"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedContent []dto.SqlServerInstance = []dto.SqlServerInstance{}

		expectedOut := dto.SqlServerInstanceResponse{
			Content: expectedContent,
			Metadata: dto.PagingMetadata{
				Empty:         true,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          0,
				TotalElements: 0,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_found", func(t *testing.T) {
		out, err := m.db.SearchSqlServerInstances([]string{"test-db2"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedContent []dto.SqlServerInstance = []dto.SqlServerInstance{
			{
				Hostname:      "test-db2",
				Environment:   "PRD",
				Name:          "MSSQLSERVER",
				Status:        "ONLINE",
				Edition:       "STD",
				CollationName: "Latin1_General_CI_AS",
				Version:       "15.0.2000.5",
			},
		}

		expectedOut := dto.SqlServerInstanceResponse{
			Content: expectedContent,
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          1,
				TotalElements: 1,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("fullmode", func(t *testing.T) {
		out, err := m.db.SearchSqlServerInstances([]string{""}, "hostname", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedContent []dto.SqlServerInstance = []dto.SqlServerInstance{
			{
				Hostname:      "test-db",
				Environment:   "TST",
				Name:          "MSSQLSERVER",
				Status:        "ONLINE",
				Edition:       "ENT",
				CollationName: "Latin1_General_CI_AS",
				Version:       "15.0.2000.5",
			},
			{
				Hostname:      "test-db2",
				Environment:   "PRD",
				Name:          "MSSQLSERVER",
				Status:        "ONLINE",
				Edition:       "STD",
				CollationName: "Latin1_General_CI_AS",
				Version:       "15.0.2000.5",
			},
		}

		expectedOut := dto.SqlServerInstanceResponse{
			Content: expectedContent,
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          2,
				TotalElements: 2,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

}
