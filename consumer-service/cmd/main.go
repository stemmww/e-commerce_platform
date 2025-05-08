package main

import (
	"consumer-service/internal/nats"
	"log"
)

func main() {
	// Connect to NATS
	natsClient, err := nats.NewSubscriber("nats://nats:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsClient.Close()

	log.Println("✅ Consumer Service connected to NATS")

	// Subscribe to "order.created" topic
	if err := natsClient.Subscribe("order.created"); err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	log.Println("🚀 Listening for events...")

	// Block forever
	select {}
}
