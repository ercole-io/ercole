// Copyright (c) 2019 Sorint.lab S.p.A.
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

// Package service is a package that provides methods for querying data
package service

import (
	"encoding/json"
	"time"

	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
)

// ListAssets returns the list of assets with some stats
func (as *APIService) ListAssets(sortBy string, sortDesc bool, location string, environment string, olderThan time.Time) ([]model.AssetStatus, utils.AdvancedErrorInterface) {
	partialList, err := as.Database.GetAssetsUsage(location, environment, olderThan)
	if err != nil {
		return nil, err
	}

	oracleLicenseListRaw, err := as.Database.ListLicenses(false, "", false, -1, -1, location, environment, olderThan)
	if err != nil {
		return nil, err
	}

	finalList := make([]model.AssetStatus, 0)

	//Oracle/Exadata
	if partialList["Oracle/Exadata"] > 0 {
		finalList = append(finalList, model.AssetStatus{
			Name:       "Oracle/Exadata",
			Used:       partialList["Oracle/Exadata"],
			Count:      partialList["Oracle/Exadata"],
			Compliance: true,
			Cost:       0.0,
		})
	}

	//Oracle/Databases
	type License struct {
		Count float32
		Used  float32
	}
	oracleLicenseList := make([]License, 0)
	json.Unmarshal([]byte(utils.ToJSON(oracleLicenseListRaw)), &oracleLicenseList)
	used := float32(0.0)
	holded := float32(0.0)
	for _, lic := range oracleLicenseList {
		used += lic.Used
		if lic.Count > lic.Used {
			holded += lic.Used
		} else {
			holded += lic.Count
		}
	}
	if used > 0 {
		finalList = append(finalList, model.AssetStatus{
			Name:       "Oracle/Database",
			Used:       used,
			Count:      holded,
			Compliance: used <= holded,
			Cost:       0.0,
		})
	}

	return finalList, nil
}
