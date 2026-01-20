-- name: CreateTenant :one
-- Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
INSERT INTO tenants (
    name, slug, plan, contact_email, default_currency
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetTenantByID :one
SELECT * FROM tenants
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetTenantBySlug :one
SELECT * FROM tenants
WHERE slug = $1 AND deleted_at IS NULL;

-- name: ListTenants :many
SELECT * FROM tenants
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateTenant :one
UPDATE tenants
SET name = $2, plan = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SoftDeleteTenant :exec
UPDATE tenants
SET deleted_at = NOW()
WHERE id = $1;
