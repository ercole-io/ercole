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
	Database string `bson:"Database"`
	User     string `bson:"User"`
	Total    int    `bson:"Total"`
	Tables   int    `bson:"Tables"`
	Indexes  int    `bson:"Indexes"`
	LOB      int    `bson:"LOB"`
}

// SchemaBsonValidatorRules contains mongodb validation rules for schema
var SchemaBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"Database",
		"User",
		"Total",
		"Tables",
		"Indexes",
		"LOB",
	}},
	{"properties", bson.D{
		{"Database", bson.D{
			{"bsonType", "string"},
		}},
		{"User", bson.D{
			{"bsonType", "string"},
		}},
		{"Total", bson.D{
			{"bsonType", "number"},
		}},
		{"Tables", bson.D{
			{"bsonType", "number"},
		}},
		{"Indexes", bson.D{
			{"bsonType", "number"},
		}},
		{"LOB", bson.D{
			{"bsonType", "number"},
		}},
	}},
}
