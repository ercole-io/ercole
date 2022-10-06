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
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	r "math/rand"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiservice_database "github.com/ercole-io/ercole/v2/api-service/database"
	apiservice_service "github.com/ercole-io/ercole/v2/api-service/service"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

var testRSAPrivateKey string = `
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAsxbIblAtTWazN2FkI6n+Gd4wCSMiXQEelxyM/8Zyo/j+dWcO
8A65rZIoKIwHwXL3+CdthCMRJKj5wfgVm7qqNn4tRG9UzLRYM0+Ks3xdq2Q6sxJk
pSDAD4I87uin3YNJqOFQJvsG9diZU3xUNnzXFiRAvsAnd9bvprQeFyR12ZtpECMY
T4zgEue68CCzzcAvN4S8Zq+5hOu5bh2GNbUlamHP1tJpNBN60ZiJMXCFU3AE9Nli
9BOEyi1YdqIVngXdMhbBHNfbWUWBplGTCuWGLpVMwIUUaI8rGn7q1hYxcAK2Fi77
ZP/im75RCyKzb9xzYr2d0VufWn5DNcmLfZpgewIDAQABAoIBAG0+I4su+0Nwtze7
/+LFakwLPdAFD4weB7Pz5YqMWhft5gJlmDYVNWxMcJSzPnPhlqNYIbTt0yJCtP9+
PmgdSIEvHJvXMaohBIBgL+JmpZjL7gaX3K7huGZ9cn/liahU0pTixArTK57Bvl2v
xIrsQiOuf5QcELdIdC2DR6ukQQM4hCjzGbBLPbci1r+sBLaMGskvtEuvAKavlAKO
bjLbt8MCQqBMsQRQSYYRDUIOVuqbfEYkYiq3gpVI6GPyRX7rLXLYTi58ngsHpdiK
V4Qypz+YjDitFm1prcp3piZBPz9D/BRpA1wwbmOn/5AVDHjzJ3aDQL2TFveotpcZ
ecefxJkCgYEA7KLSyaYUNhU6f4o9f8Hc3YTN7KM3s6wNhszzjnAMuwq4FoD9l30G
YhO6EgGQ8FWzqJ4WBbB2Dpi1RaQ3w1RrYQb9VfQY+NOJmZlvOfZZSD5Uev5v3KvO
masS5y9osfJX/a9f3svHtqD1GV7E01TBXEUHr8FbBRg8Q4J6wSoOgXUCgYEAwb5v
EqgvRXhMGGYK7eWZ/QogrJXsh3+Cb1kwrJC+qhKm9BkpKp4JU2+A4vDqOZ+urUIb
aBRM434a69ZWJ0XTzvYfilPEUXSlAfb+xGPERXRu0xhWqAKGq83+volpDV+9sTUu
Dd5mVQ1i31ckgFN16uAJ4xGIrwpWJKnmfOlerC8CgYBxKA9aJBPoJNCbapSsAh1G
xJngPdCGF5FEU79n7ob37lFHWZlqlnu17K7+q0cO1jyaNjZbtB1QL5AHZFbSDg1n
EXuVXauPWUCkda2tbvMUy9GEGyWMxY9/BkJ80Lvk0/llszZKCPJQj7mEzz+Zux7X
q57YWcLXtdYjhkKDGkRjfQKBgGt+BdA7IecQRF/xFbVCAzrCSLiYgd/3nd27hWbo
8/AWYyzhXNa5UgFJxx+ifMG118tm9x+6y6IYUEVy6N/nPQoBwiQUL8LlzbsWV+mM
VNQYMnjKcyHKLP/bTbBXOsLh0LQmBkRJlUsxHx89ERJlu/Gxlaq3CrfbK0oyPaAm
NpGfAoGAf7OAh+0DKK+yUmBbsj72ppfl3kK7SM+RvxEiFmaH7B77z1kjUziDuEl+
J0NPlNaa7q04kMXv5TY/SZbThZWUL0Wyz7nYJgEIBswKvkz305q8QaHM/k7Rrr1g
pGjyq0RrXmoLoKLqsB+d+jvcXvk0kaOJgcX50C+hUp5AxtQU/bU=
-----END RSA PRIVATE KEY-----
`

var invalidTestRSAKey string = `
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAsxbIblA@tTWazN2FkI6n+Gd4wCSMiXQEelxyM/8Zyo/j+dWcO
pGjyq0RrXmoLoKLqsB+d###+jvcXvk0kaOJgcX50C+hUp5AxtQU/bU=
-----END RSA PRIVATE KEY-----
`

var correctLoginRequest string = `
{
	"username": "foobar",
	"password": "C0rr3ctP4ssw0rd"
}`
var invalidLoginRequest string = `
{
	foo bar
}`
var incorrectLoginRequest string = `
{
	"username": "foobar",
	"password": "N0tC0rrectP4ssW0rd"
}`

var validToken string = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiZm9vYmFyIl0sImV4cCI6MTU3Mjk2MjU0MywiaWF0IjoxNTcyOTYyNTIzLCJpc3MiOiJlcmNvbGUiLCJqdGkiOiJmb29iYXIiLCJuYmYiOjE1NzI5NjI1MjMsInN1YiI6ImZvb2JhciJ9.Ki3LzVISVOgW_tNnLumDRbLEef2dAL4nTAcK0jLSgwPj1nQErkXi9YM1Fs2r2yByZ0gJLsj9ok3kxPMp1_qdwqqrvgdiDqj-hcdrL76C0poWkystVJgxtjFgi683zYyfjvLjPCqQnNuI8icIGJTogzjZoKjJg7KPeqHjBr0Qsb2epMzD-Cj5ir4PGM3LYCAX6ui7gxGMIGg-zrqcMvQ8EIa5Kn_Xp4bVGjjkvOxpql7ivGE9BlRggVHRjX9iUtD_aC6WfgKHyQDmyECtB_A0jMK3mHUBoXc28axkML9i4AXE3Nw7z_qA47HDeVY99CrvQhFRy30JCxJh_pmRbyETzw"
var invalidToken string = "sdfsdf.sdfs.df-sdf-sdfssdfdsdff-"
var tokenWithInvalidSignature string = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiZm9vYmFyIl0sImV4cCI6MTU3Mjk2MjU0MywiaWF0IjoxNTcyOTYyNTIzLCJpc3MiOiJlcmNvbGUiLCJqdGkiOiJmb29iYXIiLCJuYmYiOjE1NzI5NjI1MjMsInN1YiI6ImZvb2JhciJ9.QuK03C6rKVF0WU8GhmYwhz40FVoL1IsDEfatep-KbjS8SJBw4OojOJqfyF5Vpeu5AqJvaOqFuoQ1fjGA9Yjhk0F7TlCl-LJHE80Dlrj0W4cR1BJ2u8Mf7-xaMmTe0FQt7x12WTr04DlKfTHuBkO2__DDJDwYzUuJoNSlJbTczMs"

var ercoleConfig = config.Configuration{
	Mongodb: config.Mongodb{
		DBName: fmt.Sprintf("ercole_test_%d", r.Int()),
		URI:    "mongodb://127.0.0.1:27017",
	},
}

var db = &apiservice_database.MongoDatabase{
	Config:  ercoleConfig,
	TimeNow: time.Now,
}

var serviceAuth = &apiservice_service.APIService{
	Database: db,
}

func TestGetUserInfoIfCredentialsAreCorrect_WhenAreCredentialsAreWrong(t *testing.T) {
	db.ConnectToMongodb()

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username: "foobar",
			Password: "C0rr3ctP4ssw0rd",
		},
		Service: *serviceAuth,
	}
	res, err := bap.GetUserInfoIfCredentialsAreCorrect("foobar", "password")
	require.ErrorContains(t, err, "User not found")
	assert.Nil(t, res)
}

func TestGetUserInfoIfCredentialsAreCorrect_WhenAreCredentialsAreCorrect(t *testing.T) {
	db.ConnectToMongodb()

	defer serviceAuth.RemoveUser("foobar")

	serviceAuth.AddUser(
		model.User{
			Username: "foobar",
			Password: "C0rr3ctP4ssw0rd",
			Groups:   []string{"Test"},
		},
	)

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username: "foobar",
			Password: "C0rr3ctP4ssw0rd",
		},
		Service: *serviceAuth,
	}

	res, err := bap.GetUserInfoIfCredentialsAreCorrect("foobar", "C0rr3ctP4ssw0rd")
	require.NoError(t, err)
	assert.Equal(t, "foobar", res["Username"])

	db.Client.Database(db.Config.Mongodb.DBName).Drop(context.TODO())
	db.Client.Disconnect(context.TODO())
}

func TestInit_OK(t *testing.T) {
	// Create a temporary private key
	f, err := ioutil.TempFile("/tmp/", "ercole-*")
	require.NoError(t, err)

	_, err = f.WriteString(testRSAPrivateKey)
	require.NoError(t, err)

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			PrivateKey: f.Name(),
		},
	}

	bap.Init()
	assert.NotNil(t, bap.privateKey)
	assert.NotNil(t, bap.publicKey)

	f.Close()
	require.NoError(t, os.Remove(f.Name()))
}

func TestInit_NoFile(t *testing.T) {
	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			PrivateKey: "/tmp/path/to/atlantis",
		},
	}

	assert.Panics(t, bap.Init)
}

func TestInit_InvalidFile(t *testing.T) {
	// Create a temporary private key
	f, err := ioutil.TempFile("/tmp/", "ercole-*")
	require.NoError(t, err)

	_, err = f.WriteString(invalidTestRSAKey)
	require.NoError(t, err)

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			PrivateKey: f.Name(),
		},
	}

	assert.Panics(t, bap.Init)
}

func TestGetToken_OK(t *testing.T) {
	var err error

	db.ConnectToMongodb()

	defer serviceAuth.RemoveUser("foobar")

	serviceAuth.AddUser(
		model.User{
			Username: "foobar",
			Password: "C0rr3ctP4ssw0rd",
			Groups:   []string{"Test"},
		},
	)

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:     logger.NewLogger("TEST"),
		Service: *serviceAuth,
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(testRSAPrivateKey))
	require.NoError(t, err)

	bap.privateKey = privateKey
	bap.publicKey = &privateKey.PublicKey

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(bap.GetToken)
	req, err := http.NewRequest("POST", "/user/login", strings.NewReader(correctLoginRequest))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var claims jwt.RegisteredClaims

	jwt.TimeFunc = bap.TimeNow
	_, err = jwt.ParseWithClaims(rr.Body.String(), &claims, func(_ *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})
	require.NoError(t, err)

	assert.Equal(t, jwt.NewNumericDate(utils.P("2019-11-05T14:02:23Z").Local()), claims.ExpiresAt)
	assert.Equal(t, jwt.NewNumericDate(utils.P("2019-11-05T14:02:03Z").Local()), claims.IssuedAt)
	assert.Equal(t, jwt.NewNumericDate(utils.P("2019-11-05T14:02:03Z").Local()), claims.NotBefore)
	assert.Equal(t, "ercole", claims.Issuer)
	assert.Equal(t, "foobar", claims.Subject)

	db.Client.Database(db.Config.Mongodb.DBName).Drop(context.TODO())
	db.Client.Disconnect(context.TODO())
}

func TestGetToken_InvalidRequest(t *testing.T) {
	var err error

	db.ConnectToMongodb()

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:     logger.NewLogger("TEST"),
		Service: *serviceAuth,
	}

	bap.privateKey, bap.publicKey, err = parsePrivateKey([]byte(testRSAPrivateKey))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(bap.GetToken)
	req, err := http.NewRequest("POST", "/user/login", strings.NewReader(invalidLoginRequest))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetToken_InvalidCredentials(t *testing.T) {
	var err error

	db.ConnectToMongodb()

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:     logger.NewLogger("TEST"),
		Service: *serviceAuth,
	}

	bap.privateKey, bap.publicKey, err = parsePrivateKey([]byte(testRSAPrivateKey))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(bap.GetToken)
	req, err := http.NewRequest("POST", "/user/login", strings.NewReader(incorrectLoginRequest))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetToken_InvalidKeys(t *testing.T) {
	var err error

	db.ConnectToMongodb()

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:     logger.NewLogger("TEST"),
		Service: *serviceAuth,
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(bap.GetToken)
	req, err := http.NewRequest("POST", "/user/login", strings.NewReader(correctLoginRequest))
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthenticateMiddleware_NoAuthorizationHeader(t *testing.T) {
	var err error

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:     logger.NewLogger("TEST"),
	}

	bap.privateKey, bap.publicKey, err = parsePrivateKey([]byte(testRSAPrivateKey))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := bap.AuthenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(222)
	}))
	req, err := http.NewRequest("GET", "/myping", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthenticateMiddleware_WrongAuthorizationHeader(t *testing.T) {
	var err error

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:     logger.NewLogger("TEST"),
	}

	bap.privateKey, bap.publicKey, err = parsePrivateKey([]byte(testRSAPrivateKey))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := bap.AuthenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(222)
	}))
	req, err := http.NewRequest("GET", "/myping", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Foobar 78378943789239")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthenticateMiddleware_BasicInvalidBase64(t *testing.T) {
	var err error

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:     logger.NewLogger("TEST"),
	}

	bap.privateKey, bap.publicKey, err = parsePrivateKey([]byte(testRSAPrivateKey))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := bap.AuthenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(222)
	}))
	req, err := http.NewRequest("GET", "/myping", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Basic !!!!")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthenticateMiddleware_BasicMissingColon(t *testing.T) {
	var err error

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:     logger.NewLogger("TEST"),
	}

	bap.privateKey, bap.publicKey, err = parsePrivateKey([]byte(testRSAPrivateKey))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := bap.AuthenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(222)
	}))
	req, err := http.NewRequest("GET", "/myping", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Basic Zm9vYmFy")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthenticateMiddleware_BasicInvalidCredentials(t *testing.T) {
	var err error

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:     logger.NewLogger("TEST"),
	}

	bap.privateKey, bap.publicKey, err = parsePrivateKey([]byte(testRSAPrivateKey))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := bap.AuthenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(222)
	}))
	req, err := http.NewRequest("GET", "/myping", nil)
	require.NoError(t, err)
	req.SetBasicAuth("foobar", "NotC0rr3ctP4ssw0rd")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthenticateMiddleware_BasicOk(t *testing.T) {
	var err error

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:     logger.NewLogger("TEST"),
	}

	bap.privateKey, bap.publicKey, err = parsePrivateKey([]byte(testRSAPrivateKey))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := bap.AuthenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(222)
	}))
	req, err := http.NewRequest("GET", "/myping", nil)
	require.NoError(t, err)
	req.SetBasicAuth("foobar", "C0rr3ctP4ssw0rd")

	handler.ServeHTTP(rr, req)

	require.Equal(t, 222, rr.Code)
}

func TestAuthenticateMiddleware_BearerInvalidToken(t *testing.T) {
	var err error

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:     logger.NewLogger("TEST"),
	}

	bap.privateKey, bap.publicKey, err = parsePrivateKey([]byte(testRSAPrivateKey))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := bap.AuthenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(222)
	}))
	req, err := http.NewRequest("GET", "/myping", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+invalidToken)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthenticateMiddleware_BearerInvalidSignature(t *testing.T) {
	var err error

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:     logger.NewLogger("TEST"),
	}

	bap.privateKey, bap.publicKey, err = parsePrivateKey([]byte(testRSAPrivateKey))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := bap.AuthenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(222)
	}))
	req, err := http.NewRequest("GET", "/myping", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+tokenWithInvalidSignature)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthenticateMiddleware_BearerTokenExpired(t *testing.T) {
	var err error

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T16:02:03Z")),
		Log:     logger.NewLogger("TEST"),
	}

	bap.privateKey, bap.publicKey, err = parsePrivateKey([]byte(testRSAPrivateKey))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := bap.AuthenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(222)
	}))
	req, err := http.NewRequest("GET", "/myping", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+validToken)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthenticateMiddleware_BearerTokenFromFuture(t *testing.T) {
	var err error

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T12:02:03Z")),
		Log:     logger.NewLogger("TEST"),
	}

	bap.privateKey, bap.publicKey, err = parsePrivateKey([]byte(testRSAPrivateKey))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := bap.AuthenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(222)
	}))
	req, err := http.NewRequest("GET", "/myping", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+validToken)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthenticateMiddleware_BearerOk(t *testing.T) {
	var err error

	bap := BasicAuthenticationProvider{
		Config: config.AuthenticationProviderConfig{
			Username:             "foobar",
			Password:             "C0rr3ctP4ssw0rd",
			TokenValidityTimeout: 20,
		},
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:     logger.NewLogger("TEST", logger.LogVerbosely(true)),
	}

	bap.privateKey, bap.publicKey, err = parsePrivateKey([]byte(testRSAPrivateKey))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := bap.AuthenticateMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(222)
	}))
	req, err := http.NewRequest("GET", "/myping", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+validToken)

	handler.ServeHTTP(rr, req)

	require.Equal(t, 222, rr.Code)
}
