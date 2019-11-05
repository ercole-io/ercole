// Package controller contains structs and methods used to provide endpoints for storing hostdata informations
package controller

import (
	"net/http"
	"time"

	"github.com/amreo/ercole-services/alert-service/service"
	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/utils"
	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AlertQueueControllerInterface is a interface that wrap methods used to inserting events in the queue
type AlertQueueControllerInterface interface {
	// HostDataInsertion insert the event HostDataInsertion with the id in the queue
	HostDataInsertion(w http.ResponseWriter, r *http.Request)
	// AuthenticateMiddleware return the middleware used to authenticate users
	AuthenticateMiddleware() func(http.Handler) http.Handler
}

// AlertQueueController is the struct used to handle the requests from agents and contains the concrete implementation of AlertQueueControllerInterface
type AlertQueueController struct {
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Service contains the underlying service used to perform various logical and store operations
	Service service.AlertServiceInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
}

// AuthenticateMiddleware return the middleware used to authenticate (request) users
func (ctrl *AlertQueueController) AuthenticateMiddleware() func(http.Handler) http.Handler {
	return httpauth.SimpleBasicAuth(ctrl.Config.AlertService.PublisherUsername, ctrl.Config.AlertService.PublisherPassword)
}

// HostDataInsertion insert the event HostDataInsertion with the id in the queue
func (ctrl *AlertQueueController) HostDataInsertion(w http.ResponseWriter, r *http.Request) {
	var id primitive.ObjectID
	var err error

	//Get the id from the path variable
	if id, err = primitive.ObjectIDFromHex(mux.Vars(r)["id"]); err != nil {
		utils.WriteAndLogError(w, http.StatusUnprocessableEntity, utils.NewAdvancedErrorPtr(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	//Insert the event
	if err := ctrl.Service.HostDataInsertion(id); err != nil {
		utils.WriteAndLogError(w, http.StatusUnprocessableEntity, err)
		return
	}
}
