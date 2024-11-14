// Copyright (c) 2021 Sorint.lab S.p.A.
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

package service

import (
	"time"

	"github.com/ercole-io/ercole/v2/api-service/dto"
)

// UpdateLicenseIgnoredField update license ignored field (true/false)
func (as *APIService) UpdateLicenseIgnoredField(hostname string, dbname string, licenseTypeID string, ignored bool, ignoredComment string) error {
	if err := as.Database.UpdateLicenseIgnoredField(hostname, dbname, licenseTypeID, ignored, ignoredComment); err != nil {
		return err
	}

	return nil
}

// CanMigrateLicense If an oracle database has an enterprise license and an option less than 3 months, then the db can be migrated.
// return if the database can be migrated or not.
func (as *APIService) CanMigrateLicense(hostname string, dbname string, filter dto.GlobalFilter) (bool, error) {
	isEnt := false

	usedLicenses, err := as.getOracleDatabasesUsedLicenses(hostname, filter)
	if err != nil {
		return false, err
	}

	for _, l := range usedLicenses {
		if l.Description == "Oracle Database Enterprise Edition" {
			isEnt = true
		}
	}

	if !isEnt {
		return false, nil
	}

	opts, err := as.Database.FindOracleOptionsByDbname(hostname, dbname)
	if err != nil {
		return false, err
	}

	for _, opt := range opts {
		if opt.LastUsageDate.After(time.Now().AddDate(0, -3, 0)) {
			return true, nil
		}
	}

	return false, nil
}
