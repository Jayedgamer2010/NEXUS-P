package utils

import (
	"errors"
	"time"

	"nexus/backend/config"
	"nexus/backend/models"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    uint   `json:"user_id"`
	UUID      string `json:"uuid"`
	Role      string `json:"role"`
	RootAdmin bool   `json:"root_admin"`
	jwt.RegisteredClaims
}

func GenerateToken(user models.User, cfg *config.Config) (string, error) {
	expireHours := cfg.JWTExpire
	if expireHours <= 0 {
		expireHours = 72
	}

	claims := Claims{
		UserID:    user.ID,
		UUID:      user.UUID,
		Role:      user.Role,
		RootAdmin: user.RootAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func ValidateToken(tokenString string, cfg *config.Config) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
