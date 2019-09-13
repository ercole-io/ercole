package utils

import (
	"log"
	"runtime"
)

// AdvancedErrorInterface is a extension of the error interface with other informations
type AdvancedErrorInterface interface {
	error
	// ErrorClass returns the class of the error
	ErrorClass() string
	// SourceFilename returns the filename of the source code where the error was detected
	SourceFilename() string
	// LineNumber returns the number of the line where the error was detected
	LineNumber() int
}

// AdvancedError is a struct that contains informations and class about a error
type AdvancedError struct {
	// Err contains the base error of the AdvancedError
	Err error
	// Class contains the class of the error
	Class string
	// File contains the filename of the source code where the error was detected
	Source string
	// LineNumber contains the number of the line where the error was detected
	Line int
}

//
func (ae *AdvancedError) Error() string {
	return ae.Err.Error()
}

// ErrorClass return the class of the error
func (ae *AdvancedError) ErrorClass() string {
	return ae.Class
}

// SourceFilename return the source filename of the error
func (ae *AdvancedError) SourceFilename() string {
	return ae.Source
}

// LineNumber return the line number of the error
func (ae *AdvancedError) LineNumber() int {
	return ae.Line
}

// NewAdvancedError return a new AdvancedError using the err as base error and class as class name
func NewAdvancedError(err error, class string) AdvancedError {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return AdvancedError{Err: err, Class: class, Source: file, Line: line}
	} else {
		return AdvancedError{Err: err, Class: class, Source: "????", Line: -1}
	}
}

// NewAdvancedErrorPtr return a pointer to a new AdvancedError using the err as base error and class as class name
func NewAdvancedErrorPtr(err error, class string) *AdvancedError {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return &AdvancedError{Err: err, Class: class, Source: file, Line: line}
	} else {
		return &AdvancedError{Err: err, Class: class, Source: "????", Line: -1}
	}
}

// LogErr log the errorr to the stdout
func LogErr(err AdvancedErrorInterface) {
	//Log the error
	log.Printf("%s:%d %s: '%s'\n", err.SourceFilename(), err.LineNumber(), err.Error(), err.ErrorClass())
}
