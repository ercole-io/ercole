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

package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/plandem/xlsx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchHosts_JSONPaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := map[string]interface{}{
		"Content": []interface{}{
			map[string]interface{}{
				"CPUCores":       1,
				"CPUModel":       "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
				"CPUThreads":     2,
				"Cluster":        "Angola-1dac9f7418db9b52c259ce4ba087cdb6",
				"CreatedAt":      utils.P("2020-04-07T08:52:59.844+02:00"),
				"Databases":      "8888888-d41d8cd98f00b204e9800998ecf8427e",
				"Environment":    "PROD",
				"HostType":       "virtualization",
				"Hostname":       "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
				"Kernel":         "3.10.0-862.9.1.el7.x86_64",
				"Location":       "Italy",
				"MemTotal":       3,
				"OS":             "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
				"OracleCluster":  false,
				"PhysicalHost":   "suspended-290dce22a939f3868f8f23a6e1f57dd8",
				"Socket":         2,
				"SunCluster":     false,
				"SwapTotal":      4,
				"Type":           "VMWARE",
				"VeritasCluster": false,
				"Version":        "1.6.1",
				"Virtual":        true,
				"_id":            utils.Str2oid("5e8c234b24f648a08585bd3d"),
			},
			map[string]interface{}{
				"CPUCores":       1,
				"CPUModel":       "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
				"CPUThreads":     2,
				"Cluster":        "Puzzait",
				"CreatedAt":      utils.P("2020-04-07T08:52:59.869+02:00"),
				"Databases":      "",
				"Environment":    "PROD",
				"HostType":       "virtualization",
				"Hostname":       "test-virt",
				"Kernel":         "3.10.0-862.9.1.el7.x86_64",
				"Location":       "Italy",
				"MemTotal":       3,
				"OS":             "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
				"OracleCluster":  false,
				"PhysicalHost":   "s157-cb32c10a56c256746c337e21b3f82402",
				"Socket":         2,
				"SunCluster":     false,
				"SwapTotal":      4,
				"Type":           "VMWARE",
				"VeritasCluster": false,
				"Version":        "1.6.1",
				"Virtual":        true,
				"_id":            utils.Str2oid("5e8c234b24f648a08585bd41"),
			},
		},
		"Metadata": map[string]interface{}{
			"Empty":         false,
			"First":         true,
			"Last":          true,
			"Number":        0,
			"Size":          20,
			"TotalElements": 25,
			"TotalPages":    1,
		},
	}

	resFromService := []map[string]interface{}{
		expectedRes,
	}

	as.EXPECT().
		SearchHosts("summary", "foobar", "Hostname", true, 2, 3, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?mode=summary&search=foobar&sort-by=Hostname&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchHosts_JSONUnpaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		{
			"CPUCores":       1,
			"CPUModel":       "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
			"CPUThreads":     2,
			"Cluster":        "Angola-1dac9f7418db9b52c259ce4ba087cdb6",
			"CreatedAt":      utils.P("2020-04-07T08:52:59.844+02:00"),
			"Databases":      "8888888-d41d8cd98f00b204e9800998ecf8427e",
			"Environment":    "PROD",
			"HostType":       "virtualization",
			"Hostname":       "fb-canvas-b9b1d8fa8328fe972b1e031621e8a6c9",
			"Kernel":         "3.10.0-862.9.1.el7.x86_64",
			"Location":       "Italy",
			"MemTotal":       3,
			"OS":             "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
			"OracleCluster":  false,
			"PhysicalHost":   "suspended-290dce22a939f3868f8f23a6e1f57dd8",
			"Socket":         2,
			"SunCluster":     false,
			"SwapTotal":      4,
			"Type":           "VMWARE",
			"VeritasCluster": false,
			"Version":        "1.6.1",
			"Virtual":        true,
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd3d"),
		},
		{
			"CPUCores":       1,
			"CPUModel":       "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
			"CPUThreads":     2,
			"Cluster":        "Puzzait",
			"CreatedAt":      utils.P("2020-04-07T08:52:59.869+02:00"),
			"Databases":      "",
			"Environment":    "PROD",
			"HostType":       "virtualization",
			"Hostname":       "test-virt",
			"Kernel":         "3.10.0-862.9.1.el7.x86_64",
			"Location":       "Italy",
			"MemTotal":       3,
			"OS":             "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
			"OracleCluster":  false,
			"PhysicalHost":   "s157-cb32c10a56c256746c337e21b3f82402",
			"Socket":         2,
			"SunCluster":     false,
			"SwapTotal":      4,
			"Type":           "VMWARE",
			"VeritasCluster": false,
			"Version":        "1.6.1",
			"Virtual":        true,
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd41"),
		},
	}

	as.EXPECT().
		SearchHosts("full", "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchHosts_JSONUnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?mode=sadfsad", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_JSONUnprocessableEntity2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?sort-desc=sadfsad", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_JSONUnprocessableEntity3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?page=sadfsad", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_JSONUnprocessableEntity4(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?size=sadfsad", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_JSONUnprocessableEntity5(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?older-than=sadfsad", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_JSONInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().
		SearchHosts("full", "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchHosts_LMSSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		{
			"ConnectString":            "",
			"CoresPerProcessor":        float64(1),
			"DBInstanceName":           "ERCOLE",
			"Environment":              "TST",
			"Features":                 "Diagnostics Pack",
			"Notes":                    "",
			"OperatingSystem":          "Red Hat Enterprise Linux Server release 7.6 (Maipo)",
			"PhysicalCores":            float64(2),
			"PhysicalServerName":       "",
			"PluggableDatabaseName":    "",
			"ProcessorModel":           "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
			"ProcessorSpeed":           "2.53GHz",
			"Processors":               float64(2),
			"ProductEdition":           "Enterprise",
			"ProductVersion":           "12",
			"RacNodeNames":             "",
			"ServerPurchaseDate":       "",
			"ThreadsPerCore":           int32(2),
			"VirtualServerName":        "itl-csllab-112.sorint.localpippo",
			"VirtualizationTechnology": "VMWARE",
			"_id":                      utils.Str2oid("5e96ade270c184faca93fe20"),
		},
		{
			"ConnectString":            "",
			"CoresPerProcessor":        float64(4),
			"DBInstanceName":           "rudeboy-fb3160a04ffea22b55555bbb58137f77 007bond-f260462ca34bbd17deeda88f042e42a1 jacket-d4a157354d91bfc68fce6f45546d8f3d allstate-9a6a2a820a3f61aeb345a834abf40fba 4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Environment":              "SVIL",
			"Features":                 "",
			"Notes":                    "",
			"OperatingSystem":          "Red Hat Enterprise Linux Server release 5.5 (Tikanga)",
			"PhysicalCores":            float64(8),
			"PhysicalServerName":       "publicitate-36d06ca83eafa454423d2097f4965517",
			"PluggableDatabaseName":    "",
			"ProcessorModel":           "Intel(R) Xeon(R) CPU           X5570  @ 2.93GHz",
			"ProcessorSpeed":           "2.93GHz",
			"Processors":               float64(2),
			"ProductEdition":           "Enterprise",
			"ProductVersion":           "11",
			"RacNodeNames":             "",
			"ServerPurchaseDate":       "",
			"ThreadsPerCore":           int32(2),
			"VirtualServerName":        "",
			"VirtualizationTechnology": "PH",
			"_id":                      utils.Str2oid("5e96ade270c184faca93fe25"),
		},
	}

	as.EXPECT().
		SearchHosts("lms", "foobar", "Processors", true, -1, -1, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?search=foobar&sort-by=Processors&sort-desc=true&location=Italy&environment=TST&&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.oracle.lms+vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := xlsx.Open(rr.Body)
	require.NoError(t, err)
	sh := sp.SheetByName("Database_&_EBS")
	require.NotNil(t, sh)
	assert.Equal(t, "", sh.Cell(0, 3).String())
	assert.Equal(t, "itl-csllab-112.sorint.localpippo", sh.Cell(1, 3).String())
	assert.Equal(t, "VMWARE", sh.Cell(2, 3).String())
	assert.Equal(t, "ERCOLE", sh.Cell(3, 3).String())
	assert.Equal(t, "", sh.Cell(4, 3).String())
	assert.Equal(t, "", sh.Cell(5, 3).String())
	assert.Equal(t, "12", sh.Cell(7, 3).String())
	assert.Equal(t, "Enterprise", sh.Cell(8, 3).String())
	assert.Equal(t, "TST", sh.Cell(9, 3).String())
	assert.Equal(t, "Diagnostics Pack", sh.Cell(10, 3).String())
	assert.Equal(t, "", sh.Cell(11, 3).String())
	assert.Equal(t, "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz", sh.Cell(12, 3).String())
	AssertXLSXFloat(t, 2, sh.Cell(13, 3))
	AssertXLSXFloat(t, 1, sh.Cell(14, 3))
	AssertXLSXFloat(t, 2, sh.Cell(15, 3))
	AssertXLSXFloat(t, 2, sh.Cell(16, 3))
	assert.Equal(t, "2.53GHz", sh.Cell(17, 3).String())
	assert.Equal(t, "", sh.Cell(18, 3).String())
	assert.Equal(t, "Red Hat Enterprise Linux Server release 7.6 (Maipo)", sh.Cell(19, 3).String())
	assert.Equal(t, "", sh.Cell(20, 3).String())

	assert.Equal(t, "publicitate-36d06ca83eafa454423d2097f4965517", sh.Cell(0, 4).String())
	assert.Equal(t, "", sh.Cell(1, 4).String())
	assert.Equal(t, "PH", sh.Cell(2, 4).String())
	assert.Equal(t, "rudeboy-fb3160a04ffea22b55555bbb58137f77 007bond-f260462ca34bbd17deeda88f042e42a1 jacket-d4a157354d91bfc68fce6f45546d8f3d allstate-9a6a2a820a3f61aeb345a834abf40fba 4wcqjn-ecf040bdfab7695ab332aef7401f185c", sh.Cell(3, 4).String())
	assert.Equal(t, "", sh.Cell(4, 4).String())
	assert.Equal(t, "", sh.Cell(5, 4).String())
	assert.Equal(t, "11", sh.Cell(7, 4).String())
	assert.Equal(t, "Enterprise", sh.Cell(8, 4).String())
	assert.Equal(t, "SVIL", sh.Cell(9, 4).String())
	assert.Equal(t, "", sh.Cell(10, 4).String())
	assert.Equal(t, "", sh.Cell(11, 4).String())
	assert.Equal(t, "Intel(R) Xeon(R) CPU           X5570  @ 2.93GHz", sh.Cell(12, 4).String())
	AssertXLSXFloat(t, 2, sh.Cell(13, 4))
	AssertXLSXFloat(t, 4, sh.Cell(14, 4))
	AssertXLSXFloat(t, 8, sh.Cell(15, 4))
	AssertXLSXFloat(t, 2, sh.Cell(16, 4))
	assert.Equal(t, "2.93GHz", sh.Cell(17, 4).String())
	assert.Equal(t, "", sh.Cell(18, 4).String())
	assert.Equal(t, "Red Hat Enterprise Linux Server release 5.5 (Tikanga)", sh.Cell(19, 4).String())
	assert.Equal(t, "", sh.Cell(20, 4).String())
}

func TestSearchHosts_LMSUnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?sort-desc=sdfsdf", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.oracle.lms+vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_LMSUnprocessableEntity2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?older-than=sdfsdf", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.oracle.lms+vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_LMSSuccessInternalServerError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().
		SearchHosts("lms", "foobar", "Processors", true, -1, -1, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?search=foobar&sort-by=Processors&sort-desc=true&location=Italy&environment=TST&&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.oracle.lms+vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchHosts_LMSSuccessInternalServerError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		{
			"ConnectString":            "",
			"CoresPerProcessor":        float64(1),
			"DBInstanceName":           "ERCOLE",
			"Environment":              "TST",
			"Features":                 "Diagnostics Pack",
			"Notes":                    "",
			"OperatingSystem":          "Red Hat Enterprise Linux Server release 7.6 (Maipo)",
			"PhysicalCores":            float64(2),
			"PhysicalServerName":       "",
			"PluggableDatabaseName":    "",
			"ProcessorModel":           "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
			"ProcessorSpeed":           "2.53GHz",
			"Processors":               float64(2),
			"ProductEdition":           "Enterprise",
			"ProductVersion":           "12",
			"RacNodeNames":             "",
			"ServerPurchaseDate":       "",
			"ThreadsPerCore":           int32(2),
			"VirtualServerName":        "itl-csllab-112.sorint.localpippo",
			"VirtualizationTechnology": "VMWARE",
			"_id":                      utils.Str2oid("5e96ade270c184faca93fe20"),
		},
		{
			"ConnectString":            "",
			"CoresPerProcessor":        float64(4),
			"DBInstanceName":           "rudeboy-fb3160a04ffea22b55555bbb58137f77 007bond-f260462ca34bbd17deeda88f042e42a1 jacket-d4a157354d91bfc68fce6f45546d8f3d allstate-9a6a2a820a3f61aeb345a834abf40fba 4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Environment":              "SVIL",
			"Features":                 "",
			"Notes":                    "",
			"OperatingSystem":          "Red Hat Enterprise Linux Server release 5.5 (Tikanga)",
			"PhysicalCores":            float64(8),
			"PhysicalServerName":       "publicitate-36d06ca83eafa454423d2097f4965517",
			"PluggableDatabaseName":    "",
			"ProcessorModel":           "Intel(R) Xeon(R) CPU           X5570  @ 2.93GHz",
			"ProcessorSpeed":           "2.93GHz",
			"Processors":               float64(2),
			"ProductEdition":           "Enterprise",
			"ProductVersion":           "11",
			"RacNodeNames":             "",
			"ServerPurchaseDate":       "",
			"ThreadsPerCore":           int32(2),
			"VirtualServerName":        "",
			"VirtualizationTechnology": "PH",
			"_id":                      utils.Str2oid("5e96ade270c184faca93fe25"),
		},
	}

	as.EXPECT().
		SearchHosts("lms", "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.oracle.lms+vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchHosts_XLSXSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		{
			"CPUCores":       float64(24),
			"CPUModel":       "Intel(R) Xeon(R) Platinum 8160 CPU @ 2.10GHz",
			"CPUThreads":     float64(48),
			"Cluster":        nil,
			"CreatedAt":      utils.PDT("2020-04-15T08:46:58.461+02:00"),
			"Databases":      "8888888-d41d8cd98f00b204e9800998ecf8427e",
			"Environment":    "PROD",
			"HostType":       "exadata",
			"Hostname":       "engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804",
			"Kernel":         "4.1.12-124.26.12.el7uek.x86_64",
			"Location":       "Italy",
			"MemTotal":       float64(376),
			"OS":             "Red Hat Enterprise Linux Server release 7.6 (Maipo)",
			"OracleCluster":  true,
			"PhysicalHost":   nil,
			"Socket":         float64(1),
			"SunCluster":     false,
			"SwapTotal":      float64(23),
			"Type":           "PH",
			"VeritasCluster": false,
			"Version":        "latest",
			"Virtual":        false,
			"_id":            utils.Str2oid("5e96ade270c184faca93fe31"),
		},
		{
			"CPUCores":       float64(1),
			"CPUModel":       "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
			"CPUThreads":     float64(2),
			"Cluster":        "Puzzait",
			"CreatedAt":      utils.PDT("2020-04-15T08:46:58.471+02:00"),
			"Databases":      "ERCOLE",
			"Environment":    "TST",
			"HostType":       "oracledb",
			"Hostname":       "test-db",
			"Kernel":         "3.10.0-514.el7.x86_64",
			"Location":       "Germany",
			"MemTotal":       float64(3),
			"OS":             "Red Hat Enterprise Linux Server release 7.6 (Maipo)",
			"OracleCluster":  false,
			"PhysicalHost":   "s157-cb32c10a56c256746c337e21b3f82402",
			"Socket":         float64(2),
			"SunCluster":     false,
			"SwapTotal":      float64(1),
			"Type":           "VMWARE",
			"VeritasCluster": false,
			"Version":        "latest",
			"Virtual":        true,
			"_id":            utils.Str2oid("5e96ade270c184faca93fe36"),
		},
	}

	as.EXPECT().
		SearchHosts("summary", "foobar", "Processors", true, -1, -1, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?search=foobar&sort-by=Processors&sort-desc=true&location=Italy&environment=TST&&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	sp, err := xlsx.Open(rr.Body)
	require.NoError(t, err)
	sh := sp.SheetByName("Hosts")
	require.NotNil(t, sh)
	assert.Equal(t, "engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804", sh.Cell(0, 1).String())
	assert.Equal(t, "PROD", sh.Cell(1, 1).String())
	assert.Equal(t, "exadata", sh.Cell(2, 1).String())
	assert.Equal(t, "latest", sh.Cell(5, 1).String())
	assert.Equal(t, utils.P("2020-04-15T08:46:58.461+02:00").UTC().String(), sh.Cell(6, 1).String())
	assert.Equal(t, "8888888-d41d8cd98f00b204e9800998ecf8427e", sh.Cell(7, 1).String())
	assert.Equal(t, "Red Hat Enterprise Linux Server release 7.6 (Maipo)", sh.Cell(8, 1).String())
	assert.Equal(t, "4.1.12-124.26.12.el7uek.x86_64", sh.Cell(9, 1).String())
	AssertXLSXBool(t, true, sh.Cell(10, 1))
	AssertXLSXBool(t, false, sh.Cell(11, 1))
	AssertXLSXBool(t, false, sh.Cell(12, 1))
	AssertXLSXBool(t, false, sh.Cell(13, 1))
	assert.Equal(t, "PH", sh.Cell(14, 1).String())
	AssertXLSXInt(t, 48, sh.Cell(15, 1))
	AssertXLSXInt(t, 24, sh.Cell(16, 1))
	AssertXLSXInt(t, 1, sh.Cell(17, 1))
	AssertXLSXInt(t, 376, sh.Cell(18, 1))
	AssertXLSXInt(t, 23, sh.Cell(19, 1))
	assert.Equal(t, "Intel(R) Xeon(R) Platinum 8160 CPU @ 2.10GHz", sh.Cell(20, 1).String())

	assert.Equal(t, "test-db", sh.Cell(0, 2).String())
	assert.Equal(t, "TST", sh.Cell(1, 2).String())
	assert.Equal(t, "oracledb", sh.Cell(2, 2).String())
	assert.Equal(t, "Puzzait", sh.Cell(3, 2).String())
	assert.Equal(t, "s157-cb32c10a56c256746c337e21b3f82402", sh.Cell(4, 2).String())
	assert.Equal(t, "latest", sh.Cell(5, 2).String())
	assert.Equal(t, utils.P("2020-04-15T08:46:58.471+02:00").UTC().String(), sh.Cell(6, 2).String())
	assert.Equal(t, "ERCOLE", sh.Cell(7, 2).String())
	assert.Equal(t, "Red Hat Enterprise Linux Server release 7.6 (Maipo)", sh.Cell(8, 2).String())
	assert.Equal(t, "3.10.0-514.el7.x86_64", sh.Cell(9, 2).String())
	AssertXLSXBool(t, false, sh.Cell(10, 2))
	AssertXLSXBool(t, false, sh.Cell(11, 2))
	AssertXLSXBool(t, false, sh.Cell(12, 2))
	AssertXLSXBool(t, true, sh.Cell(13, 2))
	assert.Equal(t, "VMWARE", sh.Cell(14, 2).String())
	AssertXLSXInt(t, 2, sh.Cell(15, 2))
	AssertXLSXInt(t, 1, sh.Cell(16, 2))
	AssertXLSXInt(t, 2, sh.Cell(17, 2))
	AssertXLSXInt(t, 3, sh.Cell(18, 2))
	AssertXLSXInt(t, 1, sh.Cell(19, 2))
	assert.Equal(t, "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz", sh.Cell(20, 2).String())
}

func TestSearchHosts_XLSXUnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?sort-desc=dsasdasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_XLSXUnprocessableEntity2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts?older-than=asasdasd", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchHosts_XLSXInternalServerError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().
		SearchHosts("summary", "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSearchHosts_XLSXInternalServerError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []map[string]interface{}{
		{
			"CPUCores":       float64(24),
			"CPUModel":       "Intel(R) Xeon(R) Platinum 8160 CPU @ 2.10GHz",
			"CPUThreads":     float64(48),
			"Cluster":        nil,
			"CreatedAt":      utils.PDT("2020-04-15T08:46:58.461+02:00"),
			"Databases":      "8888888-d41d8cd98f00b204e9800998ecf8427e",
			"Environment":    "PROD",
			"HostType":       "exadata",
			"Hostname":       "engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804",
			"Kernel":         "4.1.12-124.26.12.el7uek.x86_64",
			"Location":       "Italy",
			"MemTotal":       float64(376),
			"OS":             "Red Hat Enterprise Linux Server release 7.6 (Maipo)",
			"OracleCluster":  true,
			"PhysicalHost":   nil,
			"Socket":         float64(1),
			"SunCluster":     false,
			"SwapTotal":      float64(23),
			"Type":           "PH",
			"VeritasCluster": false,
			"Version":        "latest",
			"Virtual":        false,
			"_id":            utils.Str2oid("5e96ade270c184faca93fe31"),
		},
		{
			"CPUCores":       float64(1),
			"CPUModel":       "Intel(R) Xeon(R) CPU           E5630  @ 2.53GHz",
			"CPUThreads":     float64(2),
			"Cluster":        "Puzzait",
			"CreatedAt":      utils.PDT("2020-04-15T08:46:58.471+02:00"),
			"Databases":      "ERCOLE",
			"Environment":    "TST",
			"HostType":       "oracledb",
			"Hostname":       "test-db",
			"Kernel":         "3.10.0-514.el7.x86_64",
			"Location":       "Germany",
			"MemTotal":       float64(3),
			"OS":             "Red Hat Enterprise Linux Server release 7.6 (Maipo)",
			"OracleCluster":  false,
			"PhysicalHost":   "s157-cb32c10a56c256746c337e21b3f82402",
			"Socket":         float64(2),
			"SunCluster":     false,
			"SwapTotal":      float64(1),
			"Type":           "VMWARE",
			"VeritasCluster": false,
			"Version":        "latest",
			"Virtual":        true,
			"_id":            utils.Str2oid("5e96ade270c184faca93fe36"),
		},
	}

	as.EXPECT().
		SearchHosts("summary", "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchHosts)
	req, err := http.NewRequest("GET", "/hosts", nil)
	require.NoError(t, err)
	req.Header.Add("Accept", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetHost_JSONSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := map[string]interface{}{
		"Archived":    false,
		"Cluster":     "Puzzait",
		"CreatedAt":   utils.P("2020-04-15T08:46:58.466+02:00"),
		"Databases":   "",
		"Environment": "PROD",
		"Extra": map[string]interface{}{
			"Clusters": []interface{}{
				map[string]interface{}{
					"CPU":     140,
					"Name":    "Puzzait",
					"Sockets": 10,
					"Type":    "vmware",
					"VMs": []interface{}{
						map[string]interface{}{
							"CappedCPU":    false,
							"ClusterName":  "Puzzait",
							"Hostname":     "test-virt",
							"Name":         "test-virt",
							"PhysicalHost": "s157-cb32c10a56c256746c337e21b3f82402",
						},
						map[string]interface{}{
							"CappedCPU":    false,
							"ClusterName":  "Puzzait",
							"Hostname":     "test-db",
							"Name":         "test-db",
							"PhysicalHost": "s157-cb32c10a56c256746c337e21b3f82402",
						},
					},
				},
			},
			"Databases": []interface{}{},
			"Filesystems": []interface{}{
				map[string]interface{}{
					"Available":  "4.6G",
					"Filesystem": "/dev/mapper/vg_os-lv_root",
					"FsType":     "xfs",
					"MountedOn":  "/",
					"Size":       "8.0G",
					"Used":       "3.5G",
					"UsedPerc":   "43%",
				},
			},
		},
		"HostDataSchemaVersion": 3,
		"HostType":              "virtualization",
		"Hostname":              "test-virt",
		"Info": map[string]interface{}{
			"AixCluster":     false,
			"CPUCores":       1,
			"CPUModel":       "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
			"CPUThreads":     2,
			"Environment":    "PROD",
			"Hostname":       "test-virt",
			"Kernel":         "3.10.0-862.9.1.el7.x86_64",
			"Location":       "Italy",
			"MemoryTotal":    3,
			"OS":             "Red Hat Enterprise Linux Server release 7.5 (Maipo)",
			"OracleCluster":  false,
			"Socket":         2,
			"SunCluster":     false,
			"SwapTotal":      4,
			"Type":           "VMWARE",
			"VeritasCluster": false,
			"Virtual":        true,
		},
		"Location":      "Italy",
		"PhysicalHost":  "s157-cb32c10a56c256746c337e21b3f82402",
		"SchemaVersion": 1,
		"Schemas":       "",
		"ServerVersion": "latest",
		"Version":       "1.6.1",
		"_id":           utils.Str2oid("5e96ade270c184faca93fe34"),
	}

	as.EXPECT().
		GetHost("foobar", utils.P("2020-06-10T11:54:59Z"), false).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar?older-than=2020-06-10T11%3A54%3A59Z", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestGetHost_JSONFailUnprocessableEntity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar?older-than=fgfggf", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetHost_JSONFailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().
		GetHost("foobar", utils.MAX_TIME, false).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetHost_JSONFailNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().
		GetHost("foobar", utils.MAX_TIME, false).
		Return(nil, utils.AerrHostNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetHost_MongoJSONSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	res := utils.LoadFixtureHostDataMap(t, "../../fixture/test_dataservice_mongohostdata_02.json")
	expectedRes, err := ioutil.ReadFile("../../fixture/test_dataservice_mongohostdata_02.json")
	require.NoError(t, err)

	as.EXPECT().
		GetHost("foobar", utils.P("2020-06-10T11:54:59Z"), true).
		Return(res, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar?older-than=2020-06-10T11%3A54%3A59Z", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	req.Header.Add("Accept", "application/vnd.ercole.mongohostdata+json")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, string(expectedRes), rr.Body.String())
}

func TestGetHost_MongoJSONFailUnprocessableEntity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar?older-than=fgfggf", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	req.Header.Add("Accept", "application/vnd.ercole.mongohostdata+json")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetHost_MongoJSONFailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().
		GetHost("foobar", utils.MAX_TIME, true).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	req.Header.Add("Accept", "application/vnd.ercole.mongohostdata+json")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetHost_MongoJSONFailNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().
		GetHost("foobar", utils.MAX_TIME, true).
		Return(nil, utils.AerrHostNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHost)
	req, err := http.NewRequest("GET", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	req.Header.Add("Accept", "application/vnd.ercole.mongohostdata+json")
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestListLocations_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []string{"Italy", "German", "France"}

	as.EXPECT().
		ListLocations("Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListLocations)
	req, err := http.NewRequest("GET", "/locations?environment=TST&location=Italy&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestListLocations_FailUnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListLocations)
	req, err := http.NewRequest("GET", "/locations?older-than=dfsgdfsg", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestListLocations_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().
		ListLocations("", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListLocations)
	req, err := http.NewRequest("GET", "/locations", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestListEnvironments_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []string{"TST", "PRD", "DEV"}

	as.EXPECT().
		ListEnvironments("Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListEnvironments)
	req, err := http.NewRequest("GET", "/environments?environment=TST&location=Italy&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestListEnvironments_FailUnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListEnvironments)
	req, err := http.NewRequest("GET", "/environments?older-than=dfsgdfsg", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestListEnvironments_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().
		ListEnvironments("", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ListEnvironments)
	req, err := http.NewRequest("GET", "/environments", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestArchiveHost_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().ArchiveHost("foobar").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ArchiveHost)
	req, err := http.NewRequest("DELETE", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestArchiveHost_FailReadOnly(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: true,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ArchiveHost)
	req, err := http.NewRequest("DELETE", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestArchiveHost_FailNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().ArchiveHost("foobar").Return(utils.AerrHostNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ArchiveHost)
	req, err := http.NewRequest("DELETE", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestArchiveHost_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	as.EXPECT().ArchiveHost("foobar").Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.ArchiveHost)
	req, err := http.NewRequest("DELETE", "/hosts/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "foobar",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
