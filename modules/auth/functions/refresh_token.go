package auth_functions

import (
	"context"
	"fmt"
	"net/http"
	auth_token "oauth2/core/helpers/auth"
	responses "oauth2/core/helpers/response"
	"oauth2/data/db"
	"oauth2/data/models"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	conn := db.SetupDB()
	defer conn.Close(context.Background())

	// Parse do corpo da requisição x-www-form-urlencoded
	err := r.ParseForm()
	if err != nil {
		responses.Err(w, http.StatusBadRequest, fmt.Errorf("invalid form data"))
		return
	}

	// Pegar o refresh_token do formulário
	refreshToken := r.FormValue("refresh_token")
	if refreshToken == "" {
		responses.Err(w, http.StatusBadRequest, fmt.Errorf("refresh_token is required"))
		return
	}

	// Buscar usuário pelo refresh token
	var user models.User
	query := "SELECT users.id, users.name, users.email FROM users JOIN tokens ON users.id = tokens.user_id WHERE tokens.refresh_token = $1"
	err = conn.QueryRow(context.Background(), query, refreshToken).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, fmt.Errorf("invalid refresh token"))
		return
	}

	// Gerar novos tokens
	newAccessToken, newRefreshToken, err := auth_token.CreateToken(user.ID, refreshToken)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Atualizar o refresh token apenas se ele tiver mudado
	if newRefreshToken != refreshToken {
		updateTokenQuery := "UPDATE tokens SET refresh_token = $1 WHERE user_id = $2"
		_, err = conn.Exec(context.Background(), updateTokenQuery, newRefreshToken, user.ID)
		if err != nil {
			responses.Err(w, http.StatusInternalServerError, err)
			return
		}
	}

	// Responder com os novos tokens
	responseData := map[string]interface{}{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}

	responses.JSON(w, http.StatusOK, responseData)
}
