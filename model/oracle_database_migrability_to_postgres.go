// Copyright (c) 2023 Sorint.lab S.p.A.
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

type PgsqlMigrability struct {
	Metric     *string `json:"metric,omitempty" bson:"metric,omitempty"`
	Count      int
	Schema     *string `json:"schema,omitempty" bson:"schema,omitempty"`
	ObjectType *string `json:"objectType,omitempty" bson:"objectType,omitempty"`
}

func (pm PgsqlMigrability) GetMetric() string {
	if pm.Metric == nil {
		return ""
	}

	return *pm.Metric
}

func (pm PgsqlMigrability) GetSchema() string {
	if pm.Schema == nil {
		return ""
	}

	return *pm.Schema
}

func (pm PgsqlMigrability) GetObjectType() string {
	if pm.ObjectType == nil {
		return ""
	}

	return *pm.ObjectType
}
