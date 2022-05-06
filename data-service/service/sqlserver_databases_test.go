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
	"github.com/stretchr/testify/assert"
)

var hostDataSql1 model.HostDataBE = model.HostDataBE{
	ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	Hostname:  "superhost1",
	Archived:  true,
	CreatedAt: utils.P("2019-11-05T18:02:03Z"),
	Features: model.Features{
		Microsoft: &model.MicrosoftFeature{
			SQLServer: &model.MicrosoftSQLServerFeature{
				Instances: []model.MicrosoftSQLServerInstance{
					{
						Name: "acd",
						License: model.MicrosoftSQLServerLicense{
							LicenseTypeID: "SqlServer ENT",
							Name:          "SqlServer ENT",
							Count:         10,
							Ignored:       false,
						},
					},
				},
			},
		},
	},
}

var hostDataSql2 model.HostDataBE = model.HostDataBE{
	ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	Hostname:  "superhost1",
	Archived:  true,
	CreatedAt: utils.P("2019-11-05T18:02:03Z"),
	Features: model.Features{
		Microsoft: &model.MicrosoftFeature{
			SQLServer: &model.MicrosoftSQLServerFeature{
				Instances: []model.MicrosoftSQLServerInstance{
					{
						Name: "acd",
						License: model.MicrosoftSQLServerLicense{
							LicenseTypeID: "SqlServer ENT",
							Name:          "SqlServer ENT",
							Count:         10,
							Ignored:       true,
						},
					},
				},
			},
		},
	},
}

var hostDataSql3 model.HostDataBE = model.HostDataBE{
	ID:        utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	Hostname:  "superhost1",
	Archived:  true,
	CreatedAt: utils.P("2019-11-05T18:02:03Z"),
	Features: model.Features{
		Microsoft: &model.MicrosoftFeature{
			SQLServer: &model.MicrosoftSQLServerFeature{
				Instances: []model.MicrosoftSQLServerInstance{
					{
						Name:    "acd",
						Version: "15.0.2080.9",
					},
				},
			},
		},
	},
}

func TestSqlServerIgnorePreviousLicences_SuccessNoPreviousIgnored(t *testing.T) {
	hds := HostDataService{
		Log: logger.NewLogger("TEST"),
	}

	hds.ignorePreviousLicences(&hostData6, &hostData6)

	for _, db := range hostDataSql1.Features.Microsoft.SQLServer.Instances {
		if db.License.LicenseTypeID == "SqlServer ENT" {
			if db.License.Ignored {
				t.Fatal("unexpected ignored license")
			}
		}
	}
}

func TestSqlServerIgnorePreviousLicences_SuccessWithPreviousIgnored(t *testing.T) {
	hds := HostDataService{
		Log: logger.NewLogger("TEST"),
	}

	hds.ignorePreviousLicences(&hostDataSql2, &hostDataSql1)

	for _, db := range hostDataSql1.Features.Microsoft.SQLServer.Instances {
		if db.License.LicenseTypeID == "SqlServer ENT" {
			if db.License.Ignored {
				t.Fatal("unexpected ignored license")
			}
		}
	}
}

func TestSqlServerVersion_Success(t *testing.T) {
	hds := HostDataService{
		Log: logger.NewLogger("TEST"),
	}

	host := hds.setSqlServerVersion(&hostDataSql3)

	assert.Contains(t, host.Features.Microsoft.SQLServer.Instances[0].Version, "2019")
}
