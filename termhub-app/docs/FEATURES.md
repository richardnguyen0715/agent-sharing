# TermHub — Feature Documentation

TermHub is a web-based, multi-user terminal platform. Each user gets isolated PTY
(shell) sessions on a shared server, streamed live to the browser over WebSockets,
with role-based access control, command governance, session recording, audit
logging, and an admin dashboard.

> The real application lives in `termhub-app/`. The repo-root `README.md` is an
> unrelated template and should be ignored.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | React 19, Vite 6, TypeScript 5.8, TailwindCSS 4, Zustand 5, react-router-dom 7, xterm.js (`@xterm/xterm` + fit/webgl addons), Tabler icons |
| Backend | Go 1.21, chi/v5 router, gorilla/websocket, creack/pty, golang-jwt/v5, pgx/v5, zerolog, godotenv, bcrypt |
| Database | PostgreSQL 16 (pgx pool) |
| Cache | Redis 7 — provisioned in compose/config but **not used** by backend code yet |
| Proxy | Nginx (SPA serving + reverse proxy for `/api` and WS upgrade) |
| Deploy | Docker, Docker Compose, Tailscale Funnel |

**Directory layout:**
- Backend (clean layering): `cmd/server` → `internal/{config,domain,handler,middleware,migrate,policy,pty,recording,repository,service,ws}`; SQL in `backend/migrations/` (embedded copy in `internal/migrate/sql/`).
- Frontend: `src/{components/{layout,terminal,modals,ui},hooks,pages,services,stores,types}`.

---

## 1. Authentication & Authorization

- **JWT login** — bcrypt verify (cost 12), HS256 access + refresh tokens with
  user_id/username/role/type claims. `POST /api/auth/login`.
- **Registration with auto-login** — validates username (3–32 chars), password
  (≥8), role; issues tokens immediately. `POST /api/auth/register`.
- **Token refresh** — validates refresh token, re-fetches user (catches deleted
  users / role changes), issues new pair. `POST /api/auth/refresh`.
- **Current user** — `GET /api/auth/me`.
- **JWT auth middleware** — Bearer parsing, signing-method check, rejects refresh
  tokens used as access tokens, injects claims into context.
- **RBAC middleware** — `RequireRole(...)` set-based gate; admin routes require `admin`.
- **Four roles** — `admin`, `operator`, `viewer`, `restricted` (DB CHECK constraint).
- **Default admin seeding** — creates `admin`/`admin123` on first boot if no users exist.
- **Frontend route guards** — `ProtectedRoute` (auth) and `AdminRoute` (admin-only).
- **Auto token-refresh on 401** — API client retries once after refresh, force-logout on failure.

## 2. Terminal Management

- **Create terminal** — persists DB record then spawns PTY; rolls back DB row if
  spawn fails. Defaults: title=`Terminal`, cwd=`/tmp`, shell=`/bin/bash`,
  colorTag=`#00FF00`. `POST /api/terminals`.
- **List own terminals** — owner-scoped. `GET /api/terminals`.
- **Get terminal by ID** — ownership-checked. `GET /api/terminals/{id}`.
- **Delete / kill** — audit-logs before delete, kills PTY, removes DB row. `DELETE /api/terminals/{id}`.
- **Resize** — ownership-checked, validates cols/rows>0. `POST /api/terminals/{id}/resize`.
- **Model** — id, owner, title, cwd, shell, status (active/idle/terminated), colorTag, timestamps.
- **Frontend tab management** — add/remove/activate, dedup, neighbor-selection on close (Zustand store).

## 3. PTY / Terminal Engine

- **PTY lifecycle manager** — create/get/resize/kill/cleanup, in-memory session map, active count.
- **Shell + CWD resolution** — `$SHELL` or `/bin/bash`; expands `~`, validates dir,
  falls back to `$HOME`. Sets `TERM=xterm-256color`, default 24×80.
- **Graceful kill** — close master (SIGHUP) → wait 2s → SIGKILL. (creack/pty handles
  setsid internally; `Setpgid`/`Setsid` deliberately unset for Docker/macOS.)
- **Graceful shutdown** — concurrent PTY + recording cleanup, 3s timeout.

## 4. WebSocket / Real-time Streaming

- **WS terminal endpoint** — `GET /api/ws/terminals/{id}`, auth via `?token=` query
  param, rejects refresh tokens, verifies ownership.
- **Session recovery** — if backend restarted (PTY gone but DB record exists),
  transparently re-spawns PTY on connect.
- **WS Hub** — register/unregister, broadcast-to-terminal (multi-client fan-out), client-count metric.
- **Read/Write pumps** — ping/pong keepalive (60s/54s), 8KB read limit, JSON
  `stdin`/`resize` messages + raw-input fallback.
- **Output broadcast** — streams PTY output to all clients; emits `error` message on process exit.
- **Frontend xterm.js** — GitHub-dark theme, JetBrains Mono, 5000 scrollback,
  FitAddon, ResizeObserver, keyboard→WS, resize→WS. WS URL derived from page location.

## 5. Command Governance (Policy Engine)

- **Heuristic engine** — plain substring matching (explicitly *not* claimed as real
  security). Actions: block / warn / allow.
- **Block rules** — `rm -rf /`, `shutdown`, `reboot`, `mkfs`, `dd if=`, fork bomb.
- **Warn rules** — `chmod 777`, `rm -rf`, `kill -9 1`, `sudo su`.
- **Enforcement at WS input** — blocked commands return error (not written to PTY);
  warns write a `policy.violation` audit log.
- **Runtime rule addition** — `AddRule` exists, not yet exposed via API.

## 6. Session Recording

- **Recorder manager** — start/stop/list/get; JSONL append-only files under
  `recordings/{terminalId}/{timestamp}.jsonl`; flush every 100 events; path-traversal guard.
- **Control API** — `POST /api/terminals/{id}/recording/start` and `/stop`.
- **Browse/playback API** — `GET /api/recordings` (list),
  `GET /api/recordings/{terminalId}/{filename}` (events).
- **Capture** — `RecordOutput`/`RecordInput` wired into WS loop (best-effort).
- ⚠️ **Frontend Recordings page is a placeholder** — shows hardcoded data, does not
  call the (functional) recording API yet.

## 7. Audit Logging

- **Audit service** — inserts events with JSONB metadata; paginated list (DESC) + list-by-user.
- **Actions** — login, logout, terminal create/delete/connect/disconnect, user
  create/update/login/register, policy.violation.
- **Emitted on** — login, register, terminal create/delete, WS connect/disconnect, policy violation.
- **Admin listing** — `GET /api/admin/audit?limit=&offset=`.

## 8. Admin Dashboard

- **List all users** — paginated. `GET /api/admin/users`.
- **List all terminals (cross-user)** — `GET /api/admin/terminals`.
- **System metrics** — totalUsers, totalTerminals, activeTerminals, activeConnections.
  `GET /api/admin/metrics`.
- **Admin UI** — live dashboard (10s auto-refresh), metric cards, terminal table
  (status/owner/role/CWD/uptime), search + filter pills, audit log + active-warnings panels.
- ⚠️ **UI-only controls** — "Add user", "Export logs", per-row View/Stop have no handlers.

## 9. Database & Migrations

- **Auto-migrations on startup** — embedded SQL runner with `schema_migrations`
  tracking; applies `*.up.sql` in version order.
- **Schema** — `users`, `terminals` (FK→users, ON DELETE CASCADE),
  `terminal_sessions` (FK→terminals), `audit_logs` (FK SET NULL), with indexes and matching `.down.sql`.
- **Repositories** — pgx-based CRUD, count/count-by-status (`repository/user.go`, `repository/terminal.go`).

## 10. Configuration / Env Vars

Required: `DATABASE_URL`, `JWT_SECRET`.

| Var | Default |
|---|---|
| `PORT` | 8080 |
| `JWT_ACCESS_TTL` | 15m |
| `JWT_REFRESH_TTL` | 168h |
| `REDIS_URL` | (loaded, unused) |
| `CORS_ORIGINS` | — |
| `LOG_LEVEL` | — |
| `VITE_API_URL` / `VITE_WS_URL` | (frontend) |
| `POSTGRES_PASSWORD` / `FRONTEND_PORT` | — / 3000 |

- **CORS middleware** — origin allowlist from config.
- **Request middleware** — RequestID, RealIP, Recoverer, request logging.

## 11. Deployment

- **Docker Compose stack** — frontend (Nginx), backend, postgres, redis; only
  frontend port exposed in prod; named volumes for pgdata/redis/recordings.
- **Multi-stage Dockerfiles** — backend Go-Alpine→Debian-slim (real shell for PTY,
  non-root, healthcheck `/api/health`); frontend Node-Alpine→Nginx-Alpine with
  `VITE_*` build args + security headers.
- **Tailscale Funnel one-command deploy** — installs/auths Tailscale, auto-generates
  `JWT_SECRET`/`POSTGRES_PASSWORD`, sets CORS to tailnet FQDN, builds stack,
  configures `tailscale serve`+`funnel` for public HTTPS. Teardown script included.
- **Makefile workflow** — dev (`dev-infra/dev-backend/dev-frontend`), prod
  (`up/down/logs/status`), funnel (`deploy/teardown/funnel-*`), quality
  (`test-backend/lint-backend/typecheck`).
- **Health check** — `GET /api/health` (pings DB).

---

## API Surface

```
POST   /api/auth/login
POST   /api/auth/register
POST   /api/auth/refresh
GET    /api/auth/me

GET    /api/terminals
POST   /api/terminals
GET    /api/terminals/{id}
DELETE /api/terminals/{id}
POST   /api/terminals/{id}/resize
POST   /api/terminals/{id}/recording/start
POST   /api/terminals/{id}/recording/stop

GET    /api/ws/terminals/{id}          (WebSocket, ?token=)

GET    /api/recordings
GET    /api/recordings/{terminalId}/{filename}

GET    /api/admin/users
GET    /api/admin/terminals
GET    /api/admin/audit
GET    /api/admin/metrics

GET    /api/health
```

---

## Known Gaps / Caveats

1. **Redis provisioned but unused** — `REDIS_URL` loaded into config, never read.
2. **Recordings page is placeholder UI** — backend API works; frontend shows hardcoded data.
3. **Admin write controls non-functional** — Add user, Export logs, View/Stop have no
   handlers; no user-management write API (no create/update/delete user endpoints).
4. **No `restricted`-role command limiting enforced** — policy engine applies the same
   global rules to everyone regardless of role.
5. **Policy engine is heuristic substring matching** by design (not real security).
6. **WS `CheckOrigin` allows all origins** — fine for dev, harden for production.
7. **`terminal_sessions` table unused** — sessions tracked only in-memory by the WS hub.
