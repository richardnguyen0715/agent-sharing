<div align="center">

# 🧠 AI-Native Repository Template

**The last template you'll ever clone.**
A production-ready, AI-agent-first repository scaffold — built for developers who ship with Copilot, Claude, Cursor, and every AI coding tool that comes next.

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](#)
[![AI Ready](https://img.shields.io/badge/AI-native-8A2BE2)](#)
[![Skills](https://img.shields.io/badge/skills-35-green)](#-skills--35)
[![Agents](https://img.shields.io/badge/agents-23-orange)](#-agents--23)
[![Stack](https://img.shields.io/badge/stack-gold--standard-gold)](#-architecture)

</div>

---

## The problem with every other template

Most repo templates give you folder structure and a README. Then you're on your own.

You still have to:
- Explain your codebase to every AI tool from scratch
- Prompt-engineer the same workflows every session
- Rediscover decisions made 3 months ago
- Hope your AI assistant doesn't invent architecture

**This template is different.** It's built around a single idea:

> _The repository itself is the AI's operating system._

---

## What's inside

### 🤖 35 Native Skills

Structured, prompt-based workflows your AI runs — not guesses, but **proven engineering playbooks** baked directly into `.github/skills/`.

**Engineering**
| | |
|---|---|
| `spec-driven-development` | Write spec before writing a single line of code |
| `planning-and-task-breakdown` | Break any feature into atomic, verifiable tasks |
| `test-driven-development` | Red → Green → Refactor, every time |
| `debugging-and-error-recovery` | Systematic root-cause workflow, not vibes |
| `security-and-hardening` | Threat model before you ship |
| `doubt-driven-development` | Verify before touching high-stakes code |
| + 11 more | `ci-cd-and-automation`, `shipping-and-launch`, `api-and-interface-design`... |

**Academic Research** _(powered by [academic-research-skills v3.9.4.1](https://github.com/Imbad0202/academic-research-skills))_
| | |
|---|---|
| `deep-research` | 13-agent pipeline: literature search, PRISMA, synthesis, Socratic guided research |
| `academic-paper` | 12-agent pipeline: full paper writing with citations, LaTeX/DOCX, peer review |
| `academic-paper-reviewer` | 7-agent simulation: EIC + 3 domain reviewers + Devil's Advocate |
| `academic-pipeline` | Full orchestrator: research → write → integrity → review → revise |

**Product & PM**
`deliver-prd` · `deliver-user-stories` · `define-problem-statement` · `foundation-lean-canvas` · `measure-experiment-design` · and more

---

### 🦾 23 Native Agents

Specialized AI personas pre-loaded into `.github/agents/` — no setup, no prompting, just `/agent <name>`.

| Category | Agents |
|---|---|
| **Knowledge** ⚡ | `research-context-agent` _(auto-spawns)_, `academic-research-orchestrator` |
| **Engineering** | `backend-developer`, `frontend-developer`, `fullstack-developer`, `python-pro`, `typescript-pro` |
| **Architecture** | `api-designer`, `api-documenter`, `database-administrator` |
| **Quality** | `code-reviewer`, `test-engineer`, `test-automator`, `security-auditor`, `performance-engineer` |
| **DevOps** | `devops-engineer`, `docker-expert`, `terraform-engineer` |
| **AI & ML** | `ml-engineer` |
| **Debug** | `debugger` |

> ⚡ The `research-context-agent` **auto-spawns** before complex tasks — it detects when domain research would improve your implementation, gathers evidence, and injects a brief before your engineering agent starts. Your AI builds on proof, not assumptions.

---

### 🧬 `.ai/` — The Operational Brain

The `.ai/` directory is the persistent memory system that makes your AI a disciplined engineering collaborator across every session, every context reset, every new model.

```
.ai/
├── context/
│   ├── policies-you-must-follow-first.md   ← Session bootstrap. Read before anything.
│   ├── spec/                               ← Product requirements & system design
│   ├── instructions/                       ← Historical decisions & directives
│   └── analysis-results/                  ← Research & competitive analysis
├── rules/
│   ├── coding-policies.md                 ← Engineering laws. No exceptions.
│   ├── how-to-structure.md               ← Gold-standard repo architecture
│   └── testing-debugging-frameworks.md   ← Test strategy & debug procedures
├── runtime/
│   ├── checkpoints/                       ← Where we are. Resume in < 10 min.
│   ├── phases-docs/                       ← Why decisions were made
│   ├── adr/                               ← Architecture Decision Records
│   └── planing-tasks/plan/ + tasks/       ← Approved plans & atomic task files
└── templates/
    ├── checkpoint-template.md
    ├── adr-template.md
    └── phase-template.md
```

**The protocol guarantees:**
- ✅ Resume any project in <10 minutes, no context needed
- ✅ No silent architecture changes — every decision is traceable
- ✅ Any AI tool reads `.ai/` and knows exactly where things stand
- ✅ Planning gate: no implementation without an approved plan

---

### 🏗️ Gold-Standard Repository Architecture

This template follows the 2026 gold-standard monorepo structure — ready to scale to web, mobile, desktop, AI systems, and multi-agent workflows:

```
repo/
├── apps/           ← Next.js · Expo · Tauri · Astro
├── packages/       ← ui · auth · database · ai · agents · prompts
├── services/       ← API · workers · jobs
├── infrastructure/ ← Terraform · Kubernetes · Helm
├── platform/       ← CI/CD · Docker · secrets · environments
├── docs/           ← architecture · ADRs · runbooks · onboarding
├── .ai/            ← AI operational brain (see above)
├── .github/        ← Skills · Agents · Copilot instructions
└── .devcontainer/  ← Reproducible dev environment
```

---

## How it works

**Session 1 — New project**
```
You: "I want to build a RAG-based customer support bot"

Copilot reads .ai/context/policies-you-must-follow-first.md
↓
research-context-agent auto-spawns → deep-research (lit-review mode)
↓
Evidence brief: RAG patterns, known failure modes, validated architectures
↓
spec-driven-development skill → structured spec
↓
planning-and-task-breakdown skill → approved task list
↓
implementation begins with architecture that won't embarrass you in 6 months
```

**Session 47 — New context window, same project**
```
New agent reads .ai/runtime/checkpoints/ → knows current state in 30 seconds
Reads .ai/runtime/adr/ → knows why decisions were made
Reads .ai/runtime/planing-tasks/ → knows what's left
Continues exactly where things left off
```

---

## Who this is for

- **AI-first engineers** who use Copilot, Claude, Cursor, or Cline daily
- **Solo builders** who want their AI to remember everything between sessions
- **Teams** who are tired of AI hallucinating architecture that doesn't match the codebase
- **Researchers & technical writers** who need rigorous, citation-backed documentation
- **Anyone who's ever said** _"I spent 20 minutes re-explaining my project to my AI"_

---

## Quick start

```bash
# Clone and rename
git clone https://github.com/your-username/repo-template.git my-new-project
cd my-new-project

# Tell your AI what project this is
# Edit .ai/context/spec/ with your requirements
# Edit .ai/context/policies-you-must-follow-first.md with your stack

# Start a session — your AI bootstraps from .ai/ automatically
```

The first thing your AI will do in any session: read `.ai/context/policies-you-must-follow-first.md`. That file is your project's ground truth.

---

## Compatibility

Works with any AI coding tool that reads instruction files:

| Tool | Config file used |
|---|---|
| GitHub Copilot | `.github/copilot-instructions.md` + `.github/skills/` + `.github/agents/` |
| Claude Code | `.claude/` (add your CLAUDE.md pointing to `.ai/`) |
| Cursor | `.cursor/rules/` (symlink to `.ai/rules/`) |
| Windsurf | `.windsurf/` |
| Cline / Roo / Continue | Point to `.ai/context/policies-you-must-follow-first.md` |

---

## Philosophy

**Documentation is not overhead. It's the operating system.**

Every AI tool you'll use in the next 5 years will be stateless. Context windows will reset. Models will change. New tools will arrive.

The only thing that persists is the repository.

This template is built on that premise: the repo carries the memory, the reasoning, and the rules. The AI is the execution engine. You are the architect.

---

<div align="center">

If this saves you even one "wait, why did we build it this way?" — give it a ⭐

**[Use this template →](https://github.com/your-username/repo-template/generate)**

</div>
