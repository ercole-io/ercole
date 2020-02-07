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

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PatchingFunction holds all informations about a patching function
type PatchingFunction struct {
	ID        primitive.ObjectID `bson:"_id"`
	Hostname  string             `bson:"Hostname"`
	CreatedAt time.Time          `bson:"CreatedAt"`
	// PatchingFunction contains the javascript code that patch the hostdata
	// the hostdata is given via the hostdata global variable.
	// the static vars is given via the vars global variable
	// The function should be idempotent and reversible
	// e.g. PF(hostdata) == PF(PF(hostdata)) && ∃ PF⁻¹ | PF⁻¹(patchedHostData) == hostdata
	Code string                 `bson:"Code"`
	Vars map[string]interface{} `bson:"Vars"`
}

// HostDataBsonValidatorRules contains mongodb validation rules for hostData
var PatchingFunctionBsonValidatorRules = bson.D{
	{"bsonType", "object"},
	{"required", bson.A{
		"Hostname",
		"CreatedAt",
		"Code",
	}},
	{"properties", bson.D{
		{"Hostname", bson.D{
			{"bsonType", "string"},
		}},
		{"CreatedAt", bson.D{
			{"bsonType", "string"},
		}},
		{"Code", bson.D{
			{"bsonType", "string"},
		}},
		{"Vars", bson.D{
			{"bsonType", "object"},
		}},
	}},
}
