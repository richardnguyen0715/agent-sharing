# AI Skills

Workflow skills native to this repository. Skills are structured prompts that guide Copilot through proven engineering workflows.

Invoke with `"Use the <skill-name> skill"` in chat.

---

## 🔬 Academic Research Skills

> **Source:** [academic-research-skills v3.9.4.1](https://github.com/Imbad0202/academic-research-skills)
> Multi-agent pipelines for rigorous research, paper writing, peer review, and full publication workflows.

| Skill | Agents | When to use |
|---|---|---|
| `deep-research` | 13 agents | Research any topic — literature search, synthesis, PRISMA, systematic review, fact-check, Socratic guided research |
| `academic-paper` | 12 agents | Write academic papers — IMRAD, LaTeX/DOCX, citations, bilingual abstract, peer review, revision coaching |
| `academic-paper-reviewer` | 7 agents | Multi-perspective peer review — EIC + 3 domain reviewers + Devil's Advocate, quality rubrics |
| `academic-pipeline` | Full orchestrator | End-to-end: research → write → integrity check → review → revise → publish (10 stages) |

**Auto-trigger keywords:**
- `deep-research`: research, literature review, systematic review, fact-check, guide my research, 研究, 文獻回顧
- `academic-paper`: write paper, academic paper, paper outline, guide my paper, 寫論文, 學術論文
- `academic-paper-reviewer`: review paper, peer review, critique paper, manuscript review
- `academic-pipeline`: academic pipeline, research to paper, full paper workflow, end-to-end paper

**Shared resources:** `.github/skills/shared/` — common protocols, schemas, templates used across all 4 skills.

---

## 🔧 Engineering Workflow

| Skill | When to use |
|---|---|
| `spec-driven-development` | Starting any new feature — write spec before code |
| `planning-and-task-breakdown` | Decompose work into ordered, verifiable tasks |
| `incremental-implementation` | Build slice-by-slice with validation at each step |
| `test-driven-development` | Implementing logic or fixing bugs — tests first |
| `debugging-and-error-recovery` | Something is broken — root cause workflow |
| `code-review-and-quality` | Before merging — quality and correctness check |
| `security-and-hardening` | Auth, input validation, secrets, threat modeling |
| `performance-optimization` | Performance concerns — measure then optimize |
| `git-workflow-and-versioning` | Committing, branching, PR workflow |
| `documentation-and-adrs` | Writing docs, ADRs, runbooks |
| `api-and-interface-design` | Designing APIs or module interfaces |
| `frontend-ui-engineering` | Building UI components and views |
| `ci-cd-and-automation` | Pipeline setup and automation |
| `shipping-and-launch` | Deploying to production safely |
| `deprecation-and-migration` | Removing old systems, safe migrations |
| `doubt-driven-development` | High-stakes or unfamiliar code — verify before acting |
| `source-driven-development` | Doc-verified framework code — no guessing |

---

## 📚 Reference & Knowledge

| Skill | When to use |
|---|---|
| `repo-analysis` | Analyze a cloned GitHub repo — extract patterns |
| `context-extraction` | Distill repo analysis into reusable KB entries |
| `reference-lookup` | Search KB for patterns before implementing |
| `pattern-adoption` | Integrate a pattern from references into current codebase |

---

## 🎯 Product & PM

| Skill | When to use |
|---|---|
| `deliver-prd` | Write a Product Requirements Document |
| `deliver-user-stories` | Generate well-formed user stories |
| `deliver-acceptance-criteria` | Define acceptance criteria for stories |
| `define-problem-statement` | Frame the problem before solutioning |
| `define-hypothesis` | Write a testable hypothesis |
| `foundation-lean-canvas` | Build a lean canvas for a product idea |
| `measure-experiment-design` | Design a measurable experiment |
| `iterate-retrospective` | Run a structured retrospective |

---

## 🏛️ Architecture

| Skill | When to use |
|---|---|
| `develop-adr` | Write an Architecture Decision Record |

---

## Adding New Skills

1. Create a new directory under `.github/skills/<skill-name>/`
2. Add a `SKILL.md` file with frontmatter:
   ```md
   ---
   name: skill-name
   description: When to use this skill (shown in skill picker)
   ---
   ```
3. Update this README

