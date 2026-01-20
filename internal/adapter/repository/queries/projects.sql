-- name: CreateProject :one
-- Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
INSERT INTO projects (
    tenant_id, name, code, description, status,
    contract_amount_cents, currency, start_date, estimated_end_date,
    labor_retainage_rate, material_retainage_rate
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetProjectByCode :one
SELECT * FROM projects
WHERE tenant_id = $1 AND code = $2 AND deleted_at IS NULL;

-- name: ListProjectsByTenant :many
SELECT * FROM projects
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListActiveProjects :many
SELECT * FROM projects
WHERE tenant_id = $1 AND status = 'ACTIVE' AND deleted_at IS NULL
ORDER BY name ASC;

-- name: UpdateProject :one
UPDATE projects
SET 
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    status = COALESCE($4, status),
    contract_amount_cents = COALESCE($5, contract_amount_cents),
    start_date = COALESCE($6, start_date),
    estimated_end_date = COALESCE($7, estimated_end_date),
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateProjectStatus :one
UPDATE projects
SET status = $2, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteProject :exec
UPDATE projects
SET deleted_at = NOW()
WHERE id = $1;

-- name: GetProjectCount :one
SELECT COUNT(*) FROM projects
WHERE tenant_id = $1 AND deleted_at IS NULL;
