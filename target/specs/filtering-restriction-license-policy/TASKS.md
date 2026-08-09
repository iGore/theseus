# Implementation Tasks: Filtering, Restriction, and License Policy Decisions

**Date**: 2026-06-22 | **Plan**: [PLAN.md](./PLAN.md) | **Spec**: [SPEC.md](./SPEC.md)
**Scope**: Reduced — F-005 only for `target/specs/filtering-restriction-license-policy/`; no renderer, dependency scanning, license classification, `BUILD.md`, `PLAN.md`, or `SPEC.md` work.

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 40    |
| Phase 1 (Setup)            | 4     |
| Phase 2 (Foundation)       | 6     |
| Story phases               | 22    |
| Polish / cross-cutting     | 3     |
| Parallel opportunities     | 16    |
| MVP scope (P1 stories)     | 7     |

## Dependency Graph

```text
Setup:        T001 -> T002 -> T003 -> T004
Foundation:   T005 -> T006 -> T007 -> T010
              T005 -> T008 [P] -> T009 [P] -> T010
US1 P1:       T010 -> T011 -> T012 -> T013 -> T014 -> T015 -> T016 -> T017
US2 P2:       T010,T017 -> T018 -> T019 -> T020 -> T021 -> T022
US3 P3:       T010 -> T023 -> T024 -> T025 -> T026
US4 P4:       T010,T022 -> T027 -> T028 -> T029 -> T030 -> T031 -> T032
Validation:   T017,T022,T026,T032 -> T033 [P], T034 [P], T035 [P], T036 [P] -> T037
Polish:       T037 -> T038 [P], T039 [P] -> T040
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

- Context7 used: `/golang/go` for `internal/` package visibility and standard `testing` / `go test` workflow guidance.
- Extension hooks: none; `.specify/extensions.yml` was not found at the repository root.
- Setup scripts: none found matching repository setup or agent-context update script patterns.
- Optional supporting artifacts: no `research.md`, `data-model.md`, `quickstart.md`, `contracts/`, or `/memory/constitution.md` found for this use case.

---

## Phase 1 — Setup

- [x] T001 Verify or create the Go module boundary for F-005 dependencies in `target/go.mod`
  Depends: none
  Acceptance: PLAN.md Technical Context lines 35-40 and ARCHITECTURE.md Structure Decision lines 30-53; module remains rooted under `target/` and includes only dependencies needed by F-005 and prior shared packages.

- [x] T002 [P] Confirm `github.com/github/go-spdx/v2/spdxexp v2.7.0` dependency availability in `target/go.mod`
  Depends: T001
  Acceptance: PLAN.md Decision Log lines 22-26 and ARCHITECTURE.md dependency decision lines 178-179; SPDX behavior stays behind `internal/spdxcompat`.

- [x] T003 [P] Create or verify the F-005 package directories in `target/internal/policy/doc.go`
  Depends: T001
  Acceptance: ARCHITECTURE.md boundaries lines 63-64 and PLAN.md Project Structure lines 73-80; no business logic is placed outside `target/internal/policy` for policy decisions.

- [x] T004 [P] Create F-005 test fixture notes for source-derived cases in `target/internal/policy/testdata/f005-fixtures.md`
  Depends: T001
  Acceptance: SPEC.md Output Contract lines 154-160 and PLAN.md WP-008 lines 143-148; fixture notes derive from committed artifacts only and do not run the live source system.

---

## Phase 2 — Foundation

- [x] T005 Define policy options and result contract in `target/internal/policy/types.go`
  Depends: T001
  Acceptance: SPEC.md FR-001 and Key Entities lines 208-213; accepts ordered package records keyed by exact `name@version` with mutable `licenses` and optional `private` fields.

- [x] T006 Define ordered record adapter expectations used by policy in `target/internal/policy/records.go`
  Depends: T005
  Acceptance: PLAN.md WP-001 lines 91-96 and ARCHITECTURE.md ordered data decision line 22; policy never relies on observable plain Go map iteration.

- [x] T007 Add baseline policy engine entrypoint in `target/internal/policy/apply.go`
  Depends: T006
  Acceptance: SPEC.md Ordered F-005 Pipeline lines 166-172 and PLAN.md Work Packages lines 91-148; entrypoint composes sorted pass, exclusion, restrictions, and fail checks without renderer behavior.

- [x] T008 [P] Define SPDX compatibility interface required by policy in `target/internal/spdxcompat/policy.go`
  Depends: T005
  Acceptance: SPEC.md FR-027 and ARCHITECTURE.md lines 63, 169, 178-179; policy depends on a narrow facade for correct/valid/satisfies behavior.

- [x] T009 [P] Lock F-005 SPDX adapter examples in `target/internal/spdxcompat/policy_test.go`
  Depends: T008
  Acceptance: SPEC.md Source Dependency Behavior lines 203-206 and Success Criteria SC-006; covers MIT/ISC, escaped-comma Apache literal, BSD alias, Public Domain, and satisfaction direction.

- [x] T010 Add shared in-memory policy test builders in `target/internal/policy/policy_test.go`
  Depends: T006
  Acceptance: SPEC.md User Scenarios independent tests lines 47, 65, 82, and 99; tests can exercise deterministic ordered package records without scanner or renderer dependencies.

---

## Phase 3 — Story: Exclude packages by license policy (Priority: P1)

- [x] T011 [US1] Implement source-compatible `--exclude` escaped-comma token parsing in `target/internal/policy/exclude.go`
  Depends: T007
  Acceptance: SPEC.md FR-008 and Exclusion Rule Table rows 1-2; only `--exclude` uses comma/escaped-comma parsing, unescaping, and trimming.

- [x] T012 [US1] Implement exact BSD alias expansion for exclusion tokens in `target/internal/policy/exclude.go`
  Depends: T011
  Acceptance: SPEC.md FR-009 and Exclusion Rule Table row 3; exact token `BSD` expands to `(0BSD OR BSD-2-Clause OR BSD-3-Clause OR BSD-4-Clause)` and no other token receives this alias.

- [x] T013 [US1] Partition valid SPDX and invalid literal exclusion tokens in `target/internal/policy/exclude.go`
  Depends: T012, T008
  Acceptance: SPEC.md FR-010 and PLAN.md WP-003 lines 105-110; valid tokens satisfy `spdxCorrect(token) == token`, invalid tokens remain literal matches, and valid expressions build `( <joined by OR> )`.

- [x] T014 [US1] Implement exclusion candidate evaluation over sorted records in `target/internal/policy/exclude.go`
  Depends: T013
  Acceptance: SPEC.md FR-011 through FR-016 and Exclusion Rule Table rows 7-14; scalar licenses become one-item candidate lists, `UNKNOWN` is always kept, guessed `*` loses its final character, candidate `BSD` aliases, invalid literals remove, and SPDX satisfaction removes.

- [x] T015 [P] [US1] Add table tests for MIT/ISC, escaped Apache, BSD alias, UNKNOWN preservation, and custom-license no-match in `target/internal/policy/exclude_test.go`
  Depends: T014, T010
  Acceptance: SPEC.md User Story 1 acceptance scenarios 1-5 and Success Criteria SC-001; key sets match the source-derived exclusion examples.

- [x] T016 [P] [US1] Add tests for guessed-license normalization and Public Domain exclusion behavior in `target/internal/policy/exclude_test.go`
  Depends: T014, T010
  Acceptance: SPEC.md FR-013, FR-015, Worked Transform Examples lines 195-200, and SC-006; `MIT*` normalizes to `MIT` for exclusion and invalid literals remain source-compatible.

- [x] T017 [US1] Integrate exclusion filtering into the policy pipeline in `target/internal/policy/apply.go`
  Depends: T014, T015, T016
  Acceptance: SPEC.md Ordered F-005 Pipeline step 4 and PLAN.md WP-004 lines 112-117; when `--exclude` is absent, `filtered = sorted`, and when present, output order is preserved after removals.

---

## Phase 4 — Story: Restrict result set by package identity and privacy (Priority: P2)

- [x] T018 [US2] Implement lexical sorted pass with private `UNLICENSED` override in `target/internal/policy/sort.go`
  Depends: T007
  Acceptance: SPEC.md FR-002, FR-003, FR-006, and User Story 2 scenario 3; keys process lexically and truthy-private records overwrite `licenses` to exact `UNLICENSED` before later filters.

- [x] T019 [US2] Implement `--packages` whitelist restriction in `target/internal/policy/restrict.go`
  Depends: T017, T018
  Acceptance: SPEC.md FR-017 and User Story 2 scenario 1; split only on `;`, do not trim tokens, and preserve filtered map iteration order for exact key matches.

- [x] T020 [US2] Implement `--excludePackages` blacklist rebuild semantics in `target/internal/policy/restrict.go`
  Depends: T019
  Acceptance: SPEC.md FR-018 and Edge Cases lines 112-114; blacklist rebuilds from the full filtered map, not from a prior whitelist result, and uses exact untrimmed `name@version` tokens.

- [x] T021 [US2] Implement `--excludePrivatePackages` deletion in `target/internal/policy/restrict.go`
  Depends: T020
  Acceptance: SPEC.md FR-019 and User Story 2 scenario 4; truthy-private records are deleted from the current restricted map after package whitelist/blacklist processing.

- [x] T022 [P] [US2] Add package restriction and private exclusion tests in `target/internal/policy/restrict_test.go`
  Depends: T019, T020, T021, T010
  Acceptance: SPEC.md User Story 2 acceptance scenarios 1-4 and Success Criteria SC-002, SC-005; verifies exact order, blacklist-overrides-whitelist construction, no trimming, `UNLICENSED`, and private deletion.

---

## Phase 5 — Story: Focus on unknown or guessed licenses (Priority: P3)

- [x] T023 [US3] Implement `--unknown` guessed-license rewrite in the sorted pass in `target/internal/policy/sort.go`
  Depends: T018
  Acceptance: SPEC.md FR-004 and User Story 3 scenarios 2-3; non-`UNKNOWN` current licenses containing `*` become exact `UNKNOWN` before only-unknown, exclusion, restrictions, and fail policies.

- [x] T024 [US3] Implement `--onlyunknown` inclusion gate in the sorted pass in `target/internal/policy/sort.go`
  Depends: T023
  Acceptance: SPEC.md FR-005 and User Story 3 scenarios 1 and 4; only current license values containing `*` or `UNKNOWN` are included, and `UNLICENSED` is not treated as unknown.

- [x] T025 [US3] Surface empty sorted result condition for downstream diagnostics in `target/internal/policy/apply.go`
  Depends: T024
  Acceptance: SPEC.md FR-007 and Edge Cases line 111; policy returns an equivalent `No packages found in this path..` condition without formatting downstream stderr.

- [x] T026 [P] [US3] Add unknown and only-unknown table tests in `target/internal/policy/unknown_test.go`
  Depends: T023, T024, T025, T010
  Acceptance: SPEC.md User Story 3 acceptance scenarios 1-4 and Success Criteria SC-003; fixtures cover `MIT`, `MIT*`, `UNKNOWN`, and `UNLICENSED` with and without the flags.

---

## Phase 6 — Story: Fail fast on forbidden or non-allowed licenses (Priority: P4)

- [x] T027 [US4] Implement fail policy list construction in `target/internal/policy/fail.go`
  Depends: T007
  Acceptance: SPEC.md FR-020, FR-022, and FR-026; `failOn` and `onlyAllow` split on `;`, trim tokens, ignore empty tokens, and programmatic both-values behavior leaves the active fail list populated.

- [x] T028 [US4] Implement restricted-record `--failOn` exact-match checks in `target/internal/policy/fail.go`
  Depends: T022, T027
  Acceptance: SPEC.md FR-020 and FR-021; first exact current-license match emits `Found license defined by the --failOn flag: "<license>". Exiting.` and returns exit code `1`.

- [x] T029 [US4] Implement restricted-record `--onlyAllow` substring checks in `target/internal/policy/fail.go`
  Depends: T028
  Acceptance: SPEC.md FR-022 and FR-023; first current-license value without any allowed token by source-compatible `indexOf` semantics emits `Package "<item>" is licensed under "<license>" which is not permitted by the --onlyAllow flag. Exiting.` and returns exit code `1`.

- [x] T030 [US4] Add exact F-005 diagnostics helper literals in `target/internal/diagnostics/policy.go`
  Depends: T027
  Acceptance: SPEC.md FR-021, FR-023, FR-024, FR-025 and Output Contract line 160; strings are centralized and preserve exact punctuation plus misspelling `delimeters`.

- [x] T031 [US4] Wire CLI preflight conflict and comma warning behavior for F-005 flags in `target/internal/diagnostics/preflight.go`
  Depends: T030
  Acceptance: SPEC.md FR-024 and FR-025; simultaneous `--failOn` and `--onlyAllow` exits before scanning with code `1`, while comma-containing values warn and continue.

- [x] T032 [P] [US4] Add fail policy and preflight diagnostic tests in `target/internal/policy/fail_test.go`
  Depends: T028, T029, T030, T031, T010
  Acceptance: SPEC.md User Story 4 acceptance scenarios 1-5 and Success Criteria SC-004; asserts first violation, exact stderr strings, warning continuation, mutual-exclusion stop, and exit code `1`.

---

## Phase 7 — Validation

- [x] T033 [P] Add end-to-end in-memory F-005 pipeline order tests in `target/internal/policy/apply_test.go`
  Depends: T017, T022, T026, T032
  Acceptance: SPEC.md Algorithm Fidelity lines 164-172; validates sorted pass, exclusion, restrictions, private deletion, and fail policies run in source order.

- [x] T034 [P] Add CLI/app reachability tests for preflight versus policy execution in `target/internal/app/policy_integration_test.go`
  Depends: T031, T032
  Acceptance: PLAN.md WP-007 lines 136-141 and Validation Approach line 168; conflict prevents policy reachability and comma warning allows policy reachability.

- [x] T035 [P] Add SPDX facade integration tests consumed by policy in `target/internal/spdxcompat/policy_integration_test.go`
  Depends: T009, T017
  Acceptance: SPEC.md FR-027 and ARCHITECTURE.md risk lines 236-239; adapter examples pass before broader policy use is considered complete.

- [x] T036 [P] Add source-derived fixture coverage documentation in `target/internal/policy/testdata/f005-fixtures.md`
  Depends: T015, T016, T022, T026, T032
  Acceptance: SPEC.md Golden baseline lines 154-160 and Success Criteria SC-001 through SC-007; documents each fixture's SPEC/Analyzer provenance without live source execution.

- [x] T037 Run full Go test validation from the target module using `target/go.mod`
  Depends: T033, T034, T035, T036
  Acceptance: PLAN.md Dependencies and Execution Order lines 152-159 and Validation Approach lines 161-169; `go test ./...` passes from `target/` and no F-006 renderer byte assertions are introduced.

---

## Phase 8 — Polish & Cross-Cutting

- [x] T038 [P] Review F-005 public comments and package boundaries in `target/internal/policy/doc.go`
  Depends: T037
  Acceptance: ARCHITECTURE.md Boundaries lines 55-66 and Context7 `/golang/go` internal package guidance; exported names are documented only where needed and F-005 internals remain under `internal/`.

- [x] T039 [P] Review centralized diagnostics usage to avoid duplicated literals in `target/internal/diagnostics/policy.go`
  Depends: T037
  Acceptance: PLAN.md Handoff Notes lines 184-188; policy and preflight tests assert strings through centralized diagnostics helpers without divergent copies.

- [x] T040 Final F-005 conformance sweep for policy rows in `target/specs/filtering-restriction-license-policy/TASKS.md`
  Depends: T038, T039
  Acceptance: ARCHITECTURE.md Conformance Matrix rows 215-222 and SPEC.md Success Criteria SC-001 through SC-007; all F-005 task checkboxes remain traceable to SPEC/PLAN and no unresolved clarification markers remain.

---

## Implementation Strategy

- **Recommended start**: T001 through T010 establish the Go module/package boundaries, policy contracts, ordered record adapters, SPDX facade, and reusable in-memory fixtures.
- **Critical path**: T001 → T005 → T006 → T007 → T011 → T012 → T013 → T014 → T017 → T018 → T019 → T020 → T021 → T027 → T028 → T029 → T032 → T037 → T040.
- **Risk areas**: T009/T035 for SPDX compatibility, T014 for exclusion evaluator semantics, T020 for blacklist rebuild order, T027 for programmatic both-values behavior, and T031 for exact preflight routing.
- **Handoff checklist**:
  - [x] All P1 story tasks complete: T011-T017
  - [x] Validation phase tasks pass: T033-T037
  - [x] BUILD.md evidence artifact started by Builder, not Task Decomposer
  - [x] No unresolved clarification items in task list

## Output Summary

- **Total task count**: 40 tasks; Phase 1: 4, Phase 2: 6, US1: 7, US2: 5, US3: 4, US4: 6, Validation: 5, Polish: 3.
- **Per-story task counts**: US1/P1 exclusion policy: 7; US2/P2 package/private restrictions: 5; US3/P3 unknown handling: 4; US4/P4 fail policies and preflight: 6.
- **Parallel opportunities**: 16 `[P]` tasks in setup, foundation, story tests, validation, and polish: T002, T003, T004, T008, T009, T015, T016, T022, T026, T032, T033, T034, T035, T036, T038, T039. Builder may execute only dependency-safe subsets concurrently.
- **MVP scope**: P1 viable slice is T011-T017 after foundation T001-T010; it covers SPDX/BSD alias and `--exclude` filtering behavior.
- **Format validation**: Every task uses the strict checklist line with sequential ID, exact file path, `Depends:` note, and `Acceptance:` note; story-phase tasks carry `[US1]` through `[US4]` labels.
- **Unresolved items**: none; prior PLAN open questions are represented as executable tasks T006, T009, and T035 rather than clarification blockers.
