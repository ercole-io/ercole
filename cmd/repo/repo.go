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
	"github.com/ercole-io/ercole/v2/config"
	"github.com/spf13/cobra"
)

// githubToken contains the token used to authorize GitHub requests (to avoid rate limit)
var githubToken string
var rebuildCache bool
var ercoleConfig *config.Configuration
var verbose bool

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage the internal repository",
	Long:  `Manage the internal repository. It requires to be run where the repository is installed`,
}

func NewRepoCmd(conf *config.Configuration) *cobra.Command {
	ercoleConfig = conf
	repoCmd.PersistentFlags().StringVarP(&githubToken, "github-token", "g", "", "Github token used to perform requests")
	repoCmd.PersistentFlags().BoolVar(&rebuildCache, "rebuild-cache", false, "Force the rebuild the cache")

	return repoCmd
}

func SetRepoVerbose(v bool) {
	verbose = v
}
