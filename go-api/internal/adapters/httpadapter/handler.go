package httpadapter

import (
	"encoding/json"
	"errors"
	"net/http"

	"southpark/go-api/internal/domain"
	ports "southpark/go-api/internal/ports"
)

// Handler is the HTTP (driving) adapter for messages.
type Handler struct {
	sender ports.MessageSender
}

// NewHandler creates a Handler that forwards to the given sender.
func NewHandler(s ports.MessageSender) *Handler {
	return &Handler{sender: s}
}

// Routes registers the HTTP endpoints.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /messages", h.postMessage)
	return mux
}

func (h *Handler) postMessage(w http.ResponseWriter, r *http.Request) {
	var msg domain.Message

	// Step 1: JSON text -> domain.Message
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&msg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	// Step 2: ask the app to send it, then pick a status code
	err := h.sender.SendMessage(r.Context(), msg)
	switch {
	case errors.Is(err, domain.ErrMissingAuthor), errors.Is(err, domain.ErrMissingBody):
		writeError(w, http.StatusBadRequest, err.Error())
		return
	case err != nil:
		writeError(w, http.StatusInternalServerError, "could not send message")
		return
	}

	// Step 3: success - queued for async processing
	writeJSON(w, http.StatusAccepted, map[string]string{"status": "queued"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}