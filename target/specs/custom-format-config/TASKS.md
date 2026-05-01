# Implementation Tasks: Custom Format and JSON Config Loading (F-007)

**Date**: 2026-05-01 | **Plan**: `target/specs/custom-format-config/PLAN.md` | **Spec**: `target/specs/custom-format-config/SPEC.md`
**Scope**: Reduced — F-007 only (user-requested scope constraint)

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 19    |
| Phase 1 (Setup)            | 3     |
| Phase 2 (Foundation)       | 3     |
| Story phases               | 9     |
| Polish / cross-cutting     | 2     |
| Parallel opportunities     | 8     |
| MVP scope (P1 stories)     | 4     |

## Dependency Graph

T001 -> T002 -> T003 -> T004 -> T005 -> T006
T006 -> T007 -> T008 -> T009 -> T010
T006 -> T011 -> T012
T010 + T012 -> T013 -> T014 -> T015
T014 + T015 -> T016 -> T017
T017 -> T018 -> T019

Parallel lanes:
- [P] T008 with [P] T009 (after T007)
- [P] T011 with T007/T008/T009 (after T006)
- [P] T012 with T008/T009 (after T011)
- [P] T015 with [P] T014 (after T013)
- [P] T018 with [P] T019 planning (execution after T017)

## Parallelism Guidance

- Tasks marked `[P]` may execute concurrently with sibling tasks in the same phase.
- Tasks with explicit `Depends:` fields must wait for all listed predecessors.
- Builder may reorder within a phase when dependencies allow, but must not skip phases.

## Task Format Rules

Every task MUST use this exact checklist format:

`- [ ] T001 [P] [US1] Description with exact file path`

---

## Phase 1 — Setup

- [x] T001 Confirm F-007-only execution scope in target/specs/custom-format-config/TASKS.md
  Depends: none
  Acceptance: PLAN.md scope note (line 6), SPEC.md assumptions (line 92)

- [x] T002 Document absence of workflow extension hooks from .specify/extensions.yml in target/specs/custom-format-config/TASKS.md
  Depends: T001
  Acceptance: Task-decomposer hook check requirement satisfied (no .specify directory present)

- [x] T003 Document missing optional supporting artifacts status in target/specs/custom-format-config/TASKS.md
  Depends: T002
  Acceptance: PLAN.md project structure section (lines 65-73) and artifact loading rule coverage

---

## Phase 2 — Foundation

- [x] T004 Define CustomFormatConfig and ParsedCustomConfigResult type contracts in internal/config/custom_format.go
  Depends: T003
  Acceptance: PLAN.md WP1, SPEC.md Entities (CustomFormatConfig, ParsedCustomConfigResult), FR-001, FR-004

- [x] T005 Define module custom-field projection boundary contract in internal/flatten/flatten.go
  Depends: T004
  Acceptance: PLAN.md WP1/WP2 boundary ownership, SPEC.md Entity ModuleCustomFields

- [x] T006 Define CLI-to-config ownership boundary for customPath/customFormat in internal/cli/options.go
  Depends: T005
  Acceptance: PLAN.md WP1 ownership boundaries, ARCHITECTURE section 3 boundary rule (cli parses only)

---

## Phase 3 — Story: Add custom fields via inline format object (Priority: P1)

- [x] T007 [US1] Define deterministic custom key mapping routine specification in internal/config/custom_format.go
  Depends: T006
  Acceptance: SPEC.md FR-001, FR-002; SC-001; User Story 1 acceptance scenarios 1-2

- [x] T008 [P] [US1] Define default fallback behavior tests for missing metadata keys in internal/config/custom_format_test.go
  Depends: T007
  Acceptance: SPEC.md FR-002, SC-001, User Story 1 independent test

- [x] T009 [P] [US1] Define module output field presence assertions for configured keys in test/integration/custom_format_inline_test.go
  Depends: T007
  Acceptance: SPEC.md SC-001; PLAN.md WP5 integration coverage

- [x] T010 [US1] Define licenseText and copyright custom-field enrichment hooks in internal/flatten/flatten.go
  Depends: T008, T009
  Acceptance: SPEC.md FR-008, FR-009; PLAN.md WP2

---

## Phase 4 — Story: Load custom format from JSON file path (Priority: P2)

- [x] T011 [P] [US2] Define customPath JSON load sequence before output shaping in internal/config/custom_format.go
  Depends: T006
  Acceptance: SPEC.md FR-004, FR-005; User Story 2 acceptance scenario 1

- [x] T012 [P] [US2] Define valid customPath integration fixture coverage in test/fixtures/custom-format/valid_custom_path.json
  Depends: T011
  Acceptance: SPEC.md SC-002; User Story 2 acceptance scenario 2; PLAN.md WP5

---

## Phase 5 — Story: Preserve error and opt-out semantics in config handling (Priority: P3)

- [x] T013 [US3] Resolve FR-010 invalid customPath behavior decision record in target/specs/custom-format-config/PLAN.md
  Depends: T010, T012
  Acceptance: SPEC.md FR-010; PLAN.md WP4; mark `NEEDS CLARIFICATION` resolved or explicitly blocked

- [x] T014 [P] [US3] Define parseJson error-matrix tests for non-string/missing/malformed inputs in internal/config/custom_format_test.go
  Depends: T013
  Acceptance: SPEC.md FR-006, FR-007; SC-003; User Story 3 acceptance scenario 1

- [x] T015 [P] [US3] Define false-property exclusion behavior tests in test/integration/custom_format_exclusion_test.go
  Depends: T013
  Acceptance: SPEC.md FR-003; SC-004; User Story 3 acceptance scenario 2

---

## Phase 6 — Validation

- [x] T016 Execute bounded F-007 unit test suite definitions in internal/config/custom_format_test.go
  Depends: T014, T015
  Acceptance: PLAN.md WP5 unit coverage; SPEC.md FR-001..FR-007 mapped

- [x] T017 Execute F-007 integration regression suite in test/integration/custom_format_inline_test.go
  Depends: T016
  Acceptance: PLAN.md WP5 integration coverage; SPEC.md SC-001..SC-004

---

## Phase 7 — Polish & Cross-Cutting

- [x] T018 [P] Capture F-007 behavior evidence and decision trace in target/specs/custom-format-config/BUILD.md
  Depends: T017
  Acceptance: PLAN.md validation approach lines 149-155; BUILD evidence expectation

- [x] T019 [P] Update task execution traceability checkboxes in target/specs/custom-format-config/TASKS.md
  Depends: T017
  Acceptance: AGENTS.md rule for TASKS.md as execution checklist; no unresolved dependencies

---

## Implementation Strategy

- **Recommended start**: T001 -> T006 (scope, constraints, and contracts first)
- **Critical path**: T001 -> T002 -> T003 -> T004 -> T005 -> T006 -> T007 -> T010 -> T013 -> T014 -> T016 -> T017 -> T018
- **Risk areas**: T013 (FR-010 decision blocker), T014 (error compatibility semantics), T010 (licenseText/copyright compatibility nuance)
- **Handoff checklist**:
  - [x] All P1 story tasks complete
  - [x] Validation phase tasks pass
  - [x] BUILD.md evidence artifact started
  - [x] No unresolved NEEDS CLARIFICATION items in task list

## Task-Decomposer Summary

- **Total task count**: 19
- **Per-phase breakdown**: Setup 3, Foundation 3, Story phases 9, Validation 2, Polish/Cross-cutting 2
- **Per-story task counts**: US1 (4), US2 (2), US3 (3)
- **Parallel opportunities**: 8 tasks marked `[P]` (T008, T009, T011, T012, T014, T015, T018, T019)
- **MVP scope**: P1 story coverage via T007-T010
- **Format validation**: Every task uses checklist format with ID and exact file path, and includes `Depends:` and `Acceptance:` lines
- **Unresolved items**:
  - None for F-007 bounded slice.
