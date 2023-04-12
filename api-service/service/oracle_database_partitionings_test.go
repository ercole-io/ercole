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

func TestListOracleDatabasePartitionings_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := []dto.OracleDatabasePartitioning{
		{
			Hostname:      "hostname",
			DatabaseName:  "databasename",
			Owner:         "ownername",
			SegmentName:   "segmentname",
			Count: 1,
			Mb:            100,
		},
	}

	expectedPDB := []dto.OracleDatabasePartitioning{
		{
			Hostname:      "hostname",
			DatabaseName:  "databasename",
			Pdb:           "pdbname",
			Owner:         "ownername",
			SegmentName:   "segmentname",
			Count: 1,
			Mb:            100,
		},
	}

	db.EXPECT().FindAllOracleDatabasePartitionings(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).Return(expected, nil)
	db.EXPECT().FindAllOraclePDBPartitionings(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).Return(expectedPDB, nil)

	expectedRes := make([]dto.OracleDatabasePartitioning, 0)
	expectedRes = append(expectedRes, expected...)
	expectedRes = append(expectedRes, expectedPDB...)

	res, err := as.ListOracleDatabasePartitionings(dto.GlobalFilter{OlderThan: utils.MAX_TIME})
	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}
