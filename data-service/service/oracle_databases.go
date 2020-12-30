package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (hds *HostDataService) oracleDatabasesChecks(hostInfo model.Host, oracleFeature *model.OracleFeature) {
	if oracleFeature.Database == nil || oracleFeature.Database.Databases == nil {
		return
	}

	licenseTypes, err := hds.getOracleDatabaseLicenseTypes()
	if err != nil {
		utils.LogErr(hds.Log, utils.NewAdvancedErrorPtr(err, "INSERT_HOSTDATA_ORACLE_DATABASE"))
		licenseTypes = make([]model.OracleDatabaseLicenseType, 0)
	}

	for i := range oracleFeature.Database.Databases {
		db := &oracleFeature.Database.Databases[i]

		if db.Status == model.OracleDatabaseStatusMounted &&
			db.Role != model.OracleDatabaseRolePrimary {
			hds.addLicensesToSecondaryDbs(hostInfo, db)
		}

		setLicenseTypeIDs(licenseTypes, db)
	}
}

func (hds *HostDataService) addLicensesToSecondaryDbs(hostInfo model.Host, secondaryDb *model.OracleDatabase) {
	dbs, err := hds.getPrimaryOpenOracleDatabases()
	if err != nil {
		utils.LogErr(hds.Log, utils.NewAdvancedErrorPtr(err, "INSERT_HOSTDATA_ORACLE_DATABASE"))
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
		alert := model.Alert{
			AlertCategory:           model.AlertCategoryEngine,
			AlertAffectedTechnology: nil,
			AlertCode:               model.AlertCodeMissingPrimaryDatabase,
			AlertSeverity:           model.AlertSeverityWarning,
			AlertStatus:             model.AlertStatusNew,
			Description:             fmt.Sprintf("Missing primary database on standby database: %s", secondaryDb.Name),
			Date:                    hds.TimeNow(),
			OtherInfo: map[string]interface{}{
				"hostname": hostInfo.Hostname,
				"dbname":   secondaryDb.Name,
			},
		}

		url := utils.NewAPIUrlNoParams(
			hds.Config.AlertService.RemoteEndpoint,
			hds.Config.AlertService.PublisherUsername,
			hds.Config.AlertService.PublisherPassword,
			"/alerts")

		alertBytes, err := json.Marshal(alert)
		if err != nil {
			utils.LogErr(hds.Log, utils.NewAdvancedErrorPtr(err, "Can't marshal alert"))
			return
		}

		resp, err := http.Post(url.String(), "application/json", bytes.NewReader(alertBytes))
		if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
			utils.LogErr(hds.Log, utils.NewAdvancedErrorPtr(err, "Can't throw new alert"))
			return
		}

		return
	}

	coreFactor := secondaryDb.CoreFactor(hostInfo)

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
				})
		}
	}
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
		return nil, utils.NewAdvancedErrorPtr(err, "Can't retrieve databases")
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&dbs); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "Can't decode databases")
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

func (hds *HostDataService) getOracleDatabaseLicenseTypes() ([]model.OracleDatabaseLicenseType, error) {
	url := utils.NewAPIUrlNoParams(
		hds.Config.APIService.RemoteEndpoint,
		hds.Config.APIService.AuthenticationProvider.Username,
		hds.Config.APIService.AuthenticationProvider.Password,
		"settings/oracle/database/license-types").String()

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, utils.NewAdvancedErrorPtr(err, "Can't retrieve licenseTypes")
	}

	licenseTypes := make([]model.OracleDatabaseLicenseType, 0)

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&licenseTypes); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "Can't decode licenseTypes")
	}

	sort.Slice(licenseTypes, licenseTypesSorter(licenseTypes))

	return licenseTypes, nil
}

func licenseTypesSorter(licenseTypes []model.OracleDatabaseLicenseType) func(int, int) bool {
	return func(i, j int) bool {
		x := &licenseTypes[i]
		y := &licenseTypes[j]

		return getPriorityByMetric(x.Metric) >= getPriorityByMetric(y.Metric)
	}
}

func getPriorityByMetric(metric string) int {
	switch metric {
	case model.LicenseTypeMetricProcessorPerpetual:
		return 4
	case model.LicenseTypeMetricNamedUserPlusPerpetual:
		return 3
	case model.LicenseTypePartMetricStreamPerpetual:
		return 2
	case model.LicenseTypeMetricComputerPerpetual:
		return 1
	default:
		return 0
	}
}

func setLicenseTypeIDs(licenseTypes []model.OracleDatabaseLicenseType, database *model.OracleDatabase) {
licenses:
	for i := range database.Licenses {
		license := &database.Licenses[i]

		for _, licenseType := range licenseTypes {
			for _, alias := range licenseType.Aliases {
				if alias == license.Name {
					license.LicenseTypeID = licenseType.ID
					continue licenses
				}
			}
		}
	}
}
