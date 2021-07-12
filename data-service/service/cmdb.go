// Copyright (c) 2021 Sorint.lab S.p.A.
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package service

import (
	"fmt"

	"github.com/ercole-io/ercole/v2/data-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (hds *HostDataService) CompareCmdbInfo(cmdbInfo dto.CmdbInfo) error {
	hostnames, err := hds.Database.GetCurrentHostnames()
	if err != nil {
		return err
	}

	for _, h := range utils.Difference(cmdbInfo.Hostnames, hostnames) {
		alert := model.Alert{
			AlertCategory: model.AlertCategoryEngine,
			AlertCode:     model.AlertCodeMissingHostInErcole,
			AlertSeverity: model.AlertSeverityWarning,
			AlertStatus:   model.AlertStatusNew,
			Description:   fmt.Sprintf("Received unknown hostname %s from CMDB %s", h, cmdbInfo.Name),
			Date:          hds.TimeNow(),
		}

		if err := hds.AlertSvcClient.ThrowNewAlert(alert); err != nil {
			hds.Log.Errorf("Can't create a new alert: %s", err)
		}
	}

	for _, h := range utils.Difference(hostnames, cmdbInfo.Hostnames) {
		alert := model.Alert{
			AlertCategory: model.AlertCategoryEngine,
			AlertCode:     model.AlertCodeMissingHostInCmdb,
			AlertSeverity: model.AlertSeverityWarning,
			AlertStatus:   model.AlertStatusNew,
			Description:   fmt.Sprintf("Missing hostname %s in CMDB %s", h, cmdbInfo.Name),
			Date:          hds.TimeNow(),
		}

		if err := hds.AlertSvcClient.ThrowNewAlert(alert); err != nil {
			hds.Log.Errorf("Can't create a new alert: %s", err)
		}
	}

	return nil
}
