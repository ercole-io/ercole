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

package service

import (
	"testing"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOracleChanges_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := []dto.OracleChangesDto{
		{
			Hostname: "newdb",
			Databasenames: []dto.OracleChangesDBs{
				{
					Databasename: "pippodb",
					OracleChanges: []dto.OracleChangesGrowth{
						{
							DailyCPUUsage: 3.4,
							SegmentsSize:  50,
							Updated:       utils.P("2020-05-21T09:32:54.83Z"),
							DatafileSize:  8,
							Allocable:     129,
						},
						{
							DailyCPUUsage: 5.3,
							SegmentsSize:  100,
							Updated:       utils.P("2020-05-21T09:32:09.288Z"),
							DatafileSize:  10,
							Allocable:     129,
						},
						{
							DailyCPUUsage: 0.7,
							SegmentsSize:  3,
							Updated:       utils.P("2020-05-21T09:30:55.061Z"),
							DatafileSize:  6,
							Allocable:     129,
						},
					},
				},
			},
		},
	}
	db.EXPECT().FindOracleChangesByHostname(dto.GlobalFilter{OlderThan: utils.MAX_TIME}, "newdb").Return(expected, nil)

	res, err := as.GetOracleChanges(dto.GlobalFilter{OlderThan: utils.MAX_TIME}, "newdb")
	require.NoError(t, err)
	assert.Equal(t, expected, res)
}
