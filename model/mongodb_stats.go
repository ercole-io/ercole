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

type DBStats struct {
	DBName               string  `json:"dbName" bson:"dbName"`
	Charset              string  `json:"charset" bson:"charset"`
	Collections          int32   `json:"collections" bson:"collections"`
	Views                int32   `json:"views" bson:"views"`
	Objects              int32   `json:"objects" bson:"objects"`
	Indexes              int32   `json:"indexes" bson:"indexes"`
	DataSize             float64 `json:"dataSize" bson:"dataSize"`
	IndexSize            float64 `json:"indexSize" bson:"indexSize"`
	StorageSize          float64 `json:"storageSize" bson:"storageSize"`
	TotalSize            float64 `json:"totalSize" bson:"totalSize"`
	FsUsedSize           float64 `json:"fsUsedSize" bson:"fsUsedSize"`
	FsTotalSize          float64 `json:"fsTotalSize" bson:"fsTotalSize"`
	FreeStorageSize      float64 `json:"freeStorageSize" bson:"freeStorageSize"`
	IndexFreeStorageSize float64 `json:"indexFreeStorageSize" bson:"indexFreeStorageSize"`
	TotalFreeStorageSize float64 `json:"totalFreeStorageSize" bson:"totalFreeStorageSize"`
	Users                int     `json:"users" bson:"users"`
	ShardDBs             []struct {
		ID          string `json:"_id" bson:"_id"`
		Primary     string `json:"primary" bson:"primary"`
		Partitioned bool   `json:"partitioned" bson:"partitioned"`
	} `json:"shardDBs" bson:"shardDBs"`
}
