package main

import (
	"context"
	"fmt" //For printing
	"log"
	"net/http"
	"southpark/go-api/internal/adapters/httpadapter"
	"southpark/go-api/internal/app"
	"southpark/go-api/internal/domain"
)

// printPublisher is a fake publisher that just prints.
type  printPublisher struct{}

func (p printPublisher) Publish(ctx context.Context, msg domain.Message) error {
	fmt.Println("Published:", msg.Author, "says", msg.Body)
	return nil
}
func main() {
	service := app.NewMessageService(printPublisher{})
	handler := httpadapter.NewHandler(service)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler.Routes()))
}