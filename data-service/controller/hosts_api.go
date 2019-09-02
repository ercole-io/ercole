package controller

import (
	"encoding/json"
	"net/http"

	"github.com/amreo/ercole-services/utils"

	"github.com/amreo/ercole-services/data-service/service"
	"github.com/amreo/ercole-services/model"
	"github.com/goji/httpauth"
)

// AuthenticateMiddleware return the middleware used to authenticate (request) users
func (this *HostDataController) AuthenticateMiddleware() func(http.Handler) http.Handler {
	return httpauth.SimpleBasicAuth(this.Config.HttpServer.AgentUsername, this.Config.HttpServer.AgentPassword)
}

// UpdateHostInfo update the informations about a host using the HostData in the request
func (this *HostDataController) UpdateHostInfo(w http.ResponseWriter, r *http.Request) {
	//Parse the hostdata from the request
	var hostData model.HostData
	if err := json.NewDecoder(r.Body).Decode(&hostData); err != nil {
		WriteAndLogError(w, http.StatusUnprocessableEntity, utils.NewAdvancedError(err, http.StatusText(http.StatusUnprocessableEntity)))
		return
	}

	//Save the HostData
	id, err := service.SaveHostData(hostData)
	if err != nil {
		WriteAndLogError(w, http.StatusInternalServerError, err)
		return
	}

	//Write the created id
	utils.WriteJSONResponse(w, http.StatusOK, id)
}
