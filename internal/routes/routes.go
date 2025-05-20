package routes

import (
	"github.com/andras-szesztai/fem_fitness_project/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", app.HealthCheck)
	})

	return router
}
