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
package service

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"

	alertservice_client "github.com/ercole-io/ercole/v2/alert-service/client"
	"github.com/ercole-io/ercole/v2/data-service/database"
	"github.com/ercole-io/ercole/v2/data-service/dto"
)

type HostDataServiceInterface interface {
	InsertHostData(hostdata model.HostDataBE) error
	CompareCmdbInfo(cmdbInfo dto.CmdbInfo) error
}

type HostDataService struct {
	Config         config.Configuration
	ServerVersion  string
	Database       database.MongoDatabaseInterface
	AlertSvcClient alertservice_client.AlertSvcClientInterface
	TimeNow        func() time.Time
	Log            *logrus.Logger
}
