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

package model

type MySQLSegmentAdvisor struct {
	TableSchema string  `json:"tableSchema" bson:"tableSchema"`
	TableName   string  `json:"tableName" bson:"tableName"`
	Engine      string  `json:"engine" bson:"engine"`
	Allocation  float64 `json:"allocation" bson:"allocation"` // in MB
	Data        float64 `json:"data" bson:"data"`             // in MB
	Index       float64 `json:"index" bson:"index"`           // in MB
	Free        float64 `json:"free" bson:"free"`             // in MB
}
