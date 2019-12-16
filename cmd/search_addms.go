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
	searchAddmsCmd := simpleAPIRequestCommand("search-addms",
		"Search current addms",
		`search-addms search the most matching addms to the arguments`,
		true, false, false, true, true, true, false,
		"/addms",
		"Failed to search addms data: %v\n",
		"Failed to search addms data(Status: %d): %s\n",
	)

	apiCmd.AddCommand(searchAddmsCmd)
}
