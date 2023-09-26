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

import "time"

const (
	KVM_HOST   = "KVM_HOST"
	DOM0       = "DOM0"
	BARE_METAL = "BARE_METAL"
	VM_KVM     = "VM_KVM"
	VM_XEN     = "VM_XEN"
)

// OracleExadataInstance holds specific informations about a exadata.
type OracleExadataInstance struct {
	Hostname    string                   `json:"hostname" bson:"hostname"`
	Environment string                   `json:"environment" bson:"environment"`
	Location    string                   `json:"location" bson:"location"`
	CreatedAt   time.Time                `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time                `json:"updateAt" bson:"updateAt"`
	RackID      string                   `json:"rackID" bson:"rackID"`
	Components  []OracleExadataComponent `json:"components" bson:"components"`
}
