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
	"github.com/ercole-io/ercole/v2/utils"
)

func (m *MongodbSuite) TestSaveAndFindPatchingFunction() {
	defer m.db.Client.Database(m.dbname).Collection("patching_functions").DeleteMany(context.TODO(), bson.M{})

	pf, err := m.db.FindPatchingFunction("foobar")

	m.Require().NoError(err)
	m.Assert().Equal(model.PatchingFunction{}, pf)

	err = m.db.SavePatchingFunction(model.PatchingFunction{
		ID:        nil,
		Code:      "dfssdfsdf",
		CreatedAt: utils.P("2020-05-20T09:53:34+00:00").UTC(),
		Hostname:  "foobar",
		Vars:      map[string]interface{}{"bar": 10},
	})

	m.Require().NoError(err)

	pf, err = m.db.FindPatchingFunction("foobar")

	m.Assert().Equal("dfssdfsdf", pf.Code)
	m.Assert().Equal(utils.P("2020-05-20T09:53:34+00:00").UTC(), pf.CreatedAt)
	m.Assert().Equal("foobar", pf.Hostname)
	m.Assert().Equal(map[string]interface{}{"bar": int32(10)}, pf.Vars)

	err = m.db.SavePatchingFunction(model.PatchingFunction{
		ID:        pf.ID,
		Code:      "ffff",
		CreatedAt: utils.P("2020-05-21T09:53:34+00:00").UTC(),
		Hostname:  "foobar",
		Vars:      map[string]interface{}{"bar": 2},
	})

	m.Require().NoError(err)

	pf, err = m.db.FindPatchingFunction("foobar")

	m.Assert().Equal("ffff", pf.Code)
	m.Assert().Equal(utils.P("2020-05-21T09:53:34+00:00").UTC(), pf.CreatedAt)
	m.Assert().Equal("foobar", pf.Hostname)
	m.Assert().Equal(map[string]interface{}{"bar": int32(2)}, pf.Vars)

}

func (m *MongodbSuite) TestSearchOracleDatabaseLicenseModifiers() {
	defer m.db.Client.Database(m.dbname).Collection("patching_functions").DeleteMany(context.TODO(), bson.M{})

	id := utils.Str2oid("5ece29086437750f8b114d60")
	m.Require().NoError(m.db.SavePatchingFunction(model.PatchingFunction{
		ID:        &id,
		Code:      "//!important",
		CreatedAt: utils.P("2020-05-21T09:53:34+00:00").UTC(),
		Hostname:  "foobar",
		Vars: map[string]interface{}{
			"licenseModifiers": map[string]interface{}{
				"foobar1": map[string]int{
					"Oracle EXE": 10,
				},
				"foobar2": map[string]int{
					"Diagnostics Pack": 20,
					"Oracle EXE":       50,
				},
			},
		},
	}))

	id2 := utils.Str2oid("5ece294be12ef084764b25e6")
	m.Require().NoError(m.db.SavePatchingFunction(model.PatchingFunction{
		ID:        &id2,
		Code:      "//!important",
		CreatedAt: utils.P("2020-05-21T09:53:34+00:00").UTC(),
		Hostname:  "foobar2",
		Vars: map[string]interface{}{
			"licenseModifiers": map[string]interface{}{
				"foobar3": map[string]int{
					"Diagnostics Pack": 70,
				},
			},
		},
	}))

	m.T().Run("should_be_paging", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseLicenseModifiers([]string{""}, "newValue", false, 0, 1)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"content": []interface{}{
					map[string]interface{}{
						"_id":          utils.Str2oid("5ece29086437750f8b114d60"),
						"hostname":     "foobar",
						"databaseName": "foobar1",
						"licenseName":  "Oracle EXE",
						"newValue":     10,
					},
				},
				"metadata": map[string]interface{}{
					"empty":         false,
					"first":         true,
					"last":          false,
					"number":        0,
					"size":          1,
					"totalElements": 4,
					"totalPages":    4,
				},
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorting", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseLicenseModifiers([]string{""}, "newValue", true, -1, -1)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"_id":          utils.Str2oid("5ece294be12ef084764b25e6"),
				"hostname":     "foobar2",
				"databaseName": "foobar3",
				"licenseName":  "Diagnostics Pack",
				"newValue":     70,
			},
			map[string]interface{}{
				"_id":          utils.Str2oid("5ece29086437750f8b114d60"),
				"hostname":     "foobar",
				"databaseName": "foobar2",
				"licenseName":  "Oracle EXE",
				"newValue":     50,
			},
			map[string]interface{}{
				"_id":          utils.Str2oid("5ece29086437750f8b114d60"),
				"hostname":     "foobar",
				"databaseName": "foobar2",
				"licenseName":  "Diagnostics Pack",
				"newValue":     20,
			},
			map[string]interface{}{
				"_id":          utils.Str2oid("5ece29086437750f8b114d60"),
				"hostname":     "foobar",
				"databaseName": "foobar1",
				"licenseName":  "Oracle EXE",
				"newValue":     10,
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_anything", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseLicenseModifiers([]string{"barfoo"}, "NewValue", false, -1, -1)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_found", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseLicenseModifiers([]string{"foobar2", "foobar3", "Diagnostics Pack"}, "NewValue", false, -1, -1)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"_id":          utils.Str2oid("5ece294be12ef084764b25e6"),
				"hostname":     "foobar2",
				"databaseName": "foobar3",
				"licenseName":  "Diagnostics Pack",
				"newValue":     70,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_return_all_results", func(t *testing.T) {
		out, err := m.db.SearchOracleDatabaseLicenseModifiers([]string{""}, "newValue", false, -1, -1)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"_id":          utils.Str2oid("5ece29086437750f8b114d60"),
				"hostname":     "foobar",
				"databaseName": "foobar1",
				"licenseName":  "Oracle EXE",
				"newValue":     10,
			},
			map[string]interface{}{
				"_id":          utils.Str2oid("5ece29086437750f8b114d60"),
				"hostname":     "foobar",
				"databaseName": "foobar2",
				"licenseName":  "Diagnostics Pack",
				"newValue":     20,
			},
			map[string]interface{}{
				"_id":          utils.Str2oid("5ece29086437750f8b114d60"),
				"hostname":     "foobar",
				"databaseName": "foobar2",
				"licenseName":  "Oracle EXE",
				"newValue":     50,
			},
			map[string]interface{}{
				"_id":          utils.Str2oid("5ece294be12ef084764b25e6"),
				"hostname":     "foobar2",
				"databaseName": "foobar3",
				"licenseName":  "Diagnostics Pack",
				"newValue":     70,
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

}
