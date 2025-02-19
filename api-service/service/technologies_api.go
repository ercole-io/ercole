// Copyright (c) 2023 Sorint.lab S.p.A.
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
	"strings"
	"time"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

// ListManagedTechnologies returns the list of Technologies with some stats
func (as *APIService) ListManagedTechnologies(sortBy string, sortDesc bool, location string, environment string, olderThan time.Time) ([]model.TechnologyStatus, error) {
	hostsCountByTechnology, err := as.Database.GetHostsCountUsingTechnologies(location, environment, olderThan)
	if err != nil {
		return nil, err
	}

	statuses := make([]model.TechnologyStatus, 0)

	oracleStatus, err := createOracleTechnologyStatus(as, hostsCountByTechnology[model.TechnologyOracleDatabase])
	if err != nil {
		return nil, err
	}

	licCompliances, err := as.GetDatabaseLicensesCompliance(strings.Split(location, ","))
	if err != nil {
		return nil, err
	}

	totCompliance := float64(0.0)

	for _, lc := range licCompliances {
		totCompliance += lc.Compliance
	}

	if len(licCompliances) > 0 {
		newAvgPercentage := totCompliance / float64(len(licCompliances))

		oracleStatus.NewAvgPercentage = newAvgPercentage
	}

	statuses = append(statuses, *oracleStatus)

	mysqlStatus, err := createMySqlTechnologyStatus(as, hostsCountByTechnology[model.TechnologyOracleMySQL])
	if err != nil {
		return nil, err
	}

	statuses = append(statuses, *mysqlStatus)

	sqlServerStatus, err := createSqlServerTechnologyStatus(as, hostsCountByTechnology[model.TechnologyMicrosoftSQLServer])
	if err != nil {
		return nil, err
	}

	statuses = append(statuses, *sqlServerStatus)

	postgreSQLStatus := model.TechnologyStatus{
		Product:            model.TechnologyPostgreSQLPostgreSQL,
		ConsumedByHosts:    0,
		CoveredByContracts: 0,
		TotalCost:          0,
		PaidCost:           0,
		Compliance:         1,
		UnpaidDues:         0,
		HostsCount:         int(hostsCountByTechnology[model.TechnologyPostgreSQLPostgreSQL]),
	}

	statuses = append(statuses, postgreSQLStatus)

	mongoDBStatus := model.TechnologyStatus{
		Product:            model.TechnologyMongoDBMongoDB,
		ConsumedByHosts:    0,
		CoveredByContracts: 0,
		TotalCost:          0,
		PaidCost:           0,
		Compliance:         1,
		UnpaidDues:         0,
		HostsCount:         int(hostsCountByTechnology[model.TechnologyMongoDBMongoDB]),
	}

	statuses = append(statuses, mongoDBStatus)

	mariaDBStatus := model.TechnologyStatus{
		Product:            model.TechnologyMariaDBFoundationMariaDB,
		ConsumedByHosts:    0,
		CoveredByContracts: 0,
		TotalCost:          0,
		PaidCost:           0,
		Compliance:         1,
		UnpaidDues:         0,
		HostsCount:         0,
	}

	statuses = append(statuses, mariaDBStatus)

	return statuses, nil
}

func createOracleTechnologyStatus(as *APIService, hostsCount float64) (*model.TechnologyStatus, error) {
	filter := dto.NewGetOracleDatabaseContractsFilter()
	
	contracts, err := as.Database.ListOracleDatabaseContracts(filter)
	if err != nil {
		return nil, err
	}

	usages, err := as.getLicensesUsage(filter.Locations)
	if err != nil {
		return nil, err
	}

	err2 := as.assignOracleDatabaseContractsToHosts(contracts, usages)
	if err2 != nil {
		return nil, utils.NewError(err2, "DB ERROR")
	}

	status := model.TechnologyStatus{
		Product:    model.TechnologyOracleDatabase,
		HostsCount: int(hostsCount),
	}

	for _, usage := range usages {
		status.ConsumedByHosts += usage.OriginalCount
		status.CoveredByContracts += (usage.OriginalCount - usage.LicenseCount)
	}

	if status.ConsumedByHosts == 0 {
		status.Compliance = 1
	} else {
		status.Compliance = status.CoveredByContracts / status.ConsumedByHosts
	}

	return &status, nil
}

func createSqlServerTechnologyStatus(as *APIService, hostsCount float64) (*model.TechnologyStatus, error) {
	licensesCompliance, err := as.GetSqlServerDatabaseLicensesCompliance([]string{})
	if err != nil {
		return nil, err
	}

	status := model.TechnologyStatus{
		Product:    model.TechnologyMicrosoftSQLServer,
		HostsCount: int(hostsCount),
	}

	for _, licenseCompliance := range licensesCompliance {
		status.ConsumedByHosts += licenseCompliance.Consumed
		status.CoveredByContracts += licenseCompliance.Covered
	}

	if status.ConsumedByHosts == 0 {
		status.Compliance = 1
	} else {
		status.Compliance = status.CoveredByContracts / status.ConsumedByHosts
	}

	return &status, nil
}

func createMySqlTechnologyStatus(as *APIService, hostsCount float64) (*model.TechnologyStatus, error) {
	licensesCompliance, err := as.GetMySQLDatabaseLicensesCompliance([]string{})
	if err != nil {
		return nil, err
	}

	status := model.TechnologyStatus{
		Product:    model.TechnologyOracleMySQL,
		HostsCount: int(hostsCount),
	}

	for _, licenseCompliance := range licensesCompliance {
		status.ConsumedByHosts += licenseCompliance.Consumed
		status.CoveredByContracts += licenseCompliance.Covered
	}

	if status.ConsumedByHosts == 0 {
		status.Compliance = 1
	} else {
		status.Compliance = status.CoveredByContracts / status.ConsumedByHosts
	}

	return &status, nil
}
