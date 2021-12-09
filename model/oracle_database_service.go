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

package model

import (
	"time"
)

// OracleDatabaseService holds information about an Oracle database service
type OracleDatabaseService struct {
	Name            *string    `json:"name" bson:"name,omitempty"`
	CreationDate    *time.Time `json:"creationDate" bson:"creationDate,omitempty"`
	FailoverMethod  *string    `json:"failoverMethod" bson:"failoverMethod,omitempty"`
	FailoverType    *string    `json:"failoverType" bson:"failoverType,omitempty"`
	FailoverRetries *int       `json:"failoverRetries" bson:"failoverRetries,omitempty"`
	FailoverDelay   *int       `json:"failoverDelay" bson:"failoverDelay,omitempty"`
	Enabled         *bool      `json:"enabled" bson:"enabled,omitempty"`
}
