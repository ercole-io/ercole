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
	"reflect"

	godynstruct "github.com/amreo/go-dyn-struct"
)

// OracleExadataComponent holds informations about a device in a exadata
type OracleExadataComponent struct {
	Hostname             string                   `json:"hostname" bson:"hostname"`
	ServerType           string                   `json:"serverType" bson:"serverType"`
	Model                string                   `json:"model" bson:"model"`
	SwVersion            string                   `json:"swVersion" bson:"swVersion"`
	SwReleaseDate        string                   `json:"swReleaseDate" bson:"swReleaseDate"`
	RunningCPUCount      *int                     `json:"runningCPUCount" bson:"runningCPUCount"`
	TotalCPUCount        *int                     `json:"totalCPUCount" bson:"totalCPUCount"`
	Memory               *int                     `json:"memory" bson:"memory"`
	Status               *string                  `json:"status" bson:"status"`
	RunningPowerSupply   *int                     `json:"runningPowerSupply" bson:"runningPowerSupply"`
	TotalPowerSupply     *int                     `json:"totalPowerSupply" bson:"totalPowerSupply"`
	PowerStatus          *string                  `json:"powerStatus" bson:"powerStatus"`
	RunningFanCount      *int                     `json:"runningFanCount" bson:"runningFanCount"`
	TotalFanCount        *int                     `json:"totalFanCount" bson:"totalFanCount"`
	FanStatus            *string                  `json:"fanStatus" bson:"fanStatus"`
	TempActual           *float64                 `json:"tempActual" bson:"tempActual"`
	TempStatus           *string                  `json:"tempStatus" bson:"tempStatus"`
	CellsrvServiceStatus *string                  `json:"cellsrvServiceStatus" bson:"cellsrvServiceStatus"`
	MsServiceStatus      *string                  `json:"msServiceStatus" bson:"msServiceStatus"`
	RsServiceStatus      *string                  `json:"rsServiceStatus" bson:"rsServiceStatus"`
	FlashcacheMode       *string                  `json:"flashcacheMode" bson:"flashcacheMode"`
	CellDisks            *[]OracleExadataCellDisk `json:"cellDisks" bson:"cellDisks"`
	OtherInfo            map[string]interface{}   `json:"-" bson:"-"`
}

// MarshalJSON return the JSON rappresentation of this
func (v OracleExadataComponent) MarshalJSON() ([]byte, error) {
	return godynstruct.DynMarshalJSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalJSON parse the JSON content in data and set the fields in v appropriately
func (v *OracleExadataComponent) UnmarshalJSON(data []byte) error {
	return godynstruct.DynUnmarshalJSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}

// MarshalBSON return the BSON rappresentation of this
func (v OracleExadataComponent) MarshalBSON() ([]byte, error) {
	return godynstruct.DynMarshalBSON(reflect.ValueOf(v), v.OtherInfo, "OtherInfo")
}

// UnmarshalBSON parse the BSON content in data and set the fields in v appropriately
func (v *OracleExadataComponent) UnmarshalBSON(data []byte) error {
	return godynstruct.DynUnmarshalBSON(data, reflect.ValueOf(v), &v.OtherInfo, "OtherInfo")
}
