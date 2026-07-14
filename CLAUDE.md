# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

> **Also read `AGENTS.md`** — it has the full command reference plus code-style conventions (imports, naming, error wrapping, templ/HTMX/SQLC/Tailwind syntax cheatsheets). This file focuses on architecture and things AGENTS.md doesn't cover.

## Commands

```shell
air                                                          # run with live reload (auto: templ/sqlc/tailwind generate + migrate on each change)
go run ./cmd/server                                          # run once, no reload (run migrate.sh first, see below)
go build -o ./tmp/main ./cmd/server                          # manual build

go test -v ./...                                             # unit tests
go test -race ./...                                          # unit tests w/ race detector
go test -v ./internal/server/handler -run TestHome           # single test
go test -v ./... -tags=e2e                                   # e2e (Playwright)
go test -v ./e2e -tags=e2e -run TestHomePage                 # single e2e test
HEADFUL=1 go test -v ./e2e -tags=e2e                         # e2e with visible browser
BROWSER=firefox go test -v ./e2e -tags=e2e                   # chromium (default), firefox, webkit

golangci-lint run                                            # lint
golangci-lint run --fix
go tool sqlc vet                                              # SQL lint
govulncheck ./...

go tool templ generate -path ./internal/components && go tool sqlc generate   # codegen (both)
go tool go-tw -i ./styles/input.css -o ./internal/dist/assets/css/output@dev.css  # tailwind

./migrate.sh -p sqlite -u ./db.sqlite3                        # apply DB migrations (server does NOT migrate on startup)
```

Codegen output is gitignored: `internal/components/**/*.go` and `internal/dist/assets/css` are generated, not committed. Edit the `.templ` sources and `internal/db/queries/query.sql`/migrations instead — never hand-edit generated `.go` files, they'll be overwritten.

## Architecture

Standard Go "server project" layout: `cmd/server/main.go` is the only entrypoint, everything else lives under `internal/` (blocks external imports) except `e2e/` (tests) and `styles/` (Tailwind input).

**Wiring flow** (read in this order to understand a request's path):
1. `cmd/server/main.go` — reads env vars, opens `db.Database`, builds `router.New(...)`, hands it to `server.New(..., server.WithRouter(...))`, calls `StartAndWait()` (graceful shutdown on SIGINT/SIGTERM, 10s grace period).
2. `internal/server/router/router.go` — registers routes on stdlib `http.ServeMux` (method-prefixed patterns like `"GET /{$}"`), then wraps the mux in a middleware chain via `middleware.Chain(...)`.
3. Middleware chain order matters: `Recovery` → `Logging` → `Security` → `RateLimit` → `CSRF`. Cache middleware applies only to `/assets/`. See `internal/server/middleware/`.
4. `internal/server/handler/` — struct-based handlers (`handler.New(logger, database)`) with DI'd logger + `db.Database`; render templ components from `internal/components/`.

**Database**: `db.Database` (`internal/db/db.go`) is a small interface (`DB()`, `Queries()`, `Close()`) wrapping a sqlc-generated `*queries.Queries` — implemented against SQLite but designed to be swapped (e.g. for `sqlx`). Query source of truth is `internal/db/queries/query.sql` (`-- name: X :one/:many/:exec` annotations); schema/migrations live in `internal/db/migrations` and sqlc reads that directory to generate types. Migrations are never run by the server binary itself — always via `migrate.sh` (which `air`'s `pre_cmd` calls automatically in dev).

**Assets**: `internal/dist/assets` is `embed`-ed into the binary (`internal/dist/dist.go`). CSS under `assets/css` is Tailwind-CLI-generated from `styles/input.css` — never hand-edit generated CSS, edit `styles/input.css`.

**Security is on by default**, not opt-in — when touching routing/handlers, preserve these:
- CSRF via Go 1.25+ native `http.CrossOriginProtection` (`Sec-Fetch-Site` header, no tokens) — `internal/server/middleware/csrf.go`
- Per-IP token-bucket rate limiting, default 50 req/min, in-memory with auto-cleanup — `RATE_LIMIT` env var
- Security headers (X-Frame-Options, CSP, HSTS, etc.) and server hardening (ReadHeaderTimeout/IdleTimeout/MaxHeaderBytes against slowloris) in `internal/server/server.go`

**Env vars**: `PORT` (8080), `LOG_LEVEL` (info), `LOG_OUTPUT` (text), `DB_URL` (./db.sqlite3), `RATE_LIMIT` (50). Copy `.env.example` to `.env` for local dev.

**Version**: `internal/version` is set via ldflags at build/release time (GoReleaser), defaults to `dev`; router uses `version.Value != "dev"` to decide whether to trust proxy headers for client IP.
