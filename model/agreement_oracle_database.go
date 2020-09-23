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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OracleDatabaseAgreement holds informations about a sigle OracleDatabaseAgreement
type OracleDatabaseAgreement struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	AgreementID     string             `json:"agreementID" bson:"agreementID"`
	PartID          string             `json:"partID" bson:"partID"`
	ItemDescription string             `json:"itemDescription" bson:"itemDescription"`
	Metrics         string             `json:"metrics" bson:"metrics"`
	CSI             string             `json:"csi" bson:"csi"`
	ReferenceNumber string             `json:"referenceNumber" bson:"referenceNumber"`
	Unlimited       bool               `json:"unlimited" bson:"unlimited"`
	Count           int                `json:"count" bson:"count"`
	CatchAll        bool               `json:"catchAll" bson:"catchAll"`
	Hosts           []string           `json:"hosts" bson:"hosts"`
}

// OracleDatabaseAgreementBsonValidatorRules contains mongodb validation rules for OracleDatabaseAgreement
var OracleDatabaseAgreementBsonValidatorRules = bson.M{
	"bsonType": "object",
	"required": bson.A{
		"agreementID",
		"partID",
		"itemDescription",
		"metrics",
		"csi",
		"referenceNumber",
		"unlimited",
		"count",
		"catchAll",
		"hosts",
	},
	"properties": bson.M{
		"agreementID": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 16,
		},
		"partID": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 16,
		},
		"itemDescription": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 64,
		},
		"metrics": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 32,
		},
		"csi": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 16,
		},
		"referenceNumber": bson.M{
			"bsonType":  "string",
			"minLength": 1,
			"maxLength": 16,
		},
		"unlimited": bson.M{
			"bsonType": "bool",
		},
		"count": bson.M{
			"bsonType": "int",
			"minimum":  0,
		},
		"catchAll": bson.M{
			"bsonType": "bool",
		},
		"hosts": bson.M{
			"bsonType":    "array",
			"uniqueItems": true,
			"items": bson.M{
				"bsonType":  "string",
				"minLength": 1,
				"maxLength": 253,
				"pattern":   `^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-_]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-_]*[A-Za-z0-9])$`,
			},
		},
	},
}
