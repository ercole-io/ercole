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
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
)

type HistoricizeLicensesComplianceJob struct {
	Database database.MongoDatabaseInterface
	TimeNow  func() time.Time
	Config   config.Configuration
	Log      logger.Logger
}

//TODO Add tests
func (job *HistoricizeLicensesComplianceJob) Run() {
	url := utils.NewAPIUrlNoParams(
		job.Config.APIService.RemoteEndpoint,
		job.Config.APIService.AuthenticationProvider.Username,
		job.Config.APIService.AuthenticationProvider.Password,
		"/hosts/technologies/all/databases/licenses-compliance").String()

	resp, err := http.Get(url)
	if err != nil || resp == nil {
		err = fmt.Errorf("Error while retrieving licenses compliance: [%w], response: [%v]", err, resp)
		job.Log.Error(err)
		return
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = fmt.Errorf("Error while retrieving licenses compliance: response status code: response: [%+v]", resp)
		job.Log.Error(err)
		return
	}

	response := map[string][]dto.LicenseCompliance{}

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&response); err != nil {
		job.Log.Error(err)
		return
	}

	licenses := response["licensesCompliance"]
	err = job.Database.HistoricizeLicensesCompliance(licenses)
	if err != nil {
		job.Log.Error("Can't historicize database licenses")
		return
	}
}
