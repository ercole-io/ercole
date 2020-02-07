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

package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ThrowNewDatabaseAlert create and insert in the database a new NEW_DATABASE alert
func (as *AlertService) ThrowNewDatabaseAlert(dbname string, hostname string) utils.AdvancedErrorInterface {
	_, err := as.Database.InsertAlert(model.Alert{
		ID:            primitive.NewObjectIDFromTimestamp(as.TimeNow()),
		AlertCode:     model.AlertCodeNewDatabase,
		AlertSeverity: model.AlertSeverityNotice,
		AlertStatus:   model.AlertStatusNew,
		Date:          as.TimeNow(),
		Description:   fmt.Sprintf("The database '%s' was created on the server %s", dbname, hostname),
		OtherInfo: map[string]interface{}{
			"Hostname": hostname,
			"Dbname":   dbname,
		},
	})
	if as.Config.AlertService.LogAlertThrows {
		log.Printf("Alert NEW_DATABASE of %s/%s was thrown\n", hostname, dbname)
	}
	return err
}

// ThrowNewServerAlert create and insert in the database a new NEW_SERVER alert
func (as *AlertService) ThrowNewServerAlert(hostname string) utils.AdvancedErrorInterface {
	_, err := as.Database.InsertAlert(model.Alert{
		ID:            primitive.NewObjectIDFromTimestamp(as.TimeNow()),
		AlertCode:     model.AlertCodeNewServer,
		AlertSeverity: model.AlertSeverityNotice,
		AlertStatus:   model.AlertStatusNew,
		Date:          as.TimeNow(),
		Description:   fmt.Sprintf("The server '%s' was added to ercole", hostname),
		OtherInfo: map[string]interface{}{
			"Hostname": hostname,
		},
	})
	if as.Config.AlertService.LogAlertThrows {
		log.Printf("Alert NEW_SERVER of %s was thrown\n", hostname)
	}
	return err
}

// ThrowNewEnterpriseLicenseAlert create and insert in the database a new NEW_DATABASE alert
func (as *AlertService) ThrowNewEnterpriseLicenseAlert(hostname string) utils.AdvancedErrorInterface {
	_, err := as.Database.InsertAlert(model.Alert{
		ID:            primitive.NewObjectIDFromTimestamp(as.TimeNow()),
		AlertCode:     model.AlertCodeNewLicense,
		AlertSeverity: model.AlertSeverityCritical,
		AlertStatus:   model.AlertStatusNew,
		Date:          as.TimeNow(),
		Description:   fmt.Sprintf("A new Enterprise license has been enabled to %s", hostname),
		OtherInfo: map[string]interface{}{
			"Hostname": hostname,
		},
	})
	if as.Config.AlertService.LogAlertThrows {
		log.Printf("Alert NEW_LICENSE of %s was thrown\n", hostname)
	}
	return err
}

// ThrowActivatedFeaturesAlert create and insert in the database a new NEW_OPTION alert
func (as *AlertService) ThrowActivatedFeaturesAlert(dbname string, hostname string, activatedFeatures []string) utils.AdvancedErrorInterface {
	_, err := as.Database.InsertAlert(model.Alert{
		ID:            primitive.NewObjectIDFromTimestamp(as.TimeNow()),
		AlertCode:     model.AlertCodeNewOption,
		AlertSeverity: model.AlertSeverityCritical,
		AlertStatus:   model.AlertStatusNew,
		Date:          as.TimeNow(),
		Description:   fmt.Sprintf("The database %s on %s has enabled new features (%s) on server", dbname, hostname, strings.Join(activatedFeatures, ", ")),
		OtherInfo: map[string]interface{}{
			"Hostname": hostname,
			"Dbname":   dbname,
			"Features": activatedFeatures,
		},
	})
	if as.Config.AlertService.LogAlertThrows {
		log.Printf("Alert NEW_OPTIONS of %s/%s was thrown\n", hostname, dbname)
	}
	return err
}

// ThrowNoDataAlert create and insert in the database a new NO_DATA alert
func (as *AlertService) ThrowNoDataAlert(hostname string, freshnessThreshold int) utils.AdvancedErrorInterface {
	_, err := as.Database.InsertAlert(model.Alert{
		ID:            primitive.NewObjectIDFromTimestamp(as.TimeNow()),
		AlertCode:     model.AlertCodeNoData,
		AlertSeverity: model.AlertSeverityMajor,
		AlertStatus:   model.AlertStatusNew,
		Date:          as.TimeNow(),
		Description:   fmt.Sprintf("No data received from the host %s in the last %d days", hostname, freshnessThreshold),
		OtherInfo: map[string]interface{}{
			"Hostname": hostname,
		},
	})
	if as.Config.AlertService.LogAlertThrows {
		log.Printf("Alert NO_DATA of %s was thrown\n", hostname)
	}
	return err
}
