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
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/utils"
	"github.com/golang/gddo/httputil"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchOracleDatabaseAddms search addms data using the filters in the request
func (ctrl *APIController) SearchOracleDatabaseAddms(w http.ResponseWriter, r *http.Request) {
	choiche := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choiche {
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

	var aerr utils.AdvancedErrorInterface
	//parse the query params
	search = r.URL.Query().Get("search")

	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, aerr = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, aerr)
		return
	}

	//get the data
	addms, aerr := ctrl.Service.SearchOracleDatabaseAddms(search, "Benefit", true, -1, -1, location, environment, olderThan)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Open the sheet
	sheets, err := excelize.OpenFile(ctrl.Config.ResourceFilePath + "/templates/template_addm.xlsx")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, utils.NewAdvancedErrorPtr(err, "READ_TEMPLATE"))
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
	choiche := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choiche {
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
	var pageNumber int
	var pageSize int
	var location string
	var environment string
	var olderThan time.Time

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

	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	//get the data
	segmentAdvisors, err := ctrl.Service.SearchOracleDatabaseSegmentAdvisors(search, sortBy, sortDesc, pageNumber, pageSize, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, segmentAdvisors)
	} else {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, segmentAdvisors[0])
	}
}

// SearchOracleDatabaseSegmentAdvisorsXLSX search segment advisors data using the filters in the request returning it in XLSX format
func (ctrl *APIController) SearchOracleDatabaseSegmentAdvisorsXLSX(w http.ResponseWriter, r *http.Request) {
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
	segmentAdvisors, aerr := ctrl.Service.SearchOracleDatabaseSegmentAdvisors(search, sortBy, sortDesc, -1, -1, location, environment, olderThan)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Open the sheet
	sheets, err := excelize.OpenFile(ctrl.Config.ResourceFilePath + "/templates/template_segment_advisor.xlsx")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, utils.NewAdvancedErrorPtr(err, "READ_TEMPLATE"))
		return
	}

	//Add the data to the sheet
	for i, val := range segmentAdvisors {
		sheets.SetCellValue("Segment_Advisor", fmt.Sprintf("A%d", i+2), val["dbname"])         //Dbname column
		sheets.SetCellValue("Segment_Advisor", fmt.Sprintf("B%d", i+2), val["environment"])    //Environment column
		sheets.SetCellValue("Segment_Advisor", fmt.Sprintf("C%d", i+2), val["hostname"])       //Hostname column
		sheets.SetCellValue("Segment_Advisor", fmt.Sprintf("D%d", i+2), val["partitionName"])  //PartitionName column
		sheets.SetCellValue("Segment_Advisor", fmt.Sprintf("E%d", i+2), val["reclaimable"])    //Reclaimable column
		sheets.SetCellValue("Segment_Advisor", fmt.Sprintf("F%d", i+2), val["recommendation"]) //Recommendation column
		sheets.SetCellValue("Segment_Advisor", fmt.Sprintf("G%d", i+2), val["segmentName"])    //SegmentName column
		sheets.SetCellValue("Segment_Advisor", fmt.Sprintf("H%d", i+2), val["segmentOwner"])   //SegmentOwner column
		sheets.SetCellValue("Segment_Advisor", fmt.Sprintf("I%d", i+2), val["segmentType"])    //SegmentType column
	}

	//Write it to the response
	utils.WriteXLSXResponse(w, sheets)
}

// SearchOracleDatabasePatchAdvisors search patch advisors data using the filters in the request
func (ctrl *APIController) SearchOracleDatabasePatchAdvisors(w http.ResponseWriter, r *http.Request) {
	choiche := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choiche {
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
	if windowTime, err = utils.Str2int(r.URL.Query().Get("window-time"), 6); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}
	status = r.URL.Query().Get("status")
	if status != "" && status != "OK" && status != "KO" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(errors.New("invalid status"), "Invalid  status"))
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
	var search string
	var sortBy string
	var sortDesc bool
	var windowTime int
	var location string
	var environment string
	var olderThan time.Time
	var status string

	var aerr utils.AdvancedErrorInterface
	//parse the query params
	search = r.URL.Query().Get("search")
	sortBy = r.URL.Query().Get("sort-by")
	if sortDesc, aerr = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, aerr)
		return
	}

	if windowTime, aerr = utils.Str2int(r.URL.Query().Get("window-time"), 6); aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, aerr)
		return
	}

	location = r.URL.Query().Get("location")
	environment = r.URL.Query().Get("environment")

	if olderThan, aerr = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, aerr)
		return
	}
	status = r.URL.Query().Get("status")
	if status != "" && status != "OK" && status != "KO" {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(errors.New("invalid status"), "Invalid  status"))
		return
	}

	//get the data
	patchAdvisors, aerr := ctrl.Service.SearchOracleDatabasePatchAdvisors(search, sortBy, sortDesc, -1, -1, ctrl.TimeNow().AddDate(0, -windowTime, 0), location, environment, olderThan, status)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Open the sheet
	sheets, err := excelize.OpenFile(ctrl.Config.ResourceFilePath + "/templates/template_patch_advisor.xlsx")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, utils.NewAdvancedErrorPtr(err, "READ_TEMPLATE"))
		return
	}

	//Add the data to the sheet
	for i, val := range patchAdvisors {
		sheets.SetCellValue("Patch_Advisor", fmt.Sprintf("A%d", i+2), val["description"])                                     //Description column
		sheets.SetCellValue("Patch_Advisor", fmt.Sprintf("B%d", i+2), val["hostname"])                                        //Hostname column
		sheets.SetCellValue("Patch_Advisor", fmt.Sprintf("C%d", i+2), val["dbname"])                                          //Dbname column
		sheets.SetCellValue("Patch_Advisor", fmt.Sprintf("D%d", i+2), val["dbver"])                                           //Dbver column
		sheets.SetCellValue("Patch_Advisor", fmt.Sprintf("E%d", i+2), val["date"].(primitive.DateTime).Time().UTC().String()) //Date column
		sheets.SetCellValue("Patch_Advisor", fmt.Sprintf("F%d", i+2), val["status"])                                          //Status column
	}

	//Write it to the response
	utils.WriteXLSXResponse(w, sheets)
}

// SearchOracleDatabases search databases data using the filters in the request
func (ctrl *APIController) SearchOracleDatabases(w http.ResponseWriter, r *http.Request) {
	choiche := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choiche {
	case "application/json":
		ctrl.SearchOracleDatabasesJSON(w, r)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchOracleDatabasesXLSX(w, r)
	}
}

// SearchOracleDatabasesJSON search databases data using the filters in the request returning it in JSON
func (ctrl *APIController) SearchOracleDatabasesJSON(w http.ResponseWriter, r *http.Request) {
	var full bool
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
	databases, err := ctrl.Service.SearchOracleDatabases(full, search, sortBy, sortDesc, pageNumber, pageSize, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, databases)
	} else {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, databases[0])
	}
}

// SearchOracleDatabasesXLSX search databases data using the filters in the request returning it in XLSX
func (ctrl *APIController) SearchOracleDatabasesXLSX(w http.ResponseWriter, r *http.Request) {
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
	databases, aerr := ctrl.Service.SearchOracleDatabases(false, search, sortBy, sortDesc, -1, -1, location, environment, olderThan)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Open the sheet
	sheets, err := excelize.OpenFile(ctrl.Config.ResourceFilePath + "/templates/template_databases.xlsx")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, utils.NewAdvancedErrorPtr(err, "READ_TEMPLATE"))
		return
	}

	//Add the data to the sheet
	for i, val := range databases {
		sheets.SetCellValue("Databases", fmt.Sprintf("A%d", i+2), val["name"])         //Name column
		sheets.SetCellValue("Databases", fmt.Sprintf("B%d", i+2), val["uniqueName"])   //UniqueName column
		sheets.SetCellValue("Databases", fmt.Sprintf("C%d", i+2), val["version"])      //Version column
		sheets.SetCellValue("Databases", fmt.Sprintf("D%d", i+2), val["hostname"])     //Hostname column
		sheets.SetCellValue("Databases", fmt.Sprintf("E%d", i+2), val["status"])       //Status column
		sheets.SetCellValue("Databases", fmt.Sprintf("F%d", i+2), val["environment"])  //Environment column
		sheets.SetCellValue("Databases", fmt.Sprintf("G%d", i+2), val["location"])     //Location column
		sheets.SetCellValue("Databases", fmt.Sprintf("H%d", i+2), val["charset"])      //Charset column
		sheets.SetCellValue("Databases", fmt.Sprintf("I%d", i+2), val["blockSize"])    //BlockSize column
		sheets.SetCellValue("Databases", fmt.Sprintf("J%d", i+2), val["cpuCount"])     //CPUCount column
		sheets.SetCellValue("Databases", fmt.Sprintf("K%d", i+2), val["work"])         //Work column
		sheets.SetCellValue("Databases", fmt.Sprintf("L%d", i+2), val["memory"])       //Memory column
		sheets.SetCellValue("Databases", fmt.Sprintf("M%d", i+2), val["datafileSize"]) //DatafileSize column
		sheets.SetCellValue("Databases", fmt.Sprintf("N%d", i+2), val["segmentsSize"]) //SegmentsSize column
		sheets.SetCellValue("Databases", fmt.Sprintf("O%d", i+2), val["archivelog"])   //ArchiveLogStatus column
		sheets.SetCellValue("Databases", fmt.Sprintf("P%d", i+2), val["dataguard"])    //Dataguard column
		sheets.SetCellValue("Databases", fmt.Sprintf("Q%d", i+2), val["rac"])          //RAC column
		sheets.SetCellValue("Databases", fmt.Sprintf("R%d", i+2), val["ha"])           //HA column
	}

	//Write it to the response
	utils.WriteXLSXResponse(w, sheets)
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

	var err utils.AdvancedErrorInterface

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

// GetLicense return a certain license asked in the request
func (ctrl *APIController) GetLicense(w http.ResponseWriter, r *http.Request) {
	var err utils.AdvancedErrorInterface
	var olderThan time.Time

	//parse the query params
	if olderThan, err = utils.Str2time(r.URL.Query().Get("older-than"), utils.MAX_TIME); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}
	name := mux.Vars(r)["name"]

	//get the data
	lic, err := ctrl.Service.GetLicense(name, olderThan)
	if err == utils.AerrLicenseNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, err)
		return
	} else if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, lic)
}

// SetLicenseCostPerProcessor set the cost per processor of a certain license
func (ctrl *APIController) SetLicenseCostPerProcessor(w http.ResponseWriter, r *http.Request) {
	if ctrl.Config.APIService.ReadOnly {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusForbidden, utils.NewAdvancedErrorPtr(errors.New("The API is disabled because the service is put in read-only mode"), "FORBIDDEN_REQUEST"))
		return
	}

	//get the data
	name := mux.Vars(r)["name"]

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, utils.NewAdvancedErrorPtr(err, "BAD_REQUEST"))
		return
	}

	costPerProcessor, err := strconv.ParseFloat(string(raw), 32)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, "BAD_REQUEST"))
		return
	}

	//set the value
	aerr := ctrl.Service.SetLicenseCostPerProcessor(name, costPerProcessor)
	if aerr == utils.AerrLicenseNotFound {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotFound, aerr)
		return
	} else if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Write the data
	utils.WriteJSONResponse(w, http.StatusOK, nil)
}
