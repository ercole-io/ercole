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
	searchSegmentAdvisorsCmd := simpleAPIRequestCommand("search-segment-advisors",
		"Search current segment advisors",
		`search-segment-advisors search the most matching segment advisors to the arguments`,
		true, []apiOption{locationOption, environmentOption, sortingOptions, olderThanOptions}, true,
		"/hosts/technologies/oracle/databases/segment-advisors",
		"Failed to search segment advisors data: %v\n",
		"Failed to search segment advisors data(Status: %d): %s\n",
	)

	apiCmd.AddCommand(searchSegmentAdvisorsCmd)
}
