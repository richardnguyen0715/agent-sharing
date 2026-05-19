# Task: T06 — Admin Dashboard

## Parent Plan
001-termhub-mvp.md

## Goal
Build the admin dashboard page matching the admin_dashboard.html design with live metrics, terminal table, and audit log.

## Steps
1. Admin dashboard page component
2. Metrics cards (active terminals, users online, warnings today, session recording %)
3. All-terminals table (filterable: all/connected/idle/warnings, searchable)
4. Table columns: status dot, terminal name, owner, role, status, CWD, uptime, actions
5. Row actions (view terminal, stop terminal)
6. Right panel: active warnings + audit log (real-time)
7. Admin API endpoints: GET /api/admin/terminals, GET /api/admin/users, GET /api/admin/metrics, GET /api/admin/audit
8. Admin-only route protection (frontend + backend)
9. Add user functionality

## Acceptance Criteria
- [ ] Dashboard matches admin_dashboard.html design
- [ ] Metrics show real-time data
- [ ] Terminal table filters and search work
- [ ] Audit log shows recent events
- [ ] Only admin role can access

## Status
TODO
