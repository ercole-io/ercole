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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestUpdateMySqlLicenseIgnoredField_Success(t *testing.T) {
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
					MySQL: &model.MySQLFeature{
						Instances: []model.MySQLInstance{
							{
								Name: "TEST123",
								License: model.MySQLLicense{
									Ignored: true,
								},
							},
						},
					},
				},
			},
		}

		hostname, instancename := "foobar", "TEST123"
		ignored := true

		db.EXPECT().UpdateMySqlLicenseIgnoredField(hostname, instancename, ignored, "").Return(nil)

		err := as.UpdateMySqlLicenseIgnoredField(hostname, instancename, ignored, "")
		require.NoError(t, err)

		var resultIgnored bool
		for i := range hostdata[0].Features.MySQL.Instances {
			db := &hostdata[0].Features.MySQL.Instances[i]
			if db.Name == instancename {
				lic := &db.License
				resultIgnored = lic.Ignored
			}
		}

		require.Equal(t, ignored, resultIgnored)

		db.EXPECT().UpdateMySqlLicenseIgnoredField(hostname, instancename, !ignored, "test comment").Return(nil)

		err = as.UpdateMySqlLicenseIgnoredField(hostname, instancename, !ignored, "test comment")
		require.NoError(t, err)
	})
}

func TestUpdateMySqlLicenseIgnoredField_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	hostname, instancename := "paperino", "TEST123"
	ignored := false

	errUpd := utils.ErrLicenseNotFound

	db.EXPECT().UpdateMySqlLicenseIgnoredField(hostname, instancename, ignored, "").Return(errUpd)
	err := as.UpdateMySqlLicenseIgnoredField(hostname, instancename, ignored, "")
	assert.EqualError(t, err, utils.ErrLicenseNotFound.Error())
}
