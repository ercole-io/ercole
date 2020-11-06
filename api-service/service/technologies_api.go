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

	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
)

// ListManagedTechnologies returns the list of Technologies with some stats
func (as *APIService) ListManagedTechnologies(sortBy string, sortDesc bool, location string, environment string, olderThan time.Time) ([]model.TechnologyStatus, utils.AdvancedErrorInterface) {
	partialList, err := as.Database.GetHostsCountUsingTechnologies(location, environment, olderThan)
	if err != nil {
		return nil, err
	}

	// FIXME Readd correct values
	//oracleLicenseListRaw, err := as.Database.SearchLicenses(location, environment, olderThan)
	//if err != nil {
	//	return nil, err
	//}

	finalList := make([]model.TechnologyStatus, 0)

	//Oracle/Databases
	type License struct {
		Count     float64 `json:"count"`
		Used      float64 `json:"used"`
		PaidCost  float64 `json:"paidCost"`
		TotalCost float64 `json:"totalCost"`
		Unlimited bool    `json:"unlimited"`
	}
	oracleLicenseList := make([]License, 0)
	//json.Unmarshal([]byte(utils.ToJSON(oracleLicenseListRaw)), &oracleLicenseList)
	used := float64(0.0)
	holded := float64(0.0)
	totalCost := float64(0.0)
	paidCost := float64(0.0)
	for _, lic := range oracleLicenseList {
		used += lic.Used
		totalCost += lic.TotalCost
		if lic.Count > lic.Used || lic.Unlimited {
			holded += lic.Used
			paidCost += lic.TotalCost
		} else {
			holded += lic.Count
			paidCost += lic.PaidCost
		}
	}
	finalList = append(finalList, model.TechnologyStatus{
		Product:    model.TechnologyOracleDatabase,
		Used:       used,
		Count:      holded,
		TotalCost:  totalCost,
		PaidCost:   paidCost,
		HostsCount: int(partialList[model.TechnologyOracleDatabase]),
		Compliance: holded / used,
		UnpaidDues: totalCost - paidCost,
	})
	if used == 0 {
		finalList[len(finalList)-1].Compliance = 1
	}

	//MariaDBFoundation/MariaDB
	finalList = append(finalList, model.TechnologyStatus{
		Product:    model.TechnologyMariaDBFoundationMariaDB,
		Used:       0,
		Count:      0,
		TotalCost:  0.0,
		PaidCost:   0.0,
		HostsCount: 0.0,
		Compliance: 1.0,
		UnpaidDues: 0.0,
	})

	//PostgreSQL/PostgreSQL
	finalList = append(finalList, model.TechnologyStatus{
		Product:    model.TechnologyPostgreSQLPostgreSQL,
		Used:       0,
		Count:      0,
		TotalCost:  0.0,
		PaidCost:   0.0,
		HostsCount: 0.0,
		Compliance: 1.0,
		UnpaidDues: 0.0,
	})

	//Oracle/MySQL
	finalList = append(finalList, model.TechnologyStatus{
		Product:    model.TechnologyOracleMySQL,
		Used:       0,
		Count:      0,
		TotalCost:  0.0,
		PaidCost:   0.0,
		HostsCount: 0.0,
		Compliance: 1.0,
		UnpaidDues: 0.0,
	})

	//Microsoft/SQLServer
	finalList = append(finalList, model.TechnologyStatus{
		Product:    model.TechnologyMicrosoftSQLServer,
		Used:       0,
		Count:      0,
		TotalCost:  0.0,
		PaidCost:   0.0,
		HostsCount: 0.0,
		Compliance: 1.0,
		UnpaidDues: 0.0,
	})

	return finalList, nil
}
