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

// Package service is a package that provides methods for manipulating host informations
package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ercole-io/ercole/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

func TestHostdatas(t *testing.T) {
	files, err := filepath.Glob("test_dataservice_hostdata_v1_*.json")
	require.NoError(t, err)

	for _, f := range files {
		t.Run(f, func(t *testing.T) {
			file, err := os.Open(f)
			require.NoError(t, err)
			defer func() {
				err = file.Close()
				require.NoError(t, err)
			}()

			raw, err := ioutil.ReadAll(file)
			require.NoError(t, err)

			//Validate the data
			//TODO Refactor this in model.FrontendHostdataSchemaValidator
			documentLoader := gojsonschema.NewBytesLoader(raw)
			schemaLoader := gojsonschema.NewStringLoader(model.FrontendHostdataSchemaValidator)

			result, err := gojsonschema.Validate(schemaLoader, documentLoader)
			require.NoError(t, err)

			assert.True(t, result.Valid())

			if !result.Valid() {
				fmt.Printf("The input hostdata for file [%v] is not valid:\n", file.Name())
				for _, desc := range result.Errors() {
					fmt.Printf("- %s\n", desc)
				}

				var errorMsg strings.Builder
				for _, err := range result.Errors() {

					value := fmt.Sprintf("%v", err.Value())
					if len(value) > 80 {
						value = value[:78] + ".."
					}
					errorMsg.WriteString(fmt.Sprintf("- %s. Value: [%v]\n", err, value))
				}

				fmt.Println(errorMsg.String())
			}
		})
	}
}
