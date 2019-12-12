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
	getDatabaseDataguardStatusStatsCmd := simpleAPIRequestCommand("dataguard-status",
		"Get databases dataguard status stats",
		`Get stats about the dataguard status of the databases`,
		false, false, false, true, true, false, false,
		"/stats/databases/dataguard-status",
		"Failed to get databases dataguard status stats: %v\n",
		"Failed to get databases dataguard status stats(Status: %d): %s\n",
	)

	statsDatabasesCmd.AddCommand(getDatabaseDataguardStatusStatsCmd)
}
