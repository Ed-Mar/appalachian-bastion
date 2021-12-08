package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "server-api", log.LstdFlags)

	severRouter := mux.NewRouter()
	ServerHandler :=
		severRouter.HandleFunc("/resources", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "/resources")
		})

	srv := &http.Server{
		Handler:      severRouter,
		Addr:         ":9090",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// start the server
	go func() {
		logger.Println("Starting server on port 9090")

		err := srv.ListenAndServe()
		if err != nil {
			logger.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	srv.Shutdown(ctx)

	log.Fatal(srv.ListenAndServe())
}
