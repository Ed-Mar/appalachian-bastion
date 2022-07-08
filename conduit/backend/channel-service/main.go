package main

import (
	"backend/channel-service/handlers"
	"backend/internal"
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

	ChannelServiceLogger := log.New(os.Stdout, "channel-service | ", log.LstdFlags)

	validation := internal.NewValidation()

	channelHandler := handlers.NewChannels(ChannelServiceLogger, validation)

	serveMux := mux.NewRouter()

	//Note: I am not positive if I am going to support URL id passing for the server
	//for the extent of time but for now I will do both cause the way I have it modeled
	//one does not really need servers id if it passed in the channel obj itself
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/servers/{serverID}/channels", channelHandler.ListAllChannelsWithMatchingServerID)
	getRouter.HandleFunc("/servers/channels/{channelID}", channelHandler.ListSingle)
	getRouter.HandleFunc("/servers/channels", channelHandler.ListEveryChannel)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/servers/channel/", channelHandler.UpdateWithoutParms)
	putRouter.Use(channelHandler.MiddlewareValidateChannel)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/servers/channel", channelHandler.CreateChannelWithoutParms)
	postRouter.Use(channelHandler.MiddlewareValidateChannel)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/servers/channels/{channelID}", channelHandler.Delete)
	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	srv := &http.Server{
		Addr:         ":9393",               // configure the bind address
		Handler:      corsHandler(serveMux), // set the default handlers
		ErrorLog:     ChannelServiceLogger,  // set the severChannelServiceLogger for the servers
		ReadTimeout:  5 * time.Second,       // max time to read request from the client
		WriteTimeout: 10 * time.Second,      // max time to write response to the client
		IdleTimeout:  120 * time.Second,     // max time for connections using TCP Keep-Alive
	}

	// start the servers
	go func() {
		ChannelServiceLogger.Println("Starting servers on port ", srv.Addr)

		err := srv.ListenAndServe()
		if err != nil {
			ChannelServiceLogger.Printf("Error starting servers: %s\n", err)
			os.Exit(1)
		}
	}()

	//TODO Clean this up below

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
