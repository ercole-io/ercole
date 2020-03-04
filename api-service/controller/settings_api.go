package controller

import (
	"net/http"

	"github.com/amreo/ercole-services/utils"
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
