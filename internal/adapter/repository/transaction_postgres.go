// Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
// This source code is proprietary. Confidential and private.
// Unauthorized copying or distribution is strictly prohibited.
// Contact: iletisim@alibuyuk.net | Website: alibuyuk.net

package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qantesm/subflow/internal/core/entity"
	"github.com/qantesm/subflow/internal/core/service"
)

// PostgresTransactionRepository implements TransactionRepository for PostgreSQL
type PostgresTransactionRepository struct {
	pool      *pgxpool.Pool
	architect string
}

// NewPostgresTransactionRepository creates a new PostgreSQL transaction repository
func NewPostgresTransactionRepository(pool *pgxpool.Pool) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{
		pool:      pool,
		architect: "Muhammet-Ali-Buyuk",
	}
}

// Save stores a transaction in the database
func (r *PostgresTransactionRepository) Save(ctx context.Context, tx *entity.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, project_id, contract_id, type, amount_cents, currency,
			effective_date, description, reference_no, metadata, created_by, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.pool.Exec(ctx, query,
		tx.ID,
		tx.ProjectID,
		tx.ContractID,
		tx.Type,
		tx.AmountCents,
		tx.Currency,
		tx.EffectiveDate,
		tx.Description,
		tx.ReferenceNo,
		tx.Metadata,
		tx.CreatedBy,
		tx.CreatedAt,
	)

	return err
}

// FindByID retrieves a transaction by its ID
func (r *PostgresTransactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Transaction, error) {
	query := `
		SELECT id, project_id, contract_id, type, amount_cents, currency,
			   effective_date, description, reference_no, metadata, created_by, created_at
		FROM transactions
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	return r.scanTransaction(row)
}

// FindByProjectID retrieves all transactions for a project
func (r *PostgresTransactionRepository) FindByProjectID(ctx context.Context, projectID uuid.UUID) ([]*entity.Transaction, error) {
	query := `
		SELECT id, project_id, contract_id, type, amount_cents, currency,
			   effective_date, description, reference_no, metadata, created_by, created_at
		FROM transactions
		WHERE project_id = $1
		ORDER BY effective_date DESC, created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*entity.Transaction
	for rows.Next() {
		tx, err := r.scanTransactionFromRows(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, rows.Err()
}

// GetProjectSummary calculates the financial summary for a project
func (r *PostgresTransactionRepository) GetProjectSummary(ctx context.Context, projectID uuid.UUID) (*service.LedgerSummary, error) {
	query := `
		SELECT 
			project_id,
			COUNT(*) as transaction_count,
			COALESCE(SUM(CASE WHEN type = 'INVOICE' THEN amount_cents ELSE 0 END), 0) as total_invoiced,
			COALESCE(SUM(CASE WHEN type = 'PAYMENT' THEN amount_cents ELSE 0 END), 0) as total_paid,
			COALESCE(SUM(CASE WHEN type = 'RETAINAGE_HELD' THEN amount_cents ELSE 0 END), 0) as retainage_held,
			COALESCE(SUM(CASE WHEN type = 'RETAINAGE_RELEASE' THEN amount_cents ELSE 0 END), 0) as retainage_released,
			COALESCE(MAX(currency), 'TRY') as currency
		FROM transactions
		WHERE project_id = $1
		GROUP BY project_id
	`

	row := r.pool.QueryRow(ctx, query, projectID)

	var (
		pID              uuid.UUID
		txCount          int
		totalInvoiced    int64
		totalPaid        int64
		retainageHeld    int64
		retainageRelease int64
		currency         string
	)

	err := row.Scan(&pID, &txCount, &totalInvoiced, &totalPaid, &retainageHeld, &retainageRelease, &currency)
	if err == pgx.ErrNoRows {
		// No transactions yet, return empty summary
		return &service.LedgerSummary{
			ProjectID: projectID,
			Currency:  "TRY",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return &service.LedgerSummary{
		ProjectID:        pID,
		TotalInvoiced:    totalInvoiced,
		TotalPaid:        totalPaid,
		TotalRetained:    retainageHeld - retainageRelease,
		CurrentBalance:   totalInvoiced - totalPaid,
		Currency:         currency,
		TransactionCount: txCount,
	}, nil
}

// Helper function to scan a transaction from a row
func (r *PostgresTransactionRepository) scanTransaction(row pgx.Row) (*entity.Transaction, error) {
	tx := &entity.Transaction{}
	var metadata []byte

	err := row.Scan(
		&tx.ID,
		&tx.ProjectID,
		&tx.ContractID,
		&tx.Type,
		&tx.AmountCents,
		&tx.Currency,
		&tx.EffectiveDate,
		&tx.Description,
		&tx.ReferenceNo,
		&metadata,
		&tx.CreatedBy,
		&tx.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, entity.ErrTransactionNotFound
	}
	if err != nil {
		return nil, err
	}

	if metadata != nil {
		tx.Metadata = json.RawMessage(metadata)
	}

	return tx, nil
}

// Helper function to scan a transaction from rows
func (r *PostgresTransactionRepository) scanTransactionFromRows(rows pgx.Rows) (*entity.Transaction, error) {
	tx := &entity.Transaction{}
	var metadata []byte

	err := rows.Scan(
		&tx.ID,
		&tx.ProjectID,
		&tx.ContractID,
		&tx.Type,
		&tx.AmountCents,
		&tx.Currency,
		&tx.EffectiveDate,
		&tx.Description,
		&tx.ReferenceNo,
		&metadata,
		&tx.CreatedBy,
		&tx.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if metadata != nil {
		tx.Metadata = json.RawMessage(metadata)
	}

	return tx, nil
}

// PostgresProjectRepository implements project persistence for PostgreSQL
type PostgresProjectRepository struct {
	pool      *pgxpool.Pool
	architect string
}

// NewPostgresProjectRepository creates a new PostgreSQL project repository
func NewPostgresProjectRepository(pool *pgxpool.Pool) *PostgresProjectRepository {
	return &PostgresProjectRepository{
		pool:      pool,
		architect: "Muhammet-Ali-Buyuk",
	}
}

// Create stores a new project in the database
func (r *PostgresProjectRepository) Create(ctx context.Context, p *entity.Project) error {
	query := `
		INSERT INTO projects (
			id, tenant_id, name, code, description, status,
			contract_amount_cents, currency, start_date, estimated_end_date,
			labor_retainage_rate, material_retainage_rate, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.pool.Exec(ctx, query,
		p.ID,
		p.TenantID,
		p.Name,
		p.Code,
		p.Description,
		p.Status,
		p.ContractAmount,
		p.Currency,
		p.StartDate,
		p.EstimatedEndDate,
		p.LaborRetainageRate,
		p.MaterialRetainageRate,
		p.CreatedAt,
		p.UpdatedAt,
	)

	return err
}

// FindByID retrieves a project by its ID
func (r *PostgresProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Project, error) {
	query := `
		SELECT id, tenant_id, name, code, description, status,
			   contract_amount_cents, currency, start_date, estimated_end_date,
			   labor_retainage_rate, material_retainage_rate, created_at, updated_at, deleted_at
		FROM projects
		WHERE id = $1 AND deleted_at IS NULL
	`

	row := r.pool.QueryRow(ctx, query, id)
	return r.scanProject(row)
}

// FindByTenant retrieves all projects for a tenant
func (r *PostgresProjectRepository) FindByTenant(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entity.Project, error) {
	query := `
		SELECT id, tenant_id, name, code, description, status,
			   contract_amount_cents, currency, start_date, estimated_end_date,
			   labor_retainage_rate, material_retainage_rate, created_at, updated_at, deleted_at
		FROM projects
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*entity.Project
	for rows.Next() {
		p, err := r.scanProjectFromRows(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, rows.Err()
}

// Update modifies an existing project
func (r *PostgresProjectRepository) Update(ctx context.Context, p *entity.Project) error {
	query := `
		UPDATE projects SET
			name = $2,
			description = $3,
			status = $4,
			contract_amount_cents = $5,
			start_date = $6,
			estimated_end_date = $7,
			labor_retainage_rate = $8,
			material_retainage_rate = $9,
			updated_at = $10
		WHERE id = $1 AND deleted_at IS NULL
	`

	p.UpdatedAt = time.Now()
	_, err := r.pool.Exec(ctx, query,
		p.ID,
		p.Name,
		p.Description,
		p.Status,
		p.ContractAmount,
		p.StartDate,
		p.EstimatedEndDate,
		p.LaborRetainageRate,
		p.MaterialRetainageRate,
		p.UpdatedAt,
	)

	return err
}

// SoftDelete marks a project as deleted
func (r *PostgresProjectRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE projects SET deleted_at = NOW() WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *PostgresProjectRepository) scanProject(row pgx.Row) (*entity.Project, error) {
	p := &entity.Project{}

	err := row.Scan(
		&p.ID,
		&p.TenantID,
		&p.Name,
		&p.Code,
		&p.Description,
		&p.Status,
		&p.ContractAmount,
		&p.Currency,
		&p.StartDate,
		&p.EstimatedEndDate,
		&p.LaborRetainageRate,
		&p.MaterialRetainageRate,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, entity.ErrProjectNotFound
	}
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *PostgresProjectRepository) scanProjectFromRows(rows pgx.Rows) (*entity.Project, error) {
	p := &entity.Project{}

	err := rows.Scan(
		&p.ID,
		&p.TenantID,
		&p.Name,
		&p.Code,
		&p.Description,
		&p.Status,
		&p.ContractAmount,
		&p.Currency,
		&p.StartDate,
		&p.EstimatedEndDate,
		&p.LaborRetainageRate,
		&p.MaterialRetainageRate,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return p, nil
}
