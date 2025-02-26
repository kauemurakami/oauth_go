package users_funtcions

import (
	"context"
	"net/http"
	responses "oauth2/core/helpers/response"
	"oauth2/data/db"
	"oauth2/data/models"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	conn := db.SetupDB()
	defer conn.Close(context.Background())

	query := `
		SELECT id,name,email
		FROM users
	`
	//exec query
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			responses.Err(w, http.StatusInternalServerError, err)
			return
		}

		users = append(users, user)
	}

	// Verificar por erros no final do loop
	if err := rows.Err(); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}
