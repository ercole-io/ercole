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

package schema

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed hostdata.json
var hostdataSchema string

//go:embed mysql.json
var mysqlSchema string

var schema *gojsonschema.Schema

func ValidateHostdata(raw []byte) error {
	if schema == nil {
		sl := gojsonschema.NewSchemaLoader()
		mysql := gojsonschema.NewStringLoader(mysqlSchema)
		if err := sl.AddSchemas(mysql); err != nil {
			fmt.Println("Should never happen: wrong schema")
			panic(err)
		}

		hostdata := gojsonschema.NewStringLoader(hostdataSchema)

		var err error
		schema, err = sl.Compile(hostdata)
		if err != nil {
			fmt.Println("Should never happen: wrong schema")
			panic(err)
		}
	}

	documentLoader := gojsonschema.NewBytesLoader(raw)
	result, err := schema.Validate(documentLoader)
	syntaxErr := &json.SyntaxError{}
	if errors.As(err, &syntaxErr) {
		return fmt.Errorf("%w: %s", utils.ErrInvalidHostdata, err)
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
			errorMsg.WriteString(fmt.Sprintf("- %s. Value: [%v]\n", err, value))
		}

		return fmt.Errorf("%w: %s", utils.ErrInvalidHostdata, errorMsg.String())
	}

	return nil
}
