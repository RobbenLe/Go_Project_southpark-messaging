package main
import (
	"fmt" //For printing
	"context"
	"southpark/go-api/internal/domain"
	ports "southpark/go-api/internal/ports"
)

// printPublisher is a fake publisher that just prints.
type  printPublisher struct{}

func (p printPublisher) Publish(ctx context.Context, msg domain.Message) error {
	fmt.Println("Published:", msg.Author, "says", msg.Body)
	return nil
}
func main() {
	var pub ports.MessagePublisher = printPublisher{}

	msg := domain.Message{Author: "Cartman", Body: "Screw you guys, I'm going home!"}
	err := pub.Publish(context.Background(), msg)
	fmt.Println("error:", err)
}