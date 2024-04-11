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

package service

import (
	"github.com/ercole-io/ercole/v2/model"
)

func (as *APIService) GetOraclePsqlMigrabilities(hostname, dbname string) ([]model.PgsqlMigrability, error) {
	return as.Database.FindPsqlMigrabilities(hostname, dbname)
}

func (as *APIService) GetOraclePsqlMigrabilitiesSemaphore(hostname, dbname string) (string, error) {
	psqlMigrabilities, err := as.Database.FindPsqlMigrabilities(hostname, dbname)
	if err != nil {
		return "", err
	}

	color := ""

	for _, migrability := range psqlMigrabilities {
		if migrability.Metric != nil && *migrability.Metric == "PLSQL LINES" {
			switch {
			case migrability.Count < 1000:
				color = "green"
			case migrability.Count >= 1000 && migrability.Count <= 10000:
				color = "yellow"
			case migrability.Count > 10000:
				color = "red"
			}
		}
	}

	return color, nil
}

func (as *APIService) GetOraclePdbPsqlMigrabilities(hostname, dbname, pdbname string) ([]model.PgsqlMigrability, error) {
	return as.Database.FindPdbPsqlMigrabilities(hostname, dbname, pdbname)
}

func (as *APIService) GetOraclePdbPsqlMigrabilitiesSemaphore(hostname, dbname, pdbname string) (string, error) {
	psqlMigrabilities, err := as.Database.FindPdbPsqlMigrabilities(hostname, dbname, pdbname)
	if err != nil {
		return "", err
	}

	color := ""

	for _, migrability := range psqlMigrabilities {
		if migrability.Metric != nil && *migrability.Metric == "PLSQL LINES" {
			switch {
			case migrability.Count < 1000:
				color = "green"
			case migrability.Count >= 1000 && migrability.Count <= 10000:
				color = "yellow"
			case migrability.Count > 10000:
				color = "red"
			}
		}
	}

	return color, nil
}
