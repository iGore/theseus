# Implementation Plan: Ancillary Stack Utility Lifecycle (F-008)

**Date**: 2026-05-01 | **Spec**: `target/specs/stack-utility-lifecycle/SPEC.md`  
**Input**: Feature specification from `target/specs/stack-utility-lifecycle/SPEC.md`

**Note**: This plan is scoped only to F-008 and uses `target/specs/ARCHITECTURE` as the architecture baseline.

## Summary

Decide and document exactly one lifecycle path for the legacy `Stack` utility (Retain, Deprecate, or Remove), then plan the Go-side implementation and validation steps needed for that choice. The architecture already places this slice in `internal/legacy/stack.go` and schedules it as implementation-order step 9, so this plan focuses on: (1) decision finalization with explicit evidence, (2) compatibility handling per selected path, and (3) traceable verification artifacts.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md` (F-008 only)
- [x] Technical context filled with concrete decisions or `NEEDS CLARIFICATION`
- [x] Work packages derived from the assigned functionality scope
- [x] Dependencies and execution order documented
- [x] Validation approach captured

## Decision Log

- [x] Use-case scope fixed to F-008 only; no cross-slice behavior redesign allowed (Spec Assumptions, FR-001..FR-006).
- [x] Go target retained; no non-Go output is planned (repo default + user constraint).
- [x] `Stack` lifecycle work remains isolated under `internal/legacy/` per shared architecture to minimize blast radius.
- [x] Planning keeps lifecycle decision as an explicit gate before Builder implementation details.
- [x] Context7 guidance incorporated for Go/Cobra conventions used by planning assumptions:
  - Go deprecation marker convention: `Deprecated:` comment text is the compatibility communication baseline when Deprecate path is chosen.
  - Cobra error propagation/validation sequencing via `RunE` + command-level guardrails remains consistent with architecture expectations for CLI-facing notices.

## Open Questions

- [x] Final lifecycle status selection resolved for this build slice: **Retain** (recorded in `BUILD.md`).
- [ ] **NEEDS CLARIFICATION**: External consumer confidence remains unknown (Spec FR-006); no telemetry/downstream import inventory is available.
- [ ] If `Deprecate` or `Remove` is selected, confirm the exact release window and release-note channel used by the project.

## Technical Context

**Language/Version**: Go (target version aligned with repository/toolchain; exact version **NEEDS CLARIFICATION**)  
**Primary Dependencies**: Go stdlib; Cobra at CLI boundary for validation/error propagation consistency; testify for tests (from shared architecture)  
**Storage**: N/A for this slice beyond repository files and release artifacts  
**Testing**: `go test` with table-driven unit tests for Stack contract behavior + integration checks for selected lifecycle path  
**Target Platform**: Cross-platform Go CLI context (Linux/macOS/Windows) as defined in shared architecture  
**Project Type**: Go CLI + reusable package architecture; F-008 implemented as internal legacy utility slice  
**Performance Goals**: N/A (utility lifecycle/compatibility decision slice)  
**Constraints**: Must preserve compatibility semantics if retained/deprecated; must provide explicit migration/breaking-change guidance if removed; no scope expansion beyond F-008  
**Scale/Scope**: Single ancillary utility (`Stack`) and associated release/compatibility handling

## Constitution Check

*GATE: Must align with applicable constitution or workflow rules before planning proceeds. Re-check after design-oriented planning artifacts are produced.*

- `/memory/constitution.md`: not present in workspace; no additional constitution gates discovered.
- `.specify/extensions.yml`: not present; no before/after plan hooks discovered.
- Repository workflow rules (`AGENTS.md`) satisfied for this artifact:
  - Spec-first flow preserved (planning from SPEC + ARCHITECTURE only).
  - Artifact location preserved under `target/specs/stack-utility-lifecycle/`.
  - Go remains default target language.
  - Plan remains checklist-oriented and bounded to one use-case slice.

## Work Packages (F-008 only)

### WP1 — Finalize lifecycle decision record (FR-001, FR-002, FR-006)

- [x] Choose exactly one status: `Retain` / `Deprecate` / `Remove`.
- [x] Record rationale with source links (`SRC-F008-01..05`) and explicit statement of no in-repo call sites.
- [x] Record external usage confidence classification as required by FR-006.
- [x] Add approval note so downstream Task Decomposer/Builder can execute without ambiguity.

**Deliverable**: Approved Lifecycle Decision Record entry for F-008.

### WP2 — Path-specific implementation planning contract

- [x] **If Retain**: define contract-preservation checklist for `Stack`, `add`, `test`, `done`, callback aggregation ordering (FR-003, SC-003).
- [x] **If Deprecate**: define one-release compatibility shim + deprecation guidance requirements (FR-004, SC-004). *(N/A: Retain path selected)*
- [x] **If Remove**: define breaking-change notice + migration guidance requirements (FR-005, SC-004). *(N/A: Retain path selected)*

**Deliverable**: One selected path with explicit must-pass criteria.

### WP3 — Verification and evidence plan

- [x] Define unit/integration test evidence required for selected path (contract behavior for retain/deprecate; release artifact checks for deprecate/remove).
- [x] Define traceability matrix from FRs to planned validations (SC-001, SC-002).
- [x] Define artifact checklist for handoff to Task Decomposer (`TASKS.md` inputs) with no unresolved blockers except declared clarifications.

**Deliverable**: Validation checklist tied to FR/SC IDs.

## Dependencies & Execution Order

1. **WP1 must complete first** (decision gate).  
2. **WP2 depends on WP1** because requirements diverge by chosen status path.  
3. **WP3 depends on WP2** to validate the selected path only.  

No F-008 execution work should begin before WP1 approval.

## Validation Approach

- [x] Confirm all functional requirements (FR-001..FR-006) are explicitly covered by planned tasks and acceptance checks.
- [x] Confirm a single non-contradictory lifecycle status is selected across all planning artifacts (SC-002).
- [x] For `Retain`: verify API/behavior parity checks exist for `Stack` + `add/test/done` aggregation semantics (SC-003).
- [x] For `Deprecate`/`Remove`: verify explicit release-note migration impact text is included pre-shipment (SC-004). *(N/A: Retain path selected)*
- [x] Confirm unresolved items remain explicitly labeled `NEEDS CLARIFICATION` (no implicit assumptions).

## Project Structure

### Documentation (this feature)

```text
target/specs/stack-utility-lifecycle/
├── SPEC.md
├── PLAN.md
└── TASKS.md             # Produced later by Task Decomposer
```

### Source Code (repository root)

```text
cmd/license-checker/
internal/legacy/
pkg/licensechecker/
```

**Structure Decision**: Keep F-008 implementation isolated in `internal/legacy/` (architecture section 3), with any CLI-facing messaging flowing through existing Cobra-based command handling.

## Complexity Tracking

No constitution violations identified; table not required.
