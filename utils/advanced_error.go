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
	"runtime"

	"github.com/sirupsen/logrus"
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

// Error return the representation string of the error
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

// LogErr log the error to the stdout
func LogErr(log *logrus.Logger, err AdvancedErrorInterface) {
	log.Errorf("%s:%d %s: '%s'\n", err.SourceFilename(), err.LineNumber(), err.ErrorClass(), err.Error())
}
