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

type PostgreSQLInstance struct {
	CurrentConnections string               `json:"currentConnections" bson:"currentConnections"`
	ManualConnections  int                  `json:"manualConnections" bson:"manualConnections"`
	InstanceSize       string               `json:"instanceSize" bson:"instanceSize"`
	Charset            string               `json:"charset" bson:"charset"`
	IsInreplica        bool                 `isinreplica:"isinreplica" bson:"currentConnections"`
	IsMaster           bool                 `json:"isMaster" bson:"isMaster"`
	IsSlave            bool                 `json:"isSlave" bson:"isSlave"`
	SlavesNum          int                  `json:"slavesNum" bson:"slavesNum"`
	UsersNum           int                  `json:"usersNum" bson:"usersNum"`
	DbNum              int                  `json:"dbNum" bson:"dbNum"`
	TblspNum           int                  `json:"tblspNum" bson:"tblspNum"`
	TrustHbaEntries    int                  `json:"trustHbaEntries" bson:"trustHbaEntries"`
	Databases          []PostgreSQLDatabase `json:"databases" bson:"databases"`
	Settings           []PostgreSQLSetting  `json:"settings" bson:"settings"`
}
