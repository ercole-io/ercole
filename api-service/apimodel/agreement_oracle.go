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

package apimodel

import "go.mongodb.org/mongo-driver/bson/primitive"

// OracleDatabaseAgreementsAddRequest contains the informations needed to add new agreements
type OracleDatabaseAgreementsAddRequest struct {
	AgreementID     string   `json:"agreementID" bson:"agreementID"`
	PartsID         []string `json:"partsID" bson:"partsID"`
	CSI             string   `json:"csi" bson:"csi"`
	ReferenceNumber string   `json:"referenceNumber" bson:"referenceNumber"`
	Unlimited       bool     `json:"unlimited" bson:"unlimited"`
	Count           int      `json:"count" bson:"count"`
	CatchAll        bool     `json:"catchAll" bson:"catchAll"`
	Hosts           []string `json:"hosts" bson:"hosts"`
}

// OracleDatabaseAgreementsFE contains the informations about a agreement
type OracleDatabaseAgreementsFE struct {
	ID              primitive.ObjectID                         `json:"id" bson:"_id"`
	AgreementID     string                                     `json:"agreementID" bson:"agreementID"`
	PartID          string                                     `json:"partID" bson:"partID"`
	ItemDescription string                                     `json:"itemDescription" bson:"itemDescription"`
	Metrics         string                                     `json:"metrics" bson:"metrics"`
	CSI             string                                     `json:"csi" bson:"csi"`
	ReferenceNumber string                                     `json:"referenceNumber" bson:"referenceNumber"`
	Unlimited       bool                                       `json:"unlimited" bson:"unlimited"`
	Count           float64                                    `json:"count" bson:"count"`
	LicensesCount   float64                                    `json:"licensesCount" bson:"licensesCount"`
	UsersCount      float64                                    `json:"usersCount" bson:"usersCount"`
	AvailableCount  float64                                    `json:"availableCount" bson:"availableCount"`
	CatchAll        bool                                       `json:"catchAll" bson:"catchAll"`
	Hosts           []OracleDatabaseAgreementsAssociatedHostFE `json:"hosts" bson:"hosts"`
}

// OracleDatabaseAgreementsAssociatedHostFE contains the informations about a associated host in agreement
type OracleDatabaseAgreementsAssociatedHostFE struct {
	Hostname                  string  `json:"hostname" bson:"hostname"`
	CoveredLicensesCount      float64 `json:"coveredLicensesCount" bson:"coveredLicensesCount"`
	TotalCoveredLicensesCount float64 `json:"totalCoveredLicensesCount" bson:"totalCoveredLicensesCount"`
	ConsumedLicensesCount     float64 `json:"consumedLicensesCount" bson:"consumedLicensesCount"`
}

// SearchOracleDatabaseAgreementsFilters contains the filters used to get the list of Oracle/Database agreements
type SearchOracleDatabaseAgreementsFilters struct {
	AgreementID       string
	PartID            string
	ItemDescription   string
	CSI               string
	Metrics           string
	ReferenceNumber   string
	Unlimited         string //"" -> Ignore, "true" -> true, "false" -> false
	CatchAll          string //"" -> Ignore, "true" -> true, "false" -> false
	LicensesCountLTE  int
	LicensesCountGTE  int
	UsersCountLTE     int
	UsersCountGTE     int
	AvailableCountLTE int
	AvailableCountGTE int
}

// OracleDatabaseLicensingObjects contains the information about the objects that needed to be licensed by Oracle/Database agreements
type OracleDatabaseLicensingObjects struct {
	LicenseName   string  `json:"licenseName" bson:"licenseName"`
	Name          string  `json:"name" bson:"name"`
	Type          string  `json:"type" bson:"type"` //host, cluster
	Count         float64 `json:"count" bson:"count"`
	OriginalCount float64 `json:"originalCount" bson:"originalCount"`
}

// OracleDatabaseLicenseInfo contains the infromation about the licenses
type OracleDatabaseLicenseInfo struct {
	ID                   string                          `json:"id" bson:"_id"`
	Compliance           bool                            `json:"compliance" bson:"compliance"`
	CostPerProcessor     float64                         `json:"costPerProcessor" bson:"costPerProcessor"`
	TotalCost            float64                         `json:"totalCost" bson:"totalCost"`
	PaidCost             float64                         `json:"paidCost" bson:"paidCost"`
	Unlimited            bool                            `json:"unlimited" bson:"unlimited"`
	Count                float64                         `json:"count" bson:"count"`
	TotalCoveredLicenses float64                         `json:"totalCoveredLicenses" bson:"totalCoveredLicenses"`
	Used                 float64                         `json:"used" bson:"used"`
	Hosts                []OracleDatabaseLicenseInfoHost `json:"hosts" bson:"hosts"`
}

// OracleDatabaseLicenseInfoHost contains the information about the licensed host and the databases used by the host
type OracleDatabaseLicenseInfoHost struct {
	Hostname string   `json:"hostname" bson:"hostname"`
	DBNames  []string `json:"dbNames" bson:"dbNames"`
}
