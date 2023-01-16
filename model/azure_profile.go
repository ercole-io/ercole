// Copyright (c) 2023 Sorint.lab S.p.A.
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

type AzureProfile struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	TenantId       string             `json:"tenantid" bson:"tenantid"`
	ClientId       string             `json:"clientid" bson:"clientid"`
	ClientSecret   *string            `json:"clientsecret" bson:"clientsecret"`
	SubscriptionId string             `json:"subscriptionid" bson:"subscriptionid"`
	Region         string             `json:"region" bson:"region"`
	Selected       bool               `json:"selected" bson:"selected"`
	Name           string             `json:"name" bson:"name"`
}
