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
	"errors"
	"net/http"
	"time"

	"github.com/amreo/ercole-services/utils"
	"github.com/gorilla/mux"
	"github.com/plandem/xlsx"
)

// SearchHosts search hosts data using the filters in the request
func (ctrl *APIController) SearchHosts(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") == "" || r.Header.Get("Accept") == "application/json" {
		ctrl.SearchHostsJSON(w, r)
	} else if r.Header.Get("Accept") == "application/vnd.oracle.lms+vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		ctrl.SearchHostsLMS(w, r)
	} else {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotAcceptable,
			utils.NewAdvancedErrorPtr(
				errors.New("The mime type in the accept header is not supported"),
				http.StatusText(http.StatusNotAcceptable),
			),
		)
		return
	}
}

// SearchHostsJSON search hosts data using the filters in the request returning it in JSON
func (ctrl *APIController) SearchHostsJSON(w http.ResponseWriter, r *http.Request) {
	var mode string
	var search string
	var sortBy string
	var sortDesc bool
	var pageNumber int
	var pageSize int
	var location string
	var environment string
	var olderThan time.Time

	var err utils.AdvancedErrorInterface
	//parse the query params
	mode = r.URL.Query().Get("mode")
	if mode == "" {
		mode = "full"
	} else if mode != "full" && mode != "summary" && mode != "lms" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

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

	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	hosts, err := ctrl.Service.SearchHosts(mode, search, sortBy, sortDesc, pageNumber, pageSize, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, hosts)
	} else {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, hosts[0])
	}
}

// SearchHostsXLSX search hosts data using the filters in the request returning it in XLSX
func (ctrl *APIController) SearchHostsLMS(w http.ResponseWriter, r *http.Request) {
	var search string
	var sortBy string
	var sortDesc bool
	var location string
	var environment string
	var olderThan time.Time

	var aerr utils.AdvancedErrorInterface
	//parse the query params
	search = r.URL.Query().Get("search")
	sortBy = r.URL.Query().Get("sort-by")
	if sortDesc, aerr = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, aerr)
		return
	}

	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, aerr = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, aerr)
		return
	}

	//get the data
	hosts, aerr := ctrl.Service.SearchHosts("lms", search, sortBy, sortDesc, -1, -1, location, environment, olderThan)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Open the sheet
	sheets, err := xlsx.Open(ctrl.Config.ResourceFilePath + "/templates/template_hosts.xlsm")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, utils.NewAdvancedErrorPtr(err, "READ_TEMPLATE"))
		return
	}

	sheet := sheets.SheetByName("Database_&_EBS")

	i := 0
	//Add the data to the sheet
	for _, val := range hosts {
		if val["Databases"] == "" {
			continue
		}
		sheet.Cell(0, i+3).SetText(val["PhysicalServerName"])
		sheet.Cell(1, i+3).SetText(val["VirtualServerName"])
		sheet.Cell(2, i+3).SetText(val["VirtualizationTechnology"])
		sheet.Cell(3, i+3).SetText(val["DBInstanceName"])
		sheet.Cell(4, i+3).SetText(val["PluggableDatabaseName"])
		sheet.Cell(5, i+3).SetText(val["ConnectString"])
		sheet.Cell(7, i+3).SetText(val["ProductVersion"])
		sheet.Cell(8, i+3).SetText(val["ProductEdition"])
		sheet.Cell(9, i+3).SetText(val["Environment"])
		sheet.Cell(10, i+3).SetText(val["Features"])
		sheet.Cell(11, i+3).SetText(val["RacNodeNames"])
		sheet.Cell(12, i+3).SetText(val["ProcessorModel"])
		sheet.Cell(13, i+3).SetInt(int(val["Processors"].(float64)))
		sheet.Cell(14, i+3).SetInt(int(val["CoresPerProcessor"].(float64)))
		sheet.Cell(15, i+3).SetInt(int(val["PhysicalCores"].(float64)))
		sheet.Cell(16, i+3).SetInt(int(val["ThreadsPerCore"].(int32)))
		sheet.Cell(17, i+3).SetText(val["ProcessorSpeed"])
		sheet.Cell(18, i+3).SetText(val["ServerPurchaseDate"])
		sheet.Cell(19, i+3).SetText(val["OperatingSystem"])
		sheet.Cell(20, i+3).SetText(val["Notes"])
		i++
	}

	//Write it to the response
	utils.WriteXLSXResponse(w, sheets)
}

// GetHost return all'informations about the host requested in the id path variable
func (ctrl *APIController) GetHost(w http.ResponseWriter, r *http.Request) {
	var olderThan time.Time
	var err utils.AdvancedErrorInterface

	hostname := mux.Vars(r)["hostname"]

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	host, err := ctrl.Service.GetHost(hostname, olderThan)
	if err == utils.AerrHostNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, host)
}
