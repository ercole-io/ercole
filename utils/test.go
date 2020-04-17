package utils

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/amreo/ercole-services/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//P parse the string s and return the equivalent time
func P(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

//Str2priTime parse the string s and return the equivalent bson primitive time
func Str2priTime(s string) primitive.DateTime {
	time := P(s)
	return primitive.NewDateTimeFromTime(time)
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

//LoadFixtureHostDataMap load the hostdata in the filename and return it
func LoadFixtureHostDataMap(t *testing.T, filename string) model.HostDataMap {
	raw, err := ioutil.ReadFile(filename)
	require.NoError(t, err)

	var hd map[string]interface{}
	require.NoError(t, bson.UnmarshalExtJSON(raw, true, &hd))

	return model.HostDataMap(hd)
}
