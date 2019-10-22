package middlewares

import (
	"errors"
	"net/http"

	auth "github.com/Wallruzz9114/gen_api/api/auth"
	responses "github.com/Wallruzz9114/gen_api/api/responses"
)

// SetMiddlewareJSON ...
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// SetMiddlewareAuthentication ...
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err == nil {
			next(w, r)
			return
		}

		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	}
}
