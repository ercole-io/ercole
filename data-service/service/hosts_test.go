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

// Package service is a package that provides methods for manipulating host informations

package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
				EnablePatching:       true,
				LogInsertingHostdata: false,
			},
		},
		ServerVersion:  "1.6.6",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:            utils.NewLogger("TEST"),
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	t.Run("New host", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().FindPatchingFunction("rac1_x").Return(model.PatchingFunction{}, nil),
			db.EXPECT().FindMostRecentHostDataOlderThan(hd.Hostname, utils.P("2019-11-05T14:02:03Z")).Return(nil, nil),
			asc.EXPECT().ThrowNewAlert(gomock.Any()).Do(func(a model.Alert) {
				assert.Equal(t, "The host rac1_x was added to ercole", a.Description)
			}).Return(nil),
			db.EXPECT().ArchiveHost("rac1_x").Return(nil, nil),
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
				Return(&mongo.InsertOneResult{InsertedID: utils.Str2oid("5dd3a8db184dbf295f0376f2")}, nil),
			db.EXPECT().DeleteNoDataAlertByHost(hd.Hostname).Return(nil),
		)

		err := hds.InsertHostData(hd)
		require.NoError(t, err)
	})

	t.Run("Update dismissed host", func(t *testing.T) {
		previousHostdata := &model.HostDataBE{Archived: true} // it's dismissed!

		gomock.InOrder(
			db.EXPECT().FindPatchingFunction("rac1_x").Return(model.PatchingFunction{}, nil),
			db.EXPECT().FindMostRecentHostDataOlderThan(hd.Hostname, utils.P("2019-11-05T14:02:03Z")).
				Return(previousHostdata, nil),
			asc.EXPECT().ThrowNewAlert(gomock.Any()).Do(func(a model.Alert) {
				assert.Equal(t, "The host rac1_x was added to ercole", a.Description)
			}).Return(nil),
			db.EXPECT().ArchiveHost("rac1_x").Return(nil, nil),
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
				Return(&mongo.InsertOneResult{InsertedID: utils.Str2oid("5dd3a8db184dbf295f0376f2")}, nil),
			db.EXPECT().DeleteNoDataAlertByHost(hd.Hostname).Return(nil),
		)

		err := hds.InsertHostData(hd)
		require.NoError(t, err)
	})
	t.Run("Update host", func(t *testing.T) {
		previousHostdata := &model.HostDataBE{Archived: false}

		gomock.InOrder(
			db.EXPECT().FindPatchingFunction("rac1_x").Return(model.PatchingFunction{}, nil),
			db.EXPECT().FindMostRecentHostDataOlderThan(hd.Hostname, utils.P("2019-11-05T14:02:03Z")).
				Return(previousHostdata, nil),
			db.EXPECT().ArchiveHost("rac1_x").Return(nil, nil),
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
				Return(&mongo.InsertOneResult{InsertedID: utils.Str2oid("5dd3a8db184dbf295f0376f2")}, nil),
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
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	gomock.InOrder(
		db.EXPECT().FindMostRecentHostDataOlderThan(hd.Hostname, utils.P("2019-11-05T14:02:03Z")).Return(nil, nil),
		asc.EXPECT().ThrowNewAlert(gomock.Any()).Do(func(a model.Alert) {
			assert.Equal(t, "The host rac1_x was added to ercole", a.Description)
		}).Return(nil),
		db.EXPECT().ArchiveHost("rac1_x").Return(nil, aerrMock),
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
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	gomock.InOrder(
		db.EXPECT().FindMostRecentHostDataOlderThan(hd.Hostname, utils.P("2019-11-05T14:02:03Z")).Return(nil, nil),
		asc.EXPECT().ThrowNewAlert(gomock.Any()).Do(func(a model.Alert) {
			assert.Equal(t, "The host rac1_x was added to ercole", a.Description)
		}).Return(nil),
		db.EXPECT().ArchiveHost("rac1_x").Return(nil, nil),
		db.EXPECT().InsertHostData(gomock.Any()).Return(nil, aerrMock).Do(func(newHD model.HostDataBE) {
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

func TestInsertHostData_DatabaseError3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			AlertService: config.AlertService{
				PublisherUsername: "publ1sh3r",
				PublisherPassword: "M0stS3cretP4ssw0rd",
				RemoteEndpoint:    "http://ercole.example.org",
			},
			DataService: config.DataService{
				EnablePatching:       true,
				LogInsertingHostdata: true,
			},
		},
		ServerVersion: "1.6.6",
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	db.EXPECT().FindPatchingFunction("rac1_x").Return(model.PatchingFunction{}, aerrMock)

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
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	gomock.InOrder(
		db.EXPECT().FindMostRecentHostDataOlderThan(hd.Hostname, utils.P("2019-11-05T14:02:03Z")).Return(nil, nil),
		asc.EXPECT().ThrowNewAlert(gomock.Any()).Do(func(a model.Alert) {
			assert.Equal(t, "The host rac1_x was added to ercole", a.Description)
		}).Return(nil),
		db.EXPECT().ArchiveHost("rac1_x").Return(nil, nil),
		db.EXPECT().InsertHostData(gomock.Any()).Return(nil, aerrMock).Do(func(newHD model.HostDataBE) {
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

func TestPatchHostData_SuccessNoPatchingFunction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	db.EXPECT().FindPatchingFunction("rac1_x").Return(model.PatchingFunction{}, nil)

	res, err := hds.patchHostData(hd)
	require.NoError(t, err)
	assert.Equal(t, hd, res)
}

func TestPatchHostData_SuccessPatchingFunction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			DataService: config.DataService{
				LogDataPatching: true,
			},
		},
		Log: utils.NewLogger("TEST"),
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")
	patchedHd := hd
	patchedHd.Tags = []string{"topolino", "pluto"}

	objID := utils.Str2oid("5ef9b4bcda4e04c0c1a94e9e")
	db.EXPECT().FindPatchingFunction("rac1_x").Return(model.PatchingFunction{
		ID:        &objID,
		CreatedAt: utils.P("2020-06-29T09:30:55+00:00"),
		Hostname:  "rac1_x",
		Vars: map[string]interface{}{
			"tags": []string{"topolino", "pluto"},
		},
		Code: `
			hostdata.tags = vars.tags;
		`,
	}, nil)

	res, err := hds.patchHostData(hd)
	require.NoError(t, err)
	assert.Equal(t, patchedHd, res)
}

func TestPatchHostData_FailPatchingFunction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	db.EXPECT().FindPatchingFunction("rac1_x").Return(model.PatchingFunction{}, aerrMock)

	_, err := hds.patchHostData(hd)
	require.Equal(t, aerrMock, err)
}

func TestPatchHostData_FailPatchingFunction2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			DataService: config.DataService{
				LogDataPatching: true,
			},
		},
		Log: utils.NewLogger("TEST"),
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_00.json")

	objID := utils.Str2oid("5ef9b4bcda4e04c0c1a94e9e")
	db.EXPECT().FindPatchingFunction("rac1_x").Return(model.PatchingFunction{
		ID:        &objID,
		CreatedAt: utils.P("2020-06-29T09:30:55+00:00"),
		Hostname:  "rac1_x",
		Vars: map[string]interface{}{
			"tags": []string{"topolino", "pluto"},
		},
		Code: `
			sdfsdasdfsdf
		`,
	}, nil)

	_, err := hds.patchHostData(hd)
	assert.Error(t, err)
}
