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

type MicrosoftSQLServerProductFeature struct {
	Product    string `json:"product" bson:"product"`
	Feature    string `json:"feature" bson:"feature"`
	Language   string `json:"language" bson:"language"`
	Edition    string `json:"edition" bson:"edition"`
	Version    string `json:"version" bson:"version"`
	Clustered  bool   `json:"clustered" bson:"clustered"`
	Configured bool   `json:"configured" bson:"configured"`
}
