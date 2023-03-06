package middleware

import (
	"net/http"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/gorilla/context"
)

func Admin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, exists := context.GetOk(r, "user")
		if !exists {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user := claims.(model.User)
		if !user.IsAdmin() {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
