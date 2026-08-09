# Implementation Tasks: F-004 License Detection, Classification, and License File Precedence

**Date**: 2026-06-22 | **Plan**: [PLAN.md](./PLAN.md) | **Spec**: [SPEC.md](./SPEC.md)
**Scope**: Reduced — generate implementation tasks for F-004 only; exclude F-003 scanning, F-005 policy, F-006 rendering, live source execution, and non-F-004 CLI behavior.

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 35    |
| Phase 1 (Setup)            | 4     |
| Phase 2 (Foundation)       | 6     |
| Story phases               | 21    |
| Polish / cross-cutting     | 2     |
| Parallel opportunities     | 16    |
| MVP scope (P1 stories)     | 8     |

## Dependency Graph

```text
Setup:      T001 -> {T002 [P], T003 [P], T004 [P]}
Foundation: T002,T003,T004 -> T005 -> T006 -> {T007 [P], T008 [P], T009 [P]} -> T010
US1/P1:     T010 -> T011 -> T012 -> T013 -> {T014 [P], T015 [P], T016 [P], T017 [P], T018 [P]}
US2/P2:     T010 -> T019 -> T020 -> {T021 [P], T022 [P], T023 [P]}
US3/P3:     T018,T023 -> T024 -> T025 -> T026 -> T027 -> T028 -> T029 -> {T030 [P], T031 [P]}
Validation: T018,T023,T031 -> T032 -> T033
Polish:     T033 -> T034 -> T035
```

## Parallelism Guidance

- Tasks marked `[P]` may execute concurrently with sibling tasks in the same phase.
- Tasks with explicit `Depends:` fields must wait for all listed predecessors.
- Builder may reorder within a phase when dependencies allow, but must not skip phases.
- Context7 was consulted for Go `internal/` visibility and `go test` workflow; retain F-004 code under `target/internal/license` and SPDX calls behind `target/internal/spdxcompat`.
- No `.specify/extensions.yml`, `/memory/constitution.md`, setup scripts, or agent-context update scripts were found; no task hooks are executable for this slice.

## Task Format Rules

Every task MUST use this exact checklist format:

`- [x] T001 [P] [US1] Description with exact file path`

Format rules:

- Start every task with `- [x]`
- Use sequential task IDs in execution order: `T001`, `T002`, `T003`, ...
- Add `[P]` only when the task is parallelizable
- Add `[US1]`, `[US2]`, `[US3]`, etc. only for story-phase tasks
- Include the exact file path directly in the task description
- Keep dependency details in a `Depends:` note directly below the task when needed

---

## Phase 1 — Setup

- [x] T001 Confirm or create Go module dependency entry for `github.com/github/go-spdx/v2 v2.7.0` in `target/go.mod`
  Depends: none
  Acceptance: PLAN.md WP-001 and ARCHITECTURE.md Package Decisions require SPDX validation behind `internal/spdxcompat`; do not replace the architecture-selected dependency without approval.

- [x] T002 [P] Create F-004 package scaffold in `target/internal/license/doc.go`
  Depends: T001
  Acceptance: ARCHITECTURE.md lines 61-63 and PLAN.md Project Structure place classifier, license-file precedence, and record application logic in `internal/license`.

- [x] T003 [P] Create SPDX adapter package scaffold in `target/internal/spdxcompat/doc.go`
  Depends: T001
  Acceptance: PLAN.md WP-001 and ARCHITECTURE.md lines 23 and 190 require all SPDX parser use through `internal/spdxcompat`.

- [x] T004 [P] Create F-004 fixture directory marker in `target/internal/license/testdata/f004/README.md`
  Depends: T001
  Acceptance: SPEC.md FR-011 and PLAN.md Validation Approach prohibit live source execution and require Analyzer-recorded/source test expectations only.

---

## Phase 2 — Foundation

- [x] T005 Define classifier-facing SPDX validation API in `target/internal/spdxcompat/adapter.go`
  Depends: T003
  Acceptance: SPEC.md FR-002 requires success/failure behavior equivalent to `spdx-expression-parse ^3.0.0` while returning raw classifier input unchanged on success.

- [x] T006 Implement `github.com/github/go-spdx/v2/spdxexp` validation wrapper in `target/internal/spdxcompat/adapter.go`
  Depends: T005
  Acceptance: PLAN.md Open Questions note exact API details as `NEEDS CLARIFICATION`; resolve by validating the mandated module API without changing the selected dependency.

- [x] T007 [P] Add SPDX adapter worked-example validation tests in `target/internal/spdxcompat/adapter_test.go`
  Depends: T006
  Acceptance: SPEC.md Algorithm Fidelity examples `MIT`, `LGPL-2.0`, `Apache-2.0`, `BSD-2-Clause`, `(GPL-2.0+ WITH Bison-exception-2.2)`, `LGPL-2.0 OR (ISC AND BSD-3-Clause+)`, and `Apache-2.0 OR ISC OR MIT` must validate successfully.

- [x] T008 [P] Define F-004 nullable license value representation in `target/internal/license/value.go`
  Depends: T002
  Acceptance: SPEC.md FR-008 and PLAN.md WP-002 require preserving exact literals plus a null-equivalent fallback distinct from the string `Undefined`.

- [x] T009 [P] Define deterministic filesystem reader interface for F-004 tests in `target/internal/license/filesystem.go`
  Depends: T002
  Acceptance: PLAN.md WP-007 requires injected filesystem reads so package/readme/license-file tests remain deterministic and do not run the live source system.

- [x] T010 Define F-004 detector options and public apply seam in `target/internal/license/detector.go`
  Depends: T005, T008, T009
  Acceptance: ARCHITECTURE.md stage sequence invokes `license.Apply` during `records.Flatten`; keep F-003/F-006 ownership boundaries explicit.

---

## Phase 3 — Story: Classify license strings and text exactly (Priority: P1)

- [x] T011 [US1] Implement SPDX pass-through first rule in `target/internal/license/classifier.go`
  Depends: T010
  Acceptance: SPEC.md US1 scenarios 1-2 and FR-002 require valid SPDX identifiers/expressions to return the original raw input unchanged before regex rules run.

- [x] T012 [US1] Implement falsey-input and first-newline preprocessing in `target/internal/license/classifier.go`
  Depends: T011
  Acceptance: SPEC.md Edge Cases require falsey input to return literal `Undefined` and truthy input to remove only the first newline before regex matching.

- [x] T013 [US1] Implement ordered classifier regex table rules 4 through 14 in `target/internal/license/classifier.go`
  Depends: T012
  Acceptance: SPEC.md Algorithm Fidelity requires exact ordered results `ISC*`, `MIT*`, `BSD*`, `BSD-Source-Code*`, `WTFPL*`, `Apache*`, and `CC0-1.0*`, including the shadowed BSD source-code ordering.

- [x] T014 [P] [US1] Implement GPL and LGPL version capture rules in `target/internal/license/classifier.go`
  Depends: T013
  Acceptance: SPEC.md Algorithm Fidelity rules 15-16 require case-insensitive version capture and `.0` suffix when captured version length is one.

- [x] T015 [P] [US1] Implement Public Domain and Custom capture fallback rules in `target/internal/license/classifier.go`
  Depends: T013
  Acceptance: SPEC.md US1 scenario 5 and Algorithm Fidelity rules 17-18 require `Public Domain` and `Custom: <capture>` for URL or `SEE LICENSE IN (.*)` matches.

- [x] T016 [P] [US1] Implement no-match null fallback in `target/internal/license/classifier.go`
  Depends: T013
  Acceptance: SPEC.md US1 scenario 6 and FR-001 require fallback/no-match to preserve a null value rather than `UNKNOWN` or `Undefined`.

- [x] T017 [P] [US1] Add complete ordered classifier table tests in `target/internal/license/classifier_test.go`
  Depends: T014, T015, T016
  Acceptance: SPEC.md SC-001 requires byte-identical strings or null for committed `tests/license.js:10-181` examples across SPDX pass-through, guessed values, `Custom:`, `Public Domain`, `Undefined`, and fallback null.

- [x] T018 [P] [US1] Add classifier edge-case tests for rule order and preprocessing in `target/internal/license/classifier_test.go`
  Depends: T014, T015, T016
  Acceptance: SPEC.md Edge Cases require first-newline-only replacement, SPDX-before-regex pass-through, and BSD source-code shadowing to be preserved exactly.

---

## Phase 4 — Story: Detect license files by source precedence (Priority: P2)

- [x] T019 [US2] Implement basename uppercasing and extension-insensitive normalization in `target/internal/license/license_files.go`
  Depends: T010
  Acceptance: SPEC.md FR-005 and Edge Cases require matching on uppercased basename with extension ignored.

- [x] T020 [US2] Implement exact license-file precedence pattern table in `target/internal/license/license_files.go`
  Depends: T019
  Acceptance: SPEC.md License-file precedence table requires ordered patterns `LICENSE`, `LICENSE-*`, `LICENCE`, `LICENCE-*`, `COPYING`, `README`, with no fallback selection.

- [x] T021 [P] [US2] Implement first-filename-per-pattern selection in `target/internal/license/license_files.go`
  Depends: T020
  Acceptance: SPEC.md US2 scenario 2 and Edge Cases require only the first filename matching each precedence pattern to be pushed.

- [x] T022 [P] [US2] Add precedence and extension/case tests in `target/internal/license/license_files_test.go`
  Depends: T020
  Acceptance: SPEC.md SC-002 requires matching committed `tests/license-files-test.js:10-87` expectations for empty/no-match, extension-insensitive, case-insensitive, and `LICENSE.md` before `COPYING`/`README.txt`.

- [x] T023 [P] [US2] Add first-`LICENSE-*` and no-fallback tests in `target/internal/license/license_files_test.go`
  Depends: T021
  Acceptance: SPEC.md US2 scenarios 2 and 4 require selecting only the first `LICENSE-*` candidate and returning an empty selection when no precedence pattern matches.

---

## Phase 5 — Story: Apply package/readme/file license precedence into records (Priority: P3)

- [x] T024 [US3] Define F-003-compatible record and node consumer interfaces in `target/internal/license/record.go`
  Depends: T018, T023
  Acceptance: PLAN.md Open Questions note exact F-003 types as `NEEDS CLARIFICATION`; use consumer-side interfaces without owning F-003 scanning or renderer formatting.

- [x] T025 [US3] Implement package metadata license normalization in `target/internal/license/metadata.go`
  Depends: T024
  Acceptance: SPEC.md FR-003 requires `json.license || json.licenses`, string/object classification, and array mapping through the classifier.

- [x] T026 [US3] Implement README-derived license fill when package metadata is absent in `target/internal/license/metadata.go`
  Depends: T025
  Acceptance: SPEC.md FR-004 and US3 scenario 1 require README classification to fill only when package license fields do not provide a license value.

- [x] T027 [US3] Implement selected license-file reading and record field insertion in `target/internal/license/apply.go`
  Depends: T026
  Acceptance: SPEC.md FR-006, FR-009, and Output Contract require conditional `licenseFile` and `licenseText` fields while leaving relative-path rendering ownership outside F-004.

- [x] T028 [US3] Implement NOTICE file detection and conditional field insertion in `target/internal/license/apply.go`
  Depends: T027
  Acceptance: SPEC.md FR-010 and US3 scenario 3 require conditional `noticeFile` after license/copyright fields when a NOTICE file is found.

- [x] T029 [US3] Implement license-file override conditions in `target/internal/license/apply.go`
  Depends: T027
  Acceptance: SPEC.md FR-007 and US3 scenario 2 require selected license-file text to overwrite absent, containing-`UNKNOWN`, or `Custom:` licenses without any CLI flag.

- [x] T030 [P] [US3] Add package metadata and README precedence tests in `target/internal/license/apply_test.go`
  Depends: T025, T026
  Acceptance: PLAN.md WP-003 requires constructed tests proving package license precedence over README, README fills missing license, arrays are mapped, and falsey/missing metadata preserves expected sentinel flow.

- [x] T031 [P] [US3] Add license-file field and override flow tests in `target/internal/license/apply_test.go`
  Depends: T027, T028, T029
  Acceptance: SPEC.md SC-003 and Output Contract require missing/`UNKNOWN`/`Custom:` override to exact `MIT*`, valid package license non-override, and conditional field order for `licenseFile`, `licenseText`, `copyright`, and `noticeFile`.

---

## Phase 6 — Validation

- [x] T032 Validate F-004 unit and integration test coverage with `go test ./...` from `target/go.mod`
  Depends: T007, T017, T018, T022, T023, T030, T031
  Acceptance: PLAN.md Validation Approach and SPEC.md SC-001 through SC-005 require all F-004 classifier, SPDX adapter, license-file detector, and record-flow tests to pass without live source execution.

- [x] T033 Verify F-004 values remain renderer-ready through ordered record assertions in `target/internal/license/apply_test.go`
  Depends: T032
  Acceptance: SPEC.md SC-004 requires downstream renderers receiving this slice's values to reproduce exact `UNKNOWN`, `MIT*`, and `ISC` values in recorded output contexts.

---

## Phase 7 — Polish & Cross-Cutting

- [x] T034 Add source-rule trace comments for classifier and file precedence in `target/internal/license/classifier.go`
  Depends: T033
  Acceptance: PLAN.md WP-002 requires explicit rules and comments mapping to SPEC.md Algorithm Fidelity table order so future cleanup does not reorder behavior.

- [x] T035 Add package boundary documentation for F-003/F-006 ownership in `target/internal/license/doc.go`
  Depends: T034
  Acceptance: PLAN.md WP-007 and ARCHITECTURE.md stage sequence require no CLI flags, no policy behavior, no renderer formatting, and mutation only of F-004-owned record fields.

---

## Implementation Strategy

- **Recommended start**: T001 through T010, especially the SPDX adapter seam before classifier code.
- **Critical path**: T001 → T005 → T006 → T010 → T011 → T012 → T013 → T014/T015/T016 → T017/T018 → T024 → T029 → T031 → T032 → T033.
- **Risk areas**: T006 SPDX API validation, T012 first-newline-only preprocessing, T013/T014/T015 ordered classifier fidelity, T021 file precedence first-match behavior, T024 F-003 interface seam, and T029 license-file override conditions.
- **Handoff checklist**:
  - [x] All P1 story tasks complete: T011-T018
  - [x] Validation phase tasks pass: T032-T033
  - [x] BUILD.md evidence artifact started by Builder, not Task Decomposer
  - [x] No NEEDS CLARIFICATION item remains unaddressed by the bounded adapter/interface tasks

## Output Summary

- **Total task count**: 35 tasks — Phase 1: 4, Phase 2: 6, US1: 8, US2: 5, US3: 8, Validation: 2, Polish/Cross-cutting: 2.
- **Per-story task counts**: US1 P1 classifier: 8 tasks (T011-T018); US2 P2 license-file precedence: 5 tasks (T019-T023); US3 P3 record precedence/application: 8 tasks (T024-T031).
- **Parallel opportunities**: 16 `[P]` tasks across setup (T002-T004), foundation (T007-T009), US1 tests/rules (T014-T018), US2 selection/tests (T021-T023), and US3 tests (T030-T031).
- **MVP scope**: P1 viable slice is T001-T018 plus T032 for classifier/SPDX validation; full F-004 requires US2/US3 before Builder handoff completion.
- **Format validation**: Every task line uses checklist format, sequential task ID, exact file path, explicit `Depends`, and `Acceptance` detail.
- **Unresolved items**: `NEEDS CLARIFICATION` is carried from PLAN.md into T006 for exact `go-spdx` API verification and T024 for future F-003 concrete type wiring; both are bounded by architecture-approved adapter/interface seams and do not authorize scope changes.
