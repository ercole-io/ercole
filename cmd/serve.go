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

	//TODO: api-service, repo-service

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
		http.ListenAndServe(fmt.Sprintf("%s:%d", config.DataService.BindIP, config.DataService.Port), cors.AllowAll().Handler(logRouter))
		log.Println("Stopping data-service")
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
	//Start the data-service
	go func() {
		log.Println("Start alert-service: listening at", config.AlertService.Port)
		http.ListenAndServe(fmt.Sprintf("%s:%d", config.AlertService.BindIP, config.AlertService.Port), cors.AllowAll().Handler(logRouter))
		log.Println("Stopping alert-service")
		wg.Done()
	}()
}
