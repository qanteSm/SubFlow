// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/qantesm/subflow/internal/core/service"
)

// ProjectHandler handles HTTP requests for project operations
type ProjectHandler struct {
	// Services would be injected here
	calculator *service.Calculator
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(calc *service.Calculator) *ProjectHandler {
	return &ProjectHandler{
		calculator: calc,
	}
}

// RegisterRoutes registers all project-related routes
func (h *ProjectHandler) RegisterRoutes(router fiber.Router) {
	projects := router.Group("/projects")
	
	projects.Get("/", h.ListProjects)
	projects.Post("/", h.CreateProject)
	projects.Get("/:id", h.GetProject)
	projects.Put("/:id", h.UpdateProject)
	projects.Delete("/:id", h.DeleteProject)
	projects.Get("/:id/financials/summary", h.GetFinancialSummary)
}

// ListProjects returns all projects for the current tenant
// @Summary List all projects
// @Tags Projects
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /projects [get]
func (h *ProjectHandler) ListProjects(c *fiber.Ctx) error {
	// TODO: Implement with repository
	return c.JSON(fiber.Map{
		"data":    []interface{}{},
		"message": "Projects retrieved successfully",
	})
}

// CreateProject creates a new project
// @Summary Create a new project
// @Tags Projects
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Router /projects [post]
func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	// TODO: Implement with validation and repository
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Project created successfully",
	})
}

// GetProject retrieves a single project by ID
// @Summary Get project by ID
// @Tags Projects
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} map[string]interface{}
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetProject(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{
		"id":      id,
		"message": "Project retrieved",
	})
}

// UpdateProject updates an existing project
// @Summary Update project
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} map[string]interface{}
// @Router /projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{
		"id":      id,
		"message": "Project updated",
	})
}

// DeleteProject soft-deletes a project
// @Summary Delete project
// @Tags Projects
// @Param id path string true "Project ID"
// @Success 204
// @Router /projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// CalculateBillingRequest represents the input for AIA billing calculation
type CalculateBillingRequest struct {
	OriginalContractSum   int64 `json:"original_contract_sum"`
	ApprovedChangeOrders  int64 `json:"approved_change_orders"`
	PreviousWorkCompleted int64 `json:"previous_work_completed"`
	CurrentWorkCompleted  int64 `json:"current_work_completed"`
	StoredMaterials       int64 `json:"stored_materials"`
	PreviousCertificates  int64 `json:"previous_certificates"`
	LaborRetainageRate    int64 `json:"labor_retainage_rate"`    // Basis points
	MaterialRetainageRate int64 `json:"material_retainage_rate"` // Basis points
}

// GetFinancialSummary returns the financial snapshot of a project
// @Summary Get project financial summary
// @Tags Projects
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} service.AIABillingResult
// @Router /projects/{id}/financials/summary [get]
func (h *ProjectHandler) GetFinancialSummary(c *fiber.Ctx) error {
	// Example calculation with sample data
	input := service.AIABillingInput{
		OriginalContractSum:   100000000, // $1,000,000.00
		ApprovedChangeOrders:  5000000,   // $50,000.00
		PreviousWorkCompleted: 30000000,  // $300,000.00
		CurrentWorkCompleted:  15000000,  // $150,000.00
		StoredMaterials:       5000000,   // $50,000.00
		PreviousCertificates:  25000000,  // $250,000.00
		LaborRetainageRate:    1000,      // 10%
		MaterialRetainageRate: 500,       // 5%
	}

	result, err := h.calculator.Calculate(input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"project_id": c.Params("id"),
		"summary":    result,
		"formatted": fiber.Map{
			"contract_sum":        service.FormatCurrency(result.ContractSum, "USD"),
			"current_payment_due": service.FormatCurrency(result.CurrentPaymentDue, "USD"),
			"total_retainage":     service.FormatCurrency(result.TotalRetainage, "USD"),
		},
	})
}
