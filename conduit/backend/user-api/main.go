package main

import (
	"backend/internal"
	"backend/internal/database"
	"backend/user-api/handlers/users"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	userAPILogger := log.New(os.Stdout, "user-api | ", log.LstdFlags)

	validation := internal.NewValidation()

	userHandler := users.NewUsers(userAPILogger, validation)

	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", userHandler.ListAll)
	getRouter.HandleFunc("/users/{id:[0-9]+}", userHandler.ListSingle)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/users", userHandler.Update)
	putRouter.Use(userHandler.MiddlewareValidateUser)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/users", userHandler.Create)
	postRouter.Use(userHandler.MiddlewareValidateUser)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/users/{id:[0-9]+}", userHandler.Delete)

	// handler for documentation
	//opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	//sh := middleware.Redoc(opts, nil)

	//getRouter.Handle("/docs", sh)
	//getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	//getRouter.Handle("/swagger.json", http.FileServer(http.Dir("./")))

	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	srv := &http.Server{
		Addr:         ":9090",               // configure the bind address
		Handler:      corsHandler(serveMux), // set the default handlers
		ErrorLog:     userAPILogger,         // set the severAPILogger for the servers
		ReadTimeout:  5 * time.Second,       // max time to read request from the client
		WriteTimeout: 10 * time.Second,      // max time to write response to the client
		IdleTimeout:  120 * time.Second,     // max time for connections using TCP Keep-Alive
	}

	// start the servers
	go func() {
		userAPILogger.Println("Starting servers on port 9090")

		err := srv.ListenAndServe()
		if err != nil {
			userAPILogger.Printf("Error starting servers: %s\n", err)
			os.Exit(1)
		}
	}()
	//Make sure the db tables and model of the severs match up
	database.AutoMigrateDB()

	// trap sigterm or interrupt and gracefully shutdown the servers
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the servers, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Fatal(srv.ListenAndServe())
}
