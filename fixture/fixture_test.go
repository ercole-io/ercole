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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ercole-io/ercole/v2/schema"
	"github.com/stretchr/testify/require"
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

			err = schema.ValidateHostdata(raw)
			require.NoError(t, err)
		})
	}
}
