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
				"content": []interface{}{
					map[string]interface{}{
						"cpu":                         140,
						"environment":                 "PROD",
						"hostname":                    "test-virt",
						"hostnameAgentVirtualization": "test-virt",
						"location":                    "Italy",
						"name":                        "Puzzait",
						"virtualizationNodes":         "s157-cb32c10a56c256746c337e21b3f82402",
						"sockets":                     10,
						"type":                        "vmware",
						"vmsCount":                    2,
						"vmsErcoleAgentCount":         1,
						"_id":                         utils.Str2oid("5eb0222a45d85f4193704944"),
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
		out, err := m.db.SearchClusters(false, []string{""}, "sockets", true, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut interface{} = []interface{}{
			map[string]interface{}{
				"cpu":                         130,
				"environment":                 "PROD",
				"hostname":                    "test-virt",
				"hostnameAgentVirtualization": "test-virt",
				"location":                    "Italy",
				"name":                        "Puzzait2",
				"virtualizationNodes":         "s157-cb32c10a56c256746c337e21b3fffeua s157-cb32c10a56c256746c337e21b3ffffff",
				"sockets":                     13,
				"type":                        "vmware",
				"vmsCount":                    2,
				"vmsErcoleAgentCount":         0,
				"_id":                         utils.Str2oid("5eb0222a45d85f4193704944"),
			}, map[string]interface{}{
				"cpu":                         140,
				"environment":                 "PROD",
				"hostname":                    "test-virt",
				"hostnameAgentVirtualization": "test-virt",
				"location":                    "Italy",
				"name":                        "Puzzait",
				"virtualizationNodes":         "s157-cb32c10a56c256746c337e21b3f82402",
				"sockets":                     10,
				"type":                        "vmware",
				"vmsCount":                    2,
				"vmsErcoleAgentCount":         1,
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
				"cpu":                         130,
				"environment":                 "PROD",
				"hostname":                    "test-virt",
				"hostnameAgentVirtualization": "test-virt",
				"location":                    "Italy",
				"name":                        "Puzzait2",
				"virtualizationNodes":         "s157-cb32c10a56c256746c337e21b3fffeua s157-cb32c10a56c256746c337e21b3ffffff",
				"sockets":                     13,
				"type":                        "vmware",
				"vmsCount":                    2,
				"vmsErcoleAgentCount":         0,
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
				"cpu":                         140,
				"environment":                 "PROD",
				"hostname":                    "test-virt",
				"hostnameAgentVirtualization": "test-virt",
				"location":                    "Italy",
				"name":                        "Puzzait",
				"fetchEndpoint":               "???",
				"virtualizationNodes":         []string{"s157-cb32c10a56c256746c337e21b3f82402"},
				"sockets":                     10,
				"type":                        "vmware",
				"_id":                         utils.Str2oid("5eb0222a45d85f4193704944"),
				"createdAt":                   utils.P("2020-05-04T16:09:46.608+02:00").Local(),
				"vmsCount":                    2,
				"vmsErcoleAgentCount":         1,
				"vms": []interface{}{
					map[string]interface{}{
						"cappedCPU":          false,
						"hostname":           "test-virt",
						"name":               "test-virt",
						"virtualizationNode": "s157-cb32c10a56c256746c337e21b3f82402",
					},
					map[string]interface{}{
						"cappedCPU":          false,
						"hostname":           "test-db",
						"name":               "test-db",
						"virtualizationNode": "s157-cb32c10a56c256746c337e21b3f82402",
					},
				},
			},
			map[string]interface{}{
				"cpu":                         130,
				"environment":                 "PROD",
				"hostname":                    "test-virt",
				"hostnameAgentVirtualization": "test-virt",
				"location":                    "Italy",
				"fetchEndpoint":               "???",
				"name":                        "Puzzait2",
				"virtualizationNodes":         []string{"s157-cb32c10a56c256746c337e21b3fffeua", "s157-cb32c10a56c256746c337e21b3ffffff"},
				"sockets":                     13,
				"type":                        "vmware",
				"_id":                         utils.Str2oid("5eb0222a45d85f4193704944"),
				"createdAt":                   utils.P("2020-05-04T16:09:46.608+02:00").Local(),
				"vmsCount":                    2,
				"vmsErcoleAgentCount":         0,
				"vms": []interface{}{
					map[string]interface{}{
						"cappedCPU":          false,
						"hostname":           "test-virt2",
						"name":               "test-virt2",
						"virtualizationNode": "s157-cb32c10a56c256746c337e21b3ffffff",
					},
					map[string]interface{}{
						"cappedCPU":          false,
						"hostname":           "test-db2",
						"name":               "test-db2",
						"virtualizationNode": "s157-cb32c10a56c256746c337e21b3fffeua",
					},
				},
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetCluster() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(utils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))

	m.T().Run("HostNotFound", func(t *testing.T) {
		clusterName := ""
		olderThan := utils.MAX_TIME
		out, err := m.db.GetCluster(clusterName, olderThan)
		m.Require().Equal(err, utils.ErrHostNotFound)

		var expected *dto.Cluster
		assert.EqualValues(m.T(), expected, out)
	})

	m.T().Run("Found", func(t *testing.T) {
		clusterName := "Puzzait"
		olderThan := utils.MAX_TIME
		out, err := m.db.GetCluster(clusterName, olderThan)
		m.Require().NoError(err)

		expected := &dto.Cluster{
			ID:                          utils.Str2oid("5eb0222a45d85f4193704944"),
			CPU:                         140,
			CreatedAt:                   utils.P("2020-05-04T14:09:46.608Z"),
			Environment:                 "PROD",
			FetchEndpoint:               "???",
			Hostname:                    "test-virt",
			HostnameAgentVirtualization: "test-virt",
			Location:                    "Italy",
			Name:                        "Puzzait",
			Sockets:                     10,
			Type:                        "vmware",
			VirtualizationNodes:         []string{"s157-cb32c10a56c256746c337e21b3f82402"},
			VirtualizationNodesCount:    1,
			VirtualizationNodesStats: []dto.VirtualizationNodesStat{
				{
					TotalVMsCount:                   2,
					TotalVMsWithErcoleAgentCount:    1,
					TotalVMsWithoutErcoleAgentCount: 1,
					VirtualizationNode:              "s157-cb32c10a56c256746c337e21b3f82402"}},

			VMs: []dto.VM{
				{
					CappedCPU:          false,
					Hostname:           "test-virt",
					Name:               "test-virt",
					VirtualizationNode: "s157-cb32c10a56c256746c337e21b3f82402"},

				{
					CappedCPU:          false,
					Hostname:           "test-db",
					Name:               "test-db",
					VirtualizationNode: "s157-cb32c10a56c256746c337e21b3f82402"},
			},
			VMsCount:            2,
			VMsErcoleAgentCount: 1,
		}

		assert.EqualValues(m.T(), expected, out)
	})
}
