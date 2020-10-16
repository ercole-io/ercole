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

package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
)

// LoadOracleDatabaseAgreementParts loads the list of Oracle/Database agreement parts and store it to as.OracleDatabaseAgreementParts.
func (as *APIService) LoadOracleDatabaseAgreementParts() {
	filename := "oracle_database_agreement_parts.json"
	path := filepath.Join(as.Config.ResourceFilePath, filename)

	bytes, err := ioutil.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		as.Log.Warnf("No %s file exists in resources (%s), no agreement parts set\n",
			filename, as.Config.ResourceFilePath)
		as.OracleDatabaseAgreementParts = make([]model.OracleDatabaseAgreementPart, 0)

		return
	} else if err != nil {
		as.Log.Errorf("Unable to read %s: %v\n", path, err)

		return
	}

	err = json.Unmarshal(bytes, &as.OracleDatabaseAgreementParts)
	if err != nil {
		as.Log.Errorf("Unable to unmarshal %s: %v\n", path, err)
		return
	}
}

// GetOracleDatabaseAgreementPartsList return the list of Oracle/Database agreement parts
func (as *APIService) GetOracleDatabaseAgreementPartsList() ([]model.OracleDatabaseAgreementPart, utils.AdvancedErrorInterface) {
	return as.OracleDatabaseAgreementParts, nil
}
