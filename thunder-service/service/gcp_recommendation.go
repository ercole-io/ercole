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
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/job"
	"github.com/ercole-io/ercole/v2/utils/exutils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ts *ThunderService) ListGcpRecommendations() ([]model.GcpRecommendation, error) {
	selectedProfiles, err := ts.Database.GetActiveGcpProfiles()
	if err != nil {
		return nil, err
	}

	profileIDs := make([]primitive.ObjectID, 0, len(selectedProfiles))
	for _, p := range selectedProfiles {
		profileIDs = append(profileIDs, p.ID)
	}

	return ts.Database.ListGcpRecommendationsByProfiles(profileIDs)
}

func (ts *ThunderService) ForceGetGcpRecommendations() {
	j := &job.GcpDataRetrieveJob{
		Database: ts.Database,
		Config:   ts.Config,
		Log:      ts.Log,
		Opt:      nil,
	}

	j.Run()
}

func (ts *ThunderService) ListGcpError() ([]model.GcpError, error) {
	selectedProfiles, err := ts.Database.GetActiveGcpProfiles()
	if err != nil {
		return nil, err
	}

	profileIDs := make([]primitive.ObjectID, 0, len(selectedProfiles))
	for _, p := range selectedProfiles {
		profileIDs = append(profileIDs, p.ID)
	}

	return ts.Database.ListGcpErrorsByProfiles(profileIDs)
}

func (ts *ThunderService) CreateGcpRecommendationsXlsx() (*excelize.File, error) {
	recommendations, err := ts.ListGcpRecommendations()
	if err != nil {
		return nil, err
	}

	sheet := "Gcp recommendations"
	headers := []string{
		"Category",
		"Object Type",
		"Suggestion",
		"Project Name",
		"Profile ID",
		"Resource Name",
		"Resource ID",
		"Optimization Score",
		"Cpu Average",
		"Cpu Max",
		"Instance Name",
		"Mem Max",
		"Block Storage Name",
		"IOPS R MAX 5DD",
		"IOPS W MAX 5DD",
		"Size GB",
		"THROUGHPUT R MAX 5DD (MiBps)",
		"THROUGHPUT W MAX 5DD (MiBps)",
		"storage type",
	}

	sheets, err := exutils.NewXLSX(ts.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	for i, val := range recommendations {
		i += 2
		sheets.SetCellValue(sheet, fmt.Sprintf("A%d", i), val.Category)
		sheets.SetCellValue(sheet, fmt.Sprintf("B%d", i), val.ObjectType)
		sheets.SetCellValue(sheet, fmt.Sprintf("C%d", i), val.Suggestion)
		sheets.SetCellValue(sheet, fmt.Sprintf("D%d", i), val.ProjectName)
		sheets.SetCellValue(sheet, fmt.Sprintf("E%d", i), val.ProfileID.Hex())
		sheets.SetCellValue(sheet, fmt.Sprintf("F%d", i), val.ResourceName)
		sheets.SetCellValue(sheet, fmt.Sprintf("G%d", i), fmt.Sprintf("%v", val.ResourceID))
		sheets.SetCellValue(sheet, fmt.Sprintf("H%d", i), val.OptimizationScore)

		if v, ok := val.Details["Cpu Average"]; ok {
			sheets.SetCellValue(sheet, fmt.Sprintf("I%d", i), v)
		}

		if v, ok := val.Details["Cpu Max"]; ok {
			sheets.SetCellValue(sheet, fmt.Sprintf("J%d", i), v)
		}

		if v, ok := val.Details["Instance Name"]; ok {
			sheets.SetCellValue(sheet, fmt.Sprintf("K%d", i), v)
		}

		if v, ok := val.Details["Mem Max"]; ok {
			sheets.SetCellValue(sheet, fmt.Sprintf("L%d", i), v)
		}

		if v, ok := val.Details["Block Storage Name"]; ok {
			sheets.SetCellValue(sheet, fmt.Sprintf("M%d", i), v)
		}

		if v, ok := val.Details["IOPS R MAX 5DD"]; ok {
			sheets.SetCellValue(sheet, fmt.Sprintf("N%d", i), v)
		}

		if v, ok := val.Details["IOPS W MAX 5DD"]; ok {
			sheets.SetCellValue(sheet, fmt.Sprintf("O%d", i), v)
		}

		if v, ok := val.Details["Size GB"]; ok {
			sheets.SetCellValue(sheet, fmt.Sprintf("P%d", i), v)
		}

		if v, ok := val.Details["THROUGHPUT R MAX 5DD (MiBps)"]; ok {
			sheets.SetCellValue(sheet, fmt.Sprintf("Q%d", i), v)
		}

		if v, ok := val.Details["THROUGHPUT W MAX 5DD (MiBps)"]; ok {
			sheets.SetCellValue(sheet, fmt.Sprintf("R%d", i), v)
		}

		if v, ok := val.Details["storage type"]; ok {
			sheets.SetCellValue(sheet, fmt.Sprintf("S%d", i), v)
		}
	}

	return sheets, err
}
