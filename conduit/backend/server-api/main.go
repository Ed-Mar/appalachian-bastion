package main

import (
	"backend/internal"
	"backend/internal/database"
	"backend/server-api/handlers/channels"
	"backend/server-api/handlers/servers"
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
	severAPILogger := log.New(os.Stdout, "server-api | ", log.LstdFlags)
	channelAPILogger := log.New(os.Stdout, "channel-api | ", log.LstdFlags)
	//databaseAPILogger := log.New(os.Stdout, "postgres-db | ", log.LstdFlags)

	validation := internal.NewValidation()

	serverHandler := servers.NewServers(severAPILogger, validation)
	channelHandler := channels.NewChannels(channelAPILogger, validation)

	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/servers", serverHandler.ListAll)
	getRouter.HandleFunc("/servers/{id:[0-9]+}", serverHandler.ListSingle)
	getRouter.HandleFunc("/servers/{id:[0-9]+}/channels", channelHandler.ListAll)
	getRouter.HandleFunc("/servers/{id:[0-9]+}/channels/{id:[0-9]+}", channelHandler.ListSingle)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/servers", serverHandler.Update)
	putRouter.Use(serverHandler.MiddlewareValidateServer)

	putChannelRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putChannelRouter.HandleFunc("/servers/{id:[0-9]+}/channels", channelHandler.Update)
	putChannelRouter.Use(channelHandler.MiddlewareValidateChannel)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/servers", serverHandler.Create)
	postRouter.Use(serverHandler.MiddlewareValidateServer)

	postChannelRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postChannelRouter.HandleFunc("/servers/{id:[0-9]+}/channels", channelHandler.Create)
	postChannelRouter.Use(channelHandler.MiddlewareValidateChannel)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/servers/{id:[0-9]+}", serverHandler.Delete)
	deleteRouter.HandleFunc("/servers/{id:[0-9]+}/channels/{id:[0-9]+}", channelHandler.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	getRouter.Handle("/swagger.json", http.FileServer(http.Dir("./")))

	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	srv := &http.Server{
		Addr:         ":9090",               // configure the bind address
		Handler:      corsHandler(serveMux), // set the default handlers
		ErrorLog:     severAPILogger,        // set the severAPILogger for the server
		ReadTimeout:  5 * time.Second,       // max time to read request from the client
		WriteTimeout: 10 * time.Second,      // max time to write response to the client
		IdleTimeout:  120 * time.Second,     // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		severAPILogger.Println("Starting server on port 9090")

		err := srv.ListenAndServe()
		if err != nil {
			severAPILogger.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()
	//Make sure the db tables and model of the severs match up
	database.AutoMigrateDB()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Fatal(srv.ListenAndServe())
}
