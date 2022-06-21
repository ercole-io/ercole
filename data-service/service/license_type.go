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

package service

import (
	"encoding/json"
	"fmt"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/schema"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/sanitizer"
)

func (hds *HostDataService) InsertOracleLicenseTypes(licenseTypes []model.OracleDatabaseLicenseType) error {
	for _, lt := range licenseTypes {
		if err := hds.Database.InsertOracleLicenseType(lt); err != nil {
			return err
		}
	}

	return nil
}

func (hds *HostDataService) SanitizeLicenseTypes(raw []byte) ([]model.OracleDatabaseLicenseType, error) {
	var m []map[string]interface{}

	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, utils.ErrInvalidLicenseType
	}

	sanitizer := sanitizer.NewSanitizer(hds.Log)

	sanitizedInt, err := sanitizer.Sanitize(m)
	if err != nil {
		return nil, fmt.Errorf("Unable to sanitize: %w", err)
	}

	if raw, err = json.Marshal(sanitizedInt); err != nil {
		return nil, fmt.Errorf("Unable to marshal: %w", err)
	}

	if validationErr := schema.ValidateLicenseType(raw); validationErr != nil {
		return nil, validationErr
	}

	licenseTypes := make([]model.OracleDatabaseLicenseType, 0)

	err = json.Unmarshal(raw, &licenseTypes)
	if err != nil {
		return nil, err
	}

	return licenseTypes, nil
}
