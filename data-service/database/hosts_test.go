// Copyright (c) 2021 Sorint.lab S.p.A.
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
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
)

func (m *MongodbSuite) TestDismissHost() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_03.json"))

	list, err := m.db.FindOldCurrentHostnames(utils.MAX_TIME)
	require.NoError(m.T(), err)
	require.Equal(m.T(), []string{"test-small"}, list)

	err = m.db.DismissHost("test-small")
	require.NoError(m.T(), err)

	list, err = m.db.FindOldCurrentHostnames(utils.MAX_TIME)
	require.NoError(m.T(), err)
	require.Equal(m.T(), []string{}, list)
}

func (m *MongodbSuite) TestInsertHostData() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := model.HostDataBE{
		ID:                  utils.Str2oid("5ef9d239a1d25d1e8703c4d3"),
		Archived:            false,
		CreatedAt:           utils.P("2020-06-29T13:36:25.589708509+02:00"),
		ServerVersion:       "latest",
		ServerSchemaVersion: 1,
		Hostname:            "rac1-x",
		Location:            "Germany",
		Environment:         "TST",
		AgentVersion:        "1.6.5",
		Tags:                []string{},
		Info: model.Host{
			Hostname:                      "rac1",
			CPUModel:                      "Intel(R) Xeon(R) CPU E5-2609 v4 @ 1.70GHz",
			CPUFrequency:                  "1.70GHz",
			CPUSockets:                    1,
			CPUCores:                      1,
			CPUThreads:                    2,
			ThreadsPerCore:                2,
			CoresPerSocket:                1,
			HardwareAbstraction:           "VIRT",
			HardwareAbstractionTechnology: "OVM",
			Kernel:                        "4.1.12-103.3.8.1.el7uek.x86_64",
			KernelVersion:                 "4.1.12-103.3.8.1.el7uek.x86_64",
			OS:                            "Red Hat Enterprise Linux",
			OSVersion:                     "7.4",
			MemoryTotal:                   7,
			SwapTotal:                     7,
		},
		ClusterMembershipStatus: model.ClusterMembershipStatus{OracleClusterware: true,
			VeritasClusterServer: false,
			SunCluster:           false,
			HACMP:                false,
		},
		Features: model.Features{
			Oracle: nil,
		},
		Filesystems: []model.Filesystem{},
		Clusters:    nil,
	}

	list, err := m.db.FindOldCurrentHostnames(utils.MAX_TIME)
	require.NoError(m.T(), err)
	assert.Equal(m.T(), []string{}, list)

	err = m.db.InsertHostData(hd)
	require.NoError(m.T(), err)

	list, err = m.db.FindOldCurrentHostnames(utils.MAX_TIME)
	require.NoError(m.T(), err)
	assert.Equal(m.T(), []string{"rac1-x"}, list)
}

func (m *MongodbSuite) TestGetCurrentHostnames() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.T().Run("should_return_empty", func(t *testing.T) {
		list, err := m.db.GetCurrentHostnames()
		require.NoError(t, err)
		assert.Equal(t, []string{}, list)
	})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_01.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_04.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_06.json"))

	m.T().Run("should_return_all_current_hosts", func(t *testing.T) {
		list, err := m.db.GetCurrentHostnames()
		require.NoError(t, err)
		assert.Equal(t, []string{"test-small2", "test-small3"}, list)
	})
}

func (m *MongodbSuite) TestFindOldCurrentHostnames() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_01.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_04.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_06.json"))

	m.T().Run("should_return_all_current_hosts", func(t *testing.T) {
		list, err := m.db.FindOldCurrentHostnames(utils.MAX_TIME)
		require.NoError(t, err)
		assert.Equal(t, []string{"test-small2", "test-small3"}, list)
	})

	m.T().Run("should_return_no_current_hosts", func(t *testing.T) {
		list, err := m.db.FindOldCurrentHostnames(utils.MIN_TIME)
		require.NoError(t, err)
		assert.Equal(t, []string{}, list)
	})

	m.T().Run("should_return_all_current_hosts2", func(t *testing.T) {
		list, err := m.db.FindOldCurrentHostnames(utils.P("2020-04-24T13:42:46+00:00"))
		require.NoError(t, err)
		assert.Equal(t, []string{"test-small2", "test-small3"}, list)
	})

	m.T().Run("should_return_no_current_hosts2", func(t *testing.T) {
		list, err := m.db.FindOldCurrentHostnames(utils.P("2020-04-24T10:55:49+00:00"))
		require.NoError(t, err)
		assert.Equal(t, []string{}, list)
	})

	m.T().Run("should_return_only_test_small2", func(t *testing.T) {
		list, err := m.db.FindOldCurrentHostnames(utils.P("2020-04-24T12:00:49+00:00"))
		require.NoError(t, err)
		assert.Equal(t, []string{"test-small2"}, list)
	})
}

func (m *MongodbSuite) TestFindOldCurrentHostdata() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_dataservice_mongohostdata_01.json")
	m.InsertHostData(hd)

	hd2 := mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_dataservice_mongohostdata_03.json")
	m.InsertHostData(hd2)

	m.T().Run("Should not find any", func(t *testing.T) {
		hosts, err := m.db.FindOldCurrentHostdata(utils.P("2019-11-18T16:38:58Z"))
		require.NoError(t, err)
		assert.Empty(t, hosts)
	})

	m.T().Run("Should find one", func(t *testing.T) {
		hosts, err := m.db.FindOldCurrentHostdata(utils.P("2019-12-04T16:38:58Z"))
		require.NoError(t, err)

		assert.Len(t, hosts, 1)

		assert.Equal(m.T(), "itl-csllab-112.sorint.localpippo", hosts[0].Hostname)
		assert.False(m.T(), hosts[0].Archived)
		assert.Equal(m.T(), utils.P("2019-11-19T16:38:58Z"), hosts[0].CreatedAt)
	})

	m.T().Run("Should find two", func(t *testing.T) {
		hosts, err := m.db.FindOldCurrentHostdata(utils.P("2020-01-14T15:38:58Z"))
		require.NoError(t, err)

		assert.Len(t, hosts, 2)

		sort.Slice(hosts, func(i, j int) bool {
			return hosts[i].ID.String() < hosts[j].ID.String()
		})

		assert.Equal(m.T(), "itl-csllab-112.sorint.localpippo", hosts[0].Hostname)
		assert.False(m.T(), hosts[0].Archived)
		assert.Equal(m.T(), utils.P("2019-11-19T16:38:58Z"), hosts[0].CreatedAt)

		assert.Equal(m.T(), "itl-csllab-223.sorint.localpippo", hosts[1].Hostname)
		assert.False(m.T(), hosts[1].Archived)
		assert.Equal(m.T(), utils.P("2020-01-13T15:38:58Z"), hosts[1].CreatedAt)
	})
}

func (m *MongodbSuite) TestFindOldArchivedHost() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_01.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_04.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_05.json"))

	m.T().Run("should_return_all_archived_hosts", func(t *testing.T) {
		list, err := m.db.FindOldArchivedHosts(utils.MAX_TIME)
		require.NoError(t, err)
		assert.Equal(t, []primitive.ObjectID{utils.Str2oid("5ea2d15320d55cbdc35022b1"), utils.Str2oid("5ea2d3c920d55cbdc35022b7")}, list)
	})

	m.T().Run("should_return_no_archived_hosts", func(t *testing.T) {
		list, err := m.db.FindOldArchivedHosts(utils.MIN_TIME)
		require.NoError(t, err)
		assert.Equal(t, []primitive.ObjectID{}, list)
	})

	m.T().Run("should_return_all_archived_hosts2", func(t *testing.T) {
		list, err := m.db.FindOldArchivedHosts(utils.P("2020-04-24T12:00:53+00:00"))
		require.NoError(t, err)
		assert.Equal(t, []primitive.ObjectID{utils.Str2oid("5ea2d15320d55cbdc35022b1"), utils.Str2oid("5ea2d3c920d55cbdc35022b7")}, list)
	})

	m.T().Run("should_return_no_archived_hosts2", func(t *testing.T) {
		list, err := m.db.FindOldArchivedHosts(utils.P("2020-04-24T10:00:49+00:00"))
		require.NoError(t, err)
		assert.Equal(t, []primitive.ObjectID{}, list)
	})

	m.T().Run("should_return_only_test_small1", func(t *testing.T) {
		list, err := m.db.FindOldArchivedHosts(utils.P("2020-04-24T11:55:49+00:00"))
		require.NoError(t, err)
		assert.Equal(t, []primitive.ObjectID{utils.Str2oid("5ea2d15320d55cbdc35022b1")}, list)
	})
}

func (m *MongodbSuite) TestDeleteHostData() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_01.json"))

	list, err := m.db.FindOldArchivedHosts(utils.MAX_TIME)
	require.NoError(m.T(), err)
	require.Equal(m.T(), []primitive.ObjectID{utils.Str2oid("5ea2d15320d55cbdc35022b1")}, list)

	err = m.db.DeleteHostData(utils.Str2oid("5ea2d15320d55cbdc35022b1"))
	require.NoError(m.T(), err)

	list, err = m.db.FindOldArchivedHosts(utils.MAX_TIME)
	require.NoError(m.T(), err)
	require.Equal(m.T(), []primitive.ObjectID{}, list)
}

func (m *MongodbSuite) TestGetHostnames() {
	m.T().Run("should_return_all", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
		m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_01.json"))
		m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_04.json"))
		m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_06.json"))

		list, err := m.db.GetHostnames()
		require.NoError(t, err)
		assert.Equal(t, []string{"test-small", "test-small2", "test-small3"}, list)
	})

	m.T().Run("should_return_empty", func(t *testing.T) {
		list, err := m.db.GetHostnames()
		require.NoError(t, err)
		assert.Equal(t, []string{}, list)
	})

}
