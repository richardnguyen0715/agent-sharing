# Gold-Standard Repository Architecture (2026)

For a modern platform that can build:

* Web apps
* Mobile apps
* Desktop apps
* APIs/services
* AI-powered systems
* Multi-agent workflows

…the current highest-standard approach is:

* **Monorepo**
* **Workspace-based**
* **Strong platform engineering**
* **Infrastructure-as-code**
* **AI-agent-aware**
* **Hermetic/devcontainer-first**
* **Observability-first**
* **Contract-driven**
* **Automation-heavy**

The strongest current ecosystem pattern is:

* `pnpm` workspaces
* `Turborepo` or `Nx`
* shared packages
* app/platform separation
* infrastructure versioned in repo
* CI/CD pipelines as code
* AI instruction files
* containerized development environments
* fully reproducible environments

This aligns with current monorepo and AI-agent repository standards. ([Turborepo][1])

---

# Recommended Stack

## Monorepo Engine

### Best default

* [Turborepo](https://turborepo.com?utm_source=chatgpt.com)
* `pnpm`

### Enterprise-heavy alternative

* [Nx](https://nx.dev?utm_source=chatgpt.com)

### Why

Turborepo currently provides the best balance of:

* speed
* caching
* DX
* CI optimization
* AI friendliness
* simplicity
* scaling

Nx is stronger for:

* enterprise governance
* generators
* graph tooling
* very large organizations

Turborepo structure conventions are now widely adopted. ([Turborepo][1])

---

# The Ideal Repository Structure

```txt
repo/
│
├── apps/
├── packages/
├── services/
├── infrastructure/
├── platform/
├── tooling/
├── scripts/
├── docs/
├── ai/
├── .github/
├── .devcontainer/
├── .vscode/
├── .husky/
├── tests/
├── examples/
├── tmp/
├── .env
├── .env.example
├── .gitignore
├── .editorconfig
├── .gitattributes
├── README.md
├── LICENSE
├── SECURITY.md
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
├── AGENTS.md
├── CLAUDE.md
├── GEMINI.md
├── turbo.json
├── pnpm-workspace.yaml
├── package.json
└── tsconfig.base.json
```

---

# Directory Standards

---

## `/apps`

User-facing applications.

```txt
apps/
├── web/
├── admin/
├── mobile/
├── desktop/
├── docs/
└── marketing/
```

### Recommended technologies

| App       | Stack               |
| --------- | ------------------- |
| web       | Next.js             |
| admin     | Next.js             |
| mobile    | React Native + Expo |
| desktop   | Tauri               |
| docs      | Nextra / Docusaurus |
| marketing | Astro               |

---

## `/packages`

Shared libraries.

```txt
packages/
├── ui/
├── design-system/
├── auth/
├── database/
├── api-client/
├── validation/
├── config/
├── logging/
├── observability/
├── ai/
├── agents/
├── prompts/
├── workflows/
├── types/
├── testing/
└── eslint-config/
```

This is the heart of scaling.

Everything reusable belongs here.

---

# AI-Native Repository Structure

This is now essential.

Current AI coding tools increasingly depend on repository instruction/configuration files. ([DeployHQ][2])

---

## `/ai`

```txt
ai/
├── prompts/
├── memory/
├── evals/
├── datasets/
├── tools/
├── skills/
├── agents/
├── policies/
├── workflows/
└── context/
```

---

# AI Agent Configuration Files

Modern repositories should support:

* Claude Code
* Cursor
* Copilot
* Gemini
* Codex
* Windsurf

Repository-level instruction files are becoming standard. ([DeployHQ][2])

---

## Root AI files

```txt
AGENTS.md
CLAUDE.md
GEMINI.md
.github/copilot-instructions.md
.cursor/rules/
.windsurf/
```

---

# Ideal `AGENTS.md`

```md
# AGENTS.md

## Architecture Rules

- Feature-first architecture
- Shared contracts only from packages/*
- No circular dependencies
- API schemas are source of truth

## Coding Rules

- Strict TypeScript
- No any
- Functional-first
- Pure business logic
- Zod validation everywhere

## Repository Workflow

- Run tests before commit
- Run lint before PR
- Use changesets
- Conventional commits

## Security

- Never expose secrets
- Never hardcode credentials
- Use environment contracts

## AI Agent Workflow

- Read architecture docs first
- Update tests with implementation
- Never bypass CI
- Never edit generated files
```

---

# Dev Environment (Highest Standard)

## Devcontainers

Modern elite repositories are now:

* reproducible
* containerized
* cloud-dev ready

Use:

* `.devcontainer/`
* GitHub Codespaces
* Docker-based development

Development container standards are increasingly mainstream. ([GitHub][3])

---

## `.devcontainer/`

```txt
.devcontainer/
├── devcontainer.json
├── Dockerfile
├── compose.yaml
└── scripts/
```

---

# Environment Strategy

Never use ad-hoc env management.

---

## Recommended

```txt
.env
.env.local
.env.development
.env.production
.env.test
.env.example
.env.schema
```

---

# Environment Validation

Use:

* `zod`
* `dotenv-safe`

Example:

```ts
export const env = envSchema.parse(process.env)
```

---

# Configuration Layer

Centralize ALL config.

```txt
packages/config/
├── eslint/
├── prettier/
├── typescript/
├── tailwind/
├── vitest/
├── playwright/
├── docker/
└── ci/
```

---

# Infrastructure-as-Code

Never separate infra from repo.

```txt
infrastructure/
├── terraform/
├── pulumi/
├── kubernetes/
├── helm/
├── cloudflare/
├── aws/
├── gcp/
└── monitoring/
```

---

# Deployment Structure

```txt
platform/
├── ci/
├── cd/
├── docker/
├── k8s/
├── releases/
├── secrets/
└── environments/
```

---

# CI/CD Standard

## `.github/`

```txt
.github/
├── workflows/
├── actions/
├── ISSUE_TEMPLATE/
├── PULL_REQUEST_TEMPLATE.md
├── CODEOWNERS
└── copilot-instructions.md
```

---

# Essential CI Pipelines

```txt
ci.yml
lint.yml
test.yml
typecheck.yml
security.yml
release.yml
deploy.yml
e2e.yml
```

---

# Git Standards

## Conventional commits

```txt
feat:
fix:
refactor:
docs:
test:
chore:
```

---

# Release Management

Use:

* Changesets
* Semantic Release

```txt
.changeset/
```

---

# Database Structure

```txt
packages/database/
├── prisma/
├── migrations/
├── seeds/
├── schema/
└── generated/
```

---

# API Architecture

## Strongest current approach

```txt
services/api/
├── src/
├── contracts/
├── modules/
├── domain/
├── application/
├── infrastructure/
└── tests/
```

Recommended:

* tRPC OR OpenAPI
* contract-first
* generated clients

---

# Observability (Mandatory)

Most repositories fail here.

You should build observability from day 1.

---

## Structure

```txt
packages/observability/
├── logging/
├── tracing/
├── metrics/
├── analytics/
├── alerts/
└── dashboards/
```

---

# Recommended Stack

| Concern    | Tool          |
| ---------- | ------------- |
| logs       | Pino          |
| metrics    | Prometheus    |
| tracing    | OpenTelemetry |
| dashboards | Grafana       |
| errors     | Sentry        |

---

# Security Structure

```txt
security/
├── policies/
├── threat-models/
├── audits/
├── scans/
└── secrets/
```

---

# Documentation Architecture

```txt
docs/
├── architecture/
├── adr/
├── api/
├── onboarding/
├── playbooks/
├── runbooks/
├── standards/
├── diagrams/
└── decisions/
```

---

# Architecture Decision Records

Very important for scaling teams.

```txt
docs/adr/
├── ADR-001-monorepo.md
├── ADR-002-auth.md
└── ADR-003-ai-agents.md
```

---

# Testing Strategy

```txt
tests/
├── unit/
├── integration/
├── e2e/
├── performance/
├── security/
├── contracts/
└── fixtures/
```

---

# Recommended Testing Stack

| Type      | Tool       |
| --------- | ---------- |
| unit      | Vitest     |
| e2e       | Playwright |
| mobile    | Detox      |
| API       | Supertest  |
| contracts | Pact       |

---

# Scripts Standard

```txt
scripts/
├── bootstrap/
├── setup/
├── release/
├── db/
├── infra/
├── ai/
├── codegen/
└── maintenance/
```

---

# Code Generation

Critical for scale.

```txt
tooling/
├── generators/
├── schematics/
├── templates/
└── codemods/
```

---

# Human-Friendly Standards

## Every repo should include

```txt
README.md
CONTRIBUTING.md
SECURITY.md
ARCHITECTURE.md
```

---

# README Structure

```md
# Project

## Architecture
## Getting Started
## Local Development
## Environment Variables
## Testing
## Deployment
## Repository Structure
## Troubleshooting
## AI Agent Instructions
```

---

# AI-Friendly Practices

These matter enormously now.

---

## Use predictable naming

GOOD:

```txt
auth.service.ts
user.repository.ts
create-user.usecase.ts
```

BAD:

```txt
stuff.ts
helpers.ts
misc.ts
```

---

# AI-Friendly Codebase Rules

## Mandatory

* small files
* single responsibility
* typed interfaces
* explicit boundaries
* no magic
* deterministic naming
* avoid giant classes
* avoid hidden side effects

---

# Scaling Standards

## Use domain boundaries

```txt
modules/
├── billing/
├── auth/
├── users/
└── notifications/
```

Not:

```txt
controllers/
services/
utils/
```

---

# Recommended Package Rules

Every package should contain:

```txt
src/
tests/
README.md
package.json
tsconfig.json
vitest.config.ts
```

---

# Mobile + Desktop Strategy

## Best modern unified setup

| Platform | Stack               |
| -------- | ------------------- |
| web      | Next.js             |
| mobile   | Expo                |
| desktop  | Tauri               |
| backend  | Node/Nest/Fastify   |
| shared   | TypeScript packages |

---

# Highest-Standard Build Pipeline

```txt
Developer
  ↓
Pre-commit hooks
  ↓
Lint + typecheck
  ↓
Unit tests
  ↓
Build graph
  ↓
Integration tests
  ↓
Preview deployment
  ↓
E2E tests
  ↓
Production deploy
  ↓
Observability validation
```

---

# Modern Repository Principles

## 1. Repository is the platform

Not just source code.

---

## 2. Everything-as-code

* infra
* CI
* environments
* policies
* observability
* AI rules

---

## 3. Reproducibility first

Every developer should get identical environments.

---

## 4. AI agents are first-class contributors

This is becoming standard practice across modern agentic tooling ecosystems. ([DeployHQ][2])

---

# Final Recommendation Stack (Best Overall)

| Layer           | Recommendation        |
| --------------- | --------------------- |
| Monorepo        | Turborepo             |
| Package manager | pnpm                  |
| Language        | TypeScript            |
| Web             | Next.js               |
| Mobile          | Expo                  |
| Desktop         | Tauri                 |
| API             | Fastify / NestJS      |
| DB              | PostgreSQL            |
| ORM             | Prisma                |
| Infra           | Terraform             |
| Containers      | Docker                |
| Dev env         | Devcontainers         |
| CI/CD           | GitHub Actions        |
| Observability   | OpenTelemetry         |
| Auth            | Better Auth / Auth.js |
| AI              | AGENTS.md + CLAUDE.md |
| Testing         | Vitest + Playwright   |
| Releases        | Changesets            |

---

# What Elite Repositories Optimize For

## Technical

* deterministic builds
* fast CI
* cacheability
* isolation
* typed contracts
* reproducibility

## Human

* discoverability
* onboarding
* maintainability
* documentation
* ownership clarity

## AI

* explicit context
* predictable structure
* modular boundaries
* instruction files
* codified workflows

---

# If You Want the Actual “Best-in-Class” Setup

The current top-tier pattern is approximately:

```txt
pnpm
+ turborepo
+ nextjs
+ expo
+ tauri
+ typescript
+ prisma
+ opentofu/terraform
+ github actions
+ devcontainers
+ AGENTS.md
+ OpenTelemetry
```

That is currently the closest thing to an industry “gold standard” for:

* AI-friendly
* scalable
* maintainable
* cross-platform
* enterprise-capable
* developer-friendly architecture.

[1]: https://turborepo.ai/docs/crafting-your-repository/structuring-a-repository?utm_source=chatgpt.com "Structuring a repository | Turborepo"
[2]: https://www.deployhq.com/blog/ai-coding-config-files-guide?utm_source=chatgpt.com "CLAUDE.md, AGENTS.md & Copilot Instructions: Configure Every AI Coding Assistant"
[3]: https://github.com/devcontainers?utm_source=chatgpt.com "devcontainers · GitHub"
