# Web-Based Multi-User Terminal Platform

## System Specification & Architecture Design

---

# 1. Project Goal

Xây dựng một nền tảng web cho phép nhiều người dùng truy cập và sử dụng nhiều terminal khác nhau trên cùng một máy chủ đã được đăng nhập sẵn bằng một tài khoản dùng chung (ví dụ tài khoản công ty).

Hệ thống:

* không phải SSH gateway truyền thống
* không virtualize machine hoàn toàn
* không container-per-user
* không multi-tenant security-grade isolation

Mà là:

> Logical terminal isolation managed entirely by the web application layer.

---

# 2. Core Use Case

Ví dụ:

* Máy chủ đã login:

  * GitHub
  * AWS
  * VPN công ty
  * Internal tools
  * SSH company account

* Gia đình có 3 người dùng:

  * User A
  * User B
  * User C

* Mỗi người:

  * có terminal riêng
  * không thấy terminal người khác
  * không attach session người khác
  * có quyền riêng
  * cùng sử dụng shared authenticated environment

---

# 3. System Philosophy

## Shared Runtime Identity

Toàn bộ terminal:

* chạy cùng Linux user
* dùng cùng credentials
* dùng cùng environment
* dùng cùng login state

Nhưng:

## Application-Level Isolation

Isolation được quản lý bởi:

* web backend
* terminal ownership
* session routing
* RBAC
* policy engine

KHÔNG phải:

* Linux users
* containers
* VM isolation

---

# 4. High-Level Architecture

```text id="g4x07y"
┌──────────────────────────────┐
│           Browser            │
│  React + xterm.js frontend   │
└──────────────┬───────────────┘
               │
        HTTPS / WebSocket
               │
┌──────────────▼───────────────┐
│         Web Backend          │
│                              │
│ - Authentication             │
│ - RBAC                       │
│ - Terminal Manager           │
│ - Session Router             │
│ - Policy Engine              │
│ - Audit Logging              │
└──────────────┬───────────────┘
               │
       PTY / Process Control
               │
┌──────────────▼───────────────┐
│       Shared Linux User      │
│                              │
│  PTY #1 -> User A            │
│  PTY #2 -> User A            │
│  PTY #3 -> User B            │
│  PTY #4 -> User C            │
└──────────────────────────────┘
```

---

# 5. Core Concepts

---

# 5.1 Terminal Session

Mỗi terminal là:

* một PTY riêng
* shell process riêng
* process tree riêng
* working directory riêng

Ví dụ:

```text id="l33ltx"
Terminal #12
 ├── bash
 ├── vim
 └── node server.js
```

---

# 5.2 Shared Host Identity

Tất cả terminal chạy dưới:

```text id="h9o8e3"
company_user
```

Environment được share:

```bash id="8tm6fz"
AWS_PROFILE=company
GITHUB_TOKEN=...
VPN_ACTIVE=true
```

---

# 5.3 Logical Ownership

Terminal ownership được quản lý bởi backend.

Ví dụ:

| Terminal | Owner  |
| -------- | ------ |
| T1       | User A |
| T2       | User A |
| T3       | User B |

User:

* không thể attach terminal người khác
* không thể list terminal người khác

---

# 6. System Components

---

# 6.1 Frontend

## Responsibilities

* terminal rendering
* websocket communication
* authentication
* terminal tabs
* resize handling
* admin dashboard

---

## Recommended Stack

| Layer            | Technology       |
| ---------------- | ---------------- |
| Framework        | React            |
| Build Tool       | Vite             |
| Terminal UI      | xterm.js         |
| State Management | Zustand          |
| UI Library       | TailwindCSS      |
| Realtime         | Native WebSocket |

---

## Terminal Rendering

Sử dụng:

[xterm.js](https://xtermjs.org/?utm_source=chatgpt.com)

Flow:

```text id="zhntn9"
Keyboard Input
 -> websocket
 -> backend
 -> PTY stdin

PTY stdout
 -> websocket
 -> xterm.js render
```

---

# 6.2 Backend

## Responsibilities

### Authentication

* login
* session management
* JWT

### Authorization

* RBAC
* terminal ownership
* admin permissions

### Terminal Management

* create PTY
* kill PTY
* resize PTY
* reconnect PTY

### WebSocket Gateway

* stream PTY I/O
* route messages

### Policy Engine

* optional command filtering
* heuristics
* warning system

### Audit Logging

* session history
* commands
* events

---

## Recommended Stack

| Layer     | Technology        |
| --------- | ----------------- |
| Language  | Go                |
| Router    | chi / gin         |
| WebSocket | gorilla/websocket |
| PTY       | creack/pty        |
| Auth      | JWT               |
| Database  | PostgreSQL        |
| Cache     | Redis             |

---

# 6.3 PTY Runtime Layer

## Responsibilities

* spawn shell
* maintain PTY
* manage process lifecycle

---

## Recommended PTY Library

### Go

[creack/pty](https://github.com/creack/pty?utm_source=chatgpt.com)

Example:

```go id="fyxvr0"
cmd := exec.Command("bash")
cmd.Dir = "/workspace/user_a"

ptyFile, err := pty.Start(cmd)
```

---

# 7. Session Lifecycle

---

# 7.1 User Login

```text id="c8wwep"
User Login
 -> JWT issued
 -> permissions loaded
```

---

# 7.2 Terminal Creation

```text id="r4m9kn"
POST /terminal/create
```

Backend:

* verifies permissions
* allocates PTY
* spawns shell
* stores ownership

---

# 7.3 WebSocket Attach

```text id="bj9f1m"
Browser
 -> websocket connect
 -> backend validates ownership
 -> attach PTY stream
```

---

# 7.4 Terminal Streaming

```text id="4duqvq"
stdin
 -> websocket
 -> PTY

stdout/stderr
 -> PTY
 -> websocket
 -> browser
```

---

# 8. Permission System

---

# 8.1 RBAC

## Roles

| Role       | Description                      |
| ---------- | -------------------------------- |
| admin      | full access                      |
| operator   | create/manage assigned terminals |
| viewer     | readonly                         |
| restricted | limited commands/features        |

---

# 8.2 Permission Examples

| Permission             | Description       |
| ---------------------- | ----------------- |
| terminal:create        | create terminal   |
| terminal:delete        | delete terminal   |
| terminal:view          | attach terminal   |
| terminal:write         | send input        |
| terminal:resize        | resize session    |
| session:recording:view | access recordings |

---

# 8.3 Terminal Ownership

Every terminal has:

```json id="hm6r8f"
{
  "terminal_id": "t_123",
  "owner_id": "user_a"
}
```

---

# 9. Filesystem Strategy

Không hard isolate filesystem.

Chỉ logical separation:

```text id="jyzgmk"
/workspace/user_a
/workspace/user_b
/workspace/shared
```

Shell startup:

```bash id="t0kv78"
cd /workspace/user_a
```

---

# 10. Command Governance

---

# IMPORTANT

System không thể secure shell hoàn toàn bằng command parsing.

Shell là:

* interactive
* stateful
* programmable

User có shell thật:

* eventually có arbitrary execution

---

# Recommended Strategy

## Soft Governance

### Allow:

* monitoring
* logging
* warnings
* heuristics

### Avoid:

* pretending command parsing is real security

---

# Example Heuristic Blocks

| Pattern  | Action |
| -------- | ------ |
| rm -rf / | block  |
| shutdown | warn   |
| reboot   | block  |
| mkfs     | block  |

---

# 11. Security Model

---

# Security Assumption

Users:

* semi-trusted
* not hostile attackers

System:

* UX isolation
* governance
* auditing

NOT:

* hardened multi-tenant sandbox

---

# 12. Database Design

---

# users

```sql id="7r9n4e"
users (
  id,
  username,
  password_hash,
  role,
  created_at
)
```

---

# terminals

```sql id="v5jlwm"
terminals (
  id,
  owner_id,
  title,
  cwd,
  status,
  created_at
)
```

---

# terminal_sessions

```sql id="8bajd9"
terminal_sessions (
  id,
  terminal_id,
  websocket_id,
  connected_at
)
```

---

# audit_logs

```sql id="4mha4h"
audit_logs (
  id,
  user_id,
  terminal_id,
  action,
  metadata,
  created_at
)
```

---

# 13. Realtime Protocol

---

# WebSocket Messages

## stdin

```json id="av7tzd"
{
  "type": "stdin",
  "data": "ls -la\n"
}
```

---

## resize

```json id="r3os2t"
{
  "type": "resize",
  "cols": 120,
  "rows": 40
}
```

---

## output

```json id="w0q48v"
{
  "type": "output",
  "data": "total 12\n"
}
```

---

# 14. Terminal Manager Design

---

# Core Structure

```go id="nsjlwm"
type Terminal struct {
    ID        string
    OwnerID   string
    PTY       *os.File
    Cmd       *exec.Cmd
    CreatedAt time.Time
}
```

---

# Terminal Registry

```go id="2o7un3"
map[string]*Terminal
```

---

# 15. Session Recovery

System nên support:

* reconnect websocket
* persistent PTY
* disconnected browser recovery

Ví dụ:

```text id="4st0es"
Browser disconnect
 -> PTY still running
 -> reconnect later
```

---

# 16. Audit & Recording

---

# Recommended Features

| Feature                 | Importance |
| ----------------------- | ---------- |
| command logging         | high       |
| session recording       | high       |
| connect/disconnect logs | high       |
| replay support          | medium     |

---

# Recording Formats

Có thể dùng:

* ttyrec
* asciinema-style format

---

# 17. Deployment Architecture

---

# MVP Deployment

```text id="eqbntr"
Single Linux Server
 ├── frontend
 ├── backend
 ├── PostgreSQL
 └── PTY runtime
```

---

# Recommended OS

| OS            | Recommendation |
| ------------- | -------------- |
| Ubuntu Server | recommended    |
| Debian        | recommended    |

---

# 18. Scalability Roadmap

---

# Phase 1 — MVP

* single server
* PTY sessions
* websocket streaming
* RBAC
* ownership isolation

---

# Phase 2

* session recording
* reconnect
* audit search

---

# Phase 3

* multi-node runtime
* distributed PTY agents
* centralized orchestration

---

# 19. Non-Goals

System KHÔNG nhằm:

* secure hostile attackers
* replace virtualization
* replace Docker/Kubernetes
* provide military-grade isolation
* sandbox arbitrary malicious users

---

# 20. Recommended Final Tech Stack

| Layer       | Technology        |
| ----------- | ----------------- |
| Frontend    | React             |
| Terminal UI | xterm.js          |
| Styling     | TailwindCSS       |
| Backend     | Go                |
| WebSocket   | gorilla/websocket |
| PTY         | creack/pty        |
| Auth        | JWT               |
| Database    | PostgreSQL        |
| Cache       | Redis             |
| OS          | Ubuntu Server     |

---

# 21. Final Architecture Summary

Hệ thống là:

> A web-based terminal virtualization and governance platform running on top of a shared authenticated Linux environment.

Core properties:

* multiple isolated PTY sessions
* shared runtime identity
* application-level ownership
* centralized access control
* realtime terminal streaming
* auditability
* browser-native UX

Nó KHÔNG phải:

* SSH wrapper
* VM platform
* container orchestration system

Mà là:

> Logical terminal isolation and management over a shared authenticated machine context.

# 22. Enterprise Integration & Security Extension

---

# 22.1 External Service Governance

## Objective

Mọi tương tác với:

* AI APIs
* MCP services
* Databricks
* internal APIs
* cloud services
* external execution systems

đều phải:

* được audit
* traceable
* observable
* asynchronous
* centralized

Thông qua:

* Kafka event pipeline

---

# 22.2 Architecture Extension

```text id="1oz0z7"
┌──────────────────────────┐
│        Browser           │
└────────────┬─────────────┘
             │
        HTTPS/WebSocket
             │
┌────────────▼─────────────┐
│       API Gateway        │
└────────────┬─────────────┘
             │
┌────────────▼─────────────┐
│     Backend Services     │
│                          │
│ - Auth                   │
│ - RBAC                   │
│ - Terminal Manager       │
│ - External Integrations  │
│ - Audit Engine           │
└──────┬───────────┬───────┘
       │           │
       │           │
       │      Kafka Producer
       │           │
       │           ▼
       │     Kafka Cluster
       │           │
       │           ▼
       │     Event Consumers
       │
       ▼
External Services
(AI/MCP/Databricks/etc)
```

---

# 22.3 Kafka Integration Requirements

## Every External Interaction Must Emit Events

Ví dụ:

| Service Type     | Example                |
| ---------------- | ---------------------- |
| AI API           | OpenAI / Claude        |
| MCP              | model context protocol |
| Data Platform    | Databricks             |
| External Compute | Spark                  |
| Internal API     | company systems        |
| Cloud Service    | AWS/GCP/Azure          |

---

# 22.4 Event Flow

## Example

```text id="l0n6cz"
User Request
 -> Backend Authorization
 -> Kafka Event Publish
 -> External API Call
 -> Kafka Response Event
 -> Audit Persistence
```

---

# 22.5 Kafka Topics

## Recommended Topics

| Topic            | Purpose               |
| ---------------- | --------------------- |
| terminal.events  | terminal lifecycle    |
| terminal.audit   | audit logs            |
| ai.requests      | AI requests           |
| ai.responses     | AI responses          |
| mcp.events       | MCP interactions      |
| databricks.jobs  | job execution         |
| security.events  | security violations   |
| websocket.events | connection telemetry  |
| auth.events      | authentication events |

---

# 22.6 Event Schema

## Standard Event Envelope

```json id="r3dz87"
{
  "event_id": "uuid",
  "event_type": "ai.request",
  "timestamp": "2026-05-19T10:00:00Z",
  "user_id": "user_a",
  "terminal_id": "t_001",
  "service": "openai",
  "metadata": {},
  "payload": {}
}
```

---

# 22.7 AI Request Event

```json id="hq4zrl"
{
  "event_type": "ai.request",
  "model": "gpt-5.5",
  "tokens": 1200,
  "workspace": "company_a",
  "request_id": "req_123"
}
```

---

# 22.8 Databricks Event

```json id="2kmc2n"
{
  "event_type": "databricks.job.run",
  "job_id": "db_job_1",
  "cluster": "analytics_cluster",
  "initiator": "user_b"
}
```

---

# 23. Security Architecture

---

# 23.1 Security Goals

System phải:

* production deployable
* enterprise-safe
* audit-capable
* resilient
* observable

---

# 23.2 Threat Model

## Protected Against

| Threat                       | Protection               |
| ---------------------------- | ------------------------ |
| Unauthorized terminal access | RBAC                     |
| Session hijacking            | JWT/session validation   |
| WebSocket abuse              | token verification       |
| Command abuse                | governance + audit       |
| Data leakage                 | ownership isolation      |
| Replay attacks               | expiring tokens          |
| Horizontal access            | backend ownership checks |

---

# 23.3 Not Protected Against

| Threat                       | Reason                     |
| ---------------------------- | -------------------------- |
| Full hostile shell attackers | shared runtime model       |
| Kernel privilege escalation  | no VM/container isolation  |
| Shared credential abuse      | intentional shared context |

---

# 23.4 Authentication

## Recommended Stack

| Component      | Technology |
| -------------- | ---------- |
| Authentication | Keycloak   |
| Session        | JWT        |
| Refresh Tokens | enabled    |
| MFA            | optional   |
| OAuth2         | supported  |

---

## Recommended Identity Provider

[Keycloak](https://www.keycloak.org/?utm_source=chatgpt.com)

---

# 23.5 Authorization Layer

Every request:

* validated
* ownership checked
* permission checked

Example:

```text id="b8oh1m"
User
 -> request terminal attach
 -> backend validates owner_id
 -> websocket upgrade allowed
```

---

# 23.6 WebSocket Security

## Requirements

### JWT validation during upgrade

```text id="6g9yfp"
GET /ws?token=...
```

---

### Session ownership validation

```text id="ckv7rh"
terminal.owner_id == user.id
```

---

### Connection timeout

| Type           | Value |
| -------------- | ----- |
| Idle timeout   | 30m   |
| Session expiry | 24h   |

---

# 23.7 Backend Security

---

# Mandatory Security Headers

## Via NGINX

```nginx id="c8x6bp"
add_header X-Frame-Options DENY;
add_header X-Content-Type-Options nosniff;
add_header Referrer-Policy no-referrer;
add_header Content-Security-Policy "default-src 'self'";
```

---

# Rate Limiting

```nginx id="ybv6rj"
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
```

---

# CORS Restrictions

Only trusted frontend origins.

---

# 23.8 Secret Management

## NEVER

* hardcode API keys
* store plaintext secrets
* expose env vars to frontend

---

# Recommended

| Secret Type    | Storage            |
| -------------- | ------------------ |
| API Keys       | Kubernetes Secrets |
| JWT Secrets    | Vault              |
| DB Credentials | Vault              |
| TLS Keys       | cert-manager       |

---

# Recommended Vault

[HashiCorp Vault](https://www.vaultproject.io/?utm_source=chatgpt.com)

---

# 23.9 Audit & Compliance

Every critical action must be logged.

---

# Logged Events

| Event              | Required |
| ------------------ | -------- |
| login              | yes      |
| websocket attach   | yes      |
| terminal creation  | yes      |
| terminal deletion  | yes      |
| external API usage | yes      |
| AI requests        | yes      |
| permission denial  | yes      |

---

# 24. Observability Architecture

---

# 24.1 Metrics

## Recommended Stack

| Component     | Technology    |
| ------------- | ------------- |
| Metrics       | Prometheus    |
| Visualization | Grafana       |
| Logs          | Loki          |
| Tracing       | OpenTelemetry |

---

# 24.2 Required Metrics

| Metric                 | Description      |
| ---------------------- | ---------------- |
| active_terminals       | active PTYs      |
| websocket_connections  | concurrent users |
| ai_requests_total      | AI usage         |
| kafka_publish_failures | Kafka failures   |
| terminal_spawn_latency | PTY startup time |
| auth_failures          | failed auth      |

---

# 24.3 Distributed Tracing

All external calls:

* traceable
* correlated
* request-scoped

Example:

```text id="9e8d0l"
Request ID
 -> websocket session
 -> AI request
 -> Kafka event
 -> audit log
```

---

# 25. Production Deployment Architecture

---

# 25.1 Infrastructure Layout

```text id="vs0j1o"
                Internet
                    │
              ┌─────▼─────┐
              │   NGINX   │
              └─────┬─────┘
                    │
        ┌───────────┼───────────┐
        │           │           │
 ┌──────▼─────┐ ┌───▼────┐ ┌────▼─────┐
 │ Frontend   │ │Backend │ │ WebSocket│
 └──────┬─────┘ └───┬────┘ └────┬─────┘
        │           │            │
        └───────────┼────────────┘
                    │
             ┌──────▼──────┐
             │ PostgreSQL  │
             └──────┬──────┘
                    │
             ┌──────▼──────┐
             │    Redis    │
             └──────┬──────┘
                    │
             ┌──────▼──────┐
             │    Kafka    │
             └─────────────┘
```

---

# 25.2 Docker Deployment

## Every Service Containerized

| Service   | Container      |
| --------- | -------------- |
| frontend  | nginx + static |
| backend   | go service     |
| websocket | go service     |
| postgres  | postgres       |
| redis     | redis          |
| kafka     | kafka          |

---

# Example Backend Dockerfile

```dockerfile id="djlwm4"
FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go build -o app .

FROM ubuntu:24.04

WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]
```

---

# 25.3 Kubernetes Deployment

---

# Recommended Namespace Structure

```text id="htohvv"
terminal-platform-prod
terminal-platform-staging
```

---

# Core K8s Components

| Component | Type        |
| --------- | ----------- |
| frontend  | Deployment  |
| backend   | Deployment  |
| websocket | Deployment  |
| postgres  | StatefulSet |
| kafka     | StatefulSet |
| redis     | StatefulSet |

---

# Horizontal Scaling

| Service   | Scaling         |
| --------- | --------------- |
| frontend  | horizontal      |
| backend   | horizontal      |
| websocket | sticky sessions |
| postgres  | primary-replica |
| kafka     | clustered       |

---

# 25.4 Ingress Layer

## Recommended

### NGINX Ingress Controller

[NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/?utm_source=chatgpt.com)

---

# Requirements

* websocket support
* TLS termination
* rate limiting
* gzip
* request buffering control

---

# Example Ingress

```yaml id="cvjlwm"
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: terminal-platform
spec:
  ingressClassName: nginx
```

---

# 25.5 TLS & HTTPS

## Mandatory

* HTTPS only
* WSS only
* HSTS enabled

---

# Recommended Certificate Management

## cert-manager

[cert-manager](https://cert-manager.io/?utm_source=chatgpt.com)

---

# 26. High Availability Strategy

---

# Stateless Services

Can scale horizontally:

* frontend
* backend API
* websocket gateways

---

# Stateful Services

Require persistence:

* PostgreSQL
* Kafka
* Redis

---

# Recommended Persistence

| Service    | Storage               |
| ---------- | --------------------- |
| PostgreSQL | SSD persistent volume |
| Kafka      | SSD persistent volume |
| Redis      | memory + AOF          |

---

# 27. Future Enterprise Extensions

---

# Potential Extensions

| Feature                | Description           |
| ---------------------- | --------------------- |
| SSO                    | enterprise login      |
| LDAP                   | company directory     |
| SIEM integration       | security monitoring   |
| Session replay UI      | playback              |
| AI governance          | prompt moderation     |
| Policy engine          | OPA integration       |
| Multi-node PTY runtime | distributed execution |

---

# Recommended Policy Engine

## Open Policy Agent (OPA)

[Open Policy Agent](https://www.openpolicyagent.org/?utm_source=chatgpt.com)

---

# 28. Final Enterprise Architecture Summary

Hệ thống cuối cùng là:

> A production-grade web-based terminal virtualization platform with centralized governance, Kafka-based event auditing, enterprise authentication, realtime PTY streaming, and shared authenticated runtime management.

Core characteristics:

* logical terminal isolation
* centralized RBAC
* Kafka event governance
* realtime websocket terminal streaming
* auditability
* production-ready deployment
* Kubernetes-compatible architecture
* enterprise observability
* secure external service integration

System focus:

* collaborative infrastructure access
* shared authenticated environments
* browser-native operations
* centralized operational governance
* enterprise audit and monitoring capabilities
