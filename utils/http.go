package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/amreo/ercole-hostdata-dataservice/utils"
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
func WriteAndLogError(w http.ResponseWriter, statusCode int, err utils.AdvancedError) {
	//Build the response
	resp := ErrorResponseFE{
		Error:            err.ErrorClass(),
		ErrorDescription: err.ErrorClass(),
		LineNumber:       err.LineNumber(),
		SourceFilename:   err.SourceFilename(),
	}
	//Log the error
	log.Printf("%s:%d %s: '%s'", resp.SourceFilename, resp.LineNumber, resp.Error, resp.ErrorDescription)
	//Write the response
	WriteJSONResponse(w, statusCode, resp)
}

// WriteJSONResponse write the statuscode and the response to w
func WriteJSONResponse(w http.ResponseWriter, statusCode int, resp interface{}) {
	//Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
