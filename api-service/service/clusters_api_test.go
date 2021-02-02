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
	"strconv"
	"testing"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchClusters_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []map[string]interface{}{
		{
			"CPU":                         0,
			"Environment":                 "PROD",
			"Hostname":                    "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			"HostnameAgentVirtualization": "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			"Location":                    "Italy",
			"Name":                        "not_in_cluster",
			"VirtualizationNodes":         "aspera-b1fe49e8501c9ef031e5acff4b5e69a9",
			"Sockets":                     0,
			"Type":                        "unknown",
			"_id":                         utils.Str2oid("5e8c234b24f648a08585bd3d"),
		},
		{
			"CPU":                         140,
			"Environment":                 "PROD",
			"Hostname":                    "test-virt",
			"HostnameAgentVirtualization": "test-virt",
			"Location":                    "Italy",
			"Name":                        "Puzzait",
			"VirtualizationNodes":         "s157-cb32c10a56c256746c337e21b3f82402",
			"Sockets":                     10,
			"Type":                        "vmware",
			"_id":                         utils.Str2oid("5e8c234b24f648a08585bd41"),
		},
	}

	db.EXPECT().SearchClusters(
		false, []string{"foo", "bar", "foobarx"}, "CPU",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchClusters(
		false, "foo bar foobarx", "CPU",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchClusters_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchClusters(
		false, []string{"foo", "bar", "foobarx"}, "CPU",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchClusters(
		false, "foo bar foobarx", "CPU",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetClusterXLSX(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	cluster := &dto.Cluster{
		ID:                          utils.Str2oid("5eb0222a45d85f4193704944"),
		CPU:                         140,
		CreatedAt:                   utils.P("2020-05-04T14:09:46.608Z"),
		Environment:                 "PROD",
		FetchEndpoint:               "???",
		Hostname:                    "test-virt",
		HostnameAgentVirtualization: "test-virt",
		Location:                    "Italy",
		Name:                        "Puzzait",
		Sockets:                     10,
		Type:                        "vmware",
		VirtualizationNodes:         []string{"s157-cb32c10a56c256746c337e21b3f82402"},
		VirtualizationNodesCount:    1,
		VirtualizationNodesStats: []dto.VirtualizationNodesStat{
			{
				TotalVMsCount:                   2,
				TotalVMsWithErcoleAgentCount:    1,
				TotalVMsWithoutErcoleAgentCount: 1,
				VirtualizationNode:              "s157-cb32c10a56c256746c337e21b3f82402"}},

		VMs: []dto.VM{
			{
				CappedCPU:          false,
				Hostname:           "test-virt",
				Name:               "test-virt",
				VirtualizationNode: "s157-cb32c10a56c256746c337e21b3f82402"},

			{
				CappedCPU:          true,
				Hostname:           "test-db",
				Name:               "test-db",
				VirtualizationNode: "s157-cb32c10a56c256746c337e21b3f82402"},
		},
		VMsCount:            2,
		VMsErcoleAgentCount: 1,
	}

	var clusterName string = "pippo"
	var olderThan = utils.P("2019-11-05T14:02:03Z")

	db.EXPECT().
		GetCluster(clusterName, olderThan).
		Return(cluster, nil)

	xlsx, err := as.GetClusterXLSX(clusterName, olderThan)
	assert.NoError(t, err)

	sheet := "VMs"

	i := -1
	columns := []string{"A", "B", "C", "D"}
	nextVal := func() string {
		i += 1
		cell := columns[i%4] + strconv.Itoa(i/4+1)
		return xlsx.GetCellValue(sheet, cell)
	}

	assert.Equal(t, "Physical Hosts", nextVal())
	assert.Equal(t, "Hostname", nextVal())
	assert.Equal(t, "VirtualizationNode", nextVal())
	assert.Equal(t, "CappedCPU", nextVal())

	assert.Equal(t, "test-virt", nextVal())
	assert.Equal(t, "test-virt", nextVal())
	assert.Equal(t, "s157-cb32c10a56c256746c337e21b3f82402", nextVal())
	assert.Equal(t, "false", nextVal())

	assert.Equal(t, "test-db", nextVal())
	assert.Equal(t, "test-db", nextVal())
	assert.Equal(t, "s157-cb32c10a56c256746c337e21b3f82402", nextVal())
	assert.Equal(t, "true", nextVal())

	assert.Equal(t, "", nextVal())
	assert.Equal(t, "", nextVal())
	assert.Equal(t, "", nextVal())
	assert.Equal(t, "", nextVal())
}
