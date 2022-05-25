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

	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

var hostDataMySql1 model.HostDataBE = model.HostDataBE{
	ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	Hostname:  "superhost1",
	Archived:  true,
	CreatedAt: utils.P("2019-11-05T18:02:03Z"),
	Features: model.Features{
		MySQL: &model.MySQLFeature{
			Instances: []model.MySQLInstance{
				{
					Name: "acd",
					License: model.MySQLLicense{
						LicenseTypeID: "ENTERPRISE",
						Name:          "ENTERPRISE",
						Count:         10,
						Ignored:       false,
					},
				},
			},
		},
	},
}

var hostDataMySql2 model.HostDataBE = model.HostDataBE{
	ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	Hostname:  "superhost1",
	Archived:  true,
	CreatedAt: utils.P("2019-11-05T18:02:03Z"),
	Features: model.Features{
		MySQL: &model.MySQLFeature{
			Instances: []model.MySQLInstance{
				{
					Name: "acd",
					License: model.MySQLLicense{
						LicenseTypeID: "ENTERPRISE",
						Name:          "ENTERPRISE",
						Count:         10,
						Ignored:       false,
					},
				},
			},
		},
	},
}

func TestMySqlIgnorePreviousLicences_SuccessNoPreviousIgnored(t *testing.T) {
	hds := HostDataService{
		Log: logger.NewLogger("TEST"),
	}

	hds.ignoreMySqlPreviousLicences(&hostDataMySql1, &hostDataMySql2)

	for _, db := range hostDataMySql1.Features.MySQL.Instances {
		if db.License.LicenseTypeID == "ENTERPRISE" {
			if db.License.Ignored {
				t.Fatal("unexpected ignored license")
			}
		}
	}
}

func TestMySqlIgnorePreviousLicences_SuccessWithPreviousIgnored(t *testing.T) {
	hds := HostDataService{
		Log: logger.NewLogger("TEST"),
	}

	hds.ignoreMySqlPreviousLicences(&hostDataMySql1, &hostDataMySql2)

	for _, db := range hostDataMySql1.Features.MySQL.Instances {
		if db.License.LicenseTypeID == "ENTERPRISE" {
			if db.License.Ignored {
				t.Fatal("unexpected ignored license")
			}
		}
	}
}
