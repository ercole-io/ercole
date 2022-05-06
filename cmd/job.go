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
	"time"

	"github.com/ercole-io/ercole/v2/logger"
	thunderservice_database "github.com/ercole-io/ercole/v2/thunder-service/database"
	"github.com/ercole-io/ercole/v2/thunder-service/job"
	"github.com/spf13/cobra"
)

// showConfigCmd represents the showConfig command
var jobCmd = &cobra.Command{

	Use:   "job",
	Short: "Run job for Oracle Cloud retrieve objects number",
	Long:  `Run job for Oracle Cloud retrieve objects number`,
	Run:   OciDataRetrieve,
}

func OciDataRetrieve(cmd *cobra.Command, args []string) {
	log := logger.NewLogger("THUN", logger.LogVerbosely(verbose))

	db := &thunderservice_database.MongoDatabase{
		Config:  ercoleConfig,
		TimeNow: time.Now,
		Log:     log,
	}

	db.Init()

	j := &job.OciDataRetrieveJob{
		Database: db,
		TimeNow:  time.Now,
		Config:   ercoleConfig,
		Log:      log,
	}
	j.Run()
}

func init() {
	rootCmd.AddCommand(jobCmd)
}
