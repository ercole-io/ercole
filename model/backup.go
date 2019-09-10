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
	BackupType string `bson:"backup_type"`
	Hour       string
	WeekDays   string `bson:"week_days"`
	AvgBckSize string `bson:"avg_bck_size"`
	Retention  string
}

// BackupBsonValidatorRules contains mongodb validation rules for backup
var BackupBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"backup_type",
		"hour",
		"week_days",
		"avg_bck_size",
		"retention",
	}},
	{"properties", bson.D{
		{"backup_type", bson.D{
			{"bsonType", "string"},
		}},
		{"hour", bson.D{
			{"bsonType", "string"},
		}},
		{"week_days", bson.D{
			{"bsonType", "string"},
		}},
		{"avg_bck_size", bson.D{
			{"bsonType", "string"},
		}},
		{"retention", bson.D{
			{"bsonType", "string"},
		}},
	}},
}
