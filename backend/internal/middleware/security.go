package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"kasirpro/internal/config"
)

func SecurityMiddleware() fiber.Handler {
	return helmet.New(helmet.Config{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https: blob:;",
	})
}

func CORSMiddleware() fiber.Handler {
	cfg := config.AppConfig
	origins := cfg.CORSAllowedOrigins
	if origins == "" {
		origins = "https://fish-warming-logos-lots.trycloudflare.com,http://localhost:5173,http://localhost:3000,http://localhost,https://*.vercel.app,https://*.railway.app"
	}

	return cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			if origins == "*" {
				return true
			}
			originList := strings.Split(origins, ",")
			for _, o := range originList {
				trimmed := strings.TrimSpace(o)
				if trimmed == "*" || trimmed == origin {
					return true
				}
				// Support vercel, railway, and cloudflare tunnel subdomain matching
				if strings.Contains(trimmed, "vercel.app") && strings.HasSuffix(origin, ".vercel.app") {
					return true
				}
				if strings.Contains(trimmed, "railway.app") && strings.HasSuffix(origin, ".railway.app") {
					return true
				}
				if strings.Contains(trimmed, "trycloudflare.com") && strings.HasSuffix(origin, ".trycloudflare.com") {
					return true
				}
			}
			return false
		},
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET, POST, PUT, DELETE, PATCH, OPTIONS",
		AllowCredentials: true,
	})
}

func StripServerHeader(c *fiber.Ctx) error {
	c.Set("Server", "KasirPro Engine")
	return c.Next()
}

func SanitizeInputMiddleware(c *fiber.Ctx) error {
	// Ensure Content-Type is valid for mutations
	if c.Method() == "POST" || c.Method() == "PUT" || c.Method() == "PATCH" {
		contentType := string(c.Request().Header.ContentType())
		if contentType != "" && !strings.Contains(contentType, "application/json") && !strings.Contains(contentType, "multipart/form-data") {
			return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
				"success": false,
				"message": "Unsupported Media Type. Hanya menerima application/json atau multipart/form-data",
			})
		}
	}
	return c.Next()
}
