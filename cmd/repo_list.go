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
	"fmt"

	"github.com/spf13/cobra"
)

var filterInstalled bool
var filterUninstalled bool
var showOnlyNameColumn bool

func init() {
	repoListCmd := &cobra.Command{
		Use:   "list",
		Short: "List artifacts from repositories",
		Long:  `List artifacts from repositories`,
		Run: func(cmd *cobra.Command, args []string) {
			//Get the list of the repository
			index := readOrUpdateIndex()

			for _, f := range index {
				if filterInstalled && !filterUninstalled && !f.Installed {
					continue
				} else if filterUninstalled && !filterInstalled && f.Installed {
					continue
				}

				var installed rune
				if f.Installed {
					installed = 'i'
				} else {
					installed = '-'
				}

				if showOnlyNameColumn {
					fmt.Println(f.FullName())
				} else {
					fmt.Printf("%c\t%s\t%s\n", installed, f.ReleaseDate, f.FullName())
				}
			}
		},
	}
	repoListCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	repoListCmd.Flags().BoolVarP(&filterInstalled, "installed", "i", false, "Filter by installed artifacts")
	repoListCmd.Flags().BoolVarP(&filterUninstalled, "uninstalled", "u", false, "Filter by uninstalled artifacts")
	repoListCmd.Flags().BoolVarP(&showOnlyNameColumn, "name", "n", false, "Show only the name column")

	repoCmd.AddCommand(repoListCmd)
}
