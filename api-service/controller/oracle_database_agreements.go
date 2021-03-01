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

// AddAssociatedLicenseTypeToOracleDbAgreement add associated part to an existing agreement else it will create it
func (ctrl *APIController) AddAssociatedLicenseTypeToOracleDbAgreement(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var req dto.AssociatedLicenseTypeInOracleDbAgreementRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest,
			utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusBadRequest)))
		return
	}

	if req.ID != "" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest,
			utils.NewAdvancedErrorPtr(errors.New("ID must be empty to add a new AssociatedLicenseType"), http.StatusText(http.StatusBadRequest)))
		return
	}

	if err := req.Check(); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	id, aerr := ctrl.Service.AddAssociatedLicenseTypeToOracleDbAgreement(req)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, id)
}

// UpdateAssociatedLicenseTypeOfOracleDbAgreement edit an agreement
func (ctrl *APIController) UpdateAssociatedLicenseTypeOfOracleDbAgreement(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var req dto.AssociatedLicenseTypeInOracleDbAgreementRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest,
			utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusBadRequest)))
		return
	}

	if err := req.Check(); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	err := ctrl.Service.UpdateAssociatedLicenseTypeOfOracleDbAgreement(req)
	if err == utils.AerrOracleDatabaseAgreementNotFound ||
		err == utils.AerrOracleDatabaseLicenseTypeIDNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

// SearchAssociatedLicenseTypesInOracleDatabaseAgreements search Oracle/Database agreements
func (ctrl *APIController) SearchAssociatedLicenseTypesInOracleDatabaseAgreements(w http.ResponseWriter, r *http.Request) {
	var err utils.AdvancedErrorInterface

	searchOracleDatabaseAgreementsFilters, err := parseSearchOracleDatabaseAgreementsFilters(r.URL.Query())
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	response, err := ctrl.Service.SearchAssociatedLicenseTypesInOracleDatabaseAgreements(searchOracleDatabaseAgreementsFilters)
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
	filters.LicenseTypeID = urlValues.Get("license-type-id")
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

// AddHostToAssociatedLicenseType add an host from AssociatedLicenseType
func (ctrl *APIController) AddHostToAssociatedLicenseType(w http.ResponseWriter, r *http.Request) {
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

	if aerr = ctrl.Service.AddHostToAssociatedLicenseType(id, string(raw)); aerr == utils.AerrOracleDatabaseAgreementNotFound ||
		aerr == utils.AerrNotInClusterHostNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, aerr)
		return
	} else if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

// RemoveHostFromAssociatedLicenseType remove an host from AssociatedLicenseType
func (ctrl *APIController) RemoveHostFromAssociatedLicenseType(w http.ResponseWriter, r *http.Request) {
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

	if aerr = ctrl.Service.RemoveHostFromAssociatedLicenseType(id, hostname); aerr == utils.AerrOracleDatabaseAgreementNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, aerr)
		return
	} else if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

// DeleteAssociatedLicenseTypeFromOracleDatabaseAgreement delete AssociatedLicenseType from an OracleDatabaseAgreement
func (ctrl *APIController) DeleteAssociatedLicenseTypeFromOracleDatabaseAgreement(w http.ResponseWriter, r *http.Request) {
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

	if aerr = ctrl.Service.DeleteAssociatedLicenseTypeFromOracleDatabaseAgreement(id); aerr == utils.AerrOracleDatabaseAgreementNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, aerr)
		return
	} else if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
