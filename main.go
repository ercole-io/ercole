package main

import (
	"context"
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
	// "github.com/mongodb/amboy"
	"github.com/rs/cors"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// serverVersion holds the version of the server
var serverVersion = "latest"

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

func serve(config config.Configuration, enableDataService bool,
	enableAlertService bool, enableAPIService bool, enableRepoService bool) {

	if config.Mongodb.Migrate {
		log.Print("Migrating...")
		cl := migration.ConnectToMongodb(config.Mongodb)
		migration.Migrate(cl.Database(config.Mongodb.DBName))
		cl.Disconnect(context.TODO())
	} 

	//TODO: implement alert-service, api-service, repo-service
	mainRouter := mux.NewRouter()

	if enableDataService {
		db := &dataservice_database.MongoDatabase{
			Config: config,
		}
		
		// queue, _ := queue.NewMongoRemoteSingleQueueGroup(
		// 	context.TODO(),
		// 	queue.RemoteQueueGroupOptions{
		// 		Ordered: true,
		// 		TTL:     0,
		// 	},
		// 	db.Client,
		// 	queue.MongoDBOptions{
		// 		URI: config.Mongodb.URI,
		// 		DB:  config.Mongodb.DBName,
		// 	},
		// )
		service := &dataservice_service.HostDataService{
			Config:         config,
			Version:        "latest",
			Database:       db,
			// InsertionQueue: queue,
		}
		ctrl := &dataservice_controller.HostDataController{
			Config:  config,
			Service: service,
		}
		db.Init()
		service.Init()
		dataservice_controller.SetupRoutesForHostDataController(mainRouter, ctrl)
	}

	log.Println("Start ercole-services: listening at", config.HTTPServer.Port)
	var logRouter http.Handler
	if config.HTTPServer.LogHTTPRequest {
		logRouter = handlers.LoggingHandler(os.Stdout, mainRouter)
	} else {
		logRouter = mainRouter
	}
	http.ListenAndServe(fmt.Sprintf(":%d", config.HTTPServer.Port), cors.AllowAll().Handler(logRouter))

}
