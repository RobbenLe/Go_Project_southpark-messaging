package domain

import (
	"errors"
	"testing"
)

func TestMessageValidate(t *testing.T) {
	tests := []struct {
		name  string
		msg    Message
		wantErr error
	}{
		{"valid message", Message{Author: "Cartman", Body: "Respect my authoritah!"}, nil},
		{"missing author", Message{Author: "", Body: "Hi"}, ErrMissingAuthor},
		{"author only spaces", Message{Author: "   ", Body: "Hi"}, ErrMissingAuthor},
		{"missing body", Message{Author: "Kenny", Body: ""}, ErrMissingBody},
		{"body only spaces", Message{Author: "Kenny", Body: "   "}, ErrMissingBody},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.Validate()
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("got %v, want %v", err, tc.wantErr)
			}
		})
	}
}