package main

import (
	"log"
	"net/http"
	"os"

	"southpark/go-api/internal/adapters/httpadapter"
	"southpark/go-api/internal/adapters/rabbitmq"
	"southpark/go-api/internal/app"
)

func main() {
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	port := getEnv("PORT", "8080")

	// Driven adapter (outgoing)
	publisher, err := rabbitmq.NewPublisher(rabbitURL)
	if err != nil {
		log.Fatalf("rabbitmq: %v", err)
	}
	defer publisher.Close()

	// Core
	service := app.NewMessageService(publisher)

	// Driving adapter (incoming)
	handler := httpadapter.NewHandler(service)

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler.Routes()))
}

// getEnv returns the environment variable, or fallback if it is empty.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}