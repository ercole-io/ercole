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

package dto

import "github.com/ercole-io/ercole/v2/model"

// LicenseCompliance contains the information about usage of a license
type LicenseCompliance struct {
	LicenseTypeID   string  `json:"licenseTypeID" bson:"licenseTypeID"`
	ItemDescription string  `json:"itemDescription" bson:"itemDescription"`
	Metric          string  `json:"metric" bson:"metric"`
	Cost            float64 `json:"cost" bson:"cost"`

	Consumed   float64 `json:"consumed"`
	Covered    float64 `json:"covered"`
	Purchased  float64 `json:"purchased"`
	Compliance float64 `json:"compliance"`
	Unlimited  bool    `json:"unlimited"`
	Available  float64 `json:"available"`
}

func (l *LicenseCompliance) ToModel() model.LicenseCompliance {
	return model.LicenseCompliance{
		LicenseTypeID:   l.LicenseTypeID,
		ItemDescription: l.ItemDescription,
		Metric:          l.Metric,
		Cost:            l.Cost,
		Consumed:        l.Consumed,
		Covered:         l.Covered,
		Purchased:       l.Purchased,
		Compliance:      l.Compliance,
		Unlimited:       l.Unlimited,
		Available:       l.Available,
	}
}

type IgnoreLicenseRequest struct {
	Technology     string `json:"technology"`
	Hostname       string `json:"hostname"`
	DatabaseName   string `json:"databaseName"`
	LicenseTypeID  string `json:"licenseTypeID"`
	Ignored        bool   `json:"ignored"`
	IgnoredComment string `json:"ignoredComment"`
}

type IgnoreLicenseResponse struct {
	Updated []IgnoreLicenseRequest `json:"updated"`
	Error   []IgnoreLicenseRequest `json:"error"`
}
