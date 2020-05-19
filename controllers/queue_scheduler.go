package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"ohmytech.io/platform/amq"
	"ohmytech.io/platform/collections"
	"ohmytech.io/platform/models"
	"ohmytech.io/platform/repositories"
	"ohmytech.io/platform/validators"
)

// QueueSchedulerController :
type QueueSchedulerController struct {
}

// List :
func (h QueueSchedulerController) List(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	filters := models.QueryFilter{Offset: 0, Limit: maxLimit, Order: `name`, Sort: `asc`, User: cu.UserToken.User}

	extractFilter(&filters, queryValues, maxLimit)

	result, err := paginationResponse(repositories.QueueSchedulerRepository{}, filters)
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
func (h QueueSchedulerController) Edit(w http.ResponseWriter, r *http.Request) {
	matchVars := mux.Vars(r)

	ID, err := extractIDFromMuxVars(matchVars)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	filter := models.QueryFilter{User: cu.UserToken.User}

	entity := models.QueueScheduler{ID: ID}

	result, err := entityResponse(repositories.QueueSchedulerRepository{}, entity, filter, r.Method)
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
func (h QueueSchedulerController) Save(w http.ResponseWriter, r *http.Request) {
	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	if !cu.IsUser() {
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

	entity, err := validators.ValidateQueueScheduler(filter, body)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	result, err := entityResponse(repositories.QueueSchedulerRepository{}, entity, filter, r.Method)
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

// Update :
func (h QueueSchedulerController) Update(w http.ResponseWriter, r *http.Request) {
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

	entity, err := validators.ValidateQueueScheduler(filter, body)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if ID != entity.ID {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	result, err := entityResponse(repositories.QueueSchedulerRepository{}, *entity, filter, r.Method)
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
func (h QueueSchedulerController) Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

// Schedule :
func (h QueueSchedulerController) Schedule(w http.ResponseWriter, r *http.Request) {
	matchVars := mux.Vars(r)

	ID, err := extractIDFromMuxVars(matchVars)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	filter := models.QueryFilter{User: cu.UserToken.User}

	entity := models.QueueScheduler{ID: ID}

	repoqs := repositories.QueueSchedulerRepository{}

	_entity, err := repoqs.FindOne(entity, filter)
	if nil != err {
		ErrorResponse(w, http.StatusNotFound, "", err)
		return
	}

	entity = _entity.(models.QueueScheduler)

	queueSchedulerHistory := models.QueueSchedulerHistory{
		QueueSchedulerID: entity.ID,
	}

	repoqsh := repositories.QueueSchedulerHistoryRepository{}

	_qsh, err := repoqsh.Create(queueSchedulerHistory, filter)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, "", err)
		return
	}

	qsh := _qsh.(models.QueueSchedulerHistory)

	entity.QueueSchedulerHistory = &qsh

	err = amq.Publish(entity.QueueType.Name, entity, nil)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, "Cannot publish schedule request", nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}

// History :
func (h QueueSchedulerController) History(w http.ResponseWriter, r *http.Request) {
	matchVars := mux.Vars(r)

	scheduleID, ok := matchVars["scheduleId"]
	if !ok {
		ErrorResponse(w, http.StatusNotAcceptable, "scheduleID cannot be empty", nil)
		return
	}

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	if !cu.IsUser() {
		ErrorResponse(w, http.StatusForbidden, "", nil)
		return
	}

	queryValues := r.URL.Query()

	filters := models.QueryFilter{Offset: 0, Limit: maxLimit, Order: `name`, Sort: `asc`, User: cu.UserToken.User}

	extractFilter(&filters, queryValues, maxLimit)

	col := collections.ResumeCollection{}

	result, err := col.FindAll(map[string]interface{}{
		"scheduleId": scheduleID,
		"companyId":  cu.UserToken.User.Company.ID,
	})
	if nil != err {
		ErrorResponse(w, http.StatusNotFound, "No match", nil)
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

// Resume :
func (h QueueSchedulerController) Resume(w http.ResponseWriter, r *http.Request) {
	matchVars := mux.Vars(r)

	scheduleID, ok := matchVars["scheduleId"]
	if !ok {
		ErrorResponse(w, http.StatusNotAcceptable, "scheduleID cannot be empty", nil)
		return
	}

	workID, ok := matchVars["workId"]
	if !ok {
		ErrorResponse(w, http.StatusNotAcceptable, "workID cannot be empty", nil)
		return
	}

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	if !cu.IsUser() {
		ErrorResponse(w, http.StatusForbidden, "", nil)
		return
	}

	queryValues := r.URL.Query()

	filters := models.QueryFilter{Offset: 0, Limit: maxLimit, Order: `name`, Sort: `asc`, User: cu.UserToken.User}

	extractFilter(&filters, queryValues, maxLimit)

	col := collections.ResumeCollection{}

	result, err := col.FindOne(map[string]interface{}{
		"scheduleId": scheduleID,
		"workId":     workID,
		"companyId":  cu.UserToken.User.Company.ID,
	})
	if nil != err {
		ErrorResponse(w, http.StatusNotFound, "No match", nil)
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

// Details :
func (h QueueSchedulerController) Details(w http.ResponseWriter, r *http.Request) {
	matchVars := mux.Vars(r)

	workID, ok := matchVars["workId"]
	if !ok {
		ErrorResponse(w, http.StatusNotAcceptable, "workID cannot be empty", nil)
		return
	}

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	if !cu.IsUser() {
		ErrorResponse(w, http.StatusForbidden, "", nil)
		return
	}

	queryValues := r.URL.Query()

	filters := models.QueryFilter{Offset: 0, Limit: maxLimit, Order: `name`, Sort: `asc`, User: cu.UserToken.User}

	extractFilter(&filters, queryValues, maxLimit)

	col := collections.AnalyserCollection{}

	result, err := col.FindAll(map[string]interface{}{
		"workId":    workID,
		"companyId": cu.UserToken.User.Company.ID,
	})
	if nil != err {
		ErrorResponse(w, http.StatusNotFound, "No match", err)
		return
	}

	mconv, err := json.Marshal(result)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mconv)
	return
}
