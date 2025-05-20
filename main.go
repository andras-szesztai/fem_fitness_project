package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/andras-szesztai/fem_fitness_project/internal/app"
	"github.com/andras-szesztai/fem_fitness_project/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	app.Logger.Println("Starting application...")

	router := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      router,
	}

	app.Logger.Printf("Server is running on port %d", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}

}
