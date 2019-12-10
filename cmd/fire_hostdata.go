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

// fireHostDataCmd represents the fire-hostdata command
var fireHostDataCmd = &cobra.Command{
	Use:   "fire-hostdata",
	Short: "Fire hostdata",
	Long:  `Fire hostdata from the stdin or from the files in the args`,
	Run: func(cmd *cobra.Command, args []string) {
		dataToByFired := map[string][]byte{}

		//Load the data
		if len(args) == 0 {
			raw, _ := ioutil.ReadAll(os.Stdin)
			fireHostdata("stdin", raw)
		} else {
			for _, arg := range args {
				if raw, err := ioutil.ReadFile(arg); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to read the file %s: %v\n", arg, err)
					os.Exit(1)
				} else {
					fireHostdata(arg, raw)
				}
			}
		}

		//Send the hostdatas
		for file, data := range dataToByFired {
			fireHostdata(file, data)
		}
	},
}

func init() {
	rootCmd.AddCommand(fireHostDataCmd)
	fireHostDataCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable the verbosity")
}

func fireHostdata(filename string, content []byte) {
	resp, err := http.Post(
		utils.NewAPIUrlNoParams(ercoleConfig.DataService.RemoteEndpoint,
			ercoleConfig.DataService.AgentUsername,
			ercoleConfig.DataService.AgentPassword,
			"/hosts",
		).String(),
		"application/json", bytes.NewReader(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send hostdata from %s: %v\n", filename, err)
		os.Exit(1)
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		out, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		fmt.Fprintf(os.Stderr, "File: %s Status: %d Cause: %s\n", filename, resp.StatusCode, string(out))
		os.Exit(1)
	} else {
		if verbose {
			fmt.Printf("File: %s Status: %d\n", filename, resp.StatusCode)
		}
	}

}
