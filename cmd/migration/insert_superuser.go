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

package migration

import (
	"fmt"
	"log"
	"os"

	"github.com/ercole-io/ercole/v2/config"
	migration "github.com/ercole-io/ercole/v2/database-migration"
	"github.com/spf13/cobra"
)

// insertSuperUserCmd represents the insert Super User command
func NewInsertSuperUserCmd(conf *config.Configuration) *cobra.Command {
	return &cobra.Command{
		Use:   "insert-su",
		Short: "Insert Super User",
		Long:  `Insert super user ercole`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Insert super user's password please")
				os.Exit(1)
			} else {
				arg := args[0]
				err := migration.InsertSuperUser(conf.Mongodb, arg)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Super user created succesfully")
			}
		},
	}

}
