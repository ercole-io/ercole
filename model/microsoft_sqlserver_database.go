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

import (
	"reflect"

	godynstruct "github.com/amreo/go-dyn-struct"
)

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
	OtherInfo       map[string]interface{}                 `json:"-" bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v MicrosoftSQLServerDatabase) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *MicrosoftSQLServerDatabase) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v MicrosoftSQLServerDatabase) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *MicrosoftSQLServerDatabase) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}
