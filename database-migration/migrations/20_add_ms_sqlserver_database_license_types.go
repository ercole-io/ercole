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
		ID:              "LYQ7J4MTPCSOIKFDB",
		ItemDescription: "SQL Server 2017 Standard Core",
		Edition:         "STD",
		Version:         "2017",
	})

	licenseTypes = append(licenseTypes, model.SqlServerDatabaseLicenseType{
		ID:              "QQSWRN2BT8VM1NJIV",
		ItemDescription: "SQL Server 2017 Enterprise Core",
		Edition:         "ENT",
		Version:         "2017",
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
