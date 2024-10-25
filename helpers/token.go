package helpers

import (
	"fmt"
	"time"

	"github.com/sipkyjayaputra/ticketing-system/configuration"
	"github.com/sipkyjayaputra/ticketing-system/model/entity"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	AuthSigningKey    = []byte(configuration.CONFIG["JWT_SECRET"])
	AuthSigningMethod = jwt.SigningMethodHS256
	AuthIssuer        = "SHARING VISION JAKARTA"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWTToken(user entity.User, expired time.Duration) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(expired).Unix(),
			Issuer:    AuthIssuer,
		},
	}

	token := jwt.NewWithClaims(AuthSigningMethod, claims)
	tokenString, err := token.SignedString(AuthSigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeJWTToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return AuthSigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
