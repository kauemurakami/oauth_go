package auth_functions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	auth_token "oauth2/core/helpers/auth"
	responses "oauth2/core/helpers/response"
	"oauth2/data/db"
	"oauth2/data/models"
)

func Signup(w http.ResponseWriter, r *http.Request) {

	conn := db.SetupDB()
	defer conn.Close(context.Background())
	// Init transaction
	tx, err := conn.Begin(context.Background())
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback(context.Background())
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Err(w, http.StatusUnprocessableEntity, err)
	}
	var user models.User
	if err = json.Unmarshal(body, &user); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
	}
	if err = user.Prepare("signup"); err != nil {
		responses.Err(w, http.StatusBadRequest, err)
		return
	}
	query := "INSERT INTO users (name, email, pass) VALUES ($1, $2, $3) RETURNING *"
	var insertedUser models.User
	err = tx.QueryRow(context.Background(),
		query,
		user.Name,
		user.Email,
		user.Pass,
	).Scan(
		&insertedUser.ID,
		&insertedUser.Name,
		&insertedUser.Email,
		&insertedUser.Pass,
		&insertedUser.CreatedAt,
	)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
	}

	token, refreshToken, err := auth_token.CreateToken(user.ID)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
	}
	fmt.Println(token, refreshToken, user.ID)
	insertTokenQuery := "INSERT INTO tokens (access_token,refresh_token, user_id) VALUES ($1, $2, $3)"
	_, err = tx.Exec(context.Background(), insertTokenQuery,
		token,
		refreshToken,
		insertedUser.ID,
	)

	if err != nil {
		fmt.Println(err)
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Commit the transaction
	if err = tx.Commit(context.Background()); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"user": map[string]string{
			"id":    insertedUser.ID.String(),
			"name":  insertedUser.Name,
			"email": insertedUser.Email,
		},
		"access_token":  token,
		"refresh_token": refreshToken,
	}

	responses.JSON(w, http.StatusOK, response)
}
