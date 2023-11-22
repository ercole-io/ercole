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
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSearchPostgreSqlInstances_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedContent := []dto.PostgreSqlInstance{
		{
			Hostname:    "test-db",
			Name:        "PostgreSQL-example:1010",
			Charset:     "UTF8",
			Version:     "PostgreSQL 10.20",
			Environment: "PRD",
			Location:    "Germany",
		},
		{
			Hostname:    "test-db",
			Name:        "PostgreSQL-example:1010",
			Charset:     "UTF8",
			Version:     "PostgreSQL 10.20",
			Environment: "PRD",
			Location:    "Germany",
		},
	}

	var expectedRes = dto.PostgreSqlInstanceResponse{
		Content: expectedContent,
		Metadata: dto.PagingMetadata{
			Empty:         false,
			First:         true,
			Last:          true,
			Number:        0,
			Size:          1,
			TotalElements: 1,
			TotalPages:    0,
		},
	}

	db.EXPECT().SearchPostgreSqlInstances(
		[]string{"foo", "bar", "foobarx"}, "Hostname",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(&expectedRes, nil).Times(1)

	res, err := as.SearchPostgreSqlInstances(
		dto.SearchPostgreSqlInstancesFilter{
			dto.GlobalFilter{
				"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
			},
			"foo bar foobarx", "Hostname",
			true, 1, 1,
		},
	)

	require.NoError(t, err)
	assert.Equal(t, &expectedRes, res)
}

func TestSearchPostgreSqlInstances_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchPostgreSqlInstances(
		[]string{"foo", "bar", "foobarx"}, "Memory",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchPostgreSqlInstances(

		dto.SearchPostgreSqlInstancesFilter{
			dto.GlobalFilter{
				"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
			},
			"foo bar foobarx", "Memory",
			true, 1, 1,
		},
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}
