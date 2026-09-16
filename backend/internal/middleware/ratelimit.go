package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// AuthRateLimiter restricts login/register attempts to 5 per minute per IP
func AuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Terlalu banyak percobaan autentikasi. Silakan tunggu 1 menit lagi.",
			})
		},
	})
}

// GeneralAPIRateLimiter allows 60 requests per minute per authenticated user or IP
func GeneralAPIRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        60,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			if userID := c.Locals("user_id"); userID != nil {
				return fmt.Sprintf("user_%v", userID)
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Batas permintaan tercapai. Silakan coba sesaat lagi.",
			})
		},
	})
}

// POSRateLimiter allows 30 transactions per minute
func POSRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        30,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			if userID := c.Locals("user_id"); userID != nil {
				return fmt.Sprintf("pos_user_%v", userID)
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Kasir sibuk, antrean checkout terlalu cepat. Tunggu beberapa detik.",
			})
		},
	})
}
