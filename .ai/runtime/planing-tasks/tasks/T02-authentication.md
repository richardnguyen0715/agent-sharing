# Task: T02 — Authentication & Authorization

## Parent Plan
001-termhub-mvp.md

## Goal
Implement JWT-based authentication with RBAC, user management API, and login page.

## Steps
1. User domain model + PostgreSQL repository
2. Password hashing (bcrypt)
3. JWT service (access token 15m, refresh token 7d)
4. Auth middleware (JWT validation, user context injection)
5. RBAC middleware (role-based permission checks)
6. API endpoints: POST /api/auth/login, POST /api/auth/register, POST /api/auth/refresh, GET /api/auth/me
7. Login page (frontend) — styled to match design system
8. Auth state management (Zustand store)
9. Protected route wrapper
10. Seed admin user

## Acceptance Criteria
- [ ] Login/register flow works end-to-end
- [ ] JWT tokens are issued and validated
- [ ] RBAC blocks unauthorized access
- [ ] Login page matches design aesthetic
- [ ] Refresh token rotation works

## Status
TODO
