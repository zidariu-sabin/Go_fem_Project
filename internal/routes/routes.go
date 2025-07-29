package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/zidariu-sabin/femProject/internal/app"
)

// a struct that works as a multiplexor to pare routes and their parameters
func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", app.HealthCheck)

	router.Get("/workout/{id}", app.WorkoutHandler.HandleGetWorkoutById)
	router.Post("/workout-create", app.WorkoutHandler.HandleCreateWorkout)
	router.Put("/workout-update/{id}", app.WorkoutHandler.HandleUpdateWorkoutById)
	router.Delete("/workout-delete/{id}", app.WorkoutHandler.HandleDeleteWorkoutById)

	router.Post("/user-create", app.UserHandler.HandleRegisterUser)
	router.Get("/user", app.UserHandler.HandleGetUserByUsername)
	router.Post("/tokens/authentification", app.TokenHandler.HandleCreateToken)

	return router
}
