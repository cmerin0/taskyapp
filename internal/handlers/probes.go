package handlers

import (
	"context"
	"time"

	"github.com/cmerin0/tasky/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// Healthcheck returns the status of the application
// @Summary Healthcheck
// @Description Check the health of the application
// @Success 200 {object} fiber.Map
// @Failure 500 Internal Server Error
// @Router /health [get]
func Healthcheck(c *fiber.Ctx) error {
	log.Info("Healthcheck endpoint hit")
	return c.JSON(fiber.Map{
		"status":  "UP",
		"version": "1.0.0",
	})
}

// ReadinessProbe checks if the application is ready to serve requests
// @Summary Readiness Probe
// @Description Check if the application is ready to serve requests
// @Success 200 {object} fiber.Map
// @Failure 503 {object} fiber.Map
// @Router /readyz [get]
func ReadinessProbe(c *fiber.Ctx) error {
	// Check MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Check if the MongoDB client is connected
	if err := db.Client.Ping(ctx, nil); err != nil {
		log.Error("MongoDB not connected: ", err)
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "DOWN",
			"error":  "MongoDB not connected",
		})
	}

	log.Info("MongoDB connected")
	return c.JSON(fiber.Map{
		"status": "READY",
	})
}

// LivenessProbe checks if the application is alive
// @Summary Liveness Probe
// @Description Check if the application is alive
// @Success 200 {object} fiber.Map
// @Router /healthz [get]
func LivenessProbe(c *fiber.Ctx) error {
	log.Info("Liveness probe hit")
	return c.JSON(fiber.Map{
		"status": "ALIVE",
	})
}
