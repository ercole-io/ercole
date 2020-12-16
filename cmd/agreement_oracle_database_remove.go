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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove an Oracle database agreement",
		Long:  `Remove an Oracle database agreement`,
		Run: func(cmd *cobra.Command, args []string) {
			//Argument length check
			for i := 0; i < len(args); i++ {
				id := args[0]
				req, _ := http.NewRequest("DELETE", utils.NewAPIUrlNoParams(ercoleConfig.APIService.RemoteEndpoint,
					ercoleConfig.APIService.AuthenticationProvider.Username,
					ercoleConfig.APIService.AuthenticationProvider.Password,
					"/agreements/oracle/database/"+id,
				).String(), nil)

				//Make the http request
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to remove the agreement %q: %v\n", id, err)
					os.Exit(1)
				} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
					out, _ := ioutil.ReadAll(resp.Body)
					defer resp.Body.Close()
					fmt.Fprintf(os.Stderr, "Agreement: %q Status: %d Cause: %s\n", id, resp.StatusCode, string(out))
					os.Exit(1)
				}
			}
		},
	}

	agreementOracleDatabaseCmd.AddCommand(cmd)
}
