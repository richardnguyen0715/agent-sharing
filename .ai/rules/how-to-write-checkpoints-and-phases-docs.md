# Lean AI Agent Operational Memory Protocol

## Purpose

This protocol creates a lightweight operational memory layer for AI coding agents while minimizing token consumption.

Goals:

- preserve architectural continuity
- maintain execution traceability
- support multi-session resumability
- reduce redundant rediscovery
- minimize documentation overhead

Target:

```text
A new agent should resume safely in <10 minutes
without reading the entire repository.
```

---

# 1. CORE DESIGN PRINCIPLES

## 1.1 Documentation Must Be Minimal

Documentation exists to preserve:

- system state
- architectural intent
- execution continuity
- unresolved work

Documentation MUST NOT become:

- verbose reports
- duplicated explanations
- full development diaries

Rule:

```text
Only document information that future agents cannot reliably infer from code.
```

---

## 1.2 No Hidden Changes

Agents MUST NOT:

- modify architecture silently
- add dependencies without justification
- introduce workarounds without debt tracking
- end sessions without updating state

All meaningful changes must remain traceable.

---

## 1.3 Assume Context Loss

Agents must assume:

- context windows reset
- memory is non-persistent
- future sessions may use different agents

Therefore:

- reasoning must be explicit
- assumptions must be externalized
- important decisions must persist outside conversation memory

---

## 1.4 Resumability Over Completeness

The protocol optimizes for:

- fast continuation
- low ambiguity
- operational clarity

NOT exhaustive historical logging.

---

# 2. REPOSITORY STRUCTURE

```text
.ai/
├── checkpoint/
├── phases/
├── runtime/
└── rules/
```

Minimal structure only.

---

# 3. CHECKPOINT LAYER

Purpose:

```text
Current repository operational snapshot
```

Checkpoint files should remain compact and continuously updated.

---

# 3.1 REQUIRED FILES

## CURRENT_STATE.md

Primary repository snapshot.

Must contain only:

- active architecture
- major modules
- active priorities
- known risks
- current constraints
- next recommended actions

Avoid:

- implementation details
- historical logs
- repeated explanations

Recommended max size:

```text
300-500 lines
```

---

## STACK_RULES.md

Defines approved technologies and constraints.

Example:

```yaml
frontend:
  framework: nextjs
  styling: tailwind
  forbidden:
    - redux
    - css-modules

backend:
  framework: nestjs
  database: postgres
```

Purpose:

- prevent architectural drift
- reduce repeated decision-making

---

## SYSTEM_MAP.md

Defines:

- folder responsibilities
- ownership boundaries
- dependency rules
- unsafe coupling

Keep concise.

Only document high-level structure.

---

## ACTIVE_CONTEXT.md

Temporary working memory.

Contains:

- current focus
- active branch
- blockers
- unsafe areas
- current investigation targets

This file should be aggressively pruned.

Only active information belongs here.

---

## DECISIONS.md

Compact ADR log.

Only required for:

- architecture changes
- runtime changes
- dependency changes
- major design tradeoffs

Minimal format:

```md
## DECISION-012

Problem:
Decision:
Why:
Tradeoffs:
```

Do NOT create ADRs for trivial implementation details.

---

# 4. PHASE LAYER

Purpose:

```text
Execution intent + resumable implementation state
```

Phases represent meaningful initiatives.

Avoid creating phases for tiny changes.

---

# 4.1 PHASE STRUCTURE

```text
phases/
└── PHASE-001-auth-refactor/
    ├── overview.md
    ├── tasks.md
    ├── implementation-log.md
    └── pending.md
```

Minimize file count.

Avoid excessive fragmentation.

---

# 4.2 overview.md

Contains:

- business goal
- technical goal
- scope
- success criteria

Should remain short.

---

# 4.3 tasks.md

Single aggregated task file preferred.

Avoid one-file-per-task unless the phase is large.

Recommended structure:

```yaml
task_id:
status:
priority:
owner:
related_files:
```

And:

- objective
- implementation plan
- current progress
- unresolved issues

Avoid redundant prose.

---

# 4.4 implementation-log.md

Append-only.

Purpose:

```text
Capture non-obvious implementation reasoning.
```

Log ONLY:

- major findings
- architecture decisions
- unexpected issues
- risky modifications

Do NOT log every small edit.

Minimal format:

```md
[Timestamp]

Action:
Reason:
Impact:
Files:
```

---

# 4.5 pending.md

Most important resumability file.

Must contain:

- unfinished work
- safe continuation point
- risky areas
- recommended next steps

Optimize for fast handoff.

---

# 5. RUNTIME TRACE RULES

Only meaningful file changes require logging.

Required for:

- architecture files
- shared modules
- infra changes
- schema changes
- dependency changes

Not required for:

- formatting
- trivial renames
- local refactors
- generated files

Minimal format:

```yaml
file:
action:
reason:
impact:
```

---

# 6. CHANGE CLASSIFICATION

Every meaningful change must declare type:

| Type | Meaning |
|---|---|
| fix | bug fix |
| feature | new functionality |
| refactor | structural improvement |
| optimization | performance improvement |
| migration | architecture/runtime change |
| security | security-related |
| hotfix | urgent fix |

Avoid overclassification.

---

# 7. SESSION PROTOCOL

## 7.1 Session Start

Agents MUST read:

- CURRENT_STATE.md
- ACTIVE_CONTEXT.md
- STACK_RULES.md
- pending.md

Only these files are mandatory.

---

## 7.2 During Execution

Agents MUST update documentation only when:

- architecture changes
- assumptions change
- blockers appear
- implementation direction changes
- important discoveries occur

Avoid continuous micro-updates.

Batch updates whenever possible.

---

## 7.3 Before Session End

Mandatory updates:

1. ACTIVE_CONTEXT.md
2. pending.md
3. task status
4. unresolved issues
5. next recommended actions

CURRENT_STATE.md only needs updates if system state changed materially.

---

# 8. IMPLEMENTATION SAFETY RULES

Before risky implementation:

```md
Assumptions
Risks
Rollback Strategy
Affected Components
```

After implementation:

```md
Validation Results
Known Limitations
```

Only required for medium/high-risk changes.

---

# 9. TECHNICAL DEBT TRACKING

Required only for:

- temporary workarounds
- deferred architecture fixes
- known unstable implementations

Minimal format:

```yaml
debt_id:
reason:
proper_solution:
risk:
```

Do not track trivial debt.

---

# 10. AGENT ACCOUNTABILITY

Every meaningful change must remain attributable.

Minimal metadata:

```yaml
agent:
timestamp:
related_task:
```

Avoid excessive metadata verbosity.

Session IDs optional.

---

# 11. MACHINE-READABLE METADATA

Structured metadata should remain lightweight.

Recommended:

```yaml
id:
status:
updated_at:
related_files:
```

Only include fields actively used by automation.

Avoid metadata inflation.

---

# 12. TOKEN OPTIMIZATION RULES

## 12.1 Prefer State Over Narrative

Good:

```md
Auth migration blocked by Redis session invalidation race condition.
```

Bad:

```md
While implementing the authentication migration,
the agent discovered...
```

---

## 12.2 Avoid Duplicate Information

Never repeat the same information across:

- CURRENT_STATE
- ACTIVE_CONTEXT
- tasks
- pending

Each file should have a single responsibility.

---

## 12.3 Prefer Aggregation

Prefer:

```text
tasks.md
```

Over:

```text
TASK-001.md
TASK-002.md
TASK-003.md
```

Unless task complexity justifies separation.

---

## 12.4 Log Decisions, Not Actions

Good:

```md
Switched to centralized token refresh strategy.
```

Bad:

```md
Opened file A
Modified line B
Saved file C
```

---

## 12.5 Prune Aggressively

ACTIVE_CONTEXT.md and pending.md should continuously remove:

- completed investigations
- obsolete blockers
- resolved assumptions
- stale notes

---

# 13. ADVANCED OPTIONAL FEATURES

Only add these if repository scale justifies them.

## Optional

- architecture graph
- semantic change summaries
- vector memory index
- automated risk analysis
- autonomous continuation queue

Do NOT introduce these prematurely.

---

# 14. FINAL STANDARD

A valid AI coding environment should guarantee:

```text
Operational Memory
+
Architectural Continuity
+
Resumability
+
Explicit Reasoning
+
Low Documentation Overhead
+
Minimal Context Loss
```

The protocol succeeds when:

- future agents avoid rediscovery
- architectural drift decreases
- token usage remains controlled
- documentation stays maintainable
- repository state remains understandable