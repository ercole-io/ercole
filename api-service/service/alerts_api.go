// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"strings"
	"time"

	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchAlerts search alerts
func (as *APIService) SearchAlerts(search string, sortBy string, sortDesc bool, page int, pageSize int, severity string, status string, from time.Time, to time.Time) ([]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.SearchAlerts(strings.Split(search, " "), sortBy, sortDesc, page, pageSize, severity, status, from, to)
}

// AckAlert ack the specified alert
func (as *APIService) AckAlert(id primitive.ObjectID) utils.AdvancedErrorInterface {
	return as.Database.UpdateAlertStatus(id, model.AlertStatusAck)
}
