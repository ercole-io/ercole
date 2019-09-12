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

func (this *AdvancedError) Error() string {
	return this.Err.Error()
}

// ErrorClass return the class of the error
func (this *AdvancedError) ErrorClass() string {
	return this.Class
}

// SourceFilename return the source filename of the error
func (this *AdvancedError) SourceFilename() string {
	return this.Source
}

// LineNumber return the line number of the error
func (this *AdvancedError) LineNumber() int {
	return this.Line
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

func NewAdvancedErrorPtr(err error, class string) *AdvancedError {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return &AdvancedError{Err: err, Class: class, Source: file, Line: line}
	} else {
		return &AdvancedError{Err: err, Class: class, Source: "????", Line: -1}
	}
}

func LogErr(err AdvancedErrorInterface) {
	//Log the error
	log.Printf("%s:%d %s: '%s'\n", err.SourceFilename(), err.LineNumber(), err.Error(), err.ErrorClass())
}
