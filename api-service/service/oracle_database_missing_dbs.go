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

// Package service is a package that provides methods for querying data

package service

import (
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
)

func (as *APIService) GetMissingDatabases() ([]dto.OracleDatabaseMissingDbs, error) {
	return as.Database.GetMissingDatabases()
}

func (as *APIService) GetMissingDatabasesByHostname(hostname string) ([]model.MissingDatabase, error) {
	return as.Database.GetMissingDatabasesByHostname(hostname)
}

func (as *APIService) UpdateMissingDatabaseIgnoredField(hostname string, dbname string, ignored bool, ignoredComment string) error {
	if err := as.Database.UpdateMissingDatabaseIgnoredField(hostname, dbname, ignored, ignoredComment); err != nil {
		return err
	}

	return nil
}
