# Implementation Tasks: F-006 Output Rendering and Side-Effect Outputs

**Date**: 2026-06-22 | **Plan**: [PLAN.md](./PLAN.md) | **Spec**: [SPEC.md](./SPEC.md)
**Scope**: Reduced — F-006 only; consume already-scanned, already-filtered ordered records from F-003/F-004/F-005 and do not implement scanning, classification, filtering, F-001 preflight, F-007 exit-code policy, or live Node differential capture.

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 41    |
| Phase 1 (Setup)            | 4     |
| Phase 2 (Foundation)       | 6     |
| Story phases               | 26    |
| Polish / cross-cutting     | 2     |
| Parallel opportunities     | 13    |
| MVP scope (P1 stories)     | 13    |

## Dependency Graph

```text
Setup:      T001 ─┬─ T002 ─┬─ T005 ─┬─ T006 [P] ─┐
                  ├─ T003 [P]       ├─ T007 [P] ─┼─ T009 ─ T010
                  └─ T004 [P]       └─ T008 [P] ─┘

US1/P1:     T010 ─ T011 ─┬─ T012 [P]
                         ├─ T013 [P]
                         ├─ T014
                         └─ T015

US2/P1:     T010 ─┬─ T016 ─ T017 [P]
                  ├─ T018 ─ T019 [P]
                  ├─ T020 ─ T021 [P]
                  └─ T022 ─ T023

US3/P2:     T018 ─ T024 ─ T025 [P]
            T020 ─ T026 ─ T027 [P]
            T011 ─ T028 ─ T029
            T024/T026/T028 ─ T030 ─ T031

US4/P2:     T023/T029 ─ T032 ─┬─ T033
                               ├─ T034 ─ T035 [P]
                               └─ T036

Validation: T012/T017/T019/T021/T023/T025/T027/T029/T031/T035/T036 ─ T037 ─ T038 ─ T039
Polish:     T039 ─ T040 ─ T041
```

## Parallelism Guidance

- Tasks marked `[P]` may execute concurrently with sibling tasks in the same phase.
- Tasks with explicit `Depends:` fields must wait for all listed predecessors.
- Builder may reorder within a phase when dependencies allow, but must not skip phases.
- Context7 note: `/golang/go` was consulted for internal package boundaries, `encoding/json` indentation support, `os.MkdirAll`/file APIs, and standard `testing` table-driven workflows; this supports the stdlib-first decomposition below.
- Extension hooks: `.specify/extensions.yml` was not present, so no executable `hooks.before_tasks` or `hooks.after_tasks` entries were surfaced.
- Setup scripts: no repository setup or agent-context update scripts were found; continue without setup-script output.

## Task Format Rules

Every task MUST use this exact checklist format:

`- [x] T001 [P] [US1] Description with exact file path`

Format rules:

- Start every task with `- [x]`
- Use sequential task IDs in execution order: `T001`, `T002`, `T003`, ...
- Add `[P]` only when the task is parallelizable
- Add `[US1]`, `[US2]`, `[US3]`, etc. only for story-phase tasks
- Include the exact file path directly in the task description
- Keep dependency details in a `Depends:` note directly below the task

---

## Phase 1 — Setup

- [x] T001 Create Go module scaffold for the target implementation in `target/go.mod`
  Depends: none
  Acceptance: ARCHITECTURE.md lines 18-24 and PLAN.md Technical Context lines 37-45 require a Go module rooted under `target/` using stdlib-first testing and rendering packages.

- [x] T002 Create the thin command entry point shell for later app wiring in `target/cmd/license-checker/main.go`
  Depends: T001
  Acceptance: ARCHITECTURE.md lines 57-65 and 87-106 require `cmd/license-checker` to delegate business logic to `internal/app` and keep rendering inside `internal/render`.

- [x] T003 [P] Create filesystem adapter contracts for renderer side effects in `target/internal/files/fs.go`
  Depends: T001
  Acceptance: ARCHITECTURE.md lines 76-83 and PLAN.md WP-009 require `ReadFile`, `WriteFile`, `MkdirAll`, and `Stat` abstractions compatible with `--out` and `--files` tests.

- [x] T004 [P] Create chalk-compatible color helper contracts in `target/internal/color/chalk.go`
  Depends: T001
  Acceptance: SPEC.md FR-019 through FR-020 and ARCHITECTURE.md dependency decision lines 180-183 require no-color byte baseline with explicit forced-color role support.

---

## Phase 2 — Foundation

- [x] T005 Define ordered package-map and record field types in `target/internal/ordered/object.go`
  Depends: T001
  Acceptance: SPEC.md FR-021 and ARCHITECTURE.md lines 22 and 240 require preserving JavaScript insertion-order semantics and avoiding plain Go map iteration at renderer boundaries.

- [x] T006 [P] Implement custom ordered JSON marshaling primitives in `target/internal/ordered/json.go`
  Depends: T005
  Acceptance: SPEC.md FR-005 and ARCHITECTURE.md package decision line 164 require two-space JSON output over deliberate top-level and per-record order.

- [x] T007 [P] Define F-006 renderer option and side-effect option structs in `target/internal/render/types.go`
  Depends: T005
  Acceptance: PLAN.md WP-001 requires options for JSON, CSV, Markdown, Summary, Color, Out, Files, CustomFormat, and CSVComponentPrefix without reshaping record content.

- [x] T008 [P] Create reusable byte-comparison and ordered-record test helpers in `target/internal/render/test_helpers_test.go`
  Depends: T005
  Acceptance: SPEC.md SC-001 through SC-006 and PLAN.md Validation Approach require byte-level tests from committed snippets and source-derived contracts.

- [x] T009 Implement renderer selection facade signatures in `target/internal/render/render.go`
  Depends: T006, T007
  Acceptance: SPEC.md FR-001 and PLAN.md WP-008 require a central formatter path with precedence JSON, CSV, Markdown, summary, then tree.

- [x] T010 Create app-facing render integration boundary in `target/internal/app/render.go`
  Depends: T002, T007, T009
  Acceptance: ARCHITECTURE.md ordered stage sequence lines 104-105 and PLAN.md WP-011 require `internal/app` to call render formatting and emission without embedding renderer business logic.

---

## Phase 3 — Story: Render default tree output (Priority: P1)

- [x] T011 [US1] Implement treeify-compatible traversal and glyph rendering in `target/internal/render/tree.go`
  Depends: T005, T007
  Acceptance: SPEC.md FR-003, Output Contract lines 150-170, and PLAN.md WP-002 require `treeify.asTree(obj, true)`-compatible glyphs, indentation, and insertion-order traversal.

- [x] T012 [P] [US1] Add README tree and guessed-license byte golden tests in `target/internal/render/tree_test.go`
  Depends: T011, T008
  Acceptance: SPEC.md User Story 1 scenarios 1-2 and SC-001 require exact lines for `cli@0.4.3`, repository/license children, and `MIT*` preservation.

- [x] T013 [P] [US1] Add tree golden fixture snippets in `target/internal/render/testdata/tree_readme.golden`
  Depends: T008
  Acceptance: SPEC.md Output Contract lines 150-170 requires preserving committed default tree and guessed-license snippets as byte baselines.

- [x] T014 [US1] Implement programmatic print behavior for tree output in `target/internal/render/print.go`
  Depends: T011
  Acceptance: SPEC.md FR-004 and whitespace contract lines 242-248 require `print(sorted)` to write `asTree(sorted)` through a line-emitting stdout operation.

- [x] T015 [US1] Implement no-color default tree path through the renderer facade in `target/internal/render/render.go`
  Depends: T009, T011, T014
  Acceptance: SPEC.md FR-001 and User Story 1 Independent Test require default rendering when no JSON/CSV/Markdown/summary/files/out structured-output flag is selected.

---

## Phase 4 — Story: Emit machine-readable JSON, CSV, and Markdown (Priority: P1)

- [x] T016 [US2] Implement ordered two-space JSON formatter with CLI newline rule in `target/internal/render/json.go`
  Depends: T006, T007
  Acceptance: SPEC.md FR-005 and whitespace contract lines 244-245 require ordered `JSON.stringify(..., null, 2)`-equivalent output plus exactly one formatter newline.

- [x] T017 [P] [US2] Add JSON field-order and newline tests in `target/internal/render/json_test.go`
  Depends: T016, T008
  Acceptance: SPEC.md User Story 2 scenario 3 and SC-002 require top-level key order, record field insertion order, and stdout-vs-file newline differences.

- [x] T018 [US2] Implement default CSV header and row rendering in `target/internal/render/csv.go`
  Depends: T007
  Acceptance: SPEC.md FR-006, FR-007, FR-009, and CSV default snippet lines 172-177 require exact header order, key order, empty missing fields, `\n` joins, and no trailing newline.

- [x] T019 [P] [US2] Add default CSV byte tests for headers, rows, missing fields, and no escaping in `target/internal/render/csv_default_test.go`
  Depends: T018, T008
  Acceptance: SPEC.md User Story 2 scenario 1, Edge Cases lines 103-104, and SC-003 require exact CSV snippets and source-compatible non-RFC escaping behavior.

- [x] T020 [US2] Implement default Markdown row rendering in `target/internal/render/markdown.go`
  Depends: T007
  Acceptance: SPEC.md FR-010 and Markdown default snippet lines 193-197 require exact `[<key>](<repository>) - <licenses>` rows joined with `\n` and no trailing newline in `asMarkDown`.

- [x] T021 [P] [US2] Add default Markdown byte tests in `target/internal/render/markdown_default_test.go`
  Depends: T020, T008
  Acceptance: SPEC.md User Story 2 scenario 2 and SC-001 require byte-identical default Markdown output for `abbrev@1.0.9`.

- [x] T022 [US2] Wire JSON, CSV, and Markdown precedence into formatter selection in `target/internal/render/render.go`
  Depends: T016, T018, T020
  Acceptance: SPEC.md FR-001 and Algorithm Fidelity lines 256-260 require JSON before CSV before Markdown before summary/tree.

- [x] T023 [US2] Add structured-mode precedence and raw formatter newline matrix tests in `target/internal/render/format_test.go`
  Depends: T022, T017, T019, T021
  Acceptance: SPEC.md SC-002 and Output Contract lines 242-248 require separate raw formatter, stdout, and `--out` newline assertions for JSON, CSV, and Markdown.

---

## Phase 5 — Story: Render summaries and custom-format tables (Priority: P2)

- [x] T024 [US3] Implement custom-format CSV headers and rows with optional component prefix in `target/internal/render/csv.go`
  Depends: T018
  Acceptance: SPEC.md FR-008, FR-027, FR-028, and Output Contract lines 179-191 require custom key insertion order and optional leading component column/value.

- [x] T025 [P] [US3] Add custom CSV and component-prefix byte tests in `target/internal/render/csv_custom_test.go`
  Depends: T024, T008
  Acceptance: SPEC.md User Story 3 scenarios 1-2 and SC-003 require exact headers and rows for `name`, `description`, `pewpew`, and `main-module`.

- [x] T026 [US3] Implement custom-format Markdown package and detail rows in `target/internal/render/markdown.go`
  Depends: T020
  Acceptance: SPEC.md FR-011, FR-027, and Markdown custom snippet lines 199-203 require exact package row plus four-space-indented custom key/value rows.

- [x] T027 [P] [US3] Add custom Markdown byte tests in `target/internal/render/markdown_custom_test.go`
  Depends: T026, T008
  Acceptance: SPEC.md User Story 3 scenario 3 requires exact ` - **[abbrev@1.0.9](...)**` package row and ordered custom detail rows.

- [x] T028 [US3] Implement exact-license summary aggregation and descending count rendering in `target/internal/render/summary.go`
  Depends: T011
  Acceptance: SPEC.md FR-013 through FR-014 and PLAN.md WP-006 require exact `licenses` string counts, descending numeric sort, tree rendering, and no invented lexical tie sort.

- [x] T029 [US3] Add summary aggregation and equal-count limitation tests in `target/internal/render/summary_test.go`
  Depends: T028, T008
  Acceptance: SPEC.md User Story 3 scenario 4, Edge Case line 106, and SC-005 require count sorting while avoiding deterministic equal-count assertions.

- [x] T030 [US3] Wire summary and custom-format modes into formatter selection in `target/internal/render/render.go`
  Depends: T024, T026, T028
  Acceptance: SPEC.md FR-001 and Algorithm Fidelity lines 257-260 require custom-format-aware CSV/Markdown behavior and summary fallback after structured modes.

- [x] T031 [US3] Add custom-format and summary precedence tests in `target/internal/render/format_custom_summary_test.go`
  Depends: T030, T025, T027, T029
  Acceptance: SPEC.md SC-003 through SC-005 and PLAN.md WP-008 require locked precedence, custom-format order, summary tree output, and newline distinctions.

---

## Phase 6 — Story: Write output files and copy license files (Priority: P2)

- [x] T032 [US4] Implement side-effect precedence and emit API in `target/internal/render/emit.go`
  Depends: T023, T031, T003
  Acceptance: SPEC.md FR-002 and Algorithm Fidelity lines 261-263 require `--files` to bypass formatted output, else `--out`, else stdout line emission.

- [x] T033 [US4] Implement `--out` recursive parent creation and raw UTF-8 writes in `target/internal/render/emit.go`
  Depends: T032
  Acceptance: SPEC.md FR-015 and User Story 4 scenario 1 require mkdirp-compatible parent creation and raw formatted string writes without stdout-added newline.

- [x] T034 [US4] Implement `asFiles` copy behavior and missing-license warning routing in `target/internal/render/files.go`
  Depends: T032
  Acceptance: SPEC.md FR-016 through FR-018 and User Story 4 scenarios 2-3 require ordered iteration, `<moduleName>-LICENSE.txt`, recursive directories, copied contents, and exact warning text.

- [x] T035 [P] [US4] Add `--files` copy, nested-name directory, and missing-warning tests in `target/internal/render/files_test.go`
  Depends: T034, T008
  Acceptance: SPEC.md SC-004 and Output Contract lines 205-215 require `foo-LICENSE.txt`, copied content, warning continuation, and exact `no license file found for: <moduleName>` text.

- [x] T036 [US4] Add `--out` versus stdout side-effect newline tests in `target/internal/render/emit_test.go`
  Depends: T033, T008
  Acceptance: SPEC.md SC-002 and whitespace contract lines 242-248 require file content and stdout line-emitting behavior to differ exactly by mode.

---

## Phase 7 — Validation

- [x] T037 Add forced-color compatibility rendering tests in `target/internal/render/color_test.go`
  Depends: T004, T015, T023, T031
  Acceptance: SPEC.md FR-019 through FR-020, Edge Case line 107, and ARCHITECTURE.md conformance row line 227 require no-color byte identity plus forced-color blue/dim/green key segmentation and bold-red sentinel preservation.

- [x] T038 Add app-level render handoff tests with fake records and fake filesystem in `target/internal/app/render_test.go`
  Depends: T010, T032, T033, T034
  Acceptance: ARCHITECTURE.md stage sequence lines 104-105 and PLAN.md WP-011 require `internal/app` to call `render.Format` then `render.Emit` without F-001/F-007 behavior leaking into `render`.

- [x] T039 Add full F-006 golden parity matrix tests in `target/internal/render/golden_test.go`
  Depends: T012, T017, T019, T021, T025, T027, T029, T035, T036, T037
  Acceptance: SPEC.md SC-001 through SC-007 and ARCHITECTURE.md Acceptance reference lines 152-154 require byte tests for all committed/source-derived snippets and explicit limitation coverage for absent live fixtures.

---

## Phase 8 — Polish & Cross-Cutting

- [x] T040 Add renderer package documentation and guardrail comments in `target/internal/render/doc.go`
  Depends: T039
  Acceptance: PLAN.md Guardrails lines 203-209 require documenting ordered data, no plain maps at renderer boundaries, no `encoding/csv` behavior changes, and no invented live fixtures.

- [x] T041 Verify the target test command and keep evidence for Builder-owned reporting in `target/internal/render/render_test.go`
  Depends: T040
  Acceptance: PLAN.md Validation Approach line 201 and ARCHITECTURE.md lines 24 and 195-196 require `go test ./...` from `target/`; Builder must record evidence in its own `BUILD.md`, not in this task artifact.

- [x] T042 Correct verifier-found golden parity gaps by removing terminal-newline trimming and adding fixture-backed byte tests for every identical conformance row in `target/internal/render/render_test.go`, `target/internal/app/app_test.go`, and `target/internal/diagnostics/diagnostics_test.go`
  Depends: T039, T041
  Acceptance: ARCHITECTURE.md Conformance Matrix lines 204-229 and Builder verifier feedback require committed fixtures under `target/testdata/golden/`, exact `[]byte` comparisons without whitespace normalization, and documented contract-derived status for rows lacking live captures.

---

## Implementation Strategy

- **Recommended start**: T001 through T010 establish the `target/` module, ordered data contracts, renderer options, and app-facing render boundary.
- **Critical path**: T001 → T005 → T006/T007 → T009 → T011/T016/T018/T020/T028 → T022/T030 → T032 → T039 → T041.
- **Risk areas**: T011 tree glyph fidelity, T016 ordered JSON without Go map iteration, T018/T024 source-compatible non-escaping CSV, T028 summary equal-count instability, T034 path/name behavior for `--files`, and T037 color compatibility.
- **MVP scope**: P1 stories are T011-T023: default tree plus JSON, CSV, Markdown formatters and their raw newline/precedence tests.
- **Per-story counts**: US1 has 5 tasks; US2 has 8 tasks; US3 has 8 tasks; US4 has 5 tasks.
- **Parallel opportunities**: 13 tasks are marked `[P]`: T003, T004, T006, T007, T008, T012, T013, T017, T019, T021, T025, T027, T035.
- **Format validation**: Every task line follows the checklist format, uses sequential IDs T001-T041, includes an exact target file path, and has `Depends:` plus `Acceptance:` notes.
- **Unresolved items**: None. Parity limitations from SPEC.md lines 299-303 are preserved as validation scope limits, not implementation blockers.
- **Handoff checklist**:
  - [x] All P1 story tasks complete
  - [x] Validation phase tasks pass
  - [x] BUILD.md evidence artifact started by Builder
  - [x] No unresolved clarification items in task list
