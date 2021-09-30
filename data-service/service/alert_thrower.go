package service

import (
	"fmt"
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

// ThrowNewEnterpriseLicenseAlert create and insert in the database a new NEW_DATABASE alert
func (hds *HostDataService) throwNewLicenseAlert(hostname, dbname string, licenseType model.OracleDatabaseLicenseType,
	alreadyEnabledBefore bool) error {

	severity := model.AlertSeverityCritical
	description := fmt.Sprintf("The database %s on %s has enabled new license: %s", dbname, hostname, licenseType.ItemDescription)

	if alreadyEnabledBefore {
		severity = model.AlertSeverityInfo
		description += " (already enabled before in this host)"
	}

	alr := model.Alert{
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

	return hds.AlertSvcClient.ThrowNewAlert(alr)
}

// ThrowActivatedFeaturesAlert create and insert in the database a new NEW_OPTION alert
func (hds *HostDataService) throwNewOptionAlert(hostname, dbname string, licenseType model.OracleDatabaseLicenseType,
	alreadyEnabledBefore bool) error {

	severity := model.AlertSeverityCritical
	description := fmt.Sprintf("The database %s on %s has enabled new option: %s", dbname, hostname, licenseType.ItemDescription)

	if alreadyEnabledBefore {
		severity = model.AlertSeverityInfo
		description += " (already enabled before in this host)"
	}

	alr := model.Alert{
		ID:                      primitive.NewObjectIDFromTimestamp(hds.TimeNow()),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryLicense,
		AlertCode:               model.AlertCodeNewOption,
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

	return hds.AlertSvcClient.ThrowNewAlert(alr)
}

// ThrowUnlistedRunningDatabasesAlert create and insert in the database a new UNLISTED_RUNNING_DATABASE alert
func (hds *HostDataService) throwUnlistedRunningDatabasesAlert(dbname string, hostname string) error {
	alr := model.Alert{
		ID:                      primitive.NewObjectIDFromTimestamp(hds.TimeNow()),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryEngine,
		AlertCode:               model.AlertCodeUnlistedRunningDatabase,
		AlertSeverity:           model.AlertSeverityWarning,
		AlertStatus:             model.AlertStatusNew,
		Date:                    hds.TimeNow(),
		Description:             fmt.Sprintf("The database %s is not listed in the oratab of the host %s", dbname, hostname),
		OtherInfo: map[string]interface{}{
			"hostname": hostname,
			"dbname":   dbname,
		},
	}

	return hds.AlertSvcClient.ThrowNewAlert(alr)
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
