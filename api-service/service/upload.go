// Copyright (c) 2024 Sorint.lab S.p.A.
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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/schema"
	"github.com/gocarina/gocsv"
)

func (as *APIService) ImportOracleDatabaseContracts(reader *csv.Reader) error {
	contracts := make([]model.OracleDatabaseContract, 0)

	if err := gocsv.UnmarshalCSV(reader, &contracts); err != nil {
		return err
	}

	for _, contract := range contracts {
		raw, err := json.Marshal(contract)
		if err != nil {
			return err
		}

		if err := schema.ValidateOracleContract(raw); err != nil {
			return err
		}

		if len(contract.HostsLiteral) > 0 {
			contract.Hosts = strings.Split(string(contract.HostsLiteral), ",")
		}

		if _, err := as.AddOracleDatabaseContract(contract); err != nil {
			return err
		}
	}

	return nil
}

func (as *APIService) ImportSQLServerDatabaseContracts(reader *csv.Reader) error {
	contracts := make([]model.SqlServerDatabaseContract, 0)

	if err := gocsv.UnmarshalCSV(reader, &contracts); err != nil {
		return err
	}

	for _, contract := range contracts {
		if len(contract.HostsLiteral) > 0 {
			contract.Hosts = strings.Split(string(contract.HostsLiteral), "|||")
		}

		if len(contract.ClusterLiteral) > 0 {
			contract.Clusters = strings.Split(string(contract.ClusterLiteral), "|||")
		}

		if _, err := as.AddSqlServerDatabaseContract(contract); err != nil {
			return err
		}
	}

	return nil
}

func (as *APIService) ImportMySQLDatabaseContracts(reader *csv.Reader) error {
	contracts := make([]model.MySQLContract, 0)

	if err := gocsv.UnmarshalCSV(reader, &contracts); err != nil {
		return err
	}

	for _, contract := range contracts {
		if len(contract.HostsLiteral) > 0 {
			contract.Hosts = strings.Split(string(contract.HostsLiteral), "|||")
		}

		if len(contract.ClusterLiteral) > 0 {
			contract.Clusters = strings.Split(string(contract.ClusterLiteral), "|||")
		}

		if _, err := as.AddMySQLContract(contract); err != nil {
			return err
		}
	}

	return nil
}

func (as *APIService) GetLicenseContractSample(dbtype string) ([]byte, error) {
	switch dbtype {
	case "oracle":
		empData := []model.OracleDatabaseContract{}
		return gocsv.MarshalBytes(empData)
	case "sqlserver":
		empData := []model.SqlServerDatabaseContract{}
		return gocsv.MarshalBytes(empData)
	case "mysql":
		empData := []model.MySQLContract{}
		return gocsv.MarshalBytes(empData)
	default:
		return nil, fmt.Errorf("cannot match database type: %s", dbtype)
	}
}
