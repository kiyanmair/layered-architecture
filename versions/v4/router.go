package v4

import (
	"net/http"

	"github.com/go-chi/chi"
)

// CreateRouter creates a new router and registers all routes.
func CreateRouter(svc Service) http.Handler {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", Ping)

		r.Route("/users", func(r chi.Router) {
			r.Post("/register", RegisterUserHandler(svc))
		})
	})

	return r
}
