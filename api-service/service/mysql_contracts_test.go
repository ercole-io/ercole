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

	"github.com/ercole-io/ercole/v2/config"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestAddMySQLContract(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database:    db,
		NewObjectID: utils.NewObjectIDForTests(),
	}

	t.Run("Success", func(t *testing.T) {
		expected := model.MySQLContract{
			ID:               utils.Str2oid("000000000000000000000001"),
			Type:             "server",
			NumberOfLicenses: 42,
			LicenseTypeID:    model.MySqlPartNumber,
			Clusters:         []string{"pippo"},
			Hosts:            []string{"pluto"},
		}
		db.EXPECT().AddMySQLContract(expected).
			Return(nil).Times(1)

		contract := model.MySQLContract{
			Type:             "server",
			NumberOfLicenses: 42,
			LicenseTypeID:    model.MySqlPartNumber,
			Clusters:         []string{"pippo"},
			Hosts:            []string{"pluto"},
		}
		actual, err := as.AddMySQLContract(contract)
		require.NoError(t, err)

		assert.Equal(t, &expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		contract := model.MySQLContract{
			ID:            utils.Str2oid("000000000000000000000002"),
			LicenseTypeID: model.MySqlPartNumber,
		}
		db.EXPECT().AddMySQLContract(contract).
			Return(errMock).Times(1)

		actual, err := as.AddMySQLContract(contract)
		assert.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestUpdateMySQLContract(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		contract := model.MySQLContract{
			LicenseTypeID: model.MySqlPartNumber,
		}
		db.EXPECT().UpdateMySQLContract(contract).
			Return(nil).Times(1)

		actual, err := as.UpdateMySQLContract(contract)
		require.NoError(t, err)
		assert.Equal(t, contract, *actual)
	})

	t.Run("Error", func(t *testing.T) {
		contract := model.MySQLContract{
			LicenseTypeID: model.MySqlPartNumber,
		}
		db.EXPECT().UpdateMySQLContract(contract).
			Return(errMock).Times(1)

		actual, err := as.UpdateMySQLContract(contract)
		require.EqualError(t, err, "MockError")
		assert.Nil(t, actual)
	})
}

func TestGetMySQLContracts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		expected := []model.MySQLContract{
			{
				ID:               [12]byte{},
				Type:             "",
				NumberOfLicenses: 0,
				Clusters:         []string{},
				Hosts:            []string{},
			},
		}
		db.EXPECT().GetMySQLContracts().
			Return(expected, nil).Times(1)

		actual, err := as.GetMySQLContracts()
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("Error", func(t *testing.T) {
		db.EXPECT().GetMySQLContracts().
			Return(nil, errMock).Times(1)

		actual, err := as.GetMySQLContracts()
		require.EqualError(t, err, "MockError")

		assert.Nil(t, actual)
	})
}

func TestGetMySQLContractsAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	data := []model.MySQLContract{
		{
			Type:             "server",
			ContractID:       "",
			CSI:              "",
			NumberOfLicenses: 42,
			Clusters:         []string{"pippo"},
			Hosts:            []string{"pluto"},
		},
	}

	db.EXPECT().GetMySQLContracts().
		Return(data, nil).Times(1)

	actual, err := as.GetMySQLContractsAsXLSX()
	require.NoError(t, err)
	assert.Equal(t, "server", actual.GetCellValue("Contracts", "A2"))
	assert.Equal(t, "", actual.GetCellValue("Contracts", "B2"))
	assert.Equal(t, "", actual.GetCellValue("Contracts", "C2"))
	assert.Equal(t, "42", actual.GetCellValue("Contracts", "D2"))
	assert.Equal(t, "[pippo]", actual.GetCellValue("Contracts", "E2"))
	assert.Equal(t, "pluto", actual.GetCellValue("Contracts", "F3"))
}

func TestDeleteMySQLContract(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	t.Run("Success", func(t *testing.T) {
		id := utils.Str2oid("iiiiiiiiiiiiiiiiiiiiiiii")
		db.EXPECT().DeleteMySQLContract(id).
			Return(nil).Times(1)

		err := as.DeleteMySQLContract(id)
		require.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		id := utils.Str2oid("iiiiiiiiiiiiiiiiiiiiiiii")
		db.EXPECT().DeleteMySQLContract(id).
			Return(errMock).Times(1)

		err := as.DeleteMySQLContract(id)
		require.EqualError(t, err, "MockError")
	})
}
