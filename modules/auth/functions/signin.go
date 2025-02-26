package auth_functions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	auth_token "oauth2/core/helpers/auth"
	responses "oauth2/core/helpers/response"
	sec "oauth2/core/security"
	"oauth2/data/db"
	"oauth2/data/models"
)

func Signin(w http.ResponseWriter, r *http.Request) {
	conn := db.SetupDB()
	defer conn.Close(context.Background())
	// Init transaction
	tx, err := conn.Begin(context.Background())
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback(context.Background())

	// Ler o corpo da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
		return
	}

	var userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Deserializar os dados
	if err = json.Unmarshal(body, &userInput); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}

	// Buscar o usuário no banco de dados
	var user models.User
	query := "SELECT id, name, email, pass FROM users WHERE email = $1"
	err = tx.QueryRow(context.Background(), query, userInput.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Pass)
	if err != nil {
		if err.Error() == "no rows in result set" {
			responses.Err(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
			return
		}
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Verificar se a senha é válida utilizando o pacote sec
	err = sec.ComparePassHash(userInput.Password, user.Pass)
	if err != nil {
		responses.Err(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	// Gerar novos tokens
	accessToken, refreshToken, err := auth_token.CreateToken(user.ID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Atualizar os tokens no banco de dados
	updateTokenQuery := "UPDATE tokens SET access_token = $1, refresh_token = $2 WHERE user_id = $3"
	_, err = tx.Exec(context.Background(), updateTokenQuery, accessToken, refreshToken, user.ID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Commit the transaction
	if err = tx.Commit(context.Background()); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Resposta com os tokens
	response := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	responses.JSON(w, http.StatusOK, response)
}
