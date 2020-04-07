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
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/stretchr/testify/suite"
)

func TestMongodbSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test for mongodb database(alert-service)")
	}

	mongodbHandlerSuiteTest := &MongodbSuite{}

	suite.Run(t, mongodbHandlerSuiteTest)
}

func TestConnectToMongodb_FailToConnect(t *testing.T) {
	logger := utils.NewLogger("TEST")
	logger.ExitFunc = func(int) {
		panic("log.Fatal called by test")
	}

	db := MongoDatabase{
		Config: config.Configuration{
			Mongodb: config.Mongodb{
				URI:    "wronguri:1234/test",
				DBName: fmt.Sprintf("ercole_test_%d", rand.Int()),
			},
		},
		Log: logger,
	}

	assert.PanicsWithValue(t, "log.Fatal called by test", db.ConnectToMongodb)
}

func (m *MongodbSuite) TestFindHostData_SuccessExist() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := utils.LoadFixtureHostData(m.T(), "../../fixture/test_dataservice_mongohostdata_01.json")
	err := m.InsertHostData(hd)
	require.NoError(m.T(), err)

	hd2, err := m.db.FindHostData(hd.ID)
	require.NoError(m.T(), err)

	assert.Equal(m.T(), hd, hd2)
}

func (m *MongodbSuite) TestFindHostData_FailWrongID() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := utils.LoadFixtureHostData(m.T(), "../../fixture/test_dataservice_mongohostdata_01.json")
	err := m.InsertHostData(hd)
	require.NoError(m.T(), err)

	notExistingID := primitive.NewObjectIDFromTimestamp(time.Now())

	hd2, err := m.db.FindHostData(notExistingID)
	require.Error(m.T(), err)

	assert.Equal(m.T(), model.HostData{}, hd2)
}

func (m *MongodbSuite) TestFindMostRecentHostDataOlderThan_OneInsert_Success() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := utils.LoadFixtureHostData(m.T(), "../../fixture/test_dataservice_hostdata_01.json")
	hd.CreatedAt = time.Now().UTC().Truncate(time.Millisecond)
	err := m.InsertHostData(hd)
	require.NoError(m.T(), err)

	time.Sleep(1 * time.Second)
	afterFirstInsert := time.Now()

	hd2, err2 := m.db.FindMostRecentHostDataOlderThan("itl-csllab-112.sorint.localpippo", afterFirstInsert)
	require.NoError(m.T(), err2)
	assert.Equal(m.T(), hd, hd2)
}

func (m *MongodbSuite) TestFindMostRecentHostDataOlderThan_MoreInserts_Success() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := utils.LoadFixtureHostData(m.T(), "../../fixture/test_dataservice_mongohostdata_01.json")
	err := m.InsertHostData(hd)
	require.NoError(m.T(), err)

	aTimeAfterInsert := hd.CreatedAt.AddDate(0, 0, 1)
	m.T().Run("Should find hd", func(t *testing.T) {
		foundHd, err := m.db.FindMostRecentHostDataOlderThan("itl-csllab-112.sorint.localpippo", aTimeAfterInsert)
		require.NoError(m.T(), err)
		assert.Equal(m.T(), hd, foundHd)
	})

	m.ArchiveHost(hd.Hostname)
	hd.Archived = true

	m.T().Run("Should find hd even if archived", func(t *testing.T) {
		foundHd, err := m.db.FindMostRecentHostDataOlderThan("itl-csllab-112.sorint.localpippo", aTimeAfterInsert)
		require.NoError(m.T(), err)
		assert.Equal(m.T(), hd, foundHd)
	})

	hd2 := utils.LoadFixtureHostData(m.T(), "../../fixture/test_dataservice_mongohostdata_02.json")
	err2 := m.InsertHostData(hd2)
	require.NoError(m.T(), err2)
	assert.NotEqual(m.T(), hd, hd2)

	aTimeAfterInsert = hd2.CreatedAt.AddDate(0, 0, 1)
	m.T().Run("Should find hd2", func(t *testing.T) {
		foundHd, err := m.db.FindMostRecentHostDataOlderThan(hd.Hostname, aTimeAfterInsert)
		require.NoError(m.T(), err)
		assert.Equal(m.T(), hd2, foundHd)
	})
}

func (m *MongodbSuite) TestInsertAlert_Success() {
	_, err := m.db.InsertAlert(alert1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("alerts").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("alerts").FindOne(context.TODO(), bson.M{
		"_id": alert1.ID,
	})
	require.NoError(m.T(), val.Err())

	var out model.Alert
	val.Decode(&out)

	assert.Equal(m.T(), alert1, out)
}

func (m *MongodbSuite) TestFindOldCurrentHosts() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := utils.LoadFixtureHostData(m.T(), "../../fixture/test_dataservice_mongohostdata_01.json")
	hd.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	hd.CreatedAt = utils.P("2019-11-05T14:02:03Z")
	err := m.InsertHostData(hd)
	require.NoError(m.T(), err)

	hd2 := utils.LoadFixtureHostData(m.T(), "../../fixture/test_dataservice_mongohostdata_03.json")
	hd2.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	hd2.CreatedAt = utils.P("2020-01-05T14:02:03Z")
	err2 := m.InsertHostData(hd2)
	require.NoError(m.T(), err2)

	m.T().Run("Should not find any", func(t *testing.T) {
		hosts, err := m.db.FindOldCurrentHosts(utils.P("2019-10-05T14:02:03Z"))
		require.NoError(m.T(), err)
		assert.Empty(m.T(), hosts)
	})

	m.T().Run("Should find one", func(t *testing.T) {
		hosts, err := m.db.FindOldCurrentHosts(utils.P("2019-12-05T14:02:03Z"))
		require.NoError(m.T(), err)

		assert.Len(m.T(), hosts, 1)
		expectedHosts := append(make([]string, 0), hd.Info.Hostname)
		assert.ElementsMatch(m.T(), expectedHosts, hosts)
	})

	m.T().Run("Should find two", func(t *testing.T) {
		hosts, err := m.db.FindOldCurrentHosts(utils.P("2020-02-05T14:02:03Z"))
		require.NoError(m.T(), err)

		assert.Len(m.T(), hosts, 2)
		expectedHosts := append(make([]string, 0), hd.Info.Hostname, hd2.Info.Hostname)
		assert.ElementsMatch(m.T(), expectedHosts, hosts)
	})
}

func (m *MongodbSuite) TestExistNoDataAlert_SuccessNotExist() {
	_, err := m.db.InsertAlert(alert1)
	defer m.db.Client.Database(m.dbname).Collection("alerts").DeleteMany(context.TODO(), bson.M{})
	require.NoError(m.T(), err)

	exist, err := m.db.ExistNoDataAlertByHost("myhost")
	require.NoError(m.T(), err)

	assert.False(m.T(), exist)
}

func (m *MongodbSuite) TestExistNoDataAlert_SuccessExist() {
	_, err := m.db.InsertAlert(alert2)
	defer m.db.Client.Database(m.dbname).Collection("alerts").DeleteMany(context.TODO(), bson.M{})
	require.NoError(m.T(), err)

	exist, err := m.db.ExistNoDataAlertByHost("myhost")
	require.NoError(m.T(), err)

	assert.True(m.T(), exist)
}
