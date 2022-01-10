// Copyright (c) 2022 Sorint.lab S.p.A.
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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/gddo/httputil"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

// SearchAlerts search alerts using the filters in the request
func (ctrl *APIController) SearchAlerts(w http.ResponseWriter, r *http.Request) {
	var mode string
	var search string
	var sortBy string
	var sortDesc bool
	var pageNumber int
	var pageSize int
	var location string
	var environment string
	var severity string
	var status string
	var from time.Time
	var to time.Time

	var err error

	mode = r.URL.Query().Get("mode")
	if mode == "" {
		mode = "all"
	} else if mode != "all" && mode != "aggregated-code-severity" && mode != "aggregated-category-severity" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewError(errors.New("Invalid mode value"), http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	search = r.URL.Query().Get("search")

	sortBy = r.URL.Query().Get("sort-by")

	if sortDesc, err = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if pageNumber, err = utils.Str2int(r.URL.Query().Get("page"), -1); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if pageSize, err = utils.Str2int(r.URL.Query().Get("size"), -1); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	severity = r.URL.Query().Get("severity")
	if severity != "" && severity != model.AlertSeverityWarning && severity != model.AlertSeverityCritical && severity != model.AlertSeverityInfo {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid severity"), "Invalid  severity"))
		return
	}

	status = r.URL.Query().Get("status")
	if status != "" && status != model.AlertStatusNew && status != model.AlertStatusAck && status != model.AlertStatusDismissed {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid status"), "Invalid  status"))
		return
	}

	if from, err = utils.Str2time(r.URL.Query().Get("from"), utils.MIN_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if to, err = utils.Str2time(r.URL.Query().Get("to"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	contentType := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch contentType {
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		if mode != "all" {
			utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest,
				utils.NewError(fmt.Errorf("only mode 'all' is acceptable for xlsx"), http.StatusText(http.StatusBadRequest)))
			return
		}

		ctrl.searchAlertsXLSX(w, r, from, to)

	default:
		ctrl.searchAlertsJSON(w, r, mode, search, sortBy, sortDesc, pageNumber, pageSize, location, environment, severity, status, from, to)
	}
}

// searchAlertsJSON search alerts using the filters in the request returning it in JSON format
func (ctrl *APIController) searchAlertsJSON(w http.ResponseWriter, r *http.Request,
	mode string, search string, sortBy string, sortDesc bool, pageNumber int, pageSize int,
	location, environment, severity string, status string, from time.Time, to time.Time) {

	response, err := ctrl.Service.SearchAlerts(mode, search, sortBy, sortDesc, pageNumber, pageSize, location, environment, severity, status, from, to)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		utils.WriteJSONResponse(w, http.StatusOK, response)
	} else {
		alerts := response[0]
		utils.WriteJSONResponse(w, http.StatusOK, alerts)
	}
}

// searchAlertsXLSX search alerts using the filters in the request returning it in XLSX format
func (ctrl *APIController) searchAlertsXLSX(w http.ResponseWriter, r *http.Request, from time.Time, to time.Time) {

	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}
	xlsx, err := ctrl.Service.SearchAlertsAsXLSX(from, to, *filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, xlsx)
}

// AckAlerts ack the specified alert in the request
func (ctrl *APIController) AckAlerts(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewError(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	body := struct {
		Ids    []primitive.ObjectID `json:"ids"`
		Filter *dto.AlertsFilter    `json:"filter"`
	}{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if (body.Ids != nil && body.Filter != nil) ||
		(body.Ids == nil && body.Filter == nil) {
		err := errors.New("you must send ids or filter (but not both: they are mutually exclusive)")
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	var filter dto.AlertsFilter
	if body.Filter != nil {
		filter = *body.Filter
	} else {
		filter.IDs = body.Ids
	}

	err := ctrl.Service.AckAlerts(filter)
	if errors.Is(err, utils.ErrAlertNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
	} else if errors.Is(err, utils.ErrInvalidAck) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
