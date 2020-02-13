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

	"github.com/amreo/ercole-services/utils"
	"github.com/spf13/cobra"
)

func init() {
	tagDatabase := &cobra.Command{
		Use:   "tag-database",
		Short: "Add a tag to a database",
		Long:  `Add a tag to a database`,
		Run: func(cmd *cobra.Command, args []string) {
			//Argument lenght check
			if len(args)%3 != 0 {
				fmt.Fprintln(os.Stderr, "The argument list must be divisible by 3")
				os.Exit(1)
			}

			for i := 0; i < len(args); i += 3 {
				hostname := args[i]
				dbname := args[i+1]
				tagname := args[i+2]

				req, _ := http.NewRequest("POST", utils.NewAPIUrlNoParams(ercoleConfig.APIService.RemoteEndpoint,
					ercoleConfig.APIService.UserUsername,
					ercoleConfig.APIService.UserPassword,
					"/hosts/"+hostname+"/databases/"+dbname+"/tags",
				).String(), bytes.NewReader([]byte(tagname)))

				//Make the http request
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to set the tag %q to \"%s/%s\": %v\n", tagname, hostname, dbname, err)
					os.Exit(1)
				} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
					out, _ := ioutil.ReadAll(resp.Body)
					defer resp.Body.Close()
					fmt.Fprintf(os.Stderr, "Tag: %q DB: \"%s/%s\" Status: %d Cause: %s\n", tagname, hostname, dbname, resp.StatusCode, string(out))
					os.Exit(1)
				}
			}
		},
	}

	apiCmd.AddCommand(tagDatabase)
}
