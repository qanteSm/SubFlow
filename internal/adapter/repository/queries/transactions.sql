-- name: CreateTransaction :one
-- Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
-- Financial ledger - IMMUTABLE entries only!
INSERT INTO transactions (
    project_id, contract_id, type, amount_cents, currency,
    effective_date, description, reference_no, metadata, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetTransactionByID :one
SELECT * FROM transactions
WHERE id = $1;

-- name: ListTransactionsByProject :many
SELECT * FROM transactions
WHERE project_id = $1
ORDER BY effective_date DESC, created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListTransactionsByProjectAndType :many
SELECT * FROM transactions
WHERE project_id = $1 AND type = $2
ORDER BY effective_date DESC;

-- name: GetProjectTransactionSummary :one
SELECT 
    project_id,
    COUNT(*) as transaction_count,
    COALESCE(SUM(CASE WHEN type = 'INVOICE' THEN amount_cents ELSE 0 END), 0) as total_invoiced,
    COALESCE(SUM(CASE WHEN type = 'PAYMENT' THEN amount_cents ELSE 0 END), 0) as total_paid,
    COALESCE(SUM(CASE WHEN type = 'RETAINAGE_HELD' THEN amount_cents ELSE 0 END), 0) as retainage_held,
    COALESCE(SUM(CASE WHEN type = 'RETAINAGE_RELEASE' THEN amount_cents ELSE 0 END), 0) as retainage_released,
    COALESCE(SUM(CASE WHEN type = 'DEDUCTION' THEN amount_cents ELSE 0 END), 0) as total_deductions
FROM transactions
WHERE project_id = $1
GROUP BY project_id;

-- name: GetProjectBalance :one
SELECT 
    COALESCE(SUM(
        CASE 
            WHEN type = 'INVOICE' THEN amount_cents
            WHEN type = 'PAYMENT' THEN -amount_cents
            WHEN type = 'RETAINAGE_RELEASE' THEN -amount_cents
            WHEN type = 'DEDUCTION' THEN -amount_cents
            ELSE 0
        END
    ), 0)::bigint as balance
FROM transactions
WHERE project_id = $1;

-- name: ListRecentTransactions :many
SELECT t.*, p.name as project_name, p.code as project_code
FROM transactions t
JOIN projects p ON t.project_id = p.id
WHERE p.tenant_id = $1
ORDER BY t.created_at DESC
LIMIT $2;

-- name: GetTransactionsByDateRange :many
SELECT * FROM transactions
WHERE project_id = $1 
    AND effective_date >= $2 
    AND effective_date <= $3
ORDER BY effective_date ASC;
