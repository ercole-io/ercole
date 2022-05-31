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
	DbSize            string `json:"dbSize" bson:"dbSize"`
	DbConnections     string `json:"dbConnections" bson:"dbConnections"`
	TableCount        int    `json:"tableCount" bson:"tableCount"`
	IndexCount        int    `json:"indexCount" bson:"indexCount"`
	TablesSize        string `json:"tablesSize" bson:"tablesSize"`
	IndexesSize       string `json:"indexesSize" bson:"indexesSize"`
	MviewsSize        string `json:"mviewsSize" bson:"mviewsSize"`
	ExtensionCount    int    `json:"extensionCount" bson:"extensionCount"`
	SchemaCount       int    `json:"schemaCount" bson:"schemaCount"`
	LogicReplSetup    bool   `json:"logicReplSetup" bson:"logicReplSetup"`
	PublicationCount  int    `json:"publicationCount" bson:"publicationCount"`
	SubscriptionCount int    `json:"subscriptionCount" bson:"subscriptionCount"`
	LockCount         int    `json:"lockCount" bson:"lockCount"`
	MviewCount        int    `json:"mviewCount" bson:"mviewCount"`
	ViewCount         int    `json:"viewCount" bson:"viewCount"`
}
