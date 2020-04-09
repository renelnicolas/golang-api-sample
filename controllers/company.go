package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"ohmytech.io/platform/models"
	"ohmytech.io/platform/repositories"
)

// CompanyController :
type CompanyController struct {
}

// List :
func (h CompanyController) List(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	filters := models.QueryFilter{Offset: 0, Limit: maxLimit, Order: `name`, Sort: `asc`, User: cu.UserToken.User}

	extractFilter(&filters, queryValues, maxLimit)

	result, err := paginationResponse(repositories.CompanyRepository{}, filters)
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

// Edit :
func (h CompanyController) Edit(w http.ResponseWriter, r *http.Request) {
	matchVars := mux.Vars(r)

	ID, err := extractMuxVars(matchVars)

	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	filter := models.QueryFilter{User: cu.UserToken.User}

	entity := models.Company{ID: ID}

	repo := repositories.CompanyRepository{}

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
func (h CompanyController) Save(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

// Update :
func (h CompanyController) Update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

// Delete :
func (h CompanyController) Delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}
