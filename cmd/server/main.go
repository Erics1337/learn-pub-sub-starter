package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

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
	fmt.Println("Waiting for messages... To exit press CTRL+C")

	// Block until a signal is received
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down...")
}
