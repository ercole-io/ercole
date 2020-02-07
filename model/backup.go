// Copyright (c) 2019 Sorint.lab S.p.A.
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

import "go.mongodb.org/mongo-driver/bson"

// Backup holds informations about a backup
type Backup struct {
	BackupType string `bson:"BackupType"`
	Hour       string `bson:"Hour"`
	WeekDays   string `bson:"WeekDays"`
	AvgBckSize string `bson:"AvgBckSize"`
	Retention  string `bson:"Retention"`
}

// BackupBsonValidatorRules contains mongodb validation rules for backup
var BackupBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"BackupType",
		"Hour",
		"WeekDays",
		"AvgBckSize",
		"Retention",
	}},
	{"properties", bson.D{
		{"BackupType", bson.D{
			{"bsonType", "string"},
		}},
		{"Hour", bson.D{
			{"bsonType", "string"},
		}},
		{"WeekDays", bson.D{
			{"bsonType", "string"},
		}},
		{"AvgBckSize", bson.D{
			{"bsonType", "string"},
		}},
		{"Retention", bson.D{
			{"bsonType", "string"},
		}},
	}},
}
