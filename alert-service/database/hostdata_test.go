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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package database

import (
	"context"
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestFindHostData_SuccessExist() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_dataservice_mongohostdata_01.json")

	m.InsertHostData(hd)

	hd2, err := m.db.FindHostData(utils.Str2oid("5dd41aa2b2fa40163c878538"))
	require.NoError(m.T(), err)

	assert.Equal(m.T(), "itl-csllab-112.sorint.localpippo", hd2.Hostname)
	assert.False(m.T(), hd2.Archived)
	assert.Equal(m.T(), utils.P("2019-11-19T16:38:58Z"), hd2.CreatedAt)
}

func (m *MongodbSuite) TestFindHostData_FailWrongID() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_dataservice_mongohostdata_01.json")
	m.InsertHostData(hd)

	notExistingID := utils.Str2oid("8a46027b2ddab34ed01a8c56")

	hd2, err := m.db.FindHostData(notExistingID)
	require.Error(m.T(), err)

	assert.Equal(m.T(), model.HostDataBE{}, hd2)
}

func (m *MongodbSuite) TestFindMostRecentHostDataOlderThan_OnlyOne() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_dataservice_mongohostdata_01.json")
	m.InsertHostData(hd)

	foundHd, err := m.db.FindMostRecentHostDataOlderThan("itl-csllab-112.sorint.localpippo", utils.P("2020-04-25T11:45:23Z"))
	require.NoError(m.T(), err)
	assert.Equal(m.T(), "itl-csllab-112.sorint.localpippo", foundHd.Hostname)
	assert.False(m.T(), foundHd.Archived)
	assert.Equal(m.T(), utils.P("2019-11-19T16:38:58Z"), foundHd.CreatedAt)
}

func (m *MongodbSuite) TestFindMostRecentHostDataOlderThan_MoreThanOne() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	hd := mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_dataservice_mongohostdata_04.json")
	m.InsertHostData(hd)

	m.T().Run("Should_find_hd_even_if_archived", func(t *testing.T) {
		foundHd, err := m.db.FindMostRecentHostDataOlderThan("itl-csllab-112.sorint.localpippo", utils.P("2019-11-20T15:38:58Z"))
		require.NoError(t, err)
		assert.True(t, foundHd.Archived)
		assert.Equal(t, "itl-csllab-112.sorint.localpippo", foundHd.Hostname)
		assert.Equal(t, utils.P("2019-11-19T15:38:58Z"), foundHd.CreatedAt)
	})

	hd2 := mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_dataservice_mongohostdata_02.json")
	m.InsertHostData(hd2)
	assert.NotEqual(m.T(), hd, hd2)

	m.T().Run("Should_find_hd2_more_inserts", func(t *testing.T) {
		foundHd, err := m.db.FindMostRecentHostDataOlderThan("itl-csllab-112.sorint.localpippo", utils.P("2019-12-20T15:38:58Z"))
		require.NoError(t, err)
		assert.False(t, foundHd.Archived)
		assert.Equal(t, "itl-csllab-112.sorint.localpippo", foundHd.Hostname)
		assert.Equal(t, utils.P("2019-12-19T15:38:58Z"), foundHd.CreatedAt)
	})
}
