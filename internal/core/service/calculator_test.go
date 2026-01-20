// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package service

import (
	"testing"
)

// TestCalculator_BasicCalculation tests the AIA billing calculation
func TestCalculator_BasicCalculation(t *testing.T) {
	calc := NewCalculator()

	input := AIABillingInput{
		OriginalContractSum:   100000000, // $1,000,000.00
		ApprovedChangeOrders:  5000000,   // $50,000.00
		PreviousWorkCompleted: 30000000,  // $300,000.00
		CurrentWorkCompleted:  15000000,  // $150,000.00
		StoredMaterials:       5000000,   // $50,000.00
		PreviousCertificates:  25000000,  // $250,000.00
		LaborRetainageRate:    1000,      // 10%
		MaterialRetainageRate: 500,       // 5%
	}

	result, err := calc.Calculate(input)
	if err != nil {
		t.Fatalf("Calculate returned error: %v", err)
	}

	// Verify contract sum = original + change orders
	expectedContractSum := int64(105000000) // $1,050,000.00
	if result.ContractSum != expectedContractSum {
		t.Errorf("ContractSum = %d, want %d", result.ContractSum, expectedContractSum)
	}

	// Verify total work completed
	expectedTotalWork := int64(45000000) // $450,000.00
	if result.TotalWorkCompleted != expectedTotalWork {
		t.Errorf("TotalWorkCompleted = %d, want %d", result.TotalWorkCompleted, expectedTotalWork)
	}

	// Verify total completed and stored
	expectedTotal := int64(50000000) // $500,000.00
	if result.TotalCompletedAndStored != expectedTotal {
		t.Errorf("TotalCompletedAndStored = %d, want %d", result.TotalCompletedAndStored, expectedTotal)
	}

	// Verify labor retainage = 45000000 * 0.10 = 4500000
	expectedLaborRetainage := int64(4500000)
	if result.LaborRetainage != expectedLaborRetainage {
		t.Errorf("LaborRetainage = %d, want %d", result.LaborRetainage, expectedLaborRetainage)
	}

	// Verify material retainage = 5000000 * 0.05 = 250000
	expectedMaterialRetainage := int64(250000)
	if result.MaterialRetainage != expectedMaterialRetainage {
		t.Errorf("MaterialRetainage = %d, want %d", result.MaterialRetainage, expectedMaterialRetainage)
	}

	// Verify total retainage
	expectedTotalRetainage := int64(4750000) // $47,500.00
	if result.TotalRetainage != expectedTotalRetainage {
		t.Errorf("TotalRetainage = %d, want %d", result.TotalRetainage, expectedTotalRetainage)
	}

	// Verify current payment due
	// Total Earned = 50000000 - 4750000 = 45250000
	// Current Due = 45250000 - 25000000 = 20250000
	expectedPaymentDue := int64(20250000) // $202,500.00
	if result.CurrentPaymentDue != expectedPaymentDue {
		t.Errorf("CurrentPaymentDue = %d, want %d", result.CurrentPaymentDue, expectedPaymentDue)
	}
}

// TestCalculator_ZeroValues tests with zero input values
func TestCalculator_ZeroValues(t *testing.T) {
	calc := NewCalculator()

	input := AIABillingInput{
		OriginalContractSum:   100000000,
		ApprovedChangeOrders:  0,
		PreviousWorkCompleted: 0,
		CurrentWorkCompleted:  0,
		StoredMaterials:       0,
		PreviousCertificates:  0,
		LaborRetainageRate:    1000,
		MaterialRetainageRate: 500,
	}

	result, err := calc.Calculate(input)
	if err != nil {
		t.Fatalf("Calculate returned error: %v", err)
	}

	if result.CurrentPaymentDue != 0 {
		t.Errorf("CurrentPaymentDue = %d, want 0", result.CurrentPaymentDue)
	}

	if result.PercentComplete != 0 {
		t.Errorf("PercentComplete = %d, want 0", result.PercentComplete)
	}
}

// TestCalculator_FullCompletion tests 100% project completion
func TestCalculator_FullCompletion(t *testing.T) {
	calc := NewCalculator()

	input := AIABillingInput{
		OriginalContractSum:   100000000,
		ApprovedChangeOrders:  0,
		PreviousWorkCompleted: 100000000, // 100% complete
		CurrentWorkCompleted:  0,
		StoredMaterials:       0,
		PreviousCertificates:  90000000, // 90% already paid
		LaborRetainageRate:    1000,
		MaterialRetainageRate: 500,
	}

	result, err := calc.Calculate(input)
	if err != nil {
		t.Fatalf("Calculate returned error: %v", err)
	}

	// 100% completion
	if result.PercentComplete != 10000 {
		t.Errorf("PercentComplete = %d, want 10000 (100%%)", result.PercentComplete)
	}

	// Balance to finish should be 0
	if result.BalanceToFinish != 0 {
		t.Errorf("BalanceToFinish = %d, want 0", result.BalanceToFinish)
	}
}

// TestCalculator_NegativeValues tests validation for invalid inputs
func TestCalculator_NegativeValues(t *testing.T) {
	calc := NewCalculator()

	input := AIABillingInput{
		OriginalContractSum:   -100000000, // Invalid
		ApprovedChangeOrders:  0,
		PreviousWorkCompleted: 0,
		CurrentWorkCompleted:  0,
		StoredMaterials:       0,
		PreviousCertificates:  0,
		LaborRetainageRate:    1000,
		MaterialRetainageRate: 500,
	}

	_, err := calc.Calculate(input)
	if err == nil {
		t.Error("Expected error for negative contract sum, got nil")
	}
}

// TestFormatCurrency tests currency formatting
func TestFormatCurrency(t *testing.T) {
	tests := []struct {
		cents    int64
		currency string
		want     string
	}{
		{123456, "TRY", "₺1,234.56"},
		{100, "USD", "$1.00"},
		{99, "EUR", "€0.99"},
		{1000000, "TRY", "₺10,000.00"},
		{0, "TRY", "₺0.00"},
	}

	for _, tt := range tests {
		got := FormatCurrency(tt.cents, tt.currency)
		if got != tt.want {
			t.Errorf("FormatCurrency(%d, %s) = %s, want %s", tt.cents, tt.currency, got, tt.want)
		}
	}
}

// BenchmarkCalculator measures calculation performance
func BenchmarkCalculator(b *testing.B) {
	calc := NewCalculator()
	input := AIABillingInput{
		OriginalContractSum:   100000000,
		ApprovedChangeOrders:  5000000,
		PreviousWorkCompleted: 30000000,
		CurrentWorkCompleted:  15000000,
		StoredMaterials:       5000000,
		PreviousCertificates:  25000000,
		LaborRetainageRate:    1000,
		MaterialRetainageRate: 500,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Calculate(input)
	}
}
