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
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/model"
)

func (ts *ThunderService) WriteOciRecommendationsXlsx(recommendations []model.OciRecommendation) (*excelize.File, error) {
	f, err := excelize.OpenFile(ts.Config.ResourceFilePath + "/templates/template_oci_recommendations.xlsx")
	if err != nil {
		return nil, err
	}

	for _, recommendation := range recommendations {
		createOciRecommendationSheetXlsx(f, recommendation.Category, recommendation)
	}

	return f, nil
}

func createOciRecommendationSheetXlsx(file *excelize.File, category string, recommendation model.OciRecommendation) {
	sheetName := category

	firstEmptyRow := len(file.GetRows(sheetName)) + 1

	file.SetCellValue(sheetName, fmt.Sprintf("A%d", firstEmptyRow), recommendation.Category)
	file.SetCellValue(sheetName, fmt.Sprintf("B%d", firstEmptyRow), recommendation.ObjectType)
	file.SetCellValue(sheetName, fmt.Sprintf("C%d", firstEmptyRow), recommendation.Suggestion)
	file.SetCellValue(sheetName, fmt.Sprintf("D%d", firstEmptyRow), recommendation.Name)

	for _, detail := range recommendation.Details {
		file.SetCellValue(sheetName, fmt.Sprintf("%s%d", getColumnLetterByHeader(file.GetRows(sheetName)[0], detail.Name), 1), detail.Name)
		file.SetCellValue(sheetName, fmt.Sprintf("%s%d", getColumnLetterByHeader(file.GetRows(sheetName)[0], detail.Name), firstEmptyRow), detail.Value)
	}
}

func columnNumberToLetter(colNum int) string {
	letter := ""

	for colNum > 0 {
		colNum--
		letter = string(rune('A'+(colNum%26))) + letter
		colNum /= 26
	}

	return letter
}

func getColumnLetterByHeader(headerRow []string, header string) string {
	for colIndex, colValue := range headerRow {
		if strings.EqualFold(colValue, header) {
			return columnNumberToLetter(colIndex + 1)
		}
	}

	for colIndex, colValue := range headerRow {
		if colValue == "" {
			return columnNumberToLetter(colIndex + 1)
		}
	}

	nextCol := len(headerRow) + 1

	return columnNumberToLetter(nextCol)
}
