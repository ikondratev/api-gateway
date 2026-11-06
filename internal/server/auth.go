package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ikondratev/api-gateway/internal/apperr"
	"github.com/ikondratev/api-gateway/internal/authpb"
	"github.com/ikondratev/api-gateway/internal/logger"
)

type Registrar interface {
	Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.TokenPair, error)
}

type authHandler struct {
	log  logger.Logger
	auth Registrar
}

func (h authHandler) register(w http.ResponseWriter, r *http.Request) {
	var req authpb.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, h.log, apperr.Business(http.StatusBadRequest, "invalid_request", "invalid request body"))
		return
	}

	resp, err := h.auth.Register(r.Context(), &req)
	if err != nil {
		WriteError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
