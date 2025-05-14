package handlers

import (
	"context"
	"time"

	"github.com/cmerin0/tasky/internal/db"
	"github.com/gofiber/fiber/v2"
)

// Add these new handlers
func Healthcheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "UP",
		"version": "1.0.0",
	})
}

func ReadinessProbe(c *fiber.Ctx) error {
	// Check MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.Client.Ping(ctx, nil); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "DOWN",
			"error":  "MongoDB not connected",
		})
	}

	return c.JSON(fiber.Map{
		"status": "READY",
	})
}

func LivenessProbe(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ALIVE",
	})
}
