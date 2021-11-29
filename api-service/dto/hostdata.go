// Copyright (c) 2021 Sorint.lab S.p.A.
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

package dto

import (
	"time"

	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ServerSchemaVersion contains the version of the schema
const ServerSchemaVersion int = 1

// HostData holds all informations about a host & services
type HostData struct {
	ID                      primitive.ObjectID            `json:"id" bson:"_id"`
	Archived                bool                          `json:"archived" bson:"archived"`
	CreatedAt               time.Time                     `json:"createdAt" bson:"createdAt"`
	ServerVersion           string                        `json:"serverVersion" bson:"serverVersion"`
	SchemaVersion           int                           `json:"schemaVersion" bson:"schemaVersion"`
	ServerSchemaVersion     int                           `json:"serverSchemaVersion" bson:"serverSchemaVersion"`
	Hostname                string                        `json:"hostname" bson:"hostname"`
	Location                string                        `json:"location" bson:"location"`
	Environment             string                        `json:"environment" bson:"environment"`
	AgentVersion            string                        `json:"agentVersion" bson:"agentVersion"`
	Cluster                 string                        `json:"cluster" bson:"cluster"`
	VirtualizationNode      string                        `json:"virtualizationNode" bson:"virtualizationNode"`
	Tags                    []string                      `json:"tags" bson:"tags"`
	Info                    model.Host                    `json:"info" bson:"info"`
	ClusterMembershipStatus model.ClusterMembershipStatus `json:"clusterMembershipStatus" bson:"clusterMembershipStatus"`
	Features                model.Features                `json:"features" bson:"features"`
	Filesystems             []model.Filesystem            `json:"filesystems" bson:"filesystems"`
	Clusters                []model.ClusterInfo           `json:"clusters" bson:"clusters"`
	Cloud                   model.Cloud                   `json:"cloud" bson:"cloud"`
	Errors                  []model.AgentError            `json:"errors" bson:"errors"`
	OtherInfo               map[string]interface{}        `json:"-" bson:"-"`
	Alerts                  []model.Alert                 `json:"alerts" bson:"alerts"`
	History                 []model.History               `json:"history" bson:"history"`
}
