// Copyright (c) 2023 Sorint.lab S.p.A.
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
	"net/http"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
)

func (ctrl *APIController) ListExadata(w http.ResponseWriter, r *http.Request) {
	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	res, err := ctrl.Service.ListExadataInstances(*filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	dtos, err := dto.ToOracleExadataInstances(res)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, dtos)
}

func (ctrl *APIController) UpdateExadataVmClusterName(w http.ResponseWriter, r *http.Request) {
	rackID := mux.Vars(r)["rackID"]
	hostID := mux.Vars(r)["hostID"]
	vmname := mux.Vars(r)["name"]

	type vm struct {
		ClusterName string
	}

	c := vm{}

	if err := utils.Decode(r.Body, &c); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if err := ctrl.Service.UpdateExadataVmClusterName(rackID, hostID, vmname, c.ClusterName); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ctrl *APIController) UpdateExadataComponentClusterName(w http.ResponseWriter, r *http.Request) {
	rackID := mux.Vars(r)["rackID"]
	hostID := mux.Vars(r)["hostID"]

	type component struct {
		ClusterName string
	}

	c := component{}

	if err := utils.Decode(r.Body, &c); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if err := ctrl.Service.UpdateExadataComponentClusterName(rackID, hostID, c.ClusterName); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ctrl *APIController) UpdateExadataRdma(w http.ResponseWriter, r *http.Request) {
	rackID := mux.Vars(r)["rackID"]

	rdma := model.OracleExadataRdma{}

	if err := utils.Decode(r.Body, &rdma); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if err := ctrl.Service.UpdateExadataRdma(rackID, rdma); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
