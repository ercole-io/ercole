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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OracleDatabaseAgreement holds informations about a sigle OracleDatabaseAgreement
type OracleDatabaseAgreement struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	AgreementID string             `json:"agreementID" bson:"agreementID"`
	CSI         string             `json:"csi" bson:"csi"`
	Parts       []AssociatedPart   `json:"parts" bson:"parts"`
}

// AssociatedPart describe a Part associated to an Agreement
type AssociatedPart struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id"`
	OracleDatabasePart `bson:",inline"`
	ReferenceNumber    string   `json:"referenceNumber" bson:"referenceNumber"`
	Unlimited          bool     `json:"unlimited" bson:"unlimited"`
	Count              int      `json:"count" bson:"count"`
	CatchAll           bool     `json:"catchAll" bson:"catchAll"` //TODO Rename in IsBasket ?
	Hosts              []string `json:"hosts" bson:"hosts"`
}

// AssociatedPartByID getter
func (agreement *OracleDatabaseAgreement) AssociatedPartByID(associatedPartID primitive.ObjectID,
) (associatedPart *AssociatedPart) {
	for i := range agreement.Parts {
		if agreement.Parts[i].ID == associatedPartID {
			return &agreement.Parts[i]
		}
	}

	return
}
