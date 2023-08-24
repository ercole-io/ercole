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

//go:embed user.json
var userSchema string

func ValidateUser(raw []byte) error {
	schemaLoader, err := loadUserSchema()
	if err != nil {
		return nil
	}

	documentLoader := gojsonschema.NewBytesLoader(raw)
	result, err := schemaLoader.Validate(documentLoader)

	syntaxErr := &json.SyntaxError{}
	if errors.As(err, &syntaxErr) {
		return fmt.Errorf("%w: %s", utils.ErrInvalidUser, err)
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

		return fmt.Errorf("%w:\n%s", utils.ErrInvalidUser, errorMsg.String())
	}

	return nil
}

func loadUserSchema() (*gojsonschema.Schema, error) {
	sl := gojsonschema.NewSchemaLoader()

	schemas := []string{userSchema}
	for i := range schemas {
		jl := gojsonschema.NewStringLoader(schemas[i])
		if err := sl.AddSchemas(jl); err != nil {
			return nil, utils.NewError(err, "Wrong user schema: [%s]", schemas[i])
		}
	}

	user := gojsonschema.NewStringLoader(userSchema)

	var err error

	uSchema, err := sl.Compile(user)
	if err != nil {
		return nil, utils.NewError(err, "Wrong user schema: can't load or compile it")
	}

	return uSchema, nil
}
