package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/zidariu-sabin/femProject/internal/app"
	"github.com/zidariu-sabin/femProject/internal/routes"
)

// entry point of the application
func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "go backend port")
	flag.Parse()
	//-port *value* will set the port we will run from to value
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	//at the end of execution close the datbase connection
	defer app.DB.Close()

	router := routes.SetupRoutes(app)

	server := &http.Server{
		//print function that returns a value
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("App is successfully running on port %d \n", port)

	err = server.ListenAndServe()

	if err != nil {
		server.ErrorLog.Fatal()
		app.Logger.Fatal()
	}
}
