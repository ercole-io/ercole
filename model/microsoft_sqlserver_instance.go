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

type MicrosoftSQLServerInstance struct {
	Status        string                       `json:"status" bson:"status"`
	Name          string                       `json:"name" bson:"name"`
	DisplayName   string                       `json:"displayName" bson:"displayName"`
	ServerName    string                       `json:"serverName" bson:"serverName"`
	DatabaseID    int                          `json:"databaseID" bson:"databaseID"`
	StateDesc     string                       `json:"stateDesc" bson:"stateDesc"`
	Version       string                       `json:"version" bson:"version"`
	Platform      string                       `json:"platform" bson:"platform"`
	CollationName string                       `json:"collationName" bson:"collationName"`
	Edition       string                       `json:"edition" bson:"edition"`
	EditionType   string                       `json:"editionType" bson:"editionType"`
	ProductCode   string                       `json:"productCode" bson:"productCode"`
	LicensingInfo string                       `json:"licensingInfo" bson:"licensingInfo"`
	Databases     []MicrosoftSQLServerDatabase `json:"databases" bson:"databases"`
}
