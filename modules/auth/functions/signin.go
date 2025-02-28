package auth_functions

import (
	"context"
	"fmt"
	"net/http"

	auth_token "oauth2/core/helpers/auth"
	responses "oauth2/core/helpers/response"
	sec "oauth2/core/security"
	"oauth2/data/db"
	"oauth2/data/models"
)

func Signin(w http.ResponseWriter, r *http.Request) {
	// Garantir que a requisição seja do tipo x-www-form-urlencoded
	if err := r.ParseForm(); err != nil {
		responses.Err(w, http.StatusBadRequest, fmt.Errorf("invalid form data"))
		return
	}

	// Capturar os valores do formulário
	email := r.FormValue("email")
	password := r.FormValue("pass")

	if email == "" || password == "" {
		responses.Err(w, http.StatusBadRequest, fmt.Errorf("email and password are required"))
		return
	}

	// Conectar ao banco
	conn := db.SetupDB()
	defer conn.Close(context.Background())

	// Iniciar transação
	tx, err := conn.Begin(context.Background())
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback(context.Background())

	// Buscar o usuário no banco de dados
	var user models.User
	query := "SELECT id, name, email, pass FROM users WHERE email = $1"
	err = tx.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Pass)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	// Verificar se a senha é válida
	err = sec.ComparePassHash(password, user.Pass)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	// Buscar o refresh token no banco de dados
	var existingRefreshToken string
	getTokenQuery := "SELECT refresh_token FROM tokens WHERE user_id = $1"
	err = tx.QueryRow(context.Background(), getTokenQuery, user.ID).Scan(&existingRefreshToken)
	if err != nil {
		existingRefreshToken = "" // Se não houver token, criamos um novo
	}

	// Gerar novos tokens
	accessToken, refreshToken, err := auth_token.CreateToken(user.ID, existingRefreshToken)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Atualizar ou inserir tokens no banco
	if existingRefreshToken == "" {
		insertTokenQuery := "INSERT INTO tokens (access_token, refresh_token, user_id) VALUES ($1, $2, $3)"
		_, err = tx.Exec(context.Background(), insertTokenQuery, accessToken, refreshToken, user.ID)
	} else if existingRefreshToken != refreshToken {
		updateTokenQuery := "UPDATE tokens SET access_token = $1, refresh_token = $2 WHERE user_id = $3"
		_, err = tx.Exec(context.Background(), updateTokenQuery, accessToken, refreshToken, user.ID)
	}

	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Commit da transação
	if err = tx.Commit(context.Background()); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Responder com os tokens
	response := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	responses.JSON(w, http.StatusOK, response)
}
