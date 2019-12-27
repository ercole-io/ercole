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

func init() {
	searchHostsCmd := simpleAPIRequestCommand("search-hosts",
		"Search current hosts",
		`search-hosts search the most matching hosts to the arguments`,
		true, []apiOption{fullOption, locationOption, environmentOption, sortingOptions},
		"/hosts",
		"Failed to search hosts data: %v\n",
		"Failed to search hosts data(Status: %d): %s\n",
	)

	apiCmd.AddCommand(searchHostsCmd)
}
