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

package service

import (
	"fmt"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/model"
)

// ThrowNewDatabaseAlert create and insert in the database a new NEW_DATABASE alert
func (hds *HostDataService) throwNewDatabaseAlert(dbname string, hostname string) error {
	alr := model.Alert{
		ID:                      primitive.NewObjectIDFromTimestamp(hds.TimeNow()),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryLicense,
		AlertCode:               model.AlertCodeNewDatabase,
		AlertSeverity:           model.AlertSeverityInfo,
		AlertStatus:             model.AlertStatusNew,
		Date:                    hds.TimeNow(),
		Description:             fmt.Sprintf("The database %s was created on the host %s", dbname, hostname),
		OtherInfo: map[string]interface{}{
			"hostname": hostname,
			"dbname":   dbname,
		},
	}

	return hds.AlertSvcClient.ThrowNewAlert(alr)
}

// ThrowNewServerAlert create and insert in the database a new NEW_SERVER alert
func (hds *HostDataService) throwNewServerAlert(hostname string) error {
	alr := model.Alert{
		ID:                      primitive.NewObjectIDFromTimestamp(hds.TimeNow()),
		AlertAffectedTechnology: nil,
		AlertCategory:           model.AlertCategoryEngine,
		AlertCode:               model.AlertCodeNewServer,
		AlertSeverity:           model.AlertSeverityInfo,
		AlertStatus:             model.AlertStatusNew,
		Date:                    hds.TimeNow(),
		Description:             fmt.Sprintf("The host %s was added to ercole", hostname),
		OtherInfo: map[string]interface{}{
			"hostname": hostname,
		},
	}

	return hds.AlertSvcClient.ThrowNewAlert(alr)
}

func (hds *HostDataService) createNewLicenseAlert(hostname, dbname string, licenseType model.OracleDatabaseLicenseType,
	alreadyEnabledBefore bool) model.Alert {
	severity := model.AlertSeverityCritical
	description := fmt.Sprintf("The database %s on %s has enabled new license: %s", dbname, hostname, licenseType.ItemDescription)

	if alreadyEnabledBefore {
		severity = model.AlertSeverityInfo
		description += " (already enabled before in this host)"
	}

	return model.Alert{
		ID:                      primitive.NewObjectIDFromTimestamp(hds.TimeNow()),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryLicense,
		AlertCode:               model.AlertCodeNewLicense,
		AlertSeverity:           severity,
		AlertStatus:             model.AlertStatusNew,
		Date:                    hds.TimeNow(),
		Description:             description,
		OtherInfo: map[string]interface{}{
			"hostname":      hostname,
			"dbname":        dbname,
			"licenseTypeID": licenseType.ID,
		},
	}
}

// ThrowNewEnterpriseLicenseAlert create and insert in the database a new NEW_LICENSE alert
func (hds *HostDataService) throwNewLicenseAlert(alerts []model.Alert) error {
	if len(alerts) == 0 {
		return nil
	}

	alertOutput := model.Alert{
		ID:                      primitive.NewObjectIDFromTimestamp(hds.TimeNow()),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryLicense,
		AlertCode:               model.AlertCodeNewLicense,
		AlertSeverity:           alerts[0].AlertSeverity,
		AlertStatus:             model.AlertStatusNew,
		Date:                    hds.TimeNow(),
		Description:             alerts[0].Description,
		OtherInfo: map[string]interface{}{},
	}

	if val, ok := alerts[0].OtherInfo["hostname"]; ok {
		alertOutput.OtherInfo["hostname"] = val
	}

	if val, ok := alerts[0].OtherInfo["dbname"]; ok {
		alertOutput.OtherInfo["dbname"] = val
	}

	if val, ok := alerts[0].OtherInfo["licenseTypeID"]; ok {
		alertOutput.OtherInfo["licenseTypeID"] = val
	}

	for _, alert := range alerts[1:] {
		if alert.AlertSeverity == model.AlertSeverityCritical {
			alertOutput.AlertSeverity = model.AlertSeverityCritical
		}

		alertOutput.Description += "\n" + alert.Description

		if val, ok := alert.OtherInfo["dbname"]; ok {
			alertOutput.OtherInfo["dbname"] = alertOutput.OtherInfo["dbname"].(string) + "," + val.(string)
		}

		if val, ok := alert.OtherInfo["licenseTypeID"]; ok {
			alertOutput.OtherInfo["licenseTypeID"] = alertOutput.OtherInfo["licenseTypeID"].(string) + "," + val.(string)
		}
	}

	return hds.AlertSvcClient.ThrowNewAlert(alertOutput)
}

func (hds *HostDataService) throwNewOptionAlerts(alerts []model.Alert) error {
	if len(alerts) == 0 {
		return nil
	}

	alertOutput := model.Alert{
		ID:                      primitive.NewObjectIDFromTimestamp(hds.TimeNow()),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryLicense,
		AlertCode:               model.AlertCodeNewOption,
		AlertSeverity:           alerts[0].AlertSeverity,
		AlertStatus:             model.AlertStatusNew,
		Date:                    hds.TimeNow(),
		Description:             alerts[0].Description,
	}

	if len(alerts[0].OtherInfo) > 0 {
		alertOutput.OtherInfo = map[string]interface{}{
			"hostname":      alerts[0].OtherInfo["hostname"],
			"dbname":        alerts[0].OtherInfo["dbname"],
			"licenseTypeID": alerts[0].OtherInfo["licenseTypeID"],
		}
	}

	for _, alert := range alerts[1:] {
		if alert.AlertSeverity == model.AlertSeverityCritical {
			alertOutput.AlertSeverity = model.AlertSeverityCritical
		}

		alertOutput.Description += "\n" + alert.Description

		if len(alertOutput.OtherInfo) > 0 {
			alertOutput.OtherInfo["dbname"] = alertOutput.OtherInfo["dbname"].(string) + "," + alert.OtherInfo["dbname"].(string)
			alertOutput.OtherInfo["licenseTypeID"] = alertOutput.OtherInfo["licenseTypeID"].(string) + "," + alert.OtherInfo["licenseTypeID"].(string)
		}
	}

	return hds.AlertSvcClient.ThrowNewAlert(alertOutput)
}

// ThrowUnlistedRunningDatabasesAlert create and insert in the database a new UNLISTED_RUNNING_DATABASE alert
func (hds *HostDataService) throwUnlistedRunningDatabasesAlert(alerts []model.Alert) error {
	if len(alerts) == 0 {
		return nil
	}

	a := alerts[0]

	hostname := ""
	dbname := ""

	if val, ok := a.OtherInfo["hostname"]; ok {
		hostname = val.(string)
	}

	if val, ok := a.OtherInfo["dbname"]; ok {
		dbname = val.(string)
	}

	description := fmt.Sprintf("Some databases on the host %s aren't listed in the oratab: %s",
		hostname, dbname)

	alert := model.Alert{
		ID:                      primitive.NewObjectIDFromTimestamp(hds.TimeNow()),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryEngine,
		AlertCode:               model.AlertCodeUnlistedRunningDatabase,
		AlertSeverity:           model.AlertSeverityWarning,
		AlertStatus:             model.AlertStatusNew,
		Date:                    hds.TimeNow(),
		Description:             description,
		OtherInfo: map[string]interface{}{
			"hostname": hostname,
			"dbname":   dbname,
		},
	}

	for _, al := range alerts[1:] {
		if val, ok := al.OtherInfo["dbname"]; ok {
			alert.Description += ", " + val.(string)
			alert.OtherInfo["dbname"] = alert.OtherInfo["dbname"].(string) + "," + val.(string)
		}
	}

	return hds.AlertSvcClient.ThrowNewAlert(alert)
}

func (hds *HostDataService) throwAugmentedCPUCoresAlert(hostname string, previousCPUCores, newCPUCores int) error {
	alr := model.Alert{
		ID:                      primitive.NewObjectIDFromTimestamp(hds.TimeNow()),
		AlertAffectedTechnology: nil,
		AlertCategory:           model.AlertCategoryLicense,
		AlertCode:               model.AlertCodeIncreasedCPUCores,
		AlertSeverity:           model.AlertSeverityCritical,
		AlertStatus:             model.AlertStatusNew,
		Date:                    hds.TimeNow(),
		Description: fmt.Sprintf("The host %s has now more CPU cores: from %d to %d",
			hostname, previousCPUCores, newCPUCores),
		OtherInfo: map[string]interface{}{
			"hostname": hostname,
		},
	}

	return hds.AlertSvcClient.ThrowNewAlert(alr)
}

func (hds *HostDataService) throwMissingPrimaryDatabase(hostname, secondaryDbName string) error {
	alert := model.Alert{
		AlertCategory:           model.AlertCategoryEngine,
		AlertAffectedTechnology: nil,
		AlertCode:               model.AlertCodeMissingPrimaryDatabase,
		AlertSeverity:           model.AlertSeverityWarning,
		AlertStatus:             model.AlertStatusNew,
		Description:             fmt.Sprintf("Missing primary database on standby database: %s", secondaryDbName),
		Date:                    hds.TimeNow(),
		OtherInfo: map[string]interface{}{
			"hostname": hostname,
			"dbname":   secondaryDbName,
		},
	}

	return hds.AlertSvcClient.ThrowNewAlert(alert)
}

func (hds *HostDataService) throwAgentErrorsAlert(hostname string, errs []model.AgentError) error {
	b := strings.Builder{}
	prefix := ""

	if len(errs) > 1 {
		prefix = "- "
	}

	for _, e := range errs {
		b.WriteString(prefix)
		b.WriteString(e.Message)
		b.WriteString("\n")
	}

	alert := model.Alert{
		AlertCategory:           model.AlertCategoryEngine,
		AlertAffectedTechnology: nil,
		AlertCode:               model.AlertCodeAgentError,
		AlertSeverity:           model.AlertSeverityCritical,
		AlertStatus:             model.AlertStatusNew,
		Description:             b.String(),
		Date:                    hds.TimeNow(),
		OtherInfo: map[string]interface{}{
			"hostname": hostname,
			"errors":   errs,
		},
	}

	return hds.AlertSvcClient.ThrowNewAlert(alert)
}

const dbNamesOtherInfo = "dbNames"

func (hds *HostDataService) throwMissingDatabasesAlert(hostname string, dbNames []string, alertSeverity string) error {
	sort.Strings(dbNames)
	description := fmt.Sprintf("The databases %q on %q are missing compared to the previous hostdata",
		strings.Join(dbNames, ", "), hostname)

	alr := model.Alert{
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryLicense,
		AlertCode:               model.AlertCodeMissingDatabase,
		AlertSeverity:           alertSeverity,
		AlertStatus:             model.AlertStatusNew,
		Date:                    hds.TimeNow(),
		Description:             description,
		OtherInfo: map[string]interface{}{
			"hostname":       hostname,
			dbNamesOtherInfo: dbNames,
		},
	}

	return hds.AlertSvcClient.ThrowNewAlert(alr)
}
