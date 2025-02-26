package auth_token

import (
	"errors"
	"fmt"
	"net/http"
	app_config "oauth2/core/config"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// validar se o token passado é valido
func ValidateToken(r *http.Request) error {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, returnKeyVerify)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}
	fmt.Println(token)
	return errors.New("token inválido ou expirado, faça login novamente")
}

func returnKeyVerify(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("método de assinatura inválido, faça login novamente")
	}
	return app_config.SECRET_KEY, nil
}

func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}
