// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/qantesm/subflow/internal/core/entity"
	"github.com/qantesm/subflow/internal/core/service"
)

// TransactionHandler handles HTTP requests for transaction operations
type TransactionHandler struct {
	ledgerService *service.LedgerService
	calculator    *service.Calculator
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(ledger *service.LedgerService, calc *service.Calculator) *TransactionHandler {
	return &TransactionHandler{
		ledgerService: ledger,
		calculator:    calc,
	}
}

// RegisterRoutes registers all transaction-related routes
func (h *TransactionHandler) RegisterRoutes(router fiber.Router) {
	transactions := router.Group("/transactions")

	transactions.Get("/project/:projectId", h.ListByProject)
	transactions.Post("/invoice", h.CreateInvoice)
	transactions.Post("/payment", h.CreatePayment)
	transactions.Post("/retainage/hold", h.HoldRetainage)
	transactions.Post("/retainage/release", h.ReleaseRetainage)

	// Calculator endpoints
	router.Post("/calculate/aia", h.CalculateAIA)
}

// CreateInvoiceRequest represents the request body for creating an invoice
type CreateInvoiceRequest struct {
	ProjectID  string `json:"project_id" validate:"required,uuid"`
	Amount     int64  `json:"amount" validate:"required,gt=0"` // In cents
	Currency   string `json:"currency" validate:"required,len=3"`
	InvoiceNo  string `json:"invoice_no" validate:"required"`
	Description string `json:"description"`
}

// CreatePaymentRequest represents the request body for creating a payment
type CreatePaymentRequest struct {
	ProjectID     string `json:"project_id" validate:"required,uuid"`
	Amount        int64  `json:"amount" validate:"required,gt=0"`
	Currency      string `json:"currency" validate:"required,len=3"`
	BankReceiptNo string `json:"bank_receipt_no" validate:"required"`
	Description   string `json:"description"`
}

// RetainageRequest represents the request body for retainage operations
type RetainageRequest struct {
	ProjectID string  `json:"project_id" validate:"required,uuid"`
	Amount    int64   `json:"amount" validate:"required,gt=0"`
	Currency  string  `json:"currency" validate:"required,len=3"`
	Rate      float64 `json:"rate,omitempty"`
}

// TransactionResponse is the standard response for transaction operations
type TransactionResponse struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	Type        string `json:"type"`
	AmountCents int64  `json:"amount_cents"`
	Currency    string `json:"currency"`
	Message     string `json:"message"`
}

// ListByProject returns all transactions for a project
// @Summary List transactions by project
// @Tags Transactions
// @Produce json
// @Param projectId path string true "Project ID"
// @Success 200 {array} entity.Transaction
// @Router /transactions/project/{projectId} [get]
func (h *TransactionHandler) ListByProject(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	transactions, err := h.ledgerService.GetTransactionHistory(c.Context(), projectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":  transactions,
		"count": len(transactions),
	})
}

// CreateInvoice creates a new invoice transaction
// @Summary Create an invoice
// @Tags Transactions
// @Accept json
// @Produce json
// @Param request body CreateInvoiceRequest true "Invoice details"
// @Success 201 {object} TransactionResponse
// @Router /transactions/invoice [post]
func (h *TransactionHandler) CreateInvoice(c *fiber.Ctx) error {
	var req CreateInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	// TODO: Get actual user ID from auth context
	userID := uuid.New() // Placeholder

	tx, err := h.ledgerService.RecordInvoice(c.Context(), projectID, req.Amount, req.Currency, req.InvoiceNo, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(TransactionResponse{
		ID:          tx.ID.String(),
		ProjectID:   tx.ProjectID.String(),
		Type:        string(tx.Type),
		AmountCents: tx.AmountCents,
		Currency:    tx.Currency,
		Message:     "Invoice created successfully",
	})
}

// CreatePayment creates a new payment transaction
// @Summary Record a payment
// @Tags Transactions
// @Accept json
// @Produce json
// @Param request body CreatePaymentRequest true "Payment details"
// @Success 201 {object} TransactionResponse
// @Router /transactions/payment [post]
func (h *TransactionHandler) CreatePayment(c *fiber.Ctx) error {
	var req CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	userID := uuid.New() // Placeholder

	tx, err := h.ledgerService.RecordPayment(c.Context(), projectID, req.Amount, req.Currency, req.BankReceiptNo, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(TransactionResponse{
		ID:          tx.ID.String(),
		ProjectID:   tx.ProjectID.String(),
		Type:        string(tx.Type),
		AmountCents: tx.AmountCents,
		Currency:    tx.Currency,
		Message:     "Payment recorded successfully",
	})
}

// HoldRetainage creates a retainage held transaction
// @Summary Hold retainage
// @Tags Transactions
// @Accept json
// @Produce json
// @Param request body RetainageRequest true "Retainage details"
// @Success 201 {object} TransactionResponse
// @Router /transactions/retainage/hold [post]
func (h *TransactionHandler) HoldRetainage(c *fiber.Ctx) error {
	var req RetainageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	userID := uuid.New()

	tx, err := h.ledgerService.RecordRetainageHeld(c.Context(), projectID, req.Amount, req.Currency, req.Rate, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(TransactionResponse{
		ID:          tx.ID.String(),
		ProjectID:   tx.ProjectID.String(),
		Type:        string(tx.Type),
		AmountCents: tx.AmountCents,
		Currency:    tx.Currency,
		Message:     "Retainage held successfully",
	})
}

// ReleaseRetainage creates a retainage release transaction
// @Summary Release retainage
// @Tags Transactions
// @Accept json
// @Produce json
// @Param request body RetainageRequest true "Retainage details"
// @Success 201 {object} TransactionResponse
// @Router /transactions/retainage/release [post]
func (h *TransactionHandler) ReleaseRetainage(c *fiber.Ctx) error {
	var req RetainageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID",
		})
	}

	userID := uuid.New()

	tx, err := h.ledgerService.RecordRetainageRelease(c.Context(), projectID, req.Amount, req.Currency, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(TransactionResponse{
		ID:          tx.ID.String(),
		ProjectID:   tx.ProjectID.String(),
		Type:        string(tx.Type),
		AmountCents: tx.AmountCents,
		Currency:    tx.Currency,
		Message:     "Retainage released successfully",
	})
}

// AIACalculateRequest represents the request body for AIA calculation
type AIACalculateRequest struct {
	OriginalContractSum   int64 `json:"original_contract_sum"`
	ApprovedChangeOrders  int64 `json:"approved_change_orders"`
	PreviousWorkCompleted int64 `json:"previous_work_completed"`
	CurrentWorkCompleted  int64 `json:"current_work_completed"`
	StoredMaterials       int64 `json:"stored_materials"`
	PreviousCertificates  int64 `json:"previous_certificates"`
	LaborRetainageRate    int64 `json:"labor_retainage_rate"`    // Basis points (1000 = 10%)
	MaterialRetainageRate int64 `json:"material_retainage_rate"` // Basis points
}

// CalculateAIA performs AIA G702/G703 billing calculation
// @Summary Calculate AIA billing
// @Tags Calculator
// @Accept json
// @Produce json
// @Param request body AIACalculateRequest true "Billing input"
// @Success 200 {object} service.AIABillingResult
// @Router /calculate/aia [post]
func (h *TransactionHandler) CalculateAIA(c *fiber.Ctx) error {
	var req AIACalculateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	input := service.AIABillingInput{
		OriginalContractSum:   req.OriginalContractSum,
		ApprovedChangeOrders:  req.ApprovedChangeOrders,
		PreviousWorkCompleted: req.PreviousWorkCompleted,
		CurrentWorkCompleted:  req.CurrentWorkCompleted,
		StoredMaterials:       req.StoredMaterials,
		PreviousCertificates:  req.PreviousCertificates,
		LaborRetainageRate:    req.LaborRetainageRate,
		MaterialRetainageRate: req.MaterialRetainageRate,
	}

	result, err := h.calculator.Calculate(input)
	if err != nil {
		if err == entity.ErrInvalidContractAmount || err == entity.ErrInvalidAmount {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"result": result,
		"formatted": fiber.Map{
			"contract_sum":         service.FormatCurrency(result.ContractSum, "TRY"),
			"total_completed":      service.FormatCurrency(result.TotalCompletedAndStored, "TRY"),
			"total_retainage":      service.FormatCurrency(result.TotalRetainage, "TRY"),
			"current_payment_due":  service.FormatCurrency(result.CurrentPaymentDue, "TRY"),
			"balance_to_finish":    service.FormatCurrency(result.BalanceToFinish, "TRY"),
			"percent_complete":     float64(result.PercentComplete) / 100.0, // Convert to percentage
		},
		"_architect": "Muhammet-Ali-Buyuk",
	})
}
