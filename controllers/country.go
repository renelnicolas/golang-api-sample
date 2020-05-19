package controllers

import (
	"encoding/json"
	"net/http"

	"ohmytech.io/platform/models"
	"ohmytech.io/platform/repositories"
)

// CountryController :
type CountryController struct {
}

// List :
func (h CountryController) List(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	cval := r.Context().Value(models.ContextUserKey)
	cu := cval.(models.UserClaims)

	filter := models.QueryFilter{Order: `name`, Sort: `asc`, User: cu.UserToken.User}

	extractFilter(&filter, queryValues, maxLimit)

	result, err := paginationResponse(repositories.CountryRepository{}, filter)
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
