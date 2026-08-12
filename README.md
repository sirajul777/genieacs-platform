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

The initial migration creates an `agents` table used by the repository foundation. Agent registration is intentionally not implemented yet.
