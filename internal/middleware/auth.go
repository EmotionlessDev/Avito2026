package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userCtxKey contextKey = "user"

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WithUser(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, userCtxKey, claims)
}

func UserFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(userCtxKey).(*Claims)
	if !ok {
		return nil, common.ErrUnauthorized
	}
	return claims, nil
}
