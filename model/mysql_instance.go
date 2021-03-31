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

type MySQLInstance struct {
	Name               string  `json:"name" bson:"name"`
	Version            string  `json:"version" bson:"version"`
	Edition            string  `json:"edition" bson:"edition"`
	Platform           string  `json:"platform" bson:"platform"`
	Architecture       string  `json:"architecture" bson:"architecture"`
	Engine             string  `json:"engine" bson:"engine"`
	RedoLogEnabled     string  `json:"redoLogEnabled" bson:"redoLogEnabled"`
	CharsetServer      string  `json:"charsetServer" bson:"charsetServer"`
	CharsetSystem      string  `json:"charsetSystem" bson:"charsetSystem"`
	PageSize           float64 `json:"pageSize" bson:"pageSize"` // in KB
	ThreadsConcurrency int     `json:"threadsConcurrency" bson:"threadsConcurrency"`
	BufferPoolSize     float64 `json:"bufferPoolSize" bson:"bufferPoolSize"` // in MB
	LogBufferSize      float64 `json:"logBufferSize" bson:"logBufferSize"`   // in MB
	SortBufferSize     float64 `json:"sortBufferSize" bson:"sortBufferSize"` // in MB
	ReadOnly           bool    `json:"readOnly" bson:"readOnly"`
	LogBin             bool    `json:"logBin" bson:"logBin"`
	HighAvailability   bool    `json:"highAvailability" bson:"highAvailability"`

	UUID       string   `json:"uuid" bson:"uuid"`
	IsMaster   bool     `json:"isMaster" bson:"isMaster"`
	SlaveUUIDs []string `json:"slaveUUIDs" bson:"slaveUUIDs"`
	IsSlave    bool     `json:"isSlave" bson:"isSlave"`
	MasterUUID *string  `json:"masterUUID" bson:"masterUUID"`

	Databases       []MySQLDatabase       `json:"databases" bson:"databases"`
	TableSchemas    []MySQLTableSchema    `json:"tableSchemas" bson:"tableSchemas"`
	SegmentAdvisors []MySQLSegmentAdvisor `json:"segmentAdvisors" bson:"segmentAdvisors"`
}

const (
	MySQLEditionCommunity  = "COMMUNITY"
	MySQLEditionEnterprise = "ENTERPRISE"
)
