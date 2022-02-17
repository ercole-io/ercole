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
	"github.com/360EntSecGroup-Skylar/excelize"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils/exutils"
)

func (as *APIService) AddMySQLAgreement(agreement model.MySQLAgreement) (*model.MySQLAgreement, error) {
	agreement.ID = as.NewObjectID()

	err := as.Database.AddMySQLAgreement(agreement)
	if err != nil {
		return nil, err
	}

	return &agreement, nil
}

func (as *APIService) UpdateMySQLAgreement(agreement model.MySQLAgreement) (*model.MySQLAgreement, error) {
	if err := as.Database.UpdateMySQLAgreement(agreement); err != nil {
		return nil, err
	}

	return &agreement, nil
}

func (as *APIService) GetMySQLAgreements() ([]model.MySQLAgreement, error) {
	agreements, err := as.Database.GetMySQLAgreements()
	if err != nil {
		return nil, err
	}

	return agreements, nil
}

func (as *APIService) DeleteMySQLAgreement(id primitive.ObjectID) error {
	if err := as.Database.DeleteMySQLAgreement(id); err != nil {
		return err
	}

	return nil
}

func (as *APIService) GetMySQLAgreementsAsXLSX() (*excelize.File, error) {
	agreements, err := as.GetMySQLAgreements()
	if err != nil {
		return nil, err
	}

	sheet := "Agreements"
	headers := []string{
		"Type",
		"Agreement Number",
		"CSI",
		"Number of licenses",
		"Clusters",
		"Host",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range agreements {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue(sheet, nextAxis(), val.Type)
		sheets.SetCellValue(sheet, nextAxis(), val.AgreementID)
		sheets.SetCellValue(sheet, nextAxis(), val.CSI)
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
