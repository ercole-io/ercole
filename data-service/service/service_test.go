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

// Package service is a package that provides methods for manipulating host informations

package service

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHostInfo_Success(t *testing.T) {
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
		},
		Version: "1.6.6",
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_00.json")

	db.EXPECT().ArchiveHost("rac1_x").Return(nil, nil).Times(1)
	db.EXPECT().ArchiveHost(gomock.Any()).Times(0)
	db.EXPECT().FindPatchingFunction(gomock.Any()).Return(model.PatchingFunction{}, nil).Times(0)
	db.EXPECT().InsertHostData(gomock.Any()).Return(&mongo.InsertOneResult{InsertedID: utils.Str2oid("5dd3a8db184dbf295f0376f2")}, nil).Do(func(newHD map[string]interface{}) {
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD["_id"].(primitive.ObjectID).Timestamp())
		assert.False(t, newHD["Archived"].(bool))
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD["CreatedAt"])
		assert.Equal(t, model.SchemaVersion, newHD["SchemaVersion"])
		assert.Equal(t, "1.6.6", newHD["ServerVersion"])
		assert.Equal(t, hd.Hostname, newHD["Hostname"])
		assert.Equal(t, hd.Environment, newHD["Environment"])
		//I assume that other fields are correct
	}).Times(1)
	db.EXPECT().InsertHostData(gomock.Any()).Times(0)
	http.DefaultClient = NewHTTPTestClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "http://publ1sh3r:M0stS3cretP4ssw0rd@ercole.example.org/queue/host-data-insertion/5dd3a8db184dbf295f0376f2", req.URL.String())

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
			Header:     make(http.Header),
		}, nil
	})

	var out map[string]interface{}
	mapstructure.Decode(hd, &out)

	res, err := hds.UpdateHostInfo(out)
	require.NoError(t, err)
	assert.Equal(t, utils.Str2oid("5dd3a8db184dbf295f0376f2"), res)
}

func TestUpdateHostInfo_DatabaseError1(t *testing.T) {
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
		},
		Version: "1.6.6",
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_00.json")
	db.EXPECT().ArchiveHost("rac1_x").Return(nil, aerrMock).Times(1)
	db.EXPECT().ArchiveHost(gomock.Any()).Times(0)
	db.EXPECT().FindPatchingFunction(gomock.Any()).Return(model.PatchingFunction{}, nil).Times(0)

	var out map[string]interface{}
	mapstructure.Decode(hd, &out)

	_, err := hds.UpdateHostInfo(out)
	require.Equal(t, aerrMock, err)
}

func TestUpdateHostInfo_DatabaseError2(t *testing.T) {
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
		},
		Version: "1.6.6",
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_00.json")

	db.EXPECT().ArchiveHost("rac1_x").Return(nil, nil).Times(1)
	db.EXPECT().ArchiveHost(gomock.Any()).Times(0)
	db.EXPECT().FindPatchingFunction(gomock.Any()).Return(model.PatchingFunction{}, nil).Times(0)
	db.EXPECT().InsertHostData(gomock.Any()).Return(nil, aerrMock).Do(func(newHD map[string]interface{}) {
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD["_id"].(primitive.ObjectID).Timestamp())
		assert.False(t, newHD["Archived"].(bool))
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD["CreatedAt"])
		assert.Equal(t, model.SchemaVersion, newHD["SchemaVersion"])
		assert.Equal(t, "1.6.6", newHD["ServerVersion"])
		assert.Equal(t, hd.Hostname, newHD["Hostname"])
		assert.Equal(t, hd.Environment, newHD["Environment"])
		//I assume that other fields are correct
	}).Times(1)
	db.EXPECT().InsertHostData(gomock.Any()).Times(0)

	var out map[string]interface{}
	mapstructure.Decode(hd, &out)

	_, err := hds.UpdateHostInfo(out)
	require.Equal(t, aerrMock, err)
}

func TestUpdateHostInfo_HttpError(t *testing.T) {
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
		},
		Version: "1.6.6",
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_00.json")

	db.EXPECT().ArchiveHost("rac1_x").Return(nil, nil).Times(1)
	db.EXPECT().ArchiveHost(gomock.Any()).Times(0)
	db.EXPECT().FindPatchingFunction(gomock.Any()).Return(model.PatchingFunction{}, nil).Times(0)
	db.EXPECT().InsertHostData(gomock.Any()).Return(&mongo.InsertOneResult{InsertedID: utils.Str2oid("5dd3a8db184dbf295f0376f2")}, nil).Do(func(newHD map[string]interface{}) {
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD["_id"].(primitive.ObjectID).Timestamp())
		assert.False(t, newHD["Archived"].(bool))
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD["CreatedAt"])
		assert.Equal(t, model.SchemaVersion, newHD["SchemaVersion"])
		assert.Equal(t, "1.6.6", newHD["ServerVersion"])
		assert.Equal(t, hd.Hostname, newHD["Hostname"])
		assert.Equal(t, hd.Environment, newHD["Environment"])
		//I assume that other fields are correct
	}).Times(1)
	db.EXPECT().InsertHostData(gomock.Any()).Times(0)
	http.DefaultClient = NewHTTPTestClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "http://publ1sh3r:M0stS3cretP4ssw0rd@ercole.example.org/queue/host-data-insertion/5dd3a8db184dbf295f0376f2", req.URL.String())
		return nil, errMock
	})

	var out map[string]interface{}
	mapstructure.Decode(hd, &out)

	_, err := hds.UpdateHostInfo(out)
	require.Equal(t, "EVENT ENQUEUE", err.ErrorClass())
	require.Contains(t, err.Error(), "http://publ1sh3r:***@ercole.example.org/queue/host-data-insertion/5dd3a8db184dbf295f0376f2")
	require.Contains(t, err.Error(), "MockError")
}

func TestUpdateHostInfo_HttpError2(t *testing.T) {
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
		},
		Version: "1.6.6",
	}
	hd := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_00.json")

	db.EXPECT().ArchiveHost("rac1_x").Return(nil, nil).Times(1)
	db.EXPECT().ArchiveHost(gomock.Any()).Times(0)
	db.EXPECT().FindPatchingFunction(gomock.Any()).Return(model.PatchingFunction{}, nil).Times(0)
	db.EXPECT().InsertHostData(gomock.Any()).Return(&mongo.InsertOneResult{InsertedID: utils.Str2oid("5dd3a8db184dbf295f0376f2")}, nil).Do(func(newHD map[string]interface{}) {
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD["_id"].(primitive.ObjectID).Timestamp())
		assert.False(t, newHD["Archived"].(bool))
		assert.Equal(t, utils.P("2019-11-05T14:02:03Z"), newHD["CreatedAt"])
		assert.Equal(t, model.SchemaVersion, newHD["SchemaVersion"])
		assert.Equal(t, "1.6.6", newHD["ServerVersion"])
		assert.Equal(t, hd.Hostname, newHD["Hostname"])
		assert.Equal(t, hd.Environment, newHD["Environment"])
		//I assume that other fields are correct
	}).Times(1)
	db.EXPECT().InsertHostData(gomock.Any()).Times(0)
	http.DefaultClient = NewHTTPTestClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "http://publ1sh3r:M0stS3cretP4ssw0rd@ercole.example.org/queue/host-data-insertion/5dd3a8db184dbf295f0376f2", req.URL.String())
		return &http.Response{
			StatusCode: 500,
			Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
			Header:     make(http.Header),
		}, nil
	})

	var out map[string]interface{}
	mapstructure.Decode(hd, &out)

	_, err := hds.UpdateHostInfo(out)
	require.Equal(t, "EVENT ENQUEUE", err.ErrorClass())
	require.EqualError(t, err, "Failed to enqueue event")
}
