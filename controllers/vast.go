package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"ohmytech.io/platform/helpers"
	"ohmytech.io/platform/validators"
	"ohmytech.io/platform/vast"
)

// VastController :
type VastController struct {
}

// Converter :
func (h VastController) Converter(w http.ResponseWriter, r *http.Request) {
	var (
		mconv []byte
		err   error
	)

	convertTo := "json"
	queryValues := r.URL.Query()

	if 0 == len(queryValues["vastURL"]) || 0 == len(queryValues["vastURL"][0]) {
		ErrorResponse(w, http.StatusNotAcceptable, "vastURL query string cannot be empty", nil)
		return
	}

	vastURL := queryValues["vastURL"][0]

	if !validators.ValidateURL(vastURL) {
		ErrorResponse(w, http.StatusNotAcceptable, "It's not an URL.", nil)
		return
	}

	u, err := url.Parse(vastURL)
	if err != nil {
		log.Fatal(err)
		ErrorResponse(w, http.StatusInternalServerError, "URL cannot be parsed.", nil)
		return
	}

	ssp := helpers.ServerSideParameters{
		XForwardFor: helpers.GetIP(r.RemoteAddr, r.Header),
		UserAgent:   r.UserAgent(),
		Referer:     r.Referer(),
		URL:         r.URL.String(),
		URLQuery:    r.URL.Query(),
	}

	strVast, err := helpers.ServerSideCall(u.String(), ssp)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	vst := vast.VAST{}

	err = xml.Unmarshal(strVast, &vst)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if 0 != len(queryValues["convertTo"]) {
		convertTo = queryValues["convertTo"][0]
	}

	switch convertTo {
	case "json":
		mconv, err = json.Marshal(vst)
	case "xml":
		mconv, err = xml.Marshal(vst)
	}

	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mconv)
	return
}

// Analyser :
func (h VastController) Analyser(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()

	if 0 == len(queryValues["vastURL"]) || 0 == len(queryValues["vastURL"][0]) {
		ErrorResponse(w, http.StatusNotAcceptable, "vastURL query string cannot be empty", nil)
		return
	}

	vastURL := queryValues["vastURL"][0]

	if !validators.ValidateURL(vastURL) {
		ErrorResponse(w, http.StatusNotAcceptable, "It's not an URL.", nil)
		return
	}

	u, err := url.Parse(vastURL)
	if err != nil {
		log.Fatal(err)
		ErrorResponse(w, http.StatusInternalServerError, "URL cannot be parsed.", nil)
		return
	}

	ssp := helpers.ServerSideParameters{
		XForwardFor: helpers.GetIP(r.RemoteAddr, r.Header),
		UserAgent:   r.UserAgent(),
		Referer:     r.Referer(),
		URL:         r.URL.String(),
		URLQuery:    r.URL.Query(),
	}

	fmt.Printf("URL : %s\n", u.String())

	strVast, err := helpers.ServerSideCall(u.String(), ssp)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	fmt.Printf("RESULT : %s\n", strVast)

	vst := vast.VAST{}

	err = xml.Unmarshal(strVast, &vst)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	fmt.Printf("RESULT : %v\n", vst)

	mconv, err := json.Marshal(vst)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mconv)
	return
}

// Parser :
func (h VastController) Parser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}
