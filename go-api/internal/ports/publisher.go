package port

import (
	"context"
    "southpark/go-api/internal/domain"
)

// MessagePublisher sends a domain Message to a messaging system.
type MessagePublisher interface {
	Publish(ctx context.Context, msg domain.Message) error //with context:   Go gets a "cancelled" signal and stops once User closed browser so nothing will be sent.

}