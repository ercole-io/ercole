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
	"os"
	"testing"
	"time"

	"github.com/ercole-io/ercole/config"
	migration "github.com/ercole-io/ercole/database-migration"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"math/rand"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MongodbSuite struct {
	suite.Suite

	db     MongoDatabase
	dbname string
	log    *logrus.Logger
}

func (db *MongodbSuite) SetupSuite() {
	val, ok := os.LookupEnv("MONGODB_URI")
	rand.Seed(time.Now().UnixNano())
	db.db = MongoDatabase{
		Config: config.Configuration{
			Mongodb: config.Mongodb{
				URI:    val,
				DBName: fmt.Sprintf("ercole_test_%d", rand.Int()),
			},
		},
		TimeNow: func() time.Time { return utils.P("2020-12-05T14:02:03Z") },
	}
	if !ok {
		db.db.Config.Mongodb.URI = "mongodb://127.0.0.1:27017"
	}
	fmt.Println("DBNAME:", db.db.Config.Mongodb.DBName)
	db.db.Config.Mongodb.URI += "/" + db.db.Config.Mongodb.DBName
	db.dbname = db.db.Config.Mongodb.DBName

	//Migrations
	cl := migration.ConnectToMongodb(db.log, db.db.Config.Mongodb)
	migration.Migrate(db.log, cl.Database(db.db.Config.Mongodb.DBName))
	cl.Disconnect(context.TODO())

	db.db.ConnectToMongodb()
}

func TestMongodbSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test for mongodb database(alert-service)")
	}

	mongodbHandlerSuiteTest := &MongodbSuite{}

	suite.Run(t, mongodbHandlerSuiteTest)
}

func (db *MongodbSuite) TearDownSuite() {
	db.db.Client.Database(db.db.Config.Mongodb.DBName).Drop(context.TODO())
	db.db.Client.Disconnect(context.TODO())
}

func (db *MongodbSuite) InsertHostData(hostData model.RawObject) {
	_, err := db.db.Client.Database(db.dbname).Collection("hosts").InsertOne(context.TODO(), hostData)
	db.Require().NoError(err)
}

func (db *MongodbSuite) RunTestQuery(testName string, query bson.A, check func(out []map[string]interface{})) {
	db.Run(testName, func() {
		cur, err := db.db.Client.Database(db.db.Config.Mongodb.DBName).Collection("hosts").Aggregate(
			context.TODO(),
			query,
		)
		require.NoError(db.T(), err)

		var out []map[string]interface{} = make([]map[string]interface{}, 0)
		require.NoError(db.T(), cur.All(context.TODO(), &out))

		check(out)
	})
}

func (db *MongodbSuite) InsertPatchingFunction(pf model.PatchingFunction) {
	_, err := db.db.Client.Database(db.dbname).Collection("patching_functions").InsertOne(context.TODO(), pf)
	db.Require().NoError(err)
}
