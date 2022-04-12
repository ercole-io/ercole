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

import "go.mongodb.org/mongo-driver/bson/primitive"

type MySQLContract struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	Type             string             `json:"type" bson:"type"`
	ContractID       string             `json:"contractID" bson:"contractID"`
	CSI              string             `json:"csi" bson:"csi"`
	NumberOfLicenses uint               `json:"numberOfLicenses" bson:"numberOfLicenses"`
	Clusters         []string           `json:"clusters" bson:"clusters"`
	Hosts            []string           `json:"hosts" bson:"hosts"`
}

const (
	MySQLContractTypeHost    string = "HOST"
	MySQLContractTypeCluster string = "CLUSTER"
)

func getMySQLContractTypes() []string {
	return []string{MySQLContractTypeHost, MySQLContractTypeCluster}
}

func (agr MySQLContract) IsValid() bool {
	if agr.ContractID == "" || agr.CSI == "" || agr.NumberOfLicenses == 0 {
		return false
	}

	fields := make(map[string][]string)
	fields[agr.Type] = getMySQLContractTypes()

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
