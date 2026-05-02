# BUILD.md — F-005 Filtering, package restriction, and policy exits

## 1. Implementation Scope

- **Files changed**
  - `target/internal/core/model.go`
  - `target/internal/core/service.go`
  - `target/internal/filter/filter.go`
  - `target/internal/filter/policy.go`
  - `target/internal/cli/options.go`
  - `target/internal/filter/filter_test.go`
  - `target/internal/filter/policy_test.go`
  - `target/test/integration/filter_policy_integration_test.go`
  - `target/specs/filtering-and-policy-enforcement/PLAN.md`
  - `target/specs/filtering-and-policy-enforcement/TASKS.md`

- **Spec requirements covered (F-005 only)**
  - FR-001, FR-002: unknown transform + onlyunknown restriction.
  - FR-003: exclude filtering, escaped-comma parsing, BSD compatibility matching.
  - FR-004: semicolon-delimited package include/exclude by canonical `name@version` key.
  - FR-005: private-package exclusion.
  - FR-006, FR-007: fail-fast `failOn`/`onlyAllow` policy errors.
  - FR-008: filter stage executes before policy evaluation.
  - FR-009: `failOn` + `onlyAllow` mutual exclusion validation.

## 2. Build Steps

1. Implemented F-005 option/model contract fields in `internal/core/model.go`.
2. Implemented ordered filter pipeline in `internal/filter/filter.go`:
   - unknown transform → onlyunknown restriction → exclude matching → package include/exclude → private exclusion.
3. Implemented policy evaluation and guardrails in `internal/filter/policy.go`:
   - fail-fast policy errors,
   - `onlyAllow` matcher seam retained for OD-001 compatibility handling,
   - mutual-exclusion validator.
4. Wired service orchestration and exit hook in `internal/core/service.go` so policy violations trigger exit strategy code `1`.
5. Wired CLI-side validation in `internal/cli/options.go` (FR-009 boundary).
6. Added unit/integration tests for F-005 behavior and ordering.
7. Updated `PLAN.md` and `TASKS.md` checklists to reflect completion and explicit OD-001 defer note.

## 3. Test Evidence

- **Tests added/updated**
  - `target/internal/filter/filter_test.go`
  - `target/internal/filter/policy_test.go`
  - `target/test/integration/filter_policy_integration_test.go`

- **Command run**
  - `gofmt -w ... && go test ./...` (from `target/`)

- **Observed result**
  - `ok   theseus/target/internal/cli`
  - `ok   theseus/target/internal/core`
  - `ok   theseus/target/internal/filter`
  - `ok   theseus/target/internal/flatten`
  - `ok   theseus/target/internal/legacy`
  - `ok   theseus/target/test/integration`
  - `?    theseus/target/internal/render [no test files]`

## 4. Residual Risks

- **OD-001 remains deferred by design**: `onlyAllow` strictness (substring compatibility vs exact/SPDX-aware) is still a known open decision; implementation keeps a matcher seam and marks it `NEEDS CLARIFICATION` in code.
- Current policy matching is compatibility-first (substring/BSD handling), not strict SPDX semantics.

## 5. Traceability

- **FR-001 / SC-001** → `internal/filter/filter.go` (unknown transform), `internal/filter/filter_test.go` (`unknown rewrites guessed marker`).
- **FR-002 / SC-002** → `internal/filter/filter.go` (onlyunknown stage), `internal/filter/filter_test.go` (`onlyunknown retains only unknown or guessed`).
- **FR-003** → `internal/filter/filter.go` (`parseCommaListWithEscapes`, BSD matching branch), `internal/filter/filter_test.go` (`exclude supports BSD and escaped commas`).
- **FR-004 / SC-004** → `internal/filter/filter.go` (`parseSemicolonList`, include/exclude package restrictions), `internal/filter/filter_test.go` package restriction tests.
- **FR-005** → `internal/filter/filter.go` private exclusion stage + test.
- **FR-006 / FR-007 / SC-003** → `internal/filter/policy.go` fail-fast policy errors; `internal/filter/policy_test.go`; exit hook verified in `test/integration/filter_policy_integration_test.go`.
- **FR-008** → enforced order in `internal/filter/filter.go` and service orchestration in `internal/core/service.go`; verified by integration test (`unknown` then `failOn UNKNOWN`).
- **FR-009** → `internal/filter/policy.go` (`ValidatePolicyCombination`) + `internal/cli/options.go` (`ValidateOptions`) + policy test coverage.
