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
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OciProfile struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Profile        string             `json:"profile" bson:"profile"`
	TenancyOCID    string             `json:"tenancyOCID" bson:"tenancyOCID"`
	UserOCID       string             `json:"userOCID" bson:"userOCID"`
	KeyFingerprint string             `json:"keyFingerprint" bson:"keyFingerprint"`
	Region         string             `json:"region" bson:"region"`
	PrivateKey     string             `json:"privateKey" bson:"privateKey"`
}

func (pr OciProfile) IsValid() bool {

	if pr.TenancyOCID[0:4] != "ocid" || pr.UserOCID[0:4] != "ocid" {
		return false
	}

	var tenancy []string = strings.Split(pr.TenancyOCID, ".")
	if tenancy[1] != "tenancy" {
		return false
	}

	var user []string = strings.Split(pr.UserOCID, ".")

	return user[1] == "user"
}
