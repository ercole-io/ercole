// Copyright (c) 2022 Sorint.lab S.p.A.
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

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/exutils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (as *APIService) AddSqlServerDatabaseContract(contract model.SqlServerDatabaseContract) (*model.SqlServerDatabaseContract, error) {
	if err := checkHosts(as, contract.Hosts); err != nil {
		return nil, err
	}

	if err := as.sqlServerLicenseTypeIDExists(contract.LicenseTypeID); err != nil {
		return nil, err
	}

	contract.ID = as.NewObjectID()

	err := as.Database.InsertSqlServerDatabaseContract(contract)
	if err != nil {
		return nil, err
	}

	agrs, err := as.GetSqlServerDatabaseContracts()
	if err != nil {
		return nil, err
	}

	for _, agr := range agrs {
		if agr.ID == contract.ID {
			return &agr, nil
		}
	}

	return nil, utils.NewError(errors.New("Can't find contract which has just been saved"))
}

func (as *APIService) sqlServerLicenseTypeIDExists(licenseTypeID string) error {
	_, err := as.GetSqlServerDatabaseLicenseType(licenseTypeID)
	if err != nil {
		return err
	}

	return nil
}

func (as *APIService) GetSqlServerDatabaseContracts() ([]model.SqlServerDatabaseContract, error) {
	contracts, err := as.Database.ListSqlServerDatabaseContracts()
	if err != nil {
		return nil, err
	}

	return contracts, nil
}

func (as *APIService) GetSqlServerDatabaseContractsAsXLSX() (*excelize.File, error) {
	contracts, err := as.GetSqlServerDatabaseContracts()
	if err != nil {
		return nil, err
	}

	sheet := "Contracts"
	headers := []string{
		"Type",
		"ContractID",
		"LicensesNumber",
		"Hosts",
		"Clusters",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range contracts {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.Type)
		sheets.SetCellValue(sheet, nextAxis(), val.ContractID)
		sheets.SetCellValue(sheet, nextAxis(), val.LicensesNumber)

		for _, val2 := range val.Hosts {
			sheets.DuplicateRow(sheet, axisHelp.GetIndexRow())
			duplicateRowNextAxis := axisHelp.NewRowSincePreviousColumn()

			sheets.SetCellValue(sheet, duplicateRowNextAxis(), val2)
		}

		for _, val2 := range val.Clusters {
			sheets.DuplicateRow(sheet, axisHelp.GetIndexRow())
			duplicateRowNextAxis := axisHelp.NewRowSincePreviousColumn()

			sheets.SetCellValue(sheet, duplicateRowNextAxis(), val2)
		}
	}

	return sheets, err
}

func (as *APIService) DeleteSqlServerDatabaseContract(id primitive.ObjectID) error {
	return as.Database.RemoveSqlServerDatabaseContract(id)
}

func (as *APIService) UpdateSqlServerDatabaseContract(contract model.SqlServerDatabaseContract) (*model.SqlServerDatabaseContract, error) {
	if err := checkHosts(as, contract.Hosts); err != nil {
		return nil, err
	}

	if err := as.sqlServerLicenseTypeIDExists(contract.LicenseTypeID); err != nil {
		return nil, err
	}

	if err := as.Database.UpdateSqlServerDatabaseContract(contract); err != nil {
		return nil, err
	}

	agrs, err := as.GetSqlServerDatabaseContracts()
	if err != nil {
		return nil, err
	}

	for _, agr := range agrs {
		if agr.ID == contract.ID {
			return &agr, nil
		}
	}

	return nil, utils.NewError(errors.New("Can't find contract which has just been saved"))
}
