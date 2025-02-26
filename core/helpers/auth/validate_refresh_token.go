package auth_token

import (
	"errors"
	"time"

	app_config "oauth2/core/config"

	jwt "github.com/dgrijalva/jwt-go"
)

// ValidateRefreshToken verifica se um refresh token é válido e não expirou
func ValidateRefreshToken(tokenString string) (userID string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(app_config.REFRESH_SECRET_KEY), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// Verifica se o token expirou
	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return "", errors.New("refresh token expired")
	}

	// Retorna o userID do token
	userID, ok = claims["userID"].(string)
	if !ok {
		return "", errors.New("invalid userID in token")
	}

	return userID, nil
}
