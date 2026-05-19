# Task: T01 — Project Scaffolding & Infrastructure

## Parent Plan
001-termhub-mvp.md

## Goal
Set up the monorepo structure inside `termhub-app/` with Go backend, React frontend, Docker Compose, and database migrations.

## Steps
1. Initialize Go module (`termhub-app/backend/`)
2. Create layered directory structure (handlers, services, repositories, domain, middleware, config, websocket)
3. Initialize Vite + React + TypeScript project (`termhub-app/frontend/`)
4. Install TailwindCSS 4, Zustand, React Router, xterm.js, Tabler Icons
5. Create docker-compose.yml with PostgreSQL 16 + Redis 7
6. Write database migrations (users, terminals, terminal_sessions, audit_logs)
7. Create .env.example with all required variables
8. Create root README.md

## Acceptance Criteria
- [ ] `go build ./...` succeeds in backend/
- [ ] `npm run dev` starts frontend
- [ ] `docker compose up -d` starts PostgreSQL and Redis
- [ ] Migrations run successfully
- [ ] Directory structure follows layered architecture

## Status
TODO
