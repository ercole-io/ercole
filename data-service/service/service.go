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

// Package service is a package that provides methods for manipulating host informations
package service

import (
	"time"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"

	alertservice_client "github.com/ercole-io/ercole/v2/alert-service/client"
	apiservice_client "github.com/ercole-io/ercole/v2/api-service/client"
	"github.com/ercole-io/ercole/v2/data-service/database"
	"github.com/ercole-io/ercole/v2/data-service/dto"
)

type HostDataServiceInterface interface {
	InsertHostData(hostdata model.HostDataBE) error
	AlertInvalidHostData(validationErr error, hostdata *model.HostDataBE)
	CompareCmdbInfo(cmdbInfo dto.CmdbInfo) error
	InsertOracleLicenseTypes(licenseTypes []model.OracleDatabaseLicenseType) error
	SanitizeLicenseTypes(raw []byte) ([]model.OracleDatabaseLicenseType, error)
	SaveExadata(exadata *model.OracleExadataInstance) error
}

type HostDataService struct {
	Config         config.Configuration
	ServerVersion  string
	Database       database.MongoDatabaseInterface
	AlertSvcClient alertservice_client.AlertSvcClientInterface
	ApiSvcClient   apiservice_client.ApiSvcClientInterface
	TimeNow        func() time.Time
	Log            logger.Logger
}
