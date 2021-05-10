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
	"time"

	"github.com/ercole-io/ercole/v2/config"
	migration "github.com/ercole-io/ercole/v2/database-migration"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"

	"math/rand"

	"github.com/stretchr/testify/suite"
)

type MongodbSuite struct {
	suite.Suite

	db     MongoDatabase
	dbname string
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
	}
	if !ok {
		db.db.Config.Mongodb.URI = "mongodb://127.0.0.1:27017"
	}
	fmt.Println("DBNAME:", db.db.Config.Mongodb.DBName)
	db.db.Config.Mongodb.URI += "/" + db.db.Config.Mongodb.DBName
	db.dbname = db.db.Config.Mongodb.DBName

	log := utils.NewLogger("test")

	err := migration.Migrate(db.db.Config.Mongodb)
	if err != nil {
		log.Fatal(err)
	}

	db.db.ConnectToMongodb()
}

func (db *MongodbSuite) TearDownSuite() {
	db.db.Client.Database(db.db.Config.Mongodb.DBName).Drop(context.TODO())
	db.db.Client.Disconnect(context.TODO())
}

func (db *MongodbSuite) InsertHostData(hostData model.RawObject) {
	_, err := db.db.Client.Database(db.dbname).Collection("hosts").InsertOne(context.TODO(), hostData)
	db.Require().NoError(err)
}
