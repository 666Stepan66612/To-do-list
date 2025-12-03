package middleware

import (
    "apiservice/auth"
    "context"
    "net/http"
    "strings"
)

type contextKey string

const UserContextKey = contextKey("user")

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid authorization header format. Expected: Bearer <token>", http.StatusUnauthorized)
            return
        }

        tokenString := parts[1]

        claims, err := auth.ValidateToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), UserContextKey, claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func GetUserFromContext(r *http.Request) *auth.Claims {
    claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
    if !ok {
        return nil
    }
    return claims
}