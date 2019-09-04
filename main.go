package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/amreo/ercole-services/config"
	dataservice_controller "github.com/amreo/ercole-services/data-service/controller"
	dataservice_database "github.com/amreo/ercole-services/data-service/database"
	dataservice_service "github.com/amreo/ercole-services/data-service/service"
	migration "github.com/amreo/ercole-services/database-migration"
	"github.com/rs/cors"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	config := config.ReadConfig()

	enableDataService := flag.Bool("data-service", false, "Enable data service")
	enableAlertService := flag.Bool("alert-service", false, "Enable alert service")
	enableApiService := flag.Bool("api-service", false, "Enable api service")
	enableRepoService := flag.Bool("repo-service", false, "Enable repo service")

	flag.Parse()

	if !*enableDataService && !*enableAlertService && !*enableApiService && !*enableRepoService {
		serve(config, true, true, true, true)
	} else {
		serve(config, *enableDataService, *enableAlertService, *enableApiService, *enableRepoService)
	}

}

func serve(config config.Configuration, enableDataService bool,
	enableAlertService bool, enableApiService bool, enableRepoService bool) {

	//TODO: implement alert-service, api-service, repo-service
	mainRouter := mux.NewRouter()

	if enableDataService {
		db := &dataservice_database.MongoDatabase{
			Config: config,
		}
		service := &dataservice_service.HostDataService{
			Config:   config,
			Version:  "latest",
			Database: db,
		}
		ctrl := &dataservice_controller.HostDataController{
			Config:  config,
			Service: service,
		}
		db.Init()
		migration.Migrate(db.Client.Database(config.Mongodb.DBName))
		service.Init()
		dataservice_controller.SetupRoutesForHostDataController(mainRouter, ctrl)
	}

	log.Println("Start ercole-services: listening at", config.HttpServer.Port)
	var logRouter http.Handler
	if config.HttpServer.LogHttpRequest {
		logRouter = handlers.LoggingHandler(os.Stdout, mainRouter)
	} else {
		logRouter = mainRouter
	}
	http.ListenAndServe(fmt.Sprintf(":%d", config.HttpServer.Port), cors.AllowAll().Handler(logRouter))

}
