//it takes the port, validates the message, then publishes. This is where domain and port finally work together
 package app

import (
	"context"
	"fmt"

	"southpark/go-api/internal/domain"
	ports "southpark/go-api/internal/ports"
)

type MessageService struct {
	publisher ports.MessagePublisher
}

func NewMessageService(p ports.MessagePublisher) *MessageService {
	return &MessageService{publisher: p} //&:	"Give me the address of this struct" (a pointer)
                                         //*MessageService	The return type: "a pointer to a MessageService"
}

func (s *MessageService) SendMessage(ctx context.Context, msg domain.Message) error {
	if err := msg.Validate(); err != nil {
		return err
	}

	if err := s.publisher.Publish(ctx, msg); err != nil {
		return fmt.Errorf("publish message: %w", err)
	}

	return nil
}