// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/mabuyuk/subflow/internal/core/entity"
	"github.com/mabuyuk/subflow/internal/core/service"
)

// InMemoryTransactionRepository is a simple in-memory implementation
// Used for testing and development before PostgreSQL is set up
type InMemoryTransactionRepository struct {
	mu           sync.RWMutex
	transactions map[uuid.UUID]*entity.Transaction
	architect    string
}

// NewInMemoryTransactionRepository creates a new in-memory repository
func NewInMemoryTransactionRepository() *InMemoryTransactionRepository {
	return &InMemoryTransactionRepository{
		transactions: make(map[uuid.UUID]*entity.Transaction),
		architect:    "Muhammet-Ali-Buyuk",
	}
}

// Save stores a transaction in memory
func (r *InMemoryTransactionRepository) Save(ctx context.Context, tx *entity.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.transactions[tx.ID] = tx
	return nil
}

// FindByID retrieves a transaction by its ID
func (r *InMemoryTransactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	tx, exists := r.transactions[id]
	if !exists {
		return nil, entity.ErrTransactionNotFound
	}
	return tx, nil
}

// FindByProjectID retrieves all transactions for a project
func (r *InMemoryTransactionRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entity.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []*entity.Transaction
	for _, tx := range r.transactions {
		if tx.ProjectID == projectID {
			result = append(result, tx)
		}
	}
	return result, nil
}

// GetProjectSummary calculates the financial summary for a project
func (r *InMemoryTransactionRepository) GetProjectSummary(ctx context.Context, projectID uuid.UUID) (*service.LedgerSummary, error) {
	transactions, err := r.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	summary := &service.LedgerSummary{
		ProjectID:        projectID,
		TransactionCount: len(transactions),
	}

	for _, tx := range transactions {
		summary.Currency = tx.Currency // Use the last found currency
		
		switch tx.Type {
		case entity.TransactionTypeInvoice:
			summary.TotalInvoiced += tx.AmountCents
		case entity.TransactionTypePayment:
			summary.TotalPaid += tx.AmountCents
		case entity.TransactionTypeRetainageHeld:
			summary.TotalRetained += tx.AmountCents
		case entity.TransactionTypeRetainageRelease:
			summary.TotalRetained -= tx.AmountCents
		}
	}

	summary.CurrentBalance = summary.TotalInvoiced - summary.TotalPaid

	return summary, nil
}

// Clear removes all transactions (for testing)
func (r *InMemoryTransactionRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.transactions = make(map[uuid.UUID]*entity.Transaction)
}
