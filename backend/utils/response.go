package utils

import "github.com/gofiber/fiber/v2"

func OK(c *fiber.Ctx, data any) error {
	return c.JSON(fiber.Map{"success": true, "data": data})
}
func Fail(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{"success": false, "message": msg})
}
