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

package cmd

import (
	"errors"
	"fmt"
	"log"

	migration "github.com/ercole-io/ercole/v2/database-migration"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var insertSUCmd = &cobra.Command{
	Use:   "insert-su",
	Short: "Insert Super User",
	Long:  `Insert super user ercole`,
	Run: func(cmd *cobra.Command, args []string) {
		validate := func(input string) error {
			if len(input) < 6 {
				return errors.New("Password must have more than 6 characters")
			}
			return nil
		}

		prompt := promptui.Prompt{
			Label:       "Password",
			Validate:    validate,
			Mask:        '*',
			HideEntered: true,
		}

		password, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		prompt.Label = "Confirm"

		confirm, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		if password != confirm {
			log.Fatal(errors.New("Password is different"))
		}

		err = migration.InsertSuperUser(ercoleConfig.Mongodb, password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Super user created succesfully")
	},
}

func init() {
	rootCmd.AddCommand(insertSUCmd)
	insertSUCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable the verbosity")
}
