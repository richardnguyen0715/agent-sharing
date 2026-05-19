# AI Agents

Specialized Copilot agents native to this repository. Each agent is a domain expert with focused instructions and tools.

Invoke with `/agent <name>` or via the GitHub Copilot CLI.

---

## 📚 Knowledge & Research

> These agents enrich all other agents with evidence-based context.

| Agent | Auto-Spawn | When to use |
|---|---|---|
| `research-context-agent` | ⚡ **Yes** | Auto-spawned before complex tasks — gathers deep research context, injects evidence brief before implementation |
| `academic-research-orchestrator` | Manual | Full research → document pipeline; routes to `deep-research`, `academic-paper`, `academic-paper-reviewer`, `academic-pipeline` |
| `repo-analyst` | Manual | Deeply analyze a cloned GitHub repository |
| `context-curator` | Manual | Distill repo analyses into reusable knowledge base |
| `reference-navigator` | Manual | Find patterns from references for current task |

Location: [`knowledge/`](./knowledge/)

**⚡ Auto-spawn protocol:** `research-context-agent` activates automatically when the agent detects that implementation would benefit from domain research. See `.github/copilot-instructions.md` §0.1 for full trigger rules.

---

## 🏗️ Engineering

| Agent | When to use |
|---|---|
| `backend-developer` | Server-side APIs, microservices, backend systems |
| `frontend-developer` | React/Vue/Angular, UI components, client-side |
| `fullstack-developer` | Features spanning database, API, and frontend together |
| `python-pro` | Idiomatic Python 3.11+, async patterns, type safety |
| `typescript-pro` | Advanced TypeScript generics, type-level programming |

Location: [`engineering/`](./engineering/)

---

## 🎨 Architecture

| Agent | When to use |
|---|---|
| `api-designer` | Designing REST/GraphQL APIs, OpenAPI specs, versioning |
| `api-documenter` | API documentation, OpenAPI, interactive portals |
| `database-administrator` | DB performance, HA, disaster recovery, schema design |

Location: [`architecture/`](./architecture/)

---

## ✅ Quality

| Agent | When to use |
|---|---|
| `code-reviewer` | Comprehensive code review — bugs, security, logic |
| `test-engineer` | Test strategy, test writing, coverage analysis |
| `test-automator` | Automated test frameworks, CI test integration |
| `security-auditor` | Security audits, vulnerability assessment, compliance |
| `performance-engineer` | Bottleneck identification, profiling, optimization |

Location: [`quality/`](./quality/)

---

## ⚙️ DevOps & Infrastructure

| Agent | When to use |
|---|---|
| `devops-engineer` | CI/CD pipelines, infrastructure automation, deployment |
| `docker-expert` | Container images, Docker optimization, orchestration |
| `terraform-engineer` | IaC, multi-cloud Terraform, module architecture |

Location: [`devops/`](./devops/)

---

## 🤖 AI & ML

| Agent | When to use |
|---|---|
| `ml-engineer` | Production ML pipelines, model serving, retraining |

Location: [`ai-ml/`](./ai-ml/)

---

## 🐛 Debug

| Agent | When to use |
|---|---|
| `debugger` | Root cause analysis, error tracing, bug diagnosis |

Location: [`debug/`](./debug/)

---

## Adding New Agents

1. Create a new `.md` file in the appropriate category subdirectory
2. Follow the frontmatter format:
   ```md
   ---
   name: agent-name
   description: When to use this agent (shown in agent picker)
   tools: ["read", "write", "edit", "bash", "glob", "grep"]
   ---
   ```
3. Update this README

