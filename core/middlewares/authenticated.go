package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	auth_token "oauth2/core/helpers/auth"
	responses "oauth2/core/helpers/response"
	"oauth2/data/db"
)

// AuthMiddleware verifica se o accessToken é válido antes de permitir o acesso
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar se há um token no cabeçalho
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			responses.Err(w, http.StatusUnauthorized, errors.New("missing authorization header"))
			return
		}

		// Separar "Bearer" do token real
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			responses.Err(w, http.StatusUnauthorized, errors.New("invalid token format"))
			return
		}
		accessToken := parts[1] // O token real

		// Conectar ao banco de dados
		conn := db.SetupDB()
		defer conn.Close(context.Background())

		// Verificar se o token existe no banco
		var count int
		query := "SELECT COUNT(*) FROM tokens WHERE access_token = $1"
		err := conn.QueryRow(context.Background(), query, accessToken).Scan(&count)
		if err != nil || count == 0 {
			responses.Err(w, http.StatusUnauthorized, errors.New("token has been revoked"))
			return
		}

		// Validar o token
		if err := auth_token.ValidateToken(r); err != nil {
			responses.Err(w, http.StatusUnauthorized, errors.New("invalid or expired token"))
			return
		}

		// Se chegou até aqui, o token é válido e existe no banco -> continuar execução
		next.ServeHTTP(w, r)
	})
}
