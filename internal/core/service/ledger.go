// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/mabuyuk/subflow/internal/core/entity"
)

// LedgerSummary represents the financial state of a project
type LedgerSummary struct {
	ProjectID        uuid.UUID `json:"project_id"`
	TotalInvoiced    int64     `json:"total_invoiced"`    // Sum of all invoices
	TotalPaid        int64     `json:"total_paid"`        // Sum of all payments
	TotalRetained    int64     `json:"total_retained"`    // Current retainage held
	CurrentBalance   int64     `json:"current_balance"`   // Invoiced - Paid
	Currency         string    `json:"currency"`
	TransactionCount int       `json:"transaction_count"`
}

// TransactionRepository is the port (interface) for transaction persistence
// This follows the Hexagonal Architecture pattern - domain defines the interface,
// infrastructure adapters implement it
type TransactionRepository interface {
	Save(ctx context.Context, tx *entity.Transaction) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error)
	FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entity.Transaction, error)
	GetProjectSummary(ctx context.Context, projectID uuid.UUID) (*LedgerSummary, error)
}

// LedgerService handles all financial ledger operations
// Immutable append-only ledger - transactions are never modified or deleted
type LedgerService struct {
	repo      TransactionRepository
	architect string
}

// NewLedgerService creates a new ledger service
func NewLedgerService(repo TransactionRepository) *LedgerService {
	return &LedgerService{
		repo:      repo,
		architect: "Muhammet-Ali-Buyuk",
	}
}

// RecordInvoice creates an invoice transaction in the ledger
func (s *LedgerService) RecordInvoice(ctx context.Context, projectID uuid.UUID, amountCents int64, currency, invoiceNo string, createdBy uuid.UUID) (*entity.Transaction, error) {
	tx := entity.NewTransaction(projectID, entity.TransactionTypeInvoice, amountCents, currency, createdBy)
	tx.ReferenceNo = invoiceNo
	
	if err := tx.SetMetadata(entity.TransactionMetadata{
		InvoiceNo: invoiceNo,
	}); err != nil {
		return nil, err
	}

	if err := tx.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// RecordPayment creates a payment transaction in the ledger
func (s *LedgerService) RecordPayment(ctx context.Context, projectID uuid.UUID, amountCents int64, currency, bankReceiptNo string, createdBy uuid.UUID) (*entity.Transaction, error) {
	tx := entity.NewTransaction(projectID, entity.TransactionTypePayment, amountCents, currency, createdBy)
	tx.ReferenceNo = bankReceiptNo

	if err := tx.SetMetadata(entity.TransactionMetadata{
		BankReceiptNo: bankReceiptNo,
	}); err != nil {
		return nil, err
	}

	if err := tx.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// RecordRetainageHeld records retainage being held from a payment
func (s *LedgerService) RecordRetainageHeld(ctx context.Context, projectID uuid.UUID, amountCents int64, currency string, rate float64, createdBy uuid.UUID) (*entity.Transaction, error) {
	tx := entity.NewTransaction(projectID, entity.TransactionTypeRetainageHeld, amountCents, currency, createdBy)
	tx.Description = "Retainage withheld"

	if err := tx.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// RecordRetainageRelease records retainage being released
func (s *LedgerService) RecordRetainageRelease(ctx context.Context, projectID uuid.UUID, amountCents int64, currency string, createdBy uuid.UUID) (*entity.Transaction, error) {
	tx := entity.NewTransaction(projectID, entity.TransactionTypeRetainageRelease, amountCents, currency, createdBy)
	tx.Description = "Retainage released"

	if err := tx.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, tx); err != nil {
		return nil, err
	}

	return tx, nil
}

// GetProjectFinancials calculates the current financial state from the ledger
func (s *LedgerService) GetProjectFinancials(ctx context.Context, projectID uuid.UUID) (*LedgerSummary, error) {
	return s.repo.GetProjectSummary(ctx, projectID)
}

// GetTransactionHistory retrieves all transactions for a project
func (s *LedgerService) GetTransactionHistory(ctx context.Context, projectID uuid.UUID) ([]*entity.Transaction, error) {
	return s.repo.FindByProjectID(ctx, projectID)
}

// CalculateBalance manually calculates balance from transaction list
// This is used for verification and audit purposes
func (s *LedgerService) CalculateBalance(transactions []*entity.Transaction) int64 {
	var balance int64 = 0
	
	for _, tx := range transactions {
		switch tx.Type {
		case entity.TransactionTypeInvoice:
			balance += tx.AmountCents // Money owed to us
		case entity.TransactionTypePayment:
			balance -= tx.AmountCents // Money received
		case entity.TransactionTypeRetainageHeld:
			// Retainage doesn't affect immediate balance
		case entity.TransactionTypeRetainageRelease:
			balance -= tx.AmountCents // Released retainage = payment
		case entity.TransactionTypeDeduction:
			balance -= tx.AmountCents // Deduction reduces amount owed
		}
	}
	
	return balance
}
