package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/amreo/ercole-services/config"
	dataservice_controller "github.com/amreo/ercole-services/data-service/controller"
	dataservice_database "github.com/amreo/ercole-services/data-service/database"
	dataservice_service "github.com/amreo/ercole-services/data-service/service"

	alertservice_controller "github.com/amreo/ercole-services/alert-service/controller"
	alertservice_database "github.com/amreo/ercole-services/alert-service/database"
	alertservice_service "github.com/amreo/ercole-services/alert-service/service"

	migration "github.com/amreo/ercole-services/database-migration"

	// "github.com/mongodb/amboy"
	"github.com/rs/cors"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// serverVersion holds the version of the server
var serverVersion = "latest"

// main is main
func main() {
	//get the configuration
	config := config.ReadConfig()
	config.Version = serverVersion

	enableDataService := flag.Bool("data-service", false, "Enable data service")
	enableAlertService := flag.Bool("alert-service", false, "Enable alert service")
	enableAPIService := flag.Bool("api-service", false, "Enable api service")
	enableRepoService := flag.Bool("repo-service", false, "Enable repo service")

	flag.Parse()

	if !*enableDataService && !*enableAlertService && !*enableAPIService && !*enableRepoService {
		serve(config, true, true, true, true)
	} else {
		serve(config, *enableDataService, *enableAlertService, *enableAPIService, *enableRepoService)
	}

}

// serve setup and start the services
func serve(config config.Configuration, enableDataService bool,
	enableAlertService bool, enableAPIService bool, enableRepoService bool) {

	var wg sync.WaitGroup

	if config.Mongodb.Migrate {
		log.Print("Migrating...")
		cl := migration.ConnectToMongodb(config.Mongodb)
		migration.Migrate(cl.Database(config.Mongodb.DBName))
		cl.Disconnect(context.TODO())
	}

	if enableDataService {
		serveDataService(config, &wg)
	}

	if enableAlertService {
		serveAlertService(config, &wg)
	}

	//TODO: api-service, repo-service

	wg.Wait()
}

// serveDataService setup and start the data-service
func serveDataService(config config.Configuration, wg *sync.WaitGroup) {
	//Setup the database
	db := &dataservice_database.MongoDatabase{
		Config: config,
	}
	db.Init()

	//Setup the sevice
	service := &dataservice_service.HostDataService{
		Config:   config,
		Version:  "latest",
		Database: db,
	}
	service.Init()

	//Setup the controller
	router := mux.NewRouter()
	ctrl := &dataservice_controller.HostDataController{
		Config:  config,
		Service: service,
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
		Config: config,
	}
	db.Init()

	//Setup the service
	service := &alertservice_service.AlertService{
		Config:   config,
		Database: db,
	}
	service.Init(wg)

	//Setup the controller
	router := mux.NewRouter()
	ctrl := &alertservice_controller.AlertQueueController{
		Config:  config,
		Service: service,
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
