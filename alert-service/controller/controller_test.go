package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/utils"
	"github.com/golang/mock/gomock"
)

//go:generate mockgen -source ../service/service.go -destination=fake_service.go -package=controller

//Common data
var errMock error = errors.New("MockError")
var aerrMock utils.AdvancedErrorInterface = utils.NewAdvancedErrorPtr(errMock, "mock")

//p parse the string s and return the equivalent time
func p(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

//btc break the time continuum and return a function that return the time t
func btc(t time.Time) func() time.Time {
	return func() time.Time {
		return t
	}
}

func str2oid(str string) primitive.ObjectID {
	val, _ := primitive.ObjectIDFromHex(str)
	return val
}

func TestHostDataInsertion_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAlertServiceInterface(mockCtrl)
	aqc := AlertQueueController{
		TimeNow: btc(p("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
	}
	as.EXPECT().HostDataInsertion(str2oid("5dc3f534db7e81a98b726a52")).Return(nil).Times(1)
	as.EXPECT().HostDataInsertion(gomock.Any()).Times(0)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(aqc.HostDataInsertion)
	req, err := http.NewRequest("GET", "/queue/host-data-insertion/5dc3f534db7e81a98b726a52", nil)
	require.Nil(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "5dc3f534db7e81a98b726a52",
	})

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestHostDataInsertion_RequestError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAlertServiceInterface(mockCtrl)
	aqc := AlertQueueController{
		TimeNow: btc(p("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
	}
	as.EXPECT().HostDataInsertion(gomock.Any()).Times(0)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(aqc.HostDataInsertion)
	req, err := http.NewRequest("GET", "/queue/host-data-insertion/pippo", nil)
	require.Nil(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "pippo",
	})

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestHostDataInsertion_ServiceError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	as := NewMockAlertServiceInterface(mockCtrl)
	aqc := AlertQueueController{
		TimeNow: btc(p("2019-11-05T14:02:03Z")),
		Service: as,
		Config: config.Configuration{
			AlertService: config.AlertService{
				FreshnessCheckJob: config.FreshnessCheckJob{
					DaysThreshold: 10,
				},
			},
		},
	}
	as.EXPECT().HostDataInsertion(str2oid("5dc3f534db7e81a98b726a52")).Return(aerrMock).Times(1)
	as.EXPECT().HostDataInsertion(gomock.Any()).Times(0)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(aqc.HostDataInsertion)
	req, err := http.NewRequest("GET", "/queue/host-data-insertion/5dc3f534db7e81a98b726a52", nil)
	require.Nil(t, err)
	req = mux.SetURLVars(req, map[string]string{
		"id": "5dc3f534db7e81a98b726a52",
	})

	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
}
