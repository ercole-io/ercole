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

// Package controller contains structs and methods used to provide endpoints for storing hostdata informations
package controller

import (
	"net/http"
	"time"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/thunder-service/service"
)

// ThunderControllerInterface is a interface that wrap methods used to inserting events in the queue
type ThunderControllerInterface interface {

	// GetOciRecommendations get recommendations from Oracle Cloud Infrastructure
	GetOciRecommendations(w http.ResponseWriter, r *http.Request)
	// Get Configuration profiles for Oracle Cloud Access
	GetOciProfiles(w http.ResponseWriter, r *http.Request)
	// Add a new Configuration profile for Oracle Cloud Access
	AddOciProfile(w http.ResponseWriter, r *http.Request)
	// Update an existing Configuration profile for Oracle Cloud Access
	UpdateOciProfile(w http.ResponseWriter, r *http.Request)
	// Delete an existing Configuration profile for Oracle Cloud Access
	DeleteOciProfile(w http.ResponseWriter, r *http.Request)
	// GetOciUnusedLoadbalancers get recommendations from Oracle Cloud Infrastructure about Load Balancer health
	GetOciUnusedLoadbalancers(w http.ResponseWriter, r *http.Request)
	// GetOciComputeInstancesIdle get recommendations from Oracle Cloud Infrastructure about Idle Instances
	GetOciComputeInstancesIdle(w http.ResponseWriter, r *http.Request)
	// GetOciBlockStorageRightsizing get recommendations from Oracle Cloud Infrastructure about Optimizable Block Storage
	GetOciBlockStorageRightsizing(w http.ResponseWriter, r *http.Request)
	// GetOciUnusedStorage get recommendations from Oracle Cloud Infrastructure about Unused Storage
	GetOciUnusedStorage(w http.ResponseWriter, r *http.Request)
	// GetOciOldSnapshotDecommissioning get recommendations from Oracle Cloud Infrastructure about old snapshot
	GetOciOldSnapshotDecommissioning(w http.ResponseWriter, r *http.Request)
	// GetOciComputeInstanceRightsizing get recommendations from Oracle Cloud Infrastructure about Underutilized Instances
	GetOciComputeInstanceRightsizing(w http.ResponseWriter, r *http.Request)
}

// ThunderController is the struct used to handle the requests from agents and contains the concrete implementation of ThunderControllerInterface
type ThunderController struct {
	Config  config.Configuration
	Service service.ThunderServiceInterface
	TimeNow func() time.Time
	Log     logger.Logger
}
