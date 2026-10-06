package routes

import (
	"tracker/controllers"
	"tracker/middleware"

	"github.com/gofiber/fiber/v2"
)

func crud(r fiber.Router, path string, h controllers.CRUD) {
	r.Get("/"+path, h.List)
	r.Post("/"+path, h.Create)
	r.Put("/"+path+"/:id", h.Update)
	r.Delete("/"+path+"/:id", h.Delete)
}

func Setup(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/auth/register", controllers.Register)
	api.Post("/auth/login", controllers.Login)

	p := api.Group("", middleware.Auth)
	p.Get("/auth/me", controllers.Me)
	crud(p, "tasks", controllers.Tasks)
	crud(p, "habits", controllers.Habits)
	p.Post("/habits/:id/checkin", controllers.Checkin)
	p.Get("/habits/:id/checkins", controllers.Checkins)
	crud(p, "events", controllers.Events)
	crud(p, "reminders", controllers.Reminders)
	crud(p, "diary", controllers.Diary)
	crud(p, "notes", controllers.Notes)
	p.Get("/search/autocomplete", controllers.Autocomplete)
	p.Get("/search", controllers.Search)
	p.Get("/analytics/dashboard", controllers.Dashboard)
	p.Get("/settings", controllers.Me)
	p.Put("/settings", controllers.UpdateSettings)
}
