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

package migration

import (
	"context"
	"errors"
	"sort"
	"time"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ercole-io/ercole/v2/config"
	_ "github.com/ercole-io/ercole/v2/database-migration/migrations"
	"github.com/ercole-io/ercole/v2/utils"
)

func Migrate(conf config.Mongodb) error {
	database, err := connectToMongodb(conf)
	if err != nil {
		return err
	}

	migrate.SetDatabase(database)

	if err := migrate.Up(migrate.AllAvailable); err != nil {
		return err
	}

	err = database.Client().Disconnect(context.TODO())
	if err != nil {
		return utils.NewError(err, "Can't disconnect from the database!")
	}

	return nil
}

func connectToMongodb(conf config.Mongodb) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(conf.URI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, utils.NewError(err, "Can't connect to the database!")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, utils.NewError(err, "Can't connect to the database!")
	}

	return client.Database(conf.DBName), nil
}

func IsAtTheLatestVersion(conf config.Mongodb) (bool, error) {
	actual, latest, err := GetVersions(conf)
	if err != nil {
		return false, err
	}

	if latest < actual {
		return false, utils.NewError(errors.New("Db version is higher than last migration version"))
	} else if latest > actual {
		return false, nil
	}

	return true, nil
}

func GetVersions(conf config.Mongodb) (actual, latest uint64, err error) {
	database, err := connectToMongodb(conf)
	if err != nil {
		return
	}

	migrate.SetDatabase(database)

	migrations := migrate.RegisteredMigrations()

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	actual, _, err = migrate.Version()
	if err != nil {
		return
	}

	lastMigration := migrations[len(migrations)-1]

	return actual, lastMigration.Version, nil
}
