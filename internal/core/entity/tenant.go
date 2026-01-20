// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package entity

import (
	"time"

	"github.com/google/uuid"
)

// TenantPlan represents the subscription tier
type TenantPlan string

const (
	TenantPlanFree       TenantPlan = "FREE"
	TenantPlanPro        TenantPlan = "PRO"
	TenantPlanEnterprise TenantPlan = "ENTERPRISE"
)

// Tenant represents a customer organization in the multi-tenant system
type Tenant struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Slug          string     `json:"slug"` // URL-friendly identifier
	Plan          TenantPlan `json:"plan"`
	IsActive      bool       `json:"is_active"`
	MaxUsers      int        `json:"max_users"`
	MaxProjects   int        `json:"max_projects"`
	DefaultCurrency string   `json:"default_currency"`
	
	// Contact info
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone,omitempty"`
	Address      string `json:"address,omitempty"`
	
	// Metadata
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// NewTenant creates a new tenant with default values
func NewTenant(name, slug, email string) *Tenant {
	now := time.Now()
	return &Tenant{
		ID:              uuid.New(),
		Name:            name,
		Slug:            slug,
		Plan:            TenantPlanFree,
		IsActive:        true,
		MaxUsers:        5,     // Free tier limit
		MaxProjects:     3,     // Free tier limit
		DefaultCurrency: "TRY",
		ContactEmail:    email,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// UpgradePlan upgrades tenant to new plan with updated limits
func (t *Tenant) UpgradePlan(plan TenantPlan) {
	t.Plan = plan
	t.UpdatedAt = time.Now()
	
	switch plan {
	case TenantPlanPro:
		t.MaxUsers = 25
		t.MaxProjects = 50
	case TenantPlanEnterprise:
		t.MaxUsers = -1  // Unlimited
		t.MaxProjects = -1
	}
}

// CanAddUser checks if tenant can add more users
func (t *Tenant) CanAddUser(currentCount int) bool {
	if t.MaxUsers < 0 {
		return true // Unlimited
	}
	return currentCount < t.MaxUsers
}

// CanAddProject checks if tenant can add more projects
func (t *Tenant) CanAddProject(currentCount int) bool {
	if t.MaxProjects < 0 {
		return true // Unlimited
	}
	return currentCount < t.MaxProjects
}
