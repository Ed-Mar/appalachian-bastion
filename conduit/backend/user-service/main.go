package main

import (
	"backend/internal"
	auth "backend/internal/authentication/handlers"
	userHandlers "backend/user-service/handlers"
	"context"
	"fmt"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const serviceName = "user-service"

func main() {
	APILogger := log.New(os.Stdout, fmt.Sprintf("%s | ", serviceName), log.LstdFlags)
	validation := internal.NewValidation()

	genericHandler := auth.NewHandler(serviceName, APILogger)
	handler := userHandlers.NewUserHandler(genericHandler, validation)

	router := mux.NewRouter()
	router.Use(handler.StandardHandler.AuthenticationMiddlewareViaLocalVerification)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.Use(handler.StandardHandler.AuthenticationMiddlewareViaTokenIntrospective)
	postRouter.HandleFunc("/", handler.CreateNewUserViaHeader)

	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.Use(handler.StandardHandler.AuthenticationMiddlewareViaTokenIntrospective)
	getRouter.HandleFunc("/", handler.GetUserProfileViaExternalID)

	putRouter := router.Methods(http.MethodPut).Subrouter()
	putRouter.Use(handler.StandardHandler.AuthenticationMiddlewareViaTokenIntrospective)

	deleteRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.Use(handler.StandardHandler.AuthenticationMiddlewareViaTokenIntrospective)

	credentials := gohandlers.AllowCredentials()
	origins := gohandlers.AllowedOrigins([]string{"http://localhost:8080"})
	headers := gohandlers.AllowedHeaders([]string{"Authentication"})
	methods := gohandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	srv := &http.Server{
		Addr:         ":9661", // configure the bind address
		Handler:      gohandlers.CORS(credentials, methods, origins, headers)(router),
		ErrorLog:     APILogger,         // set the severAPILogger for the servers
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the servers
	go func() {
		APILogger.Printf("Starting servers on port %+v\n ", srv.Addr)

		err := srv.ListenAndServe()
		if err != nil {
			APILogger.Printf("Error starting servers: %s\n", err)
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
