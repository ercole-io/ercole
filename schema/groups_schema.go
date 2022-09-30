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

//go:embed groups.json
var groupSchema string

var schemaG *gojsonschema.Schema

func ValidateGroup(raw []byte) error {
	if schemaG == nil {
		if err := loadGroupSchema(); err != nil {
			return err
		}
	}

	documentLoader := gojsonschema.NewBytesLoader(raw)
	result, err := schemaG.Validate(documentLoader)

	syntaxErr := &json.SyntaxError{}
	if errors.As(err, &syntaxErr) {
		return fmt.Errorf("%w: %s", utils.ErrInvalidGroup, err)
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

		return fmt.Errorf("%w:\n%s", utils.ErrInvalidGroup, errorMsg.String())
	}

	return nil
}

func loadGroupSchema() error {
	sl := gojsonschema.NewSchemaLoader()

	schemas := []string{groupSchema}
	for i := range schemas {
		jl := gojsonschema.NewStringLoader(schemas[i])
		if err := sl.AddSchemas(jl); err != nil {
			return utils.NewError(err, "Wrong role schema: [%s]", schemas[i])
		}
	}

	role := gojsonschema.NewStringLoader(groupSchema)

	var err error

	schemaG, err = sl.Compile(role)
	if err != nil {
		return utils.NewError(err, "Wrong group schema: can't load or compile it")
	}

	return nil
}
