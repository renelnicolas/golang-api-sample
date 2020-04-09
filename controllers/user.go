package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"ohmytech.io/platform/models"
	"ohmytech.io/platform/repositories"
)

// UserController :
type UserController struct {
}

// List :
func (h UserController) List(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	filter := models.QueryFilter{Offset: 0, Limit: maxLimit, Order: `email`, Sort: `asc`, User: cu.UserToken.User}

	extractFilter(&filter, queryValues, maxLimit)

	result, err := paginationResponse(repositories.UserRepository{}, filter)
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

// Profil :
func (h UserController) Profil(w http.ResponseWriter, r *http.Request) {
	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	entity := cu.UserToken.ToUser()

	repo := repositories.UserRepository{}

	err := repo.FindOneByEmail(&entity)
	if nil != err {
		ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	mconv, err := json.Marshal(entity)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mconv)
	return
}

// Edit :
func (h UserController) Edit(w http.ResponseWriter, r *http.Request) {
	matchVars := mux.Vars(r)

	ID, err := extractMuxVars(matchVars)

	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	filter := models.QueryFilter{User: cu.UserToken.User}

	entity := models.User{ID: ID}

	repo := repositories.UserRepository{}

	result, err := repo.FindOne(entity, filter)
	if nil != err {
		ErrorResponse(w, http.StatusNotFound, "")
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

// Save :
func (h UserController) Save(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusNotImplemented)
	return
}

// Update :
func (h UserController) Update(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusNotImplemented)
	return
}

// Delete :
func (h UserController) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusNotImplemented)
	return
}