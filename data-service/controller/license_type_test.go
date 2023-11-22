package controller

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestInsertOracleLicenseTypes_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockHostDataServiceInterface(mockCtrl)
	ac := DataController{
		TimeNow: utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Service: as,
		Config:  config.Configuration{},
		Log:     logger.NewLogger("TEST"),
	}

	raw, err := ioutil.ReadFile("../../fixture/test_oracle_database_license_types.json")
	require.NoError(t, err)

	expectedRes := []model.OracleDatabaseLicenseType{
		{
			ID:              "A00001",
			ItemDescription: "Oracle Test Item Description",
			Metric:          "Test Metric",
			Cost:            47500,
			Aliases:         []string{"Test Alias #01", "Test Alias #02"},
			Option:          false,
		},
		{
			ID:              "A00002",
			ItemDescription: "Oracle Test Item Description",
			Metric:          "Test Metric",
			Cost:            47502,
			Aliases:         []string{"Test Alias #01", "Test Alias #02"},
			Option:          false,
		},
	}

	as.EXPECT().SanitizeLicenseTypes(raw).Return(expectedRes, nil)
	as.EXPECT().InsertOracleLicenseTypes(expectedRes).Return(nil)

	handler := http.HandlerFunc(ac.InsertOracleLicenseTypes)
	req, err := http.NewRequest("POST", "/", bytes.NewReader(raw))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}
