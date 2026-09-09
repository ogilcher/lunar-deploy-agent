# Lunar Deploy Agent

A lightweight deployment orchestration agent written in Go.

Lunar Deploy Agent can execute deployment pipelines, expose deployment APIs, stream deployment events in real-time, manage deployment queues, monitor node health, and integrate with process managers such as PM2.

---

## Features

### Deployment Engine

* Config-driven deployments
* Step-based execution pipeline
* Retry support
* Deployment history tracking
* Deployment result reporting

### Deployment Steps

* Git Pull
* Shell Commands
* NPM Install
* NPM Build
* PM2 Restart
* PM2 Start or Restart
* PM2 Save
* PM2 Status
* HTTP Health Checks

### HTTP API

* Deployment execution
* Deployment history
* Job management
* Queue visibility
* Node metadata
* Health reporting

### Real-Time Events

* WebSocket event streaming
* Deployment lifecycle events
* Job lifecycle events
* Step lifecycle events

### Health Monitoring

* Queue health
* Git health
* PM2 health
* Configuration validation
* Disk space monitoring
* Memory monitoring
* Capability validation
* Uptime reporting

### Documentation

* OpenAPI Specification
* Swagger UI

---

## Architecture

See:

```txt
docs/architecture.md
```

---

## Installation

### Requirements

* Go 1.26+
* Git
* Node.js (optional)
* PM2 (optional)

### Clone

```bash
git clone https://github.com/ogilcher/lunar-deploy-agent.git

cd lunar-deploy-agent
```

### Build

```bash
go build ./cmd/lunar-agent
```

### Run

```bash
./lunar-agent serve
```

---

## Example Configurations

Examples are available in:

```txt
examples/
```

Available examples:

```txt
config.basic.yaml
config.node-pm2.yaml
config.manual-pipeline.yaml
```

---

## API Documentation

Start the server:

```bash
go run ./cmd/lunar-agent serve
```

Swagger UI:

```txt
http://localhost:8080/swagger/index.html
```

OpenAPI Spec:

```txt
http://localhost:8080/openapi.yaml
```

---

## Authentication

Protected endpoints require:

```http
Authorization: Bearer <token>
```

Configured through:

```yaml
api_token: "your-token"
```

---

## Versioning

This project follows Semantic Versioning.

Examples:

```txt
v0.3.0
v1.0.0
v1.1.0
```

---

## Roadmap

### v1.x

* Process cancellation
* Rollback support
* Persistent queue

### v2.x

* Multi-node orchestration
* Distributed deployments
* Control plane integration

---

## License

MIT License