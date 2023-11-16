package middleware

import (
	"product/config"
	"product/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user models.User, expLimit time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * expLimit).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(config.Cfg.JWT_SECRET_KEY))

	return t, err
}
