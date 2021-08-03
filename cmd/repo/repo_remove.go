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
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/spf13/cobra"
)

func init() {
	repoRemoveCmd := &cobra.Command{
		Use:   "remove [artifact...]",
		Short: "Remove an artifact",
		Long:  `Remove an artifact`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			index := readOrUpdateIndex(logger.NewLogger("REPO", logger.LogVerbosely(verbose)))

			for _, arg := range args {
				f := index.searchArtifactByArg(arg)
				if f == nil {
					index.log.Warnf("Artifact %q wasn't undestood\n", arg)
					continue
				}

				if f.Installed {
					index.Uninstall(f)
				}
			}
		},
	}

	repoCmd.AddCommand(repoRemoveCmd)
}
