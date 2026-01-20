-- =============================================================================
-- SubFlow: Enterprise Construction Financial Ledger
-- PostgreSQL Database Initialization Script
-- Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
-- Contact: iletisim@alibuyuk.net | Website: alibuyuk.net
-- =============================================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =============================================================================
-- TENANTS (Multi-Tenant SaaS)
-- =============================================================================
CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    plan VARCHAR(50) DEFAULT 'FREE' CHECK (plan IN ('FREE', 'PRO', 'ENTERPRISE')),
    is_active BOOLEAN DEFAULT TRUE,
    max_users INTEGER DEFAULT 5,
    max_projects INTEGER DEFAULT 3,
    default_currency CHAR(3) DEFAULT 'TRY',
    contact_email VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(50),
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_active ON tenants(is_active) WHERE deleted_at IS NULL;

-- =============================================================================
-- USERS
-- =============================================================================
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) DEFAULT 'VIEWER' CHECK (role IN ('ADMIN', 'MANAGER', 'ACCOUNTANT', 'VIEWER')),
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, email)
);

CREATE INDEX idx_users_tenant ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);

-- =============================================================================
-- PROJECTS
-- =============================================================================
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'ACTIVE', 'ON_HOLD', 'COMPLETED', 'CANCELLED')),
    contract_amount_cents BIGINT DEFAULT 0,
    currency CHAR(3) DEFAULT 'TRY',
    start_date DATE,
    estimated_end_date DATE,
    labor_retainage_rate DECIMAL(5,4) DEFAULT 0.1000, -- 10%
    material_retainage_rate DECIMAL(5,4) DEFAULT 0.0500, -- 5%
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(tenant_id, code)
);

CREATE INDEX idx_projects_tenant ON projects(tenant_id);
CREATE INDEX idx_projects_status ON projects(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_code ON projects(code);

-- =============================================================================
-- CONTRACTS (Subcontractor Agreements)
-- =============================================================================
CREATE TABLE IF NOT EXISTS contracts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    vendor_name VARCHAR(255) NOT NULL,
    vendor_tax_id VARCHAR(50),
    contract_amount_cents BIGINT NOT NULL,
    currency CHAR(3) DEFAULT 'TRY',
    scope_of_work TEXT,
    start_date DATE,
    end_date DATE,
    retainage_rate DECIMAL(5,4) DEFAULT 0.1000,
    status VARCHAR(50) DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'ACTIVE', 'COMPLETED', 'TERMINATED')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_contracts_project ON contracts(project_id);
CREATE INDEX idx_contracts_vendor ON contracts(vendor_name);

-- =============================================================================
-- TRANSACTIONS (Immutable Financial Ledger)
-- This is the heart of the system. NEVER UPDATE OR DELETE rows here.
-- =============================================================================
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id),
    contract_id UUID REFERENCES contracts(id),
    type VARCHAR(50) NOT NULL CHECK (type IN (
        'INVOICE',
        'PAYMENT',
        'RETAINAGE_HELD',
        'RETAINAGE_RELEASE',
        'ADJUSTMENT',
        'DEDUCTION'
    )),
    amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
    currency CHAR(3) NOT NULL DEFAULT 'TRY',
    effective_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    description TEXT,
    reference_no VARCHAR(100),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id),
    
    -- Digital fingerprint - proof of ownership
    _architect_signature VARCHAR(100) DEFAULT 'Muhammet-Ali-Buyuk-SF2026'
);

-- Prevent updates and deletes on transactions (immutable ledger)
CREATE OR REPLACE FUNCTION prevent_transaction_modification()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Transactions are immutable. Cannot modify or delete.';
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_transactions_immutable_update
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION prevent_transaction_modification();

CREATE TRIGGER tr_transactions_immutable_delete
    BEFORE DELETE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION prevent_transaction_modification();

CREATE INDEX idx_transactions_project ON transactions(project_id);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_date ON transactions(effective_date);
CREATE INDEX idx_transactions_created_by ON transactions(created_by);

-- =============================================================================
-- AUDIT LOGS (Change Tracking)
-- =============================================================================
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    actor_id UUID REFERENCES users(id),
    actor_email VARCHAR(255),
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID,
    old_value JSONB,
    new_value JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_audit_tenant ON audit_logs(tenant_id);
CREATE INDEX idx_audit_actor ON audit_logs(actor_id);
CREATE INDEX idx_audit_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_date ON audit_logs(created_at);

-- =============================================================================
-- MATERIALIZED VIEW: Project Financial Summary
-- Aggregated view for fast financial snapshots
-- =============================================================================
CREATE MATERIALIZED VIEW IF NOT EXISTS mv_project_financials AS
SELECT
    p.id AS project_id,
    p.tenant_id,
    p.name AS project_name,
    p.contract_amount_cents,
    p.currency,
    COALESCE(SUM(CASE WHEN t.type = 'INVOICE' THEN t.amount_cents ELSE 0 END), 0) AS total_invoiced,
    COALESCE(SUM(CASE WHEN t.type = 'PAYMENT' THEN t.amount_cents ELSE 0 END), 0) AS total_paid,
    COALESCE(SUM(CASE WHEN t.type = 'RETAINAGE_HELD' THEN t.amount_cents ELSE 0 END), 0) AS retainage_held,
    COALESCE(SUM(CASE WHEN t.type = 'RETAINAGE_RELEASE' THEN t.amount_cents ELSE 0 END), 0) AS retainage_released,
    COUNT(t.id) AS transaction_count,
    MAX(t.created_at) AS last_transaction_at
FROM projects p
LEFT JOIN transactions t ON p.id = t.project_id
WHERE p.deleted_at IS NULL
GROUP BY p.id, p.tenant_id, p.name, p.contract_amount_cents, p.currency;

CREATE UNIQUE INDEX ON mv_project_financials(project_id);

-- =============================================================================
-- FUNCTIONS: Financial Calculations
-- =============================================================================

-- Get project balance (invoiced - paid)
CREATE OR REPLACE FUNCTION get_project_balance(p_project_id UUID)
RETURNS BIGINT AS $$
DECLARE
    v_balance BIGINT;
BEGIN
    SELECT 
        COALESCE(SUM(
            CASE 
                WHEN type = 'INVOICE' THEN amount_cents
                WHEN type = 'PAYMENT' THEN -amount_cents
                WHEN type = 'RETAINAGE_RELEASE' THEN -amount_cents
                WHEN type = 'DEDUCTION' THEN -amount_cents
                ELSE 0
            END
        ), 0)
    INTO v_balance
    FROM transactions
    WHERE project_id = p_project_id;
    
    RETURN v_balance;
END;
$$ LANGUAGE plpgsql STABLE;

-- =============================================================================
-- SEED DATA (Demo)
-- =============================================================================
INSERT INTO tenants (id, name, slug, plan, contact_email) VALUES
    ('11111111-1111-1111-1111-111111111111', 'Demo Company', 'demo', 'PRO', 'demo@example.com')
ON CONFLICT (slug) DO NOTHING;

-- Architect signature comment
COMMENT ON DATABASE subflow IS 'SubFlow Enterprise Construction Financial Ledger - Architect: Muhammet Ali Büyük (alibuyuk.net)';
