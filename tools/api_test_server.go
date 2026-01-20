// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// Standalone API Test Server - No External Dependencies
// Run: go run tools/api_test_server.go
// Test: curl http://localhost:3000/api/v1/calculate/aia -X POST -d "{...}"

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// AIABillingInput - AIA G702 inputs
type AIABillingInput struct {
	OriginalContractSum     int64 `json:"original_contract_sum"`
	NetChangeByChangeOrders int64 `json:"net_change_by_change_orders"`
	WorkCompletedPrevious   int64 `json:"work_completed_previous"`
	WorkCompletedThisPeriod int64 `json:"work_completed_this_period"`
	MaterialsStored         int64 `json:"materials_stored"`
	PreviousCertificates    int64 `json:"previous_certificates"`
	RetainagePercent        int64 `json:"retainage_percent"` // basis points (1000 = 10%)
}

// AIABillingResult - AIA G702 outputs
type AIABillingResult struct {
	ContractSumToDate        int64   `json:"contract_sum_to_date"`
	TotalCompletedAndStored  int64   `json:"total_completed_and_stored"`
	Retainage                int64   `json:"retainage"`
	TotalEarnedLessRetainage int64   `json:"total_earned_less_retainage"`
	CurrentPaymentDue        int64   `json:"current_payment_due"`
	BalanceToFinish          int64   `json:"balance_to_finish"`
	PercentComplete          float64 `json:"percent_complete"`

	// Formatted values for display
	Formatted FormattedResult `json:"formatted"`
}

type FormattedResult struct {
	ContractSumToDate        string `json:"contract_sum_to_date"`
	TotalCompletedAndStored  string `json:"total_completed_and_stored"`
	Retainage                string `json:"retainage"`
	TotalEarnedLessRetainage string `json:"total_earned_less_retainage"`
	CurrentPaymentDue        string `json:"current_payment_due"`
	BalanceToFinish          string `json:"balance_to_finish"`
	PercentComplete          string `json:"percent_complete"`
}

func formatCurrency(cents int64) string {
	isNegative := cents < 0
	if isNegative {
		cents = -cents
	}
	major := cents / 100
	minor := cents % 100
	majorStr := fmt.Sprintf("%d", major)
	var formatted strings.Builder
	for i, c := range majorStr {
		if i > 0 && (len(majorStr)-i)%3 == 0 {
			formatted.WriteRune('.')
		}
		formatted.WriteRune(c)
	}
	if isNegative {
		return fmt.Sprintf("-₺%s,%02d", formatted.String(), minor)
	}
	return fmt.Sprintf("₺%s,%02d", formatted.String(), minor)
}

func calculate(input AIABillingInput) AIABillingResult {
	result := AIABillingResult{}

	// Line 3: Contract Sum to Date
	result.ContractSumToDate = input.OriginalContractSum + input.NetChangeByChangeOrders

	// Line 4: Total Completed & Stored
	result.TotalCompletedAndStored = input.WorkCompletedPrevious +
		input.WorkCompletedThisPeriod + input.MaterialsStored

	// Line 4a: Retainage
	result.Retainage = (result.TotalCompletedAndStored * input.RetainagePercent) / 10000

	// Line 5: Total Earned Less Retainage
	result.TotalEarnedLessRetainage = result.TotalCompletedAndStored - result.Retainage

	// Line 7: Current Payment Due
	result.CurrentPaymentDue = result.TotalEarnedLessRetainage - input.PreviousCertificates

	// Line 8: Balance to Finish
	result.BalanceToFinish = result.ContractSumToDate - result.TotalCompletedAndStored + result.Retainage

	// Percent Complete
	if result.ContractSumToDate > 0 {
		result.PercentComplete = float64(result.TotalCompletedAndStored*100) / float64(result.ContractSumToDate)
	}

	// Formatted values
	result.Formatted = FormattedResult{
		ContractSumToDate:        formatCurrency(result.ContractSumToDate),
		TotalCompletedAndStored:  formatCurrency(result.TotalCompletedAndStored),
		Retainage:                formatCurrency(result.Retainage),
		TotalEarnedLessRetainage: formatCurrency(result.TotalEarnedLessRetainage),
		CurrentPaymentDue:        formatCurrency(result.CurrentPaymentDue),
		BalanceToFinish:          formatCurrency(result.BalanceToFinish),
		PercentComplete:          fmt.Sprintf("%.2f%%", result.PercentComplete),
	}

	return result
}

// API Response wrapper
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Architect string      `json:"_architect"`
	Timestamp string      `json:"timestamp"`
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Powered-By", "SubFlow-Enterprise")
		w.Header().Set("X-Architect", "Muhammet-Ali-Buyuk")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func respond(w http.ResponseWriter, success bool, data interface{}, err string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success:   success,
		Data:      data,
		Error:     err,
		Architect: "Muhammet-Ali-Buyuk",
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, true, map[string]string{
		"status":  "healthy",
		"service": "SubFlow API",
		"version": "1.0.0",
	}, "", http.StatusOK)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, true, map[string]string{
		"application": "SubFlow",
		"version":     "1.0.0",
		"architect":   "Muhammet-Ali-Buyuk",
		"build_id":    "SF-2026-ENTERPRISE",
		"compliance":  "AIA G702/G703",
	}, "", http.StatusOK)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respond(w, false, nil, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	var input AIABillingInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respond(w, false, nil, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input
	if input.OriginalContractSum <= 0 {
		respond(w, false, nil, "original_contract_sum must be positive", http.StatusBadRequest)
		return
	}

	result := calculate(input)
	respond(w, true, result, "", http.StatusOK)
}

func demoHandler(w http.ResponseWriter, r *http.Request) {
	// Demo calculation with sample data
	input := AIABillingInput{
		OriginalContractSum:     100000000, // ₺1,000,000.00
		NetChangeByChangeOrders: 5000000,   // ₺50,000.00
		WorkCompletedPrevious:   30000000,  // ₺300,000.00
		WorkCompletedThisPeriod: 15000000,  // ₺150,000.00
		MaterialsStored:         5000000,   // ₺50,000.00
		PreviousCertificates:    25000000,  // ₺250,000.00
		RetainagePercent:        1000,      // 10%
	}

	result := calculate(input)
	respond(w, true, map[string]interface{}{
		"input":  input,
		"output": result,
		"note":   "This is a demo calculation. POST to /api/v1/calculate/aia with your own data.",
	}, "", http.StatusOK)
}

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║     SubFlow API Test Server                                   ║")
	fmt.Println("║     Mimar: Muhammet Ali Büyük | alibuyuk.net                  ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Routes
	http.HandleFunc("/health", corsMiddleware(healthHandler))
	http.HandleFunc("/api/v1/system/version", corsMiddleware(versionHandler))
	http.HandleFunc("/api/v1/calculate/aia", corsMiddleware(calculateHandler))
	http.HandleFunc("/api/v1/calculate/demo", corsMiddleware(demoHandler))

	port := ":3000"
	fmt.Printf("🚀 Server starting on http://localhost%s\n", port)
	fmt.Println()
	fmt.Println("📡 Available Endpoints:")
	fmt.Println("   GET  /health                    - Health check")
	fmt.Println("   GET  /api/v1/system/version     - Version info")
	fmt.Println("   GET  /api/v1/calculate/demo     - Demo calculation")
	fmt.Println("   POST /api/v1/calculate/aia      - AIA G702 calculation")
	fmt.Println()
	fmt.Println("📝 Example POST:")
	fmt.Println(`   curl -X POST http://localhost:3000/api/v1/calculate/aia \`)
	fmt.Println(`        -H "Content-Type: application/json" \`)
	fmt.Println(`        -d '{"original_contract_sum":100000000,"retainage_percent":1000}'`)
	fmt.Println()

	log.Fatal(http.ListenAndServe(port, nil))
}
