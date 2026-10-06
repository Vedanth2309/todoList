package controllers

import (
	"strings"

	"tracker/database"
	"tracker/middleware"
	"tracker/models"
	"tracker/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func UpdateSettings(c *fiber.Ctx) error {
	var b struct {
		Name, Email, Password *string
		Settings              *models.Settings
	}
	if c.BodyParser(&b) != nil {
		return utils.Fail(c, 400, "Invalid request body")
	}
	set := bson.M{}
	if b.Name != nil && strings.TrimSpace(*b.Name) != "" {
		set["name"] = strings.TrimSpace(*b.Name)
	}
	if b.Email != nil && strings.Contains(*b.Email, "@") {
		set["email"] = strings.ToLower(strings.TrimSpace(*b.Email))
	}
	if b.Password != nil && *b.Password != "" {
		if len(*b.Password) < 6 {
			return utils.Fail(c, 400, "Password must be at least 6 characters")
		}
		h, _ := bcrypt.GenerateFromPassword([]byte(*b.Password), 10)
		set["password"] = string(h)
	}
	if b.Settings != nil {
		set["settings"] = b.Settings
	}
	if len(set) > 0 {
		ctx, cancel := database.Ctx()
		defer cancel()
		if _, err := database.DB.Collection("users").UpdateOne(ctx, bson.M{"_id": middleware.UserID(c)}, bson.M{"$set": set}); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return utils.Fail(c, 409, "Email already in use")
			}
			return utils.Fail(c, 500, "Could not update settings")
		}
	}
	return Me(c)
}
