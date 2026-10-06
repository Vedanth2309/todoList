package middleware

import (
	"errors"
	"log"

	"tracker/utils"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code, msg := fiber.StatusInternalServerError, "Internal server error"
	var fe *fiber.Error
	if errors.As(err, &fe) {
		code, msg = fe.Code, fe.Message
	} else {
		log.Println("unhandled error:", err)
	}
	return utils.Fail(c, code, msg)
}
