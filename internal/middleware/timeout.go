package middleware

import (
	"context"
	"net/http"
	"time"
)

func Timeout(d time.Duration) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
