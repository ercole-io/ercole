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

func init() {
	getExadataTotalMemorySizeStatsCmd := simpleSingleValueAPIRequestCommand("total-memory-size",
		"Get exadata total memory size stats",
		`Get stats about the total size of memory of exadata`,
		false, true, true, true,
		"/hosts/technologies/oracle/exadata/total-memory-size",
		"Failed to get exadata total memory size stats: %v\n",
		"Failed to get exadata total memory size stats(Status: %d): %s\n",
	)

	statsExadataCmd.AddCommand(getExadataTotalMemorySizeStatsCmd)
}
