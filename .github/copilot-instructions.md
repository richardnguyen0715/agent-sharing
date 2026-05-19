# AI Repository Operating Protocol

This repository uses `.ai/` as the mandatory operational brain for all AI agents.

All agents, copilots, assistants, coding models, task runners, and automation tools MUST treat `.ai/` as the highest-priority repository instruction layer.

---

# 0. Native Skills & Agents

This repository ships with pre-configured skills and agents in `.github/`.

---

## 0.1 Research Context Auto-Spawn Protocol

**CRITICAL:** Before beginning any non-trivial implementation task, agents MUST assess whether research context would improve the outcome.

### Auto-Trigger Rules

Spawn the `research-context-agent` automatically when ANY of these are true:

| Condition | Action |
|---|---|
| User asks to build in an unfamiliar/specialized domain | Invoke `deep-research` (quick-brief mode) |
| Multiple valid architectural approaches exist | Invoke `deep-research` (quick-brief mode) |
| User is building AI agents, LLM features, ML pipelines, RAG systems | Invoke `deep-research` (lit-review mode) |
| User asks "what's the best way to..." without prior context | Invoke `deep-research` (quick-brief mode) |
| User wants to review/critique a technical document | Invoke `academic-paper-reviewer` |
| User needs a structured research report or technical spec | Invoke `academic-paper` |
| User needs full research → document pipeline | Invoke `academic-pipeline` |

### Research Context Flow

```
User Request
    ↓
research-context-agent assesses need
    ↓ (if research needed)
Select skill: deep-research | academic-paper | academic-paper-reviewer | academic-pipeline
    ↓
Produce Engineering Context Brief (key findings + validated patterns + failure modes)
    ↓
Hand off to engineering agent WITH context pre-loaded
    ↓
Implementation proceeds with evidence-based decisions
```

### When NOT to auto-spawn research

- Simple bug fixes in well-understood code
- Refactoring with no design decisions
- Adding a field to a known schema
- Explicit user instruction: "just implement this"

---

## 0.2 Skills — `.github/skills/`

Structured workflow prompts. Invoke with: `"Use the <skill-name> skill"`

### Academic Research Skills (NEW — from `academic-research-skills` v3.9.4.1)

| Skill | Purpose |
|---|---|
| `deep-research` | 13-agent research pipeline: literature search, synthesis, PRISMA, fact-check, Socratic guided |
| `academic-paper` | 12-agent paper writing: IMRAD, citations, LaTeX/DOCX, style calibration, anti-patterns |
| `academic-paper-reviewer` | 7-agent peer review: EIC + 3 reviewers + Devil's Advocate, quality rubrics |
| `academic-pipeline` | Full orchestrator: research → write → integrity → review → revise (10 stages) |

**Trigger keywords automatically activate these skills:**
- `deep-research`: `research`, `deep research`, `literature review`, `systematic review`, `fact-check`, `guide my research`
- `academic-paper`: `write paper`, `academic paper`, `paper outline`, `write abstract`, `guide my paper`
- `academic-paper-reviewer`: `review paper`, `peer review`, `manuscript review`, `critique paper`
- `academic-pipeline`: `academic pipeline`, `research to paper`, `full paper workflow`, `end-to-end paper`

### Engineering Workflow Skills

| Category | Skills |
|---|---|
| Engineering Workflow | `spec-driven-development`, `planning-and-task-breakdown`, `incremental-implementation`, `test-driven-development`, `debugging-and-error-recovery`, `code-review-and-quality`, `security-and-hardening`, `performance-optimization`, `git-workflow-and-versioning`, `documentation-and-adrs`, `api-and-interface-design`, `frontend-ui-engineering`, `ci-cd-and-automation`, `shipping-and-launch`, `deprecation-and-migration`, `doubt-driven-development`, `source-driven-development` |
| Reference & Knowledge | `repo-analysis`, `context-extraction`, `reference-lookup`, `pattern-adoption` |
| Product & PM | `deliver-prd`, `deliver-user-stories`, `deliver-acceptance-criteria`, `define-problem-statement`, `define-hypothesis`, `foundation-lean-canvas`, `measure-experiment-design`, `iterate-retrospective` |
| Architecture | `develop-adr` |

See full index: [`.github/skills/README.md`](.github/skills/README.md)

---

## 0.3 Agents — `.github/agents/`

Specialized domain experts. Invoke with `/agent <name>` or via Copilot CLI.

| Category | Agents |
|---|---|
| **Knowledge (Research)** | `research-context-agent` ⚡ (auto-spawned), `academic-research-orchestrator`, `repo-analyst`, `context-curator`, `reference-navigator` |
| Engineering | `backend-developer`, `frontend-developer`, `fullstack-developer`, `python-pro`, `typescript-pro` |
| Architecture | `api-designer`, `api-documenter`, `database-administrator` |
| Quality | `code-reviewer`, `test-engineer`, `test-automator`, `security-auditor`, `performance-engineer` |
| DevOps | `devops-engineer`, `docker-expert`, `terraform-engineer` |
| AI & ML | `ml-engineer` |
| Debug | `debugger` |

⚡ = auto-spawned based on request context

See full index: [`.github/agents/README.md`](.github/agents/README.md)

---

# 1. Mandatory Session Bootstrap

At the beginning of EVERY new session, task, implementation, analysis, refactor, debugging flow, or tool execution:

MUST READ in order:

1. `.ai/context/policies-you-must-follow-first.md`
2. `.ai/rules/coding-policies.md`
3. `.ai/rules/how-to-structure.md`
4. `.ai/rules/testing-debugging-frameworks.md`

Agents MUST load these files into active reasoning context BEFORE making changes.

No implementation is allowed before completing this bootstrap phase.

---

# 2. `.ai/` Is The Source Of Truth

The `.ai/` directory is the operational source of truth for:

- architecture
- standards
- documentation
- workflows
- implementation phases
- checkpoints
- ADRs
- planning
- task decomposition
- repository analysis

If conflict exists between implementation and `.ai/` documentation:

- STOP
- analyze discrepancy
- create/update ADR
- document reasoning
- never silently diverge

---

# 3. Mandatory Operational Workflow

For EVERY non-trivial implementation:

## Step 1 — Analyze

Read:

- current architecture
- related ADRs
- active phases
- checkpoints
- related tasks

Document findings if missing.

---

## Step 2 — Plan

Create/update:

- planning tasks
- implementation phases
- dependency impact
- architectural impact

Required locations:

- `.ai/planning-tasks/`
- `.ai/runtime/phases-docs/`

---

## Step 3 — Implement

Implementation MUST follow:

- coding policies
- structure rules
- testing/debugging frameworks
- repository architecture

No shortcuts.

---

## Step 4 — Validate

Before completion:

- verify architecture consistency
- verify tests
- verify documentation consistency
- verify no hidden side-effects

---

## Step 5 — Document

After major implementation:

MUST update:

### Checkpoints
- `.ai/runtime/checkpoints/`
- template: `.ai/templates/checkpoint-template.md`

### Phase Docs
- `.ai/runtime/phases-docs/`
- template: `.ai/templates/phase-template.md`

### ADRs
- `.ai/runtime/adr/`
- template: `.ai/templates/adr-template.md`

Documentation updates are NOT optional.

---

# 4. Forbidden Actions

Agents MUST NEVER:

- modify architecture silently
- add dependencies without explicit reasoning
- bypass repository standards
- ignore ADRs
- skip documentation updates
- leave undocumented blockers
- create hidden technical debt knowingly
- introduce inconsistent patterns
- overwrite existing operational context without analysis

---

# 5. Architecture Governance

For:

- architectural changes
- dependency changes
- infra changes
- framework changes
- state management changes
- data model changes

Agents MUST:

1. analyze impact
2. document rationale
3. create/update ADR
4. explain tradeoffs
5. identify migration implications

---

# 6. Documentation Contract

All documentation MUST be:

- implementation-aware
- synchronized with current architecture
- actionable
- explicit
- operationally useful

Avoid placeholder documentation.

---

# 7. Tooling Behavior

These rules apply across:

- Copilot
- Cursor
- Cline
- Roo
- Continue
- Claude Code
- OpenAI agents
- custom agents
- automation pipelines

The agent MUST continuously reference `.ai/` during execution, not only at session start.

---

# 8. Repository Philosophy

This repository follows:

- documentation-driven engineering
- architecture-governed implementation
- ADR-first decision tracking
- phase-oriented delivery
- checkpoint-based progress tracking

Agents are expected to operate as disciplined engineering collaborators, not autocomplete tools.