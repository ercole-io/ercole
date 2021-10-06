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
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"

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
	logger := logger.NewLogger("TEST", logger.SetExitFunc(
		func(int) {
			panic("log.Fatal called by test")
		},
	))

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
