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

package model

// OracleDatabaseLicenseType holds informations about a single OracleDatabaseLicenseType
type OracleDatabaseLicenseType struct {
	ID              string   `json:"id" bson:"_id"`
	ItemDescription string   `json:"itemDescription" bson:"itemDescription"`
	Metric          string   `json:"metric" bson:"metric"`
	Cost            float64  `json:"cost" bson:"cost"`
	Aliases         []string `json:"aliases" bson:"aliases"`
	Option          bool     `json:"option" bson:"option"`
}

// OracleDatabaseContractPart known metrics
const (
	LicenseTypeMetricProcessorPerpetual     = "Processor Perpetual"
	LicenseTypeMetricNamedUserPlusPerpetual = "Named User Plus Perpetual"
	LicenseTypePartMetricStreamPerpetual    = "Stream Perpetual"
	LicenseTypeMetricComputerPerpetual      = "Computer Perpetual"
)

func GetAllLicenseTypeMetrics() []string {
	return []string{
		LicenseTypeMetricProcessorPerpetual,
		LicenseTypeMetricNamedUserPlusPerpetual,
		LicenseTypePartMetricStreamPerpetual,
		LicenseTypeMetricComputerPerpetual,
	}
}

// Factor of multiplication from a license PerCore to a license PerUser ("Named User Plus Perpetual")
const FactorNamedUser float64 = 25

func GetFactorByMetric(metric string) float64 {
	if metric == LicenseTypeMetricNamedUserPlusPerpetual {
		return FactorNamedUser
	}

	return 1
}
