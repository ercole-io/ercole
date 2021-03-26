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
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/jtblin/go-ldap-client"
	"github.com/sirupsen/logrus"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// LDAPAuthenticationProvider is the concrete implementation of AuthenticationProvider that provide a LDAP user authentication.
type LDAPAuthenticationProvider struct {
	// Config contains the dataservice global configuration
	Config config.AuthenticationProviderConfig
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Log contains logger formatted
	Log *logrus.Logger
	// privateKey contains the private key used to sign the JWT tokens
	privateKey interface{}
	// publicKey contains the public key used to check the JWT tokens
	publicKey interface{}
	Client    *ldap.LDAPClient
}

// Init initializes the service and database
func (ap *LDAPAuthenticationProvider) Init() {
	raw, err := ioutil.ReadFile(ap.Config.PrivateKey)
	if err != nil {
		ap.Log.Panic(err)
	}

	ap.privateKey, ap.publicKey, err = utils.ParsePrivateKey(raw)
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

	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: ap.privateKey}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		ap.Log.Panic(err)
	}

	cl := jwt.Claims{
		Subject:   info["Username"].(string),
		Issuer:    "ercole",
		NotBefore: jwt.NewNumericDate(ap.TimeNow()),
		Audience:  jwt.Audience{info["Username"].(string)},
		ID:        info["Username"].(string),
		Expiry:    jwt.NewNumericDate(ap.TimeNow().Add(time.Duration(ap.Config.TokenValidityTimeout) * time.Second)),
		IssuedAt:  jwt.NewNumericDate(ap.TimeNow()),
	}
	raw, err := jwt.Signed(sig).Claims(cl).CompactSerialize()
	if err != nil {
		ap.Log.Panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(raw))
}

// AuthenticateMiddleware return the middleware used to check if the users are authenticated
func (ap *LDAPAuthenticationProvider) AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Get the token
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewError(errors.New("You don't have setted the authorization header"), http.StatusText(http.StatusUnauthorized)))
			return
		}

		//Check token type
		if strings.HasPrefix(tokenStr, "Basic ") {
			tokenStr = tokenStr[len("Basic "):]
			val, err := base64.StdEncoding.DecodeString(tokenStr)
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
		} else if strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = tokenStr[len("Bearer "):]
			//Parse the token
			parsed, err := jwt.ParseSigned(tokenStr)
			if err != nil {
				utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewError(err, http.StatusText(http.StatusUnauthorized)))
				return
			}

			//Validate the token and get the claim
			claim := jwt.Claims{}
			err = parsed.Claims(ap.publicKey, &claim)
			if err != nil {
				utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewError(err, http.StatusText(http.StatusUnauthorized)))
				return
			}

			//Check exp field
			if claim.Expiry.Time().Before(ap.TimeNow()) {
				utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewError(errors.New("The token is expired"), http.StatusText(http.StatusUnauthorized)))
				return
			}

			//Check issuedat field
			if claim.IssuedAt.Time().After(ap.TimeNow()) {
				utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewError(errors.New("Futuristic tokens (from future) are invalid"), http.StatusText(http.StatusUnauthorized)))
				return
			}

			//Serve the request
			next.ServeHTTP(w, r)
		} else {
			utils.WriteAndLogError(ap.Log, w, http.StatusUnauthorized, utils.NewError(errors.New("The authorization header value doesn't begin with Basic or Bearer"), http.StatusText(http.StatusUnauthorized)))
			return
		}
	})
}
