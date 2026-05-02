# Implementation Plan: F-005 Filtering, package restriction, and policy exits

**Date**: 2026-05-01 | **Spec**: `target/specs/filtering-and-policy-enforcement/SPEC.md`  
**Input**: Feature specification from `target/specs/filtering-and-policy-enforcement/SPEC.md`

**Note**: This plan is scoped strictly to the F-005 slice and uses `target/specs/ARCHITECTURE` as the architecture baseline.

## Summary

Reimplement the F-005 filtering and policy-exit stage in Go as `internal/filter` (`filter.go` + `policy.go`), preserving compatibility behaviors from spec and architecture: apply filter transforms/restrictions first, then enforce fail-fast policy gates with exit code `1`. The plan keeps parity-first semantics (including escaped-comma exclude parsing, BSD expansion compatibility, semicolon package-key filtering, private exclusion, and `failOn`/`onlyAllow` exits) and isolates unresolved `onlyAllow` strictness as a blocking clarification.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md` (F-005 only)
- [x] Technical context filled with concrete decisions or `NEEDS CLARIFICATION`
- [x] Work packages derived from the assigned functionality scope
- [x] Dependencies and execution order documented
- [x] Validation approach captured

## Decision Log

- [x] Keep F-005 implementation bounded to filter/policy stage only; no license discovery, flattening, or renderer work (SPEC C-001, C-004).
- [x] Place implementation in architecture-defined paths: `internal/filter/filter.go` and `internal/filter/policy.go`.
- [x] Preserve canonical key handling via `name@version` map operations only (FR-004, C-002).
- [x] Preserve fail-fast compatibility semantics for policy gates with exit status `1` (FR-006, FR-007, C-003).
- [x] Sequence work as: unknown transform -> unknown-only restriction -> exclude license filter -> package include/exclude -> private exclusion -> policy evaluation (FR-008).
- [x] CLI invalid-combination handling (`failOn` + `onlyAllow`) remains required and must be validated at command/options boundary (FR-009; architecture `internal/cli`).
- [x] Context7-aligned conventions adopted for CLI error handling and guardrails: Cobra `RunE` error propagation and mutual-exclusion validation helpers (Context7 `/spf13/cobra`).

## Open Questions

- [ ] **OD-001 (blocking for exact matcher behavior)**: `onlyAllow` matching strictness remains `NEEDS CLARIFICATION` (substring compatibility vs exact/SPDX-aware semantics). Implementation tasks must parameterize matcher strategy behind one compatibility-focused interface/func boundary until decision is made.

## Technical Context

**Language/Version**: Go 1.24.x target (align with architecture Go baseline)  
**Primary Dependencies**: stdlib (`strings`, `regexp`, `context`, `errors`), `github.com/spf13/cobra` (CLI validation path for FR-009), `github.com/stretchr/testify` (tests)  
**Storage**: N/A for F-005 core logic (in-memory map filtering/policy checks)  
**Testing**: `go test` with table-driven unit tests for filter/policy functions + fixture-driven integration tests for exit behavior compatibility  
**Target Platform**: Cross-platform Go CLI execution in CI (Linux/macOS/Windows)  
**Project Type**: Go CLI + reusable core package (architecture-defined)  
**Performance Goals**: Linear pass behavior over filtered package map (no superlinear scans for typical policy lists)  
**Constraints**: Preserve legacy semantics called out in SPEC compatibility notes; fail-fast process-exit compatibility; no scope bleed into non-F-005 slices  
**Scale/Scope**: One pipeline stage only (filter + policy gate) over flattened dependency map keyed by `name@version`

## Constitution Check

*GATE: Must align with applicable constitution or workflow rules before planning proceeds. Re-check after design-oriented planning artifacts are produced.*

- `/memory/constitution.md`: **not present** in repository root; no additional constitution gates discovered.
- `.specify/extensions.yml`: **not present**; no `hooks.before_plan`/`hooks.after_plan` to surface.
- Repository setup/agent-context scripts at repo root: **none discovered** by pattern search (`*setup*.sh`, `*agent*context*`).
- Workflow gate result: **PASS** for planning phase (scope, artifact location, Go target, no implementation).

## Project Structure

### Documentation (this feature)

```text
target/specs/filtering-and-policy-enforcement/
├── PLAN.md              # This file
├── research.md          # Optional; create only if OD-001 needs external comparison research
├── data-model.md        # Optional; create only if model-level clarifications are requested
├── quickstart.md        # Optional; create only if requested for manual validation walkthrough
├── contracts/           # Optional; create only if external interface contracts are requested
└── TASKS.md             # Produced later by Task Decomposer
```

### Source Code (repository root)

```text
cmd/
└── license-checker/

internal/
├── cli/
│   ├── command.go
│   └── options.go
├── core/
│   ├── model.go
│   └── service.go
└── filter/
    ├── filter.go
    └── policy.go

pkg/
└── licensechecker/
    └── api.go

test/
├── integration/
└── fixtures/
```

**Structure Decision**: Use the architecture’s Go service/CLI layout and implement F-005 only within `internal/filter` plus minimal CLI/core touchpoints required for FR-009 and pipeline ordering enforcement.

## Work Packages (F-005 Only)

- [x] **WP1 — Define F-005 option/model contract in core boundary**  
  Traceability: FR-001..FR-009, Inputs/Outputs section, C-002/C-004.  
  Output: explicit option fields and filtered-map contract consumed by `internal/filter` and invoked by `core.Service` in stage order.

- [x] **WP2 — Implement filter transform/restriction sequence (`internal/filter/filter.go`)**  
  Traceability: FR-001, FR-002, FR-003, FR-004, FR-005, FR-008; edge cases (escaped commas, BSD compatibility, private exclusion).  
  Output: deterministic, ordered filter pipeline operating on `name@version` map entries.

- [x] **WP3 — Implement policy gates (`internal/filter/policy.go`) with fail-fast semantics**  
  Traceability: FR-006, FR-007, FR-008, SC-003, C-003.  
  Output: policy evaluation that emits violation context and triggers exit strategy with code `1` at first violation.

- [x] **WP4 — Enforce invalid policy option combination at CLI/options boundary**  
  Traceability: FR-009.  
  Output: guard that rejects concurrent `failOn` + `onlyAllow` and returns CLI failure (`RunE` path) with exit code `1`.

- [x] **WP5 — Compatibility-focused tests for F-005 behavior**  
  Traceability: User Stories 1–3, SC-001..SC-004, compatibility notes.  
  Output: table-driven unit tests + fixture integration tests for ordering, transformations, exclusions, package restrictions, private filtering, and fail-fast exits.

- [x] **WP6 — Resolve OD-001 and lock matcher behavior**  
  Traceability: Open Decision OD-001.  
  Output: explicit matcher seam retained with current compatibility substring behavior and `NEEDS CLARIFICATION` marker preserved for strictness decision before hardening.

## Dependencies and Execution Order

1. WP1 (contract boundary)  
2. WP2 (filter stage)  
3. WP3 (policy stage)  
4. WP4 (CLI guardrail for invalid policy combination)  
5. WP5 (verification matrix + regression fixtures)  
6. WP6 (finalize OD-001 decision if still open; required for completion gate)

Ordering rationale: architecture mandates pipeline stage order and spec requires filters before policy exits (FR-008). CLI guardrail can be implemented in parallel with WP2/WP3 but must be validated in the same F-005 test matrix.

## Validation Approach

- [x] **V1: Unknown transform tests** — verify 100% rewrite of `*`-suffixed licenses to `UNKNOWN` when `unknown=true` (SC-001).
- [x] **V2: onlyunknown restriction tests** — verify only `UNKNOWN`/guessed entries remain (SC-002).
- [x] **V3: Exclude matching tests** — include escaped-comma and BSD compatibility fixtures (FR-003, edge cases).
- [x] **V4: Package include/exclude tests** — semicolon list parsing + canonical `name@version` key matching (FR-004).
- [x] **V5: Private exclusion tests** — verify removal regardless of upstream `UNLICENSED` labeling (FR-005, edge cases).
- [x] **V6: Policy fail-fast tests** — `failOn` and `onlyAllow` emit violation and terminate with code `1` (FR-006/FR-007, SC-003).
- [x] **V7: Order-of-operations tests** — prove filters run before policy checks (FR-008).
- [x] **V8: CLI guardrail tests** — both `failOn` + `onlyAllow` returns invalid-combination failure (FR-009).

## Complexity Tracking

No constitution violations identified; this section is intentionally empty.
