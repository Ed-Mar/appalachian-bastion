package main

import (
	"backend/event-streaming/kafka/messages/model"
	consumerGroups "backend/server-service/event-streaming/reader/consumer-group"
	event_handler "backend/server-service/event-streaming/reader/event-handler"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
)

var ErrReadingEventMessage = fmt.Errorf("[ERROR] [EVENT] [READER] [SEVER]: Unexpected Error Occuied while attemping to Read Incoming Messages. Restart is Needed. |  ")
var ErrJSONUnmarshalling = fmt.Errorf("[ERROR] [EVENT] [READER] [SEVER] [JSON]: Error occuied while unmarshaling btyes from event message |  ")

func main() {
	serverEventMessagesLogger := log.New(os.Stdout, "server-event-messages | ", log.LstdFlags)
	serverEventHandler := event_handler.NewServerEvent(serverEventMessagesLogger)
	consumerGroup := consumerGroups.ServerServiceConsumerGroup
	go func() {
		serverEventHandler.Logger.Println("Starting Kafka Reader for topic: ", consumerGroup.Config().Topic)
		for {
			// Reading Received Message
			readMessage, err := consumerGroup.ReadMessage(context.Background())
			if err != nil {
				serverEventHandler.Logger.Println(ErrReadingEventMessage)
				break
			}
			//serverEventMessagesLogger.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", readMessage.Topic, readMessage.Partition, readMessage.Offset, string(readMessage.Key), string(readMessage.Value))

			//Converting Value field of Received message from Byte Blob into Common EventMessage
			var eventStructs model.EventMessage
			err = json.Unmarshal(readMessage.Value, &eventStructs)
			if err != nil {
				serverEventHandler.Logger.Println(ErrJSONUnmarshalling, err)
			}
			serverEventHandler.Logger.Println("%+v\n", &eventStructs)
			//Processing read message into what should be done.
			serverEventHandler.EventMux(eventStructs)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the servers
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	err := consumerGroup.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
}
