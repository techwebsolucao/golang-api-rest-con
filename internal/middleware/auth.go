package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/user/golang-api-rest/internal/config"
	"github.com/user/golang-api-rest/internal/response"
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
				response.Error(w, http.StatusUnauthorized, "authorization header ausente")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.Error(w, http.StatusUnauthorized, "formato de authorization inválido, use: Bearer <token>")
				return
			}

			claims, err := utils.ValidateToken(parts[1], cfg)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "token inválido ou expirado")
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
				response.Error(w, http.StatusForbidden, "acesso negado")
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
