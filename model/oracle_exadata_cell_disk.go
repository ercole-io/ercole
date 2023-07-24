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

import "time"

// OracleExadataStorageCell holds info about a exadata cell disk
type OracleExadataStorageCell struct {
	Type       string                  `json:"type" bson:"type"`
	Hostname   string                  `json:"hostname" bson:"hostname"`
	CellDisk   string                  `json:"cellDisk" bson:"cellDisk"`
	Cell       string                  `json:"cell" bson:"cell"`
	Size       string                  `json:"size" bson:"size"`
	FreeSpace  string                  `json:"freeSpace" bson:"freeSpace"`
	Status     string                  `json:"status" bson:"status"`
	ErrorCount int                     `json:"errorCount" bson:"errorCount"`
	GridDisks  []OracleExadataGridDisk `json:"gridDisks,omitempty" bson:"gridDisks"`
	Databases  []OracleExadataDatabase `json:"databases" bson:"databases"`
}

type OracleExadataGridDisk struct {
	Type          string `json:"type" bson:"type"`
	Hostname      string `json:"hostname" bson:"hostname"`
	GridDisk      string `json:"gridDisk" bson:"gridDisk"`
	CellDisk      string `json:"cellDisk" bson:"cellDisk"`
	Size          string `json:"size" bson:"size"`
	Status        string `json:"status" bson:"status"`
	ErrorCount    int    `json:"errorCount" bson:"errorCount"`
	CachingPolicy string `json:"cachingPolicy" bson:"cachingPolicy"`
	AsmDiskName   string `json:"asmDiskName" bson:"asmDiskName"`
	AsmDiskGroup  string `json:"asmDiskGroup" bson:"asmDiskGroup"`
	AsmDiskSize   string `json:"asmDiskSize" bson:"asmDiskSize"`
	AsmDiskStatus string `json:"asmDiskStatus" bson:"asmDiskStatus"`
}

type OracleExadataDatabase struct {
	Type            string     `json:"type" bson:"type"`
	DbName          string     `json:"dbName" bson:"dbName"`
	Cell            string     `json:"cell" bson:"cell"`
	DbID            int        `json:"dbID" bson:"dbID"`
	FlashCacheLimit int        `json:"flashCacheLimit" bson:"flashCacheLimit"`
	IormShare       int        `json:"iormShare" bson:"iormShare"`
	LastIOReq       *time.Time `json:"lastIOReq" bson:"lastIOReq"`
}
