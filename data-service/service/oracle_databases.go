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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (hds *HostDataService) oracleDatabasesChecks(previousHostdata, hostdata *model.HostDataBE) {
	if hostdata.Features.Oracle.Database == nil || hostdata.Features.Oracle.Database.Databases == nil {
		return
	}

	hds.checkSecondaryDbs(hostdata)

	licenseTypes, err := hds.getOracleDatabaseLicenseTypes(hostdata.Environment)
	if err != nil {
		hds.Log.Error(err)
		licenseTypes = make([]model.OracleDatabaseLicenseType, 0)
	}

	hds.setLicenseTypes(hostdata, licenseTypes)

	hds.checkNewLicenses(previousHostdata, hostdata, licenseTypes)

	for _, dbname := range hostdata.Features.Oracle.Database.UnlistedRunningDatabases {
		if err := hds.ackOldUnlistedRunningDatabasesAlerts(hostdata.Hostname, dbname); err != nil {
			hds.Log.Errorf("Can't ack UnlistedRunningDatabases alerts by filter")
		}
		if err := hds.throwUnlistedRunningDatabasesAlert(dbname, hostdata.Hostname); err != nil {
			hds.Log.Error(err)
		}
	}

	if previousHostdata != nil && previousHostdata.Info.CPUCores < hostdata.Info.CPUCores {
		if err := hds.throwAugmentedCPUCoresAlert(hostdata.Hostname,
			previousHostdata.Info.CPUCores,
			hostdata.Info.CPUCores); err != nil {
			hds.Log.Error(err)
		}
	}

	hds.checkMissingDatabases(previousHostdata, hostdata)
}

func (hds *HostDataService) ackOldUnlistedRunningDatabasesAlerts(hostname, dbname string) error {
	f := dto.AlertsFilter{
		AlertCategory:           utils.Str2ptr(model.AlertCategoryEngine),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCode:               utils.Str2ptr(model.AlertCodeUnlistedRunningDatabase),
		AlertSeverity:           utils.Str2ptr(model.AlertSeverityWarning),
		OtherInfo: map[string]interface{}{
			"hostname": hostname,
			"dbname":   dbname,
		},
	}
	return hds.ApiSvcClient.AckAlerts(f)
}

func (hds *HostDataService) checkSecondaryDbs(hostdata *model.HostDataBE) {
	for i := range hostdata.Features.Oracle.Database.Databases {
		db := &hostdata.Features.Oracle.Database.Databases[i]

		if db.Status == model.OracleDatabaseStatusMounted &&
			db.Role != model.OracleDatabaseRolePrimary {
			hds.addLicensesToSecondaryDb(hostdata.Info, hostdata.CoreFactor(), db)
		}
	}
}

func (hds *HostDataService) addLicensesToSecondaryDb(hostInfo model.Host, hostCoreFactor float64, secondaryDb *model.OracleDatabase) {
	dbs, err := hds.getPrimaryOpenOracleDatabases()
	if err != nil {
		hds.Log.Error("INSERT_HOSTDATA_ORACLE_DATABASE")
		return
	}

	var primaryDb *model.OracleDatabase
	for i, db := range dbs {
		if db.DbID == secondaryDb.DbID && db.Name == secondaryDb.Name {
			primaryDb = &dbs[i]
			break
		}
	}

	if primaryDb == nil {
		if err := hds.ackOldMissingPrimaryDbAlerts(hostInfo.Hostname, secondaryDb.Name); err != nil {
			hds.Log.Errorf("Can't ack MissingPrimaryDatabase alerts by filter")
		}

		if err := hds.throwMissingPrimaryDatabase(hostInfo.Hostname, secondaryDb.Name); err != nil {
			hds.Log.Errorf("Can't throw missing primary database alert, hostname %s, secondaryDbName %s",
				hostInfo.Hostname, secondaryDb.Name)
		}

		return
	}

	coreFactor, err := secondaryDb.CoreFactor(hostInfo, hostCoreFactor)
	if err != nil {
		hds.Log.Error(err.Error())
		return
	}

primaryDbLicensesCycle:
	for _, primaryDbLicense := range primaryDb.Licenses {

		if primaryDbLicense.Count > 0 {
			for i := range secondaryDb.Licenses {
				secondaryDbLicense := &secondaryDb.Licenses[i]

				if secondaryDbLicense.Name == primaryDbLicense.Name {
					secondaryDbLicense.Count = float64(hostInfo.CPUCores) * coreFactor
					continue primaryDbLicensesCycle
				}
			}

			secondaryDb.Licenses = append(secondaryDb.Licenses,
				model.OracleDatabaseLicense{
					LicenseTypeID: primaryDbLicense.LicenseTypeID,
					Name:          primaryDbLicense.Name,
					Count:         float64(hostInfo.CPUCores) * coreFactor,
					Ignored:       primaryDbLicense.Ignored,
				})
		}
	}
}

func (hds *HostDataService) ackOldMissingPrimaryDbAlerts(hostname, dbname string) error {
	f := dto.AlertsFilter{
		AlertCategory: utils.Str2ptr(model.AlertCategoryEngine),
		AlertCode:     utils.Str2ptr(model.AlertCodeMissingPrimaryDatabase),
		AlertSeverity: utils.Str2ptr(model.AlertSeverityWarning),
		OtherInfo: map[string]interface{}{
			"hostname": hostname,
			"dbname":   dbname,
		},
	}
	return hds.ApiSvcClient.AckAlerts(f)
}

func (hds *HostDataService) getPrimaryOpenOracleDatabases() (dbs []model.OracleDatabase, err error) {
	values := url.Values{}
	values.Set("full", "true")
	url := utils.NewAPIUrl(
		hds.Config.APIService.RemoteEndpoint,
		hds.Config.APIService.AuthenticationProvider.Username,
		hds.Config.APIService.AuthenticationProvider.Password,
		"/hosts/technologies/oracle/databases", values).String()

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, utils.NewError(err, "Can't retrieve databases")
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&dbs); err != nil {
		return nil, utils.NewError(err, "Can't decode databases")
	}

	for i := 0; i < len(dbs); {
		db := &dbs[i]

		if db.Role == model.OracleDatabaseRolePrimary && db.Status == model.OracleDatabaseStatusOpen {
			i += 1
			continue
		}

		dbs = removeFromDBs(dbs, i)
	}

	return dbs, nil
}

// Do not mantain order
func removeFromDBs(s []model.OracleDatabase, i int) []model.OracleDatabase {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (hds *HostDataService) setLicenseTypes(hostdata *model.HostDataBE, licenseTypes []model.OracleDatabaseLicenseType) {
	for i := range hostdata.Features.Oracle.Database.Databases {
		db := &hostdata.Features.Oracle.Database.Databases[i]
		setLicenseTypeIDs(licenseTypes, db)
	}
}

func setLicenseTypeIDs(licenseTypes []model.OracleDatabaseLicenseType, database *model.OracleDatabase) {
	lics := database.Licenses

	// remove empty licenses
	for i := 0; i < len(lics); {
		if lics[i].Count > 0 {
			i++
			continue
		}

		lics[i] = lics[len(lics)-1]
		lics = lics[:len(lics)-1]
	}

	sort.Slice(lics, func(i, j int) bool {
		return lics[i].Count > lics[j].Count ||
			lics[i].Name < lics[j].Name
	})

	licenseTypesAlreadyUsed := make(map[string]bool) // use each LicenseType once per database

licenses:
	for i := range lics {
		license := &lics[i]

		for _, licenseType := range licenseTypes {
			if licenseTypesAlreadyUsed[licenseType.ID] {
				continue
			}

			for _, alias := range licenseType.Aliases {
				if alias == license.Name {
					license.LicenseTypeID = licenseType.ID
					licenseTypesAlreadyUsed[licenseType.ID] = true

					continue licenses
				}
			}
		}
	}

	database.Licenses = lics
}

func (hds *HostDataService) getOracleDatabaseLicenseTypes(environment string) ([]model.OracleDatabaseLicenseType, error) {
	licenseTypes, err := hds.ApiSvcClient.GetOracleDatabaseLicenseTypes()
	if err != nil {
		return nil, utils.NewError(err, "Can't retrieve licenseTypes")
	}

	sort.Slice(licenseTypes, licenseTypesSorter(hds.Config.DataService, environment, licenseTypes))

	return licenseTypes, nil
}

func licenseTypesSorter(config config.DataService, environment string, licenseTypes []model.OracleDatabaseLicenseType,
) func(int, int) bool {
	orderOfPriority, ok := config.LicenseTypeMetricsByEnvironment[environment]
	if !ok {
		orderOfPriority = config.LicenseTypeMetricsDefault
	}

	priorityMap := make(map[string]int, len(orderOfPriority))
	for i, p := range orderOfPriority {
		priorityMap[p] = len(orderOfPriority) - i
	}

	return func(i, j int) bool {
		x := &licenseTypes[i]
		y := &licenseTypes[j]

		if priorityMap[x.Metric] == priorityMap[y.Metric] {
			return x.ItemDescription < y.ItemDescription
		}

		return priorityMap[x.Metric] > priorityMap[y.Metric]
	}
}

func (hds *HostDataService) checkNewLicenses(previous, new *model.HostDataBE, licenseTypes []model.OracleDatabaseLicenseType) {
	previousDbs := make(map[string]model.OracleDatabase)
	if previous != nil &&
		previous.Features.Oracle != nil &&
		previous.Features.Oracle.Database != nil &&
		previous.Features.Oracle.Database.Databases != nil {
		previousDbs = model.DatabaseSliceAsMap(previous.Features.Oracle.Database.Databases)
	}

	newDbs := make(map[string]model.OracleDatabase)
	if new.Features.Oracle.Database != nil && new.Features.Oracle.Database.Databases != nil {
		newDbs = model.DatabaseSliceAsMap(new.Features.Oracle.Database.Databases)
	}

	previousLicenseTypesEnabled := make(map[string]bool)
	for _, db := range previousDbs {
		for _, license := range db.Licenses {
			if license.Count > 0 {
				previousLicenseTypesEnabled[license.LicenseTypeID] = true
			}
		}
	}

	licenseTypesMap := make(map[string]*model.OracleDatabaseLicenseType)
	for i := range licenseTypes {
		licenseTypesMap[licenseTypes[i].ID] = &licenseTypes[i]
	}

	var newOptionAlerts []model.Alert
	for _, newDb := range newDbs {

		oldDb, ok := previousDbs[newDb.Name]

		if !ok {
			oldDb = model.OracleDatabase{
				Licenses: []model.OracleDatabaseLicense{},
			}

			if err := hds.throwNewDatabaseAlert(newDb.Name, new.Hostname); err != nil {
				hds.Log.Error(err)
			}
		}

		diffs := model.DiffLicenses(oldDb.Licenses, newDb.Licenses)

		for licenseTypeID, diffFeature := range diffs {

			if diffFeature == model.DiffFeatureActivated {
				licenseType, ok := licenseTypesMap[licenseTypeID]
				if !ok {
					hds.Log.Warnf("%v: %s", utils.ErrOracleDatabaseLicenseTypeIDNotFound, licenseTypeID)
					continue
				}

				alreadyEnabledBefore := previousLicenseTypesEnabled[licenseTypeID]
				if licenseType.Option {

					description := fmt.Sprintf("Database %s has enabled new option: %s", newDb.Name, licenseType.ItemDescription)
					severity := model.AlertSeverityCritical

					if alreadyEnabledBefore {
						severity = model.AlertSeverityInfo
						description += " (already enabled before in this host)"
					}

					newOptionAlerts = append(newOptionAlerts, model.Alert{
						AlertSeverity: severity,
						Description:   description,
						OtherInfo: map[string]interface{}{
							"hostname":      new.Hostname,
							"dbname":        newDb.Name,
							"licenseTypeID": licenseType.ID,
						},
					})

				} else {
					if err := hds.throwNewLicenseAlert(new.Hostname, newDb.Name, *licenseType, alreadyEnabledBefore); err != nil {
						hds.Log.Error(err)
					}
				}
			}
		}
	}

	if err := hds.throwNewOptionAlerts(newOptionAlerts); err != nil {
		hds.Log.Error(err)
	}
}

func (hds *HostDataService) checkMissingDatabases(previous, new *model.HostDataBE) {
	if previous == nil ||
		previous.Features.Oracle == nil ||
		previous.Features.Oracle.Database == nil ||
		previous.Features.Oracle.Database.Databases == nil {
		return
	}

	newDbs := getDbNames(new.Features.Oracle.Database.Databases)

	if err := hds.searchAndAckOldMissingDatabasesAlerts(new.Hostname, newDbs); err != nil {
		hds.Log.Error(err)
	}

	previousDbs := getDbNames(previous.Features.Oracle.Database.Databases)

	if err := hds.findMissingDatabasesAndThrowAlerts(new.Hostname, newDbs, previousDbs); err != nil {
		hds.Log.Error(err)
	}
}

func getDbNames(dbs []model.OracleDatabase) map[string]bool {
	m := make(map[string]bool, len(dbs))
	for i := range dbs {
		m[dbs[i].Name] = true
	}

	return m
}

func (hds *HostDataService) searchAndAckOldMissingDatabasesAlerts(hostname string, newDbs map[string]bool) error {
	f := dto.AlertsFilter{
		AlertCategory:           utils.Str2ptr(model.AlertCategoryLicense),
		AlertAffectedTechnology: model.TechnologyOracleDatabasePtr,
		AlertCode:               utils.Str2ptr(model.AlertCodeMissingDatabase),
		AlertStatus:             utils.Str2ptr(model.AlertStatusNew),
		OtherInfo: map[string]interface{}{
			"hostname": hostname,
		},
	}
	alerts, err := hds.ApiSvcClient.GetAlertsByFilter(f)
	if err != nil {
		return err
	}

alerts:
	for i := range alerts {
		dbNamesInterf, ok := alerts[i].OtherInfo[dbNamesOtherInfo]
		if !ok {
			continue
		}

		dbNames, ok := dbNamesInterf.([]string)
		if !ok {
			return utils.NewErrorf("Can't convert dbNames in []string: %s", dbNamesInterf)
		}

		for _, dbName := range dbNames {
			if !newDbs[dbName] {
				continue alerts
			}
		}

		// all dbNames are in newDbs
		f = dto.AlertsFilter{IDs: []primitive.ObjectID{alerts[i].ID}}
		err = hds.ApiSvcClient.AckAlerts(f)
		if err != nil {
			hds.Log.Error(err)
		}
	}

	return nil
}

func (hds *HostDataService) findMissingDatabasesAndThrowAlerts(hostname string, newDbs, previousDbs map[string]bool) error {
	severity := model.AlertSeverityCritical
	missingDbs := make([]string, 0)

	for x := range previousDbs {
		if newDbs[x] {
			severity = model.AlertSeverityWarning // at least one database is still present
		} else {
			missingDbs = append(missingDbs, x)
		}
	}

	if len(missingDbs) > 0 {
		err := hds.throwMissingDatabasesAlert(hostname, missingDbs, severity)
		if err != nil {
			return err
		}
	}

	return nil
}
