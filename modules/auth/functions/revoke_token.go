package auth_functions

import (
	"context"
	"fmt"
	"net/http"
	responses "oauth2/core/helpers/response"
	"oauth2/data/db"
)

func RevokeToken(w http.ResponseWriter, r *http.Request) {
	// Garantir que a requisição seja do tipo x-www-form-urlencoded
	if err := r.ParseForm(); err != nil {
		responses.Err(w, http.StatusBadRequest, fmt.Errorf("invalid form data"))
		return
	}

	// Capturar o refresh_token do formulário
	refreshToken := r.FormValue("refresh_token")
	if refreshToken == "" {
		responses.Err(w, http.StatusBadRequest, fmt.Errorf("refresh_token is required"))
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

	// Excluir a linha correspondente ao refresh_token
	deleteQuery := "DELETE FROM tokens WHERE refresh_token = $1"
	res, err := tx.Exec(context.Background(), deleteQuery, refreshToken)
	if err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Verificar se algum registro foi deletado
	rowsAffected := res.RowsAffected()

	if rowsAffected == 0 {
		responses.Err(w, http.StatusNotFound, fmt.Errorf("refresh token not found"))
		return
	}

	// Commit da transação
	if err = tx.Commit(context.Background()); err != nil {
		responses.Err(w, http.StatusInternalServerError, err)
		return
	}

	// Responder com sucesso
	responses.JSON(w, http.StatusOK, map[string]string{"message": "token revoked successfully"})
}
