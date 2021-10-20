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
	"encoding/json"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (c *Client) GetAlertsByFilter(filter dto.AlertsFilter) ([]model.Alert, error) {
	b := struct {
		Filter dto.AlertsFilter `json:"filter"`
	}{
		Filter: filter,
	}
	body, err := json.Marshal(b)
	if err != nil {
		return nil, utils.NewError(err, "Can't marshal")
	}

	var alerts []model.Alert
	_, err = c.getParsedResponse(context.TODO(), "/alerts", "GET", body, &alerts)
	if err != nil {
		return nil, err
	}

	return alerts, nil
}

func (c *Client) AckAlerts(filter dto.AlertsFilter) error {
	b := struct {
		Filter dto.AlertsFilter `json:"filter"`
	}{
		Filter: filter,
	}
	body, err := json.Marshal(b)
	if err != nil {
		return utils.NewError(err, "Can't marshal")
	}

	_, err = c.getResponse(context.TODO(), "/alerts/ack", "POST", body)
	if err != nil {
		return err
	}

	return nil
}
