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

package utils

import (
	"encoding/json"
	"net/http"

	"github.com/plandem/xlsx"
	"github.com/sirupsen/logrus"
)

// ErrorResponseFE is a struct that contains informations about a error
type ErrorResponseFE struct {
	// Error contains the (generic) class of the error
	Error string
	// ErrorDescription contains detailed informations about the error
	ErrorDescription string
	// File contains the filename of the source code where the error was detected
	SourceFilename string
	// LineNumber contains the number of the line where the error was detected
	LineNumber int
}

// WriteAndLogError write the error to the w with the statusCode as statusCode and log the error to the stdout
func WriteAndLogError(log *logrus.Logger, w http.ResponseWriter, statusCode int, err AdvancedErrorInterface) {
	//Build the response
	resp := ErrorResponseFE{
		Error:            err.ErrorClass(),
		ErrorDescription: err.Error(),
		LineNumber:       err.LineNumber(),
		SourceFilename:   err.SourceFilename(),
	}
	//Log the error
	LogErr(log, err)
	//Write the response
	WriteJSONResponse(w, statusCode, resp)
}

// WriteJSONResponse write the statuscode and the response to w
func WriteJSONResponse(w http.ResponseWriter, statusCode int, resp interface{}) {
	//Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

// WriteJSONResponse write the statuscode and the response to w
func WriteXLSXResponse(w http.ResponseWriter, resp *xlsx.Spreadsheet) {
	//Write the response
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.WriteHeader(http.StatusOK)
	resp.SaveAs(w)
}
