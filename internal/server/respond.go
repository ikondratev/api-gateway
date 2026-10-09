package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ikondratev/api-gateway/internal/apperr"
	"github.com/ikondratev/api-gateway/internal/logger"
)

func WriteError(w http.ResponseWriter, log logger.Logger, err error) {
	var appErr *apperr.Error
	if errors.As(err, &appErr) && appErr.Err == nil {
		writeError(w, appErr.Status, appErr.Code, appErr.Message)
		return
	}

	log.Error("request failed", "error", err)
	writeError(w, http.StatusInternalServerError, "internal", "internal error")
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
