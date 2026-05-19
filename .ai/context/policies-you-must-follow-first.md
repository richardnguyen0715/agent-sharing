# AI Agent Operational Execution Standard

==================================================
SYSTEM PURPOSE
==================================================

This repository structure exists to transform AI agents from
stateless code generators into persistent engineering operators.

The system guarantees:

- operational continuity
- architectural consistency
- resumable execution
- traceable reasoning
- multi-session stability
- low context loss

All AI agents working inside this repository MUST follow this standard.

Failure to follow these rules is considered invalid execution.

---

# 1. GLOBAL EXECUTION PRINCIPLES

## 1.1 Documentation Is Operational Memory

Documentation is NOT optional.

Documentation represents:

- persistent memory
- execution state
- architectural reasoning
- implementation continuity
- handoff context
- future session recovery

Agents MUST treat documentation as part of the system runtime.

---

## 1.2 Assume Complete Context Loss

Agents MUST assume:

- conversation memory will disappear
- context windows will reset
- future sessions may use different agents
- future agents may know nothing about prior work

Therefore:

- all critical reasoning must be externalized
- assumptions must be documented
- architectural intent must persist outside chat history

---

## 1.3 No Hidden Actions

Agents MUST NOT:

- modify architecture silently
- add/remove dependencies without justification
- introduce workarounds without debt tracking
- create/delete important files without logging
- change system behavior without updating checkpoints

All meaningful changes must remain traceable.

---

## 1.4 Resumability First

The repository must always remain resumable.

A new agent should be able to:

1. read the operational context
2. understand the current architecture
3. identify unfinished work
4. continue implementation safely

Target:

```text
Resume project understanding within <10 minutes
without reading the entire repository.
```

---

# 2. DIRECTORY STRUCTURE

```text
.ai/
├──context/
│   ├── spec/
│   ├── instructions/
│   ├── analysis-results/
│   └── references-repo/
│   └── policies-you-must-follow-first.md
├── rules/
├── runtime/
│   ├── checkpoints/
│   ├── phases-docs/
│   ├── adr/
│   ├── dept/
│   └── planing-tasks/
│       ├── plan/       ← full implementation plans
│       └── tasks/      ← individual task files
└── templates/
```

---

# 3. DIRECTORY RESPONSIBILITIES

# 3.1 @.ai/context/spec/

Priority: HIGH

Purpose:

```text
Source of truth for system requirements.
```

Contains:

- product requirements
- feature specifications
- system design
- architecture descriptions
- workflows
- API contracts
- database design
- infrastructure requirements
- performance expectations
- scalability requirements

Agents MUST:

- read relevant specs before implementation
- align implementation with specifications
- avoid speculative architecture changes

Agents MUST NOT:

- override specs without explicit justification
- invent undocumented behavior silently

---

# 3.2 @.ai/rules/

Priority: HIGH

Purpose:

```text
Repository operational laws and engineering constraints.
```

Contains:

- coding standards
- architecture constraints
- testing policies
- debugging procedures
- security requirements
- performance rules
- naming conventions
- review standards
- documentation rules
- dependency policies

Agents MUST:

- obey policies before implementation
- validate generated code against policies
- prioritize policies over convenience

If implementation conflicts with policy:

- the conflict MUST be documented
- the justification MUST be explicit

---

# 3.3 @.ai/context/instructions/

Priority: MEDIUM

Purpose:

```text
Historical instruction memory.
```

Contains:

- prior prompts
- user directives
- historical implementation requests
- workflow instructions
- strategic engineering decisions

Purpose of this directory:

- preserve long-term execution context
- avoid repeated user explanations
- maintain behavioral continuity across sessions

Agents SHOULD:

- consult this directory before major implementation work
- avoid contradicting historical instructions

Agents MUST NOT:

- treat historical instructions as higher priority than current specs or policies

Priority order:

```text
Current Request
> Rules
> Specs
> Instructions
```

---

# 3.4 @.ai/runtime/checkpoints/

Priority: HIGH

Purpose:

```text
Persistent operational checkpoint state.
```

This directory acts as the repository memory layer.

Contains:

- current implementation state
- active progress tracking
- unfinished work
- operational summaries
- active blockers
- resumable execution context

Agents MUST update checkpoints after:

- major installations
- architecture changes
- large feature completion
- migration work
- dependency changes
- session-ending implementation work

Checkpoint files should include:

- completed work
- in-progress work
- planned work
- current phase
- blockers
- known risks
- recommended next actions

Goal:

```text
Future agents must understand current repository state immediately.
```

---

# 3.5 @.ai/runtime/phases-docs/

Priority: MEDIUM

Purpose:

```text
Architectural reasoning and implementation traceability based on implementation phases and @.ai/context/spec/.
```

Contains:

- implementation phases
- architectural decisions
- technical tradeoffs
- scalability considerations
- migration reasoning
- extensibility implications
- implementation plans
- execution strategies

Agents SHOULD document:

- WHY a decision was made
- alternative approaches considered
- tradeoffs introduced
- future risks
- scalability impact

Agents MUST avoid:

- verbose implementation diaries
- low-value action logging
- redundant explanations already visible in code

Focus on:

```text
Decision reasoning, not code narration.
```

---

# 3.6 @.ai/context/analysis-results/

Priority: MEDIUM

Purpose:

```text
Research and analytical knowledge base.
```

Contains:

- competitor system analysis
- reference architecture analysis
- benchmarking reports
- security analysis
- performance analysis
- reverse engineering notes
- design observations
- implementation breakdowns of external systems

Example:

```text
When building system A,
analysis of reference system B belongs here.
```

This directory is NOT implementation memory.

It is a research knowledge repository.

Agents SHOULD:

- reuse insights from prior analyses
- avoid repeating the same investigations
- reference findings when making architectural decisions

---

# 3.7 @.ai/runtime/planing-tasks/

Priority: **CRITICAL**

Purpose:

```text
Planning memory — every significant task must be planned, documented,
and approved before implementation begins.
```

Sub-directories:

| Path | Contents |
|---|---|
| `plan/` | Full implementation plan document per request/feature |
| `tasks/` | Individual task files, one per discrete unit of work |

### plan/ file format

Each plan file maps to one user request or feature.

Required structure:

```md
# Plan: <title>

## Summary
<1–3 sentence goal>

## Context
<relevant specs, checkpoints, constraints>

## Tech Stack / Approach
<chosen stack, patterns, tools>

## Phases
<ordered list of phases>

## Open Questions
<unresolved decisions>

## Approval Status
PENDING | APPROVED | CHANGE_REQUESTED
```

### tasks/ file format

Each task file represents one discrete unit of work derived from the plan.

Required structure:

```md
# Task: <id> — <title>

## Parent Plan
<plan filename>

## Goal
<what this task achieves>

## Steps
1. ...
2. ...

## Acceptance Criteria
- [ ] ...

## Status
TODO | IN_PROGRESS | DONE | BLOCKED
```

Agents MUST:

- create a plan file before any non-trivial implementation
- break the plan into task files in `tasks/`
- wait for user approval before executing any task
- update task status as work progresses

Agents MUST NOT:

- begin implementation before the plan is approved
- skip plan creation for "small" changes (if in doubt, plan it)

---

# 3.8 @.ai/context/references-repo/

Priority: LOW

Purpose:

```text
External reference codebases and implementation resources.
```

Contains:

- cloned repositories
- extracted code samples
- implementation references
- architectural examples
- reusable patterns

Purpose:

- implementation inspiration
- comparative analysis
- pattern extraction
- migration references

Agents MUST treat these repositories as:

```text
Reference material only.
```

Agents MUST NOT:

- blindly copy implementations
- assume compatibility with current architecture
- introduce external patterns without validation

---

# 4. IMPLEMENTATION EXECUTION RULES

# 4.1 Before Implementation

Agents MUST:

1. Read relevant specs
2. Read relevant policies
3. Read active checkpoints
4. Identify current implementation phase
5. Understand architectural constraints
6. **Create a full implementation plan** saved to `@.ai/runtime/planing-tasks/plan/`
7. **Break the plan into individual task files** saved to `@.ai/runtime/planing-tasks/tasks/`
8. **Present the complete plan to the user** — clearly showing:
   - proposed approach and tech stack
   - phases and task breakdown
   - key decisions and tradeoffs
   - any open questions
9. **Wait for explicit user approval** before writing any implementation code

User response options the agent MUST handle:

| User Response | Agent Action |
|---|---|
| ✅ Approved | Proceed with implementation as planned |
| ✏️ Change request | Update plan/tasks, re-present, wait for re-approval |
| 🔄 Rebuild request | Discard plan, re-analyze, produce new plan, wait for approval |

Agents MUST NOT begin implementation blindly.

Agents MUST NOT skip planning for "quick" or "small" tasks — if the user asks for any feature, system, or multi-step change, a plan is required.

The approval gate is **mandatory and non-negotiable**.

---

> **FIRST CATCH-UP PROTOCOL** — When an agent starts a new session and receives a user task:
>
> 1. Read this policies file
> 2. Read active checkpoints in `@.ai/runtime/checkpoints/`
> 3. Check existing plans in `@.ai/runtime/planing-tasks/plan/` and tasks in `@.ai/runtime/planing-tasks/tasks/`
> 4. If no approved plan exists → create plan → present to user → wait for approval
> 5. If an approved plan exists → confirm with user before continuing or replanning

---

# 4.2 During Implementation

Agents MUST document:

- important architectural decisions
- implementation direction changes
- blockers
- risky modifications
- temporary workarounds
- dependency changes

Agents SHOULD minimize unnecessary verbosity.

Rule:

```text
Document decisions and state,
not every keystroke.
```

---

# 4.3 After Major Work Completion

Agents MUST update:

```text
@.ai/runtime/checkpoints/
```

Including:

- completed tasks
- in-progress tasks
- planned tasks
- current repository state
- blockers
- risks
- next recommended actions

---

# 4.4 Architecture Decision Logging

Major decisions MUST be documented inside:

```text
@.ai/runtime/phases-docs/
```

Required for:

- architecture changes
- dependency changes
- runtime changes
- infrastructure changes
- major refactors
- migration strategies

Minimum required structure:

```md
Problem:
Decision:
Why:
Tradeoffs:
Future Impact:
```

---

# 5. TOKEN EFFICIENCY RULES

Documentation must remain lightweight.

Agents MUST optimize for:

- high signal
- low redundancy
- resumability
- operational clarity

Agents MUST avoid:

- repetitive summaries
- excessive prose
- duplicated information
- low-value implementation narration

Preferred:

```md
Redis session race condition blocks auth migration.
```

Avoid:

```md
While implementing the authentication system,
the agent discovered...
```

---

# 6. TRACEABILITY REQUIREMENTS

All meaningful changes should remain attributable.

Minimum metadata:

```yaml
agent:
timestamp:
related_phase:
related_task:
```

Required for:

- architectural changes
- migrations
- shared module modifications
- infrastructure updates
- dependency updates

---

# 7. TECHNICAL DEBT POLICY

Temporary solutions MUST be tracked.

Required structure:

```yaml
debt_id:
reason:
temporary_why:
proper_solution:
risk:
```

Agents MUST NOT normalize temporary fixes as permanent architecture.

---

# 8. PRIORITY ORDER

Agents MUST follow this priority order:

```text
1. Current User Request
2. Rules in @.ai/rules/
3. Specs in @.ai/context/spec/
4. Active Checkpoints in @.ai/runtime/checkpoints/
5. Historical Instructions in @.ai/context/instructions/
6. Reference Repositories in @.ai/context/references-repo/ and @.ai/context/analysis-results/
```

If conflicts occur:

- higher priority overrides lower priority
- conflicts MUST be documented

---

# 9. FINAL OPERATIONAL STANDARD

A valid AI-agent-compatible repository must guarantee:

```text
Operational Memory
+
Architectural Continuity
+
Explicit Reasoning
+
Traceable Decisions
+
Resumable Execution
+
Low Context Loss
+
Controlled Documentation Overhead
```

The protocol succeeds when:

- future agents avoid rediscovery
- architecture remains stable
- implementation continuity is preserved
- documentation remains concise
- repository state stays understandable
- token usage remains controlled