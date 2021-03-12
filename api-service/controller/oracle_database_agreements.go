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

	id, err := ctrl.Service.AddAssociatedLicenseTypeToOracleDbAgreement(req)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
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
	if err == utils.ErrOracleDatabaseAgreementNotFound ||
		err == utils.ErrOracleDatabaseLicenseTypeIDNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

// SearchAssociatedLicenseTypesInOracleDatabaseAgreements search Oracle/Database agreements
func (ctrl *APIController) SearchAssociatedLicenseTypesInOracleDatabaseAgreements(w http.ResponseWriter, r *http.Request) {
	var err error

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
	error) {

	var err error

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

	if filters.LicensesCountLTE, err = utils.Str2int(urlValues.Get("licenses-count-lte"), -1); err != nil {
		return dto.SearchOracleDatabaseAgreementsFilter{}, err
	}

	if filters.LicensesCountGTE, err = utils.Str2int(urlValues.Get("licenses-count-gte"), -1); err != nil {
		return dto.SearchOracleDatabaseAgreementsFilter{}, err
	}

	if filters.UsersCountLTE, err = utils.Str2int(urlValues.Get("users-count-lte"), -1); err != nil {
		return dto.SearchOracleDatabaseAgreementsFilter{}, err
	}

	if filters.UsersCountGTE, err = utils.Str2int(urlValues.Get("users-count-gte"), -1); err != nil {
		return dto.SearchOracleDatabaseAgreementsFilter{}, err
	}

	if filters.AvailableCountLTE, err = utils.Str2int(urlValues.Get("available-count-lte"), -1); err != nil {
		return dto.SearchOracleDatabaseAgreementsFilter{}, err
	}

	if filters.AvailableCountGTE, err = utils.Str2int(urlValues.Get("available-count-gte"), -1); err != nil {
		return dto.SearchOracleDatabaseAgreementsFilter{}, err
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

	if err = ctrl.Service.AddHostToAssociatedLicenseType(id, string(raw)); err == utils.ErrOracleDatabaseAgreementNotFound ||
		err == utils.ErrNotInClusterHostNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
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
	var id primitive.ObjectID
	var hostname string

	if id, err = primitive.ObjectIDFromHex(mux.Vars(r)["id"]); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}
	hostname = mux.Vars(r)["hostname"]

	if err = ctrl.Service.RemoveHostFromAssociatedLicenseType(id, hostname); err == utils.ErrOracleDatabaseAgreementNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
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
	var id primitive.ObjectID

	if id, err = primitive.ObjectIDFromHex(mux.Vars(r)["id"]); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	if err = ctrl.Service.DeleteAssociatedLicenseTypeFromOracleDatabaseAgreement(id); err == utils.ErrOracleDatabaseAgreementNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
