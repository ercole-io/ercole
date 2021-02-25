package service

import (
	"fmt"
	"strings"

	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		Description:             fmt.Sprintf("The database '%s' was created on the server %s", dbname, hostname),
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
		Description:             fmt.Sprintf("The server '%s' was added to ercole", hostname),
		OtherInfo: map[string]interface{}{
			"hostname": hostname,
		},
	}

	return hds.AlertSvcClient.ThrowNewAlert(alr)
}

// ThrowNewEnterpriseLicenseAlert create and insert in the database a new NEW_DATABASE alert
func (hds *HostDataService) throwNewEnterpriseLicenseAlert(hostname string) error {
	alr := model.Alert{
		ID:                      primitive.NewObjectIDFromTimestamp(hds.TimeNow()),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryLicense,
		AlertCode:               model.AlertCodeNewLicense,
		AlertSeverity:           model.AlertSeverityCritical,
		AlertStatus:             model.AlertStatusNew,
		Date:                    hds.TimeNow(),
		Description:             fmt.Sprintf("A new Enterprise license has been enabled to %s", hostname),
		OtherInfo: map[string]interface{}{
			"hostname": hostname,
		},
	}

	return hds.AlertSvcClient.ThrowNewAlert(alr)
}

// ThrowActivatedFeaturesAlert create and insert in the database a new NEW_OPTION alert
func (hds *HostDataService) throwActivatedFeaturesAlert(dbname string, hostname string, activatedFeatures []string) error {
	alr := model.Alert{
		ID:                      primitive.NewObjectIDFromTimestamp(hds.TimeNow()),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCategory:           model.AlertCategoryLicense,
		AlertCode:               model.AlertCodeNewOption,
		AlertSeverity:           model.AlertSeverityCritical,
		AlertStatus:             model.AlertStatusNew,
		Date:                    hds.TimeNow(),
		Description:             fmt.Sprintf("The database %s on %s has enabled new features (%s) on server", dbname, hostname, strings.Join(activatedFeatures, ", ")),
		OtherInfo: map[string]interface{}{
			"hostname": hostname,
			"dbname":   dbname,
			"features": activatedFeatures,
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
