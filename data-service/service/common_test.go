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

package service

import (
	"errors"
	"net/http"

	"github.com/ercole-io/ercole/v2/utils"
)

//go:generate mockgen -source ../database/database.go -destination=fake_database_test.go -package=service
//go:generate mockgen -source ../../alert-service/client/client.go -destination=fake_alert_service_client_test.go -package=service
//go:generate mockgen -source ../../api-service/client/client.go -destination=fake_api_service_client_test.go -package=service

//Common data
var errMock error = errors.New("MockError")
var aerrMock error = utils.NewError(errMock, "mock")

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func NewHTTPTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
