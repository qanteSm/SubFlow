// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package entity

import "errors"

// Domain errors - these are business logic errors, not infrastructure errors
var (
	// Project errors
	ErrProjectNotFound       = errors.New("project not found")
	ErrProjectNameRequired   = errors.New("project name is required")
	ErrProjectCodeRequired   = errors.New("project code is required")
	ErrInvalidContractAmount = errors.New("contract amount cannot be negative")
	ErrProjectNotModifiable  = errors.New("project cannot be modified in current status")

	// Transaction errors
	ErrTransactionNotFound    = errors.New("transaction not found")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
	ErrInvalidAmount          = errors.New("amount must be greater than zero")
	ErrCurrencyMismatch       = errors.New("currency mismatch in transaction")

	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrUnauthorized      = errors.New("unauthorized access")

	// Tenant errors
	ErrTenantNotFound = errors.New("tenant not found")
	ErrTenantInactive = errors.New("tenant is inactive")

	// Contract errors
	ErrContractNotFound     = errors.New("contract not found")
	ErrContractAlreadyExists = errors.New("contract already exists for this vendor")

	// Calculation errors
	ErrRetainageExceedsTotal = errors.New("retainage cannot exceed total amount")
	ErrNegativeBalance       = errors.New("operation would result in negative balance")
)
