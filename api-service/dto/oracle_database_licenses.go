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

package dto

// OracleDatabaseUsedLicenseSearchResponse dto
type OracleDatabaseUsedLicenseSearchResponse struct {
	Content  []OracleDatabaseUsedLicense `json:"content" bson:"content"`
	Metadata PagingMetadata              `json:"metadata" bson:"metadata"`
}

// OracleDatabaseUsedLicense dto
type OracleDatabaseUsedLicense struct {
	LicenseName  string  `json:"licenseName" bson:"licenseName"`
	DbName       string  `json:"dbName" bson:"dbName"`
	Hostname     string  `json:"hostname" bson:"hostname"`
	UsedLicenses float64 `json:"usedLicenses" bson:"usedLicenses"`
}

// OracleDatabaseLicenseUsage contains the information about usage of a license
type OracleDatabaseLicenseUsage struct {
	LicenseTypeID   string `json:"licenseTypeID" bson:"licenseTypeID"`
	ItemDescription string `json:"itemDescription" bson:"itemDescription"`
	Metric          string `json:"metric" bson:"metric"`

	Consumed   float64 `json:"consumed"`
	Covered    float64 `json:"covered"`
	Compliance float64 `json:"compliance"`
	Unlimited  bool    `json:"unlimited"`
}
