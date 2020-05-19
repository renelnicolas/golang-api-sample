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

// QueueTypeController :
type QueueTypeController struct {
}

// List :
func (h QueueTypeController) List(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	if !cu.IsMaster() {
		ErrorResponse(w, http.StatusForbidden, "", nil)
		return
	}

	filters := models.QueryFilter{Offset: 0, Limit: maxLimit, Order: `name`, Sort: `asc`, User: cu.UserToken.User}

	extractFilter(&filters, queryValues, maxLimit)

	result, err := paginationResponse(repositories.QueueTypeRepository{}, filters)
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

// Edit :
func (h QueueTypeController) Edit(w http.ResponseWriter, r *http.Request) {
	matchVars := mux.Vars(r)

	ID, err := extractIDFromMuxVars(matchVars)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	if !cu.IsMaster() {
		ErrorResponse(w, http.StatusForbidden, "", nil)
		return
	}

	filter := models.QueryFilter{User: cu.UserToken.User}

	entity := models.QueueType{ID: ID}

	result, err := entityResponse(repositories.QueueTypeRepository{}, entity, filter, r.Method)
	if nil != err {
		ErrorResponse(w, http.StatusNotFound, "", err)
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

// Save :
func (h QueueTypeController) Save(w http.ResponseWriter, r *http.Request) {
	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	if !cu.IsMaster() {
		ErrorResponse(w, http.StatusForbidden, "", nil)
		return
	}

	filter := models.QueryFilter{User: cu.UserToken.User}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	entity, err := validators.ValidateQueueType(filter, body)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	result, err := entityResponse(repositories.QueueTypeRepository{}, *entity, filter, r.Method)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	mconv, err := json.Marshal(result)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(mconv)
	return
}

// Update :
func (h QueueTypeController) Update(w http.ResponseWriter, r *http.Request) {
	matchVars := mux.Vars(r)

	ID, err := extractIDFromMuxVars(matchVars)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	if !cu.IsMaster() {
		ErrorResponse(w, http.StatusForbidden, "", nil)
		return
	}

	filter := models.QueryFilter{User: cu.UserToken.User}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	entity, err := validators.ValidateQueueType(filter, body)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if ID != entity.ID {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	result, err := entityResponse(repositories.QueueTypeRepository{}, *entity, filter, r.Method)
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

// Delete :
func (h QueueTypeController) Delete(w http.ResponseWriter, r *http.Request) {
	matchVars := mux.Vars(r)

	_, err := extractIDFromMuxVars(matchVars)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	if !cu.IsMaster() {
		ErrorResponse(w, http.StatusForbidden, "", nil)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
	return
}
