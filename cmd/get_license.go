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

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"

	"github.com/ercole-io/ercole/v2/utils"
)

// getLicenseCmd represents the get-license command
var getLicenseCmd = &cobra.Command{
	Use:   "get-license",
	Short: "Get a current host",
	Long:  `Get from api-service a current host`,
	Run: func(cmd *cobra.Command, args []string) {
		params := url.Values{}
		olderThanOptions.addParam(params)

		resp, err := http.Get(
			utils.NewAPIUrl(
				ercoleConfig.APIService.RemoteEndpoint,
				ercoleConfig.APIService.AuthenticationProvider.Username,
				ercoleConfig.APIService.AuthenticationProvider.Password,
				"/licenses/"+args[0],
				params,
			).String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get the license: %v\n", err)
			os.Exit(1)
		} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
			out, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			fmt.Fprintf(os.Stderr, "Failed to get the license(Status: %d): %s\n", resp.StatusCode, string(out))
			os.Exit(1)
		} else {
			out, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			var res interface{}
			err = json.Unmarshal(out, &res)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to unmarshal response body: %v (%s)\n", err, string(out))
				os.Exit(1)
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "    ")
			if err := enc.Encode(res); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
		}

	},
	Args: cobra.ExactArgs(1),
}

func init() {
	apiCmd.AddCommand(getLicenseCmd)

	olderThanOptions.addOption(getLicenseCmd)
}
