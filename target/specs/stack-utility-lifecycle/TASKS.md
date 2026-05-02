# Implementation Tasks: Ancillary Stack Utility Lifecycle (F-008)

**Date**: 2026-05-01 | **Plan**: `target/specs/stack-utility-lifecycle/PLAN.md` | **Spec**: `target/specs/stack-utility-lifecycle/SPEC.md`
**Scope**: Reduced — F-008 only (user-constrained)

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 17    |
| Phase 1 (Setup)            | 2     |
| Phase 2 (Foundation)       | 2     |
| Story phases               | 9     |
| Polish / cross-cutting     | 2     |
| Parallel opportunities     | 5     |
| MVP scope (P1 stories)     | 4     |

## Dependency Graph

`T001 -> T002 -> T003 -> T004 -> T005 -> {T006 | T007 | T008} -> {T009/T010 or T011/T012 or T013/T014} -> T015 -> T016 -> T017`

Parallel lanes:
- `[P]` `T006`, `T007`, `T008` are decision-branch stubs; execute only the branch selected in `T005`.
- `[P]` `T009` and `T010` can run together after `T006` (Retain branch).
- `[P]` `T011` and `T012` can run together after `T007` (Deprecate branch).
- `[P]` `T013` and `T014` can run together after `T008` (Remove branch).

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

- [x] T001 Confirm F-008-only execution gate and branch policy in target/specs/stack-utility-lifecycle/TASKS.md
  Depends: none
  Acceptance: PLAN.md Open Questions line 32 and Dependencies lines 89-94 acknowledged as execution gate

- [x] T002 Create Builder evidence scaffold for F-008 in target/specs/stack-utility-lifecycle/BUILD.md
  Depends: T001
  Acceptance: PLAN.md WP3 deliverable line 85 traceability handoff artifact prepared

---

## Phase 2 — Foundation

- [x] T003 Add lifecycle decision record template section for F-008 in target/specs/stack-utility-lifecycle/BUILD.md
  Depends: T002
  Acceptance: SPEC.md FR-001, FR-002, FR-006 and PLAN.md WP1 lines 64-69 represented in template

- [x] T004 Add FR/SC traceability matrix template for F-008 in target/specs/stack-utility-lifecycle/BUILD.md
  Depends: T003
  Acceptance: PLAN.md WP3 lines 81-85 and SPEC.md SC-001..SC-004 mapped with placeholders

---

## Phase 3 — Story: Decide lifecycle with explicit compatibility stance (Priority: P1)

- [x] T005 [US1] Record exactly one approved lifecycle status (Retain/Deprecate/Remove) and rationale in target/specs/stack-utility-lifecycle/BUILD.md
  Depends: T004
  Acceptance: SPEC.md User Story 1 scenario 1, FR-001, FR-002, SC-002

- [x] T006 [P] [US1] If status is Retain, document selected-path gate and branch activation in target/specs/stack-utility-lifecycle/BUILD.md
  Depends: T005
  Acceptance: PLAN.md WP2 line 73 selected-path criteria captured without contradiction

- [x] T007 [P] [US1] If status is Deprecate, document selected-path gate and branch activation in target/specs/stack-utility-lifecycle/BUILD.md *(N/A-complete: Retain path selected in T005; branch explicitly closed)*
  Depends: T005
  Acceptance: PLAN.md WP2 line 74 selected-path criteria captured without contradiction

- [x] T008 [P] [US1] If status is Remove, document selected-path gate and branch activation in target/specs/stack-utility-lifecycle/BUILD.md *(N/A-complete: Retain path selected in T005; branch explicitly closed)*
  Depends: T005
  Acceptance: PLAN.md WP2 line 75 selected-path criteria captured without contradiction

---

## Phase 4 — Story: Preserve behavior if retained (Priority: P2)

- [x] T009 [P] [US2] If Retain branch active, define Stack API parity checklist for target/internal/legacy/stack.go in target/specs/stack-utility-lifecycle/BUILD.md
  Depends: T006
  Acceptance: SPEC.md FR-003 and User Story 2 scenario 1 contract elements (`Stack`, `add`, `test`, `done`) listed

- [x] T010 [P] [US2] If Retain branch active, define retain-path verification tasks for target/internal/legacy/stack_test.go in target/specs/stack-utility-lifecycle/BUILD.md
  Depends: T006
  Acceptance: SPEC.md SC-003 and PLAN.md Validation Approach lines 99-100 reflected as test evidence checklist

---

## Phase 5 — Story: De-risk removal/deprecation path (Priority: P3)

- [x] T011 [P] [US3] If Deprecate branch active, define one-release compatibility shim tasks for target/internal/legacy/stack.go in target/specs/stack-utility-lifecycle/BUILD.md *(N/A-complete: Deprecate branch inactive under Retain decision)*
  Depends: T007
  Acceptance: SPEC.md FR-004 and PLAN.md WP2 line 74 documented as must-pass criteria

- [x] T012 [P] [US3] If Deprecate branch active, define deprecation notice and migration-text tasks for target/specs/stack-utility-lifecycle/BUILD.md *(N/A-complete: Deprecate branch inactive under Retain decision)*
  Depends: T007
  Acceptance: SPEC.md User Story 3 scenario 1 and SC-004 documented with release-artifact checklist

- [x] T013 [P] [US3] If Remove branch active, define removal-impact and migration-guidance tasks for target/specs/stack-utility-lifecycle/BUILD.md *(N/A-complete: Remove branch inactive under Retain decision)*
  Depends: T008
  Acceptance: SPEC.md FR-005 and User Story 3 scenario 1 documented with breaking-change requirements

- [x] T014 [P] [US3] If Remove branch active, define release-note evidence checklist for Stack removal in target/specs/stack-utility-lifecycle/BUILD.md *(N/A-complete: Remove branch inactive under Retain decision)*
  Depends: T008
  Acceptance: SPEC.md SC-004 and PLAN.md Validation Approach line 100 represented as required evidence

---

## Phase 6 — Validation

- [x] T015 Execute F-008 traceability validation across selected branch in target/specs/stack-utility-lifecycle/BUILD.md
  Depends: T009, T010, T011, T012, T013, T014
  Acceptance: SPEC.md SC-001 and PLAN.md WP3 line 82 satisfied for active branch and marked N/A for inactive branches

- [x] T016 Verify unresolved external-usage confidence item remains explicitly classified in target/specs/stack-utility-lifecycle/BUILD.md
  Depends: T015
  Acceptance: SPEC.md FR-006 and PLAN.md Open Questions lines 33-34 retained or explicitly resolved with evidence

---

## Phase 7 — Polish & Cross-Cutting

- [x] T017 Finalize Builder handoff checklist and execution-ready status in target/specs/stack-utility-lifecycle/TASKS.md
  Depends: T016
  Acceptance: PLAN.md lines 97-102 validation checklist fully mapped and no contradictory lifecycle statements remain

---

## Implementation Strategy

- **Recommended start**: `T001 -> T005` to force lifecycle decision before any branch work.
- **Critical path**: `T001 -> T002 -> T003 -> T004 -> T005 -> (one of T006/T007/T008) -> selected branch tasks -> T015 -> T016 -> T017`.
- **Risk areas**: `T005` (decision approval), `T015` (branch-aware traceability completeness), `T016` (external usage confidence ambiguity).
- **Handoff checklist**:
  - [x] All P1 story tasks complete
  - [x] Validation phase tasks pass
  - [x] BUILD.md evidence artifact started
  - [x] No unresolved NEEDS CLARIFICATION items in task list *(for execution checklist status; branch-conditional items are N/A-complete under Retain)*

## Decomposition Report

- **Total task count**: 17 (Setup 2, Foundation 2, Story phases 9, Validation 2, Polish 1).
- **Per-story task counts**: US1 (4), US2 (2), US3 (4; branch-conditional).
- **Parallel opportunities**: 5 tasks (`T006`, `T007`, `T008`, plus two-task lanes in each active branch pair).
- **MVP scope**: P1 coverage is `T005` plus selected-branch activation task (`T006` or `T007` or `T008`) and branch-completion validation in `T015`.
- **Format validation**: Every task uses checklist format with task ID, exact file path, `Depends:`, and `Acceptance:`.
- **Unresolved items**:
  - None for execution checklist closure; lifecycle status approved as **Retain** and non-selected branch tasks are recorded as N/A-complete.
