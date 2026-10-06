package controllers

import (
	"tracker/database"
	"tracker/middleware"
	"tracker/models"
	"tracker/utils"

	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Checkin(c *fiber.Ctx) error {
	hid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.Fail(c, 400, "Invalid id")
	}
	var b struct {
		Date      string
		Completed bool
		Note      string
	}
	if c.BodyParser(&b) != nil {
		return utils.Fail(c, 400, "Invalid request body")
	}
	if b.Date == "" {
		b.Date = time.Now().Format("2006-01-02")
	}
	uid := middleware.UserID(c)
	ctx, cancel := database.Ctx()
	defer cancel()
	// the habit must belong to the caller
	if n, _ := database.DB.Collection("habits").CountDocuments(ctx, bson.M{"_id": hid, "userId": uid}); n == 0 {
		return utils.Fail(c, 404, "Habit not found")
	}
	_, err = database.DB.Collection("habit_checkins").UpdateOne(ctx, bson.M{"userId": uid, "habitId": hid, "date": b.Date},
		bson.M{"$set": bson.M{"completed": b.Completed, "note": b.Note}}, options.Update().SetUpsert(true))
	if err != nil {
		return utils.Fail(c, 500, "Could not save check-in")
	}
	return utils.OK(c, fiber.Map{"date": b.Date, "completed": b.Completed})
}

func Checkins(c *fiber.Ctx) error {
	hid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.Fail(c, 400, "Invalid id")
	}
	ctx, cancel := database.Ctx()
	defer cancel()
	cur, err := database.DB.Collection("habit_checkins").Find(ctx, bson.M{"userId": middleware.UserID(c), "habitId": hid})
	out := []models.HabitCheckin{}
	if err != nil || cur.All(ctx, &out) != nil {
		return utils.Fail(c, 500, "Could not load check-ins")
	}
	return utils.OK(c, out)
}
