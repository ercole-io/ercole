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

// Package service is a package that provides methods for querying data
package auth

import (
	"net/http"

	"github.com/amreo/ercole-services/utils"
)

// AuthenticationProvider is a interface that wrap methods used to authenticate users
type AuthenticationProvider interface {
	// Init initialize the provider
	Init()
	// AuthenticateMiddleware return the middleware used to check if the users are authenticated
	AuthenticateMiddleware(next http.Handler) http.Handler
	// TokenEndpoint return the middleware used to check if the users are authenticated
	GetToken(w http.ResponseWriter, r *http.Request)
	// GetUserInfoIfCorrect return the informations about the user if the provided credentials are correct, otherwise return nil
	GetUserInfoIfCredentialsAreCorrect(username string, password string) (map[string]interface{}, utils.AdvancedErrorInterface)
}
