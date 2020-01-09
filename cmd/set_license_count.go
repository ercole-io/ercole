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
	"strconv"

	"github.com/amreo/ercole-services/utils"
	"github.com/spf13/cobra"
)

func init() {
	setLicenseCountCmd := &cobra.Command{
		Use:   "set-license-count",
		Short: "Set the count of a license",
		Long:  `Set the count of a license`,
		Run: func(cmd *cobra.Command, args []string) {
			//Argument lenght check
			if len(args)%2 != 0 {
				fmt.Fprintln(os.Stderr, "The argument list must be even")
				os.Exit(1)
			}

			for i := 0; i < len(args); i += 2 {
				name := args[i]
				_, err := strconv.Atoi(args[i+1])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to parse the count %q, err: %v\n", args[i+1], err)
					os.Exit(1)
				}

				req, _ := http.NewRequest("PUT", utils.NewAPIUrlNoParams(ercoleConfig.APIService.RemoteEndpoint,
					ercoleConfig.APIService.UserUsername,
					ercoleConfig.APIService.UserPassword,
					"/licenses/"+name,
				).String(), bytes.NewReader([]byte(args[i+1])))

				//Make the http request
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to set count of license %q: %v\n", name, err)
					os.Exit(1)
				} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
					out, _ := ioutil.ReadAll(resp.Body)
					defer resp.Body.Close()
					fmt.Fprintf(os.Stderr, "License: %s Status: %d Cause: %s\n", name, resp.StatusCode, string(out))
					os.Exit(1)
				}
			}
		},
	}

	apiCmd.AddCommand(setLicenseCountCmd)
}
