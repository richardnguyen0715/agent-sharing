# Task: T08 — Docker & Production Readiness

## Parent Plan
001-termhub-mvp.md

## Goal
Containerize the application and add production-readiness features.

## Steps
1. Backend Dockerfile (multi-stage: golang:1.24 builder → ubuntu:24.04 runtime)
2. Frontend Dockerfile (node builder → nginx:alpine for static serving)
3. Full docker-compose.yml (frontend, backend, postgres, redis)
4. Nginx config for frontend (SPA routing, WebSocket proxy)
5. Security headers (X-Frame-Options, X-Content-Type-Options, CSP, Referrer-Policy)
6. CORS configuration (trusted origins only)
7. Health check endpoints (GET /api/health, GET /api/health/ready)
8. Graceful shutdown handling
9. README with complete setup instructions

## Acceptance Criteria
- [ ] `docker compose up` starts entire stack
- [ ] Frontend accessible at localhost
- [ ] WebSocket connections work through proxy
- [ ] Health checks respond correctly
- [ ] README covers setup, env vars, and development workflow

## Status
TODO
