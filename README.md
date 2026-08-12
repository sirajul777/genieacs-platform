# GenieACS Platform

> **A distributed orchestration platform for managing GenieACS deployments at scale.**

GenieACS Platform is an open-source control plane designed to manage, automate, and orchestrate GenieACS deployments across multiple servers. Unlike `genieacs-mod`, which focuses on installation and server management through scripts, GenieACS Platform provides a scalable architecture built around Jobs, Agents, Extensions, and Events.

## Vision

Build the best open-source control plane for managing GenieACS deployments at scale.

The platform is designed around several core principles:

* **Everything is a Job**
* **Agent-based execution**
* **Extension-first architecture**
* **Event-driven communication**
* **Clean Architecture**
* **Production-ready engineering**

---

# Why?

Managing multiple GenieACS servers manually quickly becomes difficult.

Typical operations include:

* Installing new GenieACS servers
* Upgrading existing deployments
* Backing up MongoDB
* Restoring configurations
* Restarting services
* Health verification
* Rolling upgrades

Instead of creating a dedicated API endpoint for every operation, GenieACS Platform models every operation as a **Job**.

```
Create Job
      │
      ▼
Job Queue
      │
      ▼
Agent
      │
      ▼
Worker
      │
      ▼
Result
```

This architecture makes the platform extensible and easy to maintain.

---

# Goals

* Distributed orchestration
* Multi-Agent support
* Multi-GenieACS deployment support
* Job-based execution model
* Extension ecosystem
* Production-ready architecture
* Open-source community driven

---

# Core Concepts

## Manager

The control plane responsible for:

* Agent registration
* Authentication
* Job scheduling
* Job assignment
* Job tracking
* Event publishing

The Manager **never executes commands directly**.

---

## Agent

The execution plane.

Responsibilities:

* Register to Manager
* Send heartbeat
* Poll pending jobs
* Execute workers
* Report progress
* Return execution results

Agents are intentionally lightweight and mostly stateless.

---

## Job

Everything is represented as a Job.

Examples:

* INSTALL_GENIEACS
* UPGRADE
* BACKUP
* RESTORE
* VERIFY
* RESTART

Future jobs can be added without changing the API.

---

## Worker

A Worker executes exactly one Job type.

```text
INSTALL Worker

BACKUP Worker

RESTORE Worker

VERIFY Worker
```

Workers should only focus on execution.

They should **not**:

* perform HTTP communication
* manage retries
* access global configuration
* create workspaces

Those responsibilities belong to the Runtime.

---

## Runtime

The Runtime lives inside the Agent.

Responsibilities:

* Poll Manager
* Dispatch jobs
* Execute workers
* Report progress
* Handle cancellation
* Manage workspace

---

## Extension

The core platform knows nothing about GenieACS.

All platform-specific logic lives inside extensions.

Example:

```
extensions/

    genieacs/

        installer/

        backup/

        restore/

        verify/
```

Future extensions may include:

* Docker
* MongoDB
* Nginx
* Redis
* Systemd

---

# Architecture Principles

## Everything is a Job

No API like:

```
POST /restart
POST /backup
POST /restore
```

Instead:

```
POST /api/v1/jobs
```

with different job types.

---

## Agent Executes

Manager never:

* SSH
* runs shell scripts
* installs packages

Agents execute everything locally.

---

## Core is Generic

Core only understands:

* Agent
* Job
* Worker
* Runtime
* Extension
* Event

It does not understand GenieACS.

---

## Event Driven

Every state transition emits an event.

Examples:

* JobCreated
* JobAssigned
* JobStarted
* JobCompleted
* JobFailed

These events can later power:

* Audit Log
* Dashboard
* Notifications
* Scheduler

---

## Extension First

Every new feature should become an Extension instead of modifying the Core.

---

# Repository Structure

```text
cmd/
    manager/
    agent/

internal/
    manager/
    agent/
    platform/

pkg/
    sdk/

extensions/
    genieacs/

docs/
    adr/
    rfc/

migrations/

examples/
```

---

# Development Roadmap

## Milestone 0

* Bootstrap
* Config
* Logger
* Database
* Migration

## Milestone 1

* Agent Registration
* Heartbeat
* Job Engine
* Runtime

## Milestone 2

* GenieACS Installer Extension

## Milestone 3

* Multi Instance Management

## Milestone 4

* Backup & Restore

## Milestone 5

* Scheduler

## Milestone 6

* Dashboard API

---

# Engineering Standards

* Clean Architecture
* Dependency inversion
* SQLC
* PostgreSQL
* golangci-lint
* Unit testing
* Handler testing
* Versioned API
* Structured logging
* Observability

---

# Long-term Vision

GenieACS Platform aims to become a general orchestration framework where GenieACS is simply the first official extension.

The platform should remain generic enough to manage additional infrastructure components without changing its core architecture.

---

# License

Apache-2.0 (recommended)

---

# Status

🚧 Under active development.
