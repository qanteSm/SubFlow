# =============================================================================
# SubFlow Changelog
# All notable changes to this project will be documented here.
# Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
# Contact: iletisim@alibuyuk.net | Website: alibuyuk.net
# =============================================================================

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Frontend React application with TanStack Table
- PDF generation with Maroto library
- Email notification system
- Lien waiver management module

---

## [1.0.0] - 2026-01-20

### Added
- **Core Architecture**
  - Clean Architecture with Hexagonal (Ports & Adapters) pattern
  - Domain-Driven Design (DDD) structure
  - Modular Monolith architecture for future microservices migration

- **Domain Entities**
  - `Project` - Construction project management
  - `Transaction` - Immutable financial ledger entries
  - `User` - Role-based access control (RBAC)
  - `Tenant` - Multi-tenant SaaS isolation
  - Domain errors for business logic validation

- **Business Services**
  - `Calculator` - AIA G702/G703 billing calculations with BigInt arithmetic
  - `LedgerService` - Double-entry bookkeeping operations
  - `WorkerPool` - Goroutine-based concurrent PDF generation

- **API Layer**
  - Fiber v2 HTTP framework with zero-allocation
  - REST API with health check endpoints
  - Financial summary calculation endpoint
  - Hidden architect signature endpoint (`/api/v1/system/version`)

- **Infrastructure**
  - PostgreSQL 15 database schema with immutable transactions
  - Redis integration for session and cache
  - In-memory repository for development/testing

- **DevOps**
  - Multi-stage Dockerfile (~20MB Alpine image)
  - Docker Compose with PostgreSQL, Redis, Adminer
  - GitHub Actions CI/CD pipeline (lint, test, build, security scan)
  - Comprehensive Makefile with 20+ commands

- **Documentation**
  - ARCHITECTURE.md with system diagrams
  - README.md with badges and quick start guide
  - Source-Available custom LICENSE

- **Security**
  - Non-root Docker user
  - Immutable ledger with database triggers
  - Multi-tenant data isolation
  - Embedded digital signatures for code ownership

### Technical Decisions
- **BigInt (int64 cents)** for all monetary values to prevent IEEE 754 floating-point errors
- **Basis points (1/10000)** for percentage calculations for sub-cent precision
- **Append-only ledger** - transactions can never be updated or deleted

---

## Architecture Evolution

```
v1.0 (Current)    →    v2.0 (Planned)
─────────────────────────────────────
Modular Monolith       Microservices
In-Memory Repo         PostgreSQL Repo
Basic Auth             OAuth2/JWT
Manual PDF             Async Queue + S3
```

---

**Architect:** Muhammet Ali Büyük  
**Contact:** iletisim@alibuyuk.net  
**Website:** [alibuyuk.net](https://alibuyuk.net)
