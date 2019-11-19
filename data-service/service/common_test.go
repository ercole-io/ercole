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

package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/amreo/ercole-services/model"

	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source ../database/database.go -destination=fake_database.go -package=service
//go:generate mockgen -source service.go -destination=fake_service.go -package=service

//Common data
var errMock error = errors.New("MockError")
var aerrMock utils.AdvancedErrorInterface = utils.NewAdvancedErrorPtr(errMock, "mock")

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

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func NewHTTPTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
