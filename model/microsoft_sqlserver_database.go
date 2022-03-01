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

type MicrosoftSQLServerDatabase struct {
	DatabaseID      int                                    `json:"databaseID" bson:"databaseID"`
	Name            string                                 `json:"name" bson:"name"`
	CollationName   string                                 `json:"collationName" bson:"collationName"`
	Status          string                                 `json:"status" bson:"status"`
	RecoveryModel   string                                 `json:"recoveryModel" bson:"recoveryModel"`
	BlockSize       int                                    `json:"blockSize" bson:"blockSize"`
	SchedulersCount int                                    `json:"schedulersCount" bson:"schedulersCount"`
	AffinityMask    int                                    `json:"affinityMask" bson:"affinityMask"`
	MinServerMemory int                                    `json:"minServerMemory" bson:"minServerMemory"`
	MaxServerMemory int                                    `json:"maxServerMemory" bson:"maxServerMemory"`
	CTP             int                                    `json:"ctp" bson:"ctp"`
	MaxDop          int                                    `json:"maxDop" bson:"maxDop"`
	Alloc           float64                                `json:"alloc" bson:"alloc"`
	Backups         []MicrosoftSQLServerDatabaseBackup     `json:"backups" bson:"backups"`
	Schemas         []MicrosoftSQLServerDatabaseSchema     `json:"schemas" bson:"schemas"`
	Tablespaces     []MicrosoftSQLServerDatabaseTablespace `json:"tablespaces" bson:"tablespaces"`
}
