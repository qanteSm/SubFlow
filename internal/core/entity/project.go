// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package entity

import (
	"time"

	"github.com/google/uuid"
)

// ProjectStatus represents the current state of a project
type ProjectStatus string

const (
	ProjectStatusDraft      ProjectStatus = "DRAFT"
	ProjectStatusActive     ProjectStatus = "ACTIVE"
	ProjectStatusOnHold     ProjectStatus = "ON_HOLD"
	ProjectStatusCompleted  ProjectStatus = "COMPLETED"
	ProjectStatusCancelled  ProjectStatus = "CANCELLED"
)

// Project represents a construction project in the system
// This is a domain entity, independent of any database implementation
type Project struct {
	ID          uuid.UUID     `json:"id"`
	TenantID    uuid.UUID     `json:"tenant_id"`
	Name        string        `json:"name"`
	Code        string        `json:"code"` // Unique project code (e.g., "PRJ-2026-001")
	Description string        `json:"description"`
	Status      ProjectStatus `json:"status"`

	// Contract details
	ContractAmount   int64     `json:"contract_amount"`   // Amount in cents (BigInt)
	Currency         string    `json:"currency"`          // ISO 4217 (TRY, USD, EUR)
	StartDate        time.Time `json:"start_date"`
	EstimatedEndDate time.Time `json:"estimated_end_date"`

	// Retainage settings
	LaborRetainageRate    float64 `json:"labor_retainage_rate"`    // e.g., 0.10 for 10%
	MaterialRetainageRate float64 `json:"material_retainage_rate"` // e.g., 0.05 for 5%

	// Metadata
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// NewProject creates a new project with default values
func NewProject(tenantID uuid.UUID, name, code string) *Project {
	now := time.Now()
	return &Project{
		ID:                    uuid.New(),
		TenantID:              tenantID,
		Name:                  name,
		Code:                  code,
		Status:                ProjectStatusDraft,
		Currency:              "TRY",
		LaborRetainageRate:    0.10, // Default 10%
		MaterialRetainageRate: 0.05, // Default 5%
		CreatedAt:             now,
		UpdatedAt:             now,
	}
}

// Validate checks if the project has valid data
func (p *Project) Validate() error {
	if p.Name == "" {
		return ErrProjectNameRequired
	}
	if p.Code == "" {
		return ErrProjectCodeRequired
	}
	if p.ContractAmount < 0 {
		return ErrInvalidContractAmount
	}
	return nil
}

// IsActive returns true if the project is in active status
func (p *Project) IsActive() bool {
	return p.Status == ProjectStatusActive
}

// CanBeModified returns true if financial data can be modified
func (p *Project) CanBeModified() bool {
	return p.Status == ProjectStatusDraft || p.Status == ProjectStatusActive
}
