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

// Package service is a package that provides methods for manipulating host informations

package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
)

func TestInsertHostData_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config: config.Configuration{
			AlertService: config.AlertService{
				PublisherUsername: "publ1sh3r",
				PublisherPassword: "M0stS3cretP4ssw0rd",
				RemoteEndpoint:    "http://ercole.example.org"},
			DataService: config.DataService{
				LogInsertingHostdata: false,
			},
		},
		ServerVersion:  "1.6.6",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:            logger.NewLogger("TEST"),
	}
	hd := mongoutils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	t.Run("New host", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().FindMostRecentHostDataOlderThan(hd.Hostname, utils.P("2019-11-05T14:02:03Z")).Return(nil, nil),
			asc.EXPECT().ThrowNewAlert(gomock.Any()).Do(func(a model.Alert) {
				assert.Equal(t, "The host rac1_x was added to ercole", a.Description)
			}).Return(nil),
			db.EXPECT().DismissHost("rac1_x").Return(nil),
			db.EXPECT().InsertHostData(gomock.Any()).
				Do(func(newHD model.HostDataBE) {
					assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.ID.Timestamp())
					assert.False(t, newHD.Archived)
					assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.CreatedAt)
					assert.Equal(t, model.SchemaVersion, newHD.ServerSchemaVersion)
					assert.Equal(t, "1.6.6", newHD.ServerVersion)
					assert.Equal(t, hd.Hostname, newHD.Hostname)
					assert.Equal(t, hd.Environment, newHD.Environment)
					//I assume that other fields are correct
				}).
				Return(nil),
			db.EXPECT().DeleteNoDataAlertByHost(hd.Hostname).Return(nil),
		)

		err := hds.InsertHostData(hd)
		require.NoError(t, err)
	})

	t.Run("Update dismissed host", func(t *testing.T) {
		previousHostdata := &model.HostDataBE{Archived: true} // it's dismissed!

		gomock.InOrder(
			db.EXPECT().FindMostRecentHostDataOlderThan(hd.Hostname, utils.P("2019-11-05T14:02:03Z")).
				Return(previousHostdata, nil),
			asc.EXPECT().ThrowNewAlert(gomock.Any()).Do(func(a model.Alert) {
				assert.Equal(t, "The host rac1_x was added to ercole", a.Description)
			}).Return(nil),
			db.EXPECT().DismissHost("rac1_x").Return(nil),
			db.EXPECT().InsertHostData(gomock.Any()).
				Do(func(newHD model.HostDataBE) {
					assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.ID.Timestamp())
					assert.False(t, newHD.Archived)
					assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.CreatedAt)
					assert.Equal(t, model.SchemaVersion, newHD.ServerSchemaVersion)
					assert.Equal(t, "1.6.6", newHD.ServerVersion)
					assert.Equal(t, hd.Hostname, newHD.Hostname)
					assert.Equal(t, hd.Environment, newHD.Environment)
					//I assume that other fields are correct
				}).
				Return(nil),
			db.EXPECT().DeleteNoDataAlertByHost(hd.Hostname).Return(nil),
		)

		err := hds.InsertHostData(hd)
		require.NoError(t, err)
	})
	t.Run("Update host", func(t *testing.T) {
		previousHostdata := &model.HostDataBE{Archived: false}

		gomock.InOrder(
			db.EXPECT().FindMostRecentHostDataOlderThan(hd.Hostname, utils.P("2019-11-05T14:02:03Z")).
				Return(previousHostdata, nil),
			db.EXPECT().DismissHost("rac1_x").Return(nil),
			db.EXPECT().InsertHostData(gomock.Any()).
				Do(func(newHD model.HostDataBE) {
					assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.ID.Timestamp())
					assert.False(t, newHD.Archived)
					assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.CreatedAt)
					assert.Equal(t, model.SchemaVersion, newHD.ServerSchemaVersion)
					assert.Equal(t, "1.6.6", newHD.ServerVersion)
					assert.Equal(t, hd.Hostname, newHD.Hostname)
					assert.Equal(t, hd.Environment, newHD.Environment)
					//I assume that other fields are correct
				}).
				Return(nil),
			db.EXPECT().DeleteNoDataAlertByHost(hd.Hostname).Return(nil),
		)

		err := hds.InsertHostData(hd)
		require.NoError(t, err)
	})
}

func TestInsertHostData_DatabaseError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database:       db,
		AlertSvcClient: asc,
		Config: config.Configuration{
			AlertService: config.AlertService{
				PublisherUsername: "publ1sh3r",
				PublisherPassword: "M0stS3cretP4ssw0rd",
				RemoteEndpoint:    "http://ercole.example.org",
			},
		},
		ServerVersion: "1.6.6",
	}
	hd := mongoutils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	gomock.InOrder(
		db.EXPECT().FindMostRecentHostDataOlderThan(hd.Hostname, utils.P("2019-11-05T14:02:03Z")).Return(nil, nil),
		asc.EXPECT().ThrowNewAlert(gomock.Any()).Do(func(a model.Alert) {
			assert.Equal(t, "The host rac1_x was added to ercole", a.Description)
		}).Return(nil),
		db.EXPECT().DismissHost("rac1_x").Return(aerrMock),
	)

	err := hds.InsertHostData(hd)
	require.Equal(t, aerrMock, err)
}

func TestInsertHostData_DatabaseError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database:       db,
		AlertSvcClient: asc,
		Config: config.Configuration{
			AlertService: config.AlertService{
				PublisherUsername: "publ1sh3r",
				PublisherPassword: "M0stS3cretP4ssw0rd",
				RemoteEndpoint:    "http://ercole.example.org",
			},
		},
		ServerVersion: "1.6.6",
	}
	hd := mongoutils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	gomock.InOrder(
		db.EXPECT().FindMostRecentHostDataOlderThan(hd.Hostname, utils.P("2019-11-05T14:02:03Z")).Return(nil, nil),
		asc.EXPECT().ThrowNewAlert(gomock.Any()).Do(func(a model.Alert) {
			assert.Equal(t, "The host rac1_x was added to ercole", a.Description)
		}).Return(nil),
		db.EXPECT().DismissHost("rac1_x").Return(nil),
		db.EXPECT().InsertHostData(gomock.Any()).Return(aerrMock).Do(func(newHD model.HostDataBE) {
			assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.ID.Timestamp())
			assert.False(t, newHD.Archived)
			assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.CreatedAt)
			assert.Equal(t, model.SchemaVersion, newHD.ServerSchemaVersion)
			assert.Equal(t, "1.6.6", newHD.ServerVersion)
			assert.Equal(t, hd.Hostname, newHD.Hostname)
			assert.Equal(t, hd.Environment, newHD.Environment)
			//I assume that other fields are correct
		}),
	)

	err := hds.InsertHostData(hd)
	require.Equal(t, aerrMock, err)
}

func TestInsertHostData_HttpError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database:       db,
		AlertSvcClient: asc,
		Config: config.Configuration{
			AlertService: config.AlertService{
				PublisherUsername: "publ1sh3r",
				PublisherPassword: "M0stS3cretP4ssw0rd",
				RemoteEndpoint:    "http://ercole.example.org",
			},
		},
		ServerVersion: "1.6.6",
	}
	hd := mongoutils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	gomock.InOrder(
		db.EXPECT().FindMostRecentHostDataOlderThan(hd.Hostname, utils.P("2019-11-05T14:02:03Z")).Return(nil, nil),
		asc.EXPECT().ThrowNewAlert(gomock.Any()).Do(func(a model.Alert) {
			assert.Equal(t, "The host rac1_x was added to ercole", a.Description)
		}).Return(nil),
		db.EXPECT().DismissHost("rac1_x").Return(nil),
		db.EXPECT().InsertHostData(gomock.Any()).Return(aerrMock).Do(func(newHD model.HostDataBE) {
			assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.ID.Timestamp())
			assert.False(t, newHD.Archived)
			assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD.CreatedAt)
			assert.Equal(t, model.SchemaVersion, newHD.ServerSchemaVersion)
			assert.Equal(t, "1.6.6", newHD.ServerVersion)
			assert.Equal(t, hd.Hostname, newHD.Hostname)
			assert.Equal(t, hd.Environment, newHD.Environment)
			//I assume that other fields are correct
		}),
	)

	err := hds.InsertHostData(hd)
	fmt.Println(err.Error())
	require.Contains(t, err.Error(), "MockError")
}
