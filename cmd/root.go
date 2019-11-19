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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/amreo/ercole-services/config"
	"github.com/spf13/cobra"
)

var ercoleConfig config.Configuration
var serverVersion = "latest"

// serveCmd represents the root command
var rootCmd = &cobra.Command{
	Use:     "ercole",
	Short:   "Ercole services & tools",
	Long:    `Ercole is a software for proactively managing software assets`,
	Version: serverVersion,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ercoleConfig = readConfig()
		ercoleConfig.Version = serverVersion
	},
}

// Execute executes the commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// readConfig read, parse and return a Configuration from the configuration file in config.json or /opt/ercole-hostdata-dataservice/config.json
// It also set the global Config with the read value
func readConfig() config.Configuration {
	var conf config.Configuration

	readSingleConfigFile("/opt/ercole/config.json", &conf)
	readSingleConfigFile("/etc/ercole.json", &conf)
	home, _ := os.UserHomeDir()
	readSingleConfigFile(home+"/.ercole.json", &conf)
	readSingleConfigFile("config.json", &conf)

	//Return the read configuration
	return conf
}

func readSingleConfigFile(filename string, conf *config.Configuration) {
	//Try to read the file
	var raw []byte
	var err error

	if raw, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	if err = json.Unmarshal(raw, conf); err != nil {
		log.Fatalf("Unable to parse configuration file %s (%s)", filename, err)
	}
}
