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
	"bytes"
	"fmt"
	"sync"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils/exutils"
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

func (as *APIService) ListOracleDatabasePsqlMigrabilities() ([]dto.OracleDatabasePgsqlMigrability, error) {
	return as.Database.ListOracleDatabasePsqlMigrabilities()
}

func (as *APIService) ListOracleDatabasePdbPsqlMigrabilities() ([]dto.OracleDatabasePdbPgsqlMigrability, error) {
	return as.Database.ListOracleDatabasePdbPsqlMigrabilities()
}

func (as *APIService) CreateOraclePsqlMigrabilitiesXlsx(dbs []dto.OracleDatabasePgsqlMigrability, pdbs []dto.OracleDatabasePdbPgsqlMigrability) (*excelize.File, error) {
	var buffer bytes.Buffer

	addToBuffer(&buffer, dbs, pdbs)

	sheet := "Psql Migrabilities"
	headers := []string{
		"Hostname",
		"Db Name",
		"Pdb Name",
		"Flag",
		"Metrics",
		"Object Type",
		"Schema",
		"Count",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	content := buffer.Bytes()

	for i, line := range bytes.Split([]byte(content), []byte("\n")) {
		i += 2

		splitted := bytes.Split(line, []byte("|"))

		if len(splitted) >= 8 {
			sheets.SetCellValue(sheet, fmt.Sprintf("A%d", i), string(splitted[0]))
			sheets.SetCellValue(sheet, fmt.Sprintf("B%d", i), string(splitted[1]))
			sheets.SetCellValue(sheet, fmt.Sprintf("C%d", i), string(splitted[2]))
			sheets.SetCellValue(sheet, fmt.Sprintf("D%d", i), string(splitted[3]))
			sheets.SetCellValue(sheet, fmt.Sprintf("E%d", i), string(splitted[4]))
			sheets.SetCellValue(sheet, fmt.Sprintf("F%d", i), string(splitted[5]))
			sheets.SetCellValue(sheet, fmt.Sprintf("G%d", i), string(splitted[6]))
			sheets.SetCellValue(sheet, fmt.Sprintf("H%d", i), string(splitted[7]))
		}
	}

	return sheets, err
}

func addToBuffer(buffer *bytes.Buffer, data ...interface{}) {
	var wg sync.WaitGroup

	var mu sync.Mutex

	processSlice := func(data interface{}) {
		defer wg.Done()

		var localBuffer bytes.Buffer

		switch data := data.(type) {
		case []dto.OracleDatabasePgsqlMigrability:
			for _, db := range data {
				for _, m := range db.Metrics {
					localBuffer.WriteString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%d\n",
						db.Hostname, db.Dbname, "", db.Flag, m.GetMetric(), m.GetObjectType(), m.GetSchema(), m.Count))
				}
			}

		case []dto.OracleDatabasePdbPgsqlMigrability:
			for _, pdb := range data {
				for _, m := range pdb.Metrics {
					localBuffer.WriteString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%d\n",
						pdb.Hostname, pdb.Dbname, pdb.Pdbname, pdb.Flag, m.GetMetric(), m.GetObjectType(), m.GetSchema(), m.Count))
				}
			}

		default:
			fmt.Printf("type: %T\n", data)
			return
		}

		mu.Lock()
		buffer.Write(localBuffer.Bytes())
		mu.Unlock()
	}

	for _, d := range data {
		wg.Add(1)

		go processSlice(d)
	}

	wg.Wait()
}
