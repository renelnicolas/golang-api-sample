package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"ohmytech.io/platform/models"
	"ohmytech.io/platform/validators"
)

// AuthenticationController :
type AuthenticationController struct {
}

// SignIn :
func (h AuthenticationController) SignIn(w http.ResponseWriter, r *http.Request) {
	userConnection := models.NewUserConnection(r.RemoteAddr, r.Header)

	rb, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.Unmarshal(rb, &userConnection)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	userSignIn, err := validators.Authentication(userConnection)
	if nil != err {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	mconv, err := json.Marshal(userSignIn)
	if nil != err {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mconv)
	return
}

// SignUp :
func (h AuthenticationController) SignUp(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusNotImplemented)
	return
}

// ForgotPassword :
func (h AuthenticationController) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusNotImplemented)
	return
}

// ResetPassword :
func (h AuthenticationController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusNotImplemented)
	return
}

// LogoutCheck :
func (h AuthenticationController) LogoutCheck(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusNotImplemented)
	return
}

// Logout :
func (h AuthenticationController) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusNotImplemented)
	return
}
