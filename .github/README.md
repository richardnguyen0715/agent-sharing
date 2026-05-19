# .github — AI-Native Configuration

This directory contains all AI agent configuration for GitHub Copilot and compatible tools.

---

## Structure

```
.github/
├── copilot-instructions.md   ← Mandatory session bootstrap + operational rules
├── agents/                   ← Specialized domain agents (21 agents)
│   ├── README.md             ← Agent index with categories
│   ├── engineering/          ← backend, frontend, fullstack, python-pro, typescript-pro
│   ├── architecture/         ← api-designer, api-documenter, database-administrator
│   ├── quality/              ← code-reviewer, test-engineer, test-automator, security-auditor, performance-engineer
│   ├── devops/               ← devops-engineer, docker-expert, terraform-engineer
│   ├── ai-ml/                ← ml-engineer
│   ├── knowledge/            ← repo-analyst, context-curator, reference-navigator
│   └── debug/                ← debugger
├── skills/                   ← Workflow skills (30 skills)
│   ├── README.md             ← Skill index with categories
│   └── <skill-name>/
│       └── SKILL.md
└── agent/                    ← Copilot cloud agent setup (copilot-setup-steps.yml)
```

---

## Quick Reference

### Use a skill
```
"Use the spec-driven-development skill"
"Use the planning-and-task-breakdown skill"
"Use the debugging-and-error-recovery skill"
```

### Use an agent
```
/agent backend-developer
/agent code-reviewer
/agent debugger
```

---

## Maintenance

- **Add a skill**: Create `skills/<name>/SKILL.md` → update `skills/README.md`
- **Add an agent**: Create `agents/<category>/<name>.md` → update `agents/README.md`
- **Update instructions**: Edit `copilot-instructions.md`
