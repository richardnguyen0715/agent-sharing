The design is in @.ai/context/spec/design/terminal_platform_ui_design.html

Now here's the full breakdown of every UX decision in this design and why it was made:

---

## Layout structure

The app uses a 3-column shell: sidebar + terminal canvas + info panel, with a persistent top bar. This mirrors established patterns from VS Code, JetBrains, and web terminal tools like Wetty — users familiar with any of these will orient immediately.

The fixed chrome (topbar + statusbar) keeps navigation controls out of the "content zone" so the terminal feels full-bleed. Every pixel in the center pane belongs to the terminal output.

---

## Tab system

Tabs sit in the topbar and mirror browser tab behavior, which users already know. Each tab shows:

- a colored dot (green = idle, amber = running process) — tells the user what's happening without opening the tab
- a short label derived from the CWD or script name
- a close button that only appears on hover

This solves the "multiple terminals" use case with zero learning curve.

---

## Sidebar: owned terminals + shared context

The left sidebar separates three things users actually need fast access to:

1. their own terminals (with live status dots)
2. shared files/environment they have read access to
3. host credentials (GitHub, AWS, VPN) — always visible so users know the context they're running in

This "credentials always visible" pattern is critical for your use case. Users need to trust that their commands are running with the right authentication context. Making it ambient removes the need to ask.

---

## Terminal rendering

The dark terminal theme (`#0D1117`, the GitHub dark palette) is the single most consistent expectation users bring to any terminal. Don't fight it. Semantic coloring follows standard conventions — green for prompts, blue for paths, gray for output, red for errors. Users parse these subconsciously after years of shell use.

The blinking cursor is essential — it signals "ready for input" with no text needed.

---

## Status bar

The bottom status bar borrows from VS Code's design: small, always present, shows the most operationally relevant info: connection state, PTY ID, session duration, terminal dimensions, and CPU. These are exactly the things a user wants to glance at when something is slow or broken.

---

## Right panel: session info + activity log

This is the non-terminal interface to the terminal session. It exposes:

- the process tree (so users know what's running without typing `ps`)
- a command history/activity log with timestamps — feeds into audit logging
- quick actions like recording and read-only sharing

For an app where accountability and auditability are design requirements (per your spec), this panel makes those features visible and accessible rather than buried in admin settings.