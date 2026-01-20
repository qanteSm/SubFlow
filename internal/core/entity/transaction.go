// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// TransactionType represents the type of financial transaction
// Immutable ledger entries - NEVER updated, only appended
type TransactionType string

const (
	TransactionTypeInvoice          TransactionType = "INVOICE"           // Fatura kesildi
	TransactionTypePayment          TransactionType = "PAYMENT"           // Ödeme alındı
	TransactionTypeRetainageHeld    TransactionType = "RETAINAGE_HELD"    // Teminat tutuldu
	TransactionTypeRetainageRelease TransactionType = "RETAINAGE_RELEASE" // Teminat serbest bırakıldı
	TransactionTypeAdjustment       TransactionType = "ADJUSTMENT"        // Düzeltme
	TransactionTypeDeduction        TransactionType = "DEDUCTION"         // Kesinti
)

// Transaction represents an immutable financial event in the ledger
// This follows double-entry bookkeeping principles
type Transaction struct {
	ID            uuid.UUID       `json:"id"`
	ProjectID     uuid.UUID       `json:"project_id"`
	ContractID    *uuid.UUID      `json:"contract_id,omitempty"` // Optional: for subcontractor payments
	Type          TransactionType `json:"type"`
	AmountCents   int64           `json:"amount_cents"` // Amount in cents (BigInt arithmetic)
	Currency      string          `json:"currency"`     // ISO 4217
	EffectiveDate time.Time       `json:"effective_date"`
	Description   string          `json:"description"`
	ReferenceNo   string          `json:"reference_no"` // Bank receipt, invoice number, etc.
	Metadata      json.RawMessage `json:"metadata,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	CreatedBy     uuid.UUID       `json:"created_by"` // Actor who created this transaction
}

// TransactionMetadata contains additional context for transactions
type TransactionMetadata struct {
	BankReceiptNo     string `json:"bank_receipt_no,omitempty"`
	InvoiceNo         string `json:"invoice_no,omitempty"`
	ApplicationPeriod string `json:"application_period,omitempty"` // e.g., "2026-01"
	Notes             string `json:"notes,omitempty"`
	VendorName        string `json:"vendor_name,omitempty"`
	RetainageRate     string `json:"retainage_rate,omitempty"`
}

// NewTransaction creates a new transaction for the ledger
func NewTransaction(projectID uuid.UUID, txType TransactionType, amountCents int64, currency string, createdBy uuid.UUID) *Transaction {
	now := time.Now()
	return &Transaction{
		ID:            uuid.New(),
		ProjectID:     projectID,
		Type:          txType,
		AmountCents:   amountCents,
		Currency:      currency,
		EffectiveDate: now,
		CreatedAt:     now,
		CreatedBy:     createdBy,
	}
}

// Validate checks transaction data integrity
func (t *Transaction) Validate() error {
	if t.AmountCents <= 0 {
		return ErrInvalidAmount
	}
	if !t.Type.IsValid() {
		return ErrInvalidTransactionType
	}
	return nil
}

// IsValid checks if the transaction type is valid
func (tt TransactionType) IsValid() bool {
	switch tt {
	case TransactionTypeInvoice,
		TransactionTypePayment,
		TransactionTypeRetainageHeld,
		TransactionTypeRetainageRelease,
		TransactionTypeAdjustment,
		TransactionTypeDeduction:
		return true
	}
	return false
}

// IsCredit returns true if transaction represents money coming in
func (t *Transaction) IsCredit() bool {
	return t.Type == TransactionTypePayment || t.Type == TransactionTypeRetainageRelease
}

// IsDebit returns true if transaction represents money going out
func (t *Transaction) IsDebit() bool {
	return t.Type == TransactionTypeInvoice || t.Type == TransactionTypeDeduction
}

// SetMetadata convenience method to set metadata from struct
func (t *Transaction) SetMetadata(meta TransactionMetadata) error {
	data, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	t.Metadata = data
	return nil
}

// GetMetadata parses metadata into struct
func (t *Transaction) GetMetadata() (*TransactionMetadata, error) {
	if t.Metadata == nil {
		return &TransactionMetadata{}, nil
	}
	var meta TransactionMetadata
	if err := json.Unmarshal(t.Metadata, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}
