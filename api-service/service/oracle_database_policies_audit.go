// Copyright (c) 2024 Sorint.lab S.p.A.
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
package service

import (
	"github.com/ercole-io/ercole/v2/utils"
)

func (as *APIService) GetOracleDatabasePoliciesAuditFlag(hostname, dbname string) (string, error) {
	policiesAudit, err := as.Database.FindOracleDatabasePoliciesAudit(hostname, dbname)
	if err != nil {
		return "", err
	}

	if len(as.Config.APIService.OracleDatabasePoliciesAudit) == 0 {
		return "N/A", nil
	}

	flag := true

	for _, policyAudit := range as.Config.APIService.OracleDatabasePoliciesAudit {
		if !flag {
			break
		}

		flag = utils.Contains(policiesAudit, policyAudit)
	}

	flagColor := map[bool]string{
		true:  "GREEN",
		false: "RED",
	}

	return flagColor[flag], nil
}
