package controller

import (
	"github.com/ercole-io/ercole/v2/chart-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHostsHistory_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockChartServiceInterface(mockCtrl)
	ac := ChartController{
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	host := []dto.HostCores{}
	location := "Italy"
	environment := "TST"
	OlderThan := utils.MAX_TIME
	NewerThan := utils.MIN_TIME

	as.EXPECT().GetHostCores(location, environment, OlderThan, NewerThan).
		Return(host, nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHostCores)
	req, err := http.NewRequest("GET", "/?location=Italy&environment=TST", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	expected := map[string]interface{}{
		"coresHistory": host,
	}
	assert.JSONEq(t, utils.ToJSON(expected), rr.Body.String())
}

func TestGetHostsHistoryUnprocessableEntity1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockChartServiceInterface(mockCtrl)
	ac := ChartController{
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHostCores)
	req, err := http.NewRequest("GET", "/?older-than=sadfsad", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetHostsHistoryUnprocessableEntity2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockChartServiceInterface(mockCtrl)
	ac := ChartController{
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHostCores)
	req, err := http.NewRequest("GET", "/?newer-than=sadfsad", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetHostsHistory_InternalServerError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockChartServiceInterface(mockCtrl)
	ac := ChartController{
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	location := ""
	environment := ""
	olderThan := utils.MAX_TIME
	newerThan := utils.MIN_TIME

	as.EXPECT().GetHostCores(location, environment, olderThan, newerThan).
		Return(nil, aerrMock)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ac.GetHostCores)
	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
