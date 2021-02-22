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

// Package service is a package that provides methods for manipulating host informations
package job

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/data-service/database"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/sirupsen/logrus"
)

//TODO Throw away this job when agreement history will be implemented

// Save historical data of Oracle Databases licenses
type OracleDbsLicensesHistory struct {
	Database database.MongoDatabaseInterface
	TimeNow  func() time.Time
	Config   config.Configuration
	Log      *logrus.Logger
}

// Run archive every archived hostdata that is older than a amount
func (job *OracleDbsLicensesHistory) Run() {
	licensesCompliance := utils.NewAPIUrlNoParams(
		job.Config.APIService.RemoteEndpoint,
		job.Config.APIService.AuthenticationProvider.Username,
		job.Config.APIService.AuthenticationProvider.Password,
		"/hosts/technologies/oracle/databases/licenses-compliance").String()

	resp, err := http.Get(licensesCompliance)
	if err != nil || resp == nil {
		err = fmt.Errorf("Error while retrieving licenses compliance: [%w], response: [%v]", err, resp)
		utils.LogErr(job.Log, utils.NewAdvancedErrorPtr(err, ""))
		return
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = fmt.Errorf("Error while retrieving licenses compliance: response status code: response: [%+v]", resp)
		utils.LogErr(job.Log, utils.NewAdvancedErrorPtr(err, ""))
		return
	}

	var licenses []dto.OracleDatabaseLicenseUsage
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&licenses); err != nil {
		utils.LogErr(job.Log, utils.NewAdvancedErrorPtr(err, "Decode ERROR"))
		return
	}

	err = job.Database.HistoricizeOracleDbsLicenses(licenses)
	if err != nil {
		utils.LogErr(job.Log, utils.NewAdvancedErrorPtr(err, "Can't historicize Oracle database licenses"))
		return
	}
}
