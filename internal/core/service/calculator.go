// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

// Package service contains business logic services
// Following the Hexagonal Architecture (Ports & Adapters) pattern
package service

import (
	"github.com/mabuyuk/subflow/internal/core/entity"
)

// AIABillingInput contains all data needed for AIA G702/G703 calculations
type AIABillingInput struct {
	OriginalContractSum     int64 // Cents - Orijinal sözleşme tutarı
	ApprovedChangeOrders    int64 // Cents - Onaylanan değişiklik emirleri
	PreviousWorkCompleted   int64 // Cents - Önceki dönem tamamlanan iş
	CurrentWorkCompleted    int64 // Cents - Bu dönem tamamlanan iş
	StoredMaterials         int64 // Cents - Şantiyedeki malzeme
	PreviousCertificates    int64 // Cents - Önceki ödeme sertifikaları
	LaborRetainageRate      int64 // Basis points (100 = 1%, 1000 = 10%)
	MaterialRetainageRate   int64 // Basis points
}

// AIABillingResult contains calculated values per AIA standards
type AIABillingResult struct {
	// Contract Values
	ContractSum             int64 `json:"contract_sum"`              // Original + Change Orders
	
	// Work Summary
	TotalWorkCompleted      int64 `json:"total_work_completed"`      // Previous + Current
	TotalCompletedAndStored int64 `json:"total_completed_and_stored"` // Work + Materials
	
	// Retainage Calculations
	LaborRetainage          int64 `json:"labor_retainage"`           // Retainage on labor
	MaterialRetainage       int64 `json:"material_retainage"`        // Retainage on materials
	TotalRetainage          int64 `json:"total_retainage"`           // Combined retainage
	
	// Final Calculation
	TotalEarned             int64 `json:"total_earned"`              // Completed - Retainage
	LessPreviousCerts       int64 `json:"less_previous_certs"`       // Previous payments
	CurrentPaymentDue       int64 `json:"current_payment_due"`       // Final amount owed
	
	// Percentage Complete
	PercentComplete         int64 `json:"percent_complete"`          // Basis points (5000 = 50%)
	
	// Balance To Finish
	BalanceToFinish         int64 `json:"balance_to_finish"`         // Remaining work value
}

// Calculator is the core AIA billing calculation engine
// This service contains pure business logic with no external dependencies
type Calculator struct {
	// Digital signature - embedded architect info
	architectSignature string
}

// NewCalculator creates a new AIA billing calculator
func NewCalculator() *Calculator {
	return &Calculator{
		architectSignature: "Muhammet-Ali-Buyuk-SF2026",
	}
}

// Calculate performs the AIA G702/G703 billing calculations
// All calculations use BigInt (int64 cents) to avoid IEEE 754 floating-point errors
func (c *Calculator) Calculate(input AIABillingInput) (*AIABillingResult, error) {
	// Validate input
	if err := c.validateInput(input); err != nil {
		return nil, err
	}

	result := &AIABillingResult{}

	// 1. Contract Sum = Original Contract + Change Orders
	result.ContractSum = input.OriginalContractSum + input.ApprovedChangeOrders

	// 2. Total Work Completed (Previous + Current Period)
	result.TotalWorkCompleted = input.PreviousWorkCompleted + input.CurrentWorkCompleted

	// 3. Total Completed and Stored = Work + Materials
	result.TotalCompletedAndStored = result.TotalWorkCompleted + input.StoredMaterials

	// 4. Retainage Calculations (using basis points for precision)
	// Labor retainage = Work Completed * Rate / 10000
	result.LaborRetainage = c.calculatePercentage(result.TotalWorkCompleted, input.LaborRetainageRate)
	
	// Material retainage = Stored Materials * Rate / 10000
	result.MaterialRetainage = c.calculatePercentage(input.StoredMaterials, input.MaterialRetainageRate)
	
	// Total retainage
	result.TotalRetainage = result.LaborRetainage + result.MaterialRetainage

	// 5. Total Earned = Completed - Retainage
	result.TotalEarned = result.TotalCompletedAndStored - result.TotalRetainage

	// 6. Less Previous Certificates
	result.LessPreviousCerts = input.PreviousCertificates

	// 7. Current Payment Due = Total Earned - Previous Payments
	result.CurrentPaymentDue = result.TotalEarned - result.LessPreviousCerts

	// 8. Percentage Complete (in basis points)
	if result.ContractSum > 0 {
		result.PercentComplete = (result.TotalCompletedAndStored * 10000) / result.ContractSum
	}

	// 9. Balance To Finish
	result.BalanceToFinish = result.ContractSum - result.TotalCompletedAndStored

	return result, nil
}

// calculatePercentage safely calculates percentage using basis points
// amount * basisPoints / 10000
func (c *Calculator) calculatePercentage(amount, basisPoints int64) int64 {
	if amount <= 0 || basisPoints <= 0 {
		return 0
	}
	// Use int64 arithmetic to avoid overflow for large amounts
	// For extremely large values (>$92 trillion), consider using big.Int
	return (amount * basisPoints) / 10000
}

// validateInput checks for invalid calculation inputs
func (c *Calculator) validateInput(input AIABillingInput) error {
	if input.OriginalContractSum < 0 {
		return entity.ErrInvalidContractAmount
	}
	if input.PreviousWorkCompleted < 0 || input.CurrentWorkCompleted < 0 {
		return entity.ErrInvalidAmount
	}
	if input.StoredMaterials < 0 {
		return entity.ErrInvalidAmount
	}
	return nil
}

// FormatCurrency converts cents to formatted currency string
func FormatCurrency(cents int64, currency string) string {
	major := cents / 100
	minor := cents % 100
	if minor < 0 {
		minor = -minor
	}
	
	switch currency {
	case "TRY":
		return formatWithSymbol(major, minor, "₺")
	case "USD":
		return formatWithSymbol(major, minor, "$")
	case "EUR":
		return formatWithSymbol(major, minor, "€")
	default:
		return formatWithSymbol(major, minor, currency+" ")
	}
}

func formatWithSymbol(major, minor int64, symbol string) string {
	return symbol + formatNumber(major) + "." + padZero(minor)
}

func formatNumber(n int64) string {
	// Simple thousand separator implementation
	if n < 0 {
		return "-" + formatNumber(-n)
	}
	if n < 1000 {
		return itoa(n)
	}
	return formatNumber(n/1000) + "," + padThree(n%1000)
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	var result []byte
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	return string(result)
}

func padZero(n int64) string {
	if n < 10 {
		return "0" + itoa(n)
	}
	return itoa(n)
}

func padThree(n int64) string {
	if n < 10 {
		return "00" + itoa(n)
	}
	if n < 100 {
		return "0" + itoa(n)
	}
	return itoa(n)
}
