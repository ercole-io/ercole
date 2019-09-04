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

package data_service

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/amreo/ercole-hostdata-dataservice/config"
	"github.com/amreo/ercole-hostdata-dataservice/controller"
	"github.com/amreo/ercole-hostdata-dataservice/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func init() {
	config.ReadConfig()
}

func main() {
	mainRouter := mux.NewRouter()
	controller.SetupRoutes(mainRouter)
	service.SetupDatabase()
	log.Println("Start ercole-hostdata-dataservice: listening at", config.Config.HttpServer.Listen)
	var logRouter http.Handler
	if config.Config.HttpServer.LogHttpRequest {
		logRouter = handlers.LoggingHandler(os.Stdout, mainRouter)
	} else {
		logRouter = mainRouter
	}
	http.ListenAndServe(fmt.Sprintf(":%d", config.Config.HttpServer.Port), cors.AllowAll().Handler(logRouter))
}
