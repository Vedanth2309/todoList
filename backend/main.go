package main

import (
	"log"

	"tracker/config"
	"tracker/database"
	"tracker/middleware"
	"tracker/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	log.Println("Loading environment variables...")
	if err := config.Load(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}
	log.Println("Connecting to MongoDB...")
	if _, _, err := database.ConnectMongoDB(); err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}
	log.Println("MongoDB connected successfully (database:", config.C.MongoDB+")")
	database.InitIndexes()

	log.Println("Connecting to Elasticsearch...")
	if err := config.ConnectElasticsearch(); err != nil {
		log.Println("Elasticsearch unavailable (search will not work):", err)
	} else {
		log.Println("Elasticsearch connected successfully")
	}

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
	app.Use(recover.New(), cors.New())
	routes.Setup(app)
	log.Printf("Starting Fiber server on :%s", config.C.Port)
	log.Fatal(app.Listen(":" + config.C.Port))
}
