package app

import (
	"context"
	"errors"
	"testing"

	"southpark/go-api/internal/domain"
)

// fakePublisher records messages instead of sending them.
type fakePublisher struct {
	published []domain.Message
	err       error
}

func (f *fakePublisher) Publish(ctx context.Context, msg domain.Message) error {
	if f.err != nil {
		return f.err
	}
	f.published = append(f.published, msg)
	return nil
}

func TestSendMessage_Valid(t *testing.T) {
	fake := &fakePublisher{}
	svc := NewMessageService(fake)

	msg := domain.Message{Author: "Cartman", Body: "Respect my authoritah!"}
	err := svc.SendMessage(context.Background(), msg)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fake.published) != 1 {
		t.Fatalf("want 1 published message, got %d", len(fake.published))
	}
}

func TestSendMessage_Invalid_NotPublished(t *testing.T) {
	fake := &fakePublisher{}
	svc := NewMessageService(fake)

	err := svc.SendMessage(context.Background(), domain.Message{Author: "Kenny", Body: ""})

	if !errors.Is(err, domain.ErrMissingBody) {
		t.Fatalf("want ErrMissingBody, got %v", err)
	}
	if len(fake.published) != 0 {
		t.Fatalf("invalid message should not be published")
	}
}

func TestSendMessage_PublisherFails(t *testing.T) {
	boom := errors.New("rabbitmq is down")
	fake := &fakePublisher{err: boom}
	svc := NewMessageService(fake)

	msg := domain.Message{Author: "Butters", Body: "Oh hamburgers!"}
	err := svc.SendMessage(context.Background(), msg)

	if !errors.Is(err, boom) {
		t.Fatalf("want wrapped %v, got %v", boom, err)
	}
}