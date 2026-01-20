// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// AIA G702/G703 Report Generator - FIXED VERSION
// Run: go run generate_report.go

package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// AIABillingInput - AIA G702 Application for Payment inputs
type AIABillingInput struct {
	// Project Information (G702 Header)
	ProjectName    string
	ProjectCode    string
	ContractorName string
	OwnerName      string
	ArchitectName  string
	ContractDate   string
	ApplicationNo  int
	PeriodTo       string

	// G702 Line Items (values in cents to avoid float errors)
	OriginalContractSum     int64 // Line 1
	NetChangeByChangeOrders int64 // Line 2
	// ContractSumToDate = Line 1 + Line 2 (calculated)
	TotalCompletedAndStoredToDate int64 // Line 4 (from G703 Grand Total)
	RetainagePercent              int64 // basis points (1000 = 10%)
	// TotalEarnedLessRetainage = Line 4 - Retainage (calculated)
	LessPreviousCertificates int64 // Line 6
	// CurrentPaymentDue = Line 5 - Line 6 (calculated)
	// BalanceToFinish = Line 3 - Line 4 (calculated)

	// G703 Details
	WorkCompletedPreviousApp int64
	WorkCompletedThisPeriod  int64
	MaterialsStoredToDate    int64
}

// AIABillingResult - AIA G702 calculated results
type AIABillingResult struct {
	ContractSumToDate         int64 // Line 3
	TotalCompletedAndStored   int64 // Line 4
	Retainage                 int64 // Line 4a
	TotalEarnedLessRetainage  int64 // Line 5
	LessPreviousCertificates  int64 // Line 6
	CurrentPaymentDue         int64 // Line 7
	BalanceToFinishPlusRetain int64 // Line 8
	PercentComplete           int64 // basis points
}

func calculate(input AIABillingInput) AIABillingResult {
	result := AIABillingResult{}

	// Line 3: Contract Sum to Date
	result.ContractSumToDate = input.OriginalContractSum + input.NetChangeByChangeOrders

	// Line 4: Total Completed & Stored to Date
	result.TotalCompletedAndStored = input.WorkCompletedPreviousApp +
		input.WorkCompletedThisPeriod + input.MaterialsStoredToDate

	// Line 4a: Retainage (using basis points for precision)
	result.Retainage = (result.TotalCompletedAndStored * input.RetainagePercent) / 10000

	// Line 5: Total Earned Less Retainage
	result.TotalEarnedLessRetainage = result.TotalCompletedAndStored - result.Retainage

	// Line 6: Less Previous Certificates
	result.LessPreviousCertificates = input.LessPreviousCertificates

	// Line 7: Current Payment Due
	result.CurrentPaymentDue = result.TotalEarnedLessRetainage - result.LessPreviousCertificates

	// Line 8: Balance to Finish Including Retainage
	result.BalanceToFinishPlusRetain = result.ContractSumToDate - result.TotalCompletedAndStored + result.Retainage

	// Percent Complete
	if result.ContractSumToDate > 0 {
		result.PercentComplete = (result.TotalCompletedAndStored * 10000) / result.ContractSumToDate
	}

	return result
}

func formatCurrency(cents int64) string {
	isNegative := cents < 0
	if isNegative {
		cents = -cents
	}
	major := cents / 100
	minor := cents % 100

	// Add thousand separators
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

func formatPercent(basisPoints int64) string {
	whole := basisPoints / 100
	frac := basisPoints % 100
	return fmt.Sprintf("%%%d,%02d", whole, frac)
}

func generateHTMLReport(input AIABillingInput, result AIABillingResult) string {
	percentComplete := float64(result.PercentComplete) / 100.0

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="tr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SubFlow - AIA G702 Application for Payment</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, sans-serif;
            background: #f0f4f8;
            padding: 20px;
        }
        .container { 
            max-width: 800px; 
            margin: 0 auto; 
            background: white;
            border: 2px solid #1e3a8a;
        }
        .header { 
            background: #1e3a8a;
            color: white; 
            padding: 15px 20px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .header h1 { font-size: 18px; }
        .header .doc-id { font-size: 14px; opacity: 0.9; }
        .form-title {
            text-align: center;
            padding: 15px;
            border-bottom: 2px solid #1e3a8a;
            font-size: 16px;
            font-weight: bold;
            background: #dbeafe;
        }
        .info-section {
            display: grid;
            grid-template-columns: 1fr 1fr;
            border-bottom: 1px solid #ccc;
        }
        .info-box {
            padding: 10px 15px;
            border-right: 1px solid #ccc;
        }
        .info-box:last-child { border-right: none; }
        .info-box label { 
            font-size: 10px; 
            color: #666; 
            text-transform: uppercase;
            display: block;
            margin-bottom: 3px;
        }
        .info-box span { font-weight: 600; }
        .application-section {
            padding: 15px;
            border-bottom: 2px solid #1e3a8a;
        }
        .application-section h3 {
            font-size: 12px;
            text-transform: uppercase;
            color: #1e3a8a;
            margin-bottom: 15px;
            padding-bottom: 5px;
            border-bottom: 1px solid #ddd;
        }
        table { width: 100%%; border-collapse: collapse; }
        table tr { border-bottom: 1px solid #eee; }
        table td { padding: 8px 5px; font-size: 13px; }
        table td:first-child { width: 30px; color: #666; }
        table td:last-child { text-align: right; font-weight: 600; }
        .line-highlight {
            background: #fef3c7;
        }
        .payment-due {
            background: #10b981;
            color: white;
        }
        .payment-due td { 
            padding: 15px; 
            font-size: 16px; 
            font-weight: bold;
        }
        .negative { color: #dc2626; }
        .progress-section {
            padding: 15px;
            background: #f8fafc;
        }
        .progress-bar {
            height: 24px;
            background: #e2e8f0;
            border-radius: 4px;
            overflow: hidden;
            margin: 10px 0;
        }
        .progress-fill {
            height: 100%%;
            background: linear-gradient(90deg, #3b82f6, #1e3a8a);
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: bold;
            font-size: 12px;
        }
        .footer {
            padding: 15px;
            text-align: center;
            font-size: 11px;
            color: #666;
            border-top: 2px solid #1e3a8a;
        }
        .footer a { color: #1e3a8a; }
        .aia-badge {
            display: inline-block;
            padding: 3px 10px;
            background: #1e3a8a;
            color: white;
            border-radius: 3px;
            font-size: 10px;
            margin-bottom: 10px;
        }
        @media print {
            body { padding: 0; background: white; }
            .container { border: 1px solid black; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div>
                <h1>🏗️ SubFlow</h1>
                <small>Enterprise Construction Financial Ledger</small>
            </div>
            <div class="doc-id">
                Başvuru No: %d<br>
                Dönem: %s
            </div>
        </div>

        <div class="form-title">
            AIA DOCUMENT G702 - APPLICATION AND CERTIFICATE FOR PAYMENT
        </div>

        <div class="info-section">
            <div class="info-box">
                <label>Proje Adı</label>
                <span>%s</span>
            </div>
            <div class="info-box">
                <label>Proje No</label>
                <span>%s</span>
            </div>
        </div>
        <div class="info-section">
            <div class="info-box">
                <label>Yüklenici</label>
                <span>%s</span>
            </div>
            <div class="info-box">
                <label>Tarih</label>
                <span>%s</span>
            </div>
        </div>

        <div class="application-section">
            <h3>Contractor's Application for Payment</h3>
            <table>
                <tr>
                    <td>1.</td>
                    <td>Original Contract Sum</td>
                    <td>%s</td>
                </tr>
                <tr>
                    <td>2.</td>
                    <td>Net Change by Change Orders</td>
                    <td>%s</td>
                </tr>
                <tr class="line-highlight">
                    <td>3.</td>
                    <td><strong>Contract Sum to Date (Line 1 + 2)</strong></td>
                    <td><strong>%s</strong></td>
                </tr>
                <tr>
                    <td>4.</td>
                    <td>Total Completed & Stored to Date (from G703)</td>
                    <td>%s</td>
                </tr>
                <tr>
                    <td>4a.</td>
                    <td>Retainage (%s of Line 4)</td>
                    <td class="negative">-%s</td>
                </tr>
                <tr class="line-highlight">
                    <td>5.</td>
                    <td><strong>Total Earned Less Retainage (Line 4 - 4a)</strong></td>
                    <td><strong>%s</strong></td>
                </tr>
                <tr>
                    <td>6.</td>
                    <td>Less Previous Certificates for Payment</td>
                    <td class="negative">-%s</td>
                </tr>
            </table>
        </div>

        <table class="payment-due">
            <tr>
                <td>7.</td>
                <td>CURRENT PAYMENT DUE</td>
                <td>%s</td>
            </tr>
        </table>

        <div class="application-section">
            <table>
                <tr>
                    <td>8.</td>
                    <td>Balance to Finish, Including Retainage</td>
                    <td>%s</td>
                </tr>
            </table>
        </div>

        <div class="progress-section">
            <strong>Percentage of Completion:</strong>
            <div class="progress-bar">
                <div class="progress-fill" style="width: %.2f%%;">
                    %.2f%%
                </div>
            </div>
        </div>

        <div class="application-section">
            <h3>G703 Continuation Sheet Summary</h3>
            <table>
                <tr>
                    <td>A.</td>
                    <td>Work Completed - Previous Applications</td>
                    <td>%s</td>
                </tr>
                <tr>
                    <td>B.</td>
                    <td>Work Completed - This Period</td>
                    <td>%s</td>
                </tr>
                <tr>
                    <td>C.</td>
                    <td>Materials Presently Stored</td>
                    <td>%s</td>
                </tr>
                <tr class="line-highlight">
                    <td>D.</td>
                    <td><strong>Total Completed & Stored to Date (A+B+C)</strong></td>
                    <td><strong>%s</strong></td>
                </tr>
            </table>
        </div>

        <div class="footer">
            <span class="aia-badge">AIA G702/G703 UYUMLU</span><br>
            Bu form AIA (American Institute of Architects) G702 Application and Certificate for Payment 
            standardına uygun olarak hazırlanmıştır.<br><br>
            <strong>SubFlow Enterprise</strong> | Mimar: <a href="https://alibuyuk.net">Muhammet Ali Büyük</a><br>
            © 2026 Tüm hakları saklıdır.
        </div>
    </div>
</body>
</html>`,
		input.ApplicationNo,
		input.PeriodTo,
		input.ProjectName,
		input.ProjectCode,
		input.ContractorName,
		time.Now().Format("02.01.2006"),
		formatCurrency(input.OriginalContractSum),
		formatCurrency(input.NetChangeByChangeOrders),
		formatCurrency(result.ContractSumToDate),
		formatCurrency(result.TotalCompletedAndStored),
		formatPercent(input.RetainagePercent),
		formatCurrency(result.Retainage),
		formatCurrency(result.TotalEarnedLessRetainage),
		formatCurrency(result.LessPreviousCertificates),
		formatCurrency(result.CurrentPaymentDue),
		formatCurrency(result.BalanceToFinishPlusRetain),
		percentComplete,
		percentComplete,
		formatCurrency(input.WorkCompletedPreviousApp),
		formatCurrency(input.WorkCompletedThisPeriod),
		formatCurrency(input.MaterialsStoredToDate),
		formatCurrency(result.TotalCompletedAndStored),
	)
}

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║     SubFlow - AIA G702/G703 Report Generator                  ║")
	fmt.Println("║     Mimar: Muhammet Ali Büyük | alibuyuk.net                  ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Test data - Metro Station Project
	input := AIABillingInput{
		ProjectName:    "Metro İstasyonu Genişletme Projesi",
		ProjectCode:    "PRJ-2026-001",
		ContractorName: "ABC İnşaat A.Ş.",
		OwnerName:      "İstanbul Büyükşehir Belediyesi",
		ArchitectName:  "XYZ Mimarlık Ltd.",
		ApplicationNo:  3,
		PeriodTo:       "Ocak 2026",

		// G702 Values (in cents)
		OriginalContractSum:      100000000, // ₺1,000,000.00
		NetChangeByChangeOrders:  5000000,   // ₺50,000.00
		RetainagePercent:         1000,      // 10% (1000 basis points)
		LessPreviousCertificates: 25000000,  // ₺250,000.00

		// G703 Breakdown
		WorkCompletedPreviousApp: 30000000, // ₺300,000.00
		WorkCompletedThisPeriod:  15000000, // ₺150,000.00
		MaterialsStoredToDate:    5000000,  // ₺50,000.00
	}

	result := calculate(input)

	// Print calculation summary
	fmt.Println("📊 AIA G702 HESAPLAMA:")
	fmt.Println("═══════════════════════════════════════════")
	fmt.Printf("1. Original Contract Sum:       %s\n", formatCurrency(input.OriginalContractSum))
	fmt.Printf("2. Net Change by Change Orders: %s\n", formatCurrency(input.NetChangeByChangeOrders))
	fmt.Printf("3. Contract Sum to Date:        %s\n", formatCurrency(result.ContractSumToDate))
	fmt.Printf("4. Total Completed & Stored:    %s\n", formatCurrency(result.TotalCompletedAndStored))
	fmt.Printf("4a. Retainage (10%%):            -%s\n", formatCurrency(result.Retainage))
	fmt.Printf("5. Total Earned Less Retainage: %s\n", formatCurrency(result.TotalEarnedLessRetainage))
	fmt.Printf("6. Less Previous Certificates:  -%s\n", formatCurrency(result.LessPreviousCertificates))
	fmt.Println("═══════════════════════════════════════════")
	fmt.Printf("7. 💰 CURRENT PAYMENT DUE:      %s\n", formatCurrency(result.CurrentPaymentDue))
	fmt.Println("═══════════════════════════════════════════")
	fmt.Printf("8. Balance + Retainage:         %s\n", formatCurrency(result.BalanceToFinishPlusRetain))
	fmt.Printf("   Completion:                  %.2f%%\n", float64(result.PercentComplete)/100.0)
	fmt.Println()

	// Generate HTML report
	html := generateHTMLReport(input, result)

	filename := "report_aia_g702.html"
	err := os.WriteFile(filename, []byte(html), 0644)
	if err != nil {
		fmt.Printf("Hata: %v\n", err)
		return
	}

	fmt.Printf("✅ AIA G702 raporu oluşturuldu: %s\n", filename)
	fmt.Println("📂 Tarayıcıda açmak için: start " + filename)
}
