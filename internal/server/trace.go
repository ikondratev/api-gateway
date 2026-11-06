package server

import (
	"context"
	"net/http"
	"regexp"

	"github.com/google/uuid"
)

type traceIDKey struct{}

const traceIDHeader = "X-Trace-Id"

var traceIDPattern = regexp.MustCompile(`^[A-Za-z0-9-]{1,128}$`)

func Trace() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(traceIDHeader)
			if !traceIDPattern.MatchString(id) {
				id = uuid.NewString()
			}
			w.Header().Set(traceIDHeader, id)
			ctx := context.WithValue(r.Context(), traceIDKey{}, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func TraceIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(traceIDKey{}).(string)
	return id
}
