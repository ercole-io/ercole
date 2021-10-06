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

package repo

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ercole-io/ercole/v2/logger"
)

func init() {
	repoInfoCmd := &cobra.Command{
		Use:   "info [artifact...]",
		Short: "Get info about an artifact",
		Long:  `Get info about an artifact`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			index := readOrUpdateIndex(logger.NewLogger("REPO", logger.LogVerbosely(verbose)))

			for _, arg := range args {
				f := index.searchArtifactByArg(arg)
				if f == nil {
					fmt.Fprintf(os.Stderr, "The argument %q wasn't undestood\n", arg)
					os.Exit(1)
				}
				fmt.Println("Fullname:", f.FullName())
				fmt.Println("Repository:", f.Repository)
				fmt.Println("Name:", f.Name)
				fmt.Println("Version:", f.Version)
				fmt.Println("Release date:", f.ReleaseDate)
				fmt.Println("Filename:", f.Filename)
				fmt.Println("Operating system:", f.OperatingSystem)
				fmt.Println("Arch:", f.Arch)
				fmt.Println("Installed:", f.Installed)
				fmt.Println("UpstreamRepository.Type:", f.UpstreamRepository.Type)

				fmt.Println()
			}
		},
	}

	repoCmd.AddCommand(repoInfoCmd)
}
