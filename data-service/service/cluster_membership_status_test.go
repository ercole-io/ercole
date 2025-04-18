// Copyright (c) 2025 Sorint.lab S.p.A.
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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetVeritasHostsFqdn(t *testing.T) {
	testCases := []struct {
		name          string
		hostname      string
		clusterHosts  []string
		checkResponse func(t *testing.T, got []string)
	}{
		{
			name:         "With FQDN",
			hostname:     "host01.fqdn.test.domain",
			clusterHosts: []string{"host02", "host01"},
			checkResponse: func(t *testing.T, got []string) {
				require.Len(t, got, 2)
				require.Equal(t, "host02.fqdn.test.domain", got[0])
				require.Equal(t, "host01.fqdn.test.domain", got[1])
			},
		},
		{
			name:         "Without FQDN",
			hostname:     "host01",
			clusterHosts: []string{"host02", "host01"},
			checkResponse: func(t *testing.T, got []string) {
				require.Len(t, got, 2)
				require.Equal(t, "host02", got[0])
				require.Equal(t, "host01", got[1])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hds := HostDataService{}

			got := hds.getVeritasHostsFqdn(tc.hostname, tc.clusterHosts)
			tc.checkResponse(t, got)
		})
	}

}
