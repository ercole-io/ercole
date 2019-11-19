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
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/amreo/ercole-services/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var alert1 model.Alert = model.Alert{
	AlertCode:     model.AlertCodeNewServer,
	AlertSeverity: model.AlertSeverityNotice,
	AlertStatus:   model.AlertStatusNew,
	Date:          p("2019-11-05T18:02:03Z"),
	Description:   "pippo",
	OtherInfo: map[string]interface{}{
		"hostname": "myhost",
	},
	ID: str2oid("5dd40bfb12f54dfda7b1c291"),
}

var alert2 model.Alert = model.Alert{
	AlertCode:     model.AlertCodeNoData,
	AlertSeverity: model.AlertSeverityMajor,
	AlertStatus:   model.AlertStatusNew,
	Date:          p("2019-11-05T18:02:03Z"),
	Description:   "No data received from the host myhost in the last 90 days",
	OtherInfo: map[string]interface{}{
		"hostname": "myhost",
	},
	ID: str2oid("5dd4113f0085a6fac03c4fed"),
}

//p parse the string s and return the equivalent time
func p(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

//btc break the time continuum and return a function that return the time t
func btc(t time.Time) func() time.Time {
	return func() time.Time {
		return t
	}
}

func str2oid(str string) primitive.ObjectID {
	val, _ := primitive.ObjectIDFromHex(str)
	return val
}

func loadFixtureHostData(t *testing.T, filename string) model.HostData {
	var hd model.HostData
	raw, err := ioutil.ReadFile(filename)

	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, &hd))

	return hd
}
