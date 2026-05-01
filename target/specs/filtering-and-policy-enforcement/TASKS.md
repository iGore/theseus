# Implementation Tasks: F-005 Filtering, package restriction, and policy exits

**Date**: 2026-05-01 | **Plan**: `target/specs/filtering-and-policy-enforcement/PLAN.md` | **Spec**: `target/specs/filtering-and-policy-enforcement/SPEC.md`
**Scope**: Reduced — strictly F-005 only (Go target)

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 22    |
| Phase 1 (Setup)            | 2     |
| Phase 2 (Foundation)       | 4     |
| Story phases               | 11    |
| Polish / cross-cutting     | 2     |
| Parallel opportunities     | 7     |
| MVP scope (P1 stories)     | 4     |

## Dependency Graph

`T001 -> T002 -> T003 -> {T004 [P], T005 [P], T006} -> T007 -> {T008 [US1], T009 [P][US1], T010 [P][US1]} -> T011 [US1] -> {T012 [US2], T013 [P][US2], T014 [P][US2], T015 [US2]} -> {T016 [US3], T017 [P][US3]} -> T018 [US3] -> {T019, T020 [P]} -> T021 -> T022`

## Parallelism Guidance

- Tasks marked `[P]` may execute concurrently with sibling tasks in the same phase.
- Tasks with explicit `Depends:` fields must wait for all listed predecessors.
- Builder may reorder within a phase when dependencies allow, but must not skip phases.

## Task Format Rules

Every task MUST use this exact checklist format:

`- [ ] T001 [P] [US1] Description with exact file path`

Format rules:

- Start every task with `- [ ]`
- Use sequential task IDs in execution order: `T001`, `T002`, `T003`, ...
- Add `[P]` only when the task is parallelizable
- Add `[US1]`, `[US2]`, `[US3]`, etc. only for story-phase tasks
- Include the exact file path directly in the task description
- Keep dependency details in a `Depends:` note directly below the task when needed

---

## Phase 1 — Setup

- [x] T001 Confirm F-005 execution boundary and stage ordering references in target/specs/filtering-and-policy-enforcement/TASKS.md
  Depends: none
  Acceptance: PLAN.md line 6 and SPEC.md Scope and Traceability section

- [x] T002 Create task-traceability header comments for FR/SC mapping in internal/filter/filter_test.go
  Depends: T001
  Acceptance: PLAN WP5 and SPEC.md FR-001..FR-009 traceability preserved in test file comments

---

## Phase 2 — Foundation

- [x] T003 Define F-005 runtime option contract fields in internal/core/model.go
  Depends: T002
  Acceptance: PLAN WP1 and SPEC.md Inputs and Outputs section

- [x] T004 [P] Add filter-stage invocation boundary in internal/core/service.go
  Depends: T003
  Acceptance: PLAN WP1 and FR-008 (filters before policy evaluation)

- [x] T005 [P] Add policy-stage invocation boundary in internal/core/service.go
  Depends: T003
  Acceptance: PLAN WP1 and FR-008 (policy stage after filters)

- [x] T006 Add mutual-exclusion validation path for failOn+onlyAllow in internal/cli/options.go
  Depends: T003
  Acceptance: PLAN WP4 and SPEC.md FR-009

---

## Phase 3 — Story: Enforce compliance gates in CI (Priority: P1)

- [x] T007 [US1] Add policy evaluation entrypoint and result contract in internal/filter/policy.go
  Depends: T005
  Acceptance: PLAN WP3 and SPEC.md FR-006/FR-007

- [x] T008 [US1] Implement failOn first-violation detection path in internal/filter/policy.go
  Depends: T007
  Acceptance: SPEC.md Acceptance Scenario 1 for User Story 1 and SC-003

- [x] T009 [P] [US1] Add onlyAllow matcher strategy seam in internal/filter/policy.go marked NEEDS CLARIFICATION for OD-001
  Depends: T007
  Acceptance: PLAN OD-001 and SPEC.md Open Decisions OD-001

- [x] T010 [P] [US1] Wire policy violation-to-exit strategy hook for code 1 in internal/core/service.go
  Depends: T007
  Acceptance: SPEC.md FR-006/FR-007 and C-003

- [x] T011 [US1] Add CLI RunE invalid-policy-combination error return path in internal/cli/command.go
  Depends: T006, T010
  Acceptance: SPEC.md FR-009 and PLAN WP4

---

## Phase 4 — Story: Restrict scan results to relevant dependency subset (Priority: P2)

- [x] T012 [US2] Implement exclude license expression filtering pipeline step in internal/filter/filter.go
  Depends: T004
  Acceptance: SPEC.md FR-003 and User Story 2 Acceptance Scenario 1

- [x] T013 [P] [US2] Implement package include-list (`packages`) semicolon key restriction in internal/filter/filter.go
  Depends: T004
  Acceptance: SPEC.md FR-004 and User Story 2 Acceptance Scenario 2

- [x] T014 [P] [US2] Implement package exclude-list (`excludePackages`) semicolon key restriction in internal/filter/filter.go
  Depends: T004
  Acceptance: SPEC.md FR-004 and SC-004

- [x] T015 [US2] Implement private package exclusion stage in internal/filter/filter.go
  Depends: T012, T013, T014
  Acceptance: SPEC.md FR-005 and User Story 2 Acceptance Scenario 3

---

## Phase 5 — Story: Focus unknown-license investigation (Priority: P3)

- [x] T016 [US3] Implement unknown transform (`*` suffix to `UNKNOWN`) stage in internal/filter/filter.go
  Depends: T004
  Acceptance: SPEC.md FR-001 and User Story 3 Acceptance Scenario 1

- [x] T017 [P] [US3] Implement onlyunknown restriction stage in internal/filter/filter.go
  Depends: T016
  Acceptance: SPEC.md FR-002 and User Story 3 Acceptance Scenario 2

- [x] T018 [US3] Enforce full F-005 filter-stage order before policy stage in internal/core/service.go
  Depends: T008, T012, T013, T014, T015, T016, T017
  Acceptance: SPEC.md FR-008 and edge-case order note

---

## Phase 6 — Validation

- [x] T019 Add table-driven unit tests for unknown/onlyunknown/exclude/package/private filtering in internal/filter/filter_test.go
  Depends: T015, T017
  Acceptance: PLAN V1, V2, V3, V4, V5 and SPEC.md SC-001, SC-002, SC-004

- [x] T020 [P] Add policy fail-fast and invalid-combination tests in internal/filter/policy_test.go
  Depends: T009, T011
  Acceptance: PLAN V6, V8 and SPEC.md SC-003 plus FR-009

- [x] T021 Add integration fixture tests for filter-before-policy ordering and exit behavior in test/integration/filter_policy_integration_test.go
  Depends: T018, T019, T020
  Acceptance: PLAN V7 and SPEC.md FR-008/SC-003

---

## Phase 7 — Polish & Cross-Cutting

- [x] T022 Finalize OD-001 matcher decision notes and lock expected tests in target/specs/filtering-and-policy-enforcement/TASKS.md
  Depends: T021
  Acceptance: PLAN WP6 and SPEC.md OD-001 resolved or explicitly blocked as NEEDS CLARIFICATION

---

## Implementation Strategy

- **Recommended start**: T001 -> T003 -> T004/T005 to establish contract and stage boundaries first.
- **Critical path**: T001 -> T002 -> T003 -> T005 -> T007 -> T008 -> T018 -> T021 -> T022.
- **Risk areas**: T009/T022 (`onlyAllow` strictness OD-001), T012 (escaped-comma + BSD compatibility), T010/T011 (exit-code compatibility).
- **Handoff checklist**:
  - [x] All P1 story tasks complete
  - [x] Validation phase tasks pass
  - [x] BUILD.md evidence artifact started
  - [x] No unresolved NEEDS CLARIFICATION items in task list

## Decomposition Report

- **Total task count**: 22 (Setup 2, Foundation 4, Story phases 11, Validation 3, Polish/Cross-cutting 2).
- **Per-story task counts**: US1 (5 tasks: T007-T011), US2 (4 tasks: T012-T015), US3 (3 tasks: T016-T018).
- **Parallel opportunities**: 7 tasks marked `[P]` (T004, T005, T009, T010, T013, T014, T017, T020; parallel lanes in Foundation, US1, US2, US3, Validation).
- **MVP scope (P1 stories)**: T007-T011 (plus prerequisites T001-T006) deliver first viable CI compliance-gate slice.
- **Format validation**: Every task uses checklist format and includes task ID, exact file path, `Depends`, and `Acceptance`.
- **Unresolved items**: OD-001 strictness decision remains explicitly deferred; compatibility substring seam is implemented and marked in code for later hardening.
