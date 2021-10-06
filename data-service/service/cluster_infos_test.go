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
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestCheckClusterInfos(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)

	hds := HostDataService{
		Config:         config.Configuration{},
		ServerVersion:  "1.6.6",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:            logger.NewLogger("TEST"),
	}

	testCases := []struct {
		name      string
		hostnames []string
		expected  []model.ClusterInfo
		actual    []model.ClusterInfo
	}{
		{
			name:      "Empty",
			hostnames: []string{},
			expected:  []model.ClusterInfo{},
			actual:    []model.ClusterInfo{},
		},
		{
			name:      "No match",
			hostnames: []string{"qui.paperopolis.dn", "quo.paperopolis.dn", "qua.paperopolis.dn"},
			expected: []model.ClusterInfo{
				{
					FetchEndpoint: "",
					Type:          "",
					Name:          "",
					CPU:           0,
					Sockets:       0,
					VMs: []model.VMInfo{
						{
							Name:               "",
							Hostname:           "PIPPO",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
						{
							Name:               "",
							Hostname:           "PLUTO",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
					},
				},
			},
			actual: []model.ClusterInfo{
				{
					FetchEndpoint: "",
					Type:          "",
					Name:          "",
					CPU:           0,
					Sockets:       0,
					VMs: []model.VMInfo{
						{
							Name:               "",
							Hostname:           "PIPPO",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
						{
							Name:               "",
							Hostname:           "PLUTO",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
					},
				},
			},
		},
		{
			name:      "One match",
			hostnames: []string{"pippo.topolinia.dn", "topolino.topolinia.dn"},
			expected: []model.ClusterInfo{
				{
					FetchEndpoint: "",
					Type:          "",
					Name:          "",
					CPU:           0,
					Sockets:       0,
					VMs: []model.VMInfo{
						{
							Name:               "",
							Hostname:           "pippo.topolinia.dn",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
						{
							Name:               "",
							Hostname:           "PLUTO",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
					},
				},
			},
			actual: []model.ClusterInfo{
				{
					FetchEndpoint: "",
					Type:          "",
					Name:          "",
					CPU:           0,
					Sockets:       0,
					VMs: []model.VMInfo{
						{
							Name:               "",
							Hostname:           "PIPPO",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
						{
							Name:               "",
							Hostname:           "PLUTO",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
					},
				},
			},
		},
		{
			name: "Multiple match",
			hostnames: []string{"qui.paperopolis.dn", "quo.paperopolis.dn",
				"pippo.topolinia.dn", "topolino.topolinia.dn"},
			expected: []model.ClusterInfo{
				{
					FetchEndpoint: "",
					Type:          "",
					Name:          "",
					CPU:           0,
					Sockets:       0,
					VMs: []model.VMInfo{
						{
							Name:               "",
							Hostname:           "pippo.topolinia.dn",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
						{
							Name:               "",
							Hostname:           "qui.paperopolis.dn",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
						{
							Name:               "",
							Hostname:           "PLUTO",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
					},
				},
			},
			actual: []model.ClusterInfo{
				{
					FetchEndpoint: "",
					Type:          "",
					Name:          "",
					CPU:           0,
					Sockets:       0,
					VMs: []model.VMInfo{
						{
							Name:               "",
							Hostname:           "PIPPO",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
						{
							Name:               "",
							Hostname:           "PLUTO",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
						{
							Name:               "",
							Hostname:           "qui",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
					},
				},
			},
		},
		{
			name: "Multiple match with UPPER CASES, Camel Case, not fully qualified names",
			hostnames: []string{"qui.paperopolis.dn", "quo.paperopolis.dn",
				"pippo.topolinia.dn", "TOPOLINO"},
			expected: []model.ClusterInfo{
				{
					FetchEndpoint: "",
					Type:          "",
					Name:          "",
					CPU:           0,
					Sockets:       0,
					VMs: []model.VMInfo{
						{
							Name:               "",
							Hostname:           "pippo.topolinia.dn",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
						{
							Name:               "",
							Hostname:           "qui.paperopolis.dn",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
						{
							Name:               "",
							Hostname:           "TOPOLINO",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
						{
							Name:               "",
							Hostname:           "PLUTO",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
					},
				},
			},
			actual: []model.ClusterInfo{
				{
					FetchEndpoint: "",
					Type:          "",
					Name:          "",
					CPU:           0,
					Sockets:       0,
					VMs: []model.VMInfo{
						{
							Name:               "",
							Hostname:           "PIPPO",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
						{
							Name:               "",
							Hostname:           "Qui.quo.Qua",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
						{
							Name:               "",
							Hostname:           "topolino.topolinia.top",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
						{
							Name:               "",
							Hostname:           "PLUTO",
							CappedCPU:          false,
							VirtualizationNode: "",
							OtherInfo:          map[string]interface{}{},
						},
					},
				},
			},
		},
	}

	for i := range testCases {
		tc := &testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			db.EXPECT().GetHostnames().Return(tc.hostnames, nil)

			hds.clusterInfoChecks(tc.actual)

			require.Equal(t, len(tc.expected), len(tc.actual))
			for i := 0; i < len(tc.actual); i++ {
				assert.ElementsMatch(t, tc.expected[i].VMs, tc.actual[i].VMs)
			}
		})
	}
}
