# Implementation Tasks: Dependency Graph Loading and Traversal Control (F-002)

**Date**: 2026-05-01 | **Plan**: `target/specs/dependency-graph-loading/PLAN.md` | **Spec**: `target/specs/dependency-graph-loading/SPEC.md`
**Scope**: Reduced — Strictly F-002 graph-loading behavior only (start path, depth mapping, dev toggle, success/error handoff)

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 17    |
| Phase 1 (Setup)            | 2     |
| Phase 2 (Foundation)       | 3     |
| Story phases               | 3     |
| Polish / cross-cutting     | 2     |
| Parallel opportunities     | 4     |
| MVP scope (P1 stories)     | 3     |

## Dependency Graph

T001 -> T002 -> T003 -> T004 -> T005 -> (T006 || T007[P]) -> T008 -> (T009 || T010[P]) -> T011 -> T012 -> (T013 || T014[P]) -> T015

## Parallelism Guidance

- Tasks marked `[P]` may execute concurrently with sibling tasks in the same phase.
- Tasks with explicit `Depends:` fields must wait for all listed predecessors.
- Builder may reorder within a phase when dependencies allow, but must not skip phases.

## Task Format Rules

Every task MUST use this exact checklist format:

`- [ ] T001 [P] [US1] Description with exact file path`

---

## Phase 1 — Setup

- [x] T001 Confirm F-002 scope fence and exclusions in target/specs/dependency-graph-loading/TASKS.md
  Depends: none
  Acceptance: PLAN.md (scope note line 6) + SPEC.md (In scope / Out of scope)

- [x] T002 Add FR/SC traceability matrix section for F-002 tasks in target/specs/dependency-graph-loading/BUILD.md
  Depends: T001
  Acceptance: SPEC.md FR-001..FR-008 + SC-001..SC-004 mapped for execution evidence

---

## Phase 2 — Foundation

- [x] T003 Define stage boundary contract notes for core service and graph loader in internal/core/service.go
  Depends: T002
  Acceptance: PLAN.md WP1 + ARCHITECTURE (core.Service -> DependencyLoader boundary)

- [x] T004 Define normalized input contract notes for start/direct/production/development handling in internal/core/model.go
  Depends: T003
  Acceptance: SPEC.md Inputs + FR-001..FR-006

- [x] T005 Define loader options contract notes (depth/dev/logger linkage) in internal/graph/loader.go
  Depends: T004
  Acceptance: PLAN.md WP2 + SPEC.md Loader Options entity + FR-002..FR-006

---

## Phase 3 — Story: Load dependency graph from a selected start path (Priority: P1)

- [x] T006 [US1] Add execution task for exact start path passthrough verification in internal/core/service_test.go
  Depends: T005
  Acceptance: SPEC.md US1 Scenario 1 + FR-001 + SC-001

- [x] T007 [P] [US1] Add execution task for successful graph handoff verification to next stage in internal/core/service_test.go
  Depends: T005
  Acceptance: SPEC.md US1 Scenario 2 + FR-008

- [x] T008 [US1] Add fixture-backed integration task for valid project start-path load in internal/graph/npm_ls_loader_test.go
  Depends: T006, T007
  Acceptance: PLAN.md WP4 (loader invocation integrity) + SPEC.md SC-001

---

## Phase 4 — Story: Limit traversal to direct dependencies only (Priority: P2)

- [x] T009 [US2] Add options-matrix test task for direct-only depth mapping (`depth=0`) in internal/core/service_test.go
  Depends: T008
  Acceptance: SPEC.md US2 Scenario 1 + FR-003 + SC-002

- [x] T010 [P] [US2] Add options-matrix test task for recursive depth mapping (non-direct sentinel) in internal/core/service_test.go
  Depends: T008
  Acceptance: SPEC.md US2 Scenario 2 + FR-004 + SC-002

---

## Phase 5 — Story: Control dev-dependency loading behavior for scoped scans (Priority: P3)

- [x] T011 [US3] Add options-matrix test task for default `dev=true` behavior in internal/core/service_test.go
  Depends: T009, T010
  Acceptance: SPEC.md US3 Scenario 1 + FR-005 + SC-003

- [x] T012 [US3] Add options-matrix test task for `dev=false` when production/development flags are set in internal/core/service_test.go
  Depends: T011
  Acceptance: SPEC.md US3 Scenario 2 + FR-006 + SC-003

---

## Phase 6 — Validation

- [x] T013 Add failure-path test task for invalid/unresolvable start path error propagation in internal/core/service_test.go
  Depends: T012
  Acceptance: SPEC.md Edge Case 1 + FR-007 + SC-004

- [x] T014 [P] Add loader adapter cancellation/error-channel verification task in internal/graph/npm_ls_loader_test.go
  Depends: T012
  Acceptance: PLAN.md WP3 (context-aware error semantics) + SPEC.md FR-007

- [x] T015 Add F-002-only verification checklist updates for go test evidence capture in target/specs/dependency-graph-loading/BUILD.md
  Depends: T013, T014
  Acceptance: PLAN.md Validation Approach + SPEC.md SC-001..SC-004

---

## Phase 7 — Polish & Cross-Cutting

- [x] T016 Align F-002 task completion checkboxes with work-package checklist in target/specs/dependency-graph-loading/PLAN.md
  Depends: T015
  Acceptance: PLAN.md WP1..WP4 all traceably linked from executed tasks

- [x] T017 Add final out-of-scope guard note preventing F-003+ drift in target/specs/dependency-graph-loading/BUILD.md
  Depends: T016
  Acceptance: SPEC.md Out of scope section preserved during execution

---

## Implementation Strategy

- **Recommended start**: T001 -> T005 (lock scope, contracts, and traceability before any test execution work).
- **Critical path**: T001 -> T002 -> T003 -> T004 -> T005 -> T006 -> T008 -> T009 -> T011 -> T012 -> T013 -> T015 -> T016 -> T017.
- **Risk areas**: T008 (loader invocation parity), T010 (recursive sentinel assertion), T014 (context cancellation/error semantics).
- **Handoff checklist**:
  - [x] All P1 story tasks complete
  - [x] Validation phase tasks pass
  - [x] BUILD.md evidence artifact started
  - [x] No unresolved NEEDS CLARIFICATION items in task list

## Output Summary

- **Total task count**: 17 (Setup 2, Foundation 3, Story phases 7, Validation 3, Polish 2).
- **Per-story task counts**: US1 (3), US2 (2), US3 (2).
- **Parallel opportunities**: 4 tasks (`T007`, `T010`, `T014`, plus parallel lane at `T006 || T007` and `T009 || T010`).
- **MVP scope (P1)**: `T006`, `T007`, `T008`.
- **Format validation**: Every task uses checklist format with task ID, exact file path, `Depends`, and `Acceptance`.
- **Unresolved items**: None.
