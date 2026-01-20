# =============================================================================
# SubFlow: Enterprise Construction Financial Ledger
# Multi-Stage Docker Build for Minimal Production Image
# Copyright (c) 2026 Muhammet Ali Büyük. All rights reserved.
# Contact: iletisim@alibuyuk.net | Website: alibuyuk.net
# =============================================================================

# Stage 1: Build Environment
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=1.0.0 -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" \
    -o /build/subflow \
    ./cmd/api

# Stage 2: Production Image (Minimal)
FROM alpine:3.19

# Architect signature
LABEL maintainer="Muhammet Ali Büyük <iletisim@alibuyuk.net>"
LABEL architect="Muhammet-Ali-Buyuk"
LABEL website="alibuyuk.net"
LABEL description="SubFlow - Enterprise Construction Financial Ledger & Compliance Engine"

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1000 subflow && \
    adduser -u 1000 -G subflow -s /bin/sh -D subflow

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/subflow /app/subflow

# Copy any static files or templates if needed
# COPY --from=builder /build/templates /app/templates

# Change ownership
RUN chown -R subflow:subflow /app

# Switch to non-root user
USER subflow

# Expose application port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health || exit 1

# Run the application
ENTRYPOINT ["/app/subflow"]
