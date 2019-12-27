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
	getOperatingSystemStatsCmd := simpleAPIRequestCommand("operating-system",
		"Get operating system stats",
		`Get stats about the info.operating_system field`,
		false, false, false, true, false, false, false, false,
		"/stats/operating-systems",
		"Failed to get operating system stats: %v\n",
		"Failed to get operating system stats(Status: %d): %s\n",
	)

	statsCmd.AddCommand(getOperatingSystemStatsCmd)
}
