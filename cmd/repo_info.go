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

//Commands to be added
func init() {
	repoInfoCmd := &cobra.Command{
		Use:   "info",
		Short: "Get info about some files",
		Long:  `Get info about some files`,
		Run: func(cmd *cobra.Command, args []string) {
			//Get the list of the repository
			list, err := scanRepositories(verbose)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			//Calculate the fullname of wanted files in args
			var files []*fileInfo = make([]*fileInfo, len(args))
			for i, arg := range args {
				files[i], err = parseNameOfFile(arg, list)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}

			//Download and install the files
			for _, f := range files {
				fmt.Println("Repository:", f.repository)
				fmt.Println("Name:", f.repository)
				fmt.Println("Version:", f.version)
				fmt.Println("Release date:", f.releaseDate)
				fmt.Println("Filename:", f.filename)
				fmt.Println("Operating system:", f.operatingSystem)
				fmt.Println("Arch:", f.arch)
				fmt.Println("Installed:", f.installed)
				fmt.Println("InstalledCallback:", f.install == nil)

				fmt.Println()
			}
		},
	}
	repoInfoCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")

	repoCmd.AddCommand(repoInfoCmd)
}
