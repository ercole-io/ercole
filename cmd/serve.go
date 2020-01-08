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

package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/amreo/ercole-services/config"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"

	migration "github.com/amreo/ercole-services/database-migration"

	dataservice_controller "github.com/amreo/ercole-services/data-service/controller"
	dataservice_database "github.com/amreo/ercole-services/data-service/database"
	dataservice_service "github.com/amreo/ercole-services/data-service/service"

	alertservice_controller "github.com/amreo/ercole-services/alert-service/controller"
	alertservice_database "github.com/amreo/ercole-services/alert-service/database"
	alertservice_service "github.com/amreo/ercole-services/alert-service/service"

	apiservice_controller "github.com/amreo/ercole-services/api-service/controller"
	apiservice_database "github.com/amreo/ercole-services/api-service/database"
	apiservice_service "github.com/amreo/ercole-services/api-service/service"

	reposervice_service "github.com/amreo/ercole-services/repo-service/service"
)

var enableDataService bool
var enableAlertService bool
var enableAPIService bool
var enableRepoService bool

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run ercole services",
	Long:  `Run ercole services`,
	Run: func(cmd *cobra.Command, args []string) {
		if !enableDataService && !enableAlertService && !enableAPIService && !enableRepoService {
			serve(true, true, true, true)
		} else {
			serve(enableDataService, enableAlertService, enableAPIService, enableRepoService)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().BoolVarP(&enableDataService, "enable-data-service", "d", false, "Enable the data service")
	serveCmd.Flags().BoolVarP(&enableAlertService, "enable-alert-service", "a", false, "Enable the alert service")
	serveCmd.Flags().BoolVarP(&enableAPIService, "enable-api-service", "u", false, "Enable the api service")
	serveCmd.Flags().BoolVarP(&enableRepoService, "enable-repo-service", "r", false, "Enable the repo service")
}

// serve setup and start the services
func serve(enableDataService bool,
	enableAlertService bool, enableAPIService bool, enableRepoService bool) {

	s, _ := os.Readlink("/proc/self/exe")
	s = filepath.Dir(s)
	ercoleConfig.RepoService.DistributedFiles = s + filepath.Join("/", ercoleConfig.RepoService.DistributedFiles) + "/"

	if _, err := os.Stat(ercoleConfig.RepoService.DistributedFiles); os.IsNotExist(err) {
		log.Printf("WARNING: the directory %s for RepoService doesn't exist so the RepoService will be disabled\n", ercoleConfig.RepoService.DistributedFiles)
		enableRepoService = false
	}

	var wg sync.WaitGroup

	if ercoleConfig.Mongodb.Migrate {
		log.Print("Migrating...")
		cl := migration.ConnectToMongodb(ercoleConfig.Mongodb)
		migration.Migrate(cl.Database(ercoleConfig.Mongodb.DBName))
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

	if enableRepoService {
		serveRepoService(ercoleConfig, &wg)
	}

	wg.Wait()
}

// serveDataService setup and start the data-service
func serveDataService(config config.Configuration, wg *sync.WaitGroup) {
	//Setup the database
	db := &dataservice_database.MongoDatabase{
		Config:  config,
		TimeNow: time.Now,
	}
	db.Init()

	//Setup the sevice
	service := &dataservice_service.HostDataService{
		Config:   config,
		Version:  "latest",
		Database: db,
		TimeNow:  time.Now,
	}
	service.Init()

	//Setup the controller
	router := mux.NewRouter()
	ctrl := &dataservice_controller.HostDataController{
		Config:  config,
		Service: service,
		TimeNow: time.Now,
	}
	dataservice_controller.SetupRoutesForHostDataController(router, ctrl)

	//Setup the logger
	var logRouter http.Handler
	if config.DataService.LogHTTPRequest {
		logRouter = handlers.LoggingHandler(os.Stdout, router)
	} else {
		logRouter = router
	}

	wg.Add(1)
	//Start the data-service
	go func() {
		log.Println("Start data-service: listening at", config.DataService.Port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.DataService.BindIP, config.DataService.Port), cors.AllowAll().Handler(logRouter))
		log.Println("Stopping data-service", err)
		wg.Done()
	}()
}

// serveAlertService setup and start the alert-service
func serveAlertService(config config.Configuration, wg *sync.WaitGroup) {
	//Setup the database
	db := &alertservice_database.MongoDatabase{
		Config:  config,
		TimeNow: time.Now,
	}
	db.Init()

	//Setup the service
	service := &alertservice_service.AlertService{
		Config:   config,
		Database: db,
		TimeNow:  time.Now,
	}
	service.Init(wg)

	//Setup the controller
	router := mux.NewRouter()
	ctrl := &alertservice_controller.AlertQueueController{
		Config:  config,
		Service: service,
		TimeNow: time.Now,
	}
	alertservice_controller.SetupRoutesForAlertQueueController(router, ctrl)

	//Setup the logger
	var logRouter http.Handler
	if config.DataService.LogHTTPRequest {
		logRouter = handlers.LoggingHandler(os.Stdout, router)
	} else {
		logRouter = router
	}

	wg.Add(1)
	//Start the alert-service
	go func() {
		log.Println("Start alert-service: listening at", config.AlertService.Port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.AlertService.BindIP, config.AlertService.Port), cors.AllowAll().Handler(logRouter))
		log.Println("Stopping alert-service", err)
		wg.Done()
	}()
}

// serveAPIService setup and start the api-service
func serveAPIService(config config.Configuration, wg *sync.WaitGroup) {
	//Setup the database
	db := &apiservice_database.MongoDatabase{
		Config:                          config,
		TimeNow:                         time.Now,
		OperatingSystemAggregationRules: config.APIService.OperatingSystemAggregationRules,
	}
	db.Init()

	//Setup the service
	service := &apiservice_service.APIService{
		Config:   config,
		Database: db,
		TimeNow:  time.Now,
	}
	service.Init()

	//Setup the controller
	router := mux.NewRouter()
	ctrl := &apiservice_controller.APIController{
		Config:  config,
		Service: service,
		TimeNow: time.Now,
	}
	apiservice_controller.SetupRoutesForAPIController(router, ctrl)

	//Setup the logger
	var logRouter http.Handler
	if config.DataService.LogHTTPRequest {
		logRouter = handlers.LoggingHandler(os.Stdout, router)
	} else {
		logRouter = router
	}

	wg.Add(1)
	//Start the api-service
	go func() {
		log.Println("Start api-service: listening at", config.APIService.Port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.APIService.BindIP, config.APIService.Port), cors.AllowAll().Handler(logRouter))
		log.Println("Stopping api-service", err)
		wg.Done()
	}()
}

// serveRepoService setup and start the repo-service
func serveRepoService(config config.Configuration, wg *sync.WaitGroup) {
	//Setup the service
	service := &reposervice_service.RepoService{
		Config:      config,
		SubServices: []reposervice_service.SubRepoServiceInterface{},
	}
	if config.RepoService.HTTP.Enable {
		service.SubServices = append(service.SubServices, &reposervice_service.HTTPSubRepoService{
			Config: config,
		})
	}
	if config.RepoService.SFTP.Enable {
		service.SubServices = append(service.SubServices, &reposervice_service.SFTPRepoSubService{
			Config: config,
		})
	}

	//Init and serve
	service.Init(wg)
}
