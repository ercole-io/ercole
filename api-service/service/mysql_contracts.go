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

// Package service is a package that provides methods for querying data
package service

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils/exutils"
)

func (as *APIService) AddMySQLContract(contract model.MySQLContract) (*model.MySQLContract, error) {
	contract.ID = as.NewObjectID()

	err := as.Database.AddMySQLContract(contract)
	if err != nil {
		return nil, err
	}

	return &contract, nil
}

func (as *APIService) UpdateMySQLContract(contract model.MySQLContract) (*model.MySQLContract, error) {
	if err := as.Database.UpdateMySQLContract(contract); err != nil {
		return nil, err
	}

	return &contract, nil
}

func (as *APIService) GetMySQLContracts() ([]model.MySQLContract, error) {
	contracts, err := as.Database.GetMySQLContracts()
	if err != nil {
		return nil, err
	}

	return contracts, nil
}

func (as *APIService) DeleteMySQLContract(id primitive.ObjectID) error {
	if err := as.Database.DeleteMySQLContract(id); err != nil {
		return err
	}

	return nil
}

func (as *APIService) GetMySQLContractsAsXLSX() (*excelize.File, error) {
	contracts, err := as.GetMySQLContracts()
	if err != nil {
		return nil, err
	}

	sheet := "Contracts"
	headers := []string{
		"Type",
		"Contract Number",
		"CSI",
		"Support Expiration",
		"Number of licenses",
		"Clusters",
		"Host",
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
		sheets.SetCellValue(sheet, nextAxis(), val.CSI)

		if val.SupportExpiration != nil {
			sheets.SetCellValue(sheet, nextAxis(), val.SupportExpiration)
		} else {
			sheets.SetCellValue(sheet, nextAxis(), "")
		}

		sheets.SetCellValue(sheet, nextAxis(), val.NumberOfLicenses)
		sheets.SetCellValue(sheet, nextAxis(), val.Clusters)

		for _, val2 := range val.Hosts {
			sheets.DuplicateRow(sheet, axisHelp.GetIndexRow())
			duplicateRowNextAxis := axisHelp.NewRowSincePreviousColumn()

			sheets.SetCellValue(sheet, duplicateRowNextAxis(), val2)
		}
	}

	return sheets, err
}
