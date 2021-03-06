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

import "go.mongodb.org/mongo-driver/bson/primitive"

type MySQLAgreement struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	Type             string             `json:"type" bson:"type"`
	AgreementID      string             `json:"agreementID" bson:"agreementID"`
	CSI              string             `json:"csi" bson:"csi"`
	NumberOfLicenses uint               `json:"numberOfLicenses" bson:"numberOfLicenses"`
	Clusters         []string           `json:"clusters" bson:"clusters"`
	Hosts            []string           `json:"hosts" bson:"hosts"`
}

const (
	MySQLAgreementTypeHost    string = "HOST"
	MySQLAgreementTypeCluster string = "CLUSTER"
)

func getMySQLAgreementTypes() []string {
	return []string{MySQLAgreementTypeHost, MySQLAgreementTypeCluster}
}

func (agr MySQLAgreement) IsValid() bool {
	if agr.AgreementID == "" || agr.CSI == "" || agr.NumberOfLicenses == 0 {
		return false
	}

	fields := make(map[string][]string)
	fields[agr.Type] = getMySQLAgreementTypes()

fields:
	for thisValue, allValidValues := range fields {
		for _, validValue := range allValidValues {
			if thisValue == validValue {
				continue fields
			}
		}

		return false
	}

	return true
}
