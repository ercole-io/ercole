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
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddAssociatedPartToOracleDbAgreement add associated part to an existing agreement else it will create it
func (ctrl *APIController) AddAssociatedPartToOracleDbAgreement(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var req dto.AssociatedPartInOracleDbAgreementRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest,
			utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusBadRequest)))
		return
	}

	if req.ID != "" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest,
			utils.NewAdvancedErrorPtr(errors.New("ID must be empty to add a new AssociatedPart"), http.StatusText(http.StatusBadRequest)))
		return
	}

	id, aerr := ctrl.Service.AddAssociatedPartToOracleDbAgreement(req)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, id)
}

// UpdateAssociatedPartOfOracleDbAgreement edit an agreement
func (ctrl *APIController) UpdateAssociatedPartOfOracleDbAgreement(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var req dto.AssociatedPartInOracleDbAgreementRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest,
			utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusBadRequest)))
		return
	}

	err := ctrl.Service.UpdateAssociatedPartOfOracleDbAgreement(req)
	if err == utils.AerrOracleDatabaseAgreementNotFound ||
		err == utils.AerrOracleDatabaseAgreementInvalidPartID {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

// SearchAssociatedPartsInOracleDatabaseAgreements search Oracle/Database agreements
func (ctrl *APIController) SearchAssociatedPartsInOracleDatabaseAgreements(w http.ResponseWriter, r *http.Request) {
	var err utils.AdvancedErrorInterface

	searchOracleDatabaseAgreementsFilters, err := parseSearchOracleDatabaseAgreementsFilters(r.URL.Query())
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	response, err := ctrl.Service.SearchAssociatedPartsInOracleDatabaseAgreements(searchOracleDatabaseAgreementsFilters)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// parseSearchOracleDatabaseAgreementsFilters return the Oracle/Database agreement search filters in the request
func parseSearchOracleDatabaseAgreementsFilters(urlValues url.Values) (dto.SearchOracleDatabaseAgreementsFilter,
	utils.AdvancedErrorInterface) {

	var aerr utils.AdvancedErrorInterface

	filters := dto.SearchOracleDatabaseAgreementsFilter{}

	filters.AgreementID = urlValues.Get("agreement-id")
	filters.PartID = urlValues.Get("part-id")
	filters.ItemDescription = urlValues.Get("item-description")
	filters.CSI = urlValues.Get("csi")
	filters.Metric = urlValues.Get("metrics")
	filters.ReferenceNumber = urlValues.Get("reference-number")

	filters.Unlimited = urlValues.Get("unlimited")
	if filters.Unlimited == "" {
		filters.Unlimited = "NULL"
	} else if filters.Unlimited != "true" && filters.Unlimited != "false" && filters.Unlimited != "NULL" {
		return dto.SearchOracleDatabaseAgreementsFilter{},
			utils.NewAdvancedErrorPtr(errors.New("Invalid value for unlimited"), http.StatusText(http.StatusUnprocessableEntity))
	}

	filters.CatchAll = urlValues.Get("catch-all")
	if filters.CatchAll == "" {
		filters.CatchAll = "NULL"
	} else if filters.CatchAll != "true" && filters.CatchAll != "false" && filters.CatchAll != "NULL" {
		return dto.SearchOracleDatabaseAgreementsFilter{},
			utils.NewAdvancedErrorPtr(errors.New("Invalid value for catch-all"), http.StatusText(http.StatusUnprocessableEntity))
	}

	if filters.LicensesCountLTE, aerr = utils.Str2int(urlValues.Get("licenses-count-lte"), -1); aerr != nil {
		return dto.SearchOracleDatabaseAgreementsFilter{}, aerr
	}

	if filters.LicensesCountGTE, aerr = utils.Str2int(urlValues.Get("licenses-count-gte"), -1); aerr != nil {
		return dto.SearchOracleDatabaseAgreementsFilter{}, aerr
	}

	if filters.UsersCountLTE, aerr = utils.Str2int(urlValues.Get("users-count-lte"), -1); aerr != nil {
		return dto.SearchOracleDatabaseAgreementsFilter{}, aerr
	}

	if filters.UsersCountGTE, aerr = utils.Str2int(urlValues.Get("users-count-gte"), -1); aerr != nil {
		return dto.SearchOracleDatabaseAgreementsFilter{}, aerr
	}

	if filters.AvailableCountLTE, aerr = utils.Str2int(urlValues.Get("available-count-lte"), -1); aerr != nil {
		return dto.SearchOracleDatabaseAgreementsFilter{}, aerr
	}

	if filters.AvailableCountGTE, aerr = utils.Str2int(urlValues.Get("available-count-gte"), -1); aerr != nil {
		return dto.SearchOracleDatabaseAgreementsFilter{}, aerr
	}

	return filters, nil
}

// AddHostToAssociatedPart add an host from AssociatedPart
func (ctrl *APIController) AddHostToAssociatedPart(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var err error
	var aerr utils.AdvancedErrorInterface
	var id primitive.ObjectID

	if id, err = primitive.ObjectIDFromHex(mux.Vars(r)["id"]); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity,
			utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusBadRequest)))
		return
	}
	defer r.Body.Close()

	if aerr = ctrl.Service.AddHostToAssociatedPart(id, string(raw)); aerr == utils.AerrOracleDatabaseAgreementNotFound ||
		aerr == utils.AerrNotInClusterHostNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, aerr)
		return
	} else if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

// RemoveHostFromAssociatedPart remove an host from AssociatedPart
func (ctrl *APIController) RemoveHostFromAssociatedPart(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var err error
	var aerr utils.AdvancedErrorInterface
	var id primitive.ObjectID
	var hostname string

	if id, err = primitive.ObjectIDFromHex(mux.Vars(r)["id"]); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}
	hostname = mux.Vars(r)["hostname"]

	if aerr = ctrl.Service.RemoveHostFromAssociatedPart(id, hostname); aerr == utils.AerrOracleDatabaseAgreementNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, aerr)
		return
	} else if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

// DeleteAssociatedPartFromOracleDatabaseAgreement delete AssociatedPart from an OracleDatabaseAgreement
func (ctrl *APIController) DeleteAssociatedPartFromOracleDatabaseAgreement(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var err error
	var aerr utils.AdvancedErrorInterface
	var id primitive.ObjectID

	if id, err = primitive.ObjectIDFromHex(mux.Vars(r)["id"]); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	if aerr = ctrl.Service.DeleteAssociatedPartFromOracleDatabaseAgreement(id); aerr == utils.AerrOracleDatabaseAgreementNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, aerr)
		return
	} else if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
