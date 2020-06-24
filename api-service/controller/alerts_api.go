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
	"errors"
	"net/http"
	"time"

	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchAlerts search alerts using the filters in the request
func (ctrl *APIController) SearchAlerts(w http.ResponseWriter, r *http.Request) {
	var mode string
	var search string
	var sortBy string
	var sortDesc bool
	var pageNumber int
	var pageSize int
	var severity string
	var status string
	var from time.Time
	var to time.Time

	var err utils.AdvancedErrorInterface
	//parse the query params
	mode = r.URL.Query().Get("mode")
	if mode == "" {
		mode = "all"
	} else if mode != "all" && mode != "aggregated-code-severity" && mode != "aggregated-category-severity" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(errors.New("Invalid mode value"), http.StatusText(http.StatusUnprocessableEntity)))
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
	if severity != "" && severity != model.AlertSeverityMinor && severity != model.AlertSeverityWarning && severity != model.AlertSeverityMajor && severity != model.AlertSeverityCritical && severity != model.AlertSeverityNotice {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(errors.New("invalid severity"), "Invalid  severity"))
		return
	}
	status = r.URL.Query().Get("status")
	if status != "" && status != model.AlertStatusNew && status != model.AlertStatusAck {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(errors.New("invalid status"), "Invalid  status"))
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

	//get the data
	hosts, err := ctrl.Service.SearchAlerts(mode, search, sortBy, sortDesc, pageNumber, pageSize, severity, status, from, to)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, hosts)
	} else {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, hosts[0])
	}
}

// AckAlert ack the specified alert in the request
func (ctrl *APIController) AckAlert(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var id primitive.ObjectID
	var err error
	//Get the id from the path variable
	if id, err = primitive.ObjectIDFromHex(mux.Vars(r)["id"]); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	//set the value
	aerr := ctrl.Service.AckAlert(id)
	if aerr == utils.AerrAlertNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, aerr)
	} else if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
