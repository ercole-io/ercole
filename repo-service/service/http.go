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

// Package service is a package that contains varios file serving services
package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/amreo/ercole-services/config"
	"github.com/gorilla/handlers"
	"github.com/rs/cors"
)

// HTTPRepoSubService is a concrete implementation of SubRepoServiceInterface
type HTTPSubRepoService struct {
	// Config contains the reposervice global configuration
	Config config.Configuration
}

// Init start the service
func (hs *HTTPSubRepoService) Init(wg *sync.WaitGroup) {

	//Setup the logger
	router := http.NewServeMux()
	router.Handle("/", http.FileServer(http.Dir(hs.Config.RepoService.DistributedFiles)))
	var logRouter http.Handler

	if hs.Config.RepoService.HTTP.LogHTTPRequest {
		logRouter = handlers.LoggingHandler(os.Stdout, router)
	} else {
		logRouter = router
	}

	wg.Add(1)
	//Start the repo-service
	go func() {
		log.Println("Start repo-service/http: listening at", hs.Config.RepoService.HTTP.Port)
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", hs.Config.RepoService.HTTP.BindIP, hs.Config.RepoService.HTTP.Port), cors.AllowAll().Handler(logRouter))
		log.Println("Stopping repo-service/http", err)
		wg.Done()
	}()
}
