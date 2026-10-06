package controllers

import (
	"strings"

	"tracker/middleware"
	"tracker/services"
	"tracker/utils"

	"github.com/gofiber/fiber/v2"
)

func Search(c *fiber.Ctx) error {
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		return utils.OK(c, []any{})
	}
	res, err := services.Search(middleware.UserID(c).Hex(), q, c.Query("type"), c.Query("tag"), c.Query("from"), c.Query("to"))
	if err != nil {
		return utils.Fail(c, 502, err.Error())
	}
	return utils.OK(c, res)
}
