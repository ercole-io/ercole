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
	"os"

	"github.com/spf13/cobra"
)

func init() {
	repoInstallCmd := &cobra.Command{
		Use:   "install",
		Short: "Install an artifact",
		Long:  `Install an artifact`,
		Run: func(cmd *cobra.Command, args []string) {
			//Get the list of the repository
			index := readOrUpdateIndex()

			if len(args) == 0 {
				fmt.Fprintf(os.Stderr, "Please specify an artifact\n")
			}

			//Search the artifact and install it for every artifact
			for _, arg := range args {
				f := index.SearchArtifactByArg(arg)
				if f == nil {
					fmt.Fprintf(os.Stderr, "The argument %q wasn't undestood\n", arg)
					os.Exit(1)
				}

				if !f.Installed {
					f.Install(f)
				}
			}
		},
	}
	repoInstallCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")

	repoCmd.AddCommand(repoInstallCmd)
}
