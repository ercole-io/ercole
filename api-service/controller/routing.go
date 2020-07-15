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

package controller

import (
	"net/http"

	"github.com/ercole-io/ercole/api-service/auth"
	"github.com/gorilla/mux"
)

// SetupRoutesForAPIController setup the routes of the router using the handler in the controller as http handler
func SetupRoutesForAPIController(router *mux.Router, ctrl APIControllerInterface, auth auth.AuthenticationProvider) {

	//Add the routes
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})

	router.HandleFunc("/user/login", auth.GetToken).Methods("POST")
	//Enable authentication using the ctrl
	router = router.NewRoute().Subrouter()
	router.Use(auth.AuthenticateMiddleware)
	setupProtectedRoutes(router, ctrl)
}

func setupProtectedRoutes(router *mux.Router, ctrl APIControllerInterface) {
	router.HandleFunc("/hosts", ctrl.SearchHosts).Methods("GET")
	router.HandleFunc("/hosts/count", ctrl.GetHostsCountStats).Methods("GET")
	router.HandleFunc("/hosts/environments/frequency", ctrl.GetEnvironmentStats).Methods("GET")
	router.HandleFunc("/hosts/types", ctrl.GetTypeStats).Methods("GET")
	router.HandleFunc("/hosts/operating-systems", ctrl.GetOperatingSystemStats).Methods("GET")
	router.HandleFunc("/hosts/technologies", ctrl.ListTechnologies).Methods("GET")
	router.HandleFunc("/hosts/technologies/compliance", ctrl.GetTotalTechnologiesComplianceStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/license-modifiers", ctrl.SearchOracleDatabaseLicenseModifiers).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/top-unused-instance-resource", ctrl.GetTopUnusedOracleDatabaseInstanceResourceStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/environments", ctrl.GetOracleDatabaseEnvironmentStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/high-reliability", ctrl.GetOracleDatabaseHighReliabilityStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/versions", ctrl.GetOracleDatabaseVersionStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/top-reclaimable", ctrl.GetTopReclaimableOracleDatabaseStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/patch-status", ctrl.GetOracleDatabasePatchStatusStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/top-workload", ctrl.GetTopWorkloadOracleDatabaseStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/dataguard-status", ctrl.GetOracleDatabaseDataguardStatusStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/archivelog-status", ctrl.GetOracleDatabaseArchivelogStatusStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/rac-status", ctrl.GetOracleDatabaseRACStatusStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/total-work", ctrl.GetTotalOracleDatabaseWorkStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/total-memory-size", ctrl.GetTotalOracleDatabaseMemorySizeStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/total-datafile-size", ctrl.GetTotalOracleDatabaseDatafileSizeStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/total-segment-size", ctrl.GetTotalOracleDatabaseSegmentSizeStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/license-compliance", ctrl.GetOracleDatabaseLicenseComplianceStatusStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases", ctrl.SearchOracleDatabases).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/addms", ctrl.SearchOracleDatabaseAddms).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/segment-advisors", ctrl.SearchOracleDatabaseSegmentAdvisors).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/databases/patch-advisors", ctrl.SearchOracleDatabasePatchAdvisors).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/exadata", ctrl.SearchOracleExadata).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/exadata/total-memory-size", ctrl.GetTotalOracleExadataMemorySizeStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/exadata/total-cpu", ctrl.GetTotalOracleExadataCPUStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/exadata/average-storage-usage", ctrl.GetAverageOracleExadataStorageUsageStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/exadata/storage-error-count-status", ctrl.GetOracleExadataStorageErrorCountStatusStats).Methods("GET")
	router.HandleFunc("/hosts/technologies/oracle/exadata/patch-status", ctrl.GetOracleExadataPatchStatusStats).Methods("GET")
	router.HandleFunc("/hosts/locations", ctrl.ListLocations).Methods("GET")
	router.HandleFunc("/hosts/environments", ctrl.ListEnvironments).Methods("GET")
	router.HandleFunc("/hosts/clusters", ctrl.SearchClusters).Methods("GET")
	router.HandleFunc("/hosts/clusters/{name}", ctrl.GetCluster).Methods("GET")
	router.HandleFunc("/hosts/{hostname}", ctrl.GetHost).Methods("GET")
	router.HandleFunc("/hosts/{hostname}", ctrl.ArchiveHost).Methods("DELETE")
	router.HandleFunc("/hosts/{hostname}/patching-function", ctrl.GetPatchingFunction).Methods("GET")
	router.HandleFunc("/hosts/{hostname}/patching-function", ctrl.SetPatchingFunction).Methods("PUT")
	router.HandleFunc("/hosts/{hostname}/patching-function", ctrl.DeletePatchingFunction).Methods("DELETE")
	router.HandleFunc("/hosts/{hostname}/technologies/oracle/databases/{dbname}/tags", ctrl.AddTagToOracleDatabase).Methods("POST")
	router.HandleFunc("/hosts/{hostname}/technologies/oracle/databases/{dbname}/tags/{tagname}", ctrl.DeleteTagOfOracleDatabase).Methods("DELETE")
	router.HandleFunc("/hosts/{hostname}/technologies/oracle/databases/{dbname}/licenses/{licenseName}", ctrl.SetOracleDatabaseLicenseModifier).Methods("PUT")
	router.HandleFunc("/hosts/{hostname}/technologies/oracle/databases/{dbname}/license-modifiers/{licenseName}", ctrl.DeleteOracleDatabaseLicenseModifier).Methods("DELETE")
	router.HandleFunc("/licenses", ctrl.ListLicenses).Methods("GET")
	router.HandleFunc("/licenses", ctrl.SetLicensesCount).Methods("PUT")
	router.HandleFunc("/licenses/{name}", ctrl.GetLicense).Methods("GET")
	router.HandleFunc("/licenses/{name}/count", ctrl.SetLicenseCount).Methods("PUT")
	router.HandleFunc("/licenses/{name}/cost-per-processor", ctrl.SetLicenseCostPerProcessor).Methods("PUT")
	router.HandleFunc("/licenses/{name}/unlimited-status", ctrl.SetLicenseUnlimitedStatus).Methods("PUT")
	router.HandleFunc("/alerts", ctrl.SearchAlerts).Methods("GET")
	router.HandleFunc("/alerts/{id}", ctrl.AckAlert).Methods("POST")

	setupSettingsRoutes(router.PathPrefix("/settings").Subrouter(), ctrl)
	setupFrontendAPIRoutes(router.PathPrefix("/frontend").Subrouter(), ctrl)
}

func setupSettingsRoutes(router *mux.Router, ctrl APIControllerInterface) {
	router.HandleFunc("/default-database-tag-choiches", ctrl.GetDefaultDatabaseTags).Methods("GET")
	router.HandleFunc("/features", ctrl.GetErcoleFeatures).Methods("GET")
	router.HandleFunc("/technologies", ctrl.GetTechnologyList).Methods("GET")
}

func setupFrontendAPIRoutes(router *mux.Router, ctrl APIControllerInterface) {
	router.HandleFunc("/dashboard", ctrl.GetInfoForFrontendDashboard).Methods("GET")
}
