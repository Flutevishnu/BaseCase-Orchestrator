package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UserIDKey contextKey = "userId"

func RequireUserId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("X-User-Id")
		if userId == "" {
			http.Error(w, "Unauthorized: Missing X-User-Id header", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
