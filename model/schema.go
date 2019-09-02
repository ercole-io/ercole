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

// Schema holds information about Oracle database schema.
type Schema struct {
	Database string
	User     string
	Total    int
	Tables   int
	Indexes  int
	LOB      int
}

var SchemaBsonValidatorRules bson.D = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"database",
		"user",
		"total",
		"tables",
		"indexes",
		"lob",
	}},
	{"properties", bson.D{
		{"database", bson.D{
			{"bsonType", "string"},
		}},
		{"user", bson.D{
			{"bsonType", "string"},
		}},
		{"total", bson.D{
			{"bsonType", "int"},
		}},
		{"tables", bson.D{
			{"bsonType", "int"},
		}},
		{"indexes", bson.D{
			{"bsonType", "int"},
		}},
		{"lob", bson.D{
			{"bsonType", "int"},
		}},
	}},
}
