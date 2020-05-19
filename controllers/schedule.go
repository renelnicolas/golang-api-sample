package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"ohmytech.io/platform/models"
	"ohmytech.io/platform/repositories"
	"ohmytech.io/platform/validators"
)

// ScheduleController :
type ScheduleController struct {
}

// Save :
func (h ScheduleController) Save(w http.ResponseWriter, r *http.Request) {
	matchVars := mux.Vars(r)

	queue, ok := matchVars["queue"]
	if !ok {
		ErrorResponse(w, http.StatusNotAcceptable, "QueueName not found", nil)
		return
	}

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	scheduler := models.Scheduler{
		Queue: queue,
	}

	filter := models.QueryFilter{User: cu.UserToken.User, Extras: scheduler}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	entity, err := validators.ValidateScheduleQueue(filter, body, true)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	result, err := entityResponse(repositories.ScheduleRepository{}, entity, filter, r.Method)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	mconv, err := json.Marshal(result)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mconv)
	return
}
