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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPatchingFunction_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	id := utils.Str2oid("5e9fee9920d55cbdc35022ad")
	pf := model.PatchingFunction{
		ID:        &id,
		Hostname:  "test-db",
		CreatedAt: utils.P("2020-04-22T07:13:29.873Z"),
		Code:      "\n\t/*\u003cDATABASE_TAGS_ADDER\u003e*/\n\thostdata.Extra.Databases.forEach(function addTag(db) {\n\t\tif (db.Name in vars.Tags) {\n\t\t\tdb.Tags = vars.Tags[db.Name];\n\t\t}\n\t});\n\t/*\u003c/DATABASE_TAGS_ADDER\u003e*/\n",
		Vars: map[string]interface{}{
			"Tags": map[string]interface{}{
				"ERCOLE": []string{"foobar"},
			},
		},
	}

	as.EXPECT().
		GetPatchingFunction("test-db").
		Return(pf, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetPatchingFunction)
	req, err := http.NewRequest("GET", "/hosts/test-db/patching-function", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
	})

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(pf), rr.Body.String())
}

func TestGetPatchingFunction_FailNotFound(t *testing.T) {
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
		GetPatchingFunction("test-db").
		Return(nil, utils.AerrPatchingFunctionNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetPatchingFunction)
	req, err := http.NewRequest("GET", "/hosts/test-db/patching-function", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
	})

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetPatchingFunction_FailInternalServerError(t *testing.T) {
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
		GetPatchingFunction("test-db").
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetPatchingFunction)
	req, err := http.NewRequest("GET", "/hosts/test-db/patching-function", nil)
	require.NoError(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
	})

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSetPatchingFunction_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly:                              false,
				EnableInsertingCustomPatchingFunction: true,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	pf := model.PatchingFunction{
		Code: "//OK",
		Vars: map[string]interface{}{
			"Tags": map[string]interface{}{
				"ERCOLE": []interface{}{"foobar"},
			},
		},
		ID: nil,
	}

	as.EXPECT().SetPatchingFunction("test-db", pf).Return(utils.Str2oid("5e9ff545e4c53a19c79eadfd"), nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetPatchingFunction)
	req, err := http.NewRequest("PUT", "/hosts/test.db/patching-function", strings.NewReader(utils.ToJSON(pf)))
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestSetPatchingFunction_FailReadOnly(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly:                              true,
				EnableInsertingCustomPatchingFunction: true,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	pf := model.PatchingFunction{
		Code: "//OK",
		Vars: map[string]interface{}{
			"Tags": map[string]interface{}{
				"ERCOLE": []interface{}{"foobar"},
			},
		},
		ID: nil,
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetPatchingFunction)
	req, err := http.NewRequest("PUT", "/hosts/test.db/patching-function", strings.NewReader(utils.ToJSON(pf)))
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestSetPatchingFunction_FailDisabled(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly:                              false,
				EnableInsertingCustomPatchingFunction: false,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	pf := model.PatchingFunction{
		Code: "//OK",
		Vars: map[string]interface{}{
			"Tags": map[string]interface{}{
				"ERCOLE": []interface{}{"foobar"},
			},
		},
		ID: nil,
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetPatchingFunction)
	req, err := http.NewRequest("PUT", "/hosts/test.db/patching-function", strings.NewReader(utils.ToJSON(pf)))
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestSetPatchingFunction_FailUnprocessableEntity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly:                              false,
				EnableInsertingCustomPatchingFunction: true,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetPatchingFunction)
	req, err := http.NewRequest("PUT", "/hosts/test.db/patching-function", strings.NewReader("sddassd"))
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSetPatchingFunction_FailNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly:                              false,
				EnableInsertingCustomPatchingFunction: true,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	pf := model.PatchingFunction{
		Code: "//OK",
		Vars: map[string]interface{}{
			"Tags": map[string]interface{}{
				"ERCOLE": []interface{}{"foobar"},
			},
		},
		ID: nil,
	}

	as.EXPECT().SetPatchingFunction("test-db", pf).Return(nil, utils.AerrHostNotFound)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetPatchingFunction)
	req, err := http.NewRequest("PUT", "/hosts/test.db/patching-function", strings.NewReader(utils.ToJSON(pf)))
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestSetPatchingFunction_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly:                              false,
				EnableInsertingCustomPatchingFunction: true,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	pf := model.PatchingFunction{
		Code: "//OK",
		Vars: map[string]interface{}{
			"Tags": map[string]interface{}{
				"ERCOLE": []interface{}{"foobar"},
			},
		},
		ID: nil,
	}

	as.EXPECT().SetPatchingFunction("test-db", pf).Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetPatchingFunction)
	req, err := http.NewRequest("PUT", "/hosts/test.db/patching-function", strings.NewReader(utils.ToJSON(pf)))
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestAddTagToDatabase_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: false,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AddTagToDatabase("test-db", "ERCOLE", "awesome").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddTagToDatabase)
	req, err := http.NewRequest("PUT", "/hosts/test-db/databases/ERCOLE/tags", strings.NewReader("awesome"))
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
		"dbname":   "ERCOLE",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestAddTagToDatabase_FailReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.AddTagToDatabase)
	req, err := http.NewRequest("PUT", "/hosts/test-db/databases/ERCOLE/tags", strings.NewReader("awesome"))
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
		"dbname":   "ERCOLE",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestAddTagToDatabase_FailBadRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: false,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddTagToDatabase)
	req, err := http.NewRequest("PUT", "/hosts/test-db/databases/ERCOLE/tags", &FailingReader{})
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
		"dbname":   "ERCOLE",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAddTagToDatabase_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: false,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().AddTagToDatabase("test-db", "ERCOLE", "awesome").Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.AddTagToDatabase)
	req, err := http.NewRequest("PUT", "/hosts/test-db/databases/ERCOLE/tags", strings.NewReader("awesome"))
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
		"dbname":   "ERCOLE",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestDeleteTagOfDatabase_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: false,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().DeleteTagOfDatabase("test-db", "ERCOLE", "awesome").Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteTagOfDatabase)
	req, err := http.NewRequest("DELETE", "/hosts/test-db/databases/ERCOLE/tags/awesome", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
		"dbname":   "ERCOLE",
		"tagname":  "awesome",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteTagOfDatabase_FailReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.DeleteTagOfDatabase)
	req, err := http.NewRequest("DELETE", "/hosts/test-db/databases/ERCOLE/tags/awesome", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
		"dbname":   "ERCOLE",
		"tagname":  "awesome",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestDeleteTagOfDatabase_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: false,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().DeleteTagOfDatabase("test-db", "ERCOLE", "awesome").Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.DeleteTagOfDatabase)
	req, err := http.NewRequest("DELETE", "/hosts/test-db/databases/ERCOLE/tags/awesome", nil)
	req = mux.SetURLVars(req, map[string]string{
		"hostname": "test-db",
		"dbname":   "ERCOLE",
		"tagname":  "awesome",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSetLicenseModifier_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: false,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().SetLicenseModifier("test-db", "ERCOLE", "Diagnostics Pack", 20).Return(nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetLicenseModifier)
	req, err := http.NewRequest("PUT", "/hosts/test-db/databases/ERCOLE/licenses/Diagnostics%20Pack", strings.NewReader("20"))
	req = mux.SetURLVars(req, map[string]string{
		"hostname":    "test-db",
		"dbname":      "ERCOLE",
		"licenseName": "Diagnostics Pack",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestSetLicenseModifier_FailReadOnly(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SetLicenseModifier)
	req, err := http.NewRequest("PUT", "/hosts/test-db/databases/ERCOLE/licenses/Diagnostics%20Pack", strings.NewReader("20"))
	req = mux.SetURLVars(req, map[string]string{
		"hostname":    "test-db",
		"dbname":      "ERCOLE",
		"licenseName": "Diagnostics Pack",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestSetLicenseModifier_FailBadRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: false,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetLicenseModifier)
	req, err := http.NewRequest("PUT", "/hosts/test-db/databases/ERCOLE/licenses/Diagnostics%20Pack", &FailingReader{})
	req = mux.SetURLVars(req, map[string]string{
		"hostname":    "test-db",
		"dbname":      "ERCOLE",
		"licenseName": "Diagnostics Pack",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSetLicenseModifier_FailUnprocessableEntity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: false,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetLicenseModifier)
	req, err := http.NewRequest("PUT", "/hosts/test-db/databases/ERCOLE/licenses/Diagnostics%20Pack", strings.NewReader("sdfsdfsd"))
	req = mux.SetURLVars(req, map[string]string{
		"hostname":    "test-db",
		"dbname":      "ERCOLE",
		"licenseName": "Diagnostics Pack",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSetLicenseModifier_FailInternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			APIService: config.APIService{
				ReadOnly: false,
			},
		},
		Log: utils.NewLogger("TEST"),
	}

	as.EXPECT().SetLicenseModifier("test-db", "ERCOLE", "Diagnostics Pack", 20).Return(aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SetLicenseModifier)
	req, err := http.NewRequest("PUT", "/hosts/test-db/databases/ERCOLE/licenses/Diagnostics%20Pack", strings.NewReader("20"))
	req = mux.SetURLVars(req, map[string]string{
		"hostname":    "test-db",
		"dbname":      "ERCOLE",
		"licenseName": "Diagnostics Pack",
	})
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
