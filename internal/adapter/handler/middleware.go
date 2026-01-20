// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestID middleware adds a unique request ID to each request
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request ID exists in header
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in locals and response header
		c.Locals("requestID", requestID)
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}

// Logger middleware logs request details
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log format: [RequestID] STATUS METHOD PATH DURATION
		requestID := c.Locals("requestID")
		if requestID == nil {
			requestID = "-"
		}

		// You can integrate with zerolog here
		// For now, we use the default logger
		_ = duration // Use for logging

		return err
	}
}

// TenantContext middleware extracts tenant information from the request
func TenantContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get tenant from header, subdomain, or JWT
		tenantID := c.Get("X-Tenant-ID")
		if tenantID == "" {
			// Could also extract from subdomain or JWT
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Tenant ID is required",
			})
		}

		// Parse and validate UUID
		id, err := uuid.Parse(tenantID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid Tenant ID format",
			})
		}

		c.Locals("tenantID", id)
		return c.Next()
	}
}

// AuthRequired middleware checks for authentication
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// TODO: Implement JWT verification
		// For now, just check if header exists

		// In production, you would:
		// 1. Parse JWT token
		// 2. Verify signature
		// 3. Check expiration
		// 4. Extract user claims
		// 5. Set user in context

		return c.Next()
	}
}

// RateLimiter middleware implements basic rate limiting
type RateLimiterConfig struct {
	Max        int           // Maximum requests
	Window     time.Duration // Time window
	KeyGetter  func(*fiber.Ctx) string
	LimitReached func(*fiber.Ctx) error
}

// DefaultRateLimiterConfig returns default rate limiter configuration
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Max:    100,
		Window: time.Minute,
		KeyGetter: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},
	}
}

// ErrorHandler is a custom error handler for Fiber
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default error code
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Return JSON error response
	return c.Status(code).JSON(fiber.Map{
		"error":     message,
		"code":      code,
		"requestID": c.Locals("requestID"),
		"architect": "Muhammet-Ali-Buyuk", // Easter egg
	})
}

// SecurityHeaders middleware adds security headers
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("X-Powered-By", "SubFlow-Enterprise") // Custom header
		c.Set("X-Architect", "Muhammet-Ali-Buyuk")  // Signature

		return c.Next()
	}
}
