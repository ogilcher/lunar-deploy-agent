# Lunar Deploy Agent

Lunar Deploy Agent is a lightweight deployment orchestration agent written in Go.

The project is part of the larger Lunar Deploy Control ecosystem used to automate deployments, health monitoring,
environment management, and infrastructure operations for multi-tenant web platforms.

## Goals

- Lightweight and fast
- Self-hostable
- Cross-platform
- Secure deployment execution
- Real-time deployment logging
- Infrastructure health reporting
- Container and process management
- Git-based deployment workflows

## Planned Features

- CLI command system
- Deployment execution engine
- Health monitoring
- WebSocket communication
- Secure agent registration
- Multi-environment support
- Rollback support
- Docker integration
- GitHub Actions integration
- Vercel deployment hooks
- Structured logging
- Metrics and telemetry

## Current Status

Early development.

Current functionality:
- Basic CLI entry point
- Go project structure initialized

## Project Structure

```txt
cmd/
    lunar-agent/
        main.go
        
internal/
    config/
    deploy/
    health/
    logger/
```

## Running the Project

```bash
go run ./cmd/lunar-agent
```

## Requirements

- Go 1.24+ recommended

## Vision

Lunar Deploy Agent is intended to become a reusable, production-grade deployment agent suitable for:

- Personal infrastructure
- Startup operations
- Contractor deployments
- Internal platform tooling
- Edge/server deployments
- Self-hosted environments
