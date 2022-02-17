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
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// fireHostDataCmd represents the fire-hostdata command
var fireHostDataCmd = &cobra.Command{
	Use:   "fire-hostdata",
	Short: "Fire hostdata",
	Long:  `Fire hostdata from the stdin or from the files in the args`,
	Run: func(cmd *cobra.Command, args []string) {
		//Load the data
		if len(args) == 0 {
			raw, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				os.Exit(1)
			}
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
	},
}

var insecure bool

func init() {
	rootCmd.AddCommand(fireHostDataCmd)
	fireHostDataCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable the verbosity")
	fireHostDataCmd.Flags().BoolVarP(&insecure, "insecure", "i", false, "Allow insecure server connections when using SSL")
}

func fireHostdata(filename string, content []byte) {
	client := &http.Client{}
	tr := http.DefaultTransport.(*http.Transport).Clone()

	if insecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client.Transport = tr

	req, err := http.NewRequest("POST", ercoleConfig.DataService.RemoteEndpoint+"/hosts", bytes.NewReader(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create request: %s", err)
		os.Exit(1)
	}

	src := ercoleConfig.DataService.AgentUsername + ":" + ercoleConfig.DataService.AgentPassword
	bearer := "Basic " + base64.StdEncoding.EncodeToString([]byte(src))
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send hostdata from %s: %v\n", filename, err)
		os.Exit(1)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		fmt.Fprintf(os.Stderr, "File: %s Status: %d Cause: %s\n", filename, resp.StatusCode, string(out))
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("File: %s Status: %d\n", filename, resp.StatusCode)
	}
}
