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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package service

import (
	"testing"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/data-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCompareCmdbInfo_DbError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)

	hds := HostDataService{
		Config:         config.Configuration{},
		ServerVersion:  "1.6.6",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:            utils.NewLogger("TEST"),
	}

	db.EXPECT().GetCurrentHostnames().
		Return(nil, aerrMock)

	actualErr := hds.CompareCmdbInfo(dto.CmdbInfo{})
	assert.Equal(t, aerrMock, actualErr)
}

func TestCompareCmdbInfo_NoAlerts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)

	hds := HostDataService{
		Config:         config.Configuration{},
		ServerVersion:  "1.6.6",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:            utils.NewLogger("TEST"),
	}

	db.EXPECT().GetCurrentHostnames().
		Return([]string{"pippo", "topolino", "pluto"}, nil)

	cmdbInfo := dto.CmdbInfo{
		Name:      "thisCmdb",
		Hostnames: []string{"pippo", "topolino", "pluto"},
	}
	actualErr := hds.CompareCmdbInfo(cmdbInfo)
	assert.Nil(t, actualErr)
}

func TestCompareCmdbInfo_MissingHostInErcole(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)

	hds := HostDataService{
		Config:         config.Configuration{},
		ServerVersion:  "1.6.6",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:            utils.NewLogger("TEST"),
	}

	db.EXPECT().GetCurrentHostnames().
		Return([]string{"PIPPO.TOPOLINIA.TOP", "topolino"}, nil)

	alert := model.Alert{
		AlertCategory: model.AlertCategoryEngine,
		AlertCode:     model.AlertCodeMissingHostInErcole,
		AlertSeverity: model.AlertSeverityWarning,
		AlertStatus:   model.AlertStatusNew,
		Description:   "Received unknown hostname pluto from CMDB thisCmdb",
		Date:          hds.TimeNow(),
	}
	asc.EXPECT().ThrowNewAlert(alert).Return(nil)

	cmdbInfo := dto.CmdbInfo{
		Name:      "thisCmdb",
		Hostnames: []string{"pippo", "topolino.topolinia.top", "pluto"},
	}
	actualErr := hds.CompareCmdbInfo(cmdbInfo)
	assert.Nil(t, actualErr)
}

func TestCompareCmdbInfo_MissingHostInCmdb(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)

	hds := HostDataService{
		Config:         config.Configuration{},
		ServerVersion:  "1.6.6",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:            utils.NewLogger("TEST"),
	}

	db.EXPECT().GetCurrentHostnames().
		Return([]string{"pippo.topolinia.top", "TOPOLINO", "pluto"}, nil)

	alert := model.Alert{
		AlertCategory: model.AlertCategoryEngine,
		AlertCode:     model.AlertCodeMissingHostInCmdb,
		AlertSeverity: model.AlertSeverityWarning,
		AlertStatus:   model.AlertStatusNew,
		Description:   "Missing hostname pluto in CMDB thisCmdb",
		Date:          hds.TimeNow(),
	}
	asc.EXPECT().ThrowNewAlert(alert).Return(nil)

	cmdbInfo := dto.CmdbInfo{
		Name:      "thisCmdb",
		Hostnames: []string{"pippo", "topolino.topolinia.top"},
	}
	actualErr := hds.CompareCmdbInfo(cmdbInfo)
	assert.Nil(t, actualErr)
}
