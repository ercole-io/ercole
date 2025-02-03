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
	"slices"
	"time"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
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

func (as *APIService) GetComplianceStats(user model.User) (*dto.ComplianceStats, error) {
	locations := model.AllLocations
	if !user.IsAdmin() {
		var err error
		if locations, err = as.ListLocations(user); err != nil {
			return nil, err
		}
	}

	oracleStats, err := as.oracleStats(locations)
	if err != nil {
		return nil, err
	}

	sqlServerStats, err := as.sqlServerStats(locations)
	if err != nil {
		return nil, err
	}

	mysqlStats, err := as.mysqlStats(locations)
	if err != nil {
		return nil, err
	}

	postgresqlStats, err := as.postgreSqlStats(locations)
	if err != nil {
		return nil, err
	}

	mongodbStats, err := as.mongoDbStats(locations)
	if err != nil {
		return nil, err
	}

	mariadbStats, err := as.mariaDbStats(locations)
	if err != nil {
		return nil, err
	}

	hostsCount := oracleStats.HostCount +
		sqlServerStats.HostCount +
		mysqlStats.HostCount +
		postgresqlStats.HostCount +
		mongodbStats.HostCount +
		mariadbStats.HostCount

	instancesCount := oracleStats.Count +
		sqlServerStats.Count +
		mysqlStats.Count +
		postgresqlStats.Count +
		mongodbStats.Count +
		mariadbStats.Count

	var weightedSum, weightedAvg float64

	if hostsCount > 0 {
		weightedSum = (oracleStats.CompliancePercentageVal * float64(oracleStats.HostCount)) +
			(mysqlStats.CompliancePercentageVal * float64(mysqlStats.HostCount)) +
			(sqlServerStats.CompliancePercentageVal * float64(sqlServerStats.HostCount)) +
			(postgresqlStats.CompliancePercentageVal * float64(postgresqlStats.HostCount)) +
			(mongodbStats.CompliancePercentageVal * float64(mongodbStats.HostCount)) +
			(mariadbStats.CompliancePercentageVal * float64(mariadbStats.HostCount))

		weightedAvg = weightedSum / float64(hostsCount)
	}

	if weightedAvg == 0 || hostsCount == 0 {
		weightedAvg = 100
	}

	totStats := dto.Stats{
		Count:                   instancesCount,
		HostCount:               hostsCount,
		CompliancePercentageVal: weightedAvg,
		CompliancePercentageStr: fmt.Sprintf("%.2f%%", weightedAvg),
	}

	res := dto.ComplianceStats{
		Ercole:     &totStats,
		Oracle:     oracleStats,
		MySql:      mysqlStats,
		SqlServer:  sqlServerStats,
		PostgreSql: postgresqlStats,
		MongoDb:    mongodbStats,
		MariaDB:    mariadbStats,
	}

	return &res, nil
}

func (as *APIService) oracleStats(locations []string) (*dto.Stats, error) {
	var count, hostCount int64
	var err error

	if slices.Contains(locations, model.AllLocation) {
		if count, err = as.Database.CountOracleInstance(); err != nil {
			return nil, err
		}
		if hostCount, err = as.Database.CountOracleHosts(); err != nil {
			return nil, err
		}
	} else {
		if count, err = as.Database.CountOracleInstanceByLocations(locations); err != nil {
			return nil, err
		}
		if hostCount, err = as.Database.CountOracleHostsByLocations(locations); err != nil {
			return nil, err
		}
	}

	compliances, err := as.GetOracleDatabaseLicensesCompliance(locations)
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

	if compliancePercentage == 0 || hostCount == 0 {
		compliancePercentage = 100
	}

	return &dto.Stats{
		Count:                   int(count),
		HostCount:               int(hostCount),
		CompliancePercentageVal: compliancePercentage,
		CompliancePercentageStr: fmt.Sprintf("%.2f%%", compliancePercentage),
	}, nil
}

func (as *APIService) mysqlStats(locations []string) (*dto.Stats, error) {
	var count, hostCount int64
	var err error

	if slices.Contains(locations, model.AllLocation) {
		if count, err = as.Database.CountMySqlInstance(); err != nil {
			return nil, err
		}
		if hostCount, err = as.Database.CountMySqlHosts(); err != nil {
			return nil, err
		}
	} else {
		if count, err = as.Database.CountMySqlInstanceByLocations(locations); err != nil {
			return nil, err
		}
		if hostCount, err = as.Database.CountMySqlHostsByLocations(locations); err != nil {
			return nil, err
		}
	}

	compliances, err := as.GetMySQLDatabaseLicensesCompliance(locations)
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

	if compliancePercentage == 0 || hostCount == 0 {
		compliancePercentage = 100
	}

	return &dto.Stats{
		Count:                   int(count),
		HostCount:               int(hostCount),
		CompliancePercentageVal: compliancePercentage,
		CompliancePercentageStr: fmt.Sprintf("%.2f%%", compliancePercentage),
	}, nil
}

func (as *APIService) sqlServerStats(locations []string) (*dto.Stats, error) {
	var count, hostCount int64
	var err error

	if slices.Contains(locations, model.AllLocation) {
		if count, err = as.Database.CountSqlServerlInstance(); err != nil {
			return nil, err
		}
		if hostCount, err = as.Database.CountSqlServerHosts(); err != nil {
			return nil, err
		}
	} else {
		if count, err = as.Database.CountSqlServerlInstanceByLocations(locations); err != nil {
			return nil, err
		}
		if hostCount, err = as.Database.CountSqlServerHostsByLocations(locations); err != nil {
			return nil, err
		}
	}

	compliances, err := as.GetSqlServerDatabaseLicensesCompliance(locations)
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

	if compliancePercentage == 0 || hostCount == 0 {
		compliancePercentage = 100
	}

	return &dto.Stats{
		Count:                   int(count),
		HostCount:               int(hostCount),
		CompliancePercentageVal: compliancePercentage,
		CompliancePercentageStr: fmt.Sprintf("%.2f%%", compliancePercentage),
	}, nil
}

func (as *APIService) postgreSqlStats(locations []string) (*dto.Stats, error) {
	var count, hostCount int64
	var err error

	if slices.Contains(locations, model.AllLocation) {
		if count, err = as.Database.CountPostgreSqlInstance(); err != nil {
			return nil, err
		}
		if hostCount, err = as.Database.CountPostgreSqlHosts(); err != nil {
			return nil, err
		}
	} else {
		if count, err = as.Database.CountPostgreSqlInstanceByLocations(locations); err != nil {
			return nil, err
		}
		if hostCount, err = as.Database.CountPostgreSqlHostsByLocations(locations); err != nil {
			return nil, err
		}
	}

	return &dto.Stats{
		Count:                   int(count),
		HostCount:               int(hostCount),
		CompliancePercentageStr: "100%",
		CompliancePercentageVal: 100,
	}, nil
}

func (as *APIService) mongoDbStats(locations []string) (*dto.Stats, error) {
	var count, hostCount int64
	var err error

	if slices.Contains(locations, model.AllLocation) {
		if count, err = as.Database.CountMongoDbInstance(); err != nil {
			return nil, err
		}
		if hostCount, err = as.Database.CountMongoDbHosts(); err != nil {
			return nil, err
		}
	} else {
		if count, err = as.Database.CountMongoDbInstanceByLocations(locations); err != nil {
			return nil, err
		}
		if hostCount, err = as.Database.CountMongoDbHostsByLocations(locations); err != nil {
			return nil, err
		}
	}

	return &dto.Stats{
		Count:                   int(count),
		HostCount:               int(hostCount),
		CompliancePercentageStr: "100%",
		CompliancePercentageVal: 100,
	}, nil
}

func (as *APIService) mariaDbStats(locations []string) (*dto.Stats, error) {
	return &dto.Stats{
		Count:                   0,
		HostCount:               0,
		CompliancePercentageStr: "100%",
		CompliancePercentageVal: 100,
	}, nil
}
