# Implementation Tasks: Diagnostics, Error Handling, and Exit Codes

**Date**: 2026-06-22 | **Plan**: [PLAN.md](./PLAN.md) | **Spec**: [SPEC.md](./SPEC.md)
**Scope**: Reduced — F-007 only for `target/specs/diagnostics-errors-exit-codes/`; preserve exact stderr/stdout/exit-code/debug diagnostics and CLI integration tests; no implementation code, `BUILD.md`, `PLAN.md`, or `SPEC.md` is produced by this decomposition.

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 42    |
| Phase 1 (Setup)            | 4     |
| Phase 2 (Foundation)       | 7     |
| Story phases               | 23    |
| Polish / cross-cutting     | 3     |
| Parallel opportunities     | 19    |
| MVP scope (P1 stories)     | 8     |

## Dependency Graph

```text
Setup:        T001 -> T002 [P], T003 [P], T004 [P]
Foundation:   T002 -> T005 -> T006 -> T007 -> T009
              T006 -> T008 [P]
              T003 -> T010 [P] -> T011 [P]
US1 P1:       T009 -> T012 -> T013 -> T014 -> T015
              T013 -> T016 [P], T017 [P], T018 [P], T019 [P]
US2 P2:       T007,T014 -> T020 -> T021 -> T022 -> T023
              T021,T022 -> T024 [P], T025 [P] -> T026
US3 P3:       T007,T010,T014 -> T027 -> T028 -> T029
              T007 -> T030 -> T031 [P]
              T010 -> T032 [P] -> T033 -> T034 [P]
Validation:   T019,T026,T029,T031,T034 -> T035 -> T036 [P], T037 [P], T038 -> T039
Polish:       T039 -> T040 [P], T041 [P] -> T042
```

## Parallelism Guidance

- Tasks marked `[P]` may execute concurrently with sibling tasks in the same phase.
- Tasks with explicit `Depends:` fields must wait for all listed predecessors.
- Builder may reorder within a phase when dependencies allow, but must not skip phases.

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

## Context and Hooks

- Context7 used: `/golang/go` for `package main`, `internal/` package boundaries, standard `testing`, `go test`, and `os/exec` stdout/stderr/status capture guidance.
- Extension hooks: none; `.specify/extensions.yml` was not found at the repository root.
- Setup scripts: none found matching repository setup or agent-context update script patterns.
- Optional supporting artifacts: no `research.md`, `data-model.md`, `quickstart.md`, `contracts/`, or `/memory/constitution.md` found for this use case.

---

## Phase 1 — Setup

- [x] T001 Verify the Go module boundary for F-007 implementation in `target/go.mod`
  Depends: none
  Acceptance: PLAN.md Technical Context lines 38-46 and ARCHITECTURE.md Structure Decision lines 30-53; implementation remains under `target/` and F-007 adds no new third-party diagnostics dependency.

- [x] T002 [P] Create or verify diagnostics package scaffolding in `target/internal/diagnostics/doc.go`
  Depends: T001
  Acceptance: PLAN.md WP-001 lines 88-101 and ARCHITECTURE.md boundary line 66; exact diagnostic strings and exit decisions are centralized in `internal/diagnostics`.

- [x] T003 [P] Create or verify debug namespace package scaffolding in `target/internal/debuglog/doc.go`
  Depends: T001
  Acceptance: PLAN.md WP-006 lines 163-176 and ARCHITECTURE.md lines 43 and 182; debug behavior is namespace-gated and does not emulate uncaptured byte formatting.

- [x] T004 [P] Add source-derived F-007 fixture notes in `target/internal/diagnostics/testdata/f007-output-contract.md`
  Depends: T001
  Acceptance: SPEC.md Output Contract lines 105-178 and PLAN.md WP-007 lines 178-191; fixture notes come only from committed artifacts and do not run the live Node source.

---

## Phase 2 — Foundation

- [x] T005 Define F-007 exit and diagnostic result types in `target/internal/diagnostics/types.go`
  Depends: T002
  Acceptance: SPEC.md Key Entities lines 196-201 and PLAN.md WP-001 lines 96-99; model exit `0`, exit `1`, and continue/no-direct-exit without calling `os.Exit` outside `cmd/license-checker`.

- [x] T006 Define exact F-007 diagnostic literals and interpolation helpers in `target/internal/diagnostics/messages.go`
  Depends: T005
  Acceptance: SPEC.md FR-002 through FR-015 and Output Contract lines 109-122; preserve `delimeters`, `can not`, `Found error`, `Exiting.`, quotes, capitalization, and punctuation exactly.

- [x] T007 Implement stderr-only writer helpers in `target/internal/diagnostics/writer.go`
  Depends: T006
  Acceptance: SPEC.md Whitespace Contract lines 174-178 and PLAN.md WP-001 lines 94-101; helpers write one line to the intended writer and never duplicate diagnostics to stdout.

- [x] T008 [P] Add byte-exact diagnostic helper tests in `target/internal/diagnostics/messages_test.go`
  Depends: T006
  Acceptance: PLAN.md WP-001 Validation lines 99-101 and SPEC.md SC-002; tests compare exact source-derived bytes for all literal and interpolation helpers including newline policy.

- [x] T009 Define preflight result and ordering contract in `target/internal/diagnostics/preflight.go`
  Depends: T007
  Acceptance: SPEC.md Output Contract field order lines 153-160 and PLAN.md WP-002 lines 103-116; help, version, conflict, comma warning, and continue outcomes are represented explicitly.

- [x] T010 [P] Implement DEBUG namespace matcher in `target/internal/debuglog/debuglog.go`
  Depends: T003
  Acceptance: SPEC.md FR-014 and PLAN.md WP-006 lines 169-172; `DEBUG=license-checker*` enables `license-checker:log` and `license-checker:error` without changing exit decisions.

- [x] T011 [P] Add DEBUG namespace matcher tests in `target/internal/debuglog/debuglog_test.go`
  Depends: T010
  Acceptance: PLAN.md WP-006 Validation line 176 and ARCHITECTURE.md Conformance Matrix line 230; tests assert namespace enablement only and do not pin uncaptured timestamp/color/prefix bytes.

---

## Phase 3 — Story: CLI preflight diagnostics are predictable (Priority: P1)

- [x] T012 [US1] Implement source help body emission helper in `target/internal/diagnostics/help.go`
  Depends: T009
  Acceptance: SPEC.md Acceptance Scenario 1 lines 37-40 and help body lines 124-151; first line is `license-checker@25.0.1`, listed body lines are in source order, stderr is used, and exit code is `0`.

- [x] T013 [US1] Implement help/version/conflict/comma preflight logic in `target/internal/diagnostics/preflight.go`
  Depends: T012
  Acceptance: SPEC.md FR-001 through FR-005 and Algorithm Fidelity lines 184-188; help short-circuits first, version second, mutual exclusion third, comma warning fourth, then scan continues.

- [x] T014 [US1] Wire diagnostics preflight before scanning in `target/internal/app/run.go`
  Depends: T013
  Acceptance: ARCHITECTURE.md Composition Root lines 95-106 and PLAN.md WP-002 lines 109-116; fatal preflight exits before custom format loading, scanning, policy, or rendering.

- [x] T015 [US1] Verify process exit remains owned by main in `target/cmd/license-checker/main.go`
  Depends: T014
  Acceptance: PLAN.md Decision Log line 26 and ARCHITECTURE.md Entry point lines 89-94; `main.go` delegates to `app.Run` and calls `os.Exit(code)` with no diagnostic business logic.

- [x] T016 [P] [US1] Add table-driven preflight unit tests in `target/internal/diagnostics/preflight_test.go`
  Depends: T013
  Acceptance: SPEC.md Acceptance Scenarios 1-5 lines 39-43 and Edge Case line 79; tests assert ordering, stderr bytes, stdout absence via writer separation, and skip of comma warnings after mutual-exclusion fatal error.

- [x] T017 [P] [US1] Add process-level help and version CLI tests in `target/cmd/license-checker/preflight_cli_test.go`
  Depends: T015
  Acceptance: SPEC.md FR-001, FR-002, SC-001, and Output Contract lines 111-112; tests use Go `os/exec` capture for stderr, stdout, and exit status.

- [x] T018 [P] [US1] Add fatal preflight no-scan app test in `target/internal/app/preflight_test.go`
  Depends: T014
  Acceptance: SPEC.md Acceptance Scenario 3 line 41 and PLAN.md WP-002 Validation line 116; scanner spy proves simultaneous `--failOn` and `--onlyAllow` exits `1` before scanning or rendering.

- [x] T019 [P] [US1] Add comma-warning continuation app tests in `target/internal/app/preflight_test.go`
  Depends: T014
  Acceptance: SPEC.md Acceptance Scenarios 4-5 lines 42-43 and FR-004-FR-005; `--failOn MIT,ISC` and `--onlyAllow MIT,ISC` write exact warning to stderr and continue without direct exit change.

---

## Phase 4 — Story: License policy failures terminate with source-compatible diagnostics (Priority: P2)

- [x] T020 [US2] Implement semicolon policy token parsing for F-007 in `target/internal/policy/fail_policy.go`
  Depends: T007, T014
  Acceptance: SPEC.md Edge Case line 80 and Algorithm Fidelity line 190; non-empty trimmed semicolon tokens are used and empty trimmed tokens emit no diagnostic.

- [x] T021 [US2] Implement `--failOn` exact-match fatal diagnostic in `target/internal/policy/fail_policy.go`
  Depends: T020
  Acceptance: SPEC.md FR-006-FR-007 and Acceptance Scenarios lines 57-58; first restricted package whose `licenses` exactly equals a token writes `Found license defined by the --failOn flag: "<license>". Exiting.` to stderr and returns exit `1` immediately.

- [x] T022 [US2] Implement `--onlyAllow` substring violation diagnostic in `target/internal/policy/fail_policy.go`
  Depends: T021
  Acceptance: SPEC.md FR-008-FR-009 and Acceptance Scenario line 59; no allowed token contained in the license writes `Package "<item>" is licensed under "<license>" which is not permitted by the --onlyAllow flag. Exiting.` to stderr and returns exit `1`.

- [x] T023 [US2] Integrate fatal policy diagnostics into policy application in `target/internal/policy/apply.go`
  Depends: T022
  Acceptance: ARCHITECTURE.md Composition Root step 7 line 103 and PLAN.md WP-003 lines 124-129; diagnostics run after F-005 restrictions and before renderer callback output.

- [x] T024 [P] [US2] Add ordered policy failure unit tests in `target/internal/policy/fail_policy_test.go`
  Depends: T021, T022
  Acceptance: PLAN.md WP-003 Validation line 131 and SPEC.md SC-003; tests cover exact match, semicolon multi-token, empty-token ignore, first-match immediate exit, onlyAllow substring acceptance, violation interpolation, stdout empty, stderr exact, and status `1`.

- [x] T025 [P] [US2] Add app-level policy exit ordering tests in `target/internal/app/policy_diagnostics_test.go`
  Depends: T023
  Acceptance: SPEC.md Output Contract field order line 159 and PLAN.md WP-003 line 127; policy-failure diagnostics occur before any renderer output for that run.

- [x] T026 [US2] Add process-level policy failure CLI tests in `target/cmd/license-checker/policy_cli_test.go`
  Depends: T024, T025
  Acceptance: SPEC.md SC-001 through SC-003 and ARCHITECTURE.md Conformance Matrix lines 218-219; tests assert stderr, stdout, and exit status for `--failOn MIT`, semicolon multi-token failOn, and onlyAllow violation fixtures.

---

## Phase 5 — Story: Non-fatal diagnostics preserve channels without inventing output (Priority: P3)

- [x] T027 [US3] Implement callback-equivalent init error handling in `target/internal/app/run.go`
  Depends: T007, T010, T014
  Acceptance: SPEC.md FR-010 and Algorithm Fidelity line 193; scanner/init errors write `Found error` to stderr before the error object and do not introduce a new explicit non-zero exit solely for that callback error.

- [x] T028 [US3] Implement compatible error-object rendering boundary in `target/internal/diagnostics/errors.go`
  Depends: T027
  Acceptance: SPEC.md FR-011 and Assumption line 219; second error line remains the declared Node-compatible target formatting and is not replaced with JSON, YAML, or a new structured error code.

- [x] T029 [US3] Add callback error app tests in `target/internal/app/init_error_test.go`
  Depends: T028
  Acceptance: SPEC.md Acceptance Scenario 1 lines 73-74 and PLAN.md WP-004 Validation line 146; tests assert `Found error` precedes compatible error rendering, stdout has no diagnostic copy, and status is not newly forced to `1` solely for init error.

- [x] T030 [US3] Implement missing `--files` license-file warning in `target/internal/render/files.go`
  Depends: T007
  Acceptance: SPEC.md FR-012 and Acceptance Scenario 2 line 74; each package lacking usable `licenseFile` writes `no license file found for: <moduleName>` to stderr and continues with remaining records.

- [x] T031 [P] [US3] Add missing license-file warning tests in `target/internal/render/files_test.go`
  Depends: T030
  Acceptance: PLAN.md WP-005 Validation line 161 and SPEC.md SC-004; tests assert one warning per missing license file, exact interpolation/order, continued copies for later records, stdout unaffected, and no direct exit-code change from the warning.

- [x] T032 [P] [US3] Implement namespaced debug logger API in `target/internal/debuglog/debuglog.go`
  Depends: T010
  Acceptance: SPEC.md Acceptance Scenario 3 line 75 and PLAN.md WP-006 lines 169-172; exposes `license-checker:log` and `license-checker:error` calls while leaving byte formatting compatible/open.

- [x] T033 [US3] Wire scan-start and init-error debug calls in `target/internal/app/run.go`
  Depends: T032
  Acceptance: SPEC.md FR-014 and ARCHITECTURE.md flag propagation line 150; enabled debug calls do not alter functional stdout/stderr diagnostics or exit-code decisions.

- [x] T034 [P] [US3] Add debug behavior tests in `target/internal/app/debug_test.go`
  Depends: T033
  Acceptance: PLAN.md WP-006 Validation line 176 and SPEC.md Output Contract lines 121-122; tests assert namespace calls occur under `DEBUG=license-checker*` and normal F-007 stderr/stdout/status assertions remain unchanged.

---

## Phase 6 — Validation

- [x] T035 Add full F-007 process conformance matrix tests in `target/cmd/license-checker/f007_cli_test.go`
  Depends: T017, T019, T026, T029, T031, T034
  Acceptance: PLAN.md WP-007 lines 184-191 and SPEC.md SC-001 through SC-005; tests capture stdout, stderr, and process status for help, version, mutual exclusion, comma warnings, failOn, onlyAllow, missing-file warning, callback prefix, and debug namespace scenarios.

- [x] T036 [P] Add exact diagnostic golden assertions in `target/internal/diagnostics/testdata/f007-output-contract.txt`
  Depends: T035
  Acceptance: SPEC.md Output Contract lines 109-122 and Whitespace Contract lines 174-178; exact strings are byte-compared where captured and documented compatible exceptions remain only help trailing whitespace, debug formatting, and callback error object formatting.

- [x] T037 [P] Add reusable CLI capture helpers for stdout/stderr/status in `target/internal/testsupport/cli_capture_test.go`
  Depends: T035
  Acceptance: Context7 `/golang/go` `os/exec` guidance and ARCHITECTURE.md Testing line 24; helpers preserve separate stdout/stderr buffers and expose non-zero exit status without losing stderr bytes.

- [x] T038 Run all Go tests for the target module from `target/`
  Depends: T035, T036, T037
  Acceptance: PLAN.md Validation Approach line 209 and ARCHITECTURE.md Implementation Order line 196; `go test ./...` passes from `target/` with F-007 tests included and no live Node source execution.

- [x] T039 Audit F-007 conformance limitations in `target/internal/diagnostics/errors.go`
  Depends: T038
  Acceptance: ARCHITECTURE.md Conformance Matrix lines 228-231 and SPEC.md Assumptions lines 219-220; code comments or tests keep compatible/open rows explicit and do not claim byte identity for uncaptured help trailing bytes, debug formatting, or Node error-object formatting.

---

## Phase 7 — Polish & Cross-Cutting

- [x] T040 [P] Remove duplicate F-007 diagnostic literals outside `target/internal/diagnostics/messages.go`
  Depends: T039
  Acceptance: PLAN.md WP-001 Goal lines 88-91 and ARCHITECTURE.md boundary line 66; exact strings are not copied into `app`, `cli`, `policy`, or `render` except via diagnostics helpers/tests.

- [x] T041 [P] Verify stderr/stdout channel separation across F-007 tests in `target/cmd/license-checker/f007_cli_test.go`
  Depends: T039
  Acceptance: SPEC.md FR-015 and Whitespace Contract line 177; no stderr diagnostics are mirrored to stdout and normal stdout output is not polluted by warnings/errors.

- [x] T042 Finalize F-007 handoff checklist in `target/specs/diagnostics-errors-exit-codes/TASKS.md`
  Depends: T040, T041
  Acceptance: All task lines retain strict checklist format, every task has a `Depends:` and `Acceptance:` note, no unresolved clarification markers remain, and Builder can update checkboxes during implementation.

---

## Implementation Strategy

- **Recommended start**: T001 through T004, then centralize diagnostics in T005 through T009 before wiring any story behavior.
- **Critical path**: T001 → T002 → T005 → T006 → T007 → T009 → T013 → T014 → T020 → T021 → T022 → T023 → T026 → T035 → T038 → T042.
- **Risk areas**: help trailing whitespace, callback error-object formatting, debug dependency formatting, and policy ordering after F-005 restrictions; keep these aligned with compatible/open conformance rows instead of inventing output.
- **Handoff checklist**:
  - [x] All P1 story tasks complete: T012-T019
  - [x] Validation phase tasks pass: T035-T039
  - [x] BUILD.md evidence artifact started by Builder only
  - [x] No unresolved clarification items in task list

## Decomposition Summary

- **Total task count**: 42 tasks — Setup 4, Foundation 7, US1 8, US2 7, US3 8, Validation 5, Polish 3.
- **Per-story task counts**: US1/P1 has 8 tasks (T012-T019); US2/P2 has 7 tasks (T020-T026); US3/P3 has 8 tasks (T027-T034).
- **Parallel opportunities**: 19 tasks marked `[P]` across setup, foundation, story test lanes, validation helpers, and polish audits.
- **MVP scope**: P1 tasks T012-T019 deliver help/version/mutual-exclusion/comma-warning behavior with process-level CLI tests.
- **Format validation**: Every task follows the checklist format, uses sequential IDs, includes an exact target path or module path, and has explicit `Depends:` and `Acceptance:` notes.
- **Unresolved items**: none; open compatibility limitations are represented as acceptance constraints, not blockers.
