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

	"github.com/amreo/ercole-services/utils"
	"github.com/spf13/cobra"
)

// getHostCmd represents the getHost command
var getHostCmd = &cobra.Command{
	Use:   "get-host",
	Short: "Get a current host",
	Long:  `Get from api-service a current host`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(
			utils.NewAPIUrlNoParams(
				ercoleConfig.APIService.RemoteEndpoint,
				ercoleConfig.APIService.UserUsername,
				ercoleConfig.APIService.UserPassword,
				"/hosts/"+args[0],
			).String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get the hostdata: %v\n", err)
			os.Exit(1)
		} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
			out, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			fmt.Fprintf(os.Stderr, "Failed to get the hostdata(Status: %d): %s\n", resp.StatusCode, string(out))
			os.Exit(1)
		} else {
			out, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			var res interface{}
			err = json.Unmarshal(out, &res)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to unmarshal response body: %v (%s)\n", err, string(out))
				os.Exit(1)
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "    ")
			enc.Encode(res)
		}

	},
	Args: cobra.ExactArgs(1),
}

func init() {
	apiCmd.AddCommand(getHostCmd)
}
