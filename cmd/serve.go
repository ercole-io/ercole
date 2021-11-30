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

	"github.com/gorilla/handlers"
	"github.com/rs/cors"
	"github.com/spf13/cobra"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"

	migration "github.com/ercole-io/ercole/v2/database-migration"

	dataservice_controller "github.com/ercole-io/ercole/v2/data-service/controller"
	dataservice_database "github.com/ercole-io/ercole/v2/data-service/database"
	dataservice_job "github.com/ercole-io/ercole/v2/data-service/job"
	dataservice_service "github.com/ercole-io/ercole/v2/data-service/service"

	alertservice_client "github.com/ercole-io/ercole/v2/alert-service/client"
	alertservice_controller "github.com/ercole-io/ercole/v2/alert-service/controller"
	alertservice_database "github.com/ercole-io/ercole/v2/alert-service/database"
	alertservice_emailer "github.com/ercole-io/ercole/v2/alert-service/emailer"
	alertservice_service "github.com/ercole-io/ercole/v2/alert-service/service"

	apiservice_auth "github.com/ercole-io/ercole/v2/api-service/auth"
	apiservice_client "github.com/ercole-io/ercole/v2/api-service/client"
	apiservice_controller "github.com/ercole-io/ercole/v2/api-service/controller"
	apiservice_database "github.com/ercole-io/ercole/v2/api-service/database"
	apiservice_service "github.com/ercole-io/ercole/v2/api-service/service"

	chartservice_controller "github.com/ercole-io/ercole/v2/chart-service/controller"
	chartservice_database "github.com/ercole-io/ercole/v2/chart-service/database"
	chartservice_service "github.com/ercole-io/ercole/v2/chart-service/service"

	thunderservice_controller "github.com/ercole-io/ercole/v2/thunder-service/controller"
	thunderservice_database "github.com/ercole-io/ercole/v2/thunder-service/database"
	thunderservice_service "github.com/ercole-io/ercole/v2/thunder-service/service"

	reposervice_service "github.com/ercole-io/ercole/v2/repo-service/service"
)

var enableDataService bool
var enableAlertService bool
var enableAPIService bool
var enableChartService bool
var enableRepoService bool
var enableThunderService bool

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run ercole services",
	Long:  `Run ercole services`,
	Run: func(cmd *cobra.Command, args []string) {
		if !enableDataService && !enableAlertService && !enableAPIService && !enableRepoService && !enableChartService && !enableThunderService {
			serve(true, true, true, true, true, true)
		} else {
			serve(enableDataService, enableAlertService, enableAPIService, enableChartService, enableRepoService, enableThunderService)
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
	serveCmd.Flags().BoolVarP(&enableThunderService, "enable-thunder-service", "s", false, "Enable the thunder service")
}

// serve setup and start the services
func serve(enableDataService bool,
	enableAlertService bool, enableAPIService bool, enableChartService bool, enableRepoService bool, enableThunderService bool) {
	log := logger.NewLogger("SERV", logger.LogVerbosely(verbose))

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
		err := migration.Migrate(ercoleConfig.Mongodb)
		if err != nil {
			log.Fatalf("Failed migrating database: %s", err)
		}
	}

	if enableDataService || enableAlertService || enableAPIService || enableChartService || enableThunderService {
		check, err := migration.IsAtTheLatestVersion(ercoleConfig.Mongodb)
		if err != nil {
			log.Fatalf("Failed checking database version: %s", err)
		}

		if !check {
			log.Fatal("Database is not at the latest version\nYou can migrate to the latest version by running `ercole migrate`")
		}
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

	if enableThunderService {
		serveThunderService(ercoleConfig, &wg)
	}

	wg.Wait()
}

func serveDataService(config config.Configuration, wg *sync.WaitGroup) {
	log := logger.NewLogger("DATA", logger.LogVerbosely(verbose))

	db := &dataservice_database.MongoDatabase{
		Config:  config,
		TimeNow: time.Now,
		Log:     log,
	}
	db.Init()

	service := &dataservice_service.HostDataService{
		Config:         config,
		ServerVersion:  config.Version,
		Database:       db,
		AlertSvcClient: alertservice_client.NewClient(config.AlertService),
		ApiSvcClient:   apiservice_client.NewClient(config.APIService),
		TimeNow:        time.Now,
		Log:            log,
	}

	job := &dataservice_job.Job{
		Config:        config,
		ServerVersion: config.Version,
		Database:      db,
		TimeNow:       time.Now,
		Log:           log,
	}
	job.Init()

	ctrl := &dataservice_controller.DataController{
		Config:  config,
		Service: service,
		TimeNow: time.Now,
		Log:     log,
	}
	h := ctrl.GetDataControllerHandler()
	h = useCommonHandlers(h, config.DataService.LogHTTPRequest, log)

	wg.Add(1)
	go func() {
		log.Info("Start data-service: listening at ", config.DataService.Port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.DataService.BindIP, config.DataService.Port), h)
		if err != nil {
			log.Error("Stopped data-service: ", err)
		}

		wg.Done()
	}()
}

func serveAlertService(config config.Configuration, wg *sync.WaitGroup) {
	log := logger.NewLogger("ALRT", logger.LogVerbosely(verbose))

	db := &alertservice_database.MongoDatabase{
		Config:  config,
		TimeNow: time.Now,
		Log:     log,
	}
	db.Init()

	emailer := &alertservice_emailer.SMTPEmailer{
		Config: config,
	}

	service := &alertservice_service.AlertService{
		Config:   config,
		Database: db,
		TimeNow:  time.Now,
		Log:      log,
		Emailer:  emailer,
	}
	ctx, cancel := context.WithCancel(context.Background())
	service.Init(ctx, wg)

	ctrl := &alertservice_controller.AlertQueueController{
		Config:  config,
		Service: service,
		TimeNow: time.Now,
		Log:     log,
	}
	h := ctrl.GetAlertControllerHandler()
	h = useCommonHandlers(h, config.AlertService.LogHTTPRequest, log)

	wg.Add(1)
	go func() {
		log.Info("Start alert-service: listening at ", config.AlertService.Port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.AlertService.BindIP, config.AlertService.Port), h)
		if err != nil {
			log.Error("Stopping alert-service: ", err)
		}

		cancel()
		wg.Done()
	}()
}

func serveAPIService(config config.Configuration, wg *sync.WaitGroup) {
	log := logger.NewLogger("APIS", logger.LogVerbosely(verbose))

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

	ctrl := &apiservice_controller.APIController{
		Config:        config,
		Service:       service,
		TimeNow:       time.Now,
		Log:           log,
		Authenticator: auth,
	}
	h := ctrl.GetApiControllerHandler(auth)
	h = useCommonHandlers(h, config.APIService.LogHTTPRequest, log)

	wg.Add(1)
	go func() {
		log.Info("Start api-service: listening at ", config.APIService.Port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.APIService.BindIP, config.APIService.Port), h)
		if err != nil {
			log.Error("Stopping api-service: ", err)
		}
		wg.Done()
	}()
}

func serveChartService(config config.Configuration, wg *sync.WaitGroup) {
	log := logger.NewLogger("CHRT", logger.LogVerbosely(verbose))

	db := &chartservice_database.MongoDatabase{
		Config:                          config,
		TimeNow:                         time.Now,
		Log:                             log,
		OperatingSystemAggregationRules: config.APIService.OperatingSystemAggregationRules,
	}
	db.Init()

	service := &chartservice_service.ChartService{
		Config:       config,
		Database:     db,
		ApiSvcClient: apiservice_client.NewClient(config.APIService),
		TimeNow:      time.Now,
		Log:          log,
	}
	service.Init()

	auth := apiservice_auth.BuildAuthenticationProvider(config.APIService.AuthenticationProvider, time.Now, log)
	auth.Init()

	ctrl := &chartservice_controller.ChartController{
		Config:        config,
		Service:       service,
		TimeNow:       time.Now,
		Log:           log,
		Authenticator: auth,
	}

	h := ctrl.GetChartControllerHandler(auth)
	h = useCommonHandlers(h, config.ChartService.LogHTTPRequest, log)

	wg.Add(1)
	go func() {
		log.Info("Start chart-service: listening at ", config.ChartService.Port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.ChartService.BindIP, config.ChartService.Port), h)
		if err != nil {
			log.Error("Stopping chart-service: ", err)
		}

		wg.Done()
	}()
}

func serveRepoService(config config.Configuration, wg *sync.WaitGroup) {
	service := &reposervice_service.RepoService{
		Config:      config,
		SubServices: []reposervice_service.SubRepoServiceInterface{},
	}

	log := logger.NewLogger("REPO", logger.LogVerbosely(verbose))

	if config.RepoService.HTTP.Enable {
		service.SubServices = append(service.SubServices,
			&reposervice_service.HTTPSubRepoService{
				Config: config,
				Log:    log,
			})
	}

	service.Init(wg)
}

func serveThunderService(config config.Configuration, wg *sync.WaitGroup) {
	log := logger.NewLogger("THUN", logger.LogVerbosely(verbose))

	db := &thunderservice_database.MongoDatabase{
		Config:  config,
		TimeNow: time.Now,
		Log:     log,
	}
	db.Init()

	service := &thunderservice_service.ThunderService{
		Config:   config,
		Database: db,
		TimeNow:  time.Now,
		Log:      log,
	}
	service.Init()

	ctrl := &thunderservice_controller.ThunderController{
		Config:  config,
		Service: service,
		TimeNow: time.Now,
		Log:     log,
	}
	h := ctrl.GetThunderControllerHandler()
	h = useCommonHandlers(h, config.ThunderService.LogHTTPRequest, log)

	wg.Add(1)
	go func() {
		log.Info("Start thunder-service: listening at ", config.ThunderService.Port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.ThunderService.BindIP, config.ThunderService.Port), h)
		if err != nil {
			log.Error("Stopping thunder-service: ", err)
		}

		wg.Done()
	}()
}

func useCommonHandlers(h http.Handler, logHTTPRequest bool, log logger.Logger) http.Handler {
	if logHTTPRequest {
		h = utils.CustomLoggingHandler(h, log)
	}

	rl := recoveryLogger{
		l: log,
	}
	h = handlers.RecoveryHandler(
		handlers.PrintRecoveryStack(true),
		handlers.RecoveryLogger(rl),
	)(h)

	return cors.AllowAll().Handler(h)
}

type recoveryLogger struct {
	l logger.Logger
}

func (rl recoveryLogger) Println(args ...interface{}) {
	args = append([]interface{}{"Panic:\n"}, args...)

	rl.l.Error(args...)
}
