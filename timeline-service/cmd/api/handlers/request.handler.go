package handlers

import "github.com/gofiber/fiber/v2"

func SuccessRouter(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": false, "message": message, "data": data})
}

func ErrorRouter(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": message})
}

func ServerErrorRotuer(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Internal Server Error"})
}

func InvalidUserRouter(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": true, "message": "Invalid user | Please login again"})
}
