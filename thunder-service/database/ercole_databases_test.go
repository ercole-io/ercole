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
	"time"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestGetErcoleDatabases_Success() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.db.TimeNow = utils.Btc(utils.P("2022-03-10T17:38:03Z"))
	defer func() { m.db.TimeNow = time.Now }()

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_thunderservice_mongohostdata_11.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_thunderservice_mongohostdata_12.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_thunderservice_mongohostdata_13.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_thunderservice_mongohostdata_09.json"))

	out, err := m.db.GetErcoleDatabases()
	m.Require().NoError(err)

	expectedOut := []model.ErcoleDatabase{
		{
			Hostname: "test-db",
			Info: model.Info{
				CpuThreads: 2,
			},
			Features: model.Feature{
				Oracle: model.Oracle{
					Database: model.Database{
						Databases: []model.Object{
							{
								Name:       "ERCOLE",
								UniqueName: "ERCOLE",
								Work:       1,
							},
						},
					},
				},
			},
			Archived: false,
		},
		{
			Hostname: "test-db",
			Info: model.Info{
				CpuThreads: 2,
			},
			Features: model.Feature{
				Oracle: model.Oracle{
					Database: model.Database{
						Databases: []model.Object{
							{
								Name:       "ERCOLE",
								UniqueName: "ERCOLE",
								Work:       1,
							},
						},
					},
				},
			},
			Archived: true,
		},
		{
			Hostname: "test-db",
			Info: model.Info{
				CpuThreads: 2,
			},
			Features: model.Feature{
				Oracle: model.Oracle{
					Database: model.Database{
						Databases: []model.Object{
							{
								Name:       "ERCOLE",
								UniqueName: "ERCOLE",
								Work:       1,
							},
						},
					},
				},
			},
			Archived: true,
		},
	}
	assert.Equal(m.T(), expectedOut, out)
}

func (m *MongodbSuite) TestGetErcoleActiveDatabases_Success() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})

	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_thunderservice_mongohostdata_01.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_thunderservice_mongohostdata_02.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_thunderservice_mongohostdata_03.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_thunderservice_mongohostdata_04.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_thunderservice_mongohostdata_07.json"))

	out, err := m.db.GetErcoleActiveDatabases()
	m.Require().NoError(err)

	expectedOut := []model.ErcoleDatabase{
		{
			Hostname: "test-db",
			Info: model.Info{
				CpuThreads: 2,
			},
			Features: model.Feature{
				Oracle: model.Oracle{
					Database: model.Database{
						Databases: []model.Object{
							{
								Name:       "ERCOLE",
								UniqueName: "ERCOLE",
								Work:       1,
							},
						},
					},
				},
			},
			Archived: false,
		},
		{
			Hostname: "test-small",
			Info: model.Info{
				CpuThreads: 2,
			},
			Archived: false,
		},
		{
			Hostname: "test-small2",
			Info: model.Info{
				CpuThreads: 2,
			},
			Archived: false,
		},
	}
	assert.Equal(m.T(), expectedOut, out)
}
