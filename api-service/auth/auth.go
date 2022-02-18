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
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	jwt "github.com/golang-jwt/jwt/v4"
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
	GetUserInfoIfCredentialsAreCorrect(username string, password string) (map[string]interface{}, error)
}

// BuildAuthenticationProvider return a authentication provider that match what is requested in the configuration
// It's initialized
func BuildAuthenticationProvider(conf config.AuthenticationProviderConfig, timeNow func() time.Time, log logger.Logger) AuthenticationProvider {
	switch conf.Type {
	case "basic":
		prov := new(BasicAuthenticationProvider)
		prov.Config = conf
		prov.Log = log
		prov.TimeNow = timeNow

		return prov
	case "ldap":
		prov := new(LDAPAuthenticationProvider)
		prov.Config = conf
		prov.Log = log
		prov.TimeNow = timeNow

		return prov
	default:
		panic("The AuthenticationProvider type wasn't recognized or supported")
	}
}

func buildToken(now time.Time, tokenValidityTimeout int, username string, privateKey *rsa.PrivateKey) (string, error) {
	if privateKey == nil {
		return "", fmt.Errorf("privateKey is nil")
	}

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(tokenValidityTimeout) * time.Second)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    "ercole",
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	ss, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func parsePrivateKey(raw []byte) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(raw)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

func validateBearerToken(tokenString string, timeNow func() time.Time, publicKey *rsa.PublicKey) error {
	tokenString = tokenString[len("Bearer "):]
	jwt.TimeFunc = timeNow
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(_ *jwt.Token) (interface{}, error) {
		return publicKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}))

	if err != nil {
		return err
	}

	_, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("Invalid token")
	}

	return nil
}
