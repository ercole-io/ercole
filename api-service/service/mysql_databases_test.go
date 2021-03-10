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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchMySQLInstances(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		filter := dto.GlobalFilter{
			Location:    "Greece",
			Environment: "TEST",
			OlderThan:   utils.P("2020-05-20T09:53:34+00:00"),
		}

		expected := []dto.MySQLInstance{
			{
				Hostname:      "pippo",
				Location:      "",
				Environment:   "",
				MySQLInstance: model.MySQLInstance{},
			},
		}

		db.EXPECT().SearchMySQLInstances(filter).
			Return(expected, nil).Times(1)

		actual, err := as.SearchMySQLInstances(filter)
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		filter := dto.GlobalFilter{
			Location:    "Greece",
			Environment: "TEST",
			OlderThan:   utils.P("2020-05-20T09:53:34+00:00"),
		}

		db.EXPECT().SearchMySQLInstances(filter).
			Return(nil, errMock).Times(1)

		actual, err := as.SearchMySQLInstances(filter)
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}
