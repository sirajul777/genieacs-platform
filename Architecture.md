# Architecture

## Overview

GenieACS Platform follows a **Control Plane / Execution Plane** architecture.

```
                        +----------------------+
                        |      Dashboard       |
                        +----------+-----------+
                                   |
                            REST API / SDK
                                   |
                                   ▼
                      +-------------------------+
                      |        Manager          |
                      +-----------+-------------+
                                  |
              +-------------------+-------------------+
              |                   |                   |
              ▼                   ▼                   ▼
      Authentication        Job Engine          Event Bus
              |                   |                   |
              +-------------------+-------------------+
                                  |
                                  ▼
                        PostgreSQL Database
                                  |
                                  ▼
                          Agent Communication
                                  |
              +-------------------+-------------------+
              ▼                                       ▼
         Agent Runtime                         Agent Runtime
              |                                       |
        Worker Registry                        Worker Registry
              |                                       |
      +-------+-------+                     +---------+--------+
      |       |       |                     |         |        |
  Installer Backup Verify              Installer Backup Verify
```

---

# Control Plane

The Control Plane is implemented by the Manager.

Responsibilities:

* Authentication
* Agent registration
* Heartbeats
* Job scheduling
* Job assignment
* Event publishing
* Audit logging

The Manager never performs infrastructure operations directly.

---

# Execution Plane

The Execution Plane consists of Agents.

Responsibilities:

* Poll pending jobs
* Execute workers
* Report progress
* Upload results
* Handle cancellation

Execution is always local to the Agent.

---

# Job Lifecycle

```
Create

↓

Pending

↓

Assigned

↓

Running

↓

Success
```

or

```
Running

↓

Failed
```

Possible states:

* PENDING
* ASSIGNED
* RUNNING
* SUCCESS
* FAILED
* CANCELED

---

# Runtime

Runtime owns execution.

Responsibilities:

* Workspace creation
* Context cancellation
* Worker dispatch
* Progress reporting
* Cleanup

Workers receive only the execution context.

---

# Worker Model

Every worker implements:

```go
type Worker interface {
    Type() JobType
    Execute(ctx context.Context, job Job) (*Result, error)
}
```

Workers should remain stateless.

---

# Extension Model

Extensions register workers.

```
Extension

↓

Register()

↓

Worker Registry

↓

Worker
```

The Core never imports extension packages.

Extensions depend on Core interfaces.

---

# Event Flow

```
Job Created
        │
        ▼
Publish Event
        │
        ▼
Subscribers

    Audit

    Scheduler

    Dashboard

    Notification
```

The event bus starts as in-memory and can later be replaced with Redis Streams, NATS, or Kafka.

---

# API Design

All APIs are versioned.

```
/api/v1
```

Examples:

```
POST /api/v1/agents/register
POST /api/v1/agents/heartbeat

POST /api/v1/jobs
POST /api/v1/jobs/poll

POST /api/v1/jobs/{id}/progress
POST /api/v1/jobs/{id}/complete
POST /api/v1/jobs/{id}/fail
```

---

# Workspace Layout

Each job has its own isolated workspace.

```
/var/lib/genieacs-platform/

jobs/

    <job-id>/

        work/

        stdout.log

        stderr.log

        metadata.json

        artifacts/
```

---

# Data Model

```
Agent

1 ─────── N

Instance

1 ─────── N

Job

1 ─────── N

Job Event
```

This allows a single Agent to manage multiple GenieACS instances.

---

# Future Components

Planned additions include:

* Scheduler
* Retry Engine
* Metrics
* Audit Log
* Web Dashboard
* RBAC
* Plugin SDK
* Fleet Management
* Rolling Upgrade
* Auto Rollback

---

# Architectural Rules

1. Core must never depend on extensions.
2. Everything is executed through Jobs.
3. Workers must be independent.
4. Business logic belongs in use cases, not HTTP handlers.
5. SQL access must go through repositories.
6. Manager never executes infrastructure commands.
7. Agents execute all workloads locally.
8. Every important state transition emits an event.
9. Public APIs must remain backward compatible within a major version.
10. New functionality should be added through Extensions whenever possible.

These rules are considered architectural constraints for the project and should guide future development decisions.
s