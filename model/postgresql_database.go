// Copyright (c) 2022 Sorint.lab S.p.A.
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

type PostgreSQLDatabase struct {
	DbName             string             `json:"dbName" bson:"dbName"`
	DbOwner            string             `json:"dbOwner" bson:"dbOwner"`
	Datconnlimit       int                `json:"datconnlimit" bson:"datconnlimit"`
	SchemasCount       int                `json:"schemasCount" bson:"schemasCount"`
	DbSize             int                `json:"dbSize" bson:"dbSize"`
	TablesCount        int                `json:"tablesCount" bson:"tablesCount"`
	TablesSize         int                `json:"tablesSize" bson:"tablesSize"`
	IndexesCount       int                `json:"indexesCount" bson:"indexesCount"`
	IndexesSize        int                `json:"indexesSize" bson:"indexesSize"`
	MatviewsCount      int                `json:"matviewsCount" bson:"matviewsCount"`
	MatviewsSize       int                `json:"matviewsSize" bson:"matviewsSize"`
	ExtensionsCount    int                `json:"extensionsCount" bson:"extensionsCount"`
	LobsCount          int                `json:"lobsCount" bson:"lobsCount"`
	LobsSize           int                `json:"lobsSize" bson:"lobsSize"`
	ViewsCount         int                `json:"viewsCount" bson:"viewsCount"`
	LogicReplSetup     bool               `json:"logicReplSetup" bson:"logicReplSetup"`
	PublicationsCount  int                `json:"publicationsCount" bson:"publicationsCount"`
	SubscriptionsCount int                `json:"subscriptionsCount" bson:"subscriptionsCount"`
	Schemas            []PostgreSQLSchema `json:"schemas" bson:"schemas"`
}
