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

package exutils

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/utils"
)

type AxisHelper struct {
	row          int
	columnOffset int
}

// NewAxisHelper is made to help feel and excel file
// headersOffset is the number of line occupied by headers
// e.g.:
// axisHelp := exutils.NewAxisHelper(1) // start from second line
// for _, val := range aSliceWithValuesForEachRow{
// 	nextAxis := axisHelp.NewRow()

// 	file.SetCellValue(sheet, nextAxis(), val.Name) // A2, A3... and so on
// 	file.SetCellValue(sheet, nextAxis(), val.Surname)
// 	...
// }
func NewAxisHelper(headersOffset int) *AxisHelper {
	return &AxisHelper{
		row: 0 + headersOffset,
	}
}

// NewRow must be used when you want to go to the next row
// It returns a `nextAxis` anonymous func that return axis for the new row
func (ah *AxisHelper) NewRow() func() string {
	ah.row++

	columnOffset := -1
	return func() string {
		columnOffset++
		ah.columnOffset = columnOffset

		return fmt.Sprintf("%c%d", rune('A'+columnOffset), ah.row)
	}
}

func (ah *AxisHelper) GetIndexRow() int {
	return ah.row
}

// NewRow must be used when you want to go to the next row
// It returns a `nextAxis` anonymous func that return axis for the new row
func (ah *AxisHelper) NewRowSincePreviousColumn() func() string {
	ah.row++

	columnOffset := ah.columnOffset
	return func() string {
		columnOffset++

		return fmt.Sprintf("%c%d", rune('A'+columnOffset), ah.row)
	}
}

// NewXLSX return *excelize.File initialized with sheet name and headers
func NewXLSX(c config.Configuration, sheet string, headers ...string) (*excelize.File, error) {
	file, err := excelize.OpenFile(c.ResourceFilePath + "/templates/template_generic.xlsx")
	if err != nil {
		return nil, utils.NewError(err, "READ_TEMPLATE")
	}

	file.SetSheetName("Sheet1", sheet)

	for i, val := range headers {
		column := rune('A' + i)
		file.SetCellValue(sheet, fmt.Sprintf("%c1", column), val)
	}

	return file, nil
}
