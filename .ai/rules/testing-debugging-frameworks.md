# Modern AI Application Engineering Playbook

## Testing, Debugging, CI/CD, Deployment, and Rollout Standards

Version: 1.0
Audience: AI agents, software engineers, platform engineers, technical leads
Scope: Modern AI repositories containing APIs, agents, workflows, retrieval systems, model orchestration, frontend apps, background workers, evaluation pipelines, and infrastructure.

---

# 1. Core Principles

## 1.1 Reliability Is a Product Feature

AI systems fail differently than traditional systems:

* Non-deterministic outputs
* Model drift
* Prompt regressions
* Retrieval instability
* Latency spikes
* Context window overflow
* Hallucinations
* Tool misuse
* Silent quality degradation
* Cost explosions
* Partial failures across chained systems

Testing and deployment standards must therefore validate:

1. Functional correctness
2. Behavioral correctness
3. Safety
4. Stability over time
5. Cost efficiency
6. Recovery behavior
7. Human trustworthiness

---

## 1.2 Every Change Must Define

Every feature, fix, prompt update, workflow modification, or infrastructure change must explicitly define:

| Area          | Required Definition                      |
| ------------- | ---------------------------------------- |
| Input         | What enters the system?                  |
| Output        | Expected result?                         |
| Context       | What information influences behavior?    |
| Impact        | What can break?                          |
| Dependencies  | Which services/models/tools are touched? |
| Failure Modes | How can it fail?                         |
| Observability | How do we detect failure?                |
| Rollback      | How do we revert safely?                 |
| Risk Level    | Low / Medium / High                      |

This is mandatory before implementation.

---

# 2. Universal Change Analysis Framework

Before building ANY feature:

## 2.1 Problem Contract

Every task must start with this structure.

```yaml
feature:
  name:
  owner:
  risk_level:

inputs:
  - user_prompt
  - uploaded_file
  - api_payload

outputs:
  - response_text
  - structured_json
  - side_effect

context:
  - memory
  - retrieval
  - user_settings
  - tool_outputs
  - model_parameters

dependencies:
  - llm_provider
  - vector_db
  - redis
  - workflow_engine

impact_surface:
  - frontend
  - backend_api
  - prompt_templates
  - evaluation_scores
  - latency
  - token_cost

failure_modes:
  - hallucination
  - timeout
  - invalid_json
  - prompt_injection
  - context_overflow
  - stale_retrieval

observability:
  metrics:
    - latency
    - token_usage
    - success_rate
    - evaluation_score

rollback:
  strategy:
  trigger_conditions:
```

This contract becomes:

* The implementation checklist
* The testing plan
* The CI validation basis
* The deployment approval reference
* The rollback decision framework

---

# 3. Repository Architecture Standards

## 3.1 Mandatory Repository Structure

```text
repo/
├── apps/
│   ├── frontend/
│   ├── api/
│   ├── workers/
│   └── agents/
│
├── packages/
│   ├── prompts/
│   ├── evaluations/
│   ├── shared/
│   ├── ai_core/
│   ├── orchestration/
│   └── sdk/
│
├── infra/
│   ├── terraform/
│   ├── kubernetes/
│   ├── monitoring/
│   └── ci/
│
├── tests/
│   ├── unit/
│   ├── integration/
│   ├── e2e/
│   ├── regression/
│   ├── safety/
│   ├── load/
│   ├── evaluations/
│   └── fixtures/
│
├── datasets/
│   ├── golden/
│   ├── adversarial/
│   └── benchmarks/
│
├── docs/
├── scripts/
└── .github/
```

---

## 3.2 AI-Specific Requirements

### Prompts Must Be Versioned

Never inline prompts in business logic.

Bad:

```python
response = llm("You are helpful...")
```

Good:

```python
prompt = prompt_registry.load("summarizer/v3")
response = llm(prompt)
```

Every prompt change must:

* Have semantic versioning
* Have evaluation scores
* Have regression comparisons
* Be rollbackable

---

## 3.3 Deterministic Interfaces

Even if the model is non-deterministic:

* APIs must remain deterministic
* Schemas must remain deterministic
* Tool interfaces must remain deterministic
* Workflow state must remain deterministic

Always enforce:

* JSON schemas
* Type validation
* Structured outputs
* Retry policies
* Timeout policies
* Fallback strategies

---

# 4. Testing Pyramid for AI Systems

Traditional testing pyramids are insufficient.

Modern AI systems require:

```text
                Business Evaluation
              / Safety & Trust Layer \
           End-to-End Behavioral Tests
         Workflow + Integration Testing
        Model & Retrieval Evaluations
             Unit & Contract Tests
```

---

# 5. What MUST Be Tested

# 5.1 Unit Tests

Purpose:
Validate isolated logic.

Required Coverage:

| Component        | Must Test                      |
| ---------------- | ------------------------------ |
| Parsers          | Edge cases, malformed input    |
| Validators       | Schema failures                |
| Utilities        | Pure function correctness      |
| Prompt builders  | Variable injection correctness |
| Tokenizers       | Token limits                   |
| Tool wrappers    | Error handling                 |
| Retry logic      | Backoff behavior               |
| Guards           | Unsafe input rejection         |
| Cache layers     | Hit/miss logic                 |
| Cost calculators | Budget enforcement             |

Target:

* 90%+ critical-path coverage
* 70%+ repository-wide minimum

What can be skipped:

* Trivial getters/setters
* Framework internals
* Auto-generated code
* Pure UI styling

---

# 5.2 Contract Tests

Critical for AI systems.

Validate:

* API schemas
* Model output schemas
* Tool input/output compatibility
* Event formats
* Queue payloads
* Webhook contracts

Mandatory:

```json
{
  "type": "object",
  "required": ["answer", "confidence"]
}
```

AI systems frequently fail through schema drift.

Contract tests prevent silent breakage.

---

# 5.3 Prompt Regression Tests

Every prompt change must compare:

| Metric             | Required |
| ------------------ | -------- |
| Accuracy           | Yes      |
| Tone consistency   | Yes      |
| JSON validity      | Yes      |
| Hallucination rate | Yes      |
| Cost               | Yes      |
| Latency            | Yes      |
| Safety score       | Yes      |

Test against:

* Golden datasets
* Edge cases
* Adversarial prompts
* Long-context prompts
* Multi-language prompts
* Injection attempts

Never merge prompt changes without regression results.

---

# 5.4 Retrieval Tests (RAG)

RAG systems require independent validation.

Must test:

| Area             | Test                      |
| ---------------- | ------------------------- |
| Chunking         | Boundary correctness      |
| Embeddings       | Semantic consistency      |
| Recall           | Relevant docs returned    |
| Precision        | Irrelevant docs filtered  |
| Ranking          | Correct ordering          |
| Freshness        | New docs indexed          |
| Metadata filters | Correct filtering         |
| Access control   | Unauthorized docs blocked |

Metrics:

* Recall@K
* MRR
* NDCG
* Retrieval latency

---

# 5.5 Agent Workflow Tests

Agent systems fail through orchestration.

Test:

| Area               | Examples                    |
| ------------------ | --------------------------- |
| Tool selection     | Correct tool chosen         |
| Multi-step plans   | Stable execution            |
| State management   | Correct memory updates      |
| Retry behavior     | Recovers gracefully         |
| Loop prevention    | Infinite loops blocked      |
| Context growth     | Window stays bounded        |
| Tool hallucination | Invalid tool calls rejected |
| Parallel tasks     | Race conditions avoided     |

Must simulate:

* Partial failures
* Network failures
* Tool timeouts
* Invalid tool responses
* Missing context
* Corrupted memory

---

# 5.6 End-to-End Behavioral Tests

These validate real workflows.

Examples:

```text
User uploads PDF
→ OCR extracts text
→ Retrieval indexes content
→ Agent answers questions
→ Citations generated
→ Response streamed to frontend
```

Validate:

* Final user experience
* Latency budgets
* Streaming correctness
* Multi-service compatibility
* Persistence behavior
* Error recovery

Required:

* Happy paths
* Expected failure paths
* Adversarial paths

---

# 5.7 Safety Tests

Mandatory for all AI systems.

Test categories:

| Category             | Examples                      |
| -------------------- | ----------------------------- |
| Prompt injection     | Ignore malicious instructions |
| Data leakage         | No secret exposure            |
| Jailbreak resistance | Unsafe requests blocked       |
| PII handling         | Redaction works               |
| Toxicity             | Harmful outputs reduced       |
| Authorization        | User isolation enforced       |
| Tool abuse           | Restricted operations blocked |
| File handling        | Malicious uploads rejected    |

Safety suites must run:

* On every merge request
* Before production deployment
* Nightly against production snapshots

---

# 5.8 Load & Scalability Tests

AI systems often fail under:

* Token spikes
* Concurrent inference
* Retrieval pressure
* Queue saturation
* GPU exhaustion

Must test:

| Area             | Metrics       |
| ---------------- | ------------- |
| API throughput   | RPS           |
| Queue latency    | Time in queue |
| GPU utilization  | Saturation    |
| Model latency    | p50/p95/p99   |
| Token throughput | Tokens/sec    |
| Database load    | Query time    |
| Cache efficiency | Hit ratio     |

Mandatory thresholds:

* p95 latency target
* Maximum queue delay
* Maximum error rate
* Cost-per-request threshold

---

# 5.9 Chaos & Failure Testing

Advanced but critical.

Inject failures:

* LLM outage
* Vector DB outage
* Redis outage
* Partial network partition
* Slow external APIs
* Corrupted retrieval results
* Invalid model responses

Validate:

* Graceful degradation
* Fallback models
* Retry strategies
* Circuit breakers
* User-visible recovery behavior

---

# 6. What Can Be Skipped

Avoid over-testing low-value areas.

## Usually Safe to Skip

| Area                   | Why                  |
| ---------------------- | -------------------- |
| Framework internals    | Owned externally     |
| CSS details            | Low-risk             |
| Generated code         | Deterministic        |
| Stable libraries       | Covered upstream     |
| Static configs         | Low change frequency |
| Non-critical analytics | No user impact       |

---

## Never Skip

| Area                | Why                 |
| ------------------- | ------------------- |
| Authentication      | Security risk       |
| Authorization       | Data exposure       |
| Prompt changes      | Behavioral drift    |
| Retrieval logic     | Silent degradation  |
| Billing logic       | Financial impact    |
| Agent orchestration | System instability  |
| Safety guards       | Compliance risk     |
| Schema validation   | Runtime instability |
| Rollback validation | Recovery failure    |

---

# 7. Debugging Standards

# 7.1 Structured Debugging Framework

Every issue investigation should answer:

| Question               | Example               |
| ---------------------- | --------------------- |
| What changed?          | New prompt deployed   |
| When started?          | After deployment v1.4 |
| Who affected?          | Enterprise users only |
| Input dependent?       | Long prompts only     |
| Model dependent?       | GPT-5 only            |
| Environment dependent? | Production only       |
| Deterministic?         | Intermittent          |

---

# 7.2 Required Observability

Every AI request should generate:

```json
{
  "request_id": "",
  "trace_id": "",
  "user_id": "",
  "model": "",
  "prompt_version": "",
  "retrieval_ids": [],
  "tool_calls": [],
  "latency_ms": 0,
  "token_usage": {},
  "cost": 0,
  "evaluation_scores": {},
  "errors": []
}
```

---

# 7.3 Logging Standards

Must log:

* Request lifecycle
* Tool execution
* Retrieval results
* Retry attempts
* Model failures
* Validation failures
* Deployment metadata

Must NOT log:

* Secrets
* Full PII
* Authentication tokens
* Sensitive enterprise data

---

# 7.4 Distributed Tracing

Required for multi-agent systems.

Every workflow must preserve:

* Correlation IDs
* Parent-child spans
* Tool execution spans
* Model invocation spans
* Queue spans

Use:

* OpenTelemetry
* Jaeger
* Datadog
* Honeycomb
* Grafana Tempo

---

# 8. Evaluation Systems

# 8.1 AI Systems Need Continuous Evaluation

Traditional pass/fail testing is insufficient.

Need:

* Human evaluations
* Model-based evaluations
* Business KPI evaluations
* Safety evaluations
* Cost evaluations

---

# 8.2 Golden Datasets

Maintain curated datasets:

| Dataset              | Purpose                 |
| -------------------- | ----------------------- |
| Golden               | Expected ideal outputs  |
| Edge                 | Rare scenarios          |
| Adversarial          | Attack attempts         |
| Regression           | Previously broken cases |
| Production snapshots | Real-world behavior     |

Every production bug must add:

1. Reproduction case
2. Regression test
3. Evaluation example

---

# 8.3 Evaluation Gates

Deployment blocked if:

| Metric                  | Example Threshold  |
| ----------------------- | ------------------ |
| Accuracy drops          | >2%                |
| Hallucination increases | >1%                |
| Cost increases          | >10%               |
| Latency increases       | >15%               |
| Safety score drops      | Any critical issue |
| JSON validity drops     | <99.5%             |

---

# 9. CI Standards

# 9.1 CI Pipeline Stages

```text
1. Lint
2. Type Check
3. Unit Tests
4. Contract Tests
5. Security Scan
6. Prompt Regression
7. Integration Tests
8. E2E Tests
9. Evaluation Suite
10. Build Artifacts
11. Deployment Validation
```

---

# 9.2 Fast vs Slow CI

## Fast CI (<10 minutes)

Runs on every commit:

* Lint
* Type checks
* Unit tests
* Contract tests
* Small regression suite

## Slow CI

Runs:

* Nightly
* Before release
* Before production deployment

Includes:

* Full evaluation suite
* Load tests
* Adversarial tests
* Chaos tests
* Cross-model comparisons

---

# 9.3 Required CI Policies

Mandatory:

| Policy              | Requirement      |
| ------------------- | ---------------- |
| Branch protection   | Enabled          |
| Required reviews    | 2+ for high-risk |
| Required passing CI | Yes              |
| Signed commits      | Recommended      |
| Secret scanning     | Mandatory        |
| Dependency scanning | Mandatory        |
| SBOM generation     | Recommended      |

---

# 9.4 AI-Specific CI Requirements

Must validate:

* Prompt diffs
* Token cost changes
* Model configuration changes
* Safety policy updates
* Retrieval index migrations
* Embedding model changes

---

# 10. CD & Deployment Standards

# 10.1 Deployment Principles

Goals:

1. Small deployments
2. Reversible deployments
3. Observable deployments
4. Gradual rollouts
5. Safe defaults

---

# 10.2 Environment Strategy

Minimum:

```text
local
→ development
→ staging
→ production
```

Recommended:

```text
local
→ ephemeral preview
→ development
→ integration
→ staging
→ canary
→ production
```

---

# 10.3 Infrastructure as Code

Everything versioned:

* Kubernetes
* Terraform
* Secrets references
* Monitoring
* Queues
* Databases
* GPU configurations

Never manually modify production infrastructure.

---

# 10.4 Immutable Deployments

Rules:

* Versioned containers only
* No mutable production servers
* Reproducible builds
* Artifact signing
* Provenance tracking

---

# 11. Rollout Strategies

# 11.1 Progressive Delivery

Never instantly deploy AI behavior changes globally.

Use:

| Strategy        | Use Case                  |
| --------------- | ------------------------- |
| Canary          | High-risk changes         |
| Blue/Green      | Infrastructure migrations |
| Feature flags   | Experimental features     |
| Shadow traffic  | Silent evaluation         |
| A/B testing     | Behavioral comparison     |
| Ring deployment | Enterprise rollout        |

---

# 11.2 Recommended Rollout Sequence

```text
1% internal
→ 5% beta users
→ 10% low-risk traffic
→ 25%
→ 50%
→ 100%
```

Advance only if:

* Error rate stable
* Cost stable
* Evaluation scores stable
* Latency acceptable
* No safety regressions

---

# 11.3 Shadow Deployments

Extremely important for AI.

New model/prompt runs silently:

```text
Production Request
├── Current System → User
└── Candidate System → Evaluation Only
```

Compare:

* Quality
* Cost
* Latency
* Safety
* Hallucination rate

Before enabling publicly.

---

# 11.4 Feature Flags

All AI behaviors should be flaggable.

Feature flags required for:

* New prompts
* New models
* New tools
* New retrieval systems
* Agent autonomy changes
* Safety systems

Must support:

* Instant disable
* User segmentation
* Gradual rollout
* Environment-specific config

---

# 12. Rollback Standards

# 12.1 Every Deployment Needs Rollback

Rollback must be:

* Tested
* Automated
* Fast
* Observable

---

# 12.2 Rollback Triggers

Automatic rollback if:

| Metric           | Threshold          |
| ---------------- | ------------------ |
| Error rate       | Spike              |
| Safety violation | Critical           |
| Latency          | Severe degradation |
| Cost             | Unexpected surge   |
| Hallucination    | Major increase     |
| Availability     | SLO breach         |

---

# 12.3 Rollback Types

| Type             | Example               |
| ---------------- | --------------------- |
| Prompt rollback  | Revert prompt version |
| Model rollback   | Revert model provider |
| Infra rollback   | Revert deployment     |
| Feature rollback | Disable flag          |
| Traffic rollback | Shift routing         |

---

# 13. Production Monitoring

# 13.1 Core Metrics

## Reliability

* Availability
* Error rate
* Timeout rate
* Retry rate

## AI Quality

* Hallucination score
* Evaluation score
* User feedback
* Citation accuracy

## Performance

* p50/p95/p99 latency
* Tokens/sec
* Queue latency
* GPU utilization

## Cost

* Cost/request
* Cost/user
* Daily spend
* Token consumption

---

# 13.2 Alerting Standards

Alerts must be:

* Actionable
* Prioritized
* Routed correctly
* Deduplicated

Never alert on noisy metrics alone.

Use SLO-based alerting.

---

# 14. Security Standards

# 14.1 AI-Specific Threats

Must defend against:

* Prompt injection
* Data exfiltration
* Tool abuse
* Agent privilege escalation
* Retrieval poisoning
* Training data leakage
* Model output abuse

---

# 14.2 Secret Management

Never:

* Hardcode secrets
* Store secrets in prompts
* Log credentials
* Put API keys in frontend

Use:

* Vault
* AWS Secrets Manager
* GCP Secret Manager
* Kubernetes secrets

---

# 15. Operational Excellence

# 15.1 Incident Response

Every incident requires:

1. Timeline
2. Root cause
3. Blast radius
4. Mitigation
5. Regression prevention
6. Monitoring improvements
7. Runbook updates

---

# 15.2 Runbooks

Required for:

* Model outages
* Queue saturation
* Deployment rollback
* Retrieval corruption
* GPU exhaustion
* Token spikes
* Safety incidents

---

# 16. Recommended Tooling Stack

# 16.1 CI/CD

| Category      | Recommended               |
| ------------- | ------------------------- |
| CI            | GitHub Actions, GitLab CI |
| CD            | ArgoCD, FluxCD            |
| IaC           | Terraform                 |
| Containers    | Docker                    |
| Orchestration | Kubernetes                |
| Feature Flags | LaunchDarkly, OpenFeature |

---

# 16.2 Observability

|
