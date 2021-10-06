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
	"os"

	"github.com/ercole-io/ercole/v2/model"

	"github.com/spf13/cobra"

	"github.com/ercole-io/ercole/v2/utils"
)

var outputMode int

// getHostCmd represents the get-host command
var getHostPatchingFunctionCmd = &cobra.Command{
	Use:   "get-patching-function",
	Short: "Get the patching function of a host",
	Long:  `Get the patching function of a host`,
	Run: func(cmd *cobra.Command, args []string) {

		resp, err := http.Get(
			utils.NewAPIUrlNoParams(
				ercoleConfig.APIService.RemoteEndpoint,
				ercoleConfig.APIService.AuthenticationProvider.Username,
				ercoleConfig.APIService.AuthenticationProvider.Password,
				"/hosts/"+args[0]+"/patching-function",
			).String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get the hostdata: %v\n", err)
			os.Exit(1)
		} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
			out, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			fmt.Fprintf(os.Stderr, "Failed to get the pf (Status: %d): %s\n", resp.StatusCode, string(out))
			os.Exit(1)
		} else {
			out, _ := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			var res model.PatchingFunction
			err = json.Unmarshal(out, &res)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to unmarshal response body: %v (%s)\n", err, string(out))
				os.Exit(1)
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "    ")
			switch outputMode {
			case 0:
				enc.Encode(res)
			case 1:
				fmt.Print(res.Code)
			case 2:
				enc.Encode(res.Vars)
			}
		}

	},
	Args: cobra.ExactArgs(1),
}

func init() {
	apiCmd.AddCommand(getHostPatchingFunctionCmd)
	getHostPatchingFunctionCmd.Flags().IntVarP(&outputMode, "output-mode", "m", 0, "Change the output mode of the command. 0 to show all informations, 1 to show the code, 2 to show the data")
}
