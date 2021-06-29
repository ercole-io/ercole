// Copyright (c) 2021 Sorint.lab S.p.A.
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ercole-io/ercole/v2/utils"
)

type AgentError struct {
	Message string `json:"message" bson:"message"`
	Source  string `json:"source" bson:"source"`
}

func NewAgentError(e error) AgentError {
	agentError := AgentError{}

	var ae *utils.AdvancedError
	if !errors.As(e, &ae) {
		agentError.Message = e.Error()
		return agentError
	}

	var b strings.Builder
	if len(ae.Message) > 0 {
		b.WriteString(ae.Message)
	}
	if ae.Err != nil {
		if b.Len() > 0 {
			b.WriteString(": ")
		}
		b.WriteString(ae.Err.Error())
	}

	agentError.Message = b.String()
	agentError.Source = fmt.Sprintf("%s:%d", ae.Source, ae.Line)

	return agentError
}
