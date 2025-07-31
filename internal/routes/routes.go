package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/zidariu-sabin/femProject/internal/app"
)

// a struct that works as a multiplexor to pare routes and their parameters
func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		router.Get("/workout/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleGetWorkoutById))
		router.Post("/workout", app.Middleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
		router.Put("/workout/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkoutById))
		router.Delete("/workout/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkoutById))
	})

	router.Get("/health", app.HealthCheck)

	router.Post("/user", app.UserHandler.HandleRegisterUser)
	router.Get("/user", app.UserHandler.HandleGetUserByUsername)
	router.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)

	return router
}
