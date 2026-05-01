# Implementation Tasks: CLI options model and normalization (F-001)

**Date**: 2026-05-01 | **Plan**: `target/specs/cli-options-model/PLAN.md` | **Spec**: `target/specs/cli-options-model/SPEC.md`
**Scope**: Reduced — F-001 only (user-constrained)

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 16    |
| Phase 1 (Setup)            | 2     |
| Phase 2 (Foundation)       | 3     |
| Story phases               | 3     |
| Polish / cross-cutting     | 1     |
| Parallel opportunities     | 5     |
| MVP scope (P1 stories)     | 4     |

## Dependency Graph

T001 -> T002 -> T003 -> T004 -> T005

T005 -> T006 -> {T007 [P], T008} -> T009

T009 -> T010 -> {T011 [P], T012} -> T013

T013 -> T014 -> {T015 [P], T016 [P]}

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

- [x] T001 Create F-001 task traceability matrix scaffold in target/specs/cli-options-model/TASKS.md
  Depends: none
  Acceptance: PLAN.md WP-1..WP-5 traceability baseline established

- [x] T002 Define F-001 CLI test fixture table skeleton in internal/cli/options_test.go
  Depends: T001
  Acceptance: PLAN.md WP-5 and SPEC.md SC-001 test-slice coverage prepared

---

## Phase 2 — Foundation

- [x] T003 Align F-001 option contracts for RawArgs and NormalizedOptions in internal/core/model.go
  Depends: T002
  Acceptance: SPEC.md FR-001 and Entities RawArgs/NormalizedOptions mapped

- [x] T004 [P] Define GuardrailResult outcome contract and compatibility marker fields in internal/core/model.go
  Depends: T003
  Acceptance: SPEC.md FR-005, FR-006, FR-008 and Entities GuardrailResult mapped

- [x] T005 Wire CLI-to-core option mapping boundary in internal/cli/options.go
  Depends: T003, T004
  Acceptance: PLAN.md WP-1 boundary decision and ARCHITECTURE section 3 cli/core separation preserved

---

## Phase 3 — Story: Parse and normalize options for a scan run (Priority: P1)

- [x] T006 [US1] Implement default start-path normalization using current working directory in internal/cli/options.go
  Depends: T005
  Acceptance: SPEC.md FR-002 and US1 Scenario 1

- [x] T007 [P] [US1] Implement direct traversal depth normalization (`direct=true -> 0`, else full-traversal sentinel) in internal/cli/options.go
  Depends: T006
  Acceptance: SPEC.md FR-003 and US1 Scenario 2

- [x] T008 [US1] Implement structured-output color default override (`json|csv|markdown` disables color) in internal/cli/options.go
  Depends: T006
  Acceptance: SPEC.md FR-004 and US1 Scenario 3

- [x] T009 [US1] Add parser/normalizer table tests for no-args defaults and deterministic normalized output in internal/cli/options_test.go
  Depends: T007, T008
  Acceptance: SPEC.md SC-001 and SC-003

---

## Phase 4 — Story: Enforce early CLI guardrails before scanning (Priority: P2)

- [x] T010 [US2] Add incompatible flag guard (`failOn` + `onlyAllow`) returning pre-scan non-zero error path in internal/cli/command.go
  Depends: T009
  Acceptance: SPEC.md FR-005 and US2 Scenario 1

- [x] T011 [P] [US2] Add comma-delimiter guidance warning behavior for `failOn` and `onlyAllow` in internal/cli/command.go
  Depends: T010
  Acceptance: SPEC.md FR-006 and US2 Scenario 2

- [x] T012 [US2] Add guardrail command-path tests for conflicting policy flags and delimiter warning output in internal/cli/command_test.go
  Depends: T010, T011
  Acceptance: SPEC.md SC-002 and Edge Cases delimiter handling

---

## Phase 5 — Story: Access utility CLI meta commands (Priority: P3)

- [x] T013 [US3] Implement deterministic help and version meta-command execution path via Cobra RunE lifecycle in internal/cli/command.go
  Depends: T012
  Acceptance: SPEC.md FR-007 and US3 Scenario 1

- [x] T014 [US3] Add compatibility toggle marker for unresolved `--version` exit semantics decision in internal/cli/command.go
  Depends: T013
  Acceptance: SPEC.md FR-008 and PLAN.md Open Questions entry; NEEDS CLARIFICATION retained

- [x] T015 [P] [US3] Add help/version deterministic command-path tests including repeated-run assertions in internal/cli/command_test.go
  Depends: T013
  Acceptance: SPEC.md SC-004 deterministic behavior

- [x] T016 [P] [US3] Add pending-branch assertion for legacy `--version` exit `1` versus normalized `0` in internal/cli/command_test.go
  Depends: T014
  Acceptance: SPEC.md FR-008 and Open Questions item 1

---

## Phase 6 — Validation

- [x] T017 Execute F-001 focused test matrix and map FR-001..FR-008 plus SC-001..SC-004 evidence in target/specs/cli-options-model/BUILD.md
  Depends: T009, T012, T015, T016
  Acceptance: PLAN.md WP-5 traceability table completed

---

## Phase 7 — Polish & Cross-Cutting

- [x] T018 Document F-001 compatibility decision gate and unresolved `--version` exit blocker in target/specs/cli-options-model/PLAN.md
  Depends: T017
  Acceptance: PLAN.md Open Questions and SPEC.md FR-008 alignment retained without silent resolution

---

## Implementation Strategy

- **Recommended start**: T001-T005 to lock contracts and boundaries before story behavior.
- **Critical path**: T001 -> T002 -> T003 -> T005 -> T006 -> T008 -> T009 -> T010 -> T012 -> T013 -> T014 -> T016 -> T017 -> T018.
- **Risk areas**: T014/T016 (`--version` exit semantics still unresolved), T010 (error-path compatibility), T013 (meta-command deterministic behavior).
- **Handoff checklist**:
  - [x] All P1 story tasks complete
  - [x] Validation phase tasks pass
  - [x] BUILD.md evidence artifact started
  - [ ] No unresolved NEEDS CLARIFICATION items in task list

## Task Decomposer Report

- **Total task count**: 18 (Setup 2, Foundation 3, Story phases 11, Validation 1, Polish 1).
- **Per-story task counts**: US1(P1)=4 (T006-T009), US2(P2)=3 (T010-T012), US3(P3)=4 (T013-T016).
- **Parallel opportunities**: 5 tasks (`T004`, `T007`, `T011`, `T015`, `T016`) across Foundation/US1/US2/US3 phases.
- **MVP scope**: P1 delivery slice is T006-T009 (after foundational T001-T005).
- **Format validation**: Every task uses checklist syntax with task ID, explicit file path, `Depends`, and `Acceptance` notes.
- **Unresolved items**: `NEEDS CLARIFICATION` remains at T014 and T016 for `--version` exit code (`1` vs `0`).
