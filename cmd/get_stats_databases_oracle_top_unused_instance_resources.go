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
	GetTopUnusedOracleDatabaseInstanceResourceStatsCmd := simpleAPIRequestCommand("top-unused-instance-resource",
		"Get top unused instance resource stats",
		`Get stats about top unused instance resource`,
		false, []apiOption{locationOption, limitOption, olderThanOptions}, false,
		"/hosts/technologies/oracle/databases/top-unused-instance-resource",
		"Failed to get top unused instance resource stats: %v\n",
		"Failed to get top unused instance resource stats(Status: %d): %s\n",
	)

	statsDatabasesOracleCmd.AddCommand(GetTopUnusedOracleDatabaseInstanceResourceStatsCmd)
}
