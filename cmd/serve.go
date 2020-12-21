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

package cmd

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	migration "github.com/ercole-io/ercole/v2/database-migration"

	dataservice_controller "github.com/ercole-io/ercole/v2/data-service/controller"
	dataservice_database "github.com/ercole-io/ercole/v2/data-service/database"
	dataservice_service "github.com/ercole-io/ercole/v2/data-service/service"

	alertservice_controller "github.com/ercole-io/ercole/v2/alert-service/controller"
	alertservice_database "github.com/ercole-io/ercole/v2/alert-service/database"
	alertservice_service "github.com/ercole-io/ercole/v2/alert-service/service"

	apiservice_auth "github.com/ercole-io/ercole/v2/api-service/auth"
	apiservice_controller "github.com/ercole-io/ercole/v2/api-service/controller"
	apiservice_database "github.com/ercole-io/ercole/v2/api-service/database"
	apiservice_service "github.com/ercole-io/ercole/v2/api-service/service"

	chartservice_controller "github.com/ercole-io/ercole/v2/chart-service/controller"
	chartservice_database "github.com/ercole-io/ercole/v2/chart-service/database"
	chartservice_service "github.com/ercole-io/ercole/v2/chart-service/service"

	reposervice_service "github.com/ercole-io/ercole/v2/repo-service/service"
)

var enableDataService bool
var enableAlertService bool
var enableAPIService bool
var enableChartService bool
var enableRepoService bool

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run ercole services",
	Long:  `Run ercole services`,
	Run: func(cmd *cobra.Command, args []string) {
		if !enableDataService && !enableAlertService && !enableAPIService && !enableRepoService && !enableChartService {
			serve(true, true, true, true, true)
		} else {
			serve(enableDataService, enableAlertService, enableAPIService, enableChartService, enableRepoService)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().BoolVarP(&enableDataService, "enable-data-service", "d", false, "Enable the data service")
	serveCmd.Flags().BoolVarP(&enableAlertService, "enable-alert-service", "a", false, "Enable the alert service")
	serveCmd.Flags().BoolVarP(&enableAPIService, "enable-api-service", "u", false, "Enable the api service")
	serveCmd.Flags().BoolVarP(&enableChartService, "enable-chart-service", "t", false, "Enable the chart service")
	serveCmd.Flags().BoolVarP(&enableRepoService, "enable-repo-service", "r", false, "Enable the repo service")
}

// serve setup and start the services
func serve(enableDataService bool,
	enableAlertService bool, enableAPIService bool, enableChartService bool, enableRepoService bool) {
	log := utils.NewLogger("SERV")

	if !utils.FileExists(ercoleConfig.RepoService.DistributedFiles) {
		log.Warnf("The directory %s for RepoService doesn't exist so the RepoService will be disabled\n", ercoleConfig.RepoService.DistributedFiles)
		enableRepoService = false
	}
	if ercoleConfig.ResourceFilePath == "RESOURCES_NOT_FOUND" {
		log.Warn("A directory for resources wasn't found so some services may not work as expected")
	}

	var wg sync.WaitGroup

	if ercoleConfig.Mongodb.Migrate {
		log.Info("Migrating...")
		cl := migration.ConnectToMongodb(log, ercoleConfig.Mongodb)
		migration.Migrate(log, cl.Database(ercoleConfig.Mongodb.DBName))
		cl.Disconnect(context.TODO())
	}

	if enableDataService {
		serveDataService(ercoleConfig, &wg)
	}

	if enableAlertService {
		serveAlertService(ercoleConfig, &wg)
	}

	if enableAPIService {
		serveAPIService(ercoleConfig, &wg)
	}

	if enableChartService {
		serveChartService(ercoleConfig, &wg)
	}

	if enableRepoService {
		serveRepoService(ercoleConfig, &wg)
	}

	wg.Wait()
}

func serveDataService(config config.Configuration, wg *sync.WaitGroup) {
	log := utils.NewLogger("DATA")

	db := &dataservice_database.MongoDatabase{
		Config:  config,
		TimeNow: time.Now,
		Log:     log,
	}
	db.Init()

	service := &dataservice_service.HostDataService{
		Config:        config,
		ServerVersion: config.Version,
		Database:      db,
		TimeNow:       time.Now,
		Log:           log,
	}
	service.Init()

	router := mux.NewRouter()
	ctrl := &dataservice_controller.DataController{
		Config:  config,
		Service: service,
		TimeNow: time.Now,
		Log:     log,
	}
	dataservice_controller.SetupRoutesForHostDataController(router, ctrl)

	var logRouter http.Handler
	if config.DataService.LogHTTPRequest {
		logRouter = utils.CustomLoggingHandler(router, log)
	} else {
		logRouter = router
	}

	wg.Add(1)
	go func() {
		log.Println("Start data-service: listening at", config.DataService.Port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.DataService.BindIP, config.DataService.Port), cors.AllowAll().Handler(logRouter))
		log.Println("Stopping data-service", err)
		wg.Done()
	}()
}

func serveAlertService(config config.Configuration, wg *sync.WaitGroup) {
	log := utils.NewLogger("ALRT")

	db := &alertservice_database.MongoDatabase{
		Config:  config,
		TimeNow: time.Now,
		Log:     log,
	}
	db.Init()

	emailer := &alertservice_service.SMTPEmailer{
		Config: config,
	}

	service := &alertservice_service.AlertService{
		Config:   config,
		Database: db,
		TimeNow:  time.Now,
		Log:      log,
		Emailer:  emailer,
	}
	service.Init(wg)

	router := mux.NewRouter()
	ctrl := &alertservice_controller.AlertQueueController{
		Config:  config,
		Service: service,
		TimeNow: time.Now,
		Log:     log,
	}
	alertservice_controller.SetupRoutesForAlertQueueController(router, ctrl)

	var logRouter http.Handler
	if config.DataService.LogHTTPRequest {
		logRouter = utils.CustomLoggingHandler(router, log)
	} else {
		logRouter = router
	}

	wg.Add(1)
	go func() {
		log.Println("Start alert-service: listening at", config.AlertService.Port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.AlertService.BindIP, config.AlertService.Port), cors.AllowAll().Handler(logRouter))
		log.Println("Stopping alert-service", err)
		wg.Done()
	}()
}

func serveAPIService(config config.Configuration, wg *sync.WaitGroup) {
	log := utils.NewLogger("APIS")

	if config.APIService.DebugOracleDatabaseAgreementsAssignmentAlgorithm {
		log.Level = logrus.DebugLevel
	}

	db := &apiservice_database.MongoDatabase{
		Config:                          config,
		TimeNow:                         time.Now,
		OperatingSystemAggregationRules: config.APIService.OperatingSystemAggregationRules,
		Log:                             log,
	}
	db.Init()

	service := &apiservice_service.APIService{
		Config:   config,
		Version:  serverVersion,
		Database: db,
		TimeNow:  time.Now,
		Log:      log,
	}
	service.Init()

	auth := apiservice_auth.BuildAuthenticationProvider(config.APIService.AuthenticationProvider, time.Now, log)
	auth.Init()

	router := mux.NewRouter()
	ctrl := &apiservice_controller.APIController{
		Config:        config,
		Service:       service,
		TimeNow:       time.Now,
		Log:           log,
		Authenticator: auth,
	}
	apiservice_controller.SetupRoutesForAPIController(router, ctrl, auth)

	var logRouter http.Handler
	if config.DataService.LogHTTPRequest {
		logRouter = utils.CustomLoggingHandler(router, log)
	} else {
		logRouter = router
	}

	wg.Add(1)
	go func() {
		log.Info("Start api-service: listening at ", config.APIService.Port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.APIService.BindIP, config.APIService.Port), cors.AllowAll().Handler(logRouter))
		log.Warn("Stopping api-service", err)
		wg.Done()
	}()
}

func serveChartService(config config.Configuration, wg *sync.WaitGroup) {
	log := utils.NewLogger("CHRT")

	db := &chartservice_database.MongoDatabase{
		Config:                          config,
		TimeNow:                         time.Now,
		Log:                             log,
		OperatingSystemAggregationRules: config.APIService.OperatingSystemAggregationRules,
	}
	db.Init()

	service := &chartservice_service.ChartService{
		Config:   config,
		Database: db,
		TimeNow:  time.Now,
		Log:      log,
	}
	service.Init()

	auth := apiservice_auth.BuildAuthenticationProvider(config.APIService.AuthenticationProvider, time.Now, log)
	auth.Init()

	router := mux.NewRouter()
	ctrl := &chartservice_controller.ChartController{
		Config:        config,
		Service:       service,
		TimeNow:       time.Now,
		Log:           log,
		Authenticator: auth,
	}
	chartservice_controller.SetupRoutesForChartController(router, ctrl, auth)

	var logRouter http.Handler
	if config.DataService.LogHTTPRequest {
		logRouter = utils.CustomLoggingHandler(router, log)
	} else {
		logRouter = router
	}

	wg.Add(1)
	go func() {
		log.Info("Start chart-service: listening at ", config.ChartService.Port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.ChartService.BindIP, config.ChartService.Port), cors.AllowAll().Handler(logRouter))
		log.Warn("Stopping chart-service", err)
		wg.Done()
	}()
}

func serveRepoService(config config.Configuration, wg *sync.WaitGroup) {
	service := &reposervice_service.RepoService{
		Config:      config,
		SubServices: []reposervice_service.SubRepoServiceInterface{},
	}

	if config.RepoService.HTTP.Enable {
		service.SubServices = append(service.SubServices,
			&reposervice_service.HTTPSubRepoService{
				Config: config,
				Log:    utils.NewLogger("REPO"),
			})
	}

	if config.RepoService.SFTP.Enable {
		service.SubServices = append(service.SubServices,
			&reposervice_service.SFTPRepoSubService{
				Config: config,
				Log:    utils.NewLogger("REPO"),
			})
	}

	service.Init(wg)
}
