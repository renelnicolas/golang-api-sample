package middlewares

import (
	"context"
	"net/http"

	"ohmytech.io/platform/controllers"
	"ohmytech.io/platform/models"
	"ohmytech.io/platform/validators"
)

var (
	signingKey = []byte("signingKey")
)

// JwtAuthentication :
func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Cache-Control", "max-age=0, no-store, no-cache, must-revalidate")

		userClaims := models.UserClaims{}
		authorization := r.Header.Get("Authorization")

		err := validators.JwtAuthentication(authorization, &userClaims)
		if nil != err {
			controllers.ErrorResponse(w, http.StatusForbidden, err.Error(), nil)
			return
		}

		// Everything went well, proceed with the request and set the caller to the user
		// retrieved from the parsed token
		ctx := context.WithValue(r.Context(), models.ContextUserKey, userClaims)

		next.ServeHTTP(w, r.WithContext(ctx)) //proceed in the middleware chain!
	})
}
