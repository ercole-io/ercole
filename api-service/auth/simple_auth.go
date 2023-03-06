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

package auth

import (
	"bytes"
	"crypto/rsa"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/api-service/service"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/context"
)

// BasicAuthenticationProvider is the concrete implementation of AuthenticationProvider that provide a simple user authentication.
type BasicAuthenticationProvider struct {
	// Config contains the dataservice global configuration
	Config config.AuthenticationProviderConfig
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Log contains logger formatted
	Log logger.Logger
	// privateKey contains the private key used to sign the JWT tokens
	privateKey *rsa.PrivateKey
	// publicKey contains the public key used to check the JWT tokens
	publicKey *rsa.PublicKey
	// Service contains the underlying service used to perform various logical and store operations
	Service service.APIService
}

// Init initializes the service and database
func (ap *BasicAuthenticationProvider) Init() {
	raw, err := os.ReadFile(ap.Config.PrivateKey)
	if err != nil {
		ap.Log.Panic(err)
	}

	ap.privateKey, ap.publicKey, err = parsePrivateKey(raw)
	if err != nil {
		ap.Log.Panic(utils.NewErrorf("Unable to parse the private key: %s", err))
	}
}

// GetUserInfoIfCredentialsAreCorrect return the informations about the user if the provided credentials are correct, otherwise return nil
func (ap *BasicAuthenticationProvider) GetUserInfoIfCredentialsAreCorrect(username string, password string) (*dto.User, error) {
	user, err := ap.Service.GetUser(username, BasicType)
	if err != nil {
		return nil, err
	}

	if ap.Service.MatchPassword(user, password) {
		userDto := dto.ToUser(user)
		return &userDto, nil
	}

	return nil, utils.ErrInvalidUser
}

// GetToken return the middleware used to check if the users are authenticated
func (ap *BasicAuthenticationProvider) GetToken(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var request LoginRequest

	//Parse the request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteAndLogError(ap.Log, w, http.StatusBadRequest, utils.NewError(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	//Check if the credentials are valid
	userInfo, err := ap.GetUserInfoIfCredentialsAreCorrect(request.Username, request.Password)
	if err != nil {
		utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, err)
		return
	}

	if userInfo == nil {
		utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewError(errors.New("Failed to login, invalid credentials"), http.StatusText(http.StatusUnauthorized)))
		return
	}

	token, err := buildToken(ap.TimeNow(), ap.Config.TokenValidityTimeout, *userInfo, ap.privateKey)
	if err != nil {
		ap.Log.Errorf("Unable to get signed token: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteAndLogError(ap.Log, w, http.StatusInternalServerError, fmt.Errorf("Unable to get signed token"))

		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, dto.ToLoginResponse(token, userInfo))
}

// AuthenticateMiddleware return the middleware used to check if the users are authenticated
func (ap *BasicAuthenticationProvider) AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewError(errors.New("You don't have setted the authorization header"), http.StatusText(http.StatusUnauthorized)))
			return
		}

		if strings.HasPrefix(tokenString, "Basic ") {
			tokenString = tokenString[len("Basic "):]
			val, err := base64.StdEncoding.DecodeString(tokenString)
			if err != nil {
				utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewError(err, http.StatusText(http.StatusUnauthorized)))
				return
			}

			if !bytes.ContainsAny(val, ":") {
				utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewError(errors.New("A : is missing in the auth header"), http.StatusText(http.StatusUnauthorized)))
				return
			}

			user := val[:bytes.IndexRune(val, ':')]
			password := val[bytes.IndexRune(val, ':')+1:]

			if subtle.ConstantTimeCompare(user, []byte(ap.Config.Username)) == 0 || subtle.ConstantTimeCompare(password, []byte(ap.Config.Password)) == 0 {
				utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewError(errors.New("Invalid credentials"), http.StatusText(http.StatusUnauthorized)))
				return
			}

			next.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(tokenString, "Bearer ") {
			claims, err := validateBearerToken(tokenString, ap.TimeNow, ap.publicKey)
			if err != nil {
				utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.ErrInvalidToken)
				return
			}

			if claims == nil {
				utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.ErrInvalidToken)
				return
			}

			ercoleGroups := ap.Service.GetMatchedGroupsName(claims.Groups)

			context.Set(r, "user", model.User{Username: claims.Subject, Groups: ercoleGroups})

			next.ServeHTTP(w, r)
			return
		}

		utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewErrorf("The authorization header value doesn't begin with Basic or Bearer"))
	})
}

func (ap *BasicAuthenticationProvider) GetType() string {
	return BasicType
}
