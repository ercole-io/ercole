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
	cmd := simpleSingleValueAPIRequestCommand("part-list",
		"Get Oracle/Database agreement parts list",
		`Get Oracle/Database agreement parts list`,
		false, false, false, false,
		"/settings/oracle/database/agreements-parts",
		"Failed to get Oracle/Database info: %v\n",
		"Failed to get Oracle/Database info(Status: %d): %s\n",
	)

	agreementOracleDatabaseCmd.AddCommand(cmd)
}
