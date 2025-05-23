package routes

import (
	"github.com/andras-szesztai/fem_fitness_project/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", app.HealthCheck)

	router.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(app.Middleware.Authenticate)
			r.Route("/workouts", func(r chi.Router) {
				r.Get("/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleGetWorkout))
				r.Post("/", app.Middleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
				r.Put("/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkout))
				r.Delete("/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkout))
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/register", app.UserHandler.HandleRegisterUser)
		})

		r.Route("/tokens", func(r chi.Router) {
			r.Post("/", app.TokenHandler.HandleCreateToken)
		})
	})

	return router
}
