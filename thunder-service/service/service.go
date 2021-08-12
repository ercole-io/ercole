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
	"math/rand"
	"time"

	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/database"

	"github.com/ercole-io/ercole/v2/config"
)

// ThunderServiceInterface is a interface that wrap methods used to insert and process alert messages
type ThunderServiceInterface interface {
	// Init initializes the service
	//Init(ctx context.Context, wg *sync.WaitGroup)
	Init()
	//GetOCRecommendations get recommendation from Oracle Cloud
	GetOCRecommendations(compartmentId string) ([]model.Recommendation, error)
	//GetOCRecommendations get recommendation from Oracle Cloud
	GetOCRecommendationsWithCategory(compartmentId string) ([]model.RecommendationWithCategory, error)

	// ProcessMsg processes the message msg
	//ProcessMsg(msg hub.Message)
	//ThrowNewAlert(alert model.Alert) error
	// ThrowNewDatabaseAlert create and insert in the database a new NEW_DATABASE alert
	//ThrowNewDatabaseAlert(dbname string, hostname string) error
	// ThrowNewServerAlert create and insert in the database a new NEW_SERVER alert
	//ThrowNewServerAlert(hostname string) error
	// ThrowNewEnterpriseLicenseAlert create and insert in the database a new NEW_DATABASE alert
	//ThrowNewEnterpriseLicenseAlert(hostname string) error
	// ThrowActivatedFeaturesAlert create and insert in the database a new NEW_OPTION alert
	//ThrowActivatedFeaturesAlert(dbname string, hostname string, activatedFeatures []string) error
	// ThrowNoDataAlert create and insert in the database a new NO_DATA alert
	//ThrowNoDataAlert(hostname string, freshnessThreshold int) error
}

// AlertService is the concrete implementation of HostDataServiceInterface. It saves data to a MongoDB database
type ThunderService struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Queue that contains all messages to be processed
	//Queue *hub.Hub
	// Database contains the database layer
	Database database.MongoDatabaseInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Log contains logger formatted
	Log logger.Logger
	// Emailer contains the emailer layer
	//Emailer emailer.Emailer
	// Random contains the generator used to generate colors
	Random *rand.Rand
}

// Init initializes the service and database
func (as *ThunderService) Init() {
	as.Random = rand.New(rand.NewSource(as.TimeNow().UnixNano()))
}
