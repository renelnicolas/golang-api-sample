package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"ohmytech.io/platform/models"
	"ohmytech.io/platform/repositories"
)

var (
	maxLimit = int64(100)
)

// ControllerControllerer :
type ControllerControllerer interface {
	List(w http.ResponseWriter, r *http.Request)
	Edit(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// JSONResponseList :
type JSONResponseList struct {
	Entities interface{} `json:"entities"`
	Counter  int64       `json:"counter"`
}

// JSONResponse :
type JSONResponse struct {
	Entity interface{} `json:"entity"`
}

// JSONResponseError :
type JSONResponseError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

// ErrorResponse :
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	re := JSONResponseError{
		Status: status,
		Error:  message,
	}

	mconv, _ := json.Marshal(re)

	log.Printf(">> Error : %s", message)

	w.WriteHeader(status)
	w.Write(mconv)
}

// paginationResponse :
func paginationResponse(repo repositories.Repositoryer, filter models.QueryFilter) (*JSONResponseList, error) {
	counter, err := repo.Count(filter)
	if nil != err {
		return nil, err
	}

	entities, err := repo.FindAll(filter)
	if nil != err {
		return nil, err
	}

	result := JSONResponseList{
		Counter:  counter,
		Entities: entities,
	}

	return &result, nil
}

// entityResponse :
func entityResponse(repo repositories.Repositoryer, entity interface{}, filter models.QueryFilter, isNew bool) (*JSONResponse, error) {
	var (
		_entity interface{}
		err     error
	)

	if isNew {
		_entity, err = repo.Create(entity, filter)
	} else {
		_entity, err = repo.FindOne(entity, filter)
	}

	if nil != err {
		return nil, err
	}

	result := JSONResponse{
		Entity: _entity,
	}

	return &result, nil
}

func extractFilter(filter *models.QueryFilter, queryValues url.Values, paginationMaxLimit int64) {
	if 0 < len(queryValues["offset"]) {
		n, err := strconv.ParseInt(queryValues["offset"][0], 10, 64)
		if nil == err {
			filter.Offset = n
		}
	}

	if 0 < len(queryValues["limit"]) {
		n, err := strconv.ParseInt(queryValues["limit"][0], 10, 64)
		if nil == err {
			if paginationMaxLimit >= n {
				filter.Limit = n
			}
		}
	}

	if 0 < len(queryValues["order"]) {
		filter.Order = queryValues["order"][0]
	}

	if 0 < len(queryValues["sort"]) {
		filter.Sort = queryValues["sort"][0]
	}

	if 0 < len(queryValues["search"]) {
		filter.Search = queryValues["search"][0]
	}

	log.Println("extractFilter", filter)
}

func extractIDFromMuxVars(matchVars map[string]string) (int64, error) {
	ID, ok := matchVars["id"]
	if !ok {
		return 0, errors.New("ID cannot be empty")
	}

	entityID, err := strconv.ParseInt(ID, 10, 64)
	if nil != err {
		return 0, fmt.Errorf("Conversion error: %s", err.Error())
	}

	return entityID, nil
}
