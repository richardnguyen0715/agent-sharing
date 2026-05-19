# High-Standard Software Engineering Policies

This is a practical policy framework for writing code that is:

* Easy to read
* Easy to scale
* Easy to maintain
* Easy to test
* Easy to onboard new developers into
* Resistant to technical debt

These policies are heavily influenced by practices from companies like [Google Engineering Practices](https://google.github.io/eng-practices/?utm_source=chatgpt.com), [Martin Fowler](https://martinfowler.com?utm_source=chatgpt.com), [Refactoring Guru](https://refactoring.guru?utm_source=chatgpt.com), [Clean Code concepts](https://www.oreilly.com/library/view/clean-code-a/9780136083238/?utm_source=chatgpt.com), and large-scale engineering organizations.

---

# 1. Core Engineering Principles

## 1.1 Readability First

Policy:

* Code is read far more often than written.
* Optimize for future maintainers, not short-term typing speed.
* Clever code is banned unless necessary and documented.

Rules:

* Prefer explicitness over magic.
* Prefer clarity over brevity.
* Minimize cognitive load.
* One responsibility per block of logic.

Bad:

```ts
const r = a.map(x=>x.v).filter(Boolean)
```

Good:

```ts
const validValues = accounts
  .map(account => account.value)
  .filter(value => value !== null)
```

---

## 1.2 Consistency Over Preference

Policy:

* Team consistency matters more than individual style.

Rules:

* Follow formatter rules automatically.
* Never mix patterns for identical problems.
* One naming convention per codebase.
* One architecture style per service.

Recommended:

* ESLint
* Prettier
* Black
* Ruff
* EditorConfig

---

## 1.3 Simplicity Wins

Policy:

* Prefer the simplest correct solution.

Rules:

* Avoid premature optimization.
* Avoid unnecessary abstraction.
* Avoid over-engineering.
* Delete complexity aggressively.

Use:

* KISS Principle
* YAGNI Principle
* Occam’s Razor

---

# 2. Naming Policies

## 2.1 Naming Must Reveal Intent

Policy:

* Names should explain WHY the thing exists.

Rules:

* Avoid vague names.
* Avoid abbreviations unless industry-standard.
* Avoid meaningless suffixes/prefixes.

Forbidden:

```ts
data
info
manager
helper
temp
misc
util
handleStuff
```

Preferred:

```ts
userProfile
invoiceCalculator
paymentRetryPolicy
emailVerificationToken
```

---

## 2.2 Variable Naming Standards

### Variables

Use nouns.

```ts
const customer
const invoice
const retryCount
```

### Booleans

Must read like questions.

```ts
isActive
hasPermission
canRetry
shouldRefresh
```

### Functions

Use verbs.

```ts
calculateTotal()
fetchUsers()
validatePassword()
sendEmail()
```

### Async Functions

Must indicate async behavior.

```ts
fetchOrders()
loadConfiguration()
syncDatabase()
```

Avoid:

```ts
doStuff()
process()
handle()
```

unless context is extremely precise.

---

## 2.3 Length Rules

Policy:

* Long enough to be clear.
* Short enough to scan.

Rules:

* Scope determines length.
* Smaller scope → shorter names acceptable.

Example:

```ts
for (const i of items)
```

acceptable.

But:

```ts
const d = calculate()
```

not acceptable.

---

# 3. Function Policies

## 3.1 Single Responsibility

Policy:

* A function should do ONE thing.

Bad:

```ts
function createUser() {
  validateInput()
  saveUser()
  sendEmail()
  logAnalytics()
}
```

Good:

```ts
function createUser()
function sendWelcomeEmail()
function trackSignupAnalytics()
```

Or orchestrate cleanly:

```ts
async function registerUser() {
  const user = await createUser()
  await sendWelcomeEmail(user)
  trackSignupAnalytics(user)
}
```

---

## 3.2 Function Size Limits

Recommended:

* Ideal: 5–20 lines
* Warning: >30 lines
* Refactor required: >50 lines

Exceptions:

* Parsers
* State machines
* Algorithms

---

## 3.3 Parameter Policies

Rules:

* Prefer ≤3 parameters.
* Use objects for complex input.
* Avoid boolean flags.

Bad:

```ts
createUser(name, email, true, false)
```

Good:

```ts
createUser({
  name,
  email,
  sendWelcomeEmail: true
})
```

---

## 3.4 Pure Functions Preferred

Policy:

* Functions should minimize side effects.

Preferred:

```ts
function calculateTax(price) {
  return price * TAX_RATE
}
```

Avoid hidden mutation:

```ts
function applyTax(cart) {
  cart.total += ...
}
```

---

## 3.5 Error Handling

Rules:

* Fail fast.
* Never silently swallow errors.
* Use typed/domain-specific errors.

Bad:

```ts
catch (e) {}
```

Good:

```ts
catch (error) {
  logger.error(error)
  throw new PaymentProcessingError()
}
```

---

# 4. OOP Policies

## 4.1 Composition Over Inheritance

Policy:

* Prefer composing behaviors instead of deep inheritance trees.

Bad:

```txt
Animal
 └── Mammal
      └── Dog
           └── PoliceDog
```

Preferred:

```ts
class Dog {
  constructor(
    private readonly movement: MovementStrategy,
    private readonly communication: CommunicationStrategy
  ) {}
}
```

---

## 4.2 SOLID Principles

## S — Single Responsibility

One reason to change.

## O — Open/Closed

Extend without modifying.

## L — Liskov Substitution

Subtypes must behave correctly.

## I — Interface Segregation

Small focused interfaces.

## D — Dependency Inversion

Depend on abstractions.

---

## 4.3 Encapsulation

Rules:

* Hide internal implementation.
* Expose minimal public API.
* Avoid public mutable state.

Bad:

```ts
user.balance = -1000
```

Good:

```ts
user.withdraw(amount)
```

---

## 4.4 Immutable Data Preferred

Policy:

* Immutable state reduces bugs.

Preferred:

```ts
return {
  ...user,
  status: 'active'
}
```

Avoid:

```ts
user.status = 'active'
```

---

# 5. Design Pattern Policies

## 5.1 Use Patterns Only When Necessary

Policy:

* Don’t force patterns.

Anti-pattern:

```txt
AbstractFactorySingletonBuilderManager
```

---

## 5.2 Recommended Core Patterns

## Strategy Pattern

For interchangeable behavior.

Use when:

* Multiple algorithms
* Payment systems
* Validation rules

---

## Factory Pattern

For object creation complexity.

Use when:

* Complex initialization
* Dependency injection
* Environment-specific objects

---

## Repository Pattern

For data access abstraction.

Use when:

* Database isolation
* Testing
* Multiple storage providers

---

## Dependency Injection

Mandatory for scalable systems.

Benefits:

* Testability
* Loose coupling
* Modularity

---

## Observer/Event Pattern

Use for:

* Event-driven systems
* Analytics
* Messaging
* Notifications

---

## Adapter Pattern

Use for external APIs.

Purpose:

* Prevent vendor lock-in leakage.

---

# 6. Architecture Policies

## 6.1 Layered Architecture

Recommended structure:

```txt
controllers/
services/
repositories/
domain/
infrastructure/
shared/
```

---

## 6.2 Domain-Driven Design (DDD)

Use for large systems.

Core concepts:

* Entities
* Value Objects
* Aggregates
* Bounded Contexts

Avoid anemic domain models.

---

## 6.3 Clean Architecture

Dependency flow:

```txt
Infrastructure → Application → Domain
```

Domain must not depend on frameworks.

---

## 6.4 Modular Monolith First

Policy:

* Prefer modular monolith before microservices.

Why:

* Easier deployment
* Easier debugging
* Lower complexity
* Better iteration speed

Microservices only when:

* Independent scaling needed
* Team scaling requires it
* Operational maturity exists

---

# 7. DSA (Data Structures & Algorithms) Policies

## 7.1 Choose Correct Complexity

Know:

* Time complexity
* Space complexity
* Tradeoffs

Mandatory knowledge:

| Structure | Average Lookup |
| --------- | -------------- |
| Array     | O(n)           |
| HashMap   | O(1)           |
| BST       | O(log n)       |
| Heap      | O(log n)       |

---

## 7.2 Prefer Readability Over Micro-Optimization

Bad:

```ts
i^=j^=i^=j
```

Good:

```ts
const temp = a
a = b
b = temp
```

---

## 7.3 Avoid Reinventing Core Structures

Policy:

* Use battle-tested libraries first.

Unless:

* Interview
* Educational purpose
* Performance-critical infrastructure

---

## 7.4 Algorithmic Standards

Rules:

* Avoid nested loops without justification.
* Benchmark before optimization.
* Document non-obvious complexity.

---

# 8. Testing Policies

## 8.1 Mandatory Automated Testing

Required:

* Unit tests
* Integration tests
* End-to-end tests for critical flows

---

## 8.2 Test Behavior, Not Implementation

Bad:

```ts
expect(mockFn).toBeCalledTimes(3)
```

Good:

```ts
expect(order.status).toEqual('completed')
```

---

## 8.3 Coverage Goals

Recommended:

| Type                  | Coverage |
| --------------------- | -------- |
| Critical Domain Logic | 90%+     |
| Services              | 80%+     |
| UI                    | 60–80%   |

Coverage is not quality.

---

# 9. API Policies

## 9.1 REST Standards

Rules:

* Use nouns.
* Version APIs.
* Consistent error shape.

Good:

```txt
GET /users
POST /orders
```

Avoid:

```txt
/getUsersData
/doOrder
```

---

## 9.2 Validation

Mandatory:

* Validate all external input.
* Never trust client data.

Use:

* Schema validation
* DTOs
* Zod/Joi/Pydantic

---

## 9.3 Backward Compatibility

Rules:

* Never break clients casually.
* Deprecate gradually.
* Version carefully.

---

# 10. Database Policies

## 10.1 Schema Standards

Rules:

* Explicit constraints
* Foreign keys
* Proper indexing
* Avoid nullable chaos

---

## 10.2 Transactions

Use for:

* Financial operations
* Multi-step consistency

---

## 10.3 Query Policies

Rules:

* Avoid N+1 queries.
* Use pagination.
* Measure query performance.

---

# 11. Security Policies

Mandatory:

* Input sanitization
* Authentication
* Authorization
* Rate limiting
* Secrets management
* Dependency scanning

Never:

* Store plaintext passwords
* Commit secrets
* Trust frontend authorization

Use:

* OWASP Top 10

Reference:
[OWASP Foundation](https://owasp.org?utm_source=chatgpt.com)

---

# 12. Logging & Observability

## Logging Rules

Must include:

* Context
* Correlation IDs
* Error traces

Avoid:

* Logging secrets
* Logging sensitive user data

---

## Monitoring

Mandatory:

* Metrics
* Alerts
* Distributed tracing

Recommended stack:

* Prometheus
* Grafana
* OpenTelemetry

---

# 13. Git Policies

## Branching

Preferred:

```txt
main
develop
feature/*
fix/*
```

---

## Commit Standards

Use:

```txt
feat:
fix:
refactor:
test:
docs:
```

Example:

```txt
feat(auth): add refresh token rotation
```

---

## PR Rules

Mandatory:

* Small PRs
* Self-review first
* CI passing
* Reviewer approval

---

# 14. Documentation Policies

Mandatory documentation:

* Architecture decisions
* Public APIs
* Setup instructions
* Complex business rules

Avoid:

* Commenting obvious code

Bad:

```ts
// increment i
i++
```

Good:

```ts
// Retry required because provider has eventual consistency delay
```

---

# 15. Performance Policies

Rules:

* Measure before optimizing.
* Optimize bottlenecks only.
* Avoid premature caching.

Performance hierarchy:

```txt
Architecture > Database > Network > Algorithm > Micro-optimization
```

---

# 16. Code Review Standards

A reviewer checks:

* Readability
* Correctness
* Simplicity
* Security
* Scalability
* Testability
* Naming
* Architecture consistency

Questions reviewers should ask:

* Is this understandable in 30 seconds?
* Is complexity justified?
* Can this break silently?
* Is this easy to modify later?

---

# 17. Anti-Patterns to Ban

Ban:

* God classes
* Massive functions
* Hidden side effects
* Circular dependencies
* Deep nesting
* Magic numbers
* Copy-paste logic
* Premature abstraction
* Global mutable state
* Overuse of inheritance

---

# 18. Golden Rules

## The “3-Year Rule”

Write code assuming:

> Someone unfamiliar will maintain this in 3 years.

---

## The “2 AM Rule”

Code should be understandable:

> By a tired engineer at 2 AM during production outage.

---

## The “Delete Resistance Rule”

Good architecture allows deletion safely.

If deleting one feature breaks everything:
architecture is too coupled.

---

# 19. Recommended Engineering Standards

## General

* [Google Engineering Practices](https://google.github.io/eng-practices/?utm_source=chatgpt.com)
* [Microsoft REST API Guidelines](https://github.com/microsoft/api-guidelines?utm_source=chatgpt.com)
* [12 Factor App](https://12factor.net?utm_source=chatgpt.com)

## Architecture

* [Refactoring Guru Design Patterns](https://refactoring.guru/design-patterns?utm_source=chatgpt.com)
* [Martin Fowler Architecture Guides](https://martinfowler.com/architecture/?utm_source=chatgpt.com)

## Security

* [OWASP Top 10](https://owasp.org/www-project-top-ten/?utm_source=chatgpt.com)

---

# 20. Elite-Level Engineering Habits

High-level engineers consistently:

* Reduce complexity
* Improve naming
* Remove duplication
* Isolate risk
* Write defensive code
* Build observability early
* Design for change
* Think in systems
* Optimize developer experience
* Prefer boring reliable solutions

---

# Recommended “Ideal Stack” for Maintainability

## TypeScript Backend

* TypeScript
* ESLint
* Prettier
* Zod
* Prisma/Drizzle
* PostgreSQL
* Vitest/Jest
* OpenTelemetry

## Python Backend

* Python
* Ruff
* Black
* Mypy
* FastAPI
* Pydantic
* Pytest

---

# Final Principle

The highest standard of coding is not:

* The smartest code
* The shortest code
* The most abstract code

It is:

> Code that remains understandable, reliable, and adaptable as the system and team grow.
