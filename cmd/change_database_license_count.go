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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ercole-io/ercole/utils"
	"github.com/spf13/cobra"
)

func init() {
	changeDatabaseLicenseCountCmd := &cobra.Command{
		Use:   "change-database-license-count",
		Short: "Change the value of a license of database",
		Long:  `Change the value of a license of database`,
		Run: func(cmd *cobra.Command, args []string) {
			//Argument length check
			if len(args)%4 != 0 {
				fmt.Fprintln(os.Stderr, "The argument list must be divisible by 4")
				os.Exit(1)
			}

			for i := 0; i < len(args); i += 4 {
				hostname := args[i]
				dbname := args[i+1]
				licenseName := args[i+2]
				newValue := args[i+3]

				req, _ := http.NewRequest("PUT", utils.NewAPIUrlNoParams(ercoleConfig.APIService.RemoteEndpoint,
					ercoleConfig.APIService.AuthenticationProvider.Username,
					ercoleConfig.APIService.AuthenticationProvider.Password,
					"/hosts/"+hostname+"/databases/"+dbname+"/licenses/"+licenseName,
				).String(), bytes.NewReader([]byte(newValue)))

				//Make the http request
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to change the value of the license %q to \"%s/%s\": %v\n", licenseName, hostname, dbname, err)
					os.Exit(1)
				} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
					out, _ := ioutil.ReadAll(resp.Body)
					defer resp.Body.Close()
					fmt.Fprintf(os.Stderr, "License: %q DB: \"%s/%s\" Status: %d Cause: %s\n", licenseName, hostname, dbname, resp.StatusCode, string(out))
					os.Exit(1)
				}
			}
		},
	}

	apiCmd.AddCommand(changeDatabaseLicenseCountCmd)
}
