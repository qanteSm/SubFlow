-- name: CreateUser :one
-- Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
INSERT INTO users (
    tenant_id, email, password_hash, first_name, last_name, role
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND is_active = true;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE tenant_id = $1 AND email = $2 AND is_active = true;

-- name: ListUsersByTenant :many
SELECT * FROM users
WHERE tenant_id = $1 AND is_active = true
ORDER BY last_name, first_name
LIMIT $2 OFFSET $3;

-- name: UpdateUser :one
UPDATE users
SET 
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    role = COALESCE($4, role),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateLastLogin :exec
UPDATE users
SET last_login_at = NOW()
WHERE id = $1;

-- name: DeactivateUser :exec
UPDATE users
SET is_active = false, updated_at = NOW()
WHERE id = $1;

-- name: GetUserCount :one
SELECT COUNT(*) FROM users
WHERE tenant_id = $1 AND is_active = true;
