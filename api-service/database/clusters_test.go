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

func (m *MongodbSuite) TestSearchClusters() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchClusters(false, []string{""}, "", false, -1, -1, "", "TST", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchClusters(false, []string{""}, "", false, -1, -1, "France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchClusters(false, []string{""}, "", false, -1, -1, "", "", utils.P("1999-05-04T16:09:46.608+02:00"))
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_paging", func(t *testing.T) {
		out, err := m.db.SearchClusters(false, []string{""}, "", false, 0, 1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"Content": []interface{}{
					map[string]interface{}{
						"CPU":                         140,
						"Environment":                 "PROD",
						"Hostname":                    "test-virt",
						"HostnameAgentVirtualization": "test-virt",
						"Location":                    "Italy",
						"Name":                        "Puzzait",
						"PhysicalHosts":               "s157-cb32c10a56c256746c337e21b3f82402",
						"Sockets":                     10,
						"Type":                        "vmware",
						"_id":                         utils.Str2oid("5eb0222a45d85f4193704944"),
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
		out, err := m.db.SearchClusters(false, []string{""}, "Sockets", true, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"CPU":                         130,
				"Environment":                 "PROD",
				"Hostname":                    "test-virt",
				"HostnameAgentVirtualization": "test-virt",
				"Location":                    "Italy",
				"Name":                        "Puzzait2",
				"PhysicalHosts":               "s157-cb32c10a56c256746c337e21b3fffeua s157-cb32c10a56c256746c337e21b3ffffff",
				"Sockets":                     13,
				"Type":                        "vmware",
				"_id":                         utils.Str2oid("5eb0222a45d85f4193704944"),
			}, map[string]interface{}{
				"CPU":                         140,
				"Environment":                 "PROD",
				"Hostname":                    "test-virt",
				"HostnameAgentVirtualization": "test-virt",
				"Location":                    "Italy",
				"Name":                        "Puzzait",
				"PhysicalHosts":               "s157-cb32c10a56c256746c337e21b3f82402",
				"Sockets":                     10,
				"Type":                        "vmware",
				"_id":                         utils.Str2oid("5eb0222a45d85f4193704944"),
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_anything", func(t *testing.T) {
		out, err := m.db.SearchClusters(false, []string{"foobar"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_found", func(t *testing.T) {
		out, err := m.db.SearchClusters(false, []string{"Puzzait2"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"CPU":                         130,
				"Environment":                 "PROD",
				"Hostname":                    "test-virt",
				"HostnameAgentVirtualization": "test-virt",
				"Location":                    "Italy",
				"Name":                        "Puzzait2",
				"PhysicalHosts":               "s157-cb32c10a56c256746c337e21b3fffeua s157-cb32c10a56c256746c337e21b3ffffff",
				"Sockets":                     13,
				"Type":                        "vmware",
				"_id":                         utils.Str2oid("5eb0222a45d85f4193704944"),
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("fullmode", func(t *testing.T) {
		out, err := m.db.SearchClusters(true, []string{""}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"CPU":                         140,
				"Environment":                 "PROD",
				"Hostname":                    "test-virt",
				"HostnameAgentVirtualization": "test-virt",
				"Location":                    "Italy",
				"Name":                        "Puzzait",
				"PhysicalHosts":               []string{"s157-cb32c10a56c256746c337e21b3f82402"},
				"Sockets":                     10,
				"Type":                        "vmware",
				"_id":                         utils.Str2oid("5eb0222a45d85f4193704944"),
				"CreatedAt":                   utils.P("2020-05-04T16:09:46.608+02:00").Local(),
				"VMs": []interface{}{
					map[string]interface{}{
						"CappedCPU":    false,
						"Hostname":     "test-virt",
						"Name":         "test-virt",
						"PhysicalHost": "s157-cb32c10a56c256746c337e21b3f82402",
					},
					map[string]interface{}{
						"CappedCPU":    false,
						"Hostname":     "test-db",
						"Name":         "test-db",
						"PhysicalHost": "s157-cb32c10a56c256746c337e21b3f82402",
					},
				},
			},
			map[string]interface{}{
				"CPU":                         130,
				"Environment":                 "PROD",
				"Hostname":                    "test-virt",
				"HostnameAgentVirtualization": "test-virt",
				"Location":                    "Italy",
				"Name":                        "Puzzait2",
				"PhysicalHosts":               []string{"s157-cb32c10a56c256746c337e21b3fffeua", "s157-cb32c10a56c256746c337e21b3ffffff"},
				"Sockets":                     13,
				"Type":                        "vmware",
				"_id":                         utils.Str2oid("5eb0222a45d85f4193704944"),
				"CreatedAt":                   utils.P("2020-05-04T16:09:46.608+02:00").Local(),
				"VMs": []interface{}{
					map[string]interface{}{
						"CappedCPU":    false,
						"Hostname":     "test-virt2",
						"Name":         "test-virt2",
						"PhysicalHost": "s157-cb32c10a56c256746c337e21b3ffffff",
					},
					map[string]interface{}{
						"CappedCPU":    false,
						"Hostname":     "test-db2",
						"Name":         "test-db2",
						"PhysicalHost": "s157-cb32c10a56c256746c337e21b3fffeua",
					},
				},
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

}
