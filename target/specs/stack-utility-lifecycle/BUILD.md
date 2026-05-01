# BUILD: Ancillary Stack Utility Lifecycle (F-008)

## 1. Implementation Scope

- **Slice**: F-008 only (`stack-utility-lifecycle`)
- **Lifecycle decision selected**: **Retain**
- **Files changed**:
  - `target/internal/legacy/stack.go`
  - `target/internal/legacy/stack_test.go`
  - `target/specs/stack-utility-lifecycle/PLAN.md`
  - `target/specs/stack-utility-lifecycle/TASKS.md`
  - `target/specs/stack-utility-lifecycle/BUILD.md`
- **Requirements covered**:
  - FR-001: Exactly one lifecycle status selected (`Retain`).
  - FR-002: Decision rationale tied to source evidence (exported utility + no observed in-repo call sites).
  - FR-003: Retained API behavior implemented in Go (`Stack` + `Add/Test/Done` aggregation semantics).
  - FR-006: External usage confidence explicitly classified as **NEEDS CLARIFICATION**.

## 2. Build Steps

1. Read `SPEC.md`, `PLAN.md`, `TASKS.md`, and shared `ARCHITECTURE` for F-008 boundaries.
2. Activate **Retain** branch for F-008 and keep implementation isolated in `internal/legacy/`.
3. Implement `legacy.Stack` with callback-slot aggregation behavior:
   - `Add` allocates positional slot and returns completion callback.
   - `Test` triggers completion only after all slots finish.
   - `Done` registers final callback and user data.
   - Error output is `nil` when no errors are present.
4. Add focused unit tests for contract parity:
   - positional aggregation independent of completion order,
   - mixed nil/non-nil error slot preservation,
   - immediate completion when no slots are added before `Done`.
5. Run tests and record observed results.
6. Update `PLAN.md` and `TASKS.md` statuses for selected branch execution, including explicit N/A-complete closure of inactive Deprecate/Remove branch tasks for workflow gate compliance.

## 3. Test Evidence

- **Added tests**: `target/internal/legacy/stack_test.go`

- **Command**: `go test ./...` (run in `target/`)
  - **Observed result**: `internal/legacy` passed; overall module failed due to pre-existing non-F-008 issues in other packages (missing Cobra dependency and import cycle under `internal/cli`/`internal/core`/`internal/filter`).

- **Command**: `go test ./internal/legacy` (run in `target/`)
  - **Observed result**: `ok   theseus/target/internal/legacy` (pass)

## 4. Residual Risks

- **External consumer uncertainty remains** (FR-006): in-repo references are absent, but downstream imports are unknown (**NEEDS CLARIFICATION**).
- Full `go test ./...` cannot yet be used as a green gate for this repo because of pre-existing out-of-scope module issues.
- Retained API naming is Go-idiomatic (`Add/Test/Done`) rather than JS export syntax; compatibility is at behavioral-contract level for this Go target slice.

## 5. Traceability

| Requirement / Decision | Implementation / Evidence |
|---|---|
| FR-001 (single lifecycle status) | Selected **Retain** in this BUILD artifact and synchronized in `PLAN.md` / `TASKS.md`. |
| FR-002 (rationale from evidence) | Decision rationale: spec source traceability shows exported Stack utility with no observed in-repo call sites; retain chosen to avoid external break risk. |
| FR-003 (preserve Stack contract) | `target/internal/legacy/stack.go` + tests in `target/internal/legacy/stack_test.go` validating aggregation and callback semantics. |
| FR-006 (external usage confidence classification) | Explicitly recorded as **NEEDS CLARIFICATION** in this BUILD and retained in planning artifacts. |
| SC-002 (single non-contradictory status) | Only Retain path marked active; Deprecate/Remove branch tasks are explicitly marked **N/A-complete** in `TASKS.md` to close conditional branches without contradiction. |
| SC-003 (retain contract verification) | `go test ./internal/legacy` passing evidence captured above. |
