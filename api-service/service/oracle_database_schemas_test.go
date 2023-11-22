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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListOracleDatabaseSchemas_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := []dto.OracleDatabaseSchema{
		{
			Hostname:      "hostname",
			DatabaseName:  "databasename",
			Indexes:       0,
			LOB:           0,
			Tables:        0,
			Total:         0,
			User:          "user",
			AccountStatus: "status",
		},
	}

	expectedPDB := []dto.OracleDatabaseSchema{
		{
			Hostname:      "hostname",
			DatabaseName:  "databasename",
			Indexes:       0,
			LOB:           0,
			Tables:        0,
			Total:         0,
			User:          "user",
			AccountStatus: "status",
			Pdb:           "",
		},
	}

	db.EXPECT().FindAllOracleDatabaseSchemas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).Return(expected, nil)
	db.EXPECT().FindAllOraclePDBSchemas(dto.GlobalFilter{OlderThan: utils.MAX_TIME}).Return(expectedPDB, nil)

	expectedRes := make([]dto.OracleDatabaseSchema, 0)
	expectedRes = append(expectedRes, expected...)
	expectedRes = append(expectedRes, expectedPDB...)

	res, err := as.ListOracleDatabaseSchemas(dto.GlobalFilter{OlderThan: utils.MAX_TIME})
	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}
