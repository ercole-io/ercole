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
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchMongoDBInstances_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedContent := []dto.MongoDBInstance{
		{
			Hostname:     "test-db",
			Environment:  "PRD",
			Location:     "Germany",
			InstanceName: "host:27017",
			DBName:       "ercole",
			Charset:      "UTF8",
			Version:      "6.0.1",
		},
		{
			Hostname:     "test-db",
			Environment:  "PRD",
			Location:     "Germany",
			InstanceName: "host:27017",
			DBName:       "ercole",
			Charset:      "UTF8",
			Version:      "6.0.1",
		},
	}

	var expectedRes = dto.MongoDBInstanceResponse{
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

	db.EXPECT().SearchMongoDBInstances(
		[]string{"foo", "bar", "foobarx"}, "Hostname",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(&expectedRes, nil).Times(1)

	res, err := as.SearchMongoDBInstances(
		dto.SearchMongoDBInstancesFilter{
			GlobalFilter: dto.GlobalFilter{
				Location: "Italy", Environment: "PROD", OlderThan: utils.P("2019-12-05T14:02:03Z"),
			},
			Search: "foo bar foobarx", SortBy: "Hostname",
			SortDesc: true, PageNumber: 1, PageSize: 1,
		},
	)

	require.NoError(t, err)
	assert.Equal(t, &expectedRes, res)
}

func TestSearchMongoDBInstances_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchMongoDBInstances(
		[]string{"foo", "bar", "foobarx"}, "Memory",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchMongoDBInstances(

		dto.SearchMongoDBInstancesFilter{
			GlobalFilter: dto.GlobalFilter{
				Location: "Italy", Environment: "PROD", OlderThan: utils.P("2019-12-05T14:02:03Z"),
			},
			Search: "foo bar foobarx", SortBy: "Memory",
			SortDesc: true, PageNumber: 1, PageSize: 1,
		},
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}
