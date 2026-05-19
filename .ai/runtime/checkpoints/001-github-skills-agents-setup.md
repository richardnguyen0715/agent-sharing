# Checkpoint 001 — .github Skills & Agents Setup

```yaml
agent: copilot
timestamp: 2026-05-19
related_phase: initial-setup
related_task: github-native-skills-agents
```

## Completed

- Populated `.github/skills/` with **30 native skills** (flat structure for skill discovery)
- Populated `.github/agents/` with **21 native agents** (categorized into 7 subdirectories)
- Created `README.md` index files for both skills and agents directories
- Created `.github/README.md` — top-level `.github/` quick reference
- Updated `.github/copilot-instructions.md` — added Section 0 documenting all skills/agents + fixed bootstrap file paths

## Structure

```
.github/
├── copilot-instructions.md   ← updated with skill/agent registry
├── README.md                 ← new: structure overview + quick reference
├── agents/                   ← 21 agents in 7 categories
│   ├── engineering/          (5: backend, frontend, fullstack, python-pro, typescript-pro)
│   ├── architecture/         (3: api-designer, api-documenter, database-administrator)
│   ├── quality/              (5: code-reviewer, test-engineer, test-automator, security-auditor, performance-engineer)
│   ├── devops/               (3: devops-engineer, docker-expert, terraform-engineer)
│   ├── ai-ml/                (1: ml-engineer)
│   ├── knowledge/            (3: repo-analyst, context-curator, reference-navigator)
│   └── debug/                (1: debugger)
└── skills/                   ← 30 skills (flat, each with SKILL.md)
    ├── [17 engineering workflow skills]
    ├── [4 reference/knowledge skills]
    ├── [8 product/PM skills]
    └── [1 architecture skill]
```

## Skills Included

Engineering: spec-driven-development, planning-and-task-breakdown, incremental-implementation,
test-driven-development, debugging-and-error-recovery, code-review-and-quality,
security-and-hardening, performance-optimization, git-workflow-and-versioning,
documentation-and-adrs, api-and-interface-design, frontend-ui-engineering,
ci-cd-and-automation, shipping-and-launch, deprecation-and-migration,
doubt-driven-development, source-driven-development

Reference: repo-analysis, context-extraction, reference-lookup, pattern-adoption

Product: deliver-prd, deliver-user-stories, deliver-acceptance-criteria,
define-problem-statement, define-hypothesis, foundation-lean-canvas,
measure-experiment-design, iterate-retrospective, develop-adr

## Next Steps

- Add domain-specific agents as the app template matures (e.g., AI/agent-specific skills)
- Add `copilot-setup-steps.yml` to `.github/agent/` for cloud agent preinstall
- Consider adding more PM skills (deliver-edge-cases, deliver-release-notes) as needed

---

## Update — Academic Research Skills Added (2026-05-19)

**Source:** https://github.com/Imbad0202/academic-research-skills v3.9.4.1

### Skills added to `.github/skills/`

| Skill | Files | Sub-agents |
|---|---|---|
| `deep-research` | 50 | 13 agents (research_architect, synthesis, bibliography, meta_analysis, ethics_review, etc.) |
| `academic-paper` | 55 | 12 agents (intake, draft_writer, citation_compliance, formatter, peer_reviewer, etc.) |
| `academic-paper-reviewer` | 25 | 7 agents (eic, domain_reviewer, devils_advocate, editorial_synthesizer, etc.) |
| `academic-pipeline` | 30 | 5 agents (pipeline_orchestrator, integrity_verification, state_tracker, etc.) |
| `shared/` | 43 | Shared protocols, schemas, templates |

### Agents added to `.github/agents/knowledge/`

- `research-context-agent` — **auto-spawned** before complex tasks; injects evidence brief before implementation
- `academic-research-orchestrator` — manual; routes to the 4 academic skills for full research → document workflows

### Auto-spawn protocol

`copilot-instructions.md` §0.1 defines trigger rules:
- Unfamiliar domain → `deep-research` quick-brief
- AI/ML feature building → `deep-research` lit-review
- Document to review → `academic-paper-reviewer`
- Full research → doc → `academic-pipeline`

### Architecture decision

Full directory structure preserved (not just SKILL.md) because:
- Sub-agents, references, templates, and examples are required for proper skill execution
- `shared/` contains cross-skill protocols (style_calibration, compliance_checkpoint, etc.)
