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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/ercole-io/ercole/v2/utils"
)

func init() {
	ackAlertCmd := &cobra.Command{
		Use:   "ack-alert [id...]",
		Short: "Ack an alert",
		Long:  `Ack an alert`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			req, err := http.NewRequest("POST",
				utils.NewAPIUrlNoParams(
					ercoleConfig.APIService.RemoteEndpoint,
					ercoleConfig.APIService.AuthenticationProvider.Username,
					ercoleConfig.APIService.AuthenticationProvider.Password,
					"/alerts/ack",
				).String(),
				bytes.NewReader([]byte(utils.ToJSON(args))),
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to acknowledge alerts %q: %v\n", args, err)
				os.Exit(1)
			} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
				out, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err)
					os.Exit(1)
				}
				defer resp.Body.Close()
				fmt.Fprintf(os.Stderr, "Alerts: %q Status: %d Cause: %s\n", args, resp.StatusCode, string(out))
				os.Exit(1)
			}
		},
	}

	apiCmd.AddCommand(ackAlertCmd)
}
