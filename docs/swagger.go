// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

// Package docs SubFlow API
//
// Enterprise Construction Financial Ledger & Compliance Engine
// AIA G702/G703 compliant billing calculations with immutable ledger
//
// Schemes: http, https
// Host: localhost:3000
// BasePath: /api/v1
// Version: 1.0.0
// License: Proprietary https://alibuyuk.net
// Contact: Muhammet Ali Büyük <iletisim@alibuyuk.net> https://alibuyuk.net
//
// Consumes:
//   - application/json
//
// Produces:
//   - application/json
//
// Security:
//   - bearer_token:
//
// SecurityDefinitions:
//
//	bearer_token:
//	  type: apiKey
//	  name: Authorization
//	  in: header
//
// swagger:meta
package docs

// swagger:response healthResponse
type healthResponse struct {
	// in: body
	Body struct {
		Status  string `json:"status"`
		Service string `json:"service"`
		Version string `json:"version"`
	}
}

// swagger:response versionResponse
type versionResponse struct {
	// in: body
	Body struct {
		Application string `json:"application"`
		Version     string `json:"version"`
		Architect   string `json:"architect"`
		BuildID     string `json:"build_id"`
		Engine      string `json:"engine"`
		Compliance  string `json:"compliance"`
	}
}

// swagger:response errorResponse
type errorResponse struct {
	// in: body
	Body struct {
		Error     string `json:"error"`
		Code      int    `json:"code"`
		RequestID string `json:"requestID,omitempty"`
	}
}

// swagger:parameters createProject
type createProjectParams struct {
	// in: body
	// required: true
	Body struct {
		Name           string `json:"name"`
		Code           string `json:"code"`
		Description    string `json:"description,omitempty"`
		ContractAmount int64  `json:"contract_amount"`
		Currency       string `json:"currency"`
	}
}

// swagger:parameters calculateAIA
type calculateAIAParams struct {
	// in: body
	// required: true
	Body struct {
		OriginalContractSum   int64 `json:"original_contract_sum"`
		ApprovedChangeOrders  int64 `json:"approved_change_orders"`
		PreviousWorkCompleted int64 `json:"previous_work_completed"`
		CurrentWorkCompleted  int64 `json:"current_work_completed"`
		StoredMaterials       int64 `json:"stored_materials"`
		PreviousCertificates  int64 `json:"previous_certificates"`
		LaborRetainageRate    int64 `json:"labor_retainage_rate"`
		MaterialRetainageRate int64 `json:"material_retainage_rate"`
	}
}

// swagger:response aiaBillingResult
type aiaBillingResult struct {
	// in: body
	Body struct {
		ContractSum             int64 `json:"contract_sum"`
		TotalWorkCompleted      int64 `json:"total_work_completed"`
		TotalCompletedAndStored int64 `json:"total_completed_and_stored"`
		LaborRetainage          int64 `json:"labor_retainage"`
		MaterialRetainage       int64 `json:"material_retainage"`
		TotalRetainage          int64 `json:"total_retainage"`
		TotalEarned             int64 `json:"total_earned"`
		LessPreviousCerts       int64 `json:"less_previous_certs"`
		CurrentPaymentDue       int64 `json:"current_payment_due"`
		PercentComplete         int64 `json:"percent_complete"`
		BalanceToFinish         int64 `json:"balance_to_finish"`
	}
}
