# Task: T05 — New Terminal Modal & Lifecycle

## Parent Plan
001-termhub-mvp.md

## Goal
Build the new terminal creation modal matching the design, and implement terminal lifecycle (create, resize, reconnect, delete).

## Steps
1. New terminal modal component (matching new_terminal_modal.html)
2. Form fields: name, starting directory, shell, color tag
3. Permission toggles: auto-record, read-only sharing, persist after disconnect
4. Live preview box
5. Create terminal flow: modal → API → PTY spawn → auto-attach WebSocket
6. Delete terminal flow: close tab → API → kill PTY
7. Resize handling: xterm-addon-fit resize → WebSocket resize message → PTY resize
8. Session recovery: reconnect to running PTY after browser refresh/disconnect

## Acceptance Criteria
- [ ] Modal matches new_terminal_modal.html design
- [ ] Creating terminal spawns PTY and opens in new tab
- [ ] Color tags work in sidebar/tabs
- [ ] Resize works on window resize
- [ ] Reconnect works after page refresh

## Status
TODO
