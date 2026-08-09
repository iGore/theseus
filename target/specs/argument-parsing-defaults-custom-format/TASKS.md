# Implementation Tasks: Argument parsing, defaults, and custom format loading

**Date**: 2026-06-22 | **Plan**: [PLAN.md](./PLAN.md) | **Spec**: [SPEC.md](./SPEC.md)
**Scope**: Reduced — F-002 only. Do not implement CLI process lifecycle/help/version exits, dependency scanning internals, license detection, policy filtering, final renderers, `BUILD.md`, `PLAN.md`, or `SPEC.md`.

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 37    |
| Phase 1 (Setup)            | 4     |
| Phase 2 (Foundation)       | 5     |
| Story phases               | 21    |
| Polish / cross-cutting     | 3     |
| Parallel opportunities     | 9     |
| MVP scope (P1 stories)     | 9     |

## Context Notes

- Context7 was checked for `/golang/go`; task breakdown follows Go `package main`, `internal/` package boundaries, standard `testing`, and `go test` conventions.
- No `.specify/extensions.yml` file was present, so there are no executable `hooks.before_tasks` or `hooks.after_tasks` entries to surface.
- No repository setup or agent-context update scripts were found under the repository root, so no setup script output was captured.
- No optional `research.md`, `data-model.md`, `contracts/`, `quickstart.md`, or `/memory/constitution.md` artifacts were present for this use case.

## Dependency Graph

```text
Setup:
  T001 -> T002 -> T003
              \-> T004 [P]

Foundation:
  T002 -> T005 -> T006
  T005 -> T007 [P]
  T005 -> T008 [P]
  T005 -> T009

US1 P1 parser/default MVP:
  T005,T006,T007 -> T010 -> T011 -> T012 -> T013 -> T014 -> T015
  T010 -> T016 [P]
  T012,T013,T014,T015,T016 -> T017
  T010,T017 -> T018

US2 P2 scanner-boundary values:
  T013,T015 -> T019 -> T020
  T019 -> T021 [P]
  T020,T021 -> T022

US3 P3 custom format:
  T005,T009 -> T023 -> T024 -> T025
  T023 -> T026 [P]
  T025,T026 -> T027
  T023 -> T028 [P]
  T027,T028 -> T029
  T023,T029 -> T030

Validation:
  T017,T018 -> T031
  T022 -> T032 [P]
  T030 -> T033 [P]
  T031,T032,T033 -> T034

Polish:
  T034 -> T035 -> T036 -> T037
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

## Phase 1 — Setup

- [x] T001 Create Go module scaffold for the F-002 target in `target/go.mod`
  Depends: none
  Acceptance: PLAN.md Project Structure lines 65-84 and ARCHITECTURE.md lines 18-24, 30-53 require a Go module rooted under `target/` using private `internal/` packages.

- [x] T002 Create thin command entry-point placeholder in `target/cmd/license-checker/main.go`
  Depends: T001
  Acceptance: ARCHITECTURE.md lines 57 and 91-93 require `package main` to delegate to `internal/app` without business logic; F-002 must not implement help/version lifecycle.

- [x] T003 Create application composition skeleton for F-002 boundaries in `target/internal/app/app.go`
  Depends: T002
  Acceptance: ARCHITECTURE.md lines 58, 95-100 and PLAN.md WP-006 require composition points for parse, custom-format load, and scanner-boundary handoff without scanner internals.

- [x] T004 [P] Create package documentation for F-002 scope in `target/internal/cli/doc.go`
  Depends: T002
  Acceptance: SPEC.md Scope lines 11-19 and Out of scope lines 20-26 are reflected so Builder keeps parser/defaults bounded to F-002.

---

## Phase 2 — Foundation

- [x] T005 Define complete F-002 option model and direct-depth sentinel in `target/internal/cli/options.go`
  Depends: T001
  Acceptance: SPEC.md FR-001 through FR-002 and PLAN.md WP-001 require every known option field, no exposed `argv`, color explicitness, `CustomFormat`, and `0`/unbounded scanner-depth representation.

- [x] T006 Define known option type metadata and short aliases in `target/internal/cli/known.go`
  Depends: T005
  Acceptance: SPEC.md FR-002 through FR-003 and Algorithm Fidelity step 1 require all declared `nopt` option names plus `-v -> --version` and `-h -> --help`.

- [x] T007 [P] Define runtime color support abstraction in `target/internal/color/color.go`
  Depends: T005
  Acceptance: SPEC.md FR-007 and ARCHITECTURE.md flag propagation row `--color` require `chalk.supportsColor` equivalent to be injectable for defaults tests.

- [x] T008 [P] Define ordered key-value primitives for observable field order in `target/internal/ordered/object.go`
  Depends: T005
  Acceptance: SPEC.md FR-017 and ARCHITECTURE.md lines 22 and 240 require ordered data instead of Go map iteration for custom-format field order.

- [x] T009 Define custom-format ordered representation and error-as-value result type in `target/internal/customformat/format.go`
  Depends: T005
  Acceptance: SPEC.md Key Entities `CustomFormat`, FR-014 through FR-021, and PLAN.md WP-004/WP-005 require parsed objects, ordered keys, exact-`false` suppression, and returned errors as values.

---

## Phase 3 — Story: Parse CLI flags with source defaults (Priority: P1)

- [x] T010 [US1] Implement nopt-compatible raw parser entry point in `target/internal/cli/parser.go`
  Depends: T005, T006
  Acceptance: SPEC.md US1 scenarios 4 and 6, FR-001 through FR-006, and Algorithm Fidelity steps 1-2 require single-command parsing, known option typing, aliases, cooked-argument presence tracking, and cleaned results without `argv`.

- [x] T011 [US1] Implement boolean negation and cooked presence tracking in `target/internal/cli/presence.go`
  Depends: T010
  Acceptance: SPEC.md FR-006 and US1 scenario 6 require `has("color")` to return true for both `--color` and `--no-color` source-equivalent cooked forms.

- [x] T012 [US1] Implement source-order defaults normalization in `target/internal/cli/defaults.go`
  Depends: T010, T007
  Acceptance: SPEC.md US1 scenarios 1-3 and 5 plus Algorithm Fidelity steps 3-9 require color default/forced false, cwd `start`, boolean `relativeLicensePath`, and direct depth normalization.

- [x] T013 [US1] Implement path and string option coercion for declared flags in `target/internal/cli/parser.go`
  Depends: T010
  Acceptance: SPEC.md FR-002, FR-004, and Source dependencies lines 205-208 require source-compatible declared string/path option behavior for `out`, `customPath`, `files`, `start`, and string filters without guessing unknown edge behavior.

- [x] T014 [US1] Implement cleaned parse result metadata removal in `target/internal/cli/parser.go`
  Depends: T010
  Acceptance: SPEC.md FR-005 and Output Contract lines 159-165 require the returned options to omit `argv` metadata while preserving source-equivalent field names and values.

- [x] T015 [US1] Implement direct-depth sentinel helpers in `target/internal/cli/depth.go`
  Depends: T012
  Acceptance: SPEC.md FR-011, FR-013, and Assumptions line 236 require truthy `direct` to become numeric depth `0` and missing/falsey `direct` to become an unbounded sentinel at scanner boundaries.

- [x] T016 [P] [US1] Add table-driven parser coverage for all known options and aliases in `target/internal/cli/parser_test.go`
  Depends: T010, T011, T013, T014
  Acceptance: SPEC.md SC-001 requires tests for every known option name, both short aliases, `--color`/`--no-color` presence, and absence of `argv` metadata.

- [x] T017 [US1] Add table-driven defaults coverage for color/start/relative/direct in `target/internal/cli/defaults_test.go`
  Depends: T012, T015, T016
  Acceptance: SPEC.md SC-002 and US1 scenarios 1-5 require missing/falsey/truthy defaults, output-mode color forcing, cwd defaulting, `direct true -> 0`, and missing direct -> unbounded.

- [x] T018 [US1] Record nopt parity fixtures and known parser gaps in `target/internal/cli/conformance_test.go`
  Depends: T010, T017
  Acceptance: SPEC.md OQ-001/OQ-002 and SC-006 require fixture-backed known behavior and explicit `NEEDS CLARIFICATION` risk markers for malformed values, unknown flags, full path coercion, and unverified `--no-<flag>` edges without live source execution.

---

## Phase 4 — Story: Select scan root and direct dependency depth (Priority: P2)

- [x] T019 [US2] Define scanner boundary interface and options used by F-002 in `target/internal/npmgraph/options.go`
  Depends: T005, T015
  Acceptance: SPEC.md User Story 2 and Key Entities `ScannerOptions boundary` require `Start` and `Depth` handoff values while leaving dependency scanner implementation to F-003.

- [x] T020 [US2] Wire parsed start and direct depth into app scanner invocation in `target/internal/app/app.go`
  Depends: T003, T012, T015, T019
  Acceptance: SPEC.md US2 scenarios 1-4 and PLAN.md WP-006 require default/provided `start` and `depth: options.direct` to be passed to scanner options.

- [x] T021 [P] [US2] Add fake scanner boundary tests for start and direct depth in `target/internal/app/app_test.go`
  Depends: T019
  Acceptance: SPEC.md SC-005 and PLAN.md Validation Approach lines 158-159 require fake scanner tests for default/provided start, direct-present depth `0`, and direct-absent unbounded depth.

- [x] T022 [US2] Add app composition test proving no scanner internals are implemented in `target/internal/app/app_test.go`
  Depends: T020, T021
  Acceptance: SPEC.md Out of scope lines 22-23 and PLAN.md WP-006 require the app boundary to pass options to a fake scanner only, without package graph reconstruction in F-002.

---

## Phase 5 — Story: Load and apply custom format fields (Priority: P3)

- [x] T023 [US3] Implement parseJson-equivalent loader in `target/internal/customformat/load.go`
  Depends: T009
  Acceptance: SPEC.md US3 scenarios 1-2, FR-014 through FR-015, and Custom JSON loading rule table require non-string path errors, UTF-8 JSON parsing, and file/parse errors returned as values without throwing or invented diagnostics.

- [x] T024 [US3] Integrate customPath assignment before scanner construction in `target/internal/app/app.go`
  Depends: T020, T023
  Acceptance: SPEC.md FR-012 and US3 scenario 1 require `options.customFormat` to be assigned from `customPath` before scanner options are built.

- [x] T025 [US3] Implement custom-format include/default/suppression rules in `target/internal/customformat/apply.go`
  Depends: T023, T008
  Acceptance: SPEC.md FR-016 through FR-020 and Custom-format field ordered rule table require include unless exactly `false`, string override, missing-field default, exact-false suppression, and truthy non-string no-default branch.

- [x] T026 [P] [US3] Implement ordered JSON object decoding for custom format keys in `target/internal/ordered/json.go`
  Depends: T008, T023
  Acceptance: SPEC.md FR-017 and Output Contract lines 166-170 require source-equivalent key order rather than Go map iteration for custom CSV/Markdown consumers.

- [x] T027 [US3] Expose record-field application helper for downstream records in `target/internal/records/custom_format.go`
  Depends: T025, T026
  Acceptance: SPEC.md FR-016 through FR-021 and PLAN.md WP-005 require F-002 to provide the custom-format field contract consumed by F-003/F-006 without implementing final renderers.

- [x] T028 [P] [US3] Add custom JSON loader tests and fixtures in `target/internal/customformat/testdata/custom-format.json`
  Depends: T023
  Acceptance: SPEC.md SC-003 and Output Contract lines 145-158 require valid JSON containing `licenseModified: "no"`, invalid JSON, missing file, and non-string path returning Error values.

- [x] T029 [US3] Add custom-format field-order and default tests in `target/internal/customformat/apply_test.go`
  Depends: T025, T026, T028
  Acceptance: SPEC.md SC-004 and US3 scenarios 3-5 require key order `name`, `description`, `pewpew`, package string override, default insertion, exact-false suppression, and truthy non-string no-default behavior.

- [x] T030 [US3] Add downstream record contract tests for custom fields in `target/internal/records/custom_format_test.go`
  Depends: T027, T029
  Acceptance: SPEC.md SC-007 and Output Contract lines 131-143 require custom-format-controlled fields and order to support later byte-identical custom CSV snippets without implementing F-006 rendering now.

---

## Phase 6 — Validation

- [x] T031 Run parser/default unit test coverage gate in `target/internal/cli/parser_test.go`
  Depends: T017, T018
  Acceptance: SPEC.md SC-001, SC-002, and SC-006 are covered by table-driven parser/default tests and explicit nopt conformance risk cases.

- [x] T032 [P] Run scanner-boundary validation gate in `target/internal/app/app_test.go`
  Depends: T022
  Acceptance: SPEC.md US2 scenarios 1-4 and PLAN.md Validation Approach lines 158-159 pass with fake scanner assertions only.

- [x] T033 [P] Run custom-format validation gate in `target/internal/customformat/apply_test.go`
  Depends: T030
  Acceptance: SPEC.md SC-003, SC-004, and SC-007 pass for JSON loading, ordered custom fields, default values, suppression, and record contract fixtures.

- [x] T034 Run full Go test suite for the target module in `target/go.mod`
  Depends: T031, T032, T033
  Acceptance: PLAN.md Validation Approach line 160 requires `go test ./...` from `target/` once implementation exists; no live source-system execution is used.

---

## Phase 7 — Polish & Cross-Cutting

- [x] T035 Document F-002 conformance rows and unresolved parser risks in `target/specs/argument-parsing-defaults-custom-format/BUILD.md`
  Depends: T034
  Acceptance: SPEC.md OQ-001 through OQ-003 and SC-006 require `NEEDS CLARIFICATION` risks for nopt edge cases and parse-error downstream behavior to be recorded by Builder evidence; Task Decomposer must not create `BUILD.md`.

- [x] T036 Review F-002 package boundaries for out-of-scope leakage in `target/internal/app/app.go`
  Depends: T035
  Acceptance: SPEC.md Out of scope lines 20-26 and PLAN.md Scope line 6 require no help/version process lifecycle, dependency scanner internals, license detection, filtering, or final renderers in this slice.

- [x] T037 Final format and traceability check for F-002 tasks in `target/specs/argument-parsing-defaults-custom-format/TASKS.md`
  Depends: T036
  Acceptance: Every task line in this file follows the required checklist format, includes an exact target file path, has `Depends:` and `Acceptance:` notes, and maps to SPEC.md/PLAN.md/ARCHITECTURE.md without invented scope.

---

## Implementation Strategy

- **Recommended start**: T001 through T005 to establish the target module and shared `Options` shape before parser behavior.
- **Critical path**: T001 → T005 → T010 → T012 → T015 → T019 → T020 → T024 → T027 → T030 → T034.
- **Risk areas**: T010/T013/T018 for `nopt ^4.0.1` parity gaps; T023/T024 for error-as-value customPath behavior; T025/T026/T027 for ordered custom-format fields and Go map-order avoidance.
- **Handoff checklist**:
  - [x] All P1 story tasks complete
  - [x] Validation phase tasks pass
  - [x] BUILD.md evidence artifact started
  - [x] No unresolved NEEDS CLARIFICATION items in task list

## Summary Report

- **Total task count**: 37 tasks — Setup 4, Foundation 5, US1 9, US2 4, US3 8, Validation 4, Polish/Cross-cutting 3.
- **Per-story task counts**: US1 Parse CLI flags with source defaults = 9; US2 Select scan root and direct dependency depth = 4; US3 Load and apply custom format fields = 8.
- **Parallel opportunities**: 9 `[P]` tasks across setup/foundation, parser tests, scanner-boundary tests, custom-format ordered decoding/tests, and validation gates.
- **MVP scope**: P1 viable first slice is T010 through T018 after setup/foundation T001 through T009; it preserves parser, defaults, `nopt` metadata removal, alias, color, `start`, `relativeLicensePath`, and `direct` parity tasks.
- **Format validation**: All 37 task lines use `- [x] TNNN`, story tasks carry `[USN]`, each description includes an exact file path, and every task includes `Depends:` and `Acceptance:` details.
- **Unresolved items**: T018 and T035 carry `NEEDS CLARIFICATION` markers for unverified `nopt` edge behavior and parse-error downstream behavior from SPEC.md OQ-001 through OQ-003; Builder must document these rather than guess or run the live source system in-flow.
