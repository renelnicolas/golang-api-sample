package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"ohmytech.io/platform/amq"
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
		ErrorResponse(w, http.StatusNotAcceptable, "")
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
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	entity, err := validators.ValidateScheduleQueue(filter, body, true)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := entityResponse(repositories.ScheduleRepository{}, entity, filter, true)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	mconv, err := json.Marshal(result)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mconv)
	return
}

// Schedule :
func (h ScheduleController) Schedule(w http.ResponseWriter, r *http.Request) {
	matchVars := mux.Vars(r)

	ID, err := extractIDFromMuxVars(matchVars)
	if nil != err {
		ErrorResponse(w, http.StatusNotAcceptable, "Id not found")
		return
	}

	queue, ok := matchVars["queue"]
	if !ok {
		ErrorResponse(w, http.StatusNotAcceptable, "QueueName not found")
		return
	}

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	scheduler := models.Scheduler{
		ID:    ID,
		Queue: queue,
	}

	filter := models.QueryFilter{User: cu.UserToken.User, Extras: scheduler}

	result, err := entityResponse(repositories.ScheduleRepository{}, scheduler, filter, false)
	if nil != err {
		ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	err = amq.Publish(queue, result, nil)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, "Cannot publish schedule request")
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}
