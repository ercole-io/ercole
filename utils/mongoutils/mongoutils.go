package mongoutils

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
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
