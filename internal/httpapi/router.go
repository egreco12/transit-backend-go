package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(arrivals *ArrivalsHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(api chi.Router) {
		api.Route("/stops", func(sr chi.Router) {
			sr.Get("/{stopID}/arrivals", arrivals.GetArrivals)
		})
	})

	return r
}
