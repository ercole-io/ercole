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
	getExadataAverageStorageUsageStatsCmd := simpleSingleValueAPIRequestCommand("average-storage-usage",
		"Get exadata average storage usage stats",
		`Get stats about the average usage of storage of exadata`,
		false, true, true, true,
		"/stats/exadata/average-storage-usage",
		"Failed to get exadata average storage usage stats: %v\n",
		"Failed to get exadata average storage usage stats(Status: %d): %s\n",
	)

	statsExadataCmd.AddCommand(getExadataAverageStorageUsageStatsCmd)
}
