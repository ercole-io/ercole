package controller

import (
	"net/http"

	"github.com/ercole-io/ercole/utils"
)

// GetDefaultDatabaseTags return the default list of database tags from configuration
func (ctrl *APIController) GetDefaultDatabaseTags(w http.ResponseWriter, r *http.Request) {
	tags, err := ctrl.Service.GetDefaultDatabaseTags()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, tags)
}

// GetErcoleFeatures return a map of active/inactive features
func (ctrl *APIController) GetErcoleFeatures(w http.ResponseWriter, r *http.Request) {
	data, err := ctrl.Service.GetErcoleFeatures()
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, data)
}
