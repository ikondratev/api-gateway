package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ikondratev/api-gateway/internal/apperr"
	"github.com/ikondratev/api-gateway/internal/authpb"
	"github.com/ikondratev/api-gateway/internal/logger"
)

type Auther interface {
	Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.TokenPair, error)
	Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.TokenPair, error)
	Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.TokenPair, error)
	Logout(ctx context.Context, req *authpb.LogoutRequest) error
}

type authHandler struct {
	log  logger.Logger
	auth Auther
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

func (h authHandler) login(w http.ResponseWriter, r *http.Request) {
	var req authpb.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, h.log, apperr.Business(http.StatusBadRequest, "invalid_request", "invalid request body"))
		return
	}

	resp, err := h.auth.Login(r.Context(), &req)
	if err != nil {
		WriteError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h authHandler) refresh(w http.ResponseWriter, r *http.Request) {
	var req authpb.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, h.log, apperr.Business(http.StatusBadRequest, "invalid_request", "invalid request body"))
		return
	}

	resp, err := h.auth.Refresh(r.Context(), &req)
	if err != nil {
		WriteError(w, h.log, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h authHandler) logout(w http.ResponseWriter, r *http.Request) {
	var req authpb.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, h.log, apperr.Business(http.StatusBadRequest, "invalid_request", "invalid request body"))
		return
	}

	if err := h.auth.Logout(r.Context(), &req); err != nil {
		WriteError(w, h.log, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
