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
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	repoFetchCmd := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch files from repositories",
		Long:  `Fetch files from repositories`,
		Run: func(cmd *cobra.Command, args []string) {
			//Fetch the list of files
			list, err := scanRepositories(verbose)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			_ = list

			//Calculate the fullname of wanted files in args
			var files []*fileInfo = make([]*fileInfo, len(args))
			if verbose {
				fmt.Print("Fetching")
			}
			for i, arg := range args {
				files[i], err = parseNameOfFile(arg, list)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if verbose {
					fmt.Printf(" %s", files[i].getFullName())
				}
			}
			if verbose {
				fmt.Println()
			}

			//Download and install the files
			for _, f := range files {
				if f.installed {
					continue
				}
				if verbose {
					fmt.Printf("Downloading %s to %s.\n", f.getFullName(), filepath.Join(ercoleConfig.RepoService.DistributedFiles, f.filename))
				}
				if err = f.download(filepath.Join(ercoleConfig.RepoService.DistributedFiles, f.filename)); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}

				if f.install != nil {
					if err = f.install(filepath.Join(ercoleConfig.RepoService.DistributedFiles, f.filename)); err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
				}
			}
		},
	}
	repoFetchCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")

	repoCmd.AddCommand(repoFetchCmd)
}
