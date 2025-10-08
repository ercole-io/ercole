// Copyright (c) 2025 Sorint.lab S.p.A.
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
package job

import (
	"time"

	"github.com/ercole-io/ercole/v2/alert-service/database"
	"github.com/ercole-io/ercole/v2/alert-service/emailer"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
)

type SimulatedHostAlertJob struct {
	Database database.MongoDatabaseInterface
	Config   config.Configuration
	Log      logger.Logger
	Emailer  emailer.Emailer
}

func (j *SimulatedHostAlertJob) Run() {
	simulatedHosts, err := j.Database.GetSimulatedHosts()
	if err != nil {
		j.Log.Error(err)
		return
	}

	j.Log.Infof("found %d simulated hosts", len(simulatedHosts))

	for _, simulatedHost := range simulatedHosts {
		if time.Since(simulatedHost.CreatedAt) < 3*time.Minute {
			continue
		}

		if err := j.Database.UpdateHostCores(simulatedHost.Host.Hostname, simulatedHost.Core); err != nil {
			j.Log.Error(err)
			return
		}

		j.Log.Infof("update host %q with original cores", simulatedHost.Host.Hostname)

		if err := j.Database.RemoveSimulatedHost(simulatedHost.ID); err != nil {
			j.Log.Error(err)
			return
		}

		j.Log.Infof("removed simulated host %q", simulatedHost.Host.Hostname)
	}
}
