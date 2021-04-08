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
	"fmt"
	"runtime"
	"strings"
)

// AdvancedError is a struct that contains informations and class about a error
type AdvancedError struct {
	// Err contains the base error of the AdvancedError
	Err error `json:"err"`
	// Message contains the class of the error
	Message string `json:"class"`
	// File contains the filename of the source code where the error was detected
	Source string `json:"source"`
	// LineNumber contains the number of the line where the error was detected
	Line int `json:"line"`
}

// Error return the representation string of the error
func (ae *AdvancedError) Error() string {
	return fmt.Sprintf("%s:%d %s: '%s'\n", ae.Source, ae.Line, ae.Message, ae.Err.Error())
}

func (ae *AdvancedError) Unwrap() error {
	return ae.Err
}

// NewError return a pointer to a new AdvancedError using the err as base error and class as class name
func NewError(err error, message ...string) *AdvancedError {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return &AdvancedError{Err: err, Message: strings.Join(message, " "), Source: file, Line: line}
	} else {
		return &AdvancedError{Err: err, Message: strings.Join(message, " "), Source: "????", Line: -1}
	}
}
