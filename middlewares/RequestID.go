package middlewares

import (
	"context"
	"net/http"
	"notes-server/constants"

	"github.com/google/uuid"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), constants.RequestIDCtxKey, requestID)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
