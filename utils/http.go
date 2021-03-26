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

package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"runtime"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

// ErrorResponseFE is a struct that contains informations about a error
type ErrorResponseFE struct {
	// Error contains detailed informations about the error
	Error string `json:"error"`
	// ErrorClass contains the (generic) class of the error
	ErrorClass string `json:"errorClass,omitempty"`
	// File contains the filename of the source code where the error was detected
	SourceFilename string `json:"sourceFilename,omitempty"`
	// LineNumber contains the number of the line where the error was detected
	LineNumber int `json:"lineNumber,omitempty"`
}

// WriteAndLogError write the error to the w with the statusCode as statusCode and log the error to the stdout
func WriteAndLogError(log *logrus.Logger, w http.ResponseWriter, statusCode int, err error) {
	var resp ErrorResponseFE

	var aerr *AdvancedError
	if errors.As(err, &aerr) {
		resp = ErrorResponseFE{
			Error:          aerr.Err.Error(),
			ErrorClass:     aerr.Class,
			LineNumber:     aerr.Line,
			SourceFilename: aerr.Source,
		}
	} else if _, file, line, ok := runtime.Caller(1); ok {
		resp = ErrorResponseFE{
			Error:          err.Error(),
			ErrorClass:     http.StatusText(statusCode),
			SourceFilename: file,
			LineNumber:     line,
		}
	} else {
		resp = ErrorResponseFE{
			Error:          err.Error(),
			ErrorClass:     "",
			SourceFilename: "",
			LineNumber:     0,
		}
	}

	if statusCode >= 500 {
		log.Error(err)
	}

	WriteJSONResponse(w, statusCode, resp)
}

// WriteJSONResponse write the statuscode and the response to w
func WriteJSONResponse(w http.ResponseWriter, statusCode int, resp interface{}) {
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Errorf("Error while encoding response: %#v", resp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(jsonResp)
	if err != nil {
		logrus.Errorf("Error while writing response: %#v", string(jsonResp))
		return
	}
}

// WriteExtJSONResponse write the statuscode and the response to w
func WriteExtJSONResponse(log *logrus.Logger, w http.ResponseWriter, statusCode int, resp interface{}) {
	raw, err := bson.MarshalExtJSON(resp, true, false)
	if err != nil {
		WriteAndLogError(log, w, http.StatusInternalServerError, NewError(err, "MARSHAL_EXT_JSON"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(raw)
}

// WriteXLSXResponse for .xlsx fils
func WriteXLSXResponse(w http.ResponseWriter, resp *excelize.File) {
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.WriteHeader(http.StatusOK)

	resp.Write(w)
}

// WriteXLSMResponse for .xlsm files
func WriteXLSMResponse(w http.ResponseWriter, resp *excelize.File) {
	w.Header().Set("Content-Type", "application/vnd.ms-excel.sheet.macroEnabled.12")
	w.WriteHeader(http.StatusOK)

	resp.Write(w)
}

func Decode(body io.ReadCloser, i interface{}) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&i); err != nil {
		return err
	}

	return nil
}
