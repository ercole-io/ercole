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
	MaxConnections  int                  `json:"maxConnections" bson:"maxConnections"`
	Port            int                  `json:"port" bson:"port"`
	InstanceSize    int                  `json:"instanceSize" bson:"instanceSize"`
	Charset         string               `json:"charset" bson:"charset"`
	Isinreplica     bool                 `json:"isinreplica" bson:"isinreplica"`
	Ismaster        bool                 `json:"ismaster" bson:"ismaster"`
	Isslave         bool                 `json:"isslave" bson:"isslave"`
	ArchiverWorking bool                 `json:"archiverWorking" bson:"archiverWorking"`
	SlavesNum       int                  `json:"slavesNum" bson:"slavesNum"`
	UsersNum        int                  `json:"usersNum" bson:"usersNum"`
	DbNum           int                  `json:"dbNum" bson:"dbNum"`
	TblspNum        int                  `json:"tblspNum" bson:"tblspNum"`
	TrustHbaEntries int                  `json:"trustHbaEntries" bson:"trustHbaEntries"`
	Setting         *PostgreSQLSetting   `json:"setting" bson:"setting"`
	Databases       []PostgreSQLDatabase `json:"databases" bson:"databases"`
}
