# TermHub — Web-Based Multi-User Terminal Platform

A web-based terminal platform enabling multiple users to access isolated terminal sessions on a shared authenticated server.

## Architecture

| Layer | Technology |
|---|---|
| Frontend | React 19 + Vite + xterm.js + Zustand + TailwindCSS |
| Backend | Go + chi + gorilla/websocket + creack/pty |
| Auth | JWT (access + refresh tokens) |
| Database | PostgreSQL 16 |
| Cache | Redis 7 |
| Proxy | Nginx (production) |
| Containers | Docker + Docker Compose |

```
┌─────────────┐     ┌──────────────┐     ┌──────────────┐
│   Browser    │────▶│  Nginx (:80) │────▶│ Backend(:8080)│
│              │     │  Frontend    │     │  Go API + WS │
└─────────────┘     └──────────────┘     └──────┬───────┘
                                                │
                                    ┌───────────┴───────────┐
                                    │                       │
                              ┌─────▼─────┐          ┌─────▼─────┐
                              │ PostgreSQL │          │   Redis   │
                              │   (:5432)  │          │  (:6379)  │
                              └───────────┘          └───────────┘
```

## Quick Start — Docker (Production-like)

Build and run the entire stack with a single command:

```bash
cd termhub-app
docker compose up --build -d
```

Open http://localhost:3000 and log in with:

- **Username:** `admin`
- **Password:** `admin123`

To stop:

```bash
docker compose down
```

To stop and remove all data (volumes):

```bash
docker compose down -v
```

## Development Setup

For local development, run only the infrastructure (PostgreSQL + Redis) in Docker and run the backend/frontend natively for hot-reload:

### 1. Start infrastructure

```bash
cd termhub-app
docker compose -f docker-compose.dev.yml up -d
```

### 2. Run database migrations

```bash
cd backend
go run cmd/migrate/main.go up
```

### 3. Start backend

```bash
cd backend
cp .env.example .env
go run cmd/server/main.go
```

Backend runs on http://localhost:8080

### 4. Start frontend

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Frontend runs on http://localhost:5173

## Project Structure

```
termhub-app/
├── backend/
│   ├── cmd/server/        # Application entrypoint
│   ├── internal/          # Business logic
│   ├── migrations/        # SQL migration files
│   ├── Dockerfile         # Multi-stage Go build
│   └── .dockerignore
├── frontend/
│   ├── src/               # React application source
│   ├── Dockerfile         # Multi-stage Node build + Nginx
│   ├── nginx.conf         # Nginx config (SPA + API proxy + WebSocket)
│   └── .dockerignore
├── docker-compose.yml     # Full stack (production-like)
├── docker-compose.dev.yml # Infrastructure only (development)
├── .env.example           # Environment variable reference
└── README.md
```

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `POSTGRES_PASSWORD` | `termhub` | PostgreSQL password |
| `JWT_SECRET` | `change-me-in-production-...` | JWT signing secret (**change in production**) |
| `CORS_ORIGINS` | `http://localhost:3000` | Allowed CORS origins |
| `LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `JWT_ACCESS_TTL` | `15m` | Access token lifetime |
| `JWT_REFRESH_TTL` | `168h` | Refresh token lifetime |

Copy `.env.example` to `.env` and customize as needed. See also `backend/.env.example` and `frontend/.env.example` for service-specific configuration.

## Access URLs

| Mode | Frontend | Backend API |
|---|---|---|
| Docker (production) | http://localhost:3000 | http://localhost:3000/api/ (via Nginx proxy) |
| Development | http://localhost:5173 | http://localhost:8080 |

## Default Credentials

| Username | Password | Role |
|---|---|---|
| `admin` | `admin123` | Administrator |

> ⚠️ Change the default admin password and `JWT_SECRET` before deploying to production.

## Features

- **Multi-user terminal access** — each user gets isolated PTY sessions
- **Real-time terminal streaming** — WebSocket-based stdin/stdout
- **RBAC** — admin, operator, viewer, restricted roles
- **Session recording** — asciinema-style capture and playback
- **Admin dashboard** — live metrics, terminal monitoring, audit logs
- **Command governance** — heuristic blocking of dangerous commands
- **Session recovery** — reconnect to running PTY after disconnect

## Docker Details

### Backend Container
- **Build:** Multi-stage — Go 1.24 Alpine (build) → Ubuntu 24.04 (runtime)
- **Why Ubuntu:** PTY/terminal sessions require proper shell support; Alpine lacks full bash/coreutils
- **Health check:** `GET /api/health` via wget

### Frontend Container
- **Build:** Multi-stage — Node 22 Alpine (build) → Nginx Alpine (serve)
- **Nginx:** Serves SPA with HTML5 history fallback, proxies `/api/` to backend with WebSocket upgrade support
- **Security headers:** X-Frame-Options, X-Content-Type-Options, Referrer-Policy, X-XSS-Protection
- **Caching:** Static assets cached for 1 year with immutable flag

### Volumes
| Volume | Purpose |
|---|---|
| `termhub-pgdata` | PostgreSQL data persistence |
| `termhub-redis` | Redis data persistence |
| `termhub-recordings` | Terminal session recordings |
