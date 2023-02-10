// Copyright (c) 2022 Sorint.lab S.p.A.
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
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OracleDatabaseContract holds informations about a sigle OracleDatabaseContract
type OracleDatabaseContract struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	ContractID        string             `json:"contractID" bson:"contractID" csv:"contract_id"`
	CSI               string             `json:"csi" bson:"csi" csv:"csi"`
	LicenseTypeID     string             `json:"licenseTypeID" bson:"licenseTypeID" csv:"license_type_id"`
	ReferenceNumber   string             `json:"referenceNumber" bson:"referenceNumber" csv:"reference_number"`
	Unlimited         bool               `json:"unlimited" bson:"unlimited" csv:"unlimited"`
	Count             int                `json:"count" bson:"count" csv:"count"`
	Basket            bool               `json:"basket" bson:"basket" csv:"basket"`
	Restricted        bool               `json:"restricted" bson:"restricted" csv:"restricted"`
	SupportExpiration *time.Time         `json:"supportExpiration" bson:"supportExpiration" csv:"-"`
	Hosts             []string           `json:"hosts" bson:"hosts" csv:"-"`
	HostsLiteral      LiteralStrSlice        `json:"-" bson:"-" csv:"hosts"`
}


func (contract OracleDatabaseContract) Check() error {
	if contract.Restricted && contract.Basket {
		return errors.New("If it's restricted it can't be basket")
	}

	return nil
}
