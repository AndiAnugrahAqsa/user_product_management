package handlers

import "github.com/gofiber/fiber/v2"

func response(c *fiber.Ctx, statusCode int, message string, data any) error {
	if data != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"message": message,
			"data":    data,
		})
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"message": message,
	})
}
