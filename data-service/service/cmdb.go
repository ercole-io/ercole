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
	"strings"

	"github.com/ercole-io/ercole/v2/data-service/dto"
	"github.com/ercole-io/ercole/v2/model"
)

func (hds *HostDataService) CompareCmdbInfo(cmdbInfo dto.CmdbInfo) error {
	hostnames, err := hds.Database.GetCurrentHostnames()
	if err != nil {
		return err
	}

	for _, h := range differenceHostnames(cmdbInfo.Hostnames, hostnames) {
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

	for _, h := range differenceHostnames(hostnames, cmdbInfo.Hostnames) {
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

// differenceHostnames returns hostnames in `a` that aren't in `b`
// If a has multiple times on item, which is in b even only once, no occurrences will be returned
func differenceHostnames(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		lcx := strings.ToLower(x)
		mb[lcx] = struct{}{}

		withoutDomain := hostnameWithoutDomain(lcx)
		if lcx != withoutDomain {
			mb[withoutDomain] = struct{}{}
		}
	}

	var diff []string
	for _, x := range a {
		lcx := strings.ToLower(x)
		if _, found := mb[lcx]; found {
			continue
		}

		if _, found := mb[hostnameWithoutDomain(lcx)]; found {
			continue
		}

		diff = append(diff, x)
	}

	return diff
}

func hostnameWithoutDomain(hostname string) string {
	return strings.Split(hostname, ".")[0]
}
