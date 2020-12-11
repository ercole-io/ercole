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

// Package dto is a package that provides struct that contains charts
package dto

type TechnologyTypesChart struct {
	Databases        []TechnologyTypeChartBubble `json:"databases"`
	Middlewares      []TechnologyTypeChartBubble `json:"middlewares"`
	OperatingSystems []TechnologyTypeChartBubble `json:"operatingSystems"`
	Legend           ChartLegend                 `json:"legend"`
}

type TechnologyTypeChartBubble struct {
	Name string  `json:"name" bson:"name"`
	Size float64 `json:"size" bson:"size"`
}
