package main

import (
	"fmt"
	"os"

	"github.com/cmerin0/tasky/internal/db"
	"github.com/cmerin0/tasky/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env only in development
	if os.Getenv("GO_ENV") != "prod" {
		_ = godotenv.Load()
	}

	// Constructing the MongoDB URI
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_DBNAME"),
	)

	// Connect to the database
	// Note: The ConnectDB function should be called only once
	// to avoid multiple connections to the database.
	log.Info("Connecting to MongoDB...")
	db.ConnectDB(mongoURI)

	// Then create the app
	app := fiber.New()
	app.Use(logger.New())

	// Routes setup
	setupRoutes(app)

	// Start the server
	log.Info("Starting server on port ", os.Getenv("APP_PORT"))
	log.Info(app.Listen(":" + os.Getenv("APP_PORT")))
}

func setupRoutes(app *fiber.App) {

	// Main Route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Tasky API")
	})

	// Main API group
	api := app.Group("/api/v1")

	// Health check routes
	api.Get("/health", handlers.Healthcheck)
	api.Get("/readyz", handlers.ReadinessProbe)
	api.Get("/healthz", handlers.LivenessProbe)

	// User routes
	users := api.Group("/users")
	users.Get("/", handlers.GetUsers)
	users.Post("/", handlers.CreateUser)
	users.Get("/:userId", handlers.GetUser)
	users.Put("/:userId", handlers.UpdateUser)
	users.Delete("/:userId", handlers.DeleteUser)

	// Task routes
	tasks := api.Group("/tasks")
	tasks.Get("/", handlers.ListTasks)
	tasks.Post("/", handlers.CreateTask)
	tasks.Get("/:taskId", handlers.GetTask)
	tasks.Get("/user/:userId", handlers.GetUserTasks)
	tasks.Put("/:taskId", handlers.UpdateTask)
	tasks.Delete("/:taskId", handlers.DeleteTask)
}
