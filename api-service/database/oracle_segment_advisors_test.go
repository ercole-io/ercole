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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestSearchOracleDatabaseSegmentAdvisors() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_12.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_13.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseSegmentAdvisors([]string{""}, "", false, "", "PROD", utils.MAX_TIME)
		m.Require().NoError(err)
		expectedOut := []dto.OracleDatabaseSegmentAdvisor{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseSegmentAdvisors([]string{""}, "", false, "France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		expectedOut := []dto.OracleDatabaseSegmentAdvisor{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseSegmentAdvisors([]string{""}, "", false, "", "", utils.P("1999-05-04T16:09:46.608+02:00"))
		m.Require().NoError(err)
		expectedOut := []dto.OracleDatabaseSegmentAdvisor{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorting", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseSegmentAdvisors([]string{""}, "dbname", true, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		expectedOut := []dto.OracleDatabaseSegmentAdvisor{
			{
				SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
				SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
				SegmentType:    "TABLE",
				SegmentsSize:   18,
				PartitionName:  "iyyiuyyoy",
				Reclaimable:    534.34,
				Retrieve:       29.685555555555556,
				Recommendation: "32b36a77e7481343ef175483c086859e",
				CreatedAt:      utils.P("2020-05-13T10:08:23.885+02:00").UTC(),
				Dbname:         "foobar4",
				Environment:    "PRD",
				Hostname:       "test-db3",
				Location:       "Germany",
			},
			{
				SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
				SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
				SegmentType:    "TABLE",
				SegmentsSize:   21,
				PartitionName:  "iyyiuyyoy",
				Reclaimable:    4.3,
				Retrieve:       0.20476190476190476,
				Recommendation: "32b36a77e7481343ef175483c086859e",
				CreatedAt:      utils.P("2020-05-13T10:08:23.885+02:00").UTC(),
				Dbname:         "foobar3",
				Environment:    "PRD",
				Hostname:       "test-db3",
				Location:       "Germany",
			},
			{
				SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
				SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
				SegmentType:    "TABLE",
				SegmentsSize:   3.0,
				PartitionName:  "iyyiuyyoy",
				Reclaimable:    0.5,
				Retrieve:       0.16666666666666666,
				Recommendation: "32b36a77e7481343ef175483c086859e",
				CreatedAt:      utils.P("2020-05-13T10:08:23.885+02:00").UTC(),
				Dbname:         "foobar2",
				Environment:    "TST",
				Hostname:       "test-db2",
				Location:       "Germany",
			},
			{
				SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
				SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
				SegmentType:    "TABLE",
				SegmentsSize:   6,
				PartitionName:  "iyyiuyyoy",
				Reclaimable:    0.5,
				Retrieve:       0.08333333333333333,
				Recommendation: "32b36a77e7481343ef175483c086859e",
				CreatedAt:      utils.P("2020-05-13T10:08:23.885+02:00").UTC(),
				Dbname:         "foobar1",
				Environment:    "TST",
				Hostname:       "test-db2",
				Location:       "Germany",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_anything", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseSegmentAdvisors([]string{"barfoo"}, "", false, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		expectedOut := []dto.OracleDatabaseSegmentAdvisor{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_found", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseSegmentAdvisors([]string{"test-db2", "foobar1"}, "", false, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		expectedOut := []dto.OracleDatabaseSegmentAdvisor{
			{
				SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
				SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
				SegmentType:    "TABLE",
				SegmentsSize:   6,
				PartitionName:  "iyyiuyyoy",
				Reclaimable:    0.5,
				Retrieve:       0.08333333333333333,
				Recommendation: "32b36a77e7481343ef175483c086859e",
				CreatedAt:      utils.P("2020-05-13T10:08:23.885+02:00").UTC(),
				Dbname:         "foobar1",
				Environment:    "TST",
				Hostname:       "test-db2",
				Location:       "Germany",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_correct_results", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseSegmentAdvisors([]string{""}, "", false, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		expectedOut := []dto.OracleDatabaseSegmentAdvisor{
			{
				SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
				SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
				SegmentType:    "TABLE",
				SegmentsSize:   6,
				PartitionName:  "iyyiuyyoy",
				Reclaimable:    0.5,
				Retrieve:       0.08333333333333333,
				Recommendation: "32b36a77e7481343ef175483c086859e",
				CreatedAt:      utils.P("2020-05-13T10:08:23.885+02:00").UTC(),
				Dbname:         "foobar1",
				Environment:    "TST",
				Hostname:       "test-db2",
				Location:       "Germany",
			},
			{
				SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
				SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
				SegmentType:    "TABLE",
				SegmentsSize:   3.0,
				PartitionName:  "iyyiuyyoy",
				Reclaimable:    0.5,
				Retrieve:       0.16666666666666666,
				Recommendation: "32b36a77e7481343ef175483c086859e",
				CreatedAt:      utils.P("2020-05-13T10:08:23.885+02:00").UTC(),
				Dbname:         "foobar2",
				Environment:    "TST",
				Hostname:       "test-db2",
				Location:       "Germany",
			},
			{
				SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
				SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
				SegmentType:    "TABLE",
				SegmentsSize:   21,
				PartitionName:  "iyyiuyyoy",
				Reclaimable:    4.3,
				Retrieve:       0.20476190476190476,
				Recommendation: "32b36a77e7481343ef175483c086859e",
				CreatedAt:      utils.P("2020-05-13T10:08:23.885+02:00").UTC(),
				Dbname:         "foobar3",
				Environment:    "PRD",
				Hostname:       "test-db3",
				Location:       "Germany",
			},
			{
				SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
				SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
				SegmentType:    "TABLE",
				SegmentsSize:   18,
				PartitionName:  "iyyiuyyoy",
				Reclaimable:    534.34,
				Retrieve:       29.685555555555556,
				Recommendation: "32b36a77e7481343ef175483c086859e",
				CreatedAt:      utils.P("2020-05-13T10:08:23.885+02:00").UTC(),
				Dbname:         "foobar4",
				Environment:    "PRD",
				Hostname:       "test-db3",
				Location:       "Germany",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}
