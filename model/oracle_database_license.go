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

// OracleDatabaseLicense holds information about an Oracle database license
type OracleDatabaseLicense struct {
	LicenseTypeID string  `json:"licenseTypeID" bson:"licenseTypeID"`
	Name          string  `json:"name" bson:"name"`
	Count         float64 `json:"count" bson:"count"`
	Ignored       bool    `json:"ignored" bson:"ignored"`
}

// DiffFeature status of each feature
const (
	// DiffFeatureInactive is used when the feature changes from (0/-) to 0
	DiffFeatureInactive int = -2
	// DiffFeatureDeactivated is used when the feature changes from 1 to (0/-)
	DiffFeatureDeactivated int = -1
	// DiffFeatureMissing is used when a feature is missing in the diff
	DiffFeatureMissing int = 0
	// DiffFeatureActivated is used when the feature changes from (0/-) to 1
	DiffFeatureActivated int = 1
	// DiffFeatureActive is used when the feature changes from 1 to 1
	DiffFeatureActive int = 2
)

// DiffLicenses return a map that contains the difference of status between the oldLicenses and newLicenses
func DiffLicenses(oldLicenses, newLicenses []OracleDatabaseLicense) map[string]int {
	result := make(map[string]int)

	// Add the features to the result assuming that the all new features are inactive
	for _, license := range oldLicenses {
		if license.Count > 0 {
			result[license.LicenseTypeID] = DiffFeatureDeactivated
		} else {
			result[license.LicenseTypeID] = DiffFeatureInactive
		}
	}

	for _, license := range newLicenses {
		id := license.LicenseTypeID

		if (result[id] == DiffFeatureInactive || result[id] == DiffFeatureMissing) && license.Count <= 0 {
			result[id] = DiffFeatureInactive
			continue
		}

		if (result[id] == DiffFeatureDeactivated) && license.Count <= 0 {
			result[id] = DiffFeatureDeactivated
			continue
		}

		if (result[id] == DiffFeatureInactive || result[id] == DiffFeatureMissing) && license.Count > 0 {
			result[id] = DiffFeatureActivated
			continue
		}

		if (result[id] == DiffFeatureDeactivated) && license.Count > 0 {
			result[id] = DiffFeatureActive
			continue
		}
	}

	return result
}
