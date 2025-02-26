package middlewares

import (
	"net/http"

	auth_token "oauth2/core/helpers/auth"
	responses "oauth2/core/helpers/response"
)

// AuthMiddleware verifica se o accessToken é válido antes de permitir o acesso
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Apenas verifica se o token é válido
		if err := auth_token.ValidateToken(r); err != nil {
			responses.Err(w, http.StatusUnauthorized, err)
			return
		}

		// Se o token for válido, continua para o próximo handler
		next.ServeHTTP(w, r)
	})
}
