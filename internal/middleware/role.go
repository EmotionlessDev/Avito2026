package middleware

import "net/http"

func RoleBased(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := UserFromContext(r.Context())
			if err != nil {
				http.Error(w, "unauthorized: user not found in context", http.StatusUnauthorized)
				return
			}

			for _, role := range allowedRoles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "forbidden: insufficient role", http.StatusForbidden)
		})
	}
}
