package main

import (
	"backend/channel-api/handlers/channels"
	"backend/internal"
	"backend/internal/database"
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

	APILogger := log.New(os.Stdout, "channel-api | ", log.LstdFlags)

	validation := internal.NewValidation()

	channelHandler := channels.NewChannels(APILogger, validation)

	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/servers/{id:[0-9]+}/channels", channelHandler.ListAll)
	getRouter.HandleFunc("/servers/{id:[0-9]+}/channels/{id:[0-9]+}", channelHandler.ListSingle)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/servers/{id:[0-9]+}/channels", channelHandler.Update)
	putRouter.Use(channelHandler.MiddlewareValidateChannel)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/servers/{id:[0-9]+}/channels", channelHandler.Create)
	postRouter.Use(channelHandler.MiddlewareValidateChannel)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/servers/{id:[0-9]+}/channels/{id:[0-9]+}", channelHandler.Delete)

	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	srv := &http.Server{
		Addr:         ":9090",               // configure the bind address
		Handler:      corsHandler(serveMux), // set the default handlers
		ErrorLog:     severAPILogger,        // set the severAPILogger for the servers
		ReadTimeout:  5 * time.Second,       // max time to read request from the client
		WriteTimeout: 10 * time.Second,      // max time to write response to the client
		IdleTimeout:  120 * time.Second,     // max time for connections using TCP Keep-Alive
	}

	// start the servers
	go func() {
		severAPILogger.Println("Starting servers on port 9090")

		err := srv.ListenAndServe()
		if err != nil {
			severAPILogger.Printf("Error starting servers: %s\n", err)
			os.Exit(1)
		}
	}()
	//Make sure the db tables and models of the severs match up
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
