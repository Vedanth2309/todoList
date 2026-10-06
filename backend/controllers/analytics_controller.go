package controllers

import (
	"time"

	"tracker/database"
	"tracker/middleware"
	"tracker/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func Dashboard(c *fiber.Ctx) error {
	uid := middleware.UserID(c)
	ctx, cancel := database.Ctx()
	defer cancel()
	today := time.Now().Format("2006-01-02")
	count := func(col string, f bson.M) int64 {
		f["userId"] = uid
		n, _ := database.DB.Collection(col).CountDocuments(ctx, f)
		return n
	}
	agg := func(col string, pipe []bson.M) []bson.M {
		out := []bson.M{}
		if cur, err := database.DB.Collection(col).Aggregate(ctx, pipe); err == nil {
			cur.All(ctx, &out)
		}
		return out
	}
	group := func(col, field string) []bson.M {
		return agg(col, []bson.M{{"$match": bson.M{"userId": uid}}, {"$group": bson.M{"_id": "$" + field, "count": bson.M{"$sum": 1}}}})
	}
	total, done := count("tasks", bson.M{}), count("tasks", bson.M{"status": "COMPLETED"})
	pct := 0.0
	if total > 0 {
		pct = float64(done) / float64(total) * 100
	}
	over := agg("tasks", []bson.M{
		{"$match": bson.M{"userId": uid, "status": "COMPLETED", "completedAt": bson.M{"$gte": time.Now().AddDate(0, 0, -14)}}},
		{"$group": bson.M{"_id": bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$completedAt"}}, "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"_id": 1}}})
	return utils.OK(c, fiber.Map{
		"totalTasks": total, "completedTasks": done, "completionPct": pct,
		"dueToday": count("tasks", bson.M{"dueDate": today}),
		"overdue":  count("tasks", bson.M{"status": bson.M{"$ne": "COMPLETED"}, "dueDate": bson.M{"$lt": today, "$ne": ""}}),
		"habits":   count("habits", bson.M{}), "habitsDoneToday": count("habit_checkins", bson.M{"completed": true, "date": today}),
		"diaryEntries": count("diary", bson.M{}),
		"byPriority":   group("tasks", "priority"), "byCategory": group("tasks", "category"), "moods": group("diary", "mood"),
		"completedOverTime": over})
}
