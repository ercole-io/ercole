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

package migration

import (
	"github.com/ercole-io/ercole/v2/config"
	migration "github.com/ercole-io/ercole/v2/database-migration"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
func NewMigrateCmd(conf *config.Configuration) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database",
		Long:  `Migrate the database to the latest known version`,
		Run: func(command *cobra.Command, args []string) {
			log := utils.NewLogger("SERV")

			err := migration.Migrate(conf.Mongodb)
			if err != nil {
				log.Fatal(err)
			}
			log.Info("Migrate successfully")
		},
	}
}
