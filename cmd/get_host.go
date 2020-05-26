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

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/ercole-io/ercole/utils"
	"github.com/spf13/cobra"
)

// getHostCmd represents the get-host command
var getHostCmd = &cobra.Command{
	Use:   "get-host",
	Short: "Get a current host",
	Long:  `Get from api-service a current host`,
	Run: func(cmd *cobra.Command, args []string) {
		params := url.Values{}
		olderThanOptions.addParam(params)

		req, _ := http.NewRequest("GET", utils.NewAPIUrl(
			ercoleConfig.APIService.RemoteEndpoint,
			ercoleConfig.APIService.AuthenticationProvider.Username,
			ercoleConfig.APIService.AuthenticationProvider.Password,
			"/hosts/"+args[0],
			params,
		).String(), bytes.NewReader([]byte{}))

		switch outputFormat {
		case "json":
			outputFormat = "application/json"
		case "mongohostdata", "mhd":
			outputFormat = "application/vnd.ercole.mongohostdata+json"
		}
		req.Header.Set("Accept", outputFormat)

		resp, err := http.DefaultClient.Do(req)
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

	olderThanOptions.addOption(getHostCmd)
	getHostCmd.Flags().StringVarP(&outputFormat, "format", "f", "application/json", "Output format")
}
