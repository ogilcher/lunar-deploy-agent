# Lunar Deploy Agent Architecture

## Overview

Lunar Deploy Agent is a lightweight deployment orchestration service written in Go. It can be used as a command-line deployment tool or as a long-running HTTP service capable of receiving deployment requests, executing deployment pipelines, streaming deployment events, and reporting node health.

The agent is designed to be reusable and platform-agnostic. It provides deployment execution, monitoring, job management, and event streaming without being tied to any specific company or infrastructure provider.

---

# High-Level Architecture

```text
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  HTTP API   │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Job Queue   │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Worker    │
└──────┬──────┘
       │
       ▼
┌──────────────────────┐
│ Deployment Engine    │
└──────┬───────────────┘
       │
       ├── Git Pull
       ├── NPM Install
       ├── NPM Build
       ├── PM2 Actions
       └── Health Checks
       │
       ▼
┌──────────────────────┐
│ History + Events     │
└──────────────────────┘
```

---

# Deployment Flow

Deployments are executed through asynchronous jobs.

```text
POST /deploy
      │
      ▼
Create Deployment Job
      │
      ▼
Queue Deployment
      │
      ▼
Worker Picks Up Job
      │
      ▼
Deployment Engine Executes Steps
      │
      ▼
History Saved
      │
      ▼
Events Published
```

This architecture prevents long-running deployments from blocking HTTP requests and enables future support for multiple workers and distributed execution.

---

# Deployment Engine

The deployment engine is responsible for executing deployment steps in order.

Each deployment is composed of one or more deployment steps.

Example:

```text
git_pull
npm_install
npm_build
pm2_restart
http_health_check
```

The engine executes steps sequentially and stops on the first failure.

Deployment results include:

* Success status
* Step results
* Error information
* Start time
* Finish time
* Duration

---

# Deployment Steps

## Git Pull

Synchronizes the local repository with the configured remote.

## Shell

Executes arbitrary shell commands.

## NPM Install

Installs Node.js dependencies.

## NPM Build

Builds Node.js applications.

## PM2 Restart

Restarts an existing PM2 process.

## PM2 Start Or Restart

Creates or restarts a PM2 process.

## PM2 Save

Persists PM2 process definitions.

## PM2 Status

Validates PM2 process state.

## HTTP Health Check

Verifies application availability after deployment.

---

# Job Queue

Deployments are processed through an asynchronous job queue.

Each job transitions through the following lifecycle:

```text
Queued
   │
   ▼
Running
   │
   ├── Succeeded
   │
   ├── Failed
   │
   └── Cancelled
```

Job information is exposed through:

```text
GET /jobs
GET /jobs/{id}
DELETE /jobs/{id}
GET /queue
```

---

# Event System

The deployment agent includes a global event bus.

Deployment operations publish events such as:

```text
deployment_started
deployment_completed
deployment_failed

step_started
step_completed
step_failed

job_queued
job_running
job_succeeded
job_failed
job_cancelled
```

Events are distributed to WebSocket subscribers.

---

# WebSocket Streaming

Clients may subscribe to deployment events using:

```text
GET /events
```

WebSocket subscribers receive real-time deployment updates without polling.

Example event:

```json
{
  "type": "job_running",
  "job_id": "123456",
  "deployment": "production",
  "message": "Deployment job started."
}
```

---

# Health System

The agent continuously evaluates its ability to execute deployments.

Current health checks include:

## Queue Check

Verifies queue backlog remains within acceptable limits.

## PM2 Check

Verifies PM2 is installed and accessible.

## Git Check

Verifies Git is installed and accessible.

## Configuration Check

Verifies deployment configuration is valid.

## Disk Check

Verifies adequate disk space is available.

## Memory Check

Reports current runtime memory usage.

## Capability Check

Verifies required deployment tools exist for configured deployments.

## Uptime Check

Reports process uptime and startup time.

Health reports are exposed through:

```text
GET /health
```

---

# Node Metadata

Every agent instance exposes metadata describing itself.

Available information includes:

* Node ID
* Node Name
* Node Region
* Service Version
* Operating System
* Architecture
* Supported Capabilities

Endpoints:

```text
GET /node
GET /heartbeat
```

---

# Authentication

Operational endpoints are protected using API token authentication.

Clients authenticate using:

```text
Authorization: Bearer <token>
```

Public endpoints:

```text
GET /health
GET /openapi.yaml
GET /swagger/*
```

Protected endpoints:

```text
GET /status
GET /deployments
GET /history
GET /jobs
GET /queue
GET /node
GET /heartbeat
POST /deploy
DELETE /jobs/{id}
GET /events
```

---

# API Surface

## Public

```text
GET /health
GET /openapi.yaml
GET /swagger/*
```

## Protected

```text
GET /status
GET /deployments
GET /history

GET /jobs
GET /jobs/{id}
DELETE /jobs/{id}

GET /queue

GET /node
GET /heartbeat

POST /deploy

GET /events
```

---

# Future Roadmap

## Version 1.x

* Process cancellation
* Deployment rollback support
* Persistent queue storage
* Additional deployment presets
* Expanded health monitoring

## Version 2.x

* Multi-node deployments
* Node registration
* Distributed orchestration
* Remote scheduling
* Control plane integration
* Cluster-wide deployment coordination

---

# Design Goals

The Lunar Deploy Agent is designed around four principles:

1. Simplicity
2. Observability
3. Extensibility
4. Reliability

The agent should remain lightweight, easy to understand, and capable of serving as the foundation for larger deployment orchestration systems.
