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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SqlServerDatabaseContract holds informations about a single SqlServerDatabaseContract
type SqlServerDatabaseContract struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Type           string             `json:"type" bson:"type"`
	ContractID     string             `json:"contractID" bson:"contractID"`
	LicenseTypeID  string             `json:"licenseTypeID" bson:"licenseTypeID"`
	LicensesNumber int                `json:"licensesNumber" bson:"licensesNumber"`
	Hosts          []string           `json:"hosts" bson:"hosts"`
	Clusters       []string           `json:"clusters" bson:"clusters"`
}

const (
	SqlServerContractTypeHost    string = "HOST"
	SqlServerContractTypeCluster string = "CLUSTER"
)
