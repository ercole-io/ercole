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
	"io/ioutil"
	"net/http"

	"github.com/ercole-io/ercole/api-service/apimodel"
	"github.com/ercole-io/ercole/utils"
	"gopkg.in/square/go-jose.v2/json"
)

// GetOracleDatabaseAgreementPartsList return the list of Oracle/Database agreement parts
func (ctrl *APIController) GetOracleDatabaseAgreementPartsList(w http.ResponseWriter, r *http.Request) {
	data, err := ctrl.Service.GetOracleDatabaseAgreementPartsList()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, data)
}

// AddOracleDatabaseAgreements add some agreements
func (ctrl *APIController) AddOracleDatabaseAgreements(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var aerr utils.AdvancedErrorInterface
	var req apimodel.OracleDatabaseAgreementsAddRequest

	//Read all bytes for the request
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusBadRequest)))
		return
	}
	defer r.Body.Close()

	//Unmarshal it to req
	if err := json.Unmarshal(raw, &req); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	//Add it!
	var ids interface{}
	if ids, aerr = ctrl.Service.AddOracleDatabaseAgreements(req); aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the created id
	utils.WriteJSONResponse(w, http.StatusOK, ids)
}

// SearchOracleDatabaseAgreements search Oracle/Database agreements data using the filters in the request
func (ctrl *APIController) SearchOracleDatabaseAgreements(w http.ResponseWriter, r *http.Request) {
	var err utils.AdvancedErrorInterface
	var search string

	search = r.URL.Query().Get("search")

	//Get the filters
	searchOracleDatabaseAgreementsFilters, err := GetSearchOracleDatabaseAgreementsFilters(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	response, err := ctrl.Service.SearchOracleDatabaseAgreements(search, searchOracleDatabaseAgreementsFilters)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetSearchOracleDatabaseAgreementsFilters return the Oracle/Database agreement search filters in the request
func GetSearchOracleDatabaseAgreementsFilters(r *http.Request) (apimodel.SearchOracleDatabaseAgreementsFilters, utils.AdvancedErrorInterface) {
	var aerr utils.AdvancedErrorInterface

	filters := apimodel.SearchOracleDatabaseAgreementsFilters{}
	filters.AgreementID = r.URL.Query().Get("agreement-id")
	filters.PartID = r.URL.Query().Get("part-id")
	filters.ItemDescription = r.URL.Query().Get("item-description")
	filters.CSI = r.URL.Query().Get("csi")
	filters.Metrics = r.URL.Query().Get("metrics")
	filters.ReferenceNumber = r.URL.Query().Get("reference-number")
	filters.Unlimited = r.URL.Query().Get("unlimited")
	if filters.Unlimited == "" {
		filters.Unlimited = "NULL"
	} else if filters.Unlimited != "true" && filters.Unlimited != "false" && filters.Unlimited != "NULL" {
		return apimodel.SearchOracleDatabaseAgreementsFilters{}, utils.NewAdvancedErrorPtr(errors.New("Invalid value for unlimited"), http.StatusText(http.StatusUnprocessableEntity))
	}
	filters.CatchAll = r.URL.Query().Get("catch-all")
	if filters.CatchAll == "" {
		filters.CatchAll = "NULL"
	} else if filters.CatchAll != "true" && filters.CatchAll != "false" && filters.CatchAll != "NULL" {
		return apimodel.SearchOracleDatabaseAgreementsFilters{}, utils.NewAdvancedErrorPtr(errors.New("Invalid value for catch-all"), http.StatusText(http.StatusUnprocessableEntity))
	}
	if filters.LicensesCountLTE, aerr = utils.Str2int(r.URL.Query().Get("licenses-count-lte"), -1); aerr != nil {
		return apimodel.SearchOracleDatabaseAgreementsFilters{}, aerr
	}
	if filters.LicensesCountGTE, aerr = utils.Str2int(r.URL.Query().Get("licenses-count-gte"), -1); aerr != nil {
		return apimodel.SearchOracleDatabaseAgreementsFilters{}, aerr
	}
	if filters.UsersCountLTE, aerr = utils.Str2int(r.URL.Query().Get("users-count-lte"), -1); aerr != nil {
		return apimodel.SearchOracleDatabaseAgreementsFilters{}, aerr
	}
	if filters.UsersCountGTE, aerr = utils.Str2int(r.URL.Query().Get("users-count-gte"), -1); aerr != nil {
		return apimodel.SearchOracleDatabaseAgreementsFilters{}, aerr
	}
	if filters.AvailableCountLTE, aerr = utils.Str2int(r.URL.Query().Get("available-count-lte"), -1); aerr != nil {
		return apimodel.SearchOracleDatabaseAgreementsFilters{}, aerr
	}
	if filters.AvailableCountGTE, aerr = utils.Str2int(r.URL.Query().Get("available-count-gte"), -1); aerr != nil {
		return apimodel.SearchOracleDatabaseAgreementsFilters{}, aerr
	}
	return filters, nil
}
