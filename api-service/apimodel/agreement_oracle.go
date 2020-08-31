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

package apimodel

// OracleDatabaseAgreementsAddRequest contains the informations needed to add new agreements
type OracleDatabaseAgreementsAddRequest struct {
	AgreementID     string   `json:"agreementID" bson:"agreementID"`
	PartsID         []string `json:"partsID" bson:"partsID"`
	CSI             string   `json:"csi" bson:"csi"`
	ReferenceNumber string   `json:"referenceNumber" bson:"referenceNumber"`
	Unlimited       bool     `json:"unlimited" bson:"unlimited"`
	Count           int      `json:"count" bson:"count"`
	CatchAll        bool     `json:"catchAll" bson:"catchAll"`
	Hosts           []string `json:"hosts" bson:"hosts"`
}
