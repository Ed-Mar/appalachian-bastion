package main

import (
	consumerGroups "backend/channel-service/event-streaming/reader/consumer-group"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	reader := consumerGroups.ChannelServiceConsumerGroup
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

	// trap sigterm or interrupt and gracefully shutdown the servers
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	err := reader.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
}
