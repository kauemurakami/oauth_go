package auth_token

import (
	app_config "oauth2/core/config"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Criar access token e refresh token
func CreateToken(userID uuid.UUID) (accessToken string, refreshToken string, err error) {
	// Claims do access token (válido por 1 minuto)
	accessClaims := jwt.MapClaims{}
	accessClaims["authorized"] = true
	accessClaims["userID"] = userID
	accessClaims["exp"] = time.Now().Add(time.Minute * 1).Unix() // 1 minuto de validade

	// Criar access token
	accessJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessJwt.SignedString([]byte(app_config.SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	// Claims do refresh token (válido por 7 dias)
	refreshClaims := jwt.MapClaims{}
	refreshClaims["userID"] = userID
	refreshClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() // 7 dias de validade

	// Criar refresh token
	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshJwt.SignedString([]byte(app_config.REFRESH_SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
