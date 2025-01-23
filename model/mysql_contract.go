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
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MySQLContract struct {
	ID                primitive.ObjectID `json:"id" bson:"_id" csv:"-"`
	Type              string             `json:"type" bson:"type" csv:"Type"`
	ContractID        string             `json:"contractID" bson:"contractID" csv:"Contract Number"`
	CSI               string             `json:"csi" bson:"csi" csv:"CSI"`
	LicenseTypeID     string             `json:"licenseTypeID" bson:"licenseTypeID" csv:"License Type"`
	NumberOfLicenses  uint               `json:"numberOfLicenses" bson:"numberOfLicenses" csv:"Number of Licenses"`
	SupportExpiration *time.Time         `json:"supportExpiration" bson:"supportExpiration" csv:"-"`
	Clusters          []string           `json:"clusters" bson:"clusters" csv:"-"`
	Hosts             []string           `json:"hosts" bson:"hosts" csv:"-"`
	HostsLiteral      LiteralStrSlice    `json:"-" bson:"-" csv:"-"`
	ClusterLiteral    LiteralStrSlice    `json:"-" bson:"-" csv:"-"`
	Location          string             `json:"location" bson:"location" csv:"Location"`
}

const (
	MySQLContractTypeHost    string = "HOST"
	MySQLContractTypeCluster string = "CLUSTER"
)

const MySqlPartNumber = "B64911"

const MySqlItemDescription = "MySQL Enterprise Edition"

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
