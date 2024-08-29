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
	"errors"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils/exutils"
)

func (as *APIService) GetOracleDatabasePoliciesAuditFlag(hostname, dbname string) (map[string][]string, error) {
	exist, err := as.Database.DbExist(hostname, dbname)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("no document found")
	}

	policiesAudit, err := as.Database.FindOracleDatabasePoliciesAudit(hostname, dbname)
	if err != nil {
		return nil, err
	}

	return policiesAudit.Response(as.Config.APIService.OracleDatabasePoliciesAudit, policiesAudit.List), err
}

func (as *APIService) ListOracleDatabasePoliciesAudit() ([]dto.OraclePoliciesAuditListResponse, error) {
	return as.Database.ListOracleDatabasePoliciesAudit()
}

func (as *APIService) CreateOraclePoliciesAuditXlsx(dbs []dto.OraclePoliciesAuditListResponse, pdbs []dto.OraclePdbPoliciesAuditListResponse) (*excelize.File, error) {
	sheet := "policies audit"
	headers := []string{
		"Hostname",
		"Db Name",
		"Pdb Name",
		"Flag",
		"Policies Audit Retrieved From DB",
		"Policies Audit Configured on Ercole",
		"Matched",
		"Not Matched",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range dbs {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.Hostname)
		sheets.SetCellValue(sheet, nextAxis(), val.DbName)
		sheets.SetCellValue(sheet, nextAxis(), "")
		sheets.SetCellValue(sheet, nextAxis(), val.Flag)
		sheets.SetCellValue(sheet, nextAxis(), strings.Join(val.PoliciesAudit, ","))
		sheets.SetCellValue(sheet, nextAxis(), strings.Join(val.PoliciesAuditConfigured, ","))
		sheets.SetCellValue(sheet, nextAxis(), strings.Join(val.Matched, ","))
		sheets.SetCellValue(sheet, nextAxis(), strings.Join(val.NotMatched, ","))
	}

	for _, val := range pdbs {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.Hostname)
		sheets.SetCellValue(sheet, nextAxis(), val.DbName)
		sheets.SetCellValue(sheet, nextAxis(), val.PdbName)
		sheets.SetCellValue(sheet, nextAxis(), val.Flag)
		sheets.SetCellValue(sheet, nextAxis(), strings.Join(val.PoliciesAudit, ","))
		sheets.SetCellValue(sheet, nextAxis(), strings.Join(val.PoliciesAuditConfigured, ","))
		sheets.SetCellValue(sheet, nextAxis(), strings.Join(val.Matched, ","))
		sheets.SetCellValue(sheet, nextAxis(), strings.Join(val.NotMatched, ","))
	}

	return sheets, err
}
