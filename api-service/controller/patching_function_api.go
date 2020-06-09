// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"strconv"

	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/gorilla/mux"
)

// GetPatchingFunction return all'informations about the patching function of the host requested in the hostname path variable
func (ctrl *APIController) GetPatchingFunction(w http.ResponseWriter, r *http.Request) {
	var err utils.AdvancedErrorInterface

	hostname := mux.Vars(r)["hostname"]

	//get the data
	pf, err := ctrl.Service.GetPatchingFunction(hostname)
	if err == utils.AerrHostNotFound || err == utils.AerrPatchingFunctionNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, pf)
}

// SetPatchingFunction set the patching function of a host specified in the hostname path variable to the content of the request body
func (ctrl *APIController) SetPatchingFunction(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}
	if !ctrl.Config.APIService.EnableInsertingCustomPatchingFunction {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the configuration property EnableInsertingCustomPatchingFunction is false"), "FORBIDDEN_REQUEST"))
		return
	}

	//get the data
	hostname := mux.Vars(r)["hostname"]
	var pf model.PatchingFunction
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&pf); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}
	//set the value
	id, aerr := ctrl.Service.SetPatchingFunction(hostname, pf)
	if aerr == utils.AerrHostNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, aerr)
		return
	} else if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, id)
}

// DeletePatchingFunction delete the patching function of a host specified in the hostname path variable
func (ctrl *APIController) DeletePatchingFunction(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}
	if !ctrl.Config.APIService.EnableInsertingCustomPatchingFunction {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the configuration property EnableInsertingCustomPatchingFunction is false"), "FORBIDDEN_REQUEST"))
		return
	}

	//get the data
	hostname := mux.Vars(r)["hostname"]

	//delete the value
	aerr := ctrl.Service.DeletePatchingFunction(hostname)
	if aerr == utils.AerrHostNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, aerr)
		return
	} else if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

// AddTagToDatabase add a tag to the database if it hasn't the tag
func (ctrl *APIController) AddTagToDatabase(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	//get the data
	hostname := mux.Vars(r)["hostname"]
	dbname := mux.Vars(r)["dbname"]

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewAdvancedErrorPtr(err, "BAD_REQUEST"))
		return
	}

	//set the value
	aerr := ctrl.Service.AddTagToDatabase(hostname, dbname, string(raw))
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

// DeleteTagOfDatabase remove a certain tag from a database if it has the tag
func (ctrl *APIController) DeleteTagOfDatabase(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	//get the data
	hostname := mux.Vars(r)["hostname"]
	dbname := mux.Vars(r)["dbname"]
	tagname := mux.Vars(r)["tagname"]

	//set the value
	aerr := ctrl.Service.DeleteTagOfDatabase(hostname, dbname, tagname)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

// SetLicenseModifier set the license modifier of specified license/db/host in the request to the value in the body
func (ctrl *APIController) SetLicenseModifier(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	//get the data
	hostname := mux.Vars(r)["hostname"]
	dbname := mux.Vars(r)["dbname"]
	licensename := mux.Vars(r)["licenseName"]

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewAdvancedErrorPtr(err, "BAD_REQUEST"))
		return
	}

	newValue, err := strconv.Atoi(string(raw))
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, "BAD_REQUEST"))
		return
	}

	//set the value
	aerr := ctrl.Service.SetLicenseModifier(hostname, dbname, licensename, newValue)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

// DeleteLicenseModifier delete the license modifier of specified license/db/host in the request
func (ctrl *APIController) DeleteLicenseModifier(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	//get the data
	hostname := mux.Vars(r)["hostname"]
	dbname := mux.Vars(r)["dbname"]
	licensename := mux.Vars(r)["licenseName"]

	//set the value
	aerr := ctrl.Service.DeleteLicenseModifier(hostname, dbname, licensename)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, nil)
}

// SearchLicenseModifiers search a license modifier using the filters in the request
func (ctrl *APIController) SearchLicenseModifiers(w http.ResponseWriter, r *http.Request) {
	var search string
	var sortBy string
	var sortDesc bool
	var pageNumber int
	var pageSize int

	var err utils.AdvancedErrorInterface
	//parse the query params
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

	//get the data
	data, err := ctrl.Service.SearchLicenseModifiers(search, sortBy, sortDesc, pageNumber, pageSize)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, data)
	} else {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, data[0])
	}
}
