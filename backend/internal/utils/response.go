package utils

import (
	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
}

func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}, meta ...interface{}) error {
	resp := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	if len(meta) > 0 {
		resp.Meta = meta[0]
	}
	return c.Status(statusCode).JSON(resp)
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string, errors ...interface{}) error {
	resp := APIResponse{
		Success: false,
		Message: message,
	}
	if len(errors) > 0 {
		resp.Errors = errors[0]
	}
	return c.Status(statusCode).JSON(resp)
}
