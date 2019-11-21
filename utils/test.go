package utils

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/amreo/ercole-services/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//P parse the string s and return the equivalent time
func P(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

//Btc break the time continuum and return a function that return the time t
func Btc(t time.Time) func() time.Time {
	return func() time.Time {
		return t
	}
}

//Str2oid parse the objectid and return the parsed value
func Str2oid(str string) primitive.ObjectID {
	val, _ := primitive.ObjectIDFromHex(str)
	return val
}

//LoadFixtureHostData load the hostdata in the filename and return it
func LoadFixtureHostData(t *testing.T, filename string) model.HostData {
	var hd model.HostData
	raw, err := ioutil.ReadFile(filename)

	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, &hd))

	return hd
}
