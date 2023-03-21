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

type MongoDBInstance struct {
	Name             string                 `json:"name" bson:"name"`
	Version          string                 `json:"version" bson:"version"`
	Dbs              int                    `json:"dbs" bson:"dbs"`
	ReplicaSet       HelloResult            `json:"replicaSet" bson:"replicaSet"`
	ShardList        ShardStatus            `json:"shardList" bson:"shardList"`
	StatusConnection ServerStatusConnection `json:"statusConnection" bson:"statusConnection"`
	Stats            []DBStats              `json:"dbStats" bson:"dbStats"`
}
