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
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/golang/gddo/httputil"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

// SearchOracleDatabaseAddms search addms data using the filters in the request
func (ctrl *APIController) SearchOracleDatabaseAddms(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		ctrl.SearchOracleDatabaseAddmsJSON(w, r)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchOracleDatabaseAddmsXLSX(w, r)
	}
}

// SearchOracleDatabaseAddmsJSON search addms data using the filters in the request returning it in JSON format
func (ctrl *APIController) SearchOracleDatabaseAddmsJSON(w http.ResponseWriter, r *http.Request) {
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
	addms, err := ctrl.Service.SearchOracleDatabaseAddms(search, sortBy, sortDesc, pageNumber, pageSize, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, addms)
	} else {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, addms[0])
	}
}

// SearchOracleDatabaseAddmsXLSX search addms data using the filters in the request returning it in XLSX format
func (ctrl *APIController) SearchOracleDatabaseAddmsXLSX(w http.ResponseWriter, r *http.Request) {
	var search string
	var location string
	var environment string
	var olderThan time.Time

	var err error
	//parse the query params
	search = r.URL.Query().Get("search")

	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	addms, err := ctrl.Service.SearchOracleDatabaseAddms(search, "Benefit", true, -1, -1, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Open the sheet
	sheets, err := excelize.OpenFile(ctrl.Config.ResourceFilePath + "/templates/template_addm.xlsx")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, utils.NewError(err, "READ_TEMPLATE"))
		return
	}

	//Add the data to the sheet
	for i, val := range addms {
		sheets.SetCellValue("Addm", fmt.Sprintf("A%d", i+2), val["action"])         //Action column
		sheets.SetCellValue("Addm", fmt.Sprintf("B%d", i+2), val["benefit"])        //Benefit column
		sheets.SetCellValue("Addm", fmt.Sprintf("C%d", i+2), val["dbname"])         //Dbname column
		sheets.SetCellValue("Addm", fmt.Sprintf("D%d", i+2), val["environment"])    //Environment column
		sheets.SetCellValue("Addm", fmt.Sprintf("E%d", i+2), val["finding"])        //Finding column
		sheets.SetCellValue("Addm", fmt.Sprintf("F%d", i+2), val["hostname"])       //Hostname column
		sheets.SetCellValue("Addm", fmt.Sprintf("G%d", i+2), val["recommendation"]) //Recommendation column
	}

	//Write it to the response
	utils.WriteXLSXResponse(w, sheets)
}

// SearchOracleDatabaseSegmentAdvisors search segment advisors data using the filters in the request
func (ctrl *APIController) SearchOracleDatabaseSegmentAdvisors(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		ctrl.SearchOracleDatabaseSegmentAdvisorsJSON(w, r)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchOracleDatabaseSegmentAdvisorsXLSX(w, r)
	}
}

// SearchOracleDatabaseSegmentAdvisorsJSON search segment advisors data using the filters in the request returning it in JSON format
func (ctrl *APIController) SearchOracleDatabaseSegmentAdvisorsJSON(w http.ResponseWriter, r *http.Request) {
	var search string
	var sortBy string
	var sortDesc bool
	var location string
	var environment string
	var olderThan time.Time

	var err error

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

	segmentAdvisors, err := ctrl.Service.SearchOracleDatabaseSegmentAdvisors(search, sortBy, sortDesc, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	resp := map[string]interface{}{
		"segmentAdvisors": segmentAdvisors,
	}
	utils.WriteJSONResponse(w, http.StatusOK, resp)
}

// SearchOracleDatabaseSegmentAdvisorsXLSX search segment advisors data using the filters in the request returning it in XLSX format
func (ctrl *APIController) SearchOracleDatabaseSegmentAdvisorsXLSX(w http.ResponseWriter, r *http.Request) {
	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	xlsx, err := ctrl.Service.SearchOracleDatabaseSegmentAdvisorsAsXLSX(*filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, xlsx)
}

// SearchOracleDatabasePatchAdvisors search patch advisors data using the filters in the request
func (ctrl *APIController) SearchOracleDatabasePatchAdvisors(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choice {
	case "application/json":
		ctrl.SearchOracleDatabasePatchAdvisorsJSON(w, r)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchOracleDatabasePatchAdvisorsXLSX(w, r)
	}
}

// SearchOracleDatabasePatchAdvisorsJSON search patch advisors data using the filters in the request returning it in JSON format
func (ctrl *APIController) SearchOracleDatabasePatchAdvisorsJSON(w http.ResponseWriter, r *http.Request) {
	var search string
	var sortBy string
	var sortDesc bool
	var pageNumber int
	var pageSize int
	var windowTime int
	var location string
	var environment string
	var olderThan time.Time
	var status string
	var err error
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
	if windowTime, err = utils.Str2int(r.URL.Query().Get("window-time"), 6); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}
	status = r.URL.Query().Get("status")
	if status != "" && status != "OK" && status != "KO" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewError(errors.New("invalid status"), "Invalid  status"))
		return
	}

	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	patchAdvisors, err := ctrl.Service.SearchOracleDatabasePatchAdvisors(search, sortBy, sortDesc, pageNumber, pageSize, ctrl.TimeNow().AddDate(0, -windowTime, 0), location, environment, olderThan, status)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, patchAdvisors)
	} else {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, patchAdvisors[0])
	}
}

// SearchOracleDatabasePatchAdvisorsXLSX search patch advisors data using the filters in the request returning it in XLSX format
func (ctrl *APIController) SearchOracleDatabasePatchAdvisorsXLSX(w http.ResponseWriter, r *http.Request) {
	var windowTime int
	filter, err := dto.GetGlobalFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	if windowTime, err = utils.Str2int(r.URL.Query().Get("window-time"), 6); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	xlsx, err := ctrl.Service.SearchOracleDatabasePatchAdvisorsAsXLSX(ctrl.TimeNow().AddDate(0, -windowTime, 0), *filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, xlsx)
}

// SearchOracleDatabases search databases data using the filters in the request
func (ctrl *APIController) SearchOracleDatabases(w http.ResponseWriter, r *http.Request) {
	choice := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	filter, err := dto.GetSearchOracleDatabasesFilter(r)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	switch choice {
	case "application/json":
		ctrl.SearchOracleDatabasesJSON(w, r, *filter)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchOracleDatabasesXLSX(w, r, *filter)
	}
}

// SearchOracleDatabasesJSON search databases data using the filters in the request returning it in JSON
func (ctrl *APIController) SearchOracleDatabasesJSON(w http.ResponseWriter, r *http.Request, filter dto.SearchOracleDatabasesFilter) {
	databases, err := ctrl.Service.SearchOracleDatabases(filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if filter.PageNumber == -1 || filter.PageSize == -1 {
		utils.WriteJSONResponse(w, http.StatusOK, databases)
	} else {
		utils.WriteJSONResponse(w, http.StatusOK, databases[0])
	}
}

// SearchOracleDatabasesXLSX search databases data using the filters in the request returning it in XLSX
func (ctrl *APIController) SearchOracleDatabasesXLSX(w http.ResponseWriter, r *http.Request, filter dto.SearchOracleDatabasesFilter) {
	file, err := ctrl.Service.SearchOracleDatabasesAsXLSX(filter)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteXLSXResponse(w, file)
}

// SearchOracleDatabaseUsedLicenses search licenses consumed by the hosts using the filters in the request
func (ctrl *APIController) SearchOracleDatabaseUsedLicenses(w http.ResponseWriter, r *http.Request) {
	var sortBy string
	var sortDesc bool
	var pageNumber int
	var pageSize int
	var location string
	var environment string
	var olderThan time.Time

	var err error

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

	response, err := ctrl.Service.SearchOracleDatabaseUsedLicenses(sortBy, sortDesc, pageNumber, pageSize, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		utils.WriteJSONResponse(w, http.StatusOK, response.Content)
	} else {
		utils.WriteJSONResponse(w, http.StatusOK, response)
	}
}
