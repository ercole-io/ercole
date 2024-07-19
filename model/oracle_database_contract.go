// Copyright (c) 2024 Sorint.lab S.p.A.
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
	"encoding/json"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OracleDatabaseContract holds informations about a sigle OracleDatabaseContract
type OracleDatabaseContract struct {
	ID                primitive.ObjectID `json:"id" bson:"_id" csv:"-"`
	ContractID        string             `json:"contractID" bson:"contractID" csv:"Contract Number"`
	CSI               string             `json:"csi" bson:"csi" csv:"CSI"`
	LicenseTypeID     string             `json:"licenseTypeID" bson:"licenseTypeID" csv:"Part Number"`
	ReferenceNumber   string             `json:"referenceNumber" bson:"referenceNumber" csv:"Reference Number,omitempty"`
	Unlimited         bool               `json:"unlimited" bson:"unlimited" csv:"ULA,omitempty"`
	Count             int                `json:"count" bson:"count" csv:"License number,omitempty"`
	Basket            bool               `json:"basket" bson:"basket" csv:"Basket,omitempty"`
	Restricted        bool               `json:"restricted" bson:"restricted" csv:"Restricted,omitempty"`
	SupportExpiration *dateTime          `json:"supportExpiration" bson:"supportExpiration" csv:"Support Expiration,omitempty"`
	Hosts             []string           `json:"hosts" bson:"hosts" csv:"-"`
	HostsLiteral      LiteralStrSlice    `json:"-" bson:"-" csv:"Hosts,omitempty"`
	Status            string             `json:"status" bson:"status" csv:"Status,omitempty"`
	ProductOrderDate  *dateTime          `json:"productOrderDate" bson:"productOrderDate" csv:"Product Order Date,omitempty"`
}

func (contract OracleDatabaseContract) Check() error {
	if contract.Restricted && contract.Basket {
		return errors.New("if it's restricted it can't be basket")
	}

	return nil
}

type dateTime struct {
	time.Time
}

func (d *dateTime) MarshalCSV() (string, error) {
	return d.Time.Format("02/01/2006"), nil
}

func (d *dateTime) UnmarshalCSV(csv string) (err error) {
	d.Time, err = time.Parse("02/01/2006", csv)
	return err
}

func (d dateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time)
}

func (d *dateTime) UnmarshalJSON(data []byte) error {
	var t time.Time

	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	d.Time = t

	return nil
}

func (d dateTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if d.IsZero() {
		return bson.MarshalValue(nil)
	}

	return bson.MarshalValue(d.Time)
}

func (d *dateTime) UnmarshalBSON(data []byte) error {
	var t time.Time

	err := bson.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	d.Time = t

	return nil
}
