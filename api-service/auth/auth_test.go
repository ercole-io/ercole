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
	"testing"

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildAuthenticationProvider_NotSupported(t *testing.T) {
	testConf := config.AuthenticationProviderConfig{
		Type: "foobar",
	}

	assert.PanicsWithValue(t, "The AuthenticationProvider type wasn't recognized or supported", func() {
		BuildAuthenticationProvider(testConf, nil, nil)
	})
}

func TestBuildAuthenticationProvider_Basic(t *testing.T) {
	testConf := config.AuthenticationProviderConfig{
		Type:       "basic",
		Username:   "foobar",
		Password:   "F0oB4r",
		PrivateKey: "/tmp/path/to/private.key",
		PublicKey:  "/tmp/path/to/public.pem",
	}
	logger := utils.NewLogger("TEST")
	time := utils.Btc(utils.P("2019-11-05T14:02:03Z"))

	ap := BuildAuthenticationProvider(testConf, time, logger)

	require.IsType(t, &BasicAuthenticationProvider{}, ap)
	bap, _ := ap.(*BasicAuthenticationProvider)

	assert.Same(t, logger, bap.Log)
	utils.AssertFuncAreTheSame(t, time, bap.TimeNow)
	assert.Equal(t, testConf, bap.Config)
}
