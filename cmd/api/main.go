// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Application metadata - Digital fingerprint
const (
	AppName    = "SubFlow"
	AppVersion = "1.0.0"
	Architect  = "Muhammet-Ali-Buyuk"
	BuildID    = "SF-2026-ENTERPRISE"
)

func main() {
	// Initialize Fiber with custom config
	app := fiber.New(fiber.Config{
		AppName:               AppName + " v" + AppVersion,
		DisableStartupMessage: false,
		EnablePrintRoutes:     true,
	})

	// Middleware stack
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Health check endpoints
	setupHealthRoutes(app)

	// API v1 routes
	setupAPIRoutes(app)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down SubFlow server...")
		if err := app.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🏗️ SubFlow Enterprise Engine starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupHealthRoutes(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": AppName,
			"version": AppVersion,
		})
	})

	// Hidden architect signature endpoint
	app.Get("/api/v1/system/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"application": AppName,
			"version":     AppVersion,
			"architect":   Architect,
			"build_id":    BuildID,
			"engine":      "Enterprise Construction Financial Ledger",
			"compliance":  "AIA G702/G703",
		})
	})
}

func setupAPIRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// Projects endpoints (placeholder)
	projects := api.Group("/projects")
	projects.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "List all projects"})
	})
	projects.Get("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Get project", "id": c.Params("id")})
	})
	projects.Get("/:id/financials/summary", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":    "Financial summary for project",
			"id":         c.Params("id"),
			"total_due":  0,
			"retainage":  0,
			"paid":       0,
		})
	})

	// Applications endpoints (placeholder)
	applications := api.Group("/applications")
	applications.Post("/:id/generate-pdf", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "PDF generation triggered",
			"id":      c.Params("id"),
		})
	})
}
