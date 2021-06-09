package utils

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//P parse the string s and return the equivalent time
// e.g.: 2019-11-05T14:02:03Z
// e.g.: 2019-11-05T14:02:03+01:00
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
