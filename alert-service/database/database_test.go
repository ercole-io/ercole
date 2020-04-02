// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"log"
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

func (m *MongodbSuite) TestFindHostData_SuccessExist() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := utils.LoadFixtureHostData(m.T(), "../../fixture/test_dataservice_mongohostdata_01.json")
	err := m.InsertHostData(hd)
	require.NoError(m.T(), err)

	hd2, err := m.db.FindHostData(hd.ID)
	require.NoError(m.T(), err)
	log.Println(hd2.CreatedAt)

	assert.Equal(m.T(), hd, hd2)
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
	hd.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	hd.CreatedAt = utils.P("2019-11-05T14:02:03Z")
	err := m.InsertHostData(hd)
	require.NoError(m.T(), err)

	m.ArchiveHost(hd.Hostname)
	hd.Archived = true

	afterFirstInsert := utils.P("2019-11-05T15:02:03Z")
	hd3, err3 := m.db.FindMostRecentHostDataOlderThan("itl-csllab-112.sorint.localpippo", afterFirstInsert)
	require.NoError(m.T(), err3)
	assert.Equal(m.T(), hd, hd3)

	hd2 := utils.LoadFixtureHostData(m.T(), "../../fixture/test_dataservice_mongohostdata_02.json")
	hd2.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	hd2.CreatedAt = utils.P("2019-11-06T14:02:03Z")
	err2 := m.InsertHostData(hd2)
	require.NoError(m.T(), err2)

	assert.NotEqual(m.T(), hd, hd2)

	afterSecondInsert := utils.P("2019-11-06T15:02:03Z")
	hd4, err4 := m.db.FindMostRecentHostDataOlderThan(hd.Hostname, afterSecondInsert)
	require.NoError(m.T(), err4)
	assert.Equal(m.T(), hd2, hd4)
}
