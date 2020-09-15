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

	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *MongodbSuite) TestSearchAlerts() {
	defer m.db.Client.Database(m.dbname).Collection("alerts").DeleteMany(context.TODO(), bson.M{})
	m.InsertAlert(model.Alert{
		ID:                      utils.Str2oid("5ea6a65bb2e36eb58da2f67c"),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryLicense,
		AlertCode:               model.AlertCodeNewOption,
		AlertSeverity:           model.AlertSeverityCritical,
		AlertStatus:             model.AlertStatusNew,
		Date:                    utils.P("2020-04-15T08:46:58.475+02:00"),
		Description:             "The database ERCOLE on test-db has enabled new features (Diagnostics Pack) on server",
		OtherInfo: map[string]interface{}{
			"dbname": "ERCOLE",
			"features": []interface{}{
				"Diagnostics Pack",
			},
			"hostname": "test-db",
		},
	})
	m.InsertAlert(model.Alert{
		ID:                      utils.Str2oid("5e96ade270c184faca93fe1b"),
		AlertAffectedTechnology: nil,
		AlertCategory:           model.AlertCategoryEngine,
		AlertCode:               model.AlertCodeNewServer,
		AlertSeverity:           model.AlertSeverityInfo,
		AlertStatus:             model.AlertStatusAck,
		Date:                    utils.P("2020-04-10T08:46:58.38+02:00"),
		Description:             "The server 'rac1_x' was added to ercole",
		OtherInfo: map[string]interface{}{
			"hostname": "rac1_x",
		},
	})
	m.InsertAlert(model.Alert{
		ID:                      utils.Str2oid("5eb5057f780da34946c353fb"),
		AlertAffectedTechnology: nil,
		AlertCategory:           model.AlertCategoryEngine,
		AlertCode:               model.AlertCodeNewServer,
		AlertSeverity:           model.AlertSeverityInfo,
		AlertStatus:             model.AlertStatusAck,
		Date:                    utils.P("2020-04-10T08:46:58.38+02:00"),
		Description:             "The server 'rac1_x' was added to ercole",
		OtherInfo: map[string]interface{}{
			"hostname": "rac1_x",
		},
	})
	m.InsertAlert(model.Alert{
		ID:                      utils.Str2oid("5eb5058de2a09300d98aab67"),
		AlertAffectedTechnology: nil,
		AlertCategory:           model.AlertCategoryEngine,
		AlertCode:               model.AlertCodeNewServer,
		AlertSeverity:           model.AlertSeverityInfo,
		AlertStatus:             model.AlertStatusAck,
		Date:                    utils.P("2020-04-10T08:46:58.38+02:00"),
		Description:             "The server 'rac2_x' was added to ercole",
		OtherInfo: map[string]interface{}{
			"hostname": "rac2_x",
		},
	})

	m.T().Run("should_be_paging", func(t *testing.T) {
		out, err := m.db.SearchAlerts("all", []string{}, "", false, 0, 1, "", "", utils.MIN_TIME, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"content": []interface{}{
					map[string]interface{}{
						"_id":                     utils.Str2oid("5ea6a65bb2e36eb58da2f67c"),
						"alertAffectedTechnology": "Oracle/Database",
						"alertCategory":           model.AlertCategoryLicense,
						"alertCode":               model.AlertCodeNewOption,
						"alertSeverity":           model.AlertSeverityCritical,
						"alertStatus":             model.AlertStatusNew,
						"date":                    utils.P("2020-04-15T08:46:58.475+02:00").Local(),
						"description":             "The database ERCOLE on test-db has enabled new features (Diagnostics Pack) on server",
						"otherInfo": map[string]interface{}{
							"dbname": "ERCOLE",
							"features": []interface{}{
								"Diagnostics Pack",
							},
							"hostname": "test-db",
						},
						"hostname": "test-db",
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
		out, err := m.db.SearchAlerts("all", []string{}, "alertSeverity", true, -1, -1, "", "", utils.MIN_TIME, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"_id":                     utils.Str2oid("5e96ade270c184faca93fe1b"),
				"alertAffectedTechnology": nil,
				"alertCategory":           model.AlertCategoryEngine,
				"alertCode":               model.AlertCodeNewServer,
				"alertSeverity":           model.AlertSeverityInfo,
				"alertStatus":             model.AlertStatusAck,
				"date":                    utils.P("2020-04-10T08:46:58.38+02:00").Local(),
				"description":             "The server 'rac1_x' was added to ercole",
				"otherInfo": map[string]interface{}{
					"hostname": "rac1_x",
				},
				"hostname": "rac1_x",
			},
			map[string]interface{}{
				"_id":                     utils.Str2oid("5eb5057f780da34946c353fb"),
				"alertAffectedTechnology": nil,
				"alertCategory":           model.AlertCategoryEngine,
				"alertCode":               model.AlertCodeNewServer,
				"alertSeverity":           model.AlertSeverityInfo,
				"alertStatus":             model.AlertStatusAck,
				"date":                    utils.P("2020-04-10T08:46:58.38+02:00").Local(),
				"description":             "The server 'rac1_x' was added to ercole",
				"otherInfo": map[string]interface{}{
					"hostname": "rac1_x",
				},
				"hostname": "rac1_x",
			},
			map[string]interface{}{
				"_id":                     utils.Str2oid("5eb5058de2a09300d98aab67"),
				"alertAffectedTechnology": nil,
				"alertCategory":           model.AlertCategoryEngine,
				"alertCode":               model.AlertCodeNewServer,
				"alertSeverity":           model.AlertSeverityInfo,
				"alertStatus":             model.AlertStatusAck,
				"date":                    utils.P("2020-04-10T08:46:58.38+02:00").Local(),
				"description":             "The server 'rac2_x' was added to ercole",
				"otherInfo": map[string]interface{}{
					"hostname": "rac2_x",
				},
				"hostname": "rac2_x",
			},
			map[string]interface{}{
				"_id":                     utils.Str2oid("5ea6a65bb2e36eb58da2f67c"),
				"alertAffectedTechnology": "Oracle/Database",
				"alertCategory":           model.AlertCategoryLicense,
				"alertCode":               model.AlertCodeNewOption,
				"alertSeverity":           model.AlertSeverityCritical,
				"alertStatus":             model.AlertStatusNew,
				"date":                    utils.P("2020-04-15T08:46:58.475+02:00").Local(),
				"description":             "The database ERCOLE on test-db has enabled new features (Diagnostics Pack) on server",
				"otherInfo": map[string]interface{}{
					"dbname": "ERCOLE",
					"features": []interface{}{
						"Diagnostics Pack",
					},
					"hostname": "test-db",
				},
				"hostname": "test-db",
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_filter_by_status", func(t *testing.T) {
		out, err := m.db.SearchAlerts("all", []string{}, "", false, -1, -1, "", model.AlertStatusNew, utils.MIN_TIME, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"_id":                     utils.Str2oid("5ea6a65bb2e36eb58da2f67c"),
				"alertAffectedTechnology": "Oracle/Database",
				"alertCategory":           model.AlertCategoryLicense,
				"alertCode":               model.AlertCodeNewOption,
				"alertSeverity":           model.AlertSeverityCritical,
				"alertStatus":             model.AlertStatusNew,
				"date":                    utils.P("2020-04-15T08:46:58.475+02:00").Local(),
				"description":             "The database ERCOLE on test-db has enabled new features (Diagnostics Pack) on server",
				"otherInfo": map[string]interface{}{
					"dbname": "ERCOLE",
					"features": []interface{}{
						"Diagnostics Pack",
					},
					"hostname": "test-db",
				},
				"hostname": "test-db",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_filter_by_severity", func(t *testing.T) {
		out, err := m.db.SearchAlerts("all", []string{}, "", false, -1, -1, model.AlertSeverityCritical, "", utils.MIN_TIME, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"_id":                     utils.Str2oid("5ea6a65bb2e36eb58da2f67c"),
				"alertAffectedTechnology": "Oracle/Database",
				"alertCategory":           model.AlertCategoryLicense,
				"alertCode":               model.AlertCodeNewOption,
				"alertSeverity":           model.AlertSeverityCritical,
				"alertStatus":             model.AlertStatusNew,
				"date":                    utils.P("2020-04-15T08:46:58.475+02:00").Local(),
				"description":             "The database ERCOLE on test-db has enabled new features (Diagnostics Pack) on server",
				"otherInfo": map[string]interface{}{
					"dbname": "ERCOLE",
					"features": []interface{}{
						"Diagnostics Pack",
					},
					"hostname": "test-db",
				},
				"hostname": "test-db",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_filter_by_from", func(t *testing.T) {
		out, err := m.db.SearchAlerts("all", []string{}, "", false, -1, -1, "", "", utils.P("2020-04-13T08:46:58.38+02:00"), utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"_id":                     utils.Str2oid("5ea6a65bb2e36eb58da2f67c"),
				"alertAffectedTechnology": "Oracle/Database",
				"alertCategory":           model.AlertCategoryLicense,
				"alertCode":               model.AlertCodeNewOption,
				"alertSeverity":           model.AlertSeverityCritical,
				"alertStatus":             model.AlertStatusNew,
				"date":                    utils.P("2020-04-15T08:46:58.475+02:00").Local(),
				"description":             "The database ERCOLE on test-db has enabled new features (Diagnostics Pack) on server",
				"otherInfo": map[string]interface{}{
					"dbname": "ERCOLE",
					"features": []interface{}{
						"Diagnostics Pack",
					},
					"hostname": "test-db",
				},
				"hostname": "test-db",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_filter_by_to", func(t *testing.T) {
		out, err := m.db.SearchAlerts("all", []string{}, "", false, -1, -1, "", "", utils.MIN_TIME, utils.P("2020-04-13T08:46:58.38+02:00"))
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"_id":                     utils.Str2oid("5e96ade270c184faca93fe1b"),
				"alertAffectedTechnology": nil,
				"alertCategory":           model.AlertCategoryEngine,
				"alertCode":               model.AlertCodeNewServer,
				"alertSeverity":           model.AlertSeverityInfo,
				"alertStatus":             model.AlertStatusAck,
				"date":                    utils.P("2020-04-10T08:46:58.38+02:00").Local(),
				"description":             "The server 'rac1_x' was added to ercole",
				"otherInfo": map[string]interface{}{
					"hostname": "rac1_x",
				},
				"hostname": "rac1_x",
			},
			map[string]interface{}{
				"_id":                     utils.Str2oid("5eb5057f780da34946c353fb"),
				"alertAffectedTechnology": nil,
				"alertCategory":           model.AlertCategoryEngine,
				"alertCode":               model.AlertCodeNewServer,
				"alertSeverity":           model.AlertSeverityInfo,
				"alertStatus":             model.AlertStatusAck,
				"date":                    utils.P("2020-04-10T08:46:58.38+02:00").Local(),
				"description":             "The server 'rac1_x' was added to ercole",
				"otherInfo": map[string]interface{}{
					"hostname": "rac1_x",
				},
				"hostname": "rac1_x",
			},
			map[string]interface{}{
				"_id":                     utils.Str2oid("5eb5058de2a09300d98aab67"),
				"alertAffectedTechnology": nil,
				"alertCategory":           model.AlertCategoryEngine,
				"alertCode":               model.AlertCodeNewServer,
				"alertSeverity":           model.AlertSeverityInfo,
				"alertStatus":             model.AlertStatusAck,
				"date":                    utils.P("2020-04-10T08:46:58.38+02:00").Local(),
				"description":             "The server 'rac2_x' was added to ercole",
				"otherInfo": map[string]interface{}{
					"hostname": "rac2_x",
				},
				"hostname": "rac2_x",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_search1", func(t *testing.T) {
		out, err := m.db.SearchAlerts("all", []string{"foobar"}, "", false, -1, -1, "", "", utils.MIN_TIME, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
	m.T().Run("should_search2", func(t *testing.T) {
		out, err := m.db.SearchAlerts("all", []string{"added", model.AlertCodeNewServer, model.AlertSeverityInfo, "rac1_x"}, "", false, -1, -1, "", "", utils.MIN_TIME, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"_id":                     utils.Str2oid("5e96ade270c184faca93fe1b"),
				"alertAffectedTechnology": nil,
				"alertCategory":           model.AlertCategoryEngine,
				"alertCode":               model.AlertCodeNewServer,
				"alertSeverity":           model.AlertSeverityInfo,
				"alertStatus":             model.AlertStatusAck,
				"date":                    utils.P("2020-04-10T08:46:58.38+02:00").Local(),
				"description":             "The server 'rac1_x' was added to ercole",
				"otherInfo": map[string]interface{}{
					"hostname": "rac1_x",
				},
				"hostname": "rac1_x",
			},
			map[string]interface{}{
				"_id":                     utils.Str2oid("5eb5057f780da34946c353fb"),
				"alertAffectedTechnology": nil,
				"alertCategory":           model.AlertCategoryEngine,
				"alertCode":               model.AlertCodeNewServer,
				"alertSeverity":           model.AlertSeverityInfo,
				"alertStatus":             model.AlertStatusAck,
				"date":                    utils.P("2020-04-10T08:46:58.38+02:00").Local(),
				"description":             "The server 'rac1_x' was added to ercole",
				"otherInfo": map[string]interface{}{
					"hostname": "rac1_x",
				},
				"hostname": "rac1_x",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search3", func(t *testing.T) {
		out, err := m.db.SearchAlerts("all", []string{"ERCOLE", "Diagnostics Pack"}, "", false, -1, -1, "", "", utils.MIN_TIME, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"_id":                     utils.Str2oid("5ea6a65bb2e36eb58da2f67c"),
				"alertAffectedTechnology": "Oracle/Database",
				"alertCategory":           model.AlertCategoryLicense,
				"alertCode":               model.AlertCodeNewOption,
				"alertSeverity":           model.AlertSeverityCritical,
				"alertStatus":             model.AlertStatusNew,
				"date":                    utils.P("2020-04-15T08:46:58.475+02:00").Local(),
				"description":             "The database ERCOLE on test-db has enabled new features (Diagnostics Pack) on server",
				"otherInfo": map[string]interface{}{
					"dbname": "ERCOLE",
					"features": []interface{}{
						"Diagnostics Pack",
					},
					"hostname": "test-db",
				},
				"hostname": "test-db",
			},
		}

		assert.JSONEq(m.T(), utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_aggregate_result_code_severity", func(t *testing.T) {
		out, err := m.db.SearchAlerts("aggregated-code-severity", []string{}, "count", false, -1, -1, "", "", utils.MIN_TIME, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"category":      "LICENSE",
				"affectedHosts": 1,
				"code":          "NEW_OPTION",
				"count":         1,
				"oldestAlert":   utils.P("2020-04-15T08:46:58.475+02:00").Local(),
				"severity":      "CRITICAL",
			},
			map[string]interface{}{
				"category":      "ENGINE",
				"affectedHosts": 2,
				"code":          "NEW_SERVER",
				"count":         3,
				"oldestAlert":   utils.P("2020-04-10T08:46:58.38+02:00").Local(),
				"severity":      "INFO",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_aggregate_result_category_severity", func(t *testing.T) {
		out, err := m.db.SearchAlerts("aggregated-category-severity", []string{}, "count", false, -1, -1, "", "", utils.MIN_TIME, utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"category":      "LICENSE",
				"affectedHosts": 1,
				"count":         1,
				"oldestAlert":   utils.P("2020-04-15T08:46:58.475+02:00").Local(),
				"severity":      "CRITICAL",
			},
			map[string]interface{}{
				"category":      "ENGINE",
				"affectedHosts": 2,
				"count":         3,
				"oldestAlert":   utils.P("2020-04-10T08:46:58.38+02:00").Local(),
				"severity":      "INFO",
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestUpdateAlertsStatus() {
	defer m.db.Client.Database(m.dbname).Collection("alerts").DeleteMany(context.TODO(), bson.M{})
	alert := model.Alert{
		ID:                      utils.Str2oid("5ea6a65bb2e36eb58da2f67c"),
		AlertCategory:           model.AlertCategoryLicense,
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCode:               model.AlertCodeNewOption,
		AlertSeverity:           model.AlertSeverityCritical,
		AlertStatus:             model.AlertStatusNew,
		Date:                    utils.P("2020-04-15T08:46:58.475Z"),
		Description:             "The database ERCOLE on test-db has enabled new features (Diagnostics Pack) on server",
		OtherInfo: map[string]interface{}{
			"dbname": "ERCOLE",
			"features": bson.A{
				"Diagnostics Pack",
			},
			"hostname": "test-db",
		},
	}
	m.InsertAlert(alert)

	m.T().Run("should_return_not_found", func(t *testing.T) {
		ids := []primitive.ObjectID{utils.Str2oid("5ea6a65bb2e36eb58daaaaaa")}
		err := m.db.UpdateAlertsStatus(ids, model.AlertStatusAck)

		assert.Equal(t, utils.AerrAlertNotFound, err)
	})

	m.T().Run("should_success", func(t *testing.T) {
		ids := []primitive.ObjectID{utils.Str2oid("5ea6a65bb2e36eb58da2f67c")}
		err := m.db.UpdateAlertsStatus(ids, model.AlertStatusAck)
		assert.NoError(t, err)

		val := m.db.Client.Database(m.dbname).Collection("alerts").FindOne(context.TODO(), bson.M{"_id": utils.Str2oid("5ea6a65bb2e36eb58da2f67c")})
		var res model.Alert
		val.Decode(&res)

		alert.AlertStatus = model.AlertStatusAck
		assert.Equal(t, alert, res)
	})
}
