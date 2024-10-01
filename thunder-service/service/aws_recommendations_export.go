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
	"github.com/ercole-io/ercole/v2/thunder-service/dto"
)

func (ts *ThunderService) WriteAwsRecommendationsXlsx(recommendations []dto.AwsRecommendationDto) (*excelize.File, error) {
	f, err := excelize.OpenFile(ts.Config.ResourceFilePath + "/templates/template_aws_recommendations.xlsx")
	if err != nil {
		return nil, err
	}

	for _, recommendation := range recommendations {
		createAwsRecommendationSheetXlsx(f, recommendation.ObjectType, recommendation)
	}

	return f, nil
}

func createAwsRecommendationSheetXlsx(file *excelize.File, objectType string, recommendation dto.AwsRecommendationDto) {
	sheetName := objectType

	firstEmptyRow := len(file.GetRows(sheetName)) + 1

	file.SetCellValue(sheetName, fmt.Sprintf("A%d", firstEmptyRow), recommendation.Category)
	file.SetCellValue(sheetName, fmt.Sprintf("B%d", firstEmptyRow), recommendation.ObjectType)
	file.SetCellValue(sheetName, fmt.Sprintf("C%d", firstEmptyRow), recommendation.Suggestion)
	file.SetCellValue(sheetName, fmt.Sprintf("D%d", firstEmptyRow), recommendation.Name)

	for i, details := range recommendation.Details {
		for k := range details {
			if v, ok := details[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("%s%d", excelize.ToAlphaString(i+4), firstEmptyRow), v)
			}
		}
	}
}
