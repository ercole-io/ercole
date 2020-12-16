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

package cmd

import (
	"context"

	"github.com/spf13/cobra"

	migration "github.com/ercole-io/ercole/v2/database-migration"
	"github.com/ercole-io/ercole/v2/utils"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database",
	Long:  `Migrate the database`,
	Run: func(cmd *cobra.Command, args []string) {
		log := utils.NewLogger("SERV")

		cl := migration.ConnectToMongodb(log, ercoleConfig.Mongodb)
		migration.Migrate(log, cl.Database(ercoleConfig.Mongodb.DBName))
		cl.Disconnect(context.TODO())
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
