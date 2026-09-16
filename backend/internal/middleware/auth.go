package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"kasirpro/internal/utils"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string

		// 1. Check Authorization header
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString = parts[1]
			}
		}

		// 2. Fallback to HttpOnly Cookie
		if tokenString == "" {
			tokenString = c.Cookies("access_token")
		}

		if tokenString == "" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Sesi login tidak ditemukan atau kedaluwarsa. Silakan login kembali.")
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil || claims.TokenType != "access" {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token autentikasi tidak valid atau sudah kedaluwarsa.")
		}

		// Inject to Locals
		c.Locals("user_id", claims.UserID)
		c.Locals("outlet_id", claims.OutletID)
		c.Locals("role", claims.Role)
		c.Locals("email", claims.Email)
		c.Locals("name", claims.Name)
		c.Locals("claims", claims)

		return c.Next()
	}
}
