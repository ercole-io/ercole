// Copyright (c) 2021 Sorint.lab S.p.A.
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
	"github.com/ercole-io/ercole/v2/utils"
)

func TestUpdateLicenseIgnoredField_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)

	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		TimeNow:     utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		NewObjectID: utils.NewObjectIDForTests(),
	}

	t.Run("Success", func(t *testing.T) {
		hostdata := []model.HostDataBE{
			{
				Hostname: "foobar",
				Features: model.Features{
					Oracle: &model.OracleFeature{
						Database: &model.OracleDatabaseFeature{
							Databases: []model.OracleDatabase{
								{
									InstanceName: "TEST123",
									Licenses: []model.OracleDatabaseLicense{
										{
											LicenseTypeID: "A90611",
											Ignored:       true,
										},
									},
								},
							},
						},
					},
				},
			},
		}

		hostname, dbname, licenseTypeID := "foobar", "TEST123", "A90611"
		ignored := true

		db.EXPECT().UpdateLicenseIgnoredField(hostname, dbname, licenseTypeID, ignored).Return(nil)

		err := as.UpdateLicenseIgnoredField(hostname, dbname, licenseTypeID, ignored)
		require.NoError(t, err)

		var resultIgnored bool
		for i := range hostdata[0].Features.Oracle.Database.Databases {
			db := &hostdata[0].Features.Oracle.Database.Databases[i]
			if db.InstanceName == dbname {
				for j := range db.Licenses {
					lic := &db.Licenses[j]
					if lic.LicenseTypeID == licenseTypeID {
						resultIgnored = lic.Ignored
					}
				}
			}
		}

		require.Equal(t, ignored, resultIgnored)

		db.EXPECT().UpdateLicenseIgnoredField(hostname, dbname, licenseTypeID, !ignored).Return(nil)

		err = as.UpdateLicenseIgnoredField(hostname, dbname, licenseTypeID, !ignored)
		require.NoError(t, err)
	})
}

func TestUpdateLicenseIgnoredField_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	hostname, dbname, licenseTypeID := "paperino", "TEST123", "A90611"
	ignored := false

	errUpd := utils.ErrLicenseNotFound

	db.EXPECT().UpdateLicenseIgnoredField(hostname, dbname, licenseTypeID, ignored).Return(errUpd)
	err := as.UpdateLicenseIgnoredField(hostname, dbname, licenseTypeID, ignored)
	assert.EqualError(t, err, utils.ErrLicenseNotFound.Error())
}
