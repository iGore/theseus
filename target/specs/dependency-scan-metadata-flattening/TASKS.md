# Implementation Tasks: F-003 Dependency Scan and Metadata Flattening

**Date**: 2026-06-22 | **Plan**: [`PLAN.md`](./PLAN.md) | **Spec**: [`SPEC.md`](./SPEC.md)
**Scope**: Full for assigned F-003 use-case folder only; excludes F-001/F-002 parsing, F-004 license classification/file extraction, F-005 policy filtering, F-006 rendering, and F-007 diagnostics except for F-003 data contracts consumed by those slices.

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 36    |
| Phase 1 (Setup)            | 4     |
| Phase 2 (Foundation)       | 6     |
| Story phases               | 3     |
| Polish / cross-cutting     | 2     |
| Parallel opportunities     | 18    |
| MVP scope (P1 stories)     | 8     |

## Context Notes

- `.specify/extensions.yml`: not present; no executable `hooks.before_tasks` or `hooks.after_tasks` entries to surface.
- Repository setup or agent-context update scripts: none found; no setup output was captured.
- `/memory/constitution.md`: not present; repository workflow constraints come from `AGENTS.md` and the mandatory artifacts.
- Context7: `/golang/go` was consulted for Go task-shaping conventions: private `internal/` packages, `package main` command entry point, standard `testing`/`go test`, and deterministic string sorting.

## Dependency Graph

```text
Setup:
  T001 -> {T002 [P], T003 [P], T004 [P]} -> Foundation

Foundation:
  T001 -> T005 -> {T006 [P], T007}
  T001 -> {T008, T009, T010 [P]}
  {T005,T007,T008,T009} -> Story phases

US1 P1 stable scan/map MVP:
  {T003,T004,T008,T009} -> T011
  T008 -> {T012 [P], T014 [P]}
  {T011,T012,T014} -> T013 -> T015 [P]
  {T005,T007,T009} -> T016 -> {T017 [P], T018}

US2 P2 metadata/order:
  T016 -> {T019, T021, T024}
  T019 -> T020 [P]
  T021 -> {T022 [P], T023 [P]}

US3 P3 scope/path/sentinels:
  T016 -> {T025, T027, T030}
  T025 -> T026 [P]
  T027 -> T028 [P]
  T011 -> T029 [P]
  T030 -> T031 [P]

Validation:
  {T013,T018,T019,T021,T025,T027,T030} -> T032
  T013 -> T033 [P]
  {T015,T017,T020,T022,T023,T026,T028,T029,T031,T032,T033} -> T034 [P]

Polish:
  {T032,T034} -> {T035, T036}
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
- Keep dependency details in a `Depends:` note directly below the task

---

## Phase 1 — Setup

- [x] T001 Create Go module skeleton for the target implementation in `target/go.mod`
  Depends: none
  Acceptance: `PLAN.md` Project Structure lines 66-82 and `ARCHITECTURE.md` Structure Decision lines 28-53 require a Go module rooted under `target/`.

- [x] T002 [P] Create thin command entry point placeholder in `target/cmd/license-checker/main.go`
  Depends: T001
  Acceptance: `ARCHITECTURE.md` Composition Root lines 87-94 requires `package main` to delegate business logic outside `main.go`.

- [x] T003 [P] Create app-level F-003 invocation seam in `target/internal/app/f003.go`
  Depends: T001
  Acceptance: `PLAN.md` WP-07 lines 160-169 and `ARCHITECTURE.md` stage sequence lines 95-106 require scanner/flatten orchestration after parsed options and before later slices.

- [x] T004 [P] Create non-behavioral debug log interface in `target/internal/debuglog/debuglog.go`
  Depends: T001
  Acceptance: `SPEC.md` Acceptance Scenario US1.1 and `PLAN.md` WP-03 lines 116-119 require a provider log option without F-003 owning debug byte formatting.

---

## Phase 2 — Foundation

- [x] T005 Implement ordered record object mutations in `target/internal/ordered/object.go`
  Depends: T001
  Acceptance: `SPEC.md` Output Contract lines 142-156 and `PLAN.md` WP-02 lines 99-108 require insert-if-new, overwrite-without-moving, deletion, and deterministic key iteration.

- [x] T006 [P] Add ordered object mutation tests in `target/internal/ordered/object_test.go`
  Depends: T005
  Acceptance: `SPEC.md` SC-003 and `ARCHITECTURE.md` Risks lines 238-241 require tests proving Go map ordering is not used for observable records.

- [x] T007 Implement ordered top-level package map in `target/internal/ordered/map.go`
  Depends: T005
  Acceptance: `SPEC.md` FR-019 and Output Contract line 166 require explicit lexical package-key order equivalent to `Object.keys(data).sort()`.

- [x] T008 Define `read-installed`-compatible node and scanner contracts in `target/internal/npmgraph/node.go`
  Depends: T001
  Acceptance: `SPEC.md` FR-002, FR-024 and Algorithm Fidelity lines 183-185 require node fields for F-003 plus room for downstream `license`, `licenses`, and `readme`.

- [x] T009 Define records flattener options and API in `target/internal/records/options.go`
  Depends: T001
  Acceptance: `SPEC.md` FR-001 and `PLAN.md` WP-01 lines 88-98 require F-003 options for production, development, unknown, custom format, start/base path, and included fields.

- [x] T010 [P] Add F-003 source-derived fixture documentation in `target/internal/records/testdata/README.md`
  Depends: T001
  Acceptance: `SPEC.md` Edge Cases lines 98-104 and `PLAN.md` Validation Approach lines 194-202 require committed examples/source-derived fixtures only, never live Node source execution.

---

## Phase 3 — Story: Scan an installed npm project into stable package records (Priority: P1)

- [x] T011 [US1] Implement F-003 provider option mapping in `target/internal/app/f003.go`
  Depends: T003, T004, T008, T009
  Acceptance: `SPEC.md` US1.1 and FR-002 through FR-004 require `read(options.start, { dev: true, log: debugLog, depth: options.direct }, cb)`-equivalent semantics, with `dev=false` for production or development and depth passed unchanged.

- [x] T012 [P] [US1] Implement npm package metadata loading in `target/internal/npmgraph/packagejson.go`
  Depends: T008
  Acceptance: `SPEC.md` FR-002 and Key Entities lines 206-211 require package node fields `name`, `version`, `private`, `repository`, `url`, `author`, `path`, `dependencies`, `extraneous`, and `root`.

- [x] T013 [US1] Implement filesystem scanner preserving `read-installed` tree compatibility in `target/internal/npmgraph/scanner.go`
  Depends: T011, T012, T014
  Acceptance: `SPEC.md` FR-024 and Algorithm Fidelity lines 183-185 require logical installed npm tree semantics, nested dependencies, duplicate/dedupe behavior, `extraneous`, `root`, and absolute `path` propagation.

- [x] T014 [P] [US1] Add installed-tree compatibility fixture root in `target/internal/npmgraph/testdata/read-installed-fixture/package.json`
  Depends: T008
  Acceptance: `PLAN.md` WP-08 lines 171-182 requires compatibility fixtures for nested dependencies, `extraneous`, `root`, absolute paths under `start`, and dedupe behavior.

- [x] T015 [P] [US1] Add scanner compatibility and provider-option tests in `target/internal/npmgraph/scanner_test.go`
  Depends: T013
  Acceptance: `SPEC.md` US1.1, SC-004, and SC-006 require tests for provider `dev`, `depth`, nested tree shape, duplicate identities, and path propagation.

- [x] T016 [US1] Implement flatten recursion, duplicate guard, and incomplete-node deletion in `target/internal/records/flatten.go`
  Depends: T005, T007, T009
  Acceptance: `SPEC.md` US1.3, US1.4, FR-005, FR-006, FR-007, FR-010, FR-017, and FR-018 require exact `name@version` keys, initial `licenses=UNKNOWN`, optional `private=true`, duplicate short-circuit, dependency recursion, and missing-name/version deletion.

- [x] T017 [P] [US1] Add flatten key, duplicate, circular, and incomplete-node tests in `target/internal/records/flatten_test.go`
  Depends: T016
  Acceptance: `SPEC.md` Edge Cases lines 98-104 and SC-006 require fixture coverage for duplicate/circular guard, first-record retention, and missing name/version deletion.

- [x] T018 [US1] Implement lexical sorted package-map finalization in `target/internal/records/sort.go`
  Depends: T016
  Acceptance: `SPEC.md` FR-019, Output Contract line 166, and `PLAN.md` WP-06 lines 148-158 require sorted top-level key order before F-005 filters.

---

## Phase 4 — Story: Preserve source metadata and repository normalization (Priority: P2)

- [x] T019 [US2] Implement ordered repository URL normalization in `target/internal/records/repository.go`
  Depends: T016
  Acceptance: `SPEC.md` US2.1, FR-011, and repository normalization table lines 170-181 require ordered replacements and fallback preservation.

- [x] T020 [P] [US2] Add repository normalization table tests in `target/internal/records/repository_test.go`
  Depends: T019
  Acceptance: `SPEC.md` SC-002 requires every repository normalization row, fallback, and `abbrev@1.0.9` value `https://github.com/isaacs/abbrev-js` to be asserted exactly.

- [x] T021 [US2] Implement metadata, author, custom-format, dependencyPath, and path field population in `target/internal/records/metadata.go`
  Depends: T016, T019
  Acceptance: `SPEC.md` US2.2, US2.3, FR-012 through FR-016, and field-order table lines 142-156 require source-order insertion and overwrite-without-moving behavior.

- [x] T022 [P] [US2] Add field insertion-order and author URL overwrite tests in `target/internal/records/metadata_order_test.go`
  Depends: T021
  Acceptance: `SPEC.md` SC-003 and Edge Case line 103 require ordered keys for `licenses`, `private`, `repository`, `url`, `publisher`, `email`, author `url`, `dependencyPath`, custom fields, and `path`.

- [x] T023 [P] [US2] Add custom-format default, exclusion, and overwrite tests in `target/internal/records/custom_format_test.go`
  Depends: T021
  Acceptance: `SPEC.md` US2.3 and FR-015 require source-order custom keys, `false` exclusions, string copy, absent/falsey defaults, and no default for non-string truthy values.

- [x] T024 [US2] Preserve F-004 append seam after F-003 fields in `target/internal/records/license_hook.go`
  Depends: T016, T021
  Acceptance: `SPEC.md` Output Contract lines 154-156 and `PLAN.md` WP-02 line 108 require later `licenseFile`, `licenseText`, `copyright`, and `noticeFile` fields to append after F-003-owned fields without reordering them.

---

## Phase 5 — Story: Honor scan-scope and path options (Priority: P3)

- [x] T025 [US3] Implement production and development pruning checks in `target/internal/records/prune.go`
  Depends: T016
  Acceptance: `SPEC.md` US3.1, US3.2, FR-008, FR-009, and pruning table lines 195-202 require production to omit `extraneous` nodes and development to omit non-`extraneous`, non-`root` nodes before insertion.

- [x] T026 [P] [US3] Add default, production, and development pruning tests in `target/internal/records/prune_test.go`
  Depends: T025
  Acceptance: `SPEC.md` SC-004 and SC-006 require scan-scope tests for default inclusion plus production/development prune outcomes.

- [x] T027 [US3] Implement `path` and `dependencyPath` option behavior in `target/internal/records/paths.go`
  Depends: T016, T021
  Acceptance: `SPEC.md` US3.4, US3.5, FR-014, FR-016, and SC-005 require absolute `path` under scan root and `dependencyPath` only when `unknown` is true.

- [x] T028 [P] [US3] Add absolute path and `dependencyPath` tests in `target/internal/records/paths_test.go`
  Depends: T027
  Acceptance: `SPEC.md` SC-005 requires path values rooted under `start` and `dependencyPath` absence unless `unknown: true`.

- [x] T029 [P] [US3] Add direct-depth app propagation tests in `target/internal/app/f003_test.go`
  Depends: T011
  Acceptance: `SPEC.md` US3.3 and FR-004 require F-002-resolved `direct` values, including `0` and unbounded sentinel, to be passed unchanged as provider `depth`.

- [x] T030 [US3] Implement F-003 sorted-pass license sentinel rewrites in `target/internal/records/sentinels.go`
  Depends: T018, T024
  Acceptance: `SPEC.md` FR-020 through FR-022 and Sentinels lines 158-164 require private `UNLICENSED`, falsey `UNKNOWN`, and `unknown` guessed `*` rewrite without changing field order and without F-004 classifier scope.

- [x] T031 [P] [US3] Add sorted-pass sentinel tests in `target/internal/records/sentinels_test.go`
  Depends: T030
  Acceptance: `SPEC.md` SC-006 and `PLAN.md` Compatibility Risks lines 204-210 require private, falsey license, and `--unknown` guessed-license behavior to be locked by tests.

---

## Phase 6 — Validation

- [x] T032 Add end-to-end F-003 scan-to-record integration test in `target/internal/app/f003_integration_test.go`
  Depends: T013, T018, T019, T021, T025, T027, T030
  Acceptance: `SPEC.md` SC-001 and `PLAN.md` Validation Approach lines 194-202 require a sorted package map containing `abbrev@1.0.9` with repository `https://github.com/isaacs/abbrev-js` using committed/source-derived fixtures only.

- [x] T033 [P] Add source-observed `abbrev@1.0.9` fixture package in `target/internal/npmgraph/testdata/abbrev-fixture/node_modules/abbrev/package.json`
  Depends: T013
  Acceptance: `SPEC.md` US1.2, FR-023, and Output Contract lines 137-140 require the `abbrev@1.0.9` package identity and normalized repository value to remain reproducible for downstream renderers.

- [x] T034 [P] Add no-live-source validation guard tests in `target/internal/npmgraph/no_live_source_test.go`
  Depends: T015, T017, T020, T022, T023, T026, T028, T029, T031, T032, T033
  Acceptance: `SPEC.md` Edge Case line 104, Assumptions lines 224-231, and `AGENTS.md` parity verification rules require automated tests to avoid installing or executing the live Node source system.

---

## Phase 7 — Polish & Cross-Cutting

- [x] T035 Add records package documentation for F-003 boundaries in `target/internal/records/doc.go`
  Depends: T032, T034
  Acceptance: `PLAN.md` Scope lines 1-10 and `SPEC.md` Out of scope lines 21-27 require the package docs to keep F-004/F-005/F-006/F-007 ownership out of F-003 implementation.

- [x] T036 Add npmgraph package documentation for `read-installed` compatibility constraints in `target/internal/npmgraph/doc.go`
  Depends: T032, T034
  Acceptance: `SPEC.md` Algorithm Fidelity lines 183-185 and `ARCHITECTURE.md` Risks lines 236-239 require scanner docs to warn against naive filesystem-walk regressions.

---

## Implementation Strategy

- **Recommended start**: T001, then parallelize T002/T003/T004 before implementing ordered primitives and contracts.
- **Critical path**: T001 → T005 → T007 → T009 → T016 → T018 → T030 → T032 → T035/T036.
- **Risk areas**: T013 (`read-installed` semantics), T005/T007/T021 (ordered record/read-installed field-order compatibility), T025/T027/T030 (source sentinel and option edge cases), and T034 (no live source-system execution guard).
- **Handoff checklist**:
  - [x] All P1 story tasks T011 through T018 complete
  - [x] Validation phase tasks T032 through T034 pass
  - [x] BUILD.md evidence artifact started by Builder, not Task Decomposer
  - [x] No unresolved NEEDS CLARIFICATION items in task list

## Output Summary

- **Total task count**: 36 tasks; Phase 1: 4, Phase 2: 6, US1: 8, US2: 6, US3: 7, Validation: 3, Polish/Cross-cutting: 2.
- **Per-story task counts**: US1/P1 has 8 tasks (T011-T018); US2/P2 has 6 tasks (T019-T024); US3/P3 has 7 tasks (T025-T031).
- **Parallel opportunities**: 18 tasks are marked `[P]`: T002, T003, T004, T006, T010, T012, T014, T015, T017, T020, T022, T023, T026, T028, T029, T031, T033, T034.
- **MVP scope**: P1 viable first slice is T011 through T018 after Setup/Foundation, covering provider invocation, scanner model, flatten recursion, duplicate/incomplete handling, and lexical sorted package records.
- **Format validation**: Every task line uses `- [x] TNNN`, story-phase tasks carry `[US1]`, `[US2]`, or `[US3]`, every task description includes an exact target file path, and every task includes `Depends:` and `Acceptance:` details.
- **Unresolved items**: none; no `NEEDS CLARIFICATION` markers are present.
