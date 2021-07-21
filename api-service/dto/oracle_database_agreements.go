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

// OracleDatabaseAgreementFE contains the informations about an AssociatedLicenseType in an Agreement for the frontend
type OracleDatabaseAgreementFE struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"` // ID of agreement - licenseType couple
	AgreementID string             `json:"agreementID" bson:"agreementID"`
	CSI         string             `json:"csi" bson:"csi"`

	// LicenseType

	LicenseTypeID   string `json:"licenseTypeID" bson:"licenseTypeID"`
	ItemDescription string `json:"itemDescription" bson:"itemDescription"`
	Metric          string `json:"metric" bson:"metric"`

	// Associated LicenseType

	ReferenceNumber string `json:"referenceNumber" bson:"referenceNumber"`
	Unlimited       bool   `json:"unlimited" bson:"unlimited"` // Or "ULA", "Unlimited License Agreement"

	CatchAll   bool                                      `json:"catchAll" bson:"catchAll"` //TODO Rename in basket
	Restricted bool                                      `json:"restricted" bson:"restricted"`
	Hosts      []OracleDatabaseAgreementAssociatedHostFE `json:"hosts" bson:"hosts"`

	LicensesPerCore float64 `json:"licensesPerCore" bson:"licensesPerCore"`
	LicensesPerUser float64 `json:"licensesPerUser" bson:"licensesPerUser"`

	// Value of licenses yet available to be assigned to hosts
	AvailableLicensesPerCore float64 `json:"availableLicensesPerCore" bson:"availableLicensesPerCore"`
	AvailableLicensesPerUser float64 `json:"availableLicensesPerUser" bson:"availableLicensesPerUser"`

	CoveredLicenses float64 `json:"-" bson:"-"`
}

// OracleDatabaseAgreementAssociatedHostFE contains the informations about an associated host in agreement
// If agreement is Named User, counts are in users
// TODO Rename: remove Count at the end of each name
type OracleDatabaseAgreementAssociatedHostFE struct {
	Hostname string `json:"hostname" bson:"hostname"`
	// Licenses which have been covered by agreement associated
	CoveredLicensesCount float64 `json:"coveredLicensesCount" bson:"coveredLicensesCount"`

	// Licenses covered by all agreements
	TotalCoveredLicensesCount float64 `json:"totalCoveredLicensesCount" bson:"totalCoveredLicensesCount"`
	// Licenses consumed (used) by this hostname, data from agents
	ConsumedLicensesCount float64 `json:"consumedLicensesCount" bson:"consumedLicensesCount"`
}

// GetOracleDatabaseAgreementsFilter contains the filter used to get the list of Oracle/Database agreements
type GetOracleDatabaseAgreementsFilter struct {
	AgreementID                 string
	LicenseTypeID               string
	ItemDescription             string
	CSI                         string
	Metric                      string
	ReferenceNumber             string
	Unlimited                   string //"" -> Ignore, "true" -> true, "false" -> false
	CatchAll                    string //"" -> Ignore, "true" -> true, "false" -> false //TODO Rename in Basket
	LicensesPerCoreLTE          int
	LicensesPerCoreGTE          int
	LicensesPerUserLTE          int
	LicensesPerUserGTE          int
	AvailableLicensesPerCoreLTE int
	AvailableLicensesPerCoreGTE int
	AvailableLicensesPerUserLTE int
	AvailableLicensesPerUserGTE int
}

func NewGetOracleDatabaseAgreementsFilter() GetOracleDatabaseAgreementsFilter {
	return GetOracleDatabaseAgreementsFilter{
		LicensesPerCoreLTE:          -1,
		LicensesPerCoreGTE:          -1,
		LicensesPerUserLTE:          -1,
		LicensesPerUserGTE:          -1,
		AvailableLicensesPerCoreLTE: -1,
		AvailableLicensesPerCoreGTE: -1,
		AvailableLicensesPerUserLTE: -1,
		AvailableLicensesPerUserGTE: -1,
	}
}

// HostUsingOracleDatabaseLicenses contains the information about the hosts that use licenses by Oracle/Database
type HostUsingOracleDatabaseLicenses struct {
	LicenseTypeID string `json:"licenseTypeID" bson:"licenseTypeID"`
	Name          string `json:"name" bson:"name"`
	//Type describe if it's an host or a cluster
	Type string `json:"type" bson:"type"`
	// TODO Rename in UncoveredLicenses // Licenses to be covered by agreement
	LicenseCount float64 `json:"licenseCount" bson:"licenseCount"`
	//TODO Rename in ConsumedLicensesCount // Original value of licenseCount (UncoveredLicenses), DO NOT EDIT!
	OriginalCount float64 `json:"originalCount" bson:"originalCount"`
}
