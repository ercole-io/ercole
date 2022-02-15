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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
)

func (m *MongodbSuite) TestSearchOracleDatabasePatchAdvisors() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_12.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_13.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{""}, "", false, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "PROD", utils.MAX_TIME, "")
		m.Require().NoError(err)
		expectedOut := &dto.PatchAdvisorResponse{Content: dto.PatchAdvisors{}, Metadata: dto.PagingMetadata{Empty: true, First: true, Last: true}}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{""}, "", false, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "France", "", utils.MAX_TIME, "")
		m.Require().NoError(err)
		expectedOut := &dto.PatchAdvisorResponse{Content: dto.PatchAdvisors{}, Metadata: dto.PagingMetadata{Empty: true, First: true, Last: true}}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{""}, "", false, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "", utils.P("1999-05-04T16:09:46.608+02:00"), "")
		m.Require().NoError(err)
		expectedOut := &dto.PatchAdvisorResponse{Content: dto.PatchAdvisors{}, Metadata: dto.PagingMetadata{Empty: true, First: true, Last: true}}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_paging", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{""}, "", false, 0, 1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "", utils.MAX_TIME, "")
		m.Require().NoError(err)

		expectedOut := &dto.PatchAdvisorResponse{
			Content: dto.PatchAdvisors{
				dto.PatchAdvisor{
					ID:          utils.Str2oid("5ebbaaf747c3fcf9dc0a1f51"),
					Hostname:    "test-db2",
					Location:    "Germany",
					Environment: "TST",
					Date:        utils.PDT("2012-04-16T02:00:00+02:00"),
					CreatedAt:   utils.PDT("2020-05-13T10:08:23.885+02:00"),
					DbName:      "foobar1",
					Dbver:       "12.2.0.1.0 Enterprise Edition",
					Description: "PSU 11.2.0.3.2",
					Status:      "KO",
				},
			},
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          false,
				Number:        0,
				Size:          1,
				TotalElements: 4,
				TotalPages:    4,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorting", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{}, "dbname", true, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "", utils.MAX_TIME, "")
		m.Require().NoError(err)
		expectedOut := &dto.PatchAdvisorResponse{
			Content: dto.PatchAdvisors{
				dto.PatchAdvisor{
					ID:          utils.Str2oid("5ec2518bbc4991e955e2cb3f"),
					Hostname:    "test-db3",
					Location:    "Germany",
					Environment: "PRD",
					Date:        utils.PDT("1970-01-01T01:00:00+01:00"),
					CreatedAt:   utils.PDT("2020-05-13T10:08:23.885+02:00"),
					DbName:      "foobar4",
					Dbver:       "18.2.0.1.0 Enterprise Edition",
					Description: "",
					Status:      "KO",
				},
				dto.PatchAdvisor{
					ID:          utils.Str2oid("5ec2518bbc4991e955e2cb3f"),
					Hostname:    "test-db3",
					Location:    "Germany",
					Environment: "PRD",
					Date:        utils.PDT("2020-04-16T02:00:00+02:00"),
					CreatedAt:   utils.PDT("2020-05-13T10:08:23.885+02:00"),
					DbName:      "foobar3",
					Dbver:       "16.2.0.1.0 Enterprise Edition",
					Description: "PSU 11.2.0.3.7",
					Status:      "OK",
				},
				dto.PatchAdvisor{
					ID:          utils.Str2oid("5ebbaaf747c3fcf9dc0a1f51"),
					Hostname:    "test-db2",
					Location:    "Germany",
					Environment: "TST",
					Date:        utils.PDT("2012-04-16T02:00:00+02:00"),
					CreatedAt:   utils.PDT("2020-05-13T10:08:23.885+02:00"),
					DbName:      "foobar2",
					Dbver:       "12.2.0.1.0 Enterprise Edition",
					Description: "PSU 11.2.0.3.2",
					Status:      "KO",
				},
				dto.PatchAdvisor{
					ID:          utils.Str2oid("5ebbaaf747c3fcf9dc0a1f51"),
					Hostname:    "test-db2",
					Location:    "Germany",
					Environment: "TST",
					Date:        utils.PDT("2012-04-16T02:00:00+02:00"),
					CreatedAt:   utils.PDT("2020-05-13T10:08:23.885+02:00"),
					DbName:      "foobar1",
					Dbver:       "12.2.0.1.0 Enterprise Edition",
					Description: "PSU 11.2.0.3.2",
					Status:      "KO",
				},
			},
			Metadata: dto.PagingMetadata{
				Empty:         false,
				First:         true,
				Last:          true,
				Number:        0,
				Size:          4,
				TotalElements: 4,
				TotalPages:    0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_anything", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{"barfoo"}, "", false, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "", utils.MAX_TIME, "")
		m.Require().NoError(err)
		expectedOut := &dto.PatchAdvisorResponse{Content: dto.PatchAdvisors{}, Metadata: dto.PagingMetadata{Empty: true, First: true, Last: true}}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_found", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{"test-db2", "foobar1"}, "", false, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "", utils.MAX_TIME, "")
		m.Require().NoError(err)
		expectedOut := &dto.PatchAdvisorResponse{
			Content: dto.PatchAdvisors{
				dto.PatchAdvisor{
					ID:          utils.Str2oid("5ebbaaf747c3fcf9dc0a1f51"),
					Hostname:    "test-db2",
					Location:    "Germany",
					Environment: "TST",
					Date:        utils.PDT("2012-04-16T02:00:00+02:00"),
					CreatedAt:   utils.PDT("2020-05-13T10:08:23.885+02:00"),
					DbName:      "foobar1",
					Dbver:       "12.2.0.1.0 Enterprise Edition",
					Description: "PSU 11.2.0.3.2",
					Status:      "KO",
				},
			},
			Metadata: dto.PagingMetadata{
				First:         true,
				Last:          true,
				Size:          1,
				TotalElements: 1,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_by_status", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{}, "", false, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "", utils.MAX_TIME, "OK")
		m.Require().NoError(err)
		expectedOut := &dto.PatchAdvisorResponse{
			Content: dto.PatchAdvisors{
				dto.PatchAdvisor{
					ID:          utils.Str2oid("5ec2518bbc4991e955e2cb3f"),
					Hostname:    "test-db3",
					Location:    "Germany",
					Environment: "PRD",
					Date:        utils.PDT("2020-04-16T02:00:00+02:00"),
					CreatedAt:   utils.PDT("2020-05-13T10:08:23.885+02:00"),
					DbName:      "foobar3",
					Dbver:       "16.2.0.1.0 Enterprise Edition",
					Description: "PSU 11.2.0.3.7",
					Status:      "OK",
				},
			},
			Metadata: dto.PagingMetadata{
				First:         true,
				Last:          true,
				Size:          1,
				TotalElements: 1,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabasePatchAdvisors([]string{}, "", false, -1, -1, utils.P("2019-10-10T08:46:58.38+02:00"), "", "", utils.MAX_TIME, "")
		m.Require().NoError(err)

		expectedOut := &dto.PatchAdvisorResponse{
			Content: dto.PatchAdvisors{
				dto.PatchAdvisor{
					ID:          utils.Str2oid("5ebbaaf747c3fcf9dc0a1f51"),
					Hostname:    "test-db2",
					Location:    "Germany",
					Environment: "TST",
					Date:        utils.PDT("2012-04-16T02:00:00+02:00"),
					CreatedAt:   utils.PDT("2020-05-13T10:08:23.885+02:00"),
					DbName:      "foobar1",
					Dbver:       "12.2.0.1.0 Enterprise Edition",
					Description: "PSU 11.2.0.3.2",
					Status:      "KO",
				},
				dto.PatchAdvisor{
					ID:          utils.Str2oid("5ebbaaf747c3fcf9dc0a1f51"),
					Hostname:    "test-db2",
					Location:    "Germany",
					Environment: "TST",
					Date:        utils.PDT("2012-04-16T02:00:00+02:00"),
					CreatedAt:   utils.PDT("2020-05-13T10:08:23.885+02:00"),
					DbName:      "foobar2",
					Dbver:       "12.2.0.1.0 Enterprise Edition",
					Description: "PSU 11.2.0.3.2",
					Status:      "KO",
				},
				dto.PatchAdvisor{
					ID:          utils.Str2oid("5ec2518bbc4991e955e2cb3f"),
					Hostname:    "test-db3",
					Location:    "Germany",
					Environment: "PRD",
					Date:        utils.PDT("2020-04-16T02:00:00+02:00"),
					CreatedAt:   utils.PDT("2020-05-13T10:08:23.885+02:00"),
					DbName:      "foobar3",
					Dbver:       "16.2.0.1.0 Enterprise Edition",
					Description: "PSU 11.2.0.3.7",
					Status:      "OK",
				},
				dto.PatchAdvisor{
					ID:          utils.Str2oid("5ec2518bbc4991e955e2cb3f"),
					Hostname:    "test-db3",
					Location:    "Germany",
					Environment: "PRD",
					Date:        utils.PDT("1970-01-01T01:00:00+01:00"),
					CreatedAt:   utils.PDT("2020-05-13T10:08:23.885+02:00"),
					DbName:      "foobar4",
					Dbver:       "18.2.0.1.0 Enterprise Edition",
					Description: "",
					Status:      "KO",
				},
			},
			Metadata: dto.PagingMetadata{First: true, Last: true, Size: 4, TotalElements: 4},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
