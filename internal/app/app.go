package app

//resources used within the application
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zidariu-sabin/femProject/internal/api"
	"github.com/zidariu-sabin/femProject/internal/middleware"
	"github.com/zidariu-sabin/femProject/internal/store"
	"github.com/zidariu-sabin/femProject/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	Middleware     middleware.UserMiddleware
	DB             *sql.DB
}

// implementing logging
func NewApplication() (*Application, error) {
	pgDB, err := store.Open()

	if err != nil {
		return nil, err
	}

	// "." means we are going to use the base of the directory
	err = store.MigrateFs(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//stores
	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	//handlers
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		Middleware:     middlewareHandler,
		DB:             pgDB,
	}

	return app, nil
}

// we want to have the request as a pointer in order to persist the incoming data in comparison to the data that we send back

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
