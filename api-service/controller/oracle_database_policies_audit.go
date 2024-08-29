// Copyright (c) 2024 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/gddo/httputil"
	"github.com/gorilla/mux"
)

func (ctrl *APIController) GetOraclePoliciesAudit(w http.ResponseWriter, r *http.Request) {
	hostname := mux.Vars(r)["hostname"]
	dbname := mux.Vars(r)["dbname"]

	resp, err := ctrl.Service.GetOracleDatabasePoliciesAuditFlag(hostname, dbname)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, resp)
}

func (ctrl *APIController) ListOraclePoliciesAudit(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		resp, err := ctrl.Service.ListOracleDatabasePoliciesAudit()
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, resp)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		db, err := ctrl.Service.ListOracleDatabasePoliciesAudit()
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, err)
			return
		}

		pdb, err := ctrl.Service.ListOracleDatabasePdbPoliciesAudit()
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, err)
			return
		}

		result, err := ctrl.Service.CreateOraclePoliciesAuditXlsx(db, pdb)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, err)
			return
		}

		utils.WriteXLSXResponse(w, result)
	}
}
