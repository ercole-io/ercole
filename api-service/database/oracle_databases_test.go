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

	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
	"github.com/stretchr/testify/assert"
)

func (m *MongodbSuite) TestSearchOracleDatabases() {
	var work float64 = 1
	enabled := false
	name := "ECXSERVER"
	creationdate := utils.P("2019-06-24T17:34:20Z")

	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_09.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabases([]string{""}, "", false, -1, -1, "", "PROD", utils.MAX_TIME)
		m.Require().NoError(err)

		expectedOut := dto.OracleDatabaseResponse{
			Content: []dto.OracleDatabase{},
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
		out, err := m.db.SearchOracleDatabases([]string{""}, "", false, -1, -1, "France", "", utils.MAX_TIME)
		m.Require().NoError(err)

		expectedOut := dto.OracleDatabaseResponse{
			Content: []dto.OracleDatabase{},
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
		out, err := m.db.SearchOracleDatabases([]string{""}, "", false, -1, -1, "", "", utils.P("1999-05-04T16:09:46.608+02:00"))
		m.Require().NoError(err)

		expectedOut := dto.OracleDatabaseResponse{
			Content: []dto.OracleDatabase{},
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
		out, err := m.db.SearchOracleDatabases([]string{""}, "", false, 0, 1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		var expectedContent []dto.OracleDatabase = []dto.OracleDatabase{
			{
				Archivelog:   false,
				BlockSize:    8192,
				CPUCount:     2,
				Charset:      "AL32UTF8",
				DatafileSize: 6,
				Dataguard:    false,
				Environment:  "TST",
				Ha:           false,
				Hostname:     "test-db",
				IsCDB:        false,
				Location:     "Germany",
				Memory:       1.484,
				MemoryTarget: 1.484,
				Name:         "ERCOLE",
				Rac:          false,
				SegmentsSize: 3,
				Services:     []model.OracleDatabaseService{},
				Status:       "OPEN",
				UniqueName:   "ERCOLE",
				Version:      "12.2.0.1.0 Enterprise Edition",
				Work:         &work,
			},
		}

		expectedOut := dto.OracleDatabaseResponse{
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
		out, err := m.db.SearchOracleDatabases([]string{""}, "memory", true, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedContent []dto.OracleDatabase = []dto.OracleDatabase{
			{
				Archivelog:   true,
				BlockSize:    8192,
				Charset:      "AL32UTF8",
				CPUCount:     2,
				DatafileSize: 6,
				Dataguard:    true,
				Environment:  "TST",
				Ha:           true,
				Hostname:     "test-db2",
				IsCDB:        true,
				Location:     "Germany",
				Memory:       90.254,
				MemoryTarget: 1.484,
				Name:         "pokemons",
				Rac:          true,
				SegmentsSize: 3,
				Services: []model.OracleDatabaseService{
					{
						CreationDate: &creationdate,
						Enabled:      &enabled,
						Name:         &name,
					},
				},
				Status:     "OPEN",
				UniqueName: "pokemons",
				Version:    "12.2.0.1.0 Enterprise Edition",
				Work:       &work,
			},
		}

		expectedOut := dto.OracleDatabaseResponse{
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
		out, err := m.db.SearchOracleDatabases([]string{"foobar"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedContent []dto.OracleDatabase = []dto.OracleDatabase{}

		expectedOut := dto.OracleDatabaseResponse{
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
		out, err := m.db.SearchOracleDatabases([]string{"pokemon", "test-db2"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedContent []dto.OracleDatabase = []dto.OracleDatabase{
			{
				Archivelog:   true,
				BlockSize:    8192,
				CPUCount:     2,
				Charset:      "AL32UTF8",
				DatafileSize: 6,
				Dataguard:    true,
				Environment:  "TST",
				Ha:           true,
				Hostname:     "test-db2",
				IsCDB:        true,
				Location:     "Germany",
				Memory:       90.254,
				MemoryTarget: 1.484,
				Name:         "pokemons",
				Rac:          true,
				SegmentsSize: 3,
				Services: []model.OracleDatabaseService{
					{
						CreationDate: &creationdate,
						Enabled:      &enabled,
						Name:         &name,
					},
				},
				Status:     "OPEN",
				UniqueName: "pokemons",
				Version:    "12.2.0.1.0 Enterprise Edition",
				Work:       &work,
			},
		}

		expectedOut := dto.OracleDatabaseResponse{
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
		out, err := m.db.SearchOracleDatabases([]string{""}, "memory", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedContent []dto.OracleDatabase = []dto.OracleDatabase{
			{
				Archivelog:   false,
				BlockSize:    8192,
				CPUCount:     2,
				Charset:      "AL32UTF8",
				DatafileSize: 6,
				Dataguard:    false,
				Environment:  "TST",
				Ha:           false,
				Hostname:     "test-db",
				IsCDB:        false,
				Location:     "Germany",
				Memory:       1.484,
				MemoryTarget: 1.484,
				Name:         "ERCOLE",
				Rac:          false,
				SegmentsSize: 3,
				Services:     []model.OracleDatabaseService{},
				Status:       "OPEN",
				UniqueName:   "ERCOLE",
				Version:      "12.2.0.1.0 Enterprise Edition",
				Work:         &work,
			},
			{
				Archivelog:   true,
				BlockSize:    8192,
				CPUCount:     2,
				Charset:      "AL32UTF8",
				DatafileSize: 6,
				Dataguard:    true,
				Environment:  "TST",
				Ha:           true,
				Hostname:     "test-db2",
				IsCDB:        true,
				Location:     "Germany",
				Memory:       90.254,
				MemoryTarget: 1.484,
				Name:         "pokemons",
				Rac:          true,
				SegmentsSize: 3,
				Services: []model.OracleDatabaseService{
					{
						CreationDate: &creationdate,
						Enabled:      &enabled,
						Name:         &name,
					},
				},
				Status:     "OPEN",
				UniqueName: "pokemons",
				Version:    "12.2.0.1.0 Enterprise Edition",
				Work:       &work,
			},
		}

		expectedOut := dto.OracleDatabaseResponse{
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
