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
	"go.mongodb.org/mongo-driver/bson"
)

// LicenseCount holds information about Oracle database license
type LicenseCount struct {
	Name             string  `json:"name" bson:"_id"`
	Count            int     `json:"count" bson:"count"`
	CostPerProcessor float64 `json:"costPerProcessor" bson:"costPerProcessor"`
	Unlimited        bool    `json:"unlimited" bson:"unlimited"`
}

// LicenseCountBsonValidatorRules contains mongodb validation rules for licenseCount
var LicenseCountBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"_id",
		"count",
		"costPerProcessor",
		"unlimited",
	},
	"properties": bson.M{
		"_id": bson.M{
			"bsonType": "string",
		},
		"count": bson.M{
			"bsonType": "number",
		},
		"costPerProcessor": bson.M{
			"bsonType": "number",
		},
		"unlimited": bson.M{
			"bsonType": "bool",
		},
	},
}
