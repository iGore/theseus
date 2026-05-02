# Implementation Tasks: F-003 Flatten dependency graph into module inventory

**Date**: 2026-05-01 | **Plan**: `target/specs/dependency-flattening/PLAN.md` | **Spec**: `target/specs/dependency-flattening/SPEC.md`
**Scope**: Reduced — user-requested scope narrowed strictly to F-003 only (Go target only)

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 22    |
| Phase 1 (Setup)            | 3     |
| Phase 2 (Foundation)       | 4     |
| Story phases               | 11    |
| Polish / cross-cutting     | 2     |
| Parallel opportunities     | 7     |
| MVP scope (P1 stories)     | 4     |

## Dependency Graph

T001 -> T002 -> T003 -> T004 -> T005 -> T006 -> T007

T007 -> T008 [US1] -> T009 [P][US1]
T008 -> T010 [US1]
T010 -> T011 [P][US1]

T010 -> T012 [US2] -> T013 [P][US2]
T012 -> T014 [US2]
T014 -> T015 [P][US2]

T010 -> T016 [US3]
T016 -> T017 [P][US3]
T016 -> T018 [US3]

T011,T015,T018 -> T019 -> T020 -> T021 -> T022

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

- [x] T001 Create F-003 implementation checklist notes in target/specs/dependency-flattening/PLAN.md
  Depends: none
  Acceptance: PLAN.md WP1 checklist alignment and scope note for F-003 only

- [x] T002 Add flattening fixture inventory plan in target/specs/dependency-flattening/quickstart.md
  Depends: T001
  Acceptance: SPEC.md AC-001..AC-005 fixture coverage plan documented

- [x] T003 Create flattening contracts directory scaffold in target/specs/dependency-flattening/contracts/flattening-contract.md
  Depends: T001
  Acceptance: PLAN.md WP1 contract-baseline artifact path established

---

## Phase 2 — Foundation

- [x] T004 Define F-003 domain type updates in internal/core/model.go
  Depends: T003
  Acceptance: SPEC.md Key Entities and FR-001/FR-006 represented in Go model contracts

- [x] T005 [P] Define flatten stage interface contract in internal/core/service.go
  Depends: T004
  Acceptance: SPEC.md FR-006 stage boundary preserved (flatten output consumed downstream)

- [x] T006 Create flatten package entrypoint skeleton in internal/flatten/flatten.go
  Depends: T004
  Acceptance: PLAN.md summary target file path and ARCHITECTURE section 3 structure decision

- [x] T007 [P] Add table-driven test file scaffold for flatten stage in internal/flatten/flatten_test.go
  Depends: T006
  Acceptance: PLAN.md WP5 validation design and Context7-backed testify testing pattern applied

---

## Phase 3 — Story: Produce canonical module inventory from nested dependencies (Priority: P1)

- [x] T008 [US1] Implement recursive dependency traversal in internal/flatten/flatten.go
  Depends: T006
  Acceptance: SPEC.md FR-002 and AC-001 deep dependency traversal behavior

- [x] T009 [P] [US1] Implement canonical key construction name@version in internal/flatten/flatten.go
  Depends: T008
  Acceptance: SPEC.md FR-001 and Compatibility Note key format `name@version`

- [x] T010 [US1] Implement flat inventory map population contract in internal/flatten/flatten.go
  Depends: T009
  Acceptance: SPEC.md FR-001/FR-006 and AC-005 downstream-consumable map output

- [x] T011 [P] [US1] Add P1 fixture tests for nested-to-flat transformation in internal/flatten/flatten_test.go
  Depends: T010
  Acceptance: SPEC.md AC-001 and SC-001 pass with table-driven tests

---

## Phase 4 — Story: Prevent duplicate/circular expansion from corrupting inventory (Priority: P2)

- [x] T012 [US2] Implement visited-key guard semantics in internal/flatten/flatten.go
  Depends: T010
  Acceptance: SPEC.md FR-003 one-entry-per-key behavior

- [x] T013 [P] [US2] Implement cycle short-circuit path handling in internal/flatten/flatten.go
  Depends: T012
  Acceptance: SPEC.md AC-002 recursion branch short-circuit without non-termination

- [x] T014 [US2] Enforce deterministic uniqueness outcomes for repeated keys in internal/flatten/flatten.go
  Depends: T012
  Acceptance: SPEC.md FR-007 deterministic key-set uniqueness outcome

- [x] T015 [P] [US2] Add duplicate-path and cycle fixtures in internal/flatten/flatten_test.go
  Depends: T013, T014
  Acceptance: SPEC.md AC-002 and SC-002 with exactly one map entry per canonical key

---

## Phase 5 — Story: Preserve module metadata needed by later stages (Priority: P3)

- [x] T016 [US3] Implement metadata extraction mapping for repository/author/url/path/private fields in internal/flatten/flatten.go
  Depends: T010
  Acceptance: SPEC.md FR-005 and AC-004 metadata propagation to inventory entries

- [x] T017 [P] [US3] Implement prod/dev inclusion gating before map insertion in internal/flatten/flatten.go
  Depends: T016
  Acceptance: SPEC.md FR-004 and AC-003 excluded nodes not inserted

- [x] T018 [US3] Add metadata and prod/dev gating fixture tests in internal/flatten/flatten_test.go
  Depends: T016, T017
  Acceptance: SPEC.md AC-003/AC-004 and SC-003/SC-004 pass

---

## Phase 6 — Validation

- [x] T019 Add flatten-stage integration contract test harness in test/integration/flatten_contract_test.go
  Depends: T011, T015, T018
  Acceptance: SPEC.md AC-005 downstream harness consumes flattened map without nested tree dependency

- [x] T020 Add deterministic key-set comparison helpers for tests in internal/flatten/flatten_test.go
  Depends: T019
  Acceptance: PLAN.md WP5 determinism rule and Context7 `/golang/go` map-order constraint honored

---

## Phase 7 — Polish & Cross-Cutting

- [x] T021 Document unresolved parity decisions OD-001 and OD-002 in target/specs/dependency-flattening/contracts/flattening-contract.md
  Depends: T020
  Acceptance: SPEC.md Open Decisions section reflected as NEEDS CLARIFICATION blockers

- [x] T022 Record F-003 verification walkthrough and evidence checklist in target/specs/dependency-flattening/quickstart.md
  Depends: T021
  Acceptance: PLAN.md validation approach and SPEC.md success criteria traceable in execution notes

---

## Implementation Strategy

- **Recommended start**: T001 -> T004 -> T006 to lock scope and contracts before algorithm work.
- **Critical path**: T001 -> T003 -> T004 -> T006 -> T008 -> T009 -> T010 -> T012 -> T014 -> T015 -> T019 -> T020 -> T021 -> T022.
- **Risk areas**: T009 (key construction parity), T013 (cycle termination correctness), T016 (metadata contract completeness vs OD-002), T017 (gate point correctness).
- **Handoff checklist**:
  - [x] All P1 story tasks complete
  - [x] Validation phase tasks pass
  - [x] BUILD.md evidence artifact started
  - [ ] No unresolved NEEDS CLARIFICATION items in task list

## Decomposition Report

- **Total task count**: 22 (Setup 3, Foundation 4, Story phases 11, Validation 2, Polish 2).
- **Per-story task counts**: US1 (4), US2 (4), US3 (3).
- **Parallel opportunities**: 7 tasks (`T005`, `T007`, `T009`, `T011`, `T013`, `T015`, `T017`).
- **MVP scope (P1)**: `T008`, `T009`, `T010`, `T011`.
- **Format validation**: Every task line follows checklist format with ID and file path, and every task includes `Depends:` and `Acceptance:` notes.
- **Unresolved items**: `T021` tracks `OD-001` and `OD-002` as `NEEDS CLARIFICATION` and blocks strict parity closure.
