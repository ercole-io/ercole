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
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/amreo/ercole-services/utils"
	"github.com/golang/gddo/httputil"
	"github.com/gorilla/mux"
	"github.com/plandem/xlsx"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchAddms search addms data using the filters in the request
func (ctrl *APIController) SearchAddms(w http.ResponseWriter, r *http.Request) {
	choiche := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choiche {
	case "application/json":
		ctrl.SearchAddmsJSON(w, r)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchAddmsXLSX(w, r)
	default:
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotAcceptable,
			utils.NewAdvancedErrorPtr(
				errors.New("The mime type in the accept header is not supported"),
				http.StatusText(http.StatusNotAcceptable),
			),
		)
	}
}

// SearchAddmsJSON search addms data using the filters in the request returning it in JSON format
func (ctrl *APIController) SearchAddmsJSON(w http.ResponseWriter, r *http.Request) {
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
	addms, err := ctrl.Service.SearchAddms(search, sortBy, sortDesc, pageNumber, pageSize, location, environment, olderThan)
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

// SearchAddmsJSON search addms data using the filters in the request returning it in JSON format
func (ctrl *APIController) SearchAddmsXLSX(w http.ResponseWriter, r *http.Request) {
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
	addms, aerr := ctrl.Service.SearchAddms(search, sortBy, sortDesc, -1, -1, location, environment, olderThan)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Open the sheet
	sheets, err := xlsx.Open(ctrl.Config.ResourceFilePath + "/templates/template_addm.xlsx")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, utils.NewAdvancedErrorPtr(err, "READ_TEMPLATE"))
		return
	}
	sheet := sheets.SheetByName("Addm")

	//Add the data to the sheet
	for i, val := range addms {
		sheet.Cell(0, i+1).SetText(val["Action"])         //Action column
		sheet.Cell(1, i+1).SetText(val["Benefit"])        //Benefit column
		sheet.Cell(2, i+1).SetText(val["Dbname"])         //Dbname column
		sheet.Cell(3, i+1).SetText(val["Environment"])    //Environment column
		sheet.Cell(4, i+1).SetText(val["Finding"])        //Finding column
		sheet.Cell(5, i+1).SetText(val["Hostname"])       //Hostname column
		sheet.Cell(6, i+1).SetText(val["Recommendation"]) //Recommendation column
	}

	//Write it to the response
	utils.WriteXLSXResponse(w, sheets)
}

// SearchSegmentAdvisors search segment advisors data using the filters in the request
func (ctrl *APIController) SearchSegmentAdvisors(w http.ResponseWriter, r *http.Request) {
	choiche := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choiche {
	case "application/json":
		ctrl.SearchSegmentAdvisorsJSON(w, r)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchSegmentAdvisorsXLSX(w, r)
	default:
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotAcceptable,
			utils.NewAdvancedErrorPtr(
				errors.New("The mime type in the accept header is not supported"),
				http.StatusText(http.StatusNotAcceptable),
			),
		)
	}
}

// SearchSegmentAdvisorsJSON search segment advisors data using the filters in the request returning it in JSON format
func (ctrl *APIController) SearchSegmentAdvisorsJSON(w http.ResponseWriter, r *http.Request) {
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
	segmentAdvisors, err := ctrl.Service.SearchSegmentAdvisors(search, sortBy, sortDesc, pageNumber, pageSize, location, environment, olderThan)
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

// SearchSegmentAdvisorsXLSX search segment advisors data using the filters in the request returning it in XLSX format
func (ctrl *APIController) SearchSegmentAdvisorsXLSX(w http.ResponseWriter, r *http.Request) {
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
	segmentAdvisors, aerr := ctrl.Service.SearchSegmentAdvisors(search, sortBy, sortDesc, -1, -1, location, environment, olderThan)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Open the sheet
	sheets, err := xlsx.Open(ctrl.Config.ResourceFilePath + "/templates/template_segment_advisor.xlsx")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, utils.NewAdvancedErrorPtr(err, "READ_TEMPLATE"))
		return
	}
	sheet := sheets.SheetByName("Segment_Advisor")

	//Add the data to the sheet
	for i, val := range segmentAdvisors {
		sheet.Cell(0, i+1).SetText(val["Dbname"])         //Dbname column
		sheet.Cell(1, i+1).SetText(val["Environment"])    //Environment column
		sheet.Cell(2, i+1).SetText(val["Hostname"])       //Hostname column
		sheet.Cell(3, i+1).SetText(val["PartitionName"])  //PartitionName column
		sheet.Cell(4, i+1).SetText(val["Reclaimable"])    //Reclaimable column
		sheet.Cell(5, i+1).SetText(val["Recommendation"]) //Recommendation column
		sheet.Cell(6, i+1).SetText(val["SegmentName"])    //SegmentName column
		sheet.Cell(7, i+1).SetText(val["SegmentOwner"])   //SegmentOwner column
		sheet.Cell(8, i+1).SetText(val["SegmentType"])    //SegmentType column
	}

	//Write it to the response
	utils.WriteXLSXResponse(w, sheets)
}

// SearchPatchAdvisors search patch advisors data using the filters in the request
func (ctrl *APIController) SearchPatchAdvisors(w http.ResponseWriter, r *http.Request) {
	choiche := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choiche {
	case "application/json":
		ctrl.SearchPatchAdvisorsJSON(w, r)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchPatchAdvisorsXLSX(w, r)
	default:
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotAcceptable,
			utils.NewAdvancedErrorPtr(
				errors.New("The mime type in the accept header is not supported"),
				http.StatusText(http.StatusNotAcceptable),
			),
		)
	}
}

// SearchPatchAdvisorsJSON search patch advisors data using the filters in the request returning it in JSON format
func (ctrl *APIController) SearchPatchAdvisorsJSON(w http.ResponseWriter, r *http.Request) {
	var search string
	var sortBy string
	var sortDesc bool
	var pageNumber int
	var pageSize int
	var windowTime int
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

	if windowTime, err = utils.Str2int(r.URL.Query().Get("window-time"), 6); err != nil {
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
	patchAdvisors, err := ctrl.Service.SearchPatchAdvisors(search, sortBy, sortDesc, pageNumber, pageSize, time.Now().AddDate(0, -windowTime, 0), location, environment, olderThan)
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

// SearchPatchAdvisorsXLSX search patch advisors data using the filters in the request returning it in XLSX format
func (ctrl *APIController) SearchPatchAdvisorsXLSX(w http.ResponseWriter, r *http.Request) {
	var search string
	var sortBy string
	var sortDesc bool
	var windowTime int
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

	//get the data
	patchAdvisors, aerr := ctrl.Service.SearchPatchAdvisors(search, sortBy, sortDesc, -1, -1, time.Now().AddDate(0, -windowTime, 0), location, environment, olderThan)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Open the sheet
	sheets, err := xlsx.Open(ctrl.Config.ResourceFilePath + "/templates/template_patch_advisor.xlsx")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, utils.NewAdvancedErrorPtr(err, "READ_TEMPLATE"))
		return
	}
	sheet := sheets.SheetByName("Patch_Advisor")

	//Add the data to the sheet
	for i, val := range patchAdvisors {
		sheet.Cell(0, i+1).SetText(val["Description"])                               //Description column
		sheet.Cell(1, i+1).SetText(val["Hostname"])                                  //Hostname column
		sheet.Cell(2, i+1).SetText(val["Dbname"])                                    //Dbname column
		sheet.Cell(3, i+1).SetText(val["Dbver"])                                     //Dbver column
		sheet.Cell(4, i+1).SetText(val["Date"].(primitive.DateTime).Time().String()) //Date column
		sheet.Cell(5, i+1).SetText(val["Status"])                                    //Status column
	}

	//Write it to the response
	utils.WriteXLSXResponse(w, sheets)
}

// SearchDatabases search databases data using the filters in the request
func (ctrl *APIController) SearchDatabases(w http.ResponseWriter, r *http.Request) {
	choiche := httputil.NegotiateContentType(r, []string{"application/json", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}, "application/json")

	switch choiche {
	case "application/json":
		ctrl.SearchDatabasesJSON(w, r)
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ctrl.SearchDatabasesXLSX(w, r)
	default:
		utils.WriteAndLogError(ctrl.Log, w, http.StatusNotAcceptable,
			utils.NewAdvancedErrorPtr(
				errors.New("The mime type in the accept header is not supported"),
				http.StatusText(http.StatusNotAcceptable),
			),
		)
	}
}

// SearchDatabasesJSON search databases data using the filters in the request returning it in JSON
func (ctrl *APIController) SearchDatabasesJSON(w http.ResponseWriter, r *http.Request) {
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
	databases, err := ctrl.Service.SearchDatabases(full, search, sortBy, sortDesc, pageNumber, pageSize, location, environment, olderThan)
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

// SearchDatabasesXLSX search databases data using the filters in the request returning it in XLSX
func (ctrl *APIController) SearchDatabasesXLSX(w http.ResponseWriter, r *http.Request) {
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
	databases, aerr := ctrl.Service.SearchDatabases(false, search, sortBy, sortDesc, -1, -1, location, environment, olderThan)
	if aerr != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, aerr)
		return
	}

	//Open the sheet
	sheets, err := xlsx.Open(ctrl.Config.ResourceFilePath + "/templates/template_databases.xlsx")
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, utils.NewAdvancedErrorPtr(err, "READ_TEMPLATE"))
		return
	}
	sheet := sheets.SheetByName("Databases")

	//Add the data to the sheet
	for i, val := range databases {
		sheet.Cell(0, i+1).SetText(val["Name"])                     //Name column
		sheet.Cell(1, i+1).SetText(val["UniqueName"])               //UniqueName column
		sheet.Cell(2, i+1).SetText(val["Version"])                  //Version column
		sheet.Cell(3, i+1).SetText(val["Hostname"])                 //Hostname column
		sheet.Cell(4, i+1).SetText(val["Status"])                   //Status column
		sheet.Cell(5, i+1).SetText(val["Environment"])              //Environment column
		sheet.Cell(6, i+1).SetText(val["Location"])                 //Location column
		sheet.Cell(7, i+1).SetText(val["Charset"])                  //Charset column
		sheet.Cell(8, i+1).SetText(val["BlockSize"])                //BlockSize column
		sheet.Cell(9, i+1).SetText(val["CPUCount"])                 //CPUCount column
		sheet.Cell(10, i+1).SetText(val["Work"])                    //Work column
		sheet.Cell(11, i+1).SetFloat(val["Memory"].(float64))       //Memory column
		sheet.Cell(12, i+1).SetText(val["DatafileSize"])            //DatafileSize column
		sheet.Cell(13, i+1).SetText(val["SegmentsSize"])            //SegmentsSize column
		sheet.Cell(14, i+1).SetBool(val["ArchiveLogStatus"].(bool)) //ArchiveLogStatus column
		sheet.Cell(15, i+1).SetBool(val["Dataguard"].(bool))        //Dataguard column
		sheet.Cell(16, i+1).SetBool(val["RAC"].(bool))              //RAC column
		sheet.Cell(17, i+1).SetBool(val["HA"].(bool))               //HA column
	}

	//Write it to the response
	utils.WriteXLSXResponse(w, sheets)
}

// ListLicenses list licenses using the filters in the request
func (ctrl *APIController) ListLicenses(w http.ResponseWriter, r *http.Request) {
	var full bool
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
	licenses, err := ctrl.Service.ListLicenses(full, sortBy, sortDesc, pageNumber, pageSize, location, environment, olderThan)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	if pageNumber == -1 || pageSize == -1 {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, licenses)
	} else {
		//Write the data
		utils.WriteJSONResponse(w, http.StatusOK, licenses[0])
	}
}

// SetLicenseCount set the count of a certain license
func (ctrl *APIController) SetLicenseCount(w http.ResponseWriter, r *http.Request) {
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

	count, err := strconv.Atoi(string(raw))
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, "BAD_REQUEST"))
		return
	}

	//set the value
	aerr := ctrl.Service.SetLicenseCount(name, count)
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

// GetDefaultDatabasesTags return the default list of database tags from configuration
func (ctrl *APIController) GetDefaultDatabasesTags(w http.ResponseWriter, r *http.Request) {
	tags, err := ctrl.Service.GetDefaultDatabasesTags()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, tags)
}
