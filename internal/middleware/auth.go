package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/user/golang-api-rest/internal/config"
	"github.com/user/golang-api-rest/internal/utils"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyRole   contextKey = "role"
)

func RequireAuth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "authorization header ausente", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "formato de authorization inválido, use: Bearer <token>", http.StatusUnauthorized)
				return
			}

			claims, err := utils.ValidateToken(parts[1], cfg)
			if err != nil {
				http.Error(w, "token inválido ou expirado", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextKeyRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(ContextKeyRole).(string)
			if !ok || userRole != role {
				http.Error(w, "acesso negado", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserIDFromContext(ctx context.Context) int {
	if id, ok := ctx.Value(ContextKeyUserID).(int); ok {
		return id
	}
	return 0
}

func GetRoleFromContext(ctx context.Context) string {
	if role, ok := ctx.Value(ContextKeyRole).(string); ok {
		return role
	}
	return ""
}