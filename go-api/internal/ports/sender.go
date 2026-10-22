package port

import (
	"context"

	"southpark/go-api/internal/domain"
)

// MessageSender is what the outside world can ask the app to do.
type MessageSender interface {
	SendMessage(ctx context.Context, msg domain.Message) error
}