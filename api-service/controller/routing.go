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
	router.HandleFunc("/hosts/{hostname}", ctrl.GetHost).Methods("GET")
	router.HandleFunc("/hosts/{hostname}", ctrl.ArchiveHost).Methods("DELETE")
	router.HandleFunc("/hosts/{hostname}/patching-function", ctrl.GetPatchingFunction).Methods("GET")
	router.HandleFunc("/hosts/{hostname}/patching-function", ctrl.SetPatchingFunction).Methods("PUT")
	router.HandleFunc("/hosts/{hostname}/patching-function", ctrl.DeletePatchingFunction).Methods("DELETE")
	router.HandleFunc("/hosts/{hostname}/databases/{dbname}/tags", ctrl.AddTagToDatabase).Methods("POST")
	router.HandleFunc("/hosts/{hostname}/databases/{dbname}/tags/{tagname}", ctrl.DeleteTagOfDatabase).Methods("DELETE")
	router.HandleFunc("/hosts/{hostname}/databases/{dbname}/licenses/{licenseName}", ctrl.SetLicenseModifier).Methods("PUT")
	router.HandleFunc("/hosts/{hostname}/databases/{dbname}/license-modifiers/{licenseName}", ctrl.DeleteLicenseModifier).Methods("DELETE")

	router.HandleFunc("/locations", ctrl.ListLocations).Methods("GET")
	router.HandleFunc("/environments", ctrl.ListEnvironments).Methods("GET")
	router.HandleFunc("/assets", ctrl.ListAssets).Methods("GET")
	router.HandleFunc("/clusters", ctrl.SearchClusters).Methods("GET")
	router.HandleFunc("/databases", ctrl.SearchDatabases).Methods("GET")
	router.HandleFunc("/addms", ctrl.SearchAddms).Methods("GET")
	router.HandleFunc("/segment-advisors", ctrl.SearchSegmentAdvisors).Methods("GET")
	router.HandleFunc("/patch-advisors", ctrl.SearchPatchAdvisors).Methods("GET")
	router.HandleFunc("/licenses", ctrl.ListLicenses).Methods("GET")
	router.HandleFunc("/licenses", ctrl.SetLicensesCount).Methods("PUT")
	router.HandleFunc("/licenses/{name}", ctrl.GetLicense).Methods("GET")
	router.HandleFunc("/licenses/{name}/count", ctrl.SetLicenseCount).Methods("PUT")
	router.HandleFunc("/licenses/{name}/cost-per-processor", ctrl.SetLicenseCostPerProcessor).Methods("PUT")
	router.HandleFunc("/licenses/{name}/unlimited-status", ctrl.SetLicenseUnlimitedStatus).Methods("PUT")
	router.HandleFunc("/exadata", ctrl.SearchExadata).Methods("GET")
	router.HandleFunc("/alerts", ctrl.SearchAlerts).Methods("GET")
	router.HandleFunc("/alerts/{id}", ctrl.AckAlert).Methods("DELETE")
	router.HandleFunc("/license-modifiers", ctrl.SearchLicenseModifiers).Methods("GET")
	router.HandleFunc("/stats/count", ctrl.GetHostsCountStats).Methods("GET")
	router.HandleFunc("/stats/environments", ctrl.GetEnvironmentStats).Methods("GET")
	router.HandleFunc("/stats/types", ctrl.GetTypeStats).Methods("GET")
	router.HandleFunc("/stats/operating-systems", ctrl.GetOperatingSystemStats).Methods("GET")
	router.HandleFunc("/stats/top-unused-instance-resource", ctrl.GetTopUnusedInstanceResourceStats).Methods("GET")
	router.HandleFunc("/stats/assets/compliance", ctrl.GetTotalAssetsComplianceStats).Methods("GET")
	router.HandleFunc("/stats/databases/environments", ctrl.GetDatabaseEnvironmentStats).Methods("GET")
	router.HandleFunc("/stats/databases/high-reliability", ctrl.GetDatabaseHighReliabilityStats).Methods("GET")
	router.HandleFunc("/stats/databases/versions", ctrl.GetDatabaseVersionStats).Methods("GET")
	router.HandleFunc("/stats/databases/top-reclaimable", ctrl.GetTopReclaimableDatabaseStats).Methods("GET")
	router.HandleFunc("/stats/databases/patch-status", ctrl.GetDatabasePatchStatusStats).Methods("GET")
	router.HandleFunc("/stats/databases/top-workload", ctrl.GetTopWorkloadDatabaseStats).Methods("GET")
	router.HandleFunc("/stats/databases/dataguard-status", ctrl.GetDatabaseDataguardStatusStats).Methods("GET")
	router.HandleFunc("/stats/databases/archivelog-status", ctrl.GetDatabaseArchivelogStatusStats).Methods("GET")
	router.HandleFunc("/stats/databases/rac-status", ctrl.GetDatabaseRACStatusStats).Methods("GET")
	router.HandleFunc("/stats/databases/total-work", ctrl.GetTotalDatabaseWorkStats).Methods("GET")
	router.HandleFunc("/stats/databases/total-memory-size", ctrl.GetTotalDatabaseMemorySizeStats).Methods("GET")
	router.HandleFunc("/stats/databases/total-datafile-size", ctrl.GetTotalDatabaseDatafileSizeStats).Methods("GET")
	router.HandleFunc("/stats/databases/total-segment-size", ctrl.GetTotalDatabaseSegmentSizeStats).Methods("GET")
	router.HandleFunc("/stats/databases/license-compliance", ctrl.GetDatabaseLicenseComplianceStatusStats).Methods("GET")
	router.HandleFunc("/stats/exadata/total-memory-size", ctrl.GetTotalExadataMemorySizeStats).Methods("GET")
	router.HandleFunc("/stats/exadata/total-cpu", ctrl.GetTotalExadataCPUStats).Methods("GET")
	router.HandleFunc("/stats/exadata/average-storage-usage", ctrl.GetAverageExadataStorageUsageStats).Methods("GET")
	router.HandleFunc("/stats/exadata/storage-error-count-status", ctrl.GetExadataStorageErrorCountStatusStats).Methods("GET")
	router.HandleFunc("/stats/exadata/patch-status", ctrl.GetExadataPatchStatusStats).Methods("GET")

	setupSettingsRoutes(router.PathPrefix("/settings").Subrouter(), ctrl)
	setupFrontendAPIRoutes(router.PathPrefix("/frontend").Subrouter(), ctrl)
}

func setupSettingsRoutes(router *mux.Router, ctrl APIControllerInterface) {
	router.HandleFunc("/default-database-tag-choiches", ctrl.GetDefaultDatabaseTags).Methods("GET")
	router.HandleFunc("/features", ctrl.GetErcoleFeatures).Methods("GET")

}

func setupFrontendAPIRoutes(router *mux.Router, ctrl APIControllerInterface) {
	router.HandleFunc("/dashboard", ctrl.GetInfoForFrontendDashboard).Methods("GET")
}
