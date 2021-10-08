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

	"github.com/golang/gddo/httputil"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (ctrl *APIController) AddOracleDatabaseAgreement(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewError(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var req model.OracleDatabaseAgreement

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest,
			utils.NewError(err, http.StatusText(http.StatusBadRequest)))
		return
	}

	if req.ID != primitive.NilObjectID {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest,
			utils.NewError(errors.New("ID must be empty to add a new AssociatedLicenseType"), http.StatusText(http.StatusBadRequest)))
		return
	}

	if err := req.Check(); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	agr, err := ctrl.Service.AddOracleDatabaseAgreement(req)
	if errors.Is(err, utils.ErrOracleDatabaseAgreementNotFound) ||
		errors.Is(err, utils.ErrOracleDatabaseLicenseTypeIDNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, agr)
}

func (ctrl *APIController) UpdateOracleDatabaseAgreement(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewError(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var req model.OracleDatabaseAgreement

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest,
			utils.NewError(err, http.StatusText(http.StatusBadRequest)))
		return
	}

	if err := req.Check(); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if req.Unlimited && !req.Basket {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewErrorf("Agreement is unlimited so it must be even Basket"))
		return
	}

	agr, err := ctrl.Service.UpdateOracleDatabaseAgreement(req)
	if errors.Is(err, utils.ErrOracleDatabaseAgreementNotFound) ||
		errors.Is(err, utils.ErrOracleDatabaseLicenseTypeIDNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, agr)
}

func (ctrl *APIController) GetOracleDatabaseAgreements(w http.ResponseWriter, r *http.Request) {
	searchOracleDatabaseAgreementsFilters, err := parseGetOracleDatabaseAgreementsFilters(r.URL.Query())
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		ctrl.GetOracleDatabaseAgreementsJSON(w, r, searchOracleDatabaseAgreementsFilters)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.GetOracleDatabaseAgreementsXLSX(w, r, searchOracleDatabaseAgreementsFilters)
	}
}

func (ctrl *APIController) GetOracleDatabaseAgreementsJSON(w http.ResponseWriter, r *http.Request, filters dto.GetOracleDatabaseAgreementsFilter) {
	agreements, err := ctrl.Service.GetOracleDatabaseAgreements(filters)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"agreements": agreements,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (ctrl *APIController) GetOracleDatabaseAgreementsXLSX(w http.ResponseWriter, r *http.Request, filters dto.GetOracleDatabaseAgreementsFilter) {
	xlsx, err := ctrl.Service.GetOracleDatabaseAgreementsAsXLSX(filters)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, xlsx)
}

func parseGetOracleDatabaseAgreementsFilters(urlValues url.Values) (dto.GetOracleDatabaseAgreementsFilter,
	error) {

	var err error

	filters := dto.NewGetOracleDatabaseAgreementsFilter()

	filters.AgreementID = urlValues.Get("agreement-id")
	filters.LicenseTypeID = urlValues.Get("license-type-id")
	filters.ItemDescription = urlValues.Get("item-description")
	filters.CSI = urlValues.Get("csi")
	filters.Metric = urlValues.Get("metrics")
	filters.ReferenceNumber = urlValues.Get("reference-number")

	filters.Unlimited = urlValues.Get("unlimited")
	if filters.Unlimited != "true" && filters.Unlimited != "false" && filters.Unlimited != "" {
		return dto.GetOracleDatabaseAgreementsFilter{},
			utils.NewError(errors.New("Invalid value for unlimited"), http.StatusText(http.StatusUnprocessableEntity))
	}

	filters.Basket = urlValues.Get("basket")
	if filters.Basket != "true" && filters.Basket != "false" && filters.Basket != "" {
		return dto.GetOracleDatabaseAgreementsFilter{},
			utils.NewError(errors.New("Invalid value for basket"), http.StatusText(http.StatusUnprocessableEntity))
	}

	if filters.LicensesPerCoreLTE, err = utils.Str2int(urlValues.Get("licenses-per-core-lte"), -1); err != nil {
		return dto.GetOracleDatabaseAgreementsFilter{}, err
	}

	if filters.LicensesPerCoreGTE, err = utils.Str2int(urlValues.Get("licenses-per-core-gte"), -1); err != nil {
		return dto.GetOracleDatabaseAgreementsFilter{}, err
	}

	if filters.LicensesPerUserLTE, err = utils.Str2int(urlValues.Get("licenses-per-user-lte"), -1); err != nil {
		return dto.GetOracleDatabaseAgreementsFilter{}, err
	}

	if filters.LicensesPerUserGTE, err = utils.Str2int(urlValues.Get("licenses-per-user-gte"), -1); err != nil {
		return dto.GetOracleDatabaseAgreementsFilter{}, err
	}

	if filters.AvailableLicensesPerCoreLTE, err = utils.Str2int(urlValues.Get("available-licenses-per-core-lte"), -1); err != nil {
		return dto.GetOracleDatabaseAgreementsFilter{}, err
	}

	if filters.AvailableLicensesPerCoreGTE, err = utils.Str2int(urlValues.Get("available-licenses-per-core-gte"), -1); err != nil {
		return dto.GetOracleDatabaseAgreementsFilter{}, err
	}

	if filters.AvailableLicensesPerUserLTE, err = utils.Str2int(urlValues.Get("available-licenses-per-user-lte"), -1); err != nil {
		return dto.GetOracleDatabaseAgreementsFilter{}, err
	}

	if filters.AvailableLicensesPerUserGTE, err = utils.Str2int(urlValues.Get("available-licenses-per-user-gte"), -1); err != nil {
		return dto.GetOracleDatabaseAgreementsFilter{}, err
	}

	return filters, nil
}

func (ctrl *APIController) DeleteOracleDatabaseAgreement(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewError(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var err error
	var id primitive.ObjectID

	if id, err = primitive.ObjectIDFromHex(mux.Vars(r)["id"]); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewError(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	if err = ctrl.Service.DeleteOracleDatabaseAgreement(id); errors.Is(err, utils.ErrOracleDatabaseAgreementNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

func (ctrl *APIController) AddHostToOracleDatabaseAgreement(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewError(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var err error
	var id primitive.ObjectID

	if id, err = primitive.ObjectIDFromHex(mux.Vars(r)["id"]); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity,
			utils.NewError(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewError(err, http.StatusText(http.StatusBadRequest)))
		return
	}
	defer r.Body.Close()

	if err = ctrl.Service.AddHostToOracleDatabaseAgreement(id, string(raw)); errors.Is(err, utils.ErrOracleDatabaseAgreementNotFound) ||
		errors.Is(err, utils.ErrNotInClusterHostNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

func (ctrl *APIController) DeleteHostFromOracleDatabaseAgreement(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewError(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	var err error
	var id primitive.ObjectID
	var hostname string

	if id, err = primitive.ObjectIDFromHex(mux.Vars(r)["id"]); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewError(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}
	hostname = mux.Vars(r)["hostname"]

	if err = ctrl.Service.DeleteHostFromOracleDatabaseAgreement(id, hostname); errors.Is(err, utils.ErrOracleDatabaseAgreementNotFound) {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
