// Copyright (c) 2023 Sorint.lab S.p.A.
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package migrations

import (
	"fmt"

	"github.com/ercole-io/ercole/v2/model"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	licenseTypes := make([]interface{}, 0, 2)

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "AU1GNGYG2DBH317JG",
		ItemDescription: "SQL Server 2016 Standard Core",
		Edition:         "STD",
		Version:         "2016",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "QNUXHXH1SKT0SXEAT",
		ItemDescription: "SQL Server 2016 Enterprise Core",
		Edition:         "ENT",
		Version:         "2016",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "0ZOQ1R0L1CV9X8KQ1",
		ItemDescription: "SQL Server 2014 Standard Core",
		Edition:         "STD",
		Version:         "2014",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "AKU4AOZ57AE7AMGV0",
		ItemDescription: "SQL Server 2014 Enterprise Core",
		Edition:         "ENT",
		Version:         "2014",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "UESCGA3LJ0YW8DM3Q",
		ItemDescription: "SQL Server 2012 Standard Core",
		Edition:         "STD",
		Version:         "2012",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "1F8UP2K6L5UEUSNT4",
		ItemDescription: "SQL Server 2012 Enterprise Core",
		Edition:         "ENT",
		Version:         "2012",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "N3KIE9EJZXR386Q2P",
		ItemDescription: "SQL Server 2008 R2 Standard Core",
		Edition:         "STD",
		Version:         "2008 R2",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "ZYI8M3I5P6KTCK45C",
		ItemDescription: "SQL Server 2008 R2 Enterprise Core",
		Edition:         "ENT",
		Version:         "2008 R2",
	})

	err := migrate.Register(func(db *mongo.Database) error {
		if err := addSqlServerDatabaseLicenseTypes(db, licenseTypes); err != nil {
			return err
		}

		return nil
	}, func(db *mongo.Database) error {
		return nil
	})

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}
