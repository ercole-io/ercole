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
	getTechnologyChangeChartCmd := simpleChartRequestCommandSingleValue("changes",
		"Get the chart about the changes",
		`changes return the chart about the changes`,
		false, []apiOption{fromOption, locationOption, environmentOption, olderThanOptions}, false,
		"/technologies/changes",
		"Failed get the chart about the changes: %v\n",
		"Failed get the chart about the changes: %d): %s\n",
	)

	chartCmd.AddCommand(getTechnologyChangeChartCmd)
}
