# TermHub — Web-Based Multi-User Terminal Platform

A web-based terminal platform enabling multiple users to access isolated terminal sessions on a shared server with role-based access control, real-time WebSocket streaming, session recording, and an admin dashboard.

---

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Quick Start (Docker)](#quick-start--docker)
- [Development Setup](#development-setup)
- [Production Deployment (Linux Server)](#production-deployment--linux-server)
- [Environment Variables](#environment-variables)
- [Project Structure](#project-structure)
- [Makefile Commands](#makefile-commands)
- [Usage Guide](#usage-guide)
- [Troubleshooting](#troubleshooting)
- [Docker Details](#docker-details)

---

## Features

- **Multi-user terminal access** — each user gets isolated PTY sessions
- **Real-time terminal streaming** — WebSocket-based stdin/stdout via xterm.js
- **RBAC** — admin, operator, viewer, restricted roles
- **Session recording** — asciinema-style capture and playback
- **Admin dashboard** — live metrics, terminal monitoring, audit logs
- **Command governance** — heuristic blocking of dangerous commands
- **Session recovery** — reconnect to running PTY after page refresh or disconnect
- **Kill terminal** — terminate running sessions from the UI toolbar
- **Auto-migrations** — database schema applied automatically on startup
- **One-command deploy** — single `docker compose up` for the entire stack

---

## Architecture

| Layer | Technology |
|---|---|
| Frontend | React 19 + Vite + xterm.js + Zustand + TailwindCSS |
| Backend | Go 1.21 + chi + gorilla/websocket + creack/pty |
| Auth | JWT (access + refresh tokens, bcrypt passwords) |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Proxy | Nginx (serves SPA + reverse-proxies API + WebSocket) |
| Containers | Docker + Docker Compose |

```
                    ┌──────────────────────────────────────────────┐
                    │              Docker Compose                   │
┌──────────┐       │  ┌──────────────┐     ┌──────────────────┐   │
│  Browser  │──────│──│  Nginx (:80) │────▶│  Backend (:8080) │   │
│           │ :3000│  │  + SPA files │     │  Go API + WS     │   │
└──────────┘       │  └──────────────┘     └────────┬─────────┘   │
                    │                               │              │
                    │                   ┌───────────┴──────────┐   │
                    │             ┌─────▼─────┐          ┌─────▼──┐│
                    │             │ PostgreSQL │          │ Redis  ││
                    │             │   (:5432)  │          │(:6379) ││
                    │             └───────────┘          └────────┘│
                    └──────────────────────────────────────────────┘
```

> **Note:** In production, only port 3000 (or your chosen `FRONTEND_PORT`) is exposed. PostgreSQL and Redis are internal-only.

---

## Quick Start — Docker

Build and run the entire stack locally with one command:

```bash
git clone <your-repo-url> termhub-app
cd termhub-app

# Create environment config
cp .env.example .env

# Build & start everything
docker compose up --build -d
```

Open **http://localhost:3000** and log in with:

| Username | Password | Role |
|---|---|---|
| `admin` | `admin123` | Administrator |

To stop:

```bash
docker compose down        # Stop containers (data preserved)
docker compose down -v     # Stop + delete all data
```

---

## Development Setup

For local development with hot-reload, run infrastructure in Docker and backend/frontend natively.

### Prerequisites

- **Go 1.21+** — [golang.org/dl](https://go.dev/dl/)
- **Node.js 22+** — [nodejs.org](https://nodejs.org/)
- **Docker & Docker Compose** — [docs.docker.com](https://docs.docker.com/get-docker/)

### 1. Start infrastructure

```bash
cd termhub-app
make dev-infra
# or: docker compose -f docker-compose.dev.yml up -d
```

### 2. Start backend

```bash
cd backend
cp .env.example .env      # Only needed once
go run cmd/server/main.go
```

Backend runs on **http://localhost:8080**. Migrations run automatically.

### 3. Start frontend

```bash
cd frontend
cp .env.example .env      # Only needed once
npm install               # Only needed once
npm run dev
```

Frontend runs on **http://localhost:5173** with hot-reload. API calls proxy to `:8080`.

### Dev Access URLs

| Service | URL |
|---|---|
| Frontend (Vite) | http://localhost:5173 |
| Backend API | http://localhost:8080/api |
| PostgreSQL | localhost:5432 (user: `termhub`, pass: `termhub`, db: `termhub`) |
| Redis | localhost:6379 |

---

## Production Deployment — Linux Server

### Prerequisites

| Requirement | Minimum Version |
|---|---|
| Linux (Ubuntu/Debian/RHEL/etc.) | Any modern distro |
| Docker Engine | 24.0+ |
| Docker Compose (plugin) | 2.20+ |
| Git | 2.x |

### Step-by-step

#### 1. Install Docker (if not installed)

```bash
# Ubuntu / Debian
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
# Log out and back in for group changes to take effect
```

#### 2. Clone the repository

```bash
git clone <your-repo-url> /opt/termhub
cd /opt/termhub
```

#### 3. Configure environment

```bash
cp .env.example .env
```

Edit `.env` with production values:

```bash
# REQUIRED — change these!
POSTGRES_PASSWORD=<strong-random-password>
JWT_SECRET=<run: openssl rand -base64 32>

# Set to your server's domain or IP
CORS_ORIGINS=https://your-domain.com

# Optional
FRONTEND_PORT=3000
LOG_LEVEL=info
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
```

> **⚠️ Security:** Never use the default `JWT_SECRET` or `POSTGRES_PASSWORD` in production. Generate a strong secret:
> ```bash
> openssl rand -base64 32
> ```

#### 4. Build and start

```bash
docker compose up --build -d
```

This builds both images, starts all 4 services (frontend, backend, postgres, redis), runs database migrations automatically, and seeds the default admin user.

#### 5. Verify deployment

```bash
# Check all services are healthy
docker compose ps

# Expected output:
# NAME               STATUS              PORTS
# termhub-backend    Up (healthy)        8080/tcp
# termhub-frontend   Up (healthy)        0.0.0.0:3000->80/tcp
# termhub-postgres   Up (healthy)        5432/tcp
# termhub-redis      Up (healthy)        6379/tcp

# Test the health endpoint
curl http://localhost:3000/api/health
```

#### 6. Access the application

Open **http://your-server-ip:3000** in your browser.

Log in with `admin` / `admin123`, then **change the password immediately** via the profile page or create a new admin user.

### Putting behind a reverse proxy (recommended)

For HTTPS support, you have three options:

#### Option A: Tailscale Funnel (recommended — zero config HTTPS)

Tailscale Funnel gives you a public HTTPS URL with auto-provisioned certificates — no domain registration, no DNS config, no cert management.

**One-command deploy:**

```bash
# On your Linux server:
git clone <your-repo-url> /opt/termhub
cd /opt/termhub
sudo ./scripts/deploy-tailscale.sh
```

The script will:
1. Install Tailscale (if not installed)
2. Authenticate with your Tailscale account (opens a login URL)
3. Generate secure secrets (JWT, Postgres password)
4. Build and start the Docker stack
5. Configure Tailscale Serve + Funnel

Your app will be live at: `https://termhub.<your-tailnet>.ts.net`

**Manage Funnel:**

```bash
make funnel-status        # Check Funnel status
make funnel-off           # Disable public access (app keeps running)
make funnel-on            # Re-enable public access
make teardown             # Stop everything
```

**Requirements:**
- A free [Tailscale account](https://tailscale.com)
- Funnel must be enabled in your Tailscale admin console: [Access Controls → Funnel](https://login.tailscale.com/admin/acls/funnel)

#### Option B: Caddy (auto-TLS with a custom domain)

```bash
# Install Caddy: https://caddyserver.com/docs/install
sudo apt install -y caddy

# /etc/caddy/Caddyfile
your-domain.com {
    reverse_proxy localhost:3000
}

sudo systemctl restart caddy
```

#### Option C: Nginx with Let's Encrypt

```nginx
# /etc/nginx/sites-available/termhub
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate     /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket support
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_read_timeout 86400;
    }
}
```

When using Caddy or Nginx with HTTPS, update your `.env`:

```bash
CORS_ORIGINS=https://your-domain.com
```

### Updating to a new version

```bash
cd /opt/termhub
git pull
docker compose up --build -d
```

Migrations run automatically on backend startup — no manual steps needed.

### Backup & restore

```bash
# Backup PostgreSQL
docker compose exec postgres pg_dump -U termhub termhub > backup_$(date +%Y%m%d).sql

# Restore
cat backup_20260520.sql | docker compose exec -T postgres psql -U termhub termhub
```

---

## Environment Variables

### Root `.env` (used by `docker compose`)

| Variable | Default | Description |
|---|---|---|
| `POSTGRES_PASSWORD` | _(must set)_ | PostgreSQL password |
| `JWT_SECRET` | _(must set)_ | JWT signing secret |
| `CORS_ORIGINS` | `http://localhost:3000` | Allowed CORS origins |
| `FRONTEND_PORT` | `3000` | Host port for the frontend |
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `JWT_ACCESS_TTL` | `15m` | Access token lifetime |
| `JWT_REFRESH_TTL` | `168h` | Refresh token lifetime (7 days) |

### Backend `.env` (for local dev only)

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | Server listen port |
| `DATABASE_URL` | _(required)_ | PostgreSQL connection string |
| `REDIS_URL` | `redis://localhost:6379` | Redis connection string |
| `JWT_SECRET` | _(required)_ | JWT signing secret |
| `CORS_ORIGINS` | `http://localhost:5173` | Allowed CORS origins |
| `LOG_LEVEL` | `debug` | Log level |

### Frontend `.env` (for local dev only)

| Variable | Default | Description |
|---|---|---|
| `VITE_API_URL` | `http://localhost:8080/api` | Backend API base URL |
| `VITE_WS_URL` | `ws://localhost:8080/api/ws` | WebSocket base URL |

> In Docker production, the frontend uses relative paths (`/api`) — Nginx proxies to the backend container.

---

## Project Structure

```
termhub-app/
├── backend/
│   ├── cmd/server/            # Application entrypoint
│   ├── internal/
│   │   ├── config/            # Environment-based configuration
│   │   ├── domain/            # Domain models (User, Terminal, AuditLog)
│   │   ├── handler/           # HTTP & WebSocket handlers
│   │   ├── middleware/        # Auth, logging, CORS middleware
│   │   ├── migrate/           # Embedded SQL migration runner
│   │   ├── policy/            # Command governance engine
│   │   ├── pty/               # PTY session lifecycle manager
│   │   ├── recording/         # Session recording (asciinema format)
│   │   ├── repository/        # PostgreSQL data access
│   │   ├── service/           # Business logic layer
│   │   └── ws/                # WebSocket hub + client management
│   ├── migrations/            # Source SQL migration files
│   ├── Dockerfile             # Multi-stage Go → Debian slim
│   └── .env.example
├── frontend/
│   ├── src/
│   │   ├── components/        # React components (layout, terminal, admin)
│   │   ├── hooks/             # useTerminal (xterm.js + WebSocket)
│   │   ├── pages/             # Route pages
│   │   ├── services/          # API client + WebSocket service
│   │   ├── stores/            # Zustand state (auth, terminal, UI)
│   │   └── types/             # TypeScript type definitions
│   ├── nginx.conf             # Production Nginx config
│   ├── Dockerfile             # Multi-stage Node → Nginx Alpine
│   └── .env.example
├── docker-compose.yml         # Full production stack
├── docker-compose.dev.yml     # Dev infrastructure (PG + Redis)
├── Makefile                   # Common commands
├── .env.example               # Root env template
├── .gitignore
└── README.md
```

---

## Makefile Commands

```bash
make help             # Show all available commands

# Development
make dev              # Start dev infrastructure + instructions
make dev-infra        # Start Postgres + Redis in Docker
make dev-backend      # Run Go backend locally
make dev-frontend     # Run Vite dev server

# Production
make build            # Build Docker images
make up               # Start production stack
make down             # Stop production stack
make restart          # Restart everything
make logs             # Tail all logs
make logs-backend     # Tail backend logs only
make status           # Show container status

# Tailscale Funnel
make deploy           # Full deploy with Tailscale Funnel (sudo)
make teardown         # Stop TermHub + disable Funnel
make funnel-status    # Show Funnel status
make funnel-off       # Disable public access (app keeps running)
make funnel-on        # Re-enable public access

# Quality
make test-backend     # Run Go tests
make lint-backend     # Run go vet
make typecheck        # TypeScript type check
make build-backend    # Build Go binary locally
make build-frontend   # Build frontend locally

# Maintenance
make clean            # Remove all containers, volumes, images
make setup-env        # Create .env from template
```

---

## Usage Guide

### First Login

1. Open the application in your browser
2. Log in with `admin` / `admin123`
3. **Change your password** — go to the admin panel and update the default credentials

### Creating a Terminal

1. Click the **"+ New Terminal"** button in the left sidebar
2. A new terminal tab opens with a shell session
3. Type commands as you would in a native terminal
4. The terminal uses your server's default shell (`$SHELL`, typically `/bin/bash`)

### Managing Terminals

- **Switch tabs** — click terminal tabs at the top to switch between sessions
- **Kill a terminal** — click the **✕** button in the terminal toolbar to terminate the PTY and remove it
- **Session recovery** — if you refresh the page, existing terminals auto-reconnect

### Admin Dashboard

Accessible by users with the `admin` role:

- **Users** — view all registered users, their roles, and last activity
- **Terminals** — monitor all active terminal sessions across users
- **Audit Logs** — review all security-relevant events (login, terminal create/delete, etc.)
- **Metrics** — live system metrics (active users, terminals, sessions)

### User Roles

| Role | Terminals | Admin Panel | Manage Users |
|---|---|---|---|
| `admin` | ✅ Unlimited | ✅ Full access | ✅ |
| `operator` | ✅ Unlimited | ✅ View only | ❌ |
| `viewer` | ✅ Read-only | ❌ | ❌ |
| `restricted` | ✅ Limited commands | ❌ | ❌ |

### API Endpoints

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/auth/login` | Login (returns JWT tokens) |
| `POST` | `/api/auth/register` | Register a new user |
| `POST` | `/api/auth/refresh` | Refresh access token |
| `GET` | `/api/auth/me` | Get current user |
| `GET` | `/api/terminals` | List user's terminals |
| `POST` | `/api/terminals` | Create a new terminal |
| `DELETE` | `/api/terminals/:id` | Kill and delete a terminal |
| `POST` | `/api/terminals/:id/resize` | Resize terminal |
| `GET` | `/api/ws/terminals/:id` | WebSocket terminal stream |
| `GET` | `/api/health` | Health check |
| `GET` | `/api/admin/users` | Admin: list users |
| `GET` | `/api/admin/terminals` | Admin: list all terminals |
| `GET` | `/api/admin/audit` | Admin: audit logs |
| `GET` | `/api/admin/metrics` | Admin: system metrics |

---

## Troubleshooting

### Backend won't start — "DATABASE_URL is required"

The backend loads config from environment variables. For local dev, copy the example env:

```bash
cd backend
cp .env.example .env
```

For Docker, the `docker-compose.yml` sets all required vars automatically.

### "relation does not exist" SQL errors

Migrations run automatically on startup. If you see this, the backend likely crashed before migrations completed. Restart the backend:

```bash
docker compose restart backend
# or locally: go run cmd/server/main.go
```

### Terminal blank / can't type / "disconnected"

1. Check backend logs: `docker compose logs backend` or `make logs-backend`
2. Verify WebSocket connectivity — the browser console should show a WS connection to `/api/ws/terminals/:id`
3. If the backend was restarted, click the terminal tab — it auto-reconnects and re-spawns the PTY

### "fork/exec /bin/bash: operation not permitted"

This happens on macOS when PTY spawn uses Linux-specific `Setpgid`. The current code handles this automatically. If you see it in Docker, ensure the backend container has the necessary capabilities (the default Dockerfile uses `debian:12-slim` which supports PTY).

### Server won't shut down (Ctrl+C hangs)

Fixed in current version. The backend closes PTY file descriptors first (sends SIGHUP), then waits max 3 seconds for cleanup. If it still hangs, use:

```bash
# Find and kill the process
kill -9 $(lsof -ti:8080)
```

### Docker build fails

```bash
# Clean rebuild
docker compose down -v
docker compose build --no-cache
docker compose up -d
```

---

## Docker Details

### Backend Container
- **Build:** Multi-stage — Go 1.23 Alpine (build) → Debian 12 slim (runtime)
- **Why Debian:** PTY/terminal sessions require proper shell support (bash, coreutils)
- **Non-root:** Runs as `termhub` user
- **Health check:** `GET /api/health` via curl
- **Binary size:** ~15MB (stripped with `-ldflags="-s -w"`)

### Frontend Container
- **Build:** Multi-stage — Node 22 Alpine (build) → Nginx Alpine (serve)
- **Build args:** `VITE_API_URL` and `VITE_WS_URL` for compile-time configuration
- **Nginx:** Serves SPA with HTML5 history fallback, proxies `/api/` to backend with WebSocket upgrade support
- **Security headers:** X-Frame-Options, X-Content-Type-Options, Referrer-Policy, X-XSS-Protection
- **Caching:** Static assets cached for 1 year with immutable flag

### Volumes

| Volume | Purpose |
|---|---|
| `termhub-pgdata` | PostgreSQL data persistence |
| `termhub-redis` | Redis data persistence |
| `termhub-recordings` | Terminal session recordings |

### Network

In production (`docker-compose.yml`), only the frontend port is exposed to the host. PostgreSQL (5432) and Redis (6379) are accessible only within the Docker network — never exposed publicly.

---

## License

MIT
