package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//P parse the string s and return the equivalent time
func P(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

//PDT parse the string s and return the equivalent bson primitive time
func PDT(s string) primitive.DateTime {
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

// AssertFuncAreTheSame tests if funcExpected is the same of funcActual
func AssertFuncAreTheSame(t *testing.T, funcExpected interface{}, funcActual interface{}) {
	funcExpectedAddress := runtime.FuncForPC(reflect.ValueOf(funcExpected).Pointer()).Name()
	funcActualAddress := runtime.FuncForPC(reflect.ValueOf(funcActual).Pointer()).Name()
	assert.Equal(t, funcExpectedAddress, funcActualAddress)
}

// NewObjectIDForTests is a function to replace NewObjectID in tests that return ids increasing
func NewObjectIDForTests() func() primitive.ObjectID {
	i := 0
	return func() primitive.ObjectID {
		i++
		objID := fmt.Sprintf("%024d", i)

		return Str2oid(objID)
	}
}
