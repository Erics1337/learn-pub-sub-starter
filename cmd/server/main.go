package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const rabbitMQConnString = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril server...")
	fmt.Println("Waiting for RabbitMQ to start...")
	time.Sleep(7 * time.Second) // Add a delay to prevent race condition from container startup time

	conn, err := amqp.Dial(rabbitMQConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Successfully connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()
	fmt.Println("Channel created")

	// Publish a pause message
	gameState := routing.PlayingState{IsPaused: true}
	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, gameState)
	if err != nil {
		log.Fatalf("failed to publish pause message: %v", err)
	}
	fmt.Println("Pause message published successfully")

	fmt.Println("Waiting for messages... To exit press CTRL+C")

	// Block until a signal is received
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down...")
}
