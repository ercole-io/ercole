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
	"github.com/spf13/cobra"

	"github.com/ercole-io/ercole/v2/logger"
)

func init() {
	repoUpdateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update all artifacts installed",
		Long:  `Install the most recent version of all installed artifacts`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			index := readOrUpdateIndex(logger.NewLogger("REPO", logger.LogVerbosely(verbose)))

			updateCandidates := make(map[*ArtifactInfo]bool)

			for _, art := range index.artifacts {
				f := index.searchLatestArtifactByRepositoryAndName(art.Repository, art.Name)

				if !f.Installed {
					updateCandidates[f] = true
				}
			}

			for art := range updateCandidates {
				index.Install(art)
			}
		},
	}

	repoCmd.AddCommand(repoUpdateCmd)
}
