# Task: T04 — Terminal UI (Frontend Core)

## Parent Plan
001-termhub-mvp.md

## Goal
Build the main terminal interface matching the design mockup — 3-column layout with sidebar, xterm.js terminal, and info panel.

## Steps
1. App layout component (topbar, body with 3 columns, status bar)
2. Topbar (brand, tab bar with dots/labels/close, user avatar/role)
3. Sidebar (My Terminals list, Shared section, Actions, Host credentials)
4. Terminal pane with xterm.js (toolbar, terminal screen, status bar)
5. xterm.js WebSocket integration (connect, stdin/stdout, reconnect)
6. Tab system (Zustand store: open tabs, active tab, add/remove/switch)
7. Right panel (session info, processes, activity log, quick actions)
8. Status bar (connection state, PTY ID, duration timer, dimensions, CPU)
9. Terminal toolbar buttons (copy, record, maximize)
10. Responsive resize handling (xterm-addon-fit)

## Acceptance Criteria
- [ ] UI matches terminal_platform_ui_design.html mockup
- [ ] xterm.js renders terminal output correctly
- [ ] Keyboard input is sent to backend PTY
- [ ] Tabs work (create, switch, close)
- [ ] Sidebar shows owned terminals with status dots
- [ ] Right panel shows session info

## Status
TODO
