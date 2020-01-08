// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	repoListCmd := &cobra.Command{
		Use:   "list",
		Short: "List files from repositories",
		Long:  `List files from repositories`,
		Run: func(cmd *cobra.Command, args []string) {
			//Get the list of the repository
			list, err := scanRepositories(verbose)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			for _, f := range list {
				var installed rune
				if f.installed {
					installed = 'i'
				} else {
					installed = '-'
				}

				if f.version == "latest" {
					continue
				}

				fmt.Printf("%c\t%s\t%s\n", installed, f.releaseDate, f.getFullName())
			}
		},
	}
	repoListCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")

	repoCmd.AddCommand(repoListCmd)
}
