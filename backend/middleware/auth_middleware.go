package middleware

import (
	"strings"

	"tracker/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Auth(c *fiber.Ctx) error {
	h := c.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return utils.Fail(c, 401, "Missing or invalid Authorization header")
	}
	id, err := utils.ParseToken(strings.TrimPrefix(h, "Bearer "))
	if err != nil {
		return utils.Fail(c, 401, "Invalid or expired token")
	}
	c.Locals("userID", id)
	return c.Next()
}

// UserID returns the authenticated user's ID (only ever taken from the JWT).
func UserID(c *fiber.Ctx) primitive.ObjectID { return c.Locals("userID").(primitive.ObjectID) }
