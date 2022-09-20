// Copyright (c) 2022 Sorint.lab S.p.A.
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

// Package service is a package that provides methods for querying data
package service

import (
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/exutils"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	alert_filter "github.com/ercole-io/ercole/v2/api-service/dto/filter"
	"github.com/ercole-io/ercole/v2/model"
)

// SearchAlerts search alerts
func (as *APIService) SearchAlerts(alertFilter alert_filter.Alert) (*dto.Pagination, error) {
	return as.Database.SearchAlerts(alertFilter)
}

func (as *APIService) GetAlerts(status string, from, to time.Time, filter dto.GlobalFilter) ([]map[string]interface{}, error) {
	alerts, err := as.Database.GetAlerts(filter.Location, filter.Environment, status, from, to)
	if err != nil {
		return nil, err
	}

	return alerts, nil
}

// SearchAlertsAsXLSX return alerts as xlxs file
func (as *APIService) SearchAlertsAsXLSX(status string, from, to time.Time, filter dto.GlobalFilter) (*excelize.File, error) {
	alerts, err := as.Database.GetAlerts(filter.Location, filter.Environment, status, from, to)
	if err != nil {
		return nil, err
	}

	sheet := "Alerts"
	headers := []string{
		"Type",
		"Date",
		"Severity",
		"Hostname",
		"Code",
		"Description",
	}

	sheets, err := exutils.NewXLSX(as.Config, sheet, headers...)
	if err != nil {
		return nil, err
	}

	axisHelp := exutils.NewAxisHelper(1)

	for _, val := range alerts {
		nextAxis := axisHelp.NewRow()
		sheets.SetCellValue("Alerts", nextAxis(), val["alertCategory"])
		sheets.SetCellValue("Alerts", nextAxis(), val["date"].(primitive.DateTime).Time().UTC().String())
		sheets.SetCellValue("Alerts", nextAxis(), val["alertSeverity"])
		sheets.SetCellValue("Alerts", nextAxis(), val["hostname"])
		sheets.SetCellValue("Alerts", nextAxis(), val["alertCode"])
		sheets.SetCellValue("Alerts", nextAxis(), val["description"])
	}

	return sheets, nil
}

func (as *APIService) AckAlerts(alertsFilter dto.AlertsFilter) error {
	if alertsFilter.AlertCode != nil && *alertsFilter.AlertCode == model.AlertCodeNoData {
		return utils.NewErrorf("%w: you are trying to ack alerts with code: %s",
			utils.ErrInvalidAck,
			model.AlertCodeNoData)
	}

	if alertsFilter.AlertCode == nil {
		count, err := as.Database.CountAlertsNODATA(alertsFilter)
		if err != nil {
			return err
		}

		if count != 0 {
			return utils.NewErrorf("%w: you are trying to ack alerts with code: %s",
				utils.ErrInvalidAck,
				model.AlertCodeNoData)
		}
	}

	return as.Database.UpdateAlertsStatus(alertsFilter, model.AlertStatusAck)
}

func (as *APIService) RemoveAlertsNODATA(alertsFilter dto.AlertsFilter) error {
	return as.Database.RemoveAlertsNODATA(alertsFilter)
}

func (as *APIService) UpdateAlertsStatus(alertsFilter dto.AlertsFilter, newStatus string) error {
	return as.Database.UpdateAlertsStatus(alertsFilter, newStatus)
}
