package v1

import (
	"net/http"

	"github.com/go-chi/chi"
)

// CreateRouter creates a new router and registers all routes.
func CreateRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", Ping)

		r.Route("/users", func(r chi.Router) {
			r.Post("/register", RegisterUserHandler)
		})
	})

	return r
}
