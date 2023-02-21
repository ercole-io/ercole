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
		ID:              "FAKE-PART-NUMBER-001",
		ItemDescription: "SQL Server 2016 Standard",
		Edition:         "STD",
		Version:         "2016",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "FAKE-PART-NUMBER-002",
		ItemDescription: "SQL Server 2016 Enterprise",
		Edition:         "ENT",
		Version:         "2016",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "FAKE-PART-NUMBER-003",
		ItemDescription: "SQL Server 2014 Standard",
		Edition:         "STD",
		Version:         "2014",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "FAKE-PART-NUMBER-004",
		ItemDescription: "SQL Server 2014 Enterprise",
		Edition:         "ENT",
		Version:         "2014",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "FAKE-PART-NUMBER-005",
		ItemDescription: "SQL Server 2012 Standard",
		Edition:         "STD",
		Version:         "2012",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "FAKE-PART-NUMBER-006",
		ItemDescription: "SQL Server 2012 Enterprise",
		Edition:         "ENT",
		Version:         "2012",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "FAKE-PART-NUMBER-007",
		ItemDescription: "SQL Server 2008 R2 Standard",
		Edition:         "STD",
		Version:         "2008 R2",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "FAKE-PART-NUMBER-008",
		ItemDescription: "SQL Server 2008 R2 Enterprise",
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
