# Go + HTMX + SQLite + AI

Fork of [Piszmog/go-htmx-template](https://github.com/Piszmog/go-htmx-template) — a Web Application built with Go (`templ`), HTMX, a SQL DB (`sqlc`), E2E testing (Playwright), and styling (Tailwind CSS).

**Live demo:** https://go-htmx-sqlite-ai-production.up.railway.app/

## Installation

**Prerequisites**

- Go (see `go.mod` for version)
- [air](https://github.com/air-verse/air#installation) for live reload

```shell
git clone https://github.com/Weerapat1993/go-htmx-sqlite-ai.git
cd go-htmx-sqlite-ai
go mod tidy
```

Generate sqlc and templ files:

```shell
go tool sqlc generate
go tool templ generate -path ./internal/components
```

Apply DB migrations:

```shell
./migrate.sh -p sqlite -u ./db.sqlite3
```

For Turso (remote libSQL) instead of local SQLite:

```shell
./migrate.sh -p libsql -u <db>-<org>.turso.io -t <auth-token>
```

Copy environment config (optional, customize as needed):

```shell
cp .env.example .env
```

## Run

`air` is the primary way to run the applications for local development. It watches for file changes. When a file changes, it will rebuild and re-run the application. On each rebuild, `air` also applies any pending DB migrations via `./migrate.sh`.

If running the binary directly (not via `air`), run migrations first:

```shell
./migrate.sh -p sqlite -u ./db.sqlite3
```

**Run Dev on air**
```shell
air
```
When the application is running, go to http://localhost:8080/

**Run on Docker Compose**
```shell
docker compose -f docker-compose.dev.yml up -d
```
When the application is running, go to http://localhost:3500/

## Environment Variables

The application can be configured using environment variables. For local development, copy `.env.example` to `.env` and customize as needed.

### Available Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `LOG_LEVEL` | `info` | Logging level: `debug`, `info`, `warn`, `error` |
| `LOG_OUTPUT` | `text` | Log format: `text` or `json` |
| `DB_URL` | `./db.sqlite3` | Local SQLite file path, or a `libsql://<db>-<org>.turso.io` URL for Turso (remote) |
| `TURSO_AUTH_TOKEN` | _(none)_ | Auth token for Turso — only used when `DB_URL` is a `libsql://` URL |
| `RATE_LIMIT` | `50` | Requests per minute per IP address |

### Example

```bash
# .env — local SQLite
PORT=3000
LOG_LEVEL=debug
DB_URL=/data/myapp.db
```

```bash
# .env — Turso (remote)
DB_URL=libsql://myapp-myorg.turso.io
TURSO_AUTH_TOKEN=eyJhbGciOi...
```

## Endpoints

| Method | Path | Description |
|--------|------|--------------|
| GET | `/` | Homepage |
| GET | `/about` | About page |
| POST | `/count` | HTMX counter demo |
| GET | `/todo-list-db` | Todo list page |
| POST | `/todos` | Create todo |
| GET | `/todos/{id}` | Get todo |
| PUT | `/todos/{id}` | Update todo |
| DELETE | `/todos/{id}` | Delete todo |
| GET | `/todos/{id}/edit` | Edit todo form |
| POST | `/todos/{id}/toggle` | Toggle todo done |
| GET | `/api/health` | Health check — returns `200 OK` with `{"version":"dev"}` |

`templ`, `sqlc`, and `tailwindcss` (via [`go-tw`](https://github.com/Piszmog/go-tw)) are included as `go tool` directives. When running
the application for the first time, it may take a little time as `templ`, `sqlc` and `go-tw` are being downloaded and installed.

## Technologies

A few different technologies are configured to help getting off the ground easier.

- [sqlc](https://sqlc.dev/) for database layer
  - Stubbed to use SQLite
  - This can be easily swapped with [sqlx](https://jmoiron.github.io/sqlx/)
- [Tailwind CSS](https://tailwindcss.com/) for styling
  - Output is generated with the [CLI](https://tailwindcss.com/docs/installation/tailwind-cli)
- [templ](https://templ.guide/) for creating HTML
- [HTMX](https://htmx.org/) for HTML interaction
  - The script `upgrade_htmx.sh` is available to make upgrading easier
- [air](https://github.com/air-verse/air) for live reloading of the application.
- [golang migrate](https://github.com/golang-migrate/migrate) for DB migrations (build/dev-time tool via `./migrate.sh`; not a runtime dependency of the server binary).
- [playwright-go](https://github.com/playwright-community/playwright-go) for E2E testing.

## Security

This project comes with comprehensive security features enabled by default.

### CSRF Protection

The application uses Go 1.25+'s native `http.CrossOriginProtection` for CSRF defense. This provides transparent protection without requiring CSRF tokens in your forms.

- **How it works:** Validates requests using the `Sec-Fetch-Site` header
- **What's protected:** POST, PUT, DELETE, PATCH requests
- **Configuration:** Enabled by default in the middleware chain
- **Implementation:** See `internal/server/middleware/csrf.go`

### Rate Limiting

Per-IP rate limiting prevents abuse and helps protect against DoS attacks.

- **Default limit:** 50 requests per minute per IP address
- **Algorithm:** Token bucket with in-memory storage
- **Auto-cleanup:** Inactive IP limiters are cleaned up every 10 minutes
- **Configuration:** Set `RATE_LIMIT` environment variable to customize
- **Headers:** Returns `Retry-After: 60` when limit exceeded
- **Implementation:** See `internal/server/middleware/ratelimit.go`

### Security Headers

The following security headers are automatically set on all responses:

- `X-Frame-Options: DENY` - Prevents clickjacking attacks
- `X-Content-Type-Options: nosniff` - Prevents MIME-type sniffing
- `Referrer-Policy: strict-origin-when-cross-origin` - Controls referrer information
- `Permissions-Policy: geolocation=(), microphone=(), camera=()` - Restricts browser features
- `Strict-Transport-Security` - Enforces HTTPS (when using TLS)

**Implementation:** See `internal/server/middleware/security.go`

### Server Hardening

The HTTP server is configured with timeouts to prevent slowloris and similar attacks:

- `ReadHeaderTimeout: 5s` - Maximum time to read request headers
- `ReadTimeout: 15s` - Maximum time to read entire request
- `WriteTimeout: 15s` - Maximum time to write response
- `IdleTimeout: 60s` - Maximum idle connection time
- `MaxHeaderBytes: 1MB` - Maximum header size

**Implementation:** See `internal/server/server.go`

### Panic Recovery

A recovery middleware catches panics and logs them with full stack traces, preventing the server from crashing while providing debugging information.

**Implementation:** See `internal/server/middleware/recovery.go`

### Testing

Comprehensive security tests are included in `e2e/security_test.go`:
- Security headers validation
- CSRF protection enforcement
- Rate limiting behavior
- Server timeout configurations

## Github Workflow

The repository comes with two GitHub workflows: `ci.yml` lints and tests the code, and `release.yml` tags, creates a GitHub Release, runs [GoReleaser](https://goreleaser.com/) to build and attach binaries, and publishes the Docker image.

See `AGENTS.md` for the full project structure and code-style conventions.

## Deploy Railway
[![Deploy on Railway](https://railway.com/button.svg)](https://railway.com/deploy/KnYa4C?referralCode=HtBp41&utm_medium=integration&utm_source=template&utm_campaign=generic)

> **Note:** the one-click template above may predate the migration entrypoint / Turso / env vars described below — if the deploy comes up without `DB_URL`/`TURSO_AUTH_TOKEN` set, follow Manual setup instead.

Railway's container filesystem is ephemeral — every deploy/restart is a fresh container. A local SQLite file (`DB_URL=./db.sqlite3`) gets wiped and re-migrated from scratch each time unless backed by a Volume, and **Railway's Trial plan does not support Volumes**. Recommended fix: use [Turso](https://turso.tech) (free-tier hosted libSQL) as the database instead — persists regardless of plan, no Volume needed.

### Manual setup

1. **New Project → Deploy from GitHub repo**, select this repo. Railway auto-detects the root `Dockerfile` — no Buildpack/Nixpacks config needed.
2. **Root Directory**: leave blank/default. `Dockerfile` is at the repo root, not a subfolder.
3. **Custom Start Command**: leave blank. The deploy image is `gcr.io/distroless/static-debian13` (no shell), so a start command that Railway would normally run via `/bin/sh -c` will crash the container. The `Dockerfile`'s `CMD ["/entrypoint"]` handles startup — it applies pending DB migrations, then execs into the server binary.
4. **Database — Turso**:
   ```shell
   turso auth login
   turso db create <name>
   turso db show <name> --url        # -> libsql://<db>-<org>.turso.io
   turso db tokens create <name>
   ```
5. **Variables**:
   ```
   DB_URL=libsql://<db>-<org>.turso.io
   TURSO_AUTH_TOKEN=<token from `turso db tokens create`>
   LOG_LEVEL=info
   LOG_OUTPUT=json
   RATE_LIMIT=50
   ```
   Do not set `PORT` — Railway injects it and the app already reads it.
6. **Build Args** (optional but recommended): set `VERSION` to a git SHA or release tag. `internal/version.Value` is compared against `"dev"` to decide whether the router trusts proxy headers for client IP — Railway sits behind a proxy, so this should not be left as `dev` in production.

**Prefer local SQLite instead?** Requires a paid Railway plan (Trial has no Volumes): attach a Volume at `/data`, set `DB_URL=/data/db.sqlite3`.

### Limits

Local SQLite uses a file lock — single replica only, no horizontal scaling. Turso removes that constraint (it's a remote DB, no local file lock) — multi-replica deploys work when `DB_URL` is a `libsql://` URL.
