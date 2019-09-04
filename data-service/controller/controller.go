// controller contains structs and methods used to provide endpoints for storing hostdata informations
package controller

import (
	"net/http"

	"github.com/amreo/ercole-services/data-service/service"

	"github.com/amreo/ercole-services/config"
)

// HostDataControllerInterface is a interface that wrap methods used to handle the request for HostData endpoints
type HostDataControllerInterface interface {
	// UpdateHostInfo update the informations about a host using the HostData in the request
	UpdateHostInfo(w http.ResponseWriter, r *http.Request)
	// AuthenticateMiddleware return the middleware used to authenticate users
	AuthenticateMiddleware() func(http.Handler) http.Handler
}

// HostDataController is the struct used to handle the requests from agents and contains the concrete implementation of HostDataControllerInterface
type HostDataController struct {
	// Config contains the dataservice global configuration
	// TODO: Should be removed?
	Config config.Configuration
	// Service contains the underlying service used to perform various logical and store operations
	Service service.HostDataServiceInterface
}
