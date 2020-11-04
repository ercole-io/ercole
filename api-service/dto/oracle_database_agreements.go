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

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//TODO Should I remove some of these?

// AssociatedPartInOracleDbAgreementRequest contains the informations needed to add or update an AssociatedPart
// in an OracleDatabaseAgreement
type AssociatedPartInOracleDbAgreementRequest struct {
	ID              string   `json:"id"`
	AgreementID     string   `json:"agreementID"`
	PartID          string   `json:"partID"`
	CSI             string   `json:"csi"`
	ReferenceNumber string   `json:"referenceNumber"`
	Unlimited       bool     `json:"unlimited"`
	Count           int      `json:"count"`
	CatchAll        bool     `json:"catchAll"`
	Hosts           []string `json:"hosts"`
}

// OracleDatabaseAgreementFE contains the informations about an AssociatedPart in an Agreement for the frontend
type OracleDatabaseAgreementFE struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"` // ID of agreement - part couple
	AgreementID string             `json:"agreementID" bson:"agreementID"`
	CSI         string             `json:"csi" bson:"csi"`

	// Part
	PartID          string `json:"partID" bson:"partID"`
	ItemDescription string `json:"itemDescription" bson:"itemDescription"`
	Metric          string `json:"metric" bson:"metric"`

	// Associated Part

	ReferenceNumber string `json:"referenceNumber" bson:"referenceNumber"`
	Unlimited       bool   `json:"unlimited" bson:"unlimited"` // Or "ULA"
	// Number of licenses or users set by user.
	// If agreement is Named User, Count number express users, not licenses
	Count float64 `json:"count" bson:"count"`

	CatchAll bool                                      `json:"catchAll" bson:"catchAll"`
	Hosts    []OracleDatabaseAgreementAssociatedHostFE `json:"hosts" bson:"hosts"`

	// Value of licenses/users yet available to be assigned to hosts
	AvailableCount float64 `json:"availableCount" bson:"availableCount"`
	// Number of licenses
	LicensesCount float64 `json:"licensesCount" bson:"licensesCount"`
	// Number of users
	UsersCount float64 `json:"usersCount" bson:"usersCount"`
}

// OracleDatabaseAgreementAssociatedHostFE contains the informations about an associated host in agreement
// If agreement is Named User, counts are in users
type OracleDatabaseAgreementAssociatedHostFE struct {
	Hostname string `json:"hostname" bson:"hostname"`
	// Licenses which have been covered by agreement associated
	CoveredLicensesCount float64 `json:"coveredLicensesCount" bson:"coveredLicensesCount"`

	// Licenses covered by all agreements
	TotalCoveredLicensesCount float64 `json:"totalCoveredLicensesCount" bson:"totalCoveredLicensesCount"`
	// Licenses consumed (used) by this hostname, data from agents
	ConsumedLicensesCount float64 `json:"consumedLicensesCount" bson:"consumedLicensesCount"`
}

// SearchOracleDatabaseAgreementsFilter contains the filter used to get the list of Oracle/Database agreements
type SearchOracleDatabaseAgreementsFilter struct {
	AgreementID       string
	PartID            string
	ItemDescription   string
	CSI               string
	Metric            string
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

// HostUsingOracleDatabaseLicenses contains the information about the hosts that use licenses by Oracle/Database
type HostUsingOracleDatabaseLicenses struct {
	LicenseName string `json:"licenseName" bson:"licenseName"`
	//TODO Add partID and use it instead of LicenseName
	Name string `json:"name" bson:"name"`
	//Type describe if it's an host or a cluster
	Type string `json:"type" bson:"type"`
	// TODO Licenses to be covered by agreement... Uncovered licenses???
	LicenseCount float64 `json:"licenseCount" bson:"licenseCount"`
	//TODO original value of licenseCount, but not to be edited?!? Equals to ConsumedLicensesCount
	OriginalCount float64 `json:"originalCount" bson:"originalCount"`
}

// OracleDatabaseLicenseUsageInfo contains the information about usage of a license
type OracleDatabaseLicenseUsageInfo struct {
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
