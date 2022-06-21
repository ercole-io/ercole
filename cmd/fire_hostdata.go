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

func init() {
	rootCmd.AddCommand(fireHostDataCmd)
	fireHostDataCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable the verbosity")
	fireHostDataCmd.Flags().BoolVarP(&insecure, "insecure", "i", false, "Allow insecure server connections when using SSL")
}

func fireHostdata(filename string, content []byte) {
	importDataRequest(filename, content, "/hosts")
}
