package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ikondratev/api-gateway/internal/apperr"
	"github.com/ikondratev/api-gateway/internal/eventpb"
	"github.com/ikondratev/api-gateway/internal/logger"
)

type EventCreator interface {
	CreateEvent(
		ctx context.Context,
		authorization string,
		idempotencyKey string,
		req *eventpb.CreateEventRequest,
	) (*eventpb.CreateEventResponse, error)
}

type eventHandler struct {
	log    logger.Logger
	events EventCreator
}

func (h eventHandler) create(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		WriteError(w, h.log, apperr.Business(http.StatusUnauthorized, "unauthorized", "unauthorized"))
		return
	}

	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		WriteError(w, h.log, apperr.Business(http.StatusBadRequest, "invalid_event", "idempotency key is required"))
		return
	}

	var req eventpb.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, h.log, apperr.Business(http.StatusBadRequest, "invalid_event", "invalid request body"))
		return
	}

	resp, err := h.events.CreateEvent(r.Context(), authorization, idempotencyKey, &req)
	if err != nil {
		WriteError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
