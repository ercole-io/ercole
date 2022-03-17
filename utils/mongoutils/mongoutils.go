// Copyright (c) 2022 Sorint.lab S.p.A.
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

package mongoutils

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/model"
)

//LoadFixtureHostData load the hostdata in the filename and return it
func LoadFixtureHostData(t *testing.T, filename string) model.HostDataBE {
	var hd model.HostDataBE
	raw, err := ioutil.ReadFile(filename)

	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, &hd))

	return hd
}

//LoadFixtureHostDataMap load the hostdata in the filename and return it
func LoadFixtureHostDataMap(t *testing.T, filename string) model.RawObject {
	raw, err := ioutil.ReadFile(filename)
	require.NoError(t, err)

	var hd model.RawObject
	require.NoError(t, bson.UnmarshalExtJSON(raw, true, &hd))

	return hd
}

//LoadFixtureMongoHostDataMap load the mongohostdata in the filename and return it
func LoadFixtureMongoHostDataMap(t *testing.T, filename string) model.RawObject {
	raw, err := ioutil.ReadFile(filename)
	require.NoError(t, err)

	var hd model.RawObject
	require.NoError(t, bson.UnmarshalExtJSON(raw, true, &hd))

	return hd
}

//LoadFixtureMongoHostDataMapAsHostData load the mongohostdata in the filename and return it as hostdata
func LoadFixtureMongoHostDataMapAsHostData(t *testing.T, filename string) model.HostDataBE {
	raw, err := ioutil.ReadFile(filename)
	require.NoError(t, err)

	var hd model.HostDataBE
	require.NoError(t, bson.UnmarshalExtJSON(raw, true, &hd))

	return hd
}
