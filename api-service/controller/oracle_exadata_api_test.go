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
	"testing"

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchOracleExadata_SuccessPaged(t *testing.T) {
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
		"content": []interface{}{
			map[string]interface{}{
				"CreatedAt": utils.P("2020-04-07T08:52:59.865+02:00"),
				"DBServers": []interface{}{
					map[string]interface{}{
						"RunningCPUCount":    48,
						"TotalCPUCount":      48,
						"SwVersion":          "19.2.4.0.0.190709",
						"Hostname":           "zombie-0d1347d47a10b673a4df7aeeecc24a8a",
						"Memory":             "376GB",
						"Model":              "X7-2",
						"RunningPowerSupply": 2,
						"TotalPowerSupply":   2,
						"TempActual":         "24.0",
					},
					map[string]interface{}{
						"RunningCPUCount":    48,
						"TotalCPUCount":      48,
						"SwVersion":          "19.2.4.0.0.190709",
						"Hostname":           "kantoor-43a6cdc54bb211eb127bca5c6651950c",
						"Memory":             "376GB",
						"Model":              "X7-2",
						"RunningPowerSupply": 2,
						"TotalPowerSupply":   2,
						"TempActual":         "24.0",
					},
				},
				"Environment": "PROD",
				"Hostname":    "engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804",
				"IBSwitches": []interface{}{
					map[string]interface{}{
						"SwVersion": "2.2.13-2.190326",
						"Hostname":  "off-df8b95a01746a464e69203c840a6a46a",
						"Model":     "SUN_DCS_36p",
					},
					map[string]interface{}{
						"SwVersion": "2.2.13-2.190326",
						"Hostname":  "aspen-8d1d1b210625b1f1024b686135f889a1",
						"Model":     "SUN_DCS_36p",
					},
				},
				"Location": "Italy",
				"StorageServers": []interface{}{
					map[string]interface{}{
						"RunningCPUCount":    20,
						"TotalCPUCount":      40,
						"SwVersion":          "19.2.4.0.0.190709",
						"Hostname":           "s75-c2449b0e89e5a0b38401636eaa07abd5",
						"Memory":             "188GB",
						"Model":              "X7-2L_High_Capacity",
						"RunningPowerSupply": 2,
						"TotalPowerSupply":   2,
						"TempActual":         "23.0",
					},
					map[string]interface{}{
						"RunningCPUCount":    20,
						"TotalCPUCount":      40,
						"SwVersion":          "19.2.4.0.0.190709",
						"Hostname":           "itl-b22fa37cad1326aba990cdec7facace2",
						"Memory":             "188GB",
						"Model":              "X7-2L_High_Capacity",
						"RunningPowerSupply": 2,
						"TotalPowerSupply":   2,
						"TempActual":         "24.0",
					},
				},
				"_id": utils.Str2oid("5e8c234b24f648a08585bd3e"),
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

	resFromService := []interface{}{
		expectedRes,
	}

	as.EXPECT().
		SearchOracleExadata(true, "foobar", "Hostname", true, 2, 3, "Italy", "TST", utils.P("2020-06-10T11:54:59Z")).
		Return(resFromService, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata?full=true&search=foobar&sort-by=Hostname&sort-desc=true&page=2&size=3&location=Italy&environment=TST&older-than=2020-06-10T11%3A54%3A59Z", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchOracleExadata_SuccessUnpaged(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAPIServiceInterface(mockCtrl)
	ac := APIController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     utils.NewLogger("TEST"),
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"CreatedAt": utils.P("2020-04-07T08:52:59.865+02:00"),
			"DBServers": []interface{}{
				map[string]interface{}{
					"RunningCPUCount":    48,
					"TotalCPUCount":      48,
					"SwVersion":          "19.2.4.0.0.190709",
					"Hostname":           "zombie-0d1347d47a10b673a4df7aeeecc24a8a",
					"Memory":             "376GB",
					"Model":              "X7-2",
					"RunningPowerSupply": 2,
					"TotalPowerSupply":   2,
					"TempActual":         "24.0",
				},
				map[string]interface{}{
					"RunningCPUCount":    48,
					"TotalCPUCount":      48,
					"SwVersion":          "19.2.4.0.0.190709",
					"Hostname":           "kantoor-43a6cdc54bb211eb127bca5c6651950c",
					"Memory":             "376GB",
					"Model":              "X7-2",
					"RunningPowerSupply": 2,
					"TotalPowerSupply":   2,
					"TempActual":         "24.0",
				},
			},
			"Environment": "PROD",
			"Hostname":    "engelsiz-ee2ceb8e1e7fc19e4aeccbae135e2804",
			"IBSwitches": []interface{}{
				map[string]interface{}{
					"SwVersion": "2.2.13-2.190326",
					"Hostname":  "off-df8b95a01746a464e69203c840a6a46a",
					"Model":     "SUN_DCS_36p",
				},
				map[string]interface{}{
					"SwVersion": "2.2.13-2.190326",
					"Hostname":  "aspen-8d1d1b210625b1f1024b686135f889a1",
					"Model":     "SUN_DCS_36p",
				},
			},
			"Location": "Italy",
			"StorageServers": []interface{}{
				map[string]interface{}{
					"RunningCPUCount":    20,
					"TotalCPUCount":      40,
					"SwVersion":          "19.2.4.0.0.190709",
					"Hostname":           "s75-c2449b0e89e5a0b38401636eaa07abd5",
					"Memory":             "188GB",
					"Model":              "X7-2L_High_Capacity",
					"RunningPowerSupply": 2,
					"TotalPowerSupply":   2,
					"TempActual":         "23.0",
				},
				map[string]interface{}{
					"RunningCPUCount":    20,
					"TotalCPUCount":      40,
					"SwVersion":          "19.2.4.0.0.190709",
					"Hostname":           "itl-b22fa37cad1326aba990cdec7facace2",
					"Memory":             "188GB",
					"Model":              "X7-2L_High_Capacity",
					"RunningPowerSupply": 2,
					"TotalPowerSupply":   2,
					"TempActual":         "24.0",
				},
			},
			"_id": utils.Str2oid("5e8c234b24f648a08585bd3e"),
		},
	}

	as.EXPECT().
		SearchOracleExadata(false, "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(expectedRes, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, utils.ToJSON(expectedRes), rr.Body.String())
}

func TestSearchOracleExadata_FailUnprocessableEntity1(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata?full=dasads", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleExadata_FailUnprocessableEntity2(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata?sort-desc=dasads", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleExadata_FailUnprocessableEntity3(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata?page=dasads", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleExadata_FailUnprocessableEntity4(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata?size=dasads", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleExadata_FailUnprocessableEntity5(t *testing.T) {
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
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata?older-than=dasads", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestSearchOracleExadata_FailInternalServerError(t *testing.T) {
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
		SearchOracleExadata(false, "", "", false, -1, -1, "", "", utils.MAX_TIME).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.SearchOracleExadata)
	req, err := http.NewRequest("GET", "/exadata", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
