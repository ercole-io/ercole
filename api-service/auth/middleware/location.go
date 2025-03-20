package middleware

import (
	"net/http"
	"strings"

	"github.com/ercole-io/ercole/v2/api-service/service"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func Location(service service.APIServiceInterface) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, exists := context.GetOk(r, "user")
			if !exists {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			user, ok := claims.(model.User)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if user.IsAdmin() {
				h.ServeHTTP(w, r)
				return
			}

			locations, err := service.ListLocations(user)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if utils.ContainsI(locations, model.AllLocation) {
				h.ServeHTTP(w, r)
				return
			}

			location := r.URL.Query().Get("location")
			if location == "" || strings.EqualFold(location, model.AllLocation) {
				query := r.URL.Query()
				query.Set("location", strings.Join(locations, ","))
				r.URL.RawQuery = query.Encode()

				h.ServeHTTP(w, r)
				return
			}

			if utils.ContainsI(locations, location) {
				h.ServeHTTP(w, r)
				return
			}

			splittedLocation := strings.Split(location, ",")
			if utils.ContainsSomeI(locations, splittedLocation...) {
				h.ServeHTTP(w, r)
				return
			}

			w.WriteHeader(http.StatusForbidden)
		})
	}
}
