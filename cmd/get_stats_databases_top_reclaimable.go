/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

func init() {
	getDatabaseTopReclaimableStatsCmd := simpleAPIRequestCommand("top-reclaimable",
		"Get databases top reclaimable stats",
		`Get stats about top reclaimable databases`,
		false, false, false, true, false, false, true,
		"/stats/databases/top-reclaimable",
		"Failed to get top reclaimable databases stats: %v\n",
		"Failed to get top reclaimable databases stats(Status: %d): %s\n",
	)

	statsDatabasesCmd.AddCommand(getDatabaseTopReclaimableStatsCmd)
}
