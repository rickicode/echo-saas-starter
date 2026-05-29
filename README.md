# Echo SaaS Starter

Go + Echo + TanStack SPA starter kit with **plugin modular architecture**.

## Features

- 🔐 **Auth** — PASETO v4 (bukan JWT) + argon2id password hashing
- 👥 **RBAC** — Predefined roles: Super Admin, Admin, User
- 💳 **Subscription** — Simple tier-based plans: Free / Pro / Enterprise
- 👤 **User Management** — CRUD, profile, avatar, search/filter
- 📝 **Docs CMS** — Markdown content with draft/publish, categories, tags
- 🏠 **Guest Pages** — Landing, pricing, docs, register/login
- 📊 **Admin Dashboard** — Stats, user/role/plan/content management
- 📱 **User Dashboard** — Profile, subscription, settings
- 🗄️ **Dual DB** — PostgreSQL (primary) + SQLite3 WAL mode
- 📦 **Single Binary** — go:embed, build once deploy anywhere
- 🧩 **Plugin Modular** — Add/remove features independently

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go + Echo v4 + go:embed |
| Frontend | React 19 + TanStack Router (file-based) |
| UI | shadcn/ui + Tailwind CSS |
| Auth | PASETO v4 local + argon2id |
| Database | PostgreSQL + SQLite3 WAL (auto-detect) |
| Bundler | Vite |
| Testing | go test + vitest |

## Quick Start

```bash
# Clone
git clone https://github.com/rickicode/echo-saas-starter.git
cd echo-saas-starter

# Setup
cp .env.example .env
# Edit .env — set DATABASE_URL and PASETO_SECRET

# Docker
make docker-run
# or
docker-compose up -d

# Open
open http://localhost:8080
```

## Plugin Architecture

Every feature is an independent plugin. Remove or add without breaking the core.

```
internal/plugins/
├── auth/          # PASETO authentication
├── roles/         # RBAC (Super Admin, Admin, User)
├── users/         # User management
├── plans/         # Subscription plans
└── docs/          # Docs/Blog CMS
```

### Remove a plugin

1. Delete `internal/plugins/<name>/`
2. Delete `ui/src/plugins/<name>/`
3. Remove registration from `internal/plugins/registry.go`

That's it. No other files need changes.

### Add a new plugin

1. Create `internal/plugins/<name>/` implementing the Plugin interface
2. Create `ui/src/plugins/<name>/` with frontend routes/components
3. Add migration files under `internal/plugins/<name>/migrations/`
4. Register in `internal/plugins/registry.go`

## Using the Prompt

This repo contains an AI agent prompt (`prompt.md`) that can build the entire starter kit autonomously.

```bash
# Download and read the prompt
curl -sL https://raw.githubusercontent.com/rickicode/echo-saas-starter/main/prompt.md
```

Feed it to any AI coding agent (Claude Code, Codex, etc.) and it will build everything.

## Project Structure

```
echo-saas-starter/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── core/                    # Plugin system, middleware, config
│   ├── database/                # DB drivers (PostgreSQL, SQLite)
│   ├── embed/                   # go:embed for SPA
│   ├── plugins/                 # Feature plugins
│   │   ├── auth/
│   │   ├── roles/
│   │   ├── users/
│   │   ├── plans/
│   │   └── docs/
│   └── server/                  # Echo server setup
├── ui/                          # Frontend (TanStack Router SPA)
│   └── src/
│       ├── routes/              # File-based routing
│       ├── plugins/             # Frontend plugin modules
│       ├── components/          # Shared components
│       └── lib/                 # API client, utils
├── prompt.md                    # AI agent prompt
├── Makefile
├── Dockerfile
└── docker-compose.yml
```

## Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| DATABASE_URL | ✅ | — | `postgres://...` or `sqlite://./data.db` |
| PASETO_SECRET | ✅ | — | Secret key for PASETO tokens |
| SERVER_PORT | ❌ | 8080 | HTTP server port |
| SERVER_HOST | ❌ | 0.0.0.0 | HTTP server bind address |
| LOG_LEVEL | ❌ | info | Log level (debug/info/warn/error) |
| CORS_ORIGINS | ❌ | * | Allowed CORS origins |

## License

MIT
