package controllers

import (
	"errors"
	"time"

	"tracker/database"
	"tracker/middleware"
	"tracker/models"
	"tracker/services"
	"tracker/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CRUD struct{ List, Create, Update, Delete fiber.Handler }

// NewCRUD builds the four handlers for a user-owned collection. Every query is filtered by the JWT user ID.
func NewCRUD[T any, P interface {
	*T
	models.Doc
}](col string, onUpdate func(bson.M)) CRUD {
	coll := func() *mongo.Collection { return database.DB.Collection(col) }
	return CRUD{
		List: func(c *fiber.Ctx) error {
			ctx, cancel := database.Ctx()
			defer cancel()
			f := bson.M{"userId": middleware.UserID(c)}
			for _, k := range []string{"status", "priority", "category", "mood", "completed", "date"} {
				if v := c.Query(k); v == "true" || v == "false" {
					f[k] = v == "true"
				} else if v != "" {
					f[k] = v
				}
			}
			if t := c.Query("tag"); t != "" {
				f["tags"] = t
			}
			cur, err := coll().Find(ctx, f, options.Find().SetSort(bson.M{"createdAt": -1}))
			if err != nil {
				return utils.Fail(c, 500, "Could not load "+col)
			}
			items := []T{}
			if err := cur.All(ctx, &items); err != nil {
				return utils.Fail(c, 500, "Could not read "+col)
			}
			return utils.OK(c, items)
		},
		Create: func(c *fiber.Ctx) error {
			var t T
			p := P(&t)
			if err := c.BodyParser(p); err != nil {
				return utils.Fail(c, 400, "Invalid request body")
			}
			if msg := p.Validate(); msg != "" {
				return utils.Fail(c, 400, msg)
			}
			b, now := p.B(), time.Now()
			b.ID, b.UserID, b.CreatedAt, b.UpdatedAt = primitive.NewObjectID(), middleware.UserID(c), now, now
			ctx, cancel := database.Ctx()
			defer cancel()
			if _, err := coll().InsertOne(ctx, p); err != nil {
				return utils.Fail(c, 500, "Could not save "+col)
			}
			services.IndexDoc(col, p)
			return utils.OK(c, p)
		},
		Update: func(c *fiber.Ctx) error {
			id, err := primitive.ObjectIDFromHex(c.Params("id"))
			if err != nil {
				return utils.Fail(c, 400, "Invalid id")
			}
			var d bson.M
			if c.BodyParser(&d) != nil || d == nil {
				return utils.Fail(c, 400, "Invalid request body")
			}
			for _, k := range []string{"_id", "id", "userId", "createdAt"} {
				delete(d, k)
			}
			d["updatedAt"] = time.Now()
			if onUpdate != nil {
				onUpdate(d)
			}
			ctx, cancel := database.Ctx()
			defer cancel()
			var out T
			err = coll().FindOneAndUpdate(ctx, bson.M{"_id": id, "userId": middleware.UserID(c)}, bson.M{"$set": d},
				options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&out)
			if errors.Is(err, mongo.ErrNoDocuments) {
				return utils.Fail(c, 404, "Not found")
			} else if err != nil {
				return utils.Fail(c, 500, "Could not update "+col)
			}
			services.IndexDoc(col, P(&out))
			return utils.OK(c, out)
		},
		Delete: func(c *fiber.Ctx) error {
			id, err := primitive.ObjectIDFromHex(c.Params("id"))
			if err != nil {
				return utils.Fail(c, 400, "Invalid id")
			}
			ctx, cancel := database.Ctx()
			defer cancel()
			res, err := coll().DeleteOne(ctx, bson.M{"_id": id, "userId": middleware.UserID(c)})
			if err != nil {
				return utils.Fail(c, 500, "Could not delete "+col)
			}
			if res.DeletedCount == 0 {
				return utils.Fail(c, 404, "Not found")
			}
			services.DeleteDoc(col, id.Hex())
			return utils.OK(c, fiber.Map{"deleted": true})
		},
	}
}

func taskUpdate(d bson.M) {
	if s, ok := d["status"]; ok {
		if s == "COMPLETED" {
			d["completedAt"] = time.Now()
		} else {
			d["completedAt"] = nil
		}
	}
}

var (
	Tasks     = NewCRUD[models.Task]("tasks", taskUpdate)
	Habits    = NewCRUD[models.Habit]("habits", nil)
	Events    = NewCRUD[models.Event]("events", nil)
	Reminders = NewCRUD[models.Reminder]("reminders", nil)
	Diary     = NewCRUD[models.Diary]("diary", nil)
	Notes     = NewCRUD[models.Note]("notes", nil)
)
