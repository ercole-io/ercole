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

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
)

func (m *MongodbSuite) TestSearchClusters() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))

	m.T().Run("should_filter_out_by_environment", func(t *testing.T) {
		out, err := m.db.SearchClusters("full", []string{""}, "", false, -1, -1, "", "TST", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut []dto.Cluster

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_location", func(t *testing.T) {
		out, err := m.db.SearchClusters("full", []string{""}, "", false, -1, -1, "France", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut []dto.Cluster

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_filter_out_by_older_than", func(t *testing.T) {
		out, err := m.db.SearchClusters("full", []string{""}, "", false, -1, -1, "", "", utils.P("1999-05-04T16:09:46.608+02:00"))
		m.Require().NoError(err)
		var expectedOut []dto.Cluster

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_be_sorting", func(t *testing.T) {
		out, err := m.db.SearchClusters("full", []string{""}, "sockets", true, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		expectedOut := []dto.Cluster{
			{
				CPU:                         130,
				Environment:                 "PROD",
				Hostname:                    "test-virt",
				FetchEndpoint:               "endpoint",
				HostnameAgentVirtualization: "test-virt",
				Location:                    "Italy",
				Name:                        "Puzzait2",
				VirtualizationNodes:         []string{"s157-cb32c10a56c256746c337e21b3fffeua", "s157-cb32c10a56c256746c337e21b3ffffff"},
				PhysicalServerModelNames:    []string{"HP ProLiant DL380 Gen10", "HP ProLiant DL380 Gen9"},
				Sockets:                     13,
				Type:                        "vmware",
				VMsCount:                    2,
				VMsErcoleAgentCount:         0,
				ID:                          utils.Str2oid("5eb0222a45d85f4193704944"),
			},
			{
				CPU:                         140,
				Environment:                 "PROD",
				Hostname:                    "test-virt",
				FetchEndpoint:               "olvmmgr",
				HostnameAgentVirtualization: "test-virt",
				Location:                    "Italy",
				Name:                        "Puzzait",
				VirtualizationNodes:         []string{"s157-cb32c10a56c256746c337e21b3f82402"},
				PhysicalServerModelNames:    []string{"HP ProLiant DL380 Gen9"},
				Sockets:                     10,
				Type:                        "vmware",
				VMsCount:                    2,
				VMsErcoleAgentCount:         1,
				ID:                          utils.Str2oid("5eb0222a45d85f4193704944"),
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_anything", func(t *testing.T) {
		out, err := m.db.SearchClusters("full", []string{"foobar"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		var expectedOut []dto.Cluster

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("should_search_return_found", func(t *testing.T) {
		out, err := m.db.SearchClusters("full", []string{"Puzzait2"}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		expectedOut := []dto.Cluster{
			{
				CPU:                         130,
				Environment:                 "PROD",
				Hostname:                    "test-virt",
				FetchEndpoint:               "endpoint",
				HostnameAgentVirtualization: "test-virt",
				Location:                    "Italy",
				Name:                        "Puzzait2",
				VirtualizationNodes:         []string{"s157-cb32c10a56c256746c337e21b3fffeua", "s157-cb32c10a56c256746c337e21b3ffffff"},
				PhysicalServerModelNames:    []string{"HP ProLiant DL380 Gen10", "HP ProLiant DL380 Gen9"},
				Sockets:                     13,
				Type:                        "vmware",
				VMsCount:                    2,
				VMsErcoleAgentCount:         0,
				ID:                          utils.Str2oid("5eb0222a45d85f4193704944"),
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("fullmode", func(t *testing.T) {
		out, err := m.db.SearchClusters("full", []string{""}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)

		expectedOut := []dto.Cluster{
			{
				CPU:                         140,
				Environment:                 "PROD",
				Hostname:                    "test-virt",
				HostnameAgentVirtualization: "test-virt",
				Location:                    "Italy",
				Name:                        "Puzzait",
				FetchEndpoint:               "olvmmgr",
				VirtualizationNodes:         []string{"s157-cb32c10a56c256746c337e21b3f82402"},
				PhysicalServerModelNames:    []string{"HP ProLiant DL380 Gen9"},
				Sockets:                     10,
				Type:                        "vmware",
				ID:                          utils.Str2oid("5eb0222a45d85f4193704944"),
				VMsCount:                    2,
				VMsErcoleAgentCount:         1,
			},
			{
				CPU:                         130,
				Environment:                 "PROD",
				Hostname:                    "test-virt",
				HostnameAgentVirtualization: "test-virt",
				Location:                    "Italy",
				FetchEndpoint:               "endpoint",
				Name:                        "Puzzait2",
				VirtualizationNodes:         []string{"s157-cb32c10a56c256746c337e21b3fffeua", "s157-cb32c10a56c256746c337e21b3ffffff"},
				PhysicalServerModelNames:    []string{"HP ProLiant DL380 Gen10", "HP ProLiant DL380 Gen9"},
				Sockets:                     13,
				Type:                        "vmware",
				ID:                          utils.Str2oid("5eb0222a45d85f4193704944"),
				VMsCount:                    2,
				VMsErcoleAgentCount:         0,
			},
		}

		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})

	m.T().Run("clusternames", func(t *testing.T) {
		out, err := m.db.SearchClusters("clusternames", []string{""}, "", false, -1, -1, "", "", utils.MAX_TIME)
		m.Require().NoError(err)
		expectedOut := []dto.Cluster{
			{
				Name: "Puzzait",
			},
			{
				Name: "Puzzait2",
			},
		}
		assert.JSONEq(t, utils.ToJSON(expectedOut), utils.ToJSON(out))
	})
}

func (m *MongodbSuite) TestGetClusters() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))

	m.T().Run("No hosts", func(t *testing.T) {
		filter := dto.GlobalFilter{
			Location:  "Germany",
			OlderThan: utils.MAX_TIME,
		}
		out, err := m.db.GetClusters(filter)
		assert.Nil(m.T(), err)

		var expected []dto.Cluster
		assert.Equal(m.T(), expected, out)
	})

	m.T().Run("Found", func(t *testing.T) {
		filter := dto.GlobalFilter{
			OlderThan: utils.MAX_TIME,
			Location:  "Italy",
		}
		out, err := m.db.GetClusters(filter)
		m.Require().NoError(err)

		expected := []dto.Cluster{
			{
				ID:                          utils.Str2oid("5eb0222a45d85f4193704944"),
				CPU:                         140,
				CreatedAt:                   utils.P("2020-05-04T14:09:46.608Z"),
				Environment:                 "PROD",
				FetchEndpoint:               "olvmmgr",
				Hostname:                    "test-virt",
				HostnameAgentVirtualization: "test-virt",
				Location:                    "Italy",
				Name:                        "Puzzait",
				Sockets:                     10,
				Type:                        "vmware",
				VirtualizationNodes:         []string{"s157-cb32c10a56c256746c337e21b3f82402"},
				PhysicalServerModelNames:    []string{"HP ProLiant DL380 Gen9"},
				VirtualizationNodesCount:    1,
				VirtualizationNodesStats: []dto.VirtualizationNodesStat{
					{
						TotalVMsCount:                   2,
						TotalVMsWithErcoleAgentCount:    1,
						TotalVMsWithoutErcoleAgentCount: 1,
						VirtualizationNode:              "s157-cb32c10a56c256746c337e21b3f82402"},
				},
				VMs: []dto.VM{
					{
						CappedCPU:               false,
						Hostname:                "test-virt",
						Name:                    "test-virt",
						VirtualizationNode:      "s157-cb32c10a56c256746c337e21b3f82402",
						PhysicalServerModelName: "HP ProLiant DL380 Gen9",
					},

					{
						CappedCPU:               false,
						Hostname:                "test-db",
						Name:                    "test-db",
						VirtualizationNode:      "s157-cb32c10a56c256746c337e21b3f82402",
						PhysicalServerModelName: "HP ProLiant DL380 Gen9",
					},
				},
				VMsCount:            2,
				VMsErcoleAgentCount: 1,
			},
			{
				ID:                          utils.Str2oid("5eb0222a45d85f4193704944"),
				CPU:                         130,
				CreatedAt:                   utils.P("2020-05-04T14:09:46.608Z"),
				Environment:                 "PROD",
				FetchEndpoint:               "endpoint",
				Hostname:                    "test-virt",
				HostnameAgentVirtualization: "test-virt",
				Location:                    "Italy",
				Name:                        "Puzzait2",
				Sockets:                     13,
				Type:                        "vmware",
				VirtualizationNodes: []string{
					"s157-cb32c10a56c256746c337e21b3fffeua",
					"s157-cb32c10a56c256746c337e21b3ffffff"},
				PhysicalServerModelNames: []string{
					"HP ProLiant DL380 Gen10",
					"HP ProLiant DL380 Gen9",
				},
				VirtualizationNodesCount: 2,
				VirtualizationNodesStats: []dto.VirtualizationNodesStat{
					{
						TotalVMsCount:                   1,
						TotalVMsWithErcoleAgentCount:    0,
						TotalVMsWithoutErcoleAgentCount: 1,
						VirtualizationNode:              "s157-cb32c10a56c256746c337e21b3fffeua"},
					{
						TotalVMsCount:                   1,
						TotalVMsWithErcoleAgentCount:    0,
						TotalVMsWithoutErcoleAgentCount: 1,
						VirtualizationNode:              "s157-cb32c10a56c256746c337e21b3ffffff"},
				},
				VMs: []dto.VM{
					{
						CappedCPU:               false,
						Hostname:                "test-virt2",
						Name:                    "test-virt2",
						VirtualizationNode:      "s157-cb32c10a56c256746c337e21b3ffffff",
						PhysicalServerModelName: "HP ProLiant DL380 Gen9",
					},

					{
						CappedCPU:               false,
						Hostname:                "test-db2",
						Name:                    "test-db2",
						VirtualizationNode:      "s157-cb32c10a56c256746c337e21b3fffeua",
						PhysicalServerModelName: "HP ProLiant DL380 Gen10",
					},
				},
				VMsCount:            2,
				VMsErcoleAgentCount: 0},
		}

		assert.EqualValues(m.T(), expected, out)
	})
}

func (m *MongodbSuite) TestGetCluster() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))

	m.T().Run("ClusterNotFound", func(t *testing.T) {
		clusterName := ""
		olderThan := utils.MAX_TIME
		out, err := m.db.GetCluster(clusterName, olderThan)
		m.Require().ErrorIs(err, utils.ErrClusterNotFound)

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
			FetchEndpoint:               "olvmmgr",
			Hostname:                    "test-virt",
			HostnameAgentVirtualization: "test-virt",
			Location:                    "Italy",
			Name:                        "Puzzait",
			Sockets:                     10,
			Type:                        "vmware",
			VirtualizationNodes:         []string{"s157-cb32c10a56c256746c337e21b3f82402"},
			PhysicalServerModelNames:    []string{"HP ProLiant DL380 Gen9"},
			VirtualizationNodesCount:    1,
			VirtualizationNodesStats: []dto.VirtualizationNodesStat{
				{
					TotalVMsCount:                   2,
					TotalVMsWithErcoleAgentCount:    1,
					TotalVMsWithoutErcoleAgentCount: 1,
					VirtualizationNode:              "s157-cb32c10a56c256746c337e21b3f82402"}},

			VMs: []dto.VM{
				{
					CappedCPU:               false,
					Hostname:                "test-virt",
					Name:                    "test-virt",
					VirtualizationNode:      "s157-cb32c10a56c256746c337e21b3f82402",
					PhysicalServerModelName: "HP ProLiant DL380 Gen9",
				},
				{
					CappedCPU:               false,
					Hostname:                "test-db",
					Name:                    "test-db",
					VirtualizationNode:      "s157-cb32c10a56c256746c337e21b3f82402",
					PhysicalServerModelName: "HP ProLiant DL380 Gen9",
				},
			},
			VMsCount:            2,
			VMsErcoleAgentCount: 1,
		}

		assert.EqualValues(m.T(), expected, out)
	})
}
