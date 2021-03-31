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
	"fmt"
	"net/http"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/gddo/httputil"
	"github.com/gorilla/mux"
)

// SearchClusters search clusters data using the filters in the request
func (ctrl *APIController) SearchClusters(w http.ResponseWriter, r *http.Request) {
	choiche := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choiche {
	case "application/json":
		ctrl.SearchClustersJSON(w, r)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchClustersXLSX(w, r)
	}
}

// SearchClustersJSON search clusters data using the filters in the request returning it in JSON format
func (ctrl *APIController) SearchClustersJSON(w http.ResponseWriter, r *http.Request) {
	var full bool
	var search string
	var sortBy string
	var sortDesc bool
	var pageNumber int
	var pageSize int
	var location string
	var environment string
	var olderThan time.Time

	var err error
	//parse the query params
	if full, err = utils.Str2bool(r.URL.Query().Get("full"), false); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
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
	clusters, err := ctrl.Service.SearchClusters(full, search, sortBy, sortDesc, pageNumber, pageSize, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, clusters)
	} else {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, clusters[0])
	}
}

// SearchClustersXLSX search clusters data using the filters in the request returning it in XLSX format
func (ctrl *APIController) SearchClustersXLSX(w http.ResponseWriter, r *http.Request) {
	var search string
	var sortBy string
	var sortDesc bool
	var location string
	var environment string
	var olderThan time.Time

	var err error
	//parse the query params
	search = r.URL.Query().Get("search")
	sortBy = r.URL.Query().Get("sort-by")
	if sortDesc, err = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); err != nil {
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
	clusters, err := ctrl.Service.SearchClusters(false, search, sortBy, sortDesc, -1, -1, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Open the sheet
	sheets, err := excelize.OpenFile(ctrl.Config.ResourceFilePath + "/templates/template_clusters.xlsx")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, utils.NewError(err, "READ_TEMPLATE"))
		return
	}

	//Add the data to the sheet
	for i, val := range clusters {
		sheets.SetCellValue("Hypervisor", fmt.Sprintf("A%d", i+2), val["name"])
		sheets.SetCellValue("Hypervisor", fmt.Sprintf("B%d", i+2), val["type"])
		sheets.SetCellValue("Hypervisor", fmt.Sprintf("C%d", i+2), val["cpu"])
		sheets.SetCellValue("Hypervisor", fmt.Sprintf("D%d", i+2), val["sockets"])
		sheets.SetCellValue("Hypervisor", fmt.Sprintf("E%d", i+2), val["virtualizationNodes"])
	}

	//Write it to the response
	utils.WriteXLSXResponse(w, sheets)
}

// GetCluster get cluster data using the filters in the request
func (ctrl *APIController) GetCluster(w http.ResponseWriter, r *http.Request) {
	clusterName := mux.Vars(r)["name"]

	olderThan, err := utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	choiche := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choiche {
	case "application/json":
		ctrl.GetClusterJSON(w, r, clusterName, olderThan)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.GetClusterXLSX(w, r, clusterName, olderThan)
	}
}

//GetClusterJSON get cluster data using the filters in the request and returns it in JSON format
func (ctrl *APIController) GetClusterJSON(w http.ResponseWriter, r *http.Request, clusterName string, olderThan time.Time) {
	data, err := ctrl.Service.GetCluster(clusterName, olderThan)
	if err == utils.ErrClusterNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, data)
}

//GetClusterXLSX get cluster data using the filters in the request and returns it in XLSX format
func (ctrl *APIController) GetClusterXLSX(w http.ResponseWriter, r *http.Request, clusterName string, olderThan time.Time) {
	xlsx, err := ctrl.Service.GetClusterXLSX(clusterName, olderThan)
	if err == utils.ErrClusterNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, xlsx)
}
