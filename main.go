package main

import (
	"flag"

	"github.com/amreo/ercole-services/config"
)

func init() {
	config.ReadConfig()
}

func main() {
	var enableDataService = flag.Bool("data-service", false, "Enable data service")
	var enableAlertService = flag.Bool("alert-service", false, "Enable alert service")
	var enableApiService = flag.Bool("api-service", false, "Enable api service")
	var enableRepoService = flag.Bool("repo-service", false, "Enable repo service")

}

func serve(config config.Configuration, enableDataService bool,
	enableAlertService bool, enableApiService bool, enableRepoService bool) {

}
