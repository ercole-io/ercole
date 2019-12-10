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

	searchClustersCmd := simpleAPIRequestCommand("search-clusters",
		"Search current clusters",
		`search-clusters search the most matching clusters to the arguments`,
		true, true, false, true, true, true,
		"/clusters",
		"Failed to search clusters data: %v\n",
		"Failed to search clusters data(Status: %d): %s\n",
	)

	apiCmd.AddCommand(searchClustersCmd)
}
