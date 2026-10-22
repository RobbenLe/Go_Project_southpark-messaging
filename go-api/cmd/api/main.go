package main
import (
	"fmt" //For printing
	"southpark/go-api/internal/domain"
)

func main() {
	good := domain.Message{Author: "Cartman", Body: "Respect my authoritah!"}
	bad := domain.Message{Author: "Kenny", Body: "  "}

	fmt.Println("good:", good.Validate())
	fmt.Println("bad: ", bad.Validate())
}