package controller

import (
	"net/http"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/context"
)

func (ctrl *APIController) GetNodes(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user")
	if user == nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnauthorized, nil)
		return
	}

	nodes, err := ctrl.Service.GetNodes(user.(*model.User))
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, nodes)
}
