# Task: T03 — Terminal Management Backend

## Parent Plan
001-termhub-mvp.md

## Goal
Implement PTY management, terminal CRUD API, ownership enforcement, and WebSocket gateway.

## Steps
1. PTY manager: spawn bash/zsh/sh with creack/pty
2. Terminal struct (ID, OwnerID, PTY, Cmd, CreatedAt, Title, CWD, Status)
3. Terminal registry (concurrent-safe map)
4. Terminal service (create, list by owner, get, delete, resize)
5. Terminal repository (PostgreSQL CRUD)
6. API endpoints: POST /api/terminals, GET /api/terminals, GET /api/terminals/:id, DELETE /api/terminals/:id
7. Ownership middleware (terminal.owner_id == user.id)
8. WebSocket endpoint: GET /api/ws/:terminal_id
9. WebSocket auth (JWT token from query param)
10. WebSocket message routing (stdin, resize, output)
11. PTY I/O streaming goroutines (PTY → WS, WS → PTY)
12. Terminal cleanup on process exit

## Acceptance Criteria
- [ ] Can create a PTY session via API
- [ ] WebSocket streams terminal I/O bidirectionally
- [ ] Users can only access their own terminals
- [ ] Terminal resize works
- [ ] PTY processes are properly cleaned up

## Status
TODO
