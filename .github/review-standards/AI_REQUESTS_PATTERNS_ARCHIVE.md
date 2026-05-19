# AI Request Patterns Archive

Generic guidance for future AI models working on software tasks.

## 1. How to Execute Requests

- Prefer doing the work over only proposing it.
- Keep changes incremental, reviewable, and reversible.
- Confirm behavior with tests/checks whenever possible.
- If a request is ambiguous, choose the safest reasonable assumption and document it.

## 2. Naming and Structure Principles

- Use explicit, intention-revealing names.
- Keep naming consistent across related layers (domain, API, UI, tests).
- Prefer one clear naming convention over mixed styles.
- Avoid broad renames unless requested or clearly beneficial.

## 3. Refactor Expectations

When refactoring:

- Preserve behavior unless change is explicitly required.
- Update all call sites in one pass.
- Remove stale references.
- Keep scope tight to requested outcome.

## 4. API/Data Contract Discipline

- Treat data shapes as contracts.
- If structure changes, update all consumers (backend, frontend, tests, mocks).
- Add/adjust tests that validate the new contract shape.
- Prefer backward-compatible migrations unless breaking change is intentional.

## 5. Testing Standards

- Use the project testing conventions, not the language default by assumption.
- Prefer the suite framework adopted by the codebase (for Go in this repo: `testify/suite`).
- Add tests for both success and failure paths.
- Cover edge cases for parsing, validation, and null/empty values.
- Validate ordering/sorting when it is part of behavior.
- Use table-driven cases when many scenarios share structure.
- Keep scenario-specific higher-level tests explicit and readable.

## 6. Test Organization

- Follow the repository’s preferred placement for tests.
- If requested, keep tests centralized in dedicated test suites/directories.
- Avoid leaving one-off tests outside the agreed test structure.

## 7. Review Standards

During code review:

- Prioritize correctness, regressions, and behavioral risk.
- Flag naming inconsistencies that reduce clarity.
- Highlight missing test coverage for changed behavior.
- Report findings with precise file references.

## 8. Validation and Verification

- Run focused checks after local edits.
- Run full suites when environment allows.
- If full validation cannot run, state exactly why and what remains unverified.

## 9. Error Handling and Observability

- Do not silently swallow recoverable failures when observability is useful.
- Add structured logging for skipped/ignored invalid inputs when they can hide data quality issues.
- Include meaningful context fields in logs (IDs, operation, error).

## 10. Cross-Layer Consistency

When one layer changes, verify all impacted layers:

- data models
- business logic
- transport/serialization
- UI adapters/rendering
- test fixtures/mocks

## 11. Communication Pattern

- Be concise and objective.
- Summarize what changed, why, and how it was verified.
- Call out remaining risks or follow-up steps.
- Avoid unnecessary verbosity when request is simple.

## 12. Completion Checklist

Before considering a task done:

1. Requested scope is fully implemented.
2. Naming is consistent.
3. Contracts are aligned across consumers.
4. Tests are updated/added as needed.
5. Relevant checks were executed.
6. Any unverified areas are explicitly documented.
7. No unrelated changes were introduced.

## 13. Backward-Compatible Worker Evolution

When introducing new jobs in an existing service:

- Preserve currently running jobs unless removal is explicitly requested.
- Prefer additive rollout (new jobs + feature flags/env toggles) before replacement.
- Keep existing behavior active by default when the team depends on it in day-to-day usage.
- Avoid "architecture cleanup" changes that disable known-good flows without explicit approval.

## 14. Naming Rules for Runtime Job Code

- Do not rely on folder context to clarify generic filenames like `types.go` or `config.go`.
- Prefer explicit and intention-revealing file names for runtime job configuration and type contracts.
- Use type names that encode domain intent (e.g., `offerScheduledJobRuntimeConfig`) rather than ambiguous names (`singleJobConfig`).
- Keep definition/spec files explicit when behavior is data-driven (e.g., status transition definitions).

## 15. Test Placement and Framework Consistency

- Keep new tests under the repository's `tests/` structure unless the project explicitly uses co-located tests.
- Follow the repository's chosen test style (`testify/suite` in this codebase) for integration-style suites.
- Use `require` for preconditions and fatal assertions; use `assert`/suite assertions for follow-up checks.
- Do not introduce new testing style variants unless requested.

## 16. Job Runtime Config Guardrails

- Config parsing should fail fast for invalid env values (invalid bool/duration).
- Duration-based schedules must reject non-positive values.
- Defaults should match delivery intent and rollout plan (e.g., enable only what is approved for the current issue).
- Add focused tests for each config key's invalid-path behavior.

## 17. Mutation-Testing Mindset (Without Tooling)

When full mutation tooling is unavailable, simulate mutant resistance by adding tests for:

- Invalid enum/transition inputs that should fail fast.
- Idempotency (second run should not reprocess already transitioned rows).
- Branches around hooks/callbacks (both success and failure paths).
- Scheduler-driven execution (state transition asserted via periodic run, not only direct method calls).
- Config invalid-paths and disabled/enabled behavior switches.

## 18. PR Feedback: Log Before Returning Infra Errors

- When a service method returns an infrastructure/repository error, prefer emitting a structured `slog.Error` first.
- Include operational context fields (job/task name, transition/status fields, and error).
- Do this especially in worker/cron flows where the caller may only see aggregate failure and not domain context.

## 19. Shared Type Placement (Service vs Domain/Data)

- Do not keep shared DTO/value types in `service` when they are used by jobs, tests, and multiple layers.
- Prefer placing cross-layer transition contracts in the domain model layer (or a data/contract layer when purely transport-oriented).
- Keep `service` focused on orchestration/behavior, not ownership of reusable contracts.

## 20. Secrets Management Procedure

- Define each secret key as a typed constant in a centralized secrets-definition layer.
- Register new keys in the secrets loading pipeline:
- Add to the required-at-startup set when needed in all environments.
- Add to an environment-specific required set when only needed in selected environments.
- Provision secrets in GCP Secret Manager using environment-aware naming:
- Preferred: `<environment>-<secret-name>` (for example `development-dojizap-zapi-client-token`).
- Fallback supported by code: `<secret-name>`.
- Load and validate secrets through typed integration/module configs (fail fast when missing).
- Never log secret values; log only metadata/context (component, operation, missing key name).
- Keep tests aligned with required secrets:
- Update mock/default secret fixtures when adding new required secrets.
- Add focused tests for missing-secret validation paths.

## 21. Integration Client File Boundaries

- Keep transport implementation files focused on concrete provider logic (HTTP request/response, headers, retries, decoding).
- Move shared abstractions (interfaces) to dedicated files when they are consumed by services/mocks/tests across packages.
- Prefer explicit names for abstraction files to make ownership and intent obvious.
- Keep request/response DTOs with the concrete integration when they are endpoint-specific, and avoid mixing orchestration concerns into client files.

## 22. Typed JSON Responses in Tests

- In HTTP handler/client tests, prefer encoding response payloads from typed structs instead of writing raw JSON strings.
- Use the same DTO/response types expected by the production decoder whenever possible.
- This reduces drift between tests and contracts, catches field/tag changes earlier, and avoids malformed literal JSON in test fixtures.

## 23. Avoid Unnecessary Input Wrappers

- For service methods that only need a single primitive input (for example, one message string), prefer passing the primitive directly.
- Avoid creating one-field input structs unless there is a near-term need for extensibility, validation grouping, or interface consistency.
- This keeps method signatures simpler and reduces boilerplate in callers and tests.

---

This archive is intentionally generic and should apply to most coding tasks regardless of language, architecture, or domain.
