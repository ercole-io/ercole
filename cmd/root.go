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
	"os"
	"strings"

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/utils"
	"github.com/spf13/cobra"
)

var ercoleConfig config.Configuration
var serverVersion = "latest"
var extraConfigFile string

// serveCmd represents the root command
var rootCmd = &cobra.Command{
	Use:     "ercole",
	Short:   "Ercole services & tools",
	Long:    `Ercole is a software for proactively managing software assets`,
	Version: serverVersion,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log := utils.NewLogger("CONF")
		extraConfigFile = strings.TrimSpace(extraConfigFile)

		if len(extraConfigFile) > 0 && !fileExists(extraConfigFile) {
			log.Fatalf("Configuration file not found: %s", extraConfigFile)
		}

		ercoleConfig = config.ReadConfig(log, extraConfigFile)
		ercoleConfig.Version = serverVersion
	},
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Execute executes the commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&extraConfigFile, "config", "c", "", "Configuration file")
}
