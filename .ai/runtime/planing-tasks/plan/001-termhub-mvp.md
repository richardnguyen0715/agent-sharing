# Plan: TermHub MVP — Web-Based Multi-User Terminal Platform

## Summary

Build the TermHub application — a web-based multi-user terminal platform with Go backend (PTY management, JWT auth, RBAC, WebSocket streaming) and React frontend (xterm.js terminal rendering, Zustand state, TailwindCSS UI) inside `termhub-app/`.

## Context

- Spec: `.ai/context/spec/web-based-terminal.md` — full system spec with 28 sections
- UI designs: `.ai/context/spec/design/` — 3 HTML mockups (main page, new terminal modal, admin dashboard)
- UI/UX doc: `.ai/context/spec/ui-ux-design.md` — references design files + describes pages
- Target directory: `termhub-app/` (currently empty)
- No existing checkpoints for this app
- Policies: coding-policies.md, how-to-structure.md, testing-debugging-frameworks.md

## Tech Stack / Approach

### Backend (Go)

| Layer | Technology |
|---|---|
| Language | Go 1.24+ |
| Router | chi |
| WebSocket | gorilla/websocket |
| PTY | creack/pty |
| Auth | JWT (golang-jwt) |
| Database | PostgreSQL (pgx) |
| Cache | Redis (go-redis) |
| Password | bcrypt |
| Migration | golang-migrate |

### Frontend (React)

| Layer | Technology |
|---|---|
| Framework | React 19 |
| Build Tool | Vite |
| Terminal UI | xterm.js + xterm-addon-fit + xterm-addon-webgl |
| State | Zustand |
| Styling | TailwindCSS 4 |
| Icons | Tabler Icons |
| HTTP | fetch (native) |
| WebSocket | Native WebSocket |
| Router | React Router v7 |

### Infrastructure

| Component | Tool |
|---|---|
| Containers | Docker + docker-compose |
| DB | PostgreSQL 16 |
| Cache | Redis 7 |
| Reverse proxy | Nginx (future) |

## Phases

### Phase 1 — Project Scaffolding & Infrastructure
- Backend Go module init, directory structure (layered architecture)
- Frontend Vite + React + TailwindCSS + Zustand scaffold
- Docker Compose for PostgreSQL + Redis
- Database migrations (users, terminals, terminal_sessions, audit_logs)
- Environment config (.env, validation)

### Phase 2 — Authentication & Authorization
- User model + repository (PostgreSQL)
- JWT token generation/validation (access + refresh)
- Login/register API endpoints
- Auth middleware
- RBAC system (admin, operator, viewer, restricted roles)
- Login page (frontend)

### Phase 3 — Terminal Management (Backend Core)
- PTY manager (spawn, kill, resize, reconnect)
- Terminal CRUD API (create, list, get, delete)
- Terminal ownership enforcement
- Terminal registry (in-memory map + DB persistence)
- WebSocket gateway (stdin/stdout streaming, resize messages)
- WebSocket auth (JWT validation on upgrade)

### Phase 4 — Terminal UI (Frontend Core)
- Main page layout (3-column: sidebar + terminal + info panel)
- xterm.js integration with WebSocket backend
- Tab system (open/close/switch terminals)
- Sidebar (owned terminals list, shared section, host credentials, actions)
- Right panel (session info, processes, activity log, quick actions)
- Status bar (connection state, PTY ID, session duration, dimensions)
- Terminal toolbar (copy, record, maximize)

### Phase 5 — New Terminal Modal & Terminal Lifecycle
- New terminal modal (name, CWD, shell, color tag, permissions)
- Terminal creation flow (frontend → API → PTY spawn → WebSocket attach)
- Terminal deletion flow
- Resize handling (xterm-addon-fit + WebSocket resize message)
- Session recovery (reconnect WebSocket to running PTY)

### Phase 6 — Admin Dashboard
- Admin dashboard page with metrics cards
- All-terminals table (filterable, searchable)
- Audit log sidebar (real-time events)
- Active warnings display
- User management (add user, view users)
- Admin-only route protection

### Phase 7 — Audit, Policy Engine & Recording
- Audit logging (login, terminal create/delete, connect/disconnect, commands)
- Command governance (heuristic blocks: rm -rf /, shutdown, reboot, mkfs)
- Warning system
- Session recording (asciinema-style format)
- Recording playback page

### Phase 8 — Docker & Production Readiness
- Backend Dockerfile (multi-stage build)
- Frontend Dockerfile (nginx + static)
- Full docker-compose.yml (all services)
- Security headers (CSP, X-Frame-Options, etc.)
- CORS configuration
- Health check endpoints
- README with setup instructions

## Open Questions

1. Should we use SQLite for MVP simplicity or PostgreSQL as spec requires? → **Using PostgreSQL per spec**
2. Redis needed for MVP or can we defer? → **Include Redis in docker-compose but defer usage to Phase 2+ (JWT blacklist, session cache)**
3. Keycloak for auth or simple JWT? → **Simple JWT for MVP, Keycloak is future extension**
4. Kafka integration? → **Deferred to post-MVP per spec Phase 2/3 scalability roadmap**

## Approval Status

APPROVED (2026-05-20)
