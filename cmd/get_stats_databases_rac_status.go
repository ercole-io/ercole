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
	getDatabaseRACStatusStatsCmd := simpleAPIRequestCommand("rac-status",
		"Get databases RAC status stats",
		`Get stats about the RAC status of the databases`,
		false, false, false, true, true, false, false,
		"/stats/databases/rac-status",
		"Failed to get databases RAC status stats: %v\n",
		"Failed to get databases RAC status stats(Status: %d): %s\n",
	)

	statsDatabasesCmd.AddCommand(getDatabaseRACStatusStatsCmd)
}
