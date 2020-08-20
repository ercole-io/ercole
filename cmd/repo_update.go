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
	"github.com/ercole-io/ercole/cmd/repo"
	"github.com/spf13/cobra"
)

func init() {
	repoUpdateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update all artifacts installed",
		Long:  `Install the most recent version of all installed artifacts`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			//Get the list of the repository
			index := readOrUpdateIndex()

			updateCandidates := make(map[*repo.ArtifactInfo]bool)

			//Search the artifact and install it for every artifact
			for _, art := range index {
				f := index.SearchArtifactByRepositoryAndName(art.Repository, art.Name)

				if !f.Installed {
					updateCandidates[f] = true
				}
			}

			//Install all updateCandidates
			for art := range updateCandidates {
				art.Install(art)
			}
		},
	}
	repoUpdateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")

	repoCmd.AddCommand(repoUpdateCmd)
}
