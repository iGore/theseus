# Implementation Tasks: Output Rendering and Export (F-006)

**Date**: 2026-05-01 | **Plan**: `target/specs/output-rendering-and-export/PLAN.md` | **Spec**: `target/specs/output-rendering-and-export/SPEC.md`
**Scope**: Reduced — strictly F-006 rendering/export slice only (Go target)

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 25    |
| Phase 1 (Setup)            | 3     |
| Phase 2 (Foundation)       | 3     |
| Story phases               | 16    |
| Polish / cross-cutting     | 2     |
| Parallel opportunities     | 9     |
| MVP scope (P1 stories)     | 8     |

## Dependency Graph

T001 → T002 → T003
T003 → T004 → T005 → T006
T006 → T007 → T008
T008 → { T009[P], T010[P], T011[P], T012[P], T013[P] } → T014
T014 → T015 → T016
T016 → { T017[P], T018[P] } → T019
T019 → T020 → { T021[P], T022[P] }
T014 + T019 + T022 → T023

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

- [x] T001 Create F-006 task baseline notes in target/specs/output-rendering-and-export/TASKS.md
  Depends: none
  Acceptance: PLAN.md Progress Checklist item “Planning scope confirmed against SPEC.md (F-006 only)”

- [x] T002 Add/confirm F-006 renderer fixture inventory list in test/fixtures/render/README.md
  Depends: T001
  Acceptance: PLAN.md WP2 and Validation Approach golden-test fixture coverage

- [x] T003 Add/confirm F-006 integration scenario inventory in test/integration/render_export/README.md
  Depends: T002
  Acceptance: PLAN.md WP3/WP4/WP5 integration scope and SPEC.md SC-002/SC-003/SC-004

---

## Phase 2 — Foundation

- [x] T004 Define render-mode selection contract in internal/core/model.go
  Depends: T003
  Acceptance: SPEC.md FR-001/FR-002/FR-007 and PLAN.md WP1

- [x] T005 Implement deterministic render-mode resolver in internal/core/service.go
  Depends: T004
  Acceptance: SPEC.md acceptance scenario US1.3 and PLAN.md WP1 precedence requirement

- [x] T006 Add table-driven mode-precedence tests in internal/core/service_render_mode_test.go
  Depends: T005
  Acceptance: SPEC.md SC-001 and PLAN.md WP1 test matrix (including default tree)

---

## Phase 3 — Story: Generate the correct output format (Priority: P1)

- [x] T007 [US1] Implement tree renderer compatibility behavior in internal/render/tree.go
  Depends: T006
  Acceptance: SPEC.md FR-001 and US1 acceptance scenario 2

- [x] T008 [US1] Implement renderer dispatch wiring for all modes in internal/core/service.go
  Depends: T007
  Acceptance: SPEC.md FR-001/FR-002 and PLAN.md WP2 renderer alignment

- [x] T009 [P] [US1] Implement JSON renderer determinism rules in internal/render/json.go
  Depends: T008
  Acceptance: SPEC.md FR-001 and PLAN.md WP2 deterministic ordering contract

- [x] T010 [P] [US1] Implement CSV renderer determinism rules in internal/render/csv.go
  Depends: T008
  Acceptance: SPEC.md FR-001 and PLAN.md WP2 deterministic ordering contract

- [x] T011 [P] [US1] Implement Markdown renderer compatibility behavior in internal/render/markdown.go
  Depends: T008
  Acceptance: SPEC.md FR-001/FR-008 and PLAN.md WP2/WP5

- [x] T012 [P] [US1] Implement Summary renderer compatibility behavior in internal/render/summary.go
  Depends: T008
  Acceptance: SPEC.md FR-001 and PLAN.md WP2

- [x] T013 [P] [US1] Add golden fixtures for all renderer modes in test/fixtures/render/f006_modes/*
  Depends: T008
  Acceptance: PLAN.md WP2 and Validation Approach golden tests

- [x] T014 [US1] Add renderer parity golden tests in test/integration/render_export/render_modes_golden_test.go
  Depends: T009, T010, T011, T012, T013
  Acceptance: SPEC.md SC-001 and US1 independent test

---

## Phase 4 — Story: Persist rendered output to file (Priority: P2)

- [x] T015 [US2] Implement output sink routing (stdout vs --out file) in internal/core/service.go
  Depends: T014
  Acceptance: SPEC.md FR-003/FR-004 and PLAN.md WP3 sink routing

- [x] T016 [US2] Implement parent-directory creation before write in internal/render/files_export.go
  Depends: T015
  Acceptance: SPEC.md FR-004 and US2 acceptance scenario 1

- [x] T017 [P] [US2] Add nested --out path integration test in test/integration/render_export/output_path_creation_test.go
  Depends: T016
  Acceptance: SPEC.md SC-002 and PLAN.md WP3 integration checks

- [x] T018 [P] [US2] Add stdout-default integration test when --out is absent in test/integration/render_export/stdout_default_test.go
  Depends: T016
  Acceptance: SPEC.md FR-003, US2 acceptance scenario 2, and SC-004

- [x] T019 [US2] Add output-content persistence parity test in test/integration/render_export/output_content_parity_test.go
  Depends: T017, T018
  Acceptance: SPEC.md SC-002 and US2 independent test

---

## Phase 5 — Story: Export detected license files per package (Priority: P3)

- [x] T020 [US3] Implement --files export copy flow and layout policy in internal/render/files_export.go
  Depends: T019
  Acceptance: SPEC.md FR-005 and US3 acceptance scenario 1

- [x] T021 [P] [US3] Add --files export integration test with fixture license paths in test/integration/render_export/files_export_enabled_test.go
  Depends: T020
  Acceptance: SPEC.md SC-003 and US3 independent test

- [x] T022 [P] [US3] Add no-export side-effect integration test when --files omitted in test/integration/render_export/files_export_omitted_test.go
  Depends: T020
  Acceptance: SPEC.md US3 acceptance scenario 2

---

## Phase 6 — Validation

- [x] T023 Execute F-006 focused test suite and capture evidence in test/integration/render_export/BUILD_EVIDENCE.md
  Depends: T014, T019, T022
  Acceptance: PLAN.md Validation Approach and SPEC.md SC-001/SC-002/SC-003/SC-004

---

## Phase 7 — Polish & Cross-Cutting

- [x] T024 Add compatibility note for colorized package-key terminal gating in target/specs/output-rendering-and-export/TASKS.md
  Depends: T023
  Acceptance: SPEC.md FR-006 and PLAN.md Open Question on TTY interactivity rule

- [x] T025 Add unresolved doc-status marker for --markdown/--files visibility in target/specs/output-rendering-and-export/TASKS.md
  Depends: T024
  Acceptance: SPEC.md FR-008 and PLAN.md Open Question on documentation alignment

---

## Implementation Strategy

- **Recommended start**: T004 → T006 (mode contract + precedence tests) before renderer implementation.
- **Critical path**: T004 → T005 → T006 → T007 → T008 → T014 → T015 → T016 → T019 → T020 → T022 → T023.
- **Risk areas**: T009–T014 (deterministic renderer parity), T020–T022 (license export layout/collision compatibility), T024 (TTY compatibility ambiguity).
- **Handoff checklist**:
  - [ ] All P1 story tasks complete
  - [ ] Validation phase tasks pass
  - [ ] BUILD.md evidence artifact started
  - [ ] No unresolved NEEDS CLARIFICATION items in task list

## Decomposition Summary

- **Total task count**: 25 (Setup 3, Foundation 3, US1 8, US2 5, US3 3, Validation 1, Polish/Cross-cutting 2).
- **Per-story task counts**: US1 (8), US2 (5), US3 (3).
- **Parallel opportunities**: 9 tasks (`T009–T013`, `T017`, `T018`, `T021`, `T022`) across story phases.
- **MVP scope (P1 only)**: `T007–T014` (with prerequisite foundation `T004–T006`) delivers viable F-006 rendering mode slice.
- **Format validation**: Every task line uses checklist format and includes task ID, exact file path, `Depends`, and `Acceptance` details.
- **Unresolved items**:
  - `T024` — **NEEDS CLARIFICATION**: final TTY interactivity rule for colorized package keys (FR-006).
  - `T025` — **NEEDS CLARIFICATION**: documentation-status handling for `--markdown` and `--files` (FR-008).

## Builder Notes (F-006 execution)

- Colorized package-key output is applied only in tree mode when `Colorize=true` and `IsTerminal=true`; structured outputs (json/csv/markdown/summary) never apply ANSI color formatting.
- `--markdown` and `--files` behavior is implemented and tested for compatibility; documentation promotion decision remains open per PLAN.md Open Questions.
