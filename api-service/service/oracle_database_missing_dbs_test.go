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
	"errors"
	"testing"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetMissingDatabases(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := []dto.OracleDatabaseMissingDbs{
		{
			Hostname: "host1",
			MissingDatabases: []model.MissingDatabase{
				{
					Name: "db1",
					Ignorable: model.Ignorable{
						Ignored:        false,
						IgnoredComment: "",
					},
				},
			},
		},
	}

	db.EXPECT().GetMissingDatabases().Return(expected, nil)

	res, err := as.GetMissingDatabases()

	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestGetMissingDatabases_Error(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	db.EXPECT().GetMissingDatabases().Return(nil, errors.New("cannot connect to database"))

	res, err := as.GetMissingDatabases()

	require.Nil(t, res)
	assert.Equal(t, errors.New("cannot connect to database"), err)
}

func TestGetMissingDatabasesByHostname(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := []model.MissingDatabase{
		{
			Name: "db1",
			Ignorable: model.Ignorable{
				Ignored:        false,
				IgnoredComment: "",
			},
		},
	}

	hostname := "host1"

	db.EXPECT().GetMissingDatabasesByHostname(hostname).Return(expected, nil)

	res, err := as.GetMissingDatabasesByHostname(hostname)

	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestGetMissingDatabasesByHostname_Error(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	hostname := "host1"

	db.EXPECT().GetMissingDatabasesByHostname(hostname).Return(nil, errors.New("cannot connect to database"))

	res, err := as.GetMissingDatabasesByHostname(hostname)

	require.Nil(t, res)
	assert.Equal(t, errors.New("cannot connect to database"), err)
}

func TestUpdateMissingDatabaseIgnoredField(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	hostname := "host1"
	dbname := "db1"
	ignored := true
	comment := "no longer needed"

	db.EXPECT().UpdateMissingDatabaseIgnoredField(hostname, dbname, ignored, comment).Return(nil)

	err := as.UpdateMissingDatabaseIgnoredField(hostname, dbname, ignored, comment)

	require.NoError(t, err)
}

func TestUpdateMissingDatabaseIgnoredField_Error(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	hostname := "host1"
	dbname := "db1"
	ignored := true
	comment := "no longer needed"

	db.EXPECT().UpdateMissingDatabaseIgnoredField(hostname, dbname, ignored, comment).Return(errors.New("cannot connect to database"))

	err := as.UpdateMissingDatabaseIgnoredField(hostname, dbname, ignored, comment)

	assert.Equal(t, errors.New("cannot connect to database"), err)
}
