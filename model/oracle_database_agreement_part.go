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

// OracleDatabaseAgreementPart holds informations about a sigle OracleDatabaseAgreementPart
type OracleDatabaseAgreementPart struct {
	PartID          string   `json:"partID" bson:"partID"`
	ItemDescription string   `json:"itemDescription" bson:"itemDescription"`
	Metrics         string   `json:"metrics" bson:"metrics"`
	Cost            float64  `json:"cost" bson:"cost"`
	Aliases         []string `json:"aliases" bson:"aliases"`
}

// OracleDatabaseAgreementPart known metrics
const (
	AgreementPartMetricProcessorPerpetual     = "Processor Perpetual"
	AgreementPartMetricNamedUserPlusPerpetual = "Named User Plus Perpetual"
	AgreementPartMetricStreamPerpetual        = "Stream Perpetual"
	AgreementPartMetricComputerPerpetual      = "Computer Perpetual"
)
