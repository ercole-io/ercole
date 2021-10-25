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

package client

import (
	"context"

	"github.com/ercole-io/ercole/v2/model"
)

func (c *Client) GetOracleDatabaseLicenseTypes() ([]model.OracleDatabaseLicenseType, error) {
	var response struct {
		LicensesTypes []model.OracleDatabaseLicenseType `json:"license-types"`
	}

	_, err := c.getParsedResponse(context.TODO(), "/settings/oracle/database/license-types", "GET", nil, &response)
	if err != nil {
		return nil, err
	}

	return response.LicensesTypes, nil
}
