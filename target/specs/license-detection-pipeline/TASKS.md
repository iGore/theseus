# Implementation Tasks: License detection and normalization pipeline (F-004)

**Date**: 2026-05-01 | **Plan**: `target/specs/license-detection-pipeline/PLAN.md` | **Spec**: `target/specs/license-detection-pipeline/SPEC.md`
**Scope**: Reduced — strictly F-004 only (license detection/normalization pipeline), Go target only

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 21    |
| Phase 1 (Setup)            | 2     |
| Phase 2 (Foundation)       | 4     |
| Story phases               | 3     |
| Polish / cross-cutting     | 3     |
| Parallel opportunities     | 10    |
| MVP scope (P1 stories)     | 4     |

## Dependency Graph

T001 -> T003 -> T005 -> T006 -> T009 -> T010 -> T011 -> T015 -> T016 -> T017 -> T018

T002 -> T004 -> T008 -> T012 -> T014 -> T015

Parallel lanes:
- [P] T004 with T003 after T001/T002
- [P] T007 and [P] T008 after T005/T006 and T004 respectively
- [P] T013 and [P] T014 after T011/T012
- [P] T016 and [P] T017 after T015

## Parallelism Guidance

- Tasks marked `[P]` may execute concurrently with sibling tasks in the same phase.
- Tasks with explicit `Depends:` fields must wait for all listed predecessors.
- Builder may reorder within a phase when dependencies allow, but must not skip phases.

## Task Format Rules

Every task MUST use this exact checklist format:

`- [ ] T001 [P] [US1] Description with exact file path`

---

## Phase 1 — Setup

- [x] T001 Confirm F-004-only pipeline boundaries and stage contract notes in target/specs/license-detection-pipeline/PLAN.md
  Depends: none
  Acceptance: PLAN.md scope note + SPEC.md C-001/C-004 referenced and unchanged

- [x] T002 [P] Create license-resolution fixture index for F-004 scenarios in test/integration/testdata/license_resolution/README.md
  Depends: none
  Acceptance: SPEC.md AC-001..AC-005 scenario matrix mapped to fixture names

---

## Phase 2 — Foundation

- [x] T003 Define ModuleLicenseInput and NormalizedLicenseValue fields required by F-004 in internal/core/model.go
  Depends: T001
  Acceptance: SPEC.md Key Entities + FR-001/FR-008 represented in model contracts

- [x] T004 [P] Define ordered LicenseCandidateSet contract and matcher inputs in internal/license/files.go
  Depends: T002
  Acceptance: SPEC.md FR-006 and LicenseCandidateSet entity are directly represented

- [x] T005 Add F-004 resolver entrypoint contract between service and license stage in internal/core/service.go
  Depends: T003
  Acceptance: PLAN.md WP1/WP5 traceability and SPEC.md FR-008 boundary preserved

- [x] T006 Implement metadata-first resolver flow orchestration point in internal/license/detect.go
  Depends: T003, T005
  Acceptance: SPEC.md FR-001/FR-002 precedence starts with metadata path

---

## Phase 3 — Story: Resolve a license value for each module (Priority: P1)

- [x] T007 [P] [US1] Implement metadata parsing for license string/object/array forms in internal/license/detect.go
  Depends: T006
  Acceptance: SPEC.md FR-001 + User Story 1 Acceptance Scenario 2

- [x] T008 [P] [US1] Implement SPDX-first normalization path before heuristics in internal/license/normalize.go
  Depends: T004, T006
  Acceptance: SPEC.md FR-003 + AC-002 direct SPDX normalization behavior

- [x] T009 [US1] Wire per-module single-value resolution output assignment in internal/core/service.go
  Depends: T007, T008
  Acceptance: SPEC.md SC-001 one deterministic normalized value per module

- [x] T010 [US1] Add table-driven metadata resolution tests in internal/license/detect_test.go
  Depends: T007, T008, T009
  Acceptance: SPEC.md User Story 1 independent test + AC-001 metadata-first proof

---

## Phase 4 — Story: Fall back deterministically when metadata is incomplete (Priority: P2)

- [x] T011 [US2] Implement metadata -> README -> file fallback chain in internal/license/detect.go
  Depends: T009
  Acceptance: SPEC.md FR-002 + User Story 2 Acceptance Scenario 1

- [x] T012 [US2] Implement deterministic file candidate precedence LICENSE*/LICENCE*/COPYING/README in internal/license/files.go
  Depends: T004, T011
  Acceptance: SPEC.md FR-006 + AC-004 deterministic selection order

- [x] T013 [P] [US2] Add README-no-match continuation behavior in internal/license/detect.go
  Depends: T011
  Acceptance: SPEC.md Edge Case (README with no signal continues to file fallback)

- [x] T014 [P] [US2] Add fallback precedence and repeat-run determinism tests in test/integration/license_resolution_fallback_test.go
  Depends: T012, T013
  Acceptance: SPEC.md AC-001/AC-004 + SC-004 verified across repeated runs

---

## Phase 5 — Story: Preserve heuristic and custom-reference normalization semantics (Priority: P3)

- [x] T015 [US3] Implement inferred marker (`*`) and `Custom: <value>` normalization in internal/license/normalize.go
  Depends: T008, T011
  Acceptance: SPEC.md FR-004/FR-005 + User Story 3 Acceptance Scenarios 1-2

- [x] T016 [P] [US3] Add heuristic-marker and custom-reference unit coverage in internal/license/normalize_test.go
  Depends: T015
  Acceptance: SPEC.md AC-002/AC-003 + SC-003

- [x] T017 [P] [US3] Add compatibility lock tests for legacy parser expectations including OD placeholders in test/integration/license_resolution_compat_test.go
  Depends: T015
  Acceptance: SPEC.md FR-007/AC-005 with explicit `NEEDS CLARIFICATION` markers for OD-001 and OD-002

---

## Phase 6 — Validation

- [x] T018 Execute full F-004 validation matrix and document deterministic pass evidence in target/specs/license-detection-pipeline/BUILD.md
  Depends: T010, T014, T016, T017
  Acceptance: SPEC.md SC-001..SC-004 and PLAN.md Validation Approach all satisfied

---

## Phase 7 — Polish & Cross-Cutting

- [x] T019 [P] Document unresolved OD-001 tie-break behavior checkpoint in target/specs/license-detection-pipeline/TASKS.md
  Depends: T017
  Acceptance: SPEC.md Open Decisions OD-001 remains explicit as `NEEDS CLARIFICATION`

- [x] T020 [P] Document unresolved OD-002 empty-signal fallback checkpoint in target/specs/license-detection-pipeline/TASKS.md
  Depends: T017
  Acceptance: SPEC.md Open Decisions OD-002 remains explicit as `NEEDS CLARIFICATION`

- [x] T021 Final F-004 traceability sweep across task checklist references in target/specs/license-detection-pipeline/TASKS.md
  Depends: T018, T019, T020
  Acceptance: Every task has ID + exact file path + Depends + Acceptance mapped to PLAN.md/SPEC.md

---

## Implementation Strategy

- **Recommended start**: T001, then T003/T005/T006 to establish F-004 contract and resolver seam.
- **Critical path**: T001 -> T003 -> T005 -> T006 -> T007/T008 -> T009 -> T011 -> T012 -> T015 -> T018 -> T021.
- **Risk areas**: T007/T011 (metadata shape + fallback correctness), T015 (normalization parity), T017 (legacy compatibility lock), OD-001/OD-002 unresolved semantics.
- **Handoff checklist**:
  - [x] All P1 story tasks complete
  - [x] Validation phase tasks pass
  - [x] BUILD.md evidence artifact started
  - [ ] No unresolved NEEDS CLARIFICATION items in task list

## Decomposer Summary

- **Total task count**: 21 (Setup 2, Foundation 4, Story tasks 11, Validation 1, Polish/Cross-cutting 3).
- **Per-story counts**: US1 (4 tasks), US2 (4 tasks), US3 (3 tasks).
- **Parallel opportunities**: 10 tasks marked `[P]` (T002, T004, T007, T008, T013, T014, T016, T017, T019, T020) across Setup, Foundation, Story, and Polish phases.
- **MVP scope (P1)**: T007-T010 provide first viable F-004 slice for metadata-driven deterministic normalized license output.
- **Format validation**: Every task uses checklist format and includes ID, exact file path, `Depends`, and `Acceptance` details.
- **Unresolved items**:
  - T017 / T019 — `NEEDS CLARIFICATION` for OD-001 (`license` vs `licenses` conflict tie-break).
  - T017 / T020 — `NEEDS CLARIFICATION` for OD-002 (no recognizable signal fallback output).
