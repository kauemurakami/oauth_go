package auth_token

import (
	"time"

	app_config "oauth2/core/config"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Criar access token e possivelmente um novo refresh token
func CreateToken(userID uuid.UUID, currentRefreshToken string) (accessToken string, refreshToken string, err error) {
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

	// Validar se o refresh token ainda está válido
	_, err = ValidateRefreshToken(currentRefreshToken)
	if err == nil {
		// Se o refresh token ainda for válido, retorná-lo sem criar um novo
		return accessToken, currentRefreshToken, nil
	}

	// Caso contrário, criar um novo refresh token
	refreshClaims := jwt.MapClaims{}
	refreshClaims["userID"] = userID
	refreshClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() // 7 dias de validade

	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshJwt.SignedString([]byte(app_config.REFRESH_SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
