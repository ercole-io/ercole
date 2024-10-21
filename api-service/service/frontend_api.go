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
	"fmt"
	"time"

	"github.com/ercole-io/ercole/v2/api-service/dto"
)

// GetInfoForFrontendDashboard return all informations needed for the frontend dashboard page
func (as *APIService) GetInfoForFrontendDashboard(location string, environment string, olderThan time.Time) (map[string]interface{}, error) {
	var err error

	out := map[string]interface{}{}

	technologiesObject := map[string]interface{}{}

	technologiesObject["total"], err = as.GetTotalTechnologiesComplianceStats(location, environment, olderThan)
	if err != nil {
		return nil, err
	}

	technologiesObject["technologies"], err = as.ListManagedTechnologies("", false, location, environment, olderThan)
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

func (as *APIService) GetComplianceStats() (*dto.ComplianceStats, error) {
	oracleStats, err := as.oracleStats()
	if err != nil {
		return nil, err
	}

	sqlServerStats, err := as.sqlServerStats()
	if err != nil {
		return nil, err
	}

	mysqlStats, err := as.mysqlStats()
	if err != nil {
		return nil, err
	}

	postgresqlStats, err := as.postgreSqlStats()
	if err != nil {
		return nil, err
	}

	mongodbStats, err := as.mongoDbStats()
	if err != nil {
		return nil, err
	}

	hostsCount, err := as.Database.CountAllHost()
	if err != nil {
		return nil, err
	}

	avg := (oracleStats.CompliancePercentageVal + mysqlStats.CompliancePercentageVal + sqlServerStats.CompliancePercentageVal +
		postgresqlStats.CompliancePercentageVal + mongodbStats.CompliancePercentageVal) / 5

	totStats := dto.Stats{
		Count:                   int(hostsCount),
		HostCount:               int(hostsCount),
		CompliancePercentageVal: avg,
		CompliancePercentageStr: fmt.Sprintf("%.2f%%", avg),
	}

	res := dto.ComplianceStats{
		Ercole:     &totStats,
		Oracle:     oracleStats,
		MySql:      mysqlStats,
		SqlServer:  sqlServerStats,
		PostgreSql: postgresqlStats,
		MongoDb:    mongodbStats,
	}

	return &res, nil
}

func (as *APIService) oracleStats() (*dto.Stats, error) {
	count, err := as.Database.CountOracleInstance()
	if err != nil {
		return nil, err
	}

	hostCount, err := as.Database.CountOracleHosts()
	if err != nil {
		return nil, err
	}

	compliances, err := as.GetOracleDatabaseLicensesCompliance()
	if err != nil {
		return nil, err
	}

	compliancePercentage := float64(0.0)

	if len(compliances) > 0 {
		totCompliance := float64(0.0)

		for _, v := range compliances {
			totCompliance += v.Compliance
		}

		compliancePercentage = (totCompliance * 100) / float64(len(compliances))
	}

	return &dto.Stats{
		Count:                   int(count),
		HostCount:               int(hostCount),
		CompliancePercentageVal: compliancePercentage,
		CompliancePercentageStr: fmt.Sprintf("%.2f%%", compliancePercentage),
	}, nil
}

func (as *APIService) mysqlStats() (*dto.Stats, error) {
	count, err := as.Database.CountMySqlInstance()
	if err != nil {
		return nil, err
	}

	hostCount, err := as.Database.CountMySqlHosts()
	if err != nil {
		return nil, err
	}

	compliances, err := as.GetMySQLDatabaseLicensesCompliance()
	if err != nil {
		return nil, err
	}

	compliancePercentage := float64(0.0)

	if len(compliances) > 0 {
		totCompliance := float64(0.0)

		for _, v := range compliances {
			totCompliance += v.Compliance
		}

		compliancePercentage = (totCompliance * 100) / float64(len(compliances))
	}

	return &dto.Stats{
		Count:                   int(count),
		HostCount:               int(hostCount),
		CompliancePercentageVal: compliancePercentage,
		CompliancePercentageStr: fmt.Sprintf("%.2f%%", compliancePercentage),
	}, nil
}

func (as *APIService) sqlServerStats() (*dto.Stats, error) {
	count, err := as.Database.CountSqlServerlInstance()
	if err != nil {
		return nil, err
	}

	hostCount, err := as.Database.CountSqlServerHosts()
	if err != nil {
		return nil, err
	}

	compliances, err := as.GetSqlServerDatabaseLicensesCompliance()
	if err != nil {
		return nil, err
	}

	compliancePercentage := float64(0.0)

	if len(compliances) > 0 {
		totCompliance := float64(0.0)

		for _, v := range compliances {
			totCompliance += v.Compliance
		}

		compliancePercentage = (totCompliance * 100) / float64(len(compliances))
	}

	return &dto.Stats{
		Count:                   int(count),
		HostCount:               int(hostCount),
		CompliancePercentageVal: compliancePercentage,
		CompliancePercentageStr: fmt.Sprintf("%.2f%%", compliancePercentage),
	}, nil
}

func (as *APIService) postgreSqlStats() (*dto.Stats, error) {
	count, err := as.Database.CountPostgreSqlInstance()
	if err != nil {
		return nil, err
	}

	hostCount, err := as.Database.CountPostgreSqlHosts()
	if err != nil {
		return nil, err
	}

	return &dto.Stats{
		Count:                   int(count),
		HostCount:               int(hostCount),
		CompliancePercentageStr: "100%",
		CompliancePercentageVal: 100,
	}, nil
}

func (as *APIService) mongoDbStats() (*dto.Stats, error) {
	count, err := as.Database.CountMongoDbInstance()
	if err != nil {
		return nil, err
	}

	hostCount, err := as.Database.CountMongoDbHosts()
	if err != nil {
		return nil, err
	}

	return &dto.Stats{
		Count:                   int(count),
		HostCount:               int(hostCount),
		CompliancePercentageStr: "100%",
		CompliancePercentageVal: 100,
	}, nil
}
