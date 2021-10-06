package controller

import (
	"net/http"
	"time"

	"github.com/ercole-io/ercole/v2/utils"
)

func (ctrl *ChartController) GetHostCores(w http.ResponseWriter, r *http.Request) {
	var err error
	var location string
	var environment string
	var olderThan time.Time
	var newerThan time.Time

	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")
	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}
	if newerThan, err = utils.Str2time(r.URL.Query().Get("newer-than"), utils.MIN_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	hosts, err := ctrl.Service.GetHostCores(location, environment, olderThan, newerThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"coresHistory": hosts,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}
