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
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
)

const (
	ORACLE    = "oracle"
	SQLSERVER = "sqlserver"
	MYSQL     = "mysql"
)

func (ctrl *APIController) ImportContractFromCSV(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	defer file.Close()

	databaseType := mux.Vars(r)["databaseType"]
	if databaseType != ORACLE && databaseType != SQLSERVER && databaseType != MYSQL {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("invalid database type in param"))
		return
	}

	c := make(chan error)

	reader := csv.NewReader(file)
	go func(reader *csv.Reader) {
		switch databaseType {
		case "oracle":
			c <- ctrl.Service.ImportOracleDatabaseContracts(reader)
		case "sqlserver":
			c <- ctrl.Service.ImportSQLServerDatabaseContracts(reader)
		case "mysql":
			c <- ctrl.Service.ImportMySQLDatabaseContracts(reader)
		}
	}(reader)

	if err := <-c; err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ctrl *APIController) GetContractSampleCSV(w http.ResponseWriter, r *http.Request) {
	databaseType := mux.Vars(r)["databaseType"]
	if databaseType != ORACLE && databaseType != SQLSERVER && databaseType != MYSQL {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, errors.New("invalid database type in param"))
		return
	}

	res, err := ctrl.Service.GetLicenseContractSample(databaseType)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=sample_%s_contracts.csv", databaseType))
	w.Header().Set("Content-Type", "text/csv")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(res); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}
}
