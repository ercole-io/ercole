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

// OracleExadataComponent holds informations about a device in a exadata
type OracleExadataComponent struct {
	HostType          string                     `json:"hostType" bson:"hostType"`
	Hostname          string                     `json:"hostname" bson:"hostname"`
	CPUEnabled        int                        `json:"cpuEnabled" bson:"cpuEnabled"`
	TotalCPU          int                        `json:"totalCPU" bson:"totalCPU"`
	Memory            int                        `json:"memory" bson:"memory"`
	ImageVersion      string                     `json:"imageVersion" bson:"imageVersion"`
	Kernel            string                     `json:"kernel" bson:"kernel"`
	Model             string                     `json:"model" bson:"model"`
	FanUsed           int                        `json:"fanUsed" bson:"fanUsed"`
	FanTotal          int                        `json:"fanTotal" bson:"fanTotal"`
	PsuUsed           int                        `json:"psuUsed" bson:"psuUsed"`
	PsuTotal          int                        `json:"psuTotal" bson:"psuTotal"`
	MsStatus          string                     `json:"msStatus" bson:"msStatus"`
	RsStatus          string                     `json:"rsStatus" bson:"rsStatus"`
	CellServiceStatus string                     `json:"cellServiceStatus" bson:"cellServiceStatus"`
	SwVersion         string                     `json:"swVersion" bson:"swVersion"`
	VMs               []OracleExadataVM          `json:"vms,omitempty" bson:"vms"`
	StorageCells      []OracleExadataStorageCell `json:"storageCells,omitempty" bson:"storageCells"`
	Database          OracleExadataDatabase      `json:"database,omitempty" bson:"database"`
}

type OracleExadataDatabase struct {
	Type            string     `json:"type" bson:"type"`
	DbName          string     `json:"dbName" bson:"dbName"`
	DbID            int        `json:"dbID" bson:"dbID"`
	FlashCacheLimit int        `json:"flashCacheLimit" bson:"flashCacheLimit"`
	IormShare       int        `json:"iormShare" bson:"iormShare"`
	LastIOReq       *time.Time `json:"lastIOReq" bson:"lastIOReq"`
}
