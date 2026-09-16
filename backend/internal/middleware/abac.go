package middleware

import (
	"github.com/gofiber/fiber/v2"
	"kasirpro/internal/dto"
	"kasirpro/internal/utils"
)

// OutletIsolation ensures that all operations are scoped to the user's assigned outlet
func OutletIsolation() fiber.Handler {
	return func(c *fiber.Ctx) error {
		outletID := c.Locals("outlet_id")
		if outletID == nil || outletID.(uint) == 0 {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Akses ditolak: outlet tenant tidak teridentifikasi.")
		}
		return c.Next()
	}
}

// DiscountAuthorizationPolicy: If discount > 20%, only "owner" can authorize
func DiscountAuthorizationPolicy() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)

		// Parse body peek if this is a checkout request
		if c.Method() == "POST" && (c.Path() == "/api/v1/pos/checkout" || c.Path() == "/api/pos/checkout") {
			var body dto.CheckoutRequest
			if err := c.BodyParser(&body); err == nil {
				// Check discount percentage
				if body.DiscountType == "percentage" && body.DiscountValue > 20.0 {
					if role != "owner" {
						return utils.ErrorResponse(c, fiber.StatusForbidden, "Diskon di atas 20% memerlukan otorisasi Owner toko.")
					}
				}
			}
		}

		return c.Next()
	}
}
