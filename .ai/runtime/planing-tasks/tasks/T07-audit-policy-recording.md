# Task: T07 — Audit, Policy Engine & Recording

## Parent Plan
001-termhub-mvp.md

## Goal
Implement audit logging, command governance (heuristic blocks), and session recording.

## Steps
1. Audit service (log events to PostgreSQL audit_logs table)
2. Audit events: login, terminal create/delete, websocket connect/disconnect, command execution
3. Policy engine: command pattern matching (rm -rf /, shutdown, reboot, mkfs)
4. Warning system (soft governance: block/warn based on pattern)
5. Session recording: capture PTY output with timestamps (asciinema v2 format)
6. Recording storage (filesystem, organized by terminal_id/timestamp)
7. Recording API: GET /api/recordings, GET /api/recordings/:id
8. Recording playback page (frontend) — basic asciinema-style playback

## Acceptance Criteria
- [ ] All critical events are logged to audit_logs
- [ ] Dangerous commands are blocked/warned
- [ ] Sessions can be recorded
- [ ] Recordings can be replayed

## Status
TODO
