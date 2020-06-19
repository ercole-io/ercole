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
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"

	migration "github.com/ercole-io/ercole/database-migration"
	"github.com/ercole-io/ercole/utils"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database",
	Long:  `Migrate the database`,
	Run: func(cmd *cobra.Command, args []string) {
		log := utils.NewLogger("SERV")

		if ercoleConfig.ResourceFilePath == "RESOURCES_NOT_FOUND" {
			log.Fatal("The directory for resources wasn't found")
		}

		//Read initial licenses list
		content, err := ioutil.ReadFile(ercoleConfig.ResourceFilePath + "/initial_oracle_licenses_list.txt")
		if err != nil {
			log.Fatalf("Cannot read the licenses list: %v\n", err)
		}
		lines := strings.Split(string(content), "\n")

		//Migrate
		cl := migration.ConnectToMongodb(log, ercoleConfig.Mongodb)
		migration.Migrate(log, cl.Database(ercoleConfig.Mongodb.DBName), lines)
		cl.Disconnect(context.TODO())
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
