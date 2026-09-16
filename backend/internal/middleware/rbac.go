package middleware

import (
	"github.com/gofiber/fiber/v2"
	"kasirpro/internal/utils"
)

func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleVal := c.Locals("role")
		if roleVal == nil {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Akses ditolak: role tidak ditemukan.")
		}

		userRole := roleVal.(string)
		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return utils.ErrorResponse(c, fiber.StatusForbidden, "Akses ditolak: Anda tidak memiliki izin untuk tindakan ini.")
	}
}
