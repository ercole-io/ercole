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

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var summary bool

// searchHostsCmd represents the getHosts command
var searchHostsCmd = &cobra.Command{
	Use:   "search-hosts",
	Short: "Search current hosts",
	Long:  `search-hosts search the most matching hosts to the arguments`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(
			fmt.Sprintf("http://%s:%s@%s/hosts?full=%v",
				ercoleConfig.APIService.UserUsername,
				ercoleConfig.APIService.UserPassword,
				ercoleConfig.APIService.RemoteEndpoint,
				!summary))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to search hosts data: %v\n", err)
			os.Exit(1)
		} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
			out, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			fmt.Fprintf(os.Stderr, "Failed to search hosts data(Status: %d): %s\n", resp.StatusCode, string(out))
			os.Exit(1)
		} else {
			out, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			var res []interface{}
			err = json.Unmarshal(out, &res)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to unmarshal response body: %v (%s)\n", err, string(out))
				os.Exit(1)
			}

			for _, item := range res {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "    ")
				enc.Encode(item)
			}
		}

	},
}

func init() {
	apiCmd.AddCommand(searchHostsCmd)
	searchHostsCmd.Flags().BoolVarP(&summary, "summary", "s", false, "Summary mode")
}
