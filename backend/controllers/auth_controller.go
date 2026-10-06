package controllers

import (
	"strings"
	"time"

	"tracker/database"
	"tracker/middleware"
	"tracker/models"
	"tracker/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func session(c *fiber.Ctx, u models.User) error {
	t, err := utils.GenerateToken(u.ID)
	if err != nil {
		return utils.Fail(c, 500, "Could not create token")
	}
	return utils.OK(c, fiber.Map{"token": t, "user": u})
}

func Register(c *fiber.Ctx) error {
	var b struct{ Name, Email, Password string }
	if c.BodyParser(&b) != nil || strings.TrimSpace(b.Name) == "" || !strings.Contains(b.Email, "@") || len(b.Password) < 6 {
		return utils.Fail(c, 400, "Name, a valid email and a password of 6+ characters are required")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(b.Password), 10)
	if err != nil {
		return utils.Fail(c, 500, "Could not hash password")
	}
	u := models.User{ID: primitive.NewObjectID(), Name: strings.TrimSpace(b.Name), Email: strings.ToLower(strings.TrimSpace(b.Email)), Password: string(h),
		Settings: models.Settings{Theme: "light", Timezone: "UTC", DefaultPriority: "MEDIUM", Notifications: true}, CreatedAt: time.Now()}
	ctx, cancel := database.Ctx()
	defer cancel()
	if _, err := database.DB.Collection("users").InsertOne(ctx, u); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return utils.Fail(c, 409, "Email already registered")
		}
		return utils.Fail(c, 500, "Could not create user")
	}
	return session(c, u)
}

func Login(c *fiber.Ctx) error {
	var b struct{ Email, Password string }
	c.BodyParser(&b)
	var u models.User
	ctx, cancel := database.Ctx()
	defer cancel()
	err := database.DB.Collection("users").FindOne(ctx, bson.M{"email": strings.ToLower(strings.TrimSpace(b.Email))}).Decode(&u)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(b.Password)) != nil {
		return utils.Fail(c, 401, "Invalid email or password")
	}
	return session(c, u)
}

func Me(c *fiber.Ctx) error {
	var u models.User
	ctx, cancel := database.Ctx()
	defer cancel()
	if database.DB.Collection("users").FindOne(ctx, bson.M{"_id": middleware.UserID(c)}).Decode(&u) != nil {
		return utils.Fail(c, 404, "User not found")
	}
	return utils.OK(c, u)
}
