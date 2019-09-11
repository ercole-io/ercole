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

// Package service is a package that provides methods for manipulating host informations
package service

import (
	"log"
	"sync"

	"github.com/amreo/ercole-services/alert-service/database"
	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/amreo/ercole-services/config"

	"github.com/leandro-lugaresi/hub"
)

// AlertServiceInterface is a interface that wrap methods used to insert and process alert messages
type AlertServiceInterface interface {
	// Init initializes the service
	Init(wg *sync.WaitGroup)
	// HostDataInsertion inserts the host data insertion in the queue
	HostDataInsertion(id primitive.ObjectID) utils.AdvancedErrorInterface
	// ProcessMsg processes the message msg
	ProcessMsg(msg hub.Message)
}

// AlertService is the concrete implementation of HostDataServiceInterface. It saves data to a MongoDB database
type AlertService struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Queue that contains all messages to be processed
	Queue *hub.Hub
	// Database contains the database layer
	Database database.MongoDatabaseInterface
}

// Init initializes the service and database
func (as *AlertService) Init(wg *sync.WaitGroup) {
	//Create a new queue
	as.Queue = hub.New()

	//Subscribe the alert-service
	sub := as.Queue.Subscribe(0, "hostdata.insertion")
	wg.Add(1)
	go func(s hub.Subscription) {
		log.Println("Start alert-service/queue")
		for msg := range s.Receiver {
			as.ProcessMsg(msg)
		}
		log.Println("Stop alert-service/queue")
		wg.Done()
	}(sub)
}

// HostDataInsertion inserts the host data insertion in the queue
func (as *AlertService) HostDataInsertion(id primitive.ObjectID) utils.AdvancedErrorInterface {
	as.Queue.Publish(hub.Message{
		Name: "hostdata.insertion",
		Fields: hub.Fields{
			"id": id,
		},
	})
	return nil
}

// ProcessMsg processes the message msg
func (as *AlertService) ProcessMsg(msg hub.Message) {
	if as.Config.AlertService.LogMessages {
		log.Printf("RECEIVED EVENT %s: %s", msg.Topic(), utils.ToJson(msg.Fields))
		switch msg.Topic() {
		case model.TopicHostDataInsertion:
			as.ProcessHostDataInsertion(msg.Fields)
		default:
		}
	}
}

// ProcessHostDataInsertion processes the host data insertion event
func (as *AlertService) ProcessHostDataInsertion(params hub.Fields) {
	id := params["id"].(primitive.ObjectID)

	data, err := as.Database.FindHostData(id)
	log.Println(data, err)
}
