// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package entity

import (
	"regexp"
	"time"

	"github.com/google/uuid"
)

// UserRole represents the user's role in the system
type UserRole string

const (
	UserRoleAdmin      UserRole = "ADMIN"
	UserRoleManager    UserRole = "MANAGER"
	UserRoleAccountant UserRole = "ACCOUNTANT"
	UserRoleViewer     UserRole = "VIEWER"
)

// User represents a system user
type User struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // Never expose in JSON
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Role         UserRole   `json:"role"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// emailRegex is a simple email validation pattern
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewUser creates a new user with default values
func NewUser(tenantID uuid.UUID, email, firstName, lastName string, role UserRole) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate checks user data integrity
func (u *User) Validate() error {
	if !emailRegex.MatchString(u.Email) {
		return ErrInvalidEmail
	}
	return nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// CanManageProjects checks if user has permission to manage projects
func (u *User) CanManageProjects() bool {
	return u.Role == UserRoleAdmin || u.Role == UserRoleManager
}

// CanViewFinancials checks if user can view financial data
func (u *User) CanViewFinancials() bool {
	return u.Role == UserRoleAdmin || u.Role == UserRoleManager || u.Role == UserRoleAccountant
}

// CanApprovePayments checks if user can approve payments
func (u *User) CanApprovePayments() bool {
	return u.Role == UserRoleAdmin || u.Role == UserRoleManager
}
