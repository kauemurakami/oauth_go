package auth_functions

import (
	"context"
	"errors"
	"net/http"
	auth_token "oauth2/core/helpers/auth"
	responses "oauth2/core/helpers/response"
	"oauth2/data/db"

	"github.com/google/uuid"
)

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Parse o formulário com o refresh token e user_id
	err := r.ParseForm()
	if err != nil {
		responses.Err(w, http.StatusBadRequest, errors.New("unable to parse form"))
		return
	}

	// Pegue o refresh_token do campo "refresh_token"
	refreshToken := r.FormValue("refresh_token")
	if refreshToken == "" {
		responses.Err(w, http.StatusBadRequest, errors.New("refresh token is required"))
		return
	}

	// Pegue o user_id do campo "user_id"
	userID := r.FormValue("user_id")
	if userID == "" {
		responses.Err(w, http.StatusBadRequest, errors.New("user_id is required"))
		return
	}

	// Verifique a validade do refresh token
	validUserID, err := auth_token.ValidateRefreshToken(refreshToken)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, errors.New("invalid or expired refresh token"))
		return
	}

	// Certifique-se de que o user_id passado corresponde ao do token
	if validUserID != userID {
		responses.Err(w, http.StatusUnauthorized, errors.New("refresh token does not match user_id"))
		return
	}

	// Obtenha a conexão com o banco de dados e inicie uma transação
	conn := db.SetupDB()
	defer conn.Close(context.Background())

	tx, err := conn.Begin(context.Background())
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback(context.Background())

	uID, err := uuid.Parse(userID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Gere apenas um novo access_token, reutilizando o refresh_token existente
	accessToken, _, err := auth_token.CreateToken(uID, refreshToken)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Atualizar apenas o access_token no banco de dados
	updateTokenQuery := "UPDATE tokens SET access_token = $1 WHERE user_id = $2"
	_, err = tx.Exec(context.Background(), updateTokenQuery, accessToken, userID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Commit a transação
	if err = tx.Commit(context.Background()); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Retorne os tokens gerados
	response := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken, // Retorna o mesmo refresh_token se ainda for válido
	}
	responses.JSON(w, http.StatusOK, response)
}
