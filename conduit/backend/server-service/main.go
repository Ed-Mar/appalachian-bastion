package main

import (
	"backend/internal"
	"backend/server-service/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	severAPILogger := log.New(os.Stdout, "server-service | ", log.LstdFlags)

	validation := internal.NewValidation()

	serverHandler := handlers.NewServers(severAPILogger, validation)

	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/servers", serverHandler.ListCollection)
	getRouter.HandleFunc("/servers/{serverID}", serverHandler.ListSingleton)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/server", serverHandler.UpdateSingleton)
	putRouter.Use(serverHandler.MiddlewareValidateServer)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/server", serverHandler.CreateSingleton)
	postRouter.Use(serverHandler.MiddlewareValidateServer)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/servers/{serverID}", serverHandler.DeleteSingleton)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	getRouter.Handle("/swagger.json", http.FileServer(http.Dir("./")))

	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	srv := &http.Server{
		Addr:         ":9292",               // configure the bind address
		Handler:      corsHandler(serveMux), // set the default handlers
		ErrorLog:     severAPILogger,        // set the severAPILogger for the servers
		ReadTimeout:  5 * time.Second,       // max time to read request from the client
		WriteTimeout: 10 * time.Second,      // max time to write response to the client
		IdleTimeout:  120 * time.Second,     // max time for connections using TCP Keep-Alive
	}

	// start the servers
	go func() {
		severAPILogger.Println("Starting servers on port ", srv.Addr)

		err := srv.ListenAndServe()
		if err != nil {
			severAPILogger.Printf("Error starting servers: %s\n", err)
			os.Exit(1)
		}
	}()
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
