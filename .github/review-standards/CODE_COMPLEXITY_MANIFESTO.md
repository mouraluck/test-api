# Universal Code Complexity Manifesto

A lightweight, universal scoring model for PR complexity.

This standard extends the repository review patterns and adds a numeric, repeatable complexity score to support planning, review depth, and release risk decisions.

## 1. Goals

- Keep complexity scoring objective and explainable.
- Allow teams to estimate review/deploy risk quickly.
- Avoid subjective labels without evidence.

## 2. Scoring Method

Complexity score = `Base Change Points + Scope/Risk Modifiers - Mitigation Credits`

Rules:
- Score each relevant change category once per PR (avoid duplicate counting).
- Use the highest applicable point when overlap exists.
- Minimum final score is `0`.

## 3. Base Change Points (Universal)

- `+5` Database schema/data migration/backfill changes.
- `+4` Base classes, core abstractions, shared interfaces, or framework foundation changes.
- `+4` External contract changes (public API, event contracts, serialization formats).
- `+4` Security/permission/authentication/authorization logic changes.
- `+3` Core business rule changes affecting decision logic.
- `+3` Integration changes with external services/providers.
- `+2` Deployment/runtime/configuration behavior changes.
- `+1` New isolated class/module with no existing contract impact.
- `+1` Internal refactor with no behavior change.
- `+0.5` Tests-only changes.
- `+0.25` Documentation-only changes.

Examples aligned with requested baseline:
- "Mexe no banco de dados" -> `+5`
- "Altera classes base" -> `+4`
- "Classe nova" -> `+1`

## 4. Scope and Risk Modifiers

Apply modifiers once per PR after base points:

- `+1` Multi-layer impact (e.g., domain + transport + persistence).
- `+1` Cross-service or cross-repository dependency impact.
- `+1` Irreversible or difficult rollback path.
- `+1` No backward compatibility plan for contract/schema changes.
- `+1` Hot path/performance-sensitive code touched.
- `+1` Concurrency/async/transactional behavior changed.

## 5. Mitigation Credits

Subtract credits only when clearly implemented in the PR:

- `-1` Strong automated test coverage for changed behavior (success + failure paths).
- `-1` Feature flag/canary strategy included.
- `-1` Backward compatibility preserved (or explicit migration path provided).

Max mitigation subtraction: `-3`.

## 6. Complexity Bands and Required Review Rigor

- `0-3` Low:
  - Standard review.
  - Basic test validation.
- `4-7` Medium:
  - Focused risk review.
  - Explicit test gap check.
- `8-12` High:
  - Deep review with rollback plan validation.
  - Reviewer should challenge assumptions and edge cases.
- `13+` Critical:
  - Mandatory staged rollout/rollback notes.
  - Prefer additional reviewer and stronger validation evidence.

## 7. PR Output Format (Required)

When reviewing a PR, report complexity explicitly:

1. `Complexity Score: <number>`
2. `Band: Low|Medium|High|Critical`
3. `Breakdown:` list each scored item and points.
4. `Mitigations:` list applied credits.
5. `Residual Risk:` short objective statement.

## 8. Guardrails

- Do not inflate score by counting trivial edits repeatedly.
- Do not reduce score without concrete mitigation in diff/tests.
- If uncertain, choose the safer (higher) score and explain why.

## 9. Adoption Notes

This manifesto is intentionally simple and language-agnostic.

Teams can calibrate thresholds later, but base point semantics should remain stable to keep historical comparability across PRs.
