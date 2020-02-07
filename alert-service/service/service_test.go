// Copyright (c) 2019 Sorint.lab S.p.A.
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

	"github.com/stretchr/testify/require"

	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestProcessHostDataInsertion_SuccessNewHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().FindHostData(utils.Str2oid("5dc3f534db7e81a98b726a52")).Return(hostData1, nil).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan("superhost1", utils.P("2019-11-05T14:02:03Z")).Return(model.HostData{}, nil).Times(1)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, "The server 'superhost1' was added to ercole", alert.Description)
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), alert.Date)
	}).Times(1)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": utils.Str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestProcessHostDataInsertion_DatabaseError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().FindHostData(utils.Str2oid("5dc3f534db7e81a98b726a52")).Return(model.HostData{}, aerrMock).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": utils.Str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestProcessHostDataInsertion_DatabaseError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().FindHostData(utils.Str2oid("5dc3f534db7e81a98b726a52")).Return(hostData1, nil).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan("superhost1", utils.P("2019-11-05T14:02:03Z")).Return(model.HostData{}, aerrMock).Times(1)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": utils.Str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestProcessHostDataInsertion_DiffHostError3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().FindHostData(utils.Str2oid("5dc3f534db7e81a98b726a52")).Return(hostData1, nil).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan("superhost1", utils.P("2019-11-05T14:02:03Z")).Return(model.HostData{}, nil).Times(1)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock).Times(1)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": utils.Str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestDiffHostDataAndGenerateAlert_SuccessNoDifferences(t *testing.T) {
	as := AlertService{}

	require.NoError(t, as.DiffHostDataAndGenerateAlert(hostData2, hostData1))
}

func TestDiffHostDataAndGenerateAlert_SuccessNewHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T16:02:03Z")),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewServer,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.NoError(t, as.DiffHostDataAndGenerateAlert(model.HostData{}, hostData1))
}

func TestDiffHostDataAndGenerateAlert_SuccessNewDatabase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T16:02:03Z")),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewDatabase,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
			"Dbname":   "acd",
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.NoError(t, as.DiffHostDataAndGenerateAlert(hostData1, hostData3))
}

func TestDiffHostDataAndGenerateAlert_SuccessNewEnterpriseLicense(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T16:02:03Z")),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewLicense,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewOption,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
			"Dbname":   "acd",
			"Features": []string{"Driving"},
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.NoError(t, as.DiffHostDataAndGenerateAlert(hostData3, hostData4))
}

func TestDiffHostDataAndGenerateAlert_DatabaseError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T16:02:03Z")),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewServer,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
		},
	}}).Return(nil, aerrMock).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Equal(t, aerrMock, as.DiffHostDataAndGenerateAlert(model.HostData{}, hostData1))
}

func TestDiffHostDataAndGenerateAlert_DatabaseError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T16:02:03Z")),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewDatabase,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
			"Dbname":   "acd",
		},
	}}).Return(nil, aerrMock).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Equal(t, aerrMock, as.DiffHostDataAndGenerateAlert(hostData1, hostData3))
}

func TestDiffHostDataAndGenerateAlert_DatabaseError3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T16:02:03Z")),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewLicense,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
		},
	}}).Return(nil, aerrMock).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Equal(t, aerrMock, as.DiffHostDataAndGenerateAlert(hostData3, hostData4))
}

func TestDiffHostDataAndGenerateAlert_DatabaseError4(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  utils.Btc(utils.P("2019-11-05T16:02:03Z")),
	}

	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewLicense,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
		},
	}}).Return(nil, nil).Times(1)
	db.EXPECT().InsertAlert(&alertSimilarTo{al: model.Alert{
		AlertCode: model.AlertCodeNewOption,
		OtherInfo: map[string]interface{}{
			"Hostname": "superhost1",
			"Dbname":   "acd",
			"Features": []string{"Driving"},
		},
	}}).Return(nil, aerrMock).Times(1)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	require.Equal(t, aerrMock, as.DiffHostDataAndGenerateAlert(hostData3, hostData4))
}
