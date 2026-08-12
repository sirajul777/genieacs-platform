# GenieACS Platform

A Go monorepo for the GenieACS platform, bootstrapped with Clean Architecture boundaries, Chi HTTP routing, Cobra CLIs, Viper configuration, and Zap logging.

## Module

```text
github.com/sirajul777/genieacs-platform
```

## Requirements

- Go 1.25+

## Project layout

```text
cmd/                         # Application entrypoints
  manager/                   # Manager CLI/main
  agent/                     # Agent CLI/main
internal/
  manager/                   # Manager application internals
    server/                  # Manager HTTP server composition
    health/                  # Manager health endpoint
    config/                  # Manager config defaults/types
  agent/                     # Agent application internals
    server/                  # Agent HTTP server composition
    health/                  # Agent health endpoint
    config/                  # Agent config defaults/types
  shared/                    # Cross-application primitives
    logger/                  # Zap logger construction
    config/                  # Viper config loader
    version/                 # Build/version metadata
    response/                # HTTP response helpers
configs/                     # Example config files
docs/                        # Documentation
deployments/                 # Deployment assets
scripts/                     # Operational scripts
```

## Commands

```bash
make tidy
make test
make build
make run-manager
make run-agent
```

## HTTP endpoints

Both applications expose:

- `GET /health` - returns service health metadata.

## Configuration

Configuration is loaded by Viper from defaults, optional config files, environment variables, and CLI flags. Environment variables use the service name as a prefix and `.` is mapped to `_`.

| Service | Default address | Env prefix |
| --- | --- | --- |
| Manager | `:8080` | `MANAGER` |
| Agent | `:8081` | `AGENT` |

Examples:

```bash
MANAGER_SERVER_ADDRESS=:9000 go run ./cmd/manager
AGENT_SERVER_ADDRESS=:9001 go run ./cmd/agent
```

## Database

The manager service uses PostgreSQL through pgx. Schema changes are managed with golang-migrate and type-safe query code is generated from SQL definitions with sqlc.

Start PostgreSQL locally:

```bash
docker compose up -d postgres
```

Run migrations:

```bash
make migrate-up
```

Regenerate query code:

```bash
make sqlc
```

The initial migrations create the `agents` and `heartbeats` tables used by registration and heartbeat persistence.

## Manager Agent API

Agent registration and heartbeat endpoints are available under `/api/v1/agents`.

Register an agent:

```bash
curl -X POST http://localhost:8080/api/v1/agents/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"agent-1","endpoint":"http://agent:8081"}'
```

The registration response includes a bearer token once. The manager stores only a SHA-256 hash of the 32-byte random token.

Send a heartbeat:

```bash
curl -X POST http://localhost:8080/api/v1/agents/heartbeat \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <token>' \
  -d '{"agent_id":"<agent-id>"}'
```

## Manager Job API

The manager includes a generic job engine. `type` is intentionally generic so future extensions can define concrete job types.

Create a job:

```bash
curl -X POST http://localhost:8080/api/v1/jobs \
  -H 'Content-Type: application/json' \
  -d '{"type":"generic.example","payload":{"target":"device-1"}}'
```

Poll the next pending job for an agent:

```bash
curl -X POST http://localhost:8080/api/v1/jobs/poll \
  -H 'Content-Type: application/json' \
  -d '{"agent_id":"<agent-id>"}'
```

Agents can then report progress, complete the job, or fail it:

```bash
curl -X POST http://localhost:8080/api/v1/jobs/<job-id>/progress -d '{"progress":50}'
curl -X POST http://localhost:8080/api/v1/jobs/<job-id>/complete -d '{"result":{"ok":true}}'
curl -X POST http://localhost:8080/api/v1/jobs/<job-id>/fail -d '{"error":"reason"}'
```
