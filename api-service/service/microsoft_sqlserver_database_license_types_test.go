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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
)

func TestGetSqlServerDatabaseLicenseTypes_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := []model.SqlServerDatabaseLicenseType{
		{
			ID:              "Pippo",
			ItemDescription: "Pippo",
			Edition:         "Pippo",
			Version:         "Pippo",
		},
	}
	db.EXPECT().GetSqlServerDatabaseLicenseTypes().Return(expected, nil)

	res, err := as.GetSqlServerDatabaseLicenseTypes()
	require.NoError(t, err)
	assert.Equal(t, []model.SqlServerDatabaseLicenseType{
		{
			ID:              "Pippo",
			ItemDescription: "Pippo",
			Edition:         "Pippo",
			Version:         "Pippo",
		},
	}, res)
}
