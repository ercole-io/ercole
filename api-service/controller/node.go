package controller

import (
	"net/http"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func (ctrl *APIController) GetNodes(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user")
	if user == nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnauthorized, nil)
		return
	}

	nodes, err := ctrl.Service.GetNodes(user.(dto.AllowedUser).Groups)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, dto.ToNodes(nodes))
}

func (ctrl *APIController) GetNode(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	node, err := ctrl.Service.GetNode(name)
	if err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, node)
}

func (ctrl *APIController) AddNode(w http.ResponseWriter, r *http.Request) {
	node := &model.Node{}

	if err := utils.Decode(r.Body, node); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	if err := ctrl.Service.AddNode(*node); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ctrl *APIController) UpdateNode(w http.ResponseWriter, r *http.Request) {
	node := &model.Node{}

	if err := utils.Decode(r.Body, node); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, err)
		return
	}

	name := mux.Vars(r)["name"]
	if name != node.Name {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusBadRequest, nil)
		return
	}

	if err := ctrl.Service.UpdateNode(*node); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ctrl *APIController) RemoveNode(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	if err := ctrl.Service.RemoveNode(name); err != nil {
		utils.WriteAndLogError(ctrl.Log, w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
