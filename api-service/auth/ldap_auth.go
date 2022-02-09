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

package auth

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jtblin/go-ldap-client"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
)

// LDAPAuthenticationProvider is the concrete implementation of AuthenticationProvider that provide a LDAP user authentication.
type LDAPAuthenticationProvider struct {
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
	Client    *ldap.LDAPClient
}

// Init initializes the service and database
func (ap *LDAPAuthenticationProvider) Init() {
	raw, err := ioutil.ReadFile(ap.Config.PrivateKey)
	if err != nil {
		ap.Log.Fatal(err)
	}

	ap.privateKey, ap.publicKey, err = parsePrivateKey(raw)
	if err != nil {
		ap.Log.Panic(err)
	}

	ap.Client = &ldap.LDAPClient{
		Base:         ap.Config.LDAPBase,
		Host:         ap.Config.Host,
		Port:         ap.Config.Port,
		UseSSL:       ap.Config.LDAPUseSSL,
		BindDN:       ap.Config.LDAPBindDN,
		BindPassword: ap.Config.LDAPBindPassword,
		UserFilter:   ap.Config.LDAPUserFilter,
		GroupFilter:  ap.Config.LDAPGroupFilter,
		Attributes:   []string{"givenName", "sn", "mail", "uid"},
	}

	if ap.Client.Connect() != nil {
		ap.Log.Fatalf("Error connecting LDAP %v", err)
	}
}

// GetUserInfoIfCredentialsAreCorrect return the informations about the user if the provided credentials are correct, otherwise return nil
func (ap *LDAPAuthenticationProvider) GetUserInfoIfCredentialsAreCorrect(username string, password string) (map[string]interface{}, error) {
	ok, _, err := ap.Client.Authenticate(username, password)
	if err != nil {
		return nil, utils.NewError(err, "AUTH")
	}
	if !ok {
		return nil, nil
	}

	return map[string]interface{}{
		"Username": username,
	}, nil
}

// GetToken return the middleware used to check if the users are authenticated
func (ap *LDAPAuthenticationProvider) GetToken(w http.ResponseWriter, r *http.Request) {
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
	info, err := ap.GetUserInfoIfCredentialsAreCorrect(request.Username, request.Password)
	if err != nil {
		utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, err)
		return
	}
	if info == nil {
		utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewError(errors.New("Failed to login, invalid credentials"), http.StatusText(http.StatusUnauthorized)))
	}

	token, err := buildToken(ap.TimeNow(), ap.Config.TokenValidityTimeout, request.Username, ap.privateKey)
	if err != nil {
		ap.Log.Errorf("Unable to get signed token: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.WriteAndLogError(ap.Log, w, http.StatusInternalServerError, fmt.Errorf("Unable to get signed token"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

// AuthenticateMiddleware return the middleware used to check if the users are authenticated
func (ap *LDAPAuthenticationProvider) AuthenticateMiddleware(next http.Handler) http.Handler {
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

			info, err := ap.GetUserInfoIfCredentialsAreCorrect(string(user), string(password))
			if err != nil {
				utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, err)
				return
			}
			if info == nil {
				utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewError(errors.New("Failed to login, invalid credentials"), http.StatusText(http.StatusUnauthorized)))
			}

			next.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(tokenString, "Bearer ") {
			err := validateBearerToken(tokenString, ap.TimeNow, ap.publicKey)
			if err != nil {
				ap.Log.Debugf("Invalid token: %s", err)
				utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, fmt.Errorf("Invalid token"))
				return
			}

			next.ServeHTTP(w, r)
			return
		}

		utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewError(errors.New("The authorization header value doesn't begin with Basic or Bearer"), http.StatusText(http.StatusUnauthorized)))
	})
}
