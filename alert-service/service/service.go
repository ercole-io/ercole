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

// Package service is a package that provides methods for manipulating host informations
package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ercole-io/ercole/v2/alert-service/database"
	"github.com/ercole-io/ercole/v2/alert-service/emailer"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/sirupsen/logrus"

	"github.com/ercole-io/ercole/v2/config"

	"github.com/leandro-lugaresi/hub"
)

// AlertServiceInterface is a interface that wrap methods used to insert and process alert messages
type AlertServiceInterface interface {
	// Init initializes the service
	Init(ctx context.Context, wg *sync.WaitGroup)
	// ProcessMsg processes the message msg
	ProcessMsg(msg hub.Message)
	ThrowNewAlert(alert model.Alert) error
	// ThrowNewDatabaseAlert create and insert in the database a new NEW_DATABASE alert
	ThrowNewDatabaseAlert(dbname string, hostname string) error
	// ThrowNewServerAlert create and insert in the database a new NEW_SERVER alert
	ThrowNewServerAlert(hostname string) error
	// ThrowNewEnterpriseLicenseAlert create and insert in the database a new NEW_DATABASE alert
	ThrowNewEnterpriseLicenseAlert(hostname string) error
	// ThrowActivatedFeaturesAlert create and insert in the database a new NEW_OPTION alert
	ThrowActivatedFeaturesAlert(dbname string, hostname string, activatedFeatures []string) error
	// ThrowNoDataAlert create and insert in the database a new NO_DATA alert
	ThrowNoDataAlert(hostname string, freshnessThreshold int) error
}

// AlertService is the concrete implementation of HostDataServiceInterface. It saves data to a MongoDB database
type AlertService struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Queue that contains all messages to be processed
	Queue *hub.Hub
	// Database contains the database layer
	Database database.MongoDatabaseInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Log contains logger formatted
	Log *logrus.Logger
	// Emailer contains the emailer layer
	Emailer emailer.Emailer
}

// Init initializes the service and database
func (as *AlertService) Init(ctx context.Context, wg *sync.WaitGroup) {
	//Create a new queue
	as.Queue = hub.New()

	//Subscribe the alert-service
	sub := as.Queue.Subscribe(as.Config.AlertService.QueueBufferSize, model.TopicHostDataInsertion, model.TopicAlertInsertion)

	wg.Add(1)

	go func(s hub.Subscription) {
		as.Log.Info("Start alert-service/queue")
		for msg := range s.Receiver {
			as.ProcessMsg(msg)
		}
		as.Log.Info("Stop alert-service/queue")

		wg.Done()
	}(sub)

	go func() {
		<-ctx.Done()
		as.Queue.Close()
	}()
}

// AlertInsertion inserts an alert insertion in the queue
func (as *AlertService) AlertInsertion(alr model.Alert) error {
	as.Queue.Publish(hub.Message{
		Name: model.TopicAlertInsertion,
		Fields: hub.Fields{
			"alert": alr,
		},
	})
	return nil
}

// ProcessMsg processes the message msg
func (as *AlertService) ProcessMsg(msg hub.Message) {
	if as.Config.AlertService.LogMessages {
		as.Log.Infof("RECEIVED EVENT %s: %s", msg.Topic(), utils.ToJSON(msg.Fields))
	}

	switch msg.Topic() {
	case model.TopicAlertInsertion:
		as.ProcessAlertInsertion(msg.Fields)
	default:
		as.Log.Warnf("Received message with unknown topic: %s", msg.Topic())
	}
}

// ProcessAlertInsertion processes the alert insertion event
func (as *AlertService) ProcessAlertInsertion(params hub.Fields) {
	alert := params["alert"].(model.Alert)

	//Create the subject and message
	var subject string
	var message string
	if val, ok := alert.OtherInfo["hostname"]; ok {
		subject = fmt.Sprintf("%s %s on %s", alert.AlertSeverity, alert.Description, val)
		message = fmt.Sprintf("Date: %s\nSeverity: %s\nHost: %s\nCode: %s\n%s", alert.Date, alert.AlertSeverity, val, alert.AlertCode, alert.Description)
	} else {
		subject = fmt.Sprintf("%s %s", alert.AlertSeverity, alert.Description)
		message = fmt.Sprintf("Date: %s\nSeverity: %s\nCode: %s\n%s", alert.Date, alert.AlertSeverity, alert.AlertCode, alert.Description)
	}

	// Send the email
	err := as.Emailer.SendEmail(subject, message, as.Config.AlertService.Emailer.To)
	if err != nil {
		as.Log.Error(err)
		return
	}
}
