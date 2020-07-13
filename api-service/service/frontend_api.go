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

// Package service is a package that provides methods for querying data
package service

import (
	"time"

	"github.com/ercole-io/ercole/utils"
)

// GetInfoForFrontendDashboard return all informations needed for the frontend dashboard page
func (as *APIService) GetInfoForFrontendDashboard(location string, environment string, olderThan time.Time) (map[string]interface{}, utils.AdvancedErrorInterface) {
	var err utils.AdvancedErrorInterface
	out := map[string]interface{}{}
	technologiesObject := map[string]interface{}{}

	technologiesObject["total"], err = as.GetTotalTechnologiesComplianceStats(location, environment, olderThan)
	if err != nil {
		return nil, err
	}
	technologiesObject["technologies"], err = as.ListTechnologies("", false, location, environment, olderThan)
	if err != nil {
		return nil, err
	}

	out["features"], err = as.GetErcoleFeatures()
	if err != nil {
		return nil, err
	}
	out["technologies"] = technologiesObject

	return out, nil
}
