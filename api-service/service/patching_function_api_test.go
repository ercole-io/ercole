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

	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchLicenseModifiers_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []map[string]interface{}{
		{
			"DatabaseName": "foobar1",
			"Hostname":     "test-db1",
			"LicenseName":  "Oracle PIPPO",
			"NewValue":     20,
			"_id":          utils.Str2oid("5ece1c740fa23c31d597d8b1"),
		},
		{
			"DatabaseName": "foobar1",
			"Hostname":     "test-db2",
			"LicenseName":  "Oracle EXE",
			"NewValue":     10,
			"_id":          utils.Str2oid("5ece1c2d0fa23c31d597d8b0"),
		},
	}

	db.EXPECT().SearchLicenseModifiers(
		[]string{"foo", "bar", "foobarx"}, "Hostname",
		true, 1, 1,
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchLicenseModifiers(
		"foo bar foobarx", "Hostname",
		true, 1, 1,
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchLicenseModifiers_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchLicenseModifiers(
		[]string{"foo", "bar", "foobarx"}, "Hostname",
		true, 1, 1,
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchLicenseModifiers(
		"foo bar foobarx", "Hostname",
		true, 1, 1,
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}
