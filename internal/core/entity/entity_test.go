// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewProject(t *testing.T) {
	tenantID := uuid.New()
	project := NewProject(tenantID, "Test Project", "PRJ-001")

	if project.Name != "Test Project" {
		t.Errorf("Name = %s, want Test Project", project.Name)
	}

	if project.Code != "PRJ-001" {
		t.Errorf("Code = %s, want PRJ-001", project.Code)
	}

	if project.Status != ProjectStatusDraft {
		t.Errorf("Status = %s, want DRAFT", project.Status)
	}

	if project.LaborRetainageRate != 0.10 {
		t.Errorf("LaborRetainageRate = %f, want 0.10", project.LaborRetainageRate)
	}

	if project.ID == uuid.Nil {
		t.Error("Project ID should not be nil")
	}
}

func TestProject_Validate(t *testing.T) {
	tenantID := uuid.New()

	tests := []struct {
		name    string
		project *Project
		wantErr error
	}{
		{
			name:    "valid project",
			project: NewProject(tenantID, "Test", "PRJ-001"),
			wantErr: nil,
		},
		{
			name: "missing name",
			project: &Project{
				ID:       uuid.New(),
				TenantID: tenantID,
				Name:     "",
				Code:     "PRJ-001",
			},
			wantErr: ErrProjectNameRequired,
		},
		{
			name: "missing code",
			project: &Project{
				ID:       uuid.New(),
				TenantID: tenantID,
				Name:     "Test",
				Code:     "",
			},
			wantErr: ErrProjectCodeRequired,
		},
		{
			name: "negative contract amount",
			project: &Project{
				ID:             uuid.New(),
				TenantID:       tenantID,
				Name:           "Test",
				Code:           "PRJ-001",
				ContractAmount: -1000,
			},
			wantErr: ErrInvalidContractAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.project.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProject_IsActive(t *testing.T) {
	project := &Project{Status: ProjectStatusActive}
	if !project.IsActive() {
		t.Error("IsActive() should return true for ACTIVE status")
	}

	project.Status = ProjectStatusDraft
	if project.IsActive() {
		t.Error("IsActive() should return false for DRAFT status")
	}
}

func TestProject_CanBeModified(t *testing.T) {
	tests := []struct {
		status ProjectStatus
		want   bool
	}{
		{ProjectStatusDraft, true},
		{ProjectStatusActive, true},
		{ProjectStatusOnHold, false},
		{ProjectStatusCompleted, false},
		{ProjectStatusCancelled, false},
	}

	for _, tt := range tests {
		project := &Project{Status: tt.status}
		if got := project.CanBeModified(); got != tt.want {
			t.Errorf("CanBeModified() for %s = %v, want %v", tt.status, got, tt.want)
		}
	}
}

func TestNewTransaction(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()

	tx := NewTransaction(projectID, TransactionTypeInvoice, 100000, "TRY", userID)

	if tx.ProjectID != projectID {
		t.Error("ProjectID mismatch")
	}

	if tx.Type != TransactionTypeInvoice {
		t.Errorf("Type = %s, want INVOICE", tx.Type)
	}

	if tx.AmountCents != 100000 {
		t.Errorf("AmountCents = %d, want 100000", tx.AmountCents)
	}

	if tx.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestTransaction_Validate(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name    string
		tx      *Transaction
		wantErr error
	}{
		{
			name:    "valid transaction",
			tx:      NewTransaction(projectID, TransactionTypeInvoice, 100000, "TRY", userID),
			wantErr: nil,
		},
		{
			name: "zero amount",
			tx: &Transaction{
				ID:          uuid.New(),
				ProjectID:   projectID,
				Type:        TransactionTypeInvoice,
				AmountCents: 0,
				CreatedAt:   time.Now(),
			},
			wantErr: ErrInvalidAmount,
		},
		{
			name: "negative amount",
			tx: &Transaction{
				ID:          uuid.New(),
				ProjectID:   projectID,
				Type:        TransactionTypeInvoice,
				AmountCents: -100,
				CreatedAt:   time.Now(),
			},
			wantErr: ErrInvalidAmount,
		},
		{
			name: "invalid type",
			tx: &Transaction{
				ID:          uuid.New(),
				ProjectID:   projectID,
				Type:        "INVALID",
				AmountCents: 100000,
				CreatedAt:   time.Now(),
			},
			wantErr: ErrInvalidTransactionType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tx.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransactionType_IsValid(t *testing.T) {
	validTypes := []TransactionType{
		TransactionTypeInvoice,
		TransactionTypePayment,
		TransactionTypeRetainageHeld,
		TransactionTypeRetainageRelease,
		TransactionTypeAdjustment,
		TransactionTypeDeduction,
	}

	for _, tt := range validTypes {
		if !tt.IsValid() {
			t.Errorf("IsValid() for %s should return true", tt)
		}
	}

	invalidType := TransactionType("INVALID")
	if invalidType.IsValid() {
		t.Error("IsValid() for INVALID should return false")
	}
}

func TestTransaction_IsCredit(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()

	paymentTx := NewTransaction(projectID, TransactionTypePayment, 100000, "TRY", userID)
	if !paymentTx.IsCredit() {
		t.Error("Payment should be credit")
	}

	invoiceTx := NewTransaction(projectID, TransactionTypeInvoice, 100000, "TRY", userID)
	if invoiceTx.IsCredit() {
		t.Error("Invoice should not be credit")
	}
}

func TestUser_Validate(t *testing.T) {
	tenantID := uuid.New()

	validUser := NewUser(tenantID, "test@example.com", "John", "Doe", UserRoleViewer)
	if err := validUser.Validate(); err != nil {
		t.Errorf("Valid user should not return error: %v", err)
	}

	invalidUser := &User{Email: "invalid-email"}
	if err := invalidUser.Validate(); err != ErrInvalidEmail {
		t.Errorf("Invalid email should return ErrInvalidEmail, got: %v", err)
	}
}

func TestUser_Permissions(t *testing.T) {
	tenantID := uuid.New()

	admin := NewUser(tenantID, "admin@test.com", "Admin", "User", UserRoleAdmin)
	if !admin.CanManageProjects() {
		t.Error("Admin should be able to manage projects")
	}
	if !admin.CanViewFinancials() {
		t.Error("Admin should be able to view financials")
	}
	if !admin.CanApprovePayments() {
		t.Error("Admin should be able to approve payments")
	}

	viewer := NewUser(tenantID, "viewer@test.com", "Viewer", "User", UserRoleViewer)
	if viewer.CanManageProjects() {
		t.Error("Viewer should not be able to manage projects")
	}
	if viewer.CanViewFinancials() {
		t.Error("Viewer should not be able to view financials")
	}
}
