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

// OracleDatabaseAgreementFE contains the informations about
type OracleDatabaseAgreementFE struct {
	AgreementID string `json:"agreementID" bson:"agreementID"`
	CSI         string `json:"csi" bson:"csi"`

	// Part
	PartID          string `json:"partID" bson:"partID"`
	ItemDescription string `json:"itemDescription" bson:"itemDescription"`
	Metric          string `json:"metric" bson:"metric"`

	// Associated Part
	ID              string                                    `json:"id" bson:"_id"`
	ReferenceNumber string                                    `json:"referenceNumber" bson:"referenceNumber"`
	Unlimited       bool                                      `json:"unlimited" bson:"unlimited"`
	Count           float64                                   `json:"count" bson:"count"`
	CatchAll        bool                                      `json:"catchAll" bson:"catchAll"`
	Hosts           []OracleDatabaseAgreementAssociatedHostFE `json:"hosts" bson:"hosts"`

	AvailableCount float64 `json:"availableCount" bson:"availableCount"`
	// availableCount and metric is model.AgreementPartMetricProcessorPerpetual
	LicensesCount float64 `json:"licensesCount" bson:"licensesCount"`
	// availableCount and metric is model.AgreementPartMetricNamedUserPlusPerpetual
	UsersCount float64 `json:"usersCount" bson:"usersCount"`
}

// OracleDatabaseAgreementAssociatedHostFE contains the informations about a associated host in agreement
type OracleDatabaseAgreementAssociatedHostFE struct {
	Hostname                  string  `json:"hostname" bson:"hostname"`
	CoveredLicensesCount      float64 `json:"coveredLicensesCount" bson:"coveredLicensesCount"`
	TotalCoveredLicensesCount float64 `json:"totalCoveredLicensesCount" bson:"totalCoveredLicensesCount"`
	ConsumedLicensesCount     float64 `json:"consumedLicensesCount" bson:"consumedLicensesCount"`
}

// SearchOracleDatabaseAgreementsFilter contains the filter used to get the list of Oracle/Database agreements
type SearchOracleDatabaseAgreementsFilter struct {
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

// HostUsingOracleDatabaseLicenses contains the information about the objects that use licenses by Oracle/Database
type HostUsingOracleDatabaseLicenses struct {
	LicenseName string `json:"licenseName" bson:"licenseName"`
	Name        string `json:"name" bson:"name"`
	//Type describe if it's an host or a cluster
	Type          string  `json:"type" bson:"type"`
	LicenseCount  float64 `json:"licenseCount" bson:"licenseCount"`
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
