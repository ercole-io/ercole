// Copyright (c) 2025 Sorint.lab S.p.A.
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package controller

import (
	"errors"
	"net/http"
	"path"
	"strings"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
)

func (ctrl *APIController) CreateScenario(w http.ResponseWriter, r *http.Request) {
	f, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	locations := strings.Split(f.Location, ",")

	req := &dto.CreateScenarioRequest{}

	if err := utils.Decode(r.Body, &req); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	scenario, err := ctrl.Service.CreateScenario(*req, locations, *f)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	response := dto.ToScenarioResponse(*scenario)

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (ctrl *APIController) ListScenario(w http.ResponseWriter, r *http.Request) {
	scenarios, err := ctrl.Service.GetScenarios()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, dto.ToScenariosResponse(scenarios))
}

func (ctrl *APIController) GetScenario(w http.ResponseWriter, r *http.Request) {
	idReq, ok := mux.Vars(r)["id"]
	if !ok || idReq == "" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("missing scenario id"))
		return
	}

	id := utils.Str2oid(idReq)

	scenario, err := ctrl.Service.GetScenario(id)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, dto.ToScenarioResponse(*scenario))
}

func (ctrl *APIController) RemoveScenario(w http.ResponseWriter, r *http.Request) {
	idReq, ok := mux.Vars(r)["id"]
	if !ok || idReq == "" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("missing scenario id"))
		return
	}

	id := utils.Str2oid(idReq)

	err := ctrl.Service.RemoveScenario(id)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ctrl *APIController) GetScenarioLicense(w http.ResponseWriter, r *http.Request) {
	idReq, ok := mux.Vars(r)["id"]
	if !ok || idReq == "" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("missing scenario id"))
		return
	}

	id := utils.Str2oid(idReq)

	scenario, err := ctrl.Service.GetScenario(id)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	base := path.Base(r.URL.Path)
	switch base {
	case "license-compliance":
		utils.WriteJSONResponse(w, http.StatusOK, dto.ToScenarioLicenseComplianceResponse(*scenario))
	case "license-used-database":
		utils.WriteJSONResponse(w, http.StatusOK, dto.ToScenarioLicenseUsedDatabaseResponse(*scenario))
	case "license-used-host":
		utils.WriteJSONResponse(w, http.StatusOK, dto.ToScenarioLicenseUsedHostResponse(*scenario))
	case "license-used-cluster":
		utils.WriteJSONResponse(w, http.StatusOK, dto.ToScenarioLicenseUsedClusterResponse(*scenario))
	case "license-used-cluster-veritas":
		utils.WriteJSONResponse(w, http.StatusOK, dto.ToScenarioLicenseUsedClusterVeritasResponse(*scenario))
	}
}
