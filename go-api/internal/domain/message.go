//A South Park message has:
  //- an author  
  //- a body     

//A message is NOT OK if:
  //- the author is empty, or
  //- the body is empty

package domain
import (
	"errors"
	"strings"
)

var (
	ErrMissingAuthor = errors.New("author is required")
	ErrMissingBody = errors.New("body is required")
)

type Message struct { // struct is a class
	Author string `json:"author"`
	Body   string `json:"body"`
}

func (m Message) Validate() error {
	if strings.TrimSpace(m.Author) == "" {
		return ErrMissingAuthor
	}
	if strings.TrimSpace(m.Body) == "" {
		return ErrMissingBody
	}
	return nil
}
