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

package utils

import "fmt"

type AxisHelper struct {
	row int
}

func NewAxisHelper(headersOffset int) *AxisHelper {
	return &AxisHelper{
		row: 0 + headersOffset,
	}
}

func (ah *AxisHelper) NewRow() func() string {
	ah.row++

	columnOffset := -1
	return func() string {
		columnOffset++

		return fmt.Sprintf("%c%d", rune('A'+columnOffset), ah.row)
	}
}
