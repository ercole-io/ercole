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

// Package service is a package that provides methods for querying data
package service

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/utils/exutils"
	"strings"
	"time"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchAlerts search alerts
func (as *APIService) SearchAlerts(mode string, search string, sortBy string, sortDesc bool,
	page, pageSize int, location, environment, severity, status string, from, to time.Time,
) ([]map[string]interface{}, error) {
	return as.Database.SearchAlerts(mode, strings.Split(search, " "), sortBy, sortDesc, page, pageSize,
		location, environment, severity, status, from, to)
}
// SearchAlertsAsXLSX return alerts as xlxs file
func (as *APIService) SearchAlertsAsXLSX(from, to time.Time, filter dto.GlobalFilter) (*excelize.File, error) {
	alerts, err := as.Database.SearchAlerts("all", []string{}, "", false, -1, -1, filter.Location, filter.Environment, "", "", from, to)
	if err != nil {
		return nil, err
	}

	sheet := "Alerts"
	headers := []string {
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

// AckAlerts ack the specified alerts
func (as *APIService) AckAlerts(ids []primitive.ObjectID) error {
	return as.Database.UpdateAlertsStatus(ids, model.AlertStatusAck)
}

func (as *APIService) AckAlertsByFilter(alertsFilter dto.AlertsFilter) error {
	return as.Database.UpdateAlertsStatusByFilter(alertsFilter, model.AlertStatusAck)
}
