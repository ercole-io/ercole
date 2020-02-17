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
	getDatabaseDataguardStatusStatsCmd := simpleAPIRequestCommand("dataguard-status",
		"Get databases dataguard status stats",
		`Get stats about the dataguard status of the databases`,
		false, []apiOption{locationOption, environmentOption, olderThanOptions}, false,
		"/stats/databases/dataguard-status",
		"Failed to get databases dataguard status stats: %v\n",
		"Failed to get databases dataguard status stats(Status: %d): %s\n",
	)

	statsDatabasesCmd.AddCommand(getDatabaseDataguardStatusStatsCmd)
}
