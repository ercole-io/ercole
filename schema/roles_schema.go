// Copyright (c) 2022 Sorint.lab S.p.A.
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

package schema

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"

	"github.com/ercole-io/ercole/v2/utils"
)

//go:embed roles.json
var roleSchema string

var schemaR *gojsonschema.Schema

func ValidateRole(raw []byte) error {
	if schemaR == nil {
		if err := loadRoleSchema(); err != nil {
			return err
		}
	}

	documentLoader := gojsonschema.NewBytesLoader(raw)
	result, err := schemaR.Validate(documentLoader)

	syntaxErr := &json.SyntaxError{}
	if errors.As(err, &syntaxErr) {
		return fmt.Errorf("%w: %s", utils.ErrInvalidRole, err)
	} else if err != nil {
		return err
	}

	if !result.Valid() {
		errorMsg := new(strings.Builder)

		for _, err := range result.Errors() {
			value := fmt.Sprintf("%v", err.Value())
			if len(value) > 80 {
				value = value[:78] + ".."
			}

			errorMsg.WriteString(fmt.Sprintf("\t- %s. Value: [%v]\n", err, value))
		}

		return fmt.Errorf("%w:\n%s", utils.ErrInvalidRole, errorMsg.String())
	}

	return nil
}

func loadRoleSchema() error {
	sl := gojsonschema.NewSchemaLoader()

	schemas := []string{roleSchema}
	for i := range schemas {
		jl := gojsonschema.NewStringLoader(schemas[i])
		if err := sl.AddSchemas(jl); err != nil {
			return utils.NewError(err, "Wrong role schema: [%s]", schemas[i])
		}
	}

	role := gojsonschema.NewStringLoader(roleSchema)

	var err error

	schemaR, err = sl.Compile(role)
	if err != nil {
		return utils.NewError(err, "Wrong role schema: can't load or compile it")
	}

	return nil
}
