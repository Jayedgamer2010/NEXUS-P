package middleware

import (
	"strings"

	"nexus/backend/config"
	"nexus/backend/database"
	"nexus/backend/models"
	"nexus/backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthConfig struct {
	JWTSecret string
}

type Claims struct {
	UserID uint   `json:"user_id"`
	UUID   string `json:"uuid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func Auth(cfg *AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.Unauthorized(c, "Authorization header required")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			return utils.Unauthorized(c, "Invalid authorization header format")
		}

		tokenString := parts[1]
		tempCfg := &config.Config{JWTSecret: cfg.JWTSecret}
		claims, err := validateToken(tokenString, tempCfg)
		if err != nil {
			return utils.Unauthorized(c, "Invalid or expired token")
		}

		var user models.User
		if err := database.DB.First(&user, "id = ?", claims.UserID).Error; err != nil {
			return utils.Unauthorized(c, "User not found")
		}

		c.Locals("user", user)
		return c.Next()
	}
}

func GetUser(c *fiber.Ctx) *models.User {
	if user, ok := c.Locals("user").(models.User); ok {
		return &user
	}
	return nil
}

func validateToken(tokenString string, cfg *config.Config) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, nil
}
