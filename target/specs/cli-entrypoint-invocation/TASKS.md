# Implementation Tasks: CLI entry point and invocation model

**Date**: 2026-06-22 | **Plan**: [PLAN.md](./PLAN.md) | **Spec**: [SPEC.md](./SPEC.md)
**Scope**: Reduced — F-001 only, covering `license-checker [flags]`, ordered entry-point preflight behavior, stderr/stdout/exit codes, and downstream handoff seams. No F-002 parser compatibility beyond F-001 fields and `-h`/`-v` aliases.

## Task Summary

| Metric                     | Count |
|----------------------------|-------|
| Total tasks                | 26    |
| Phase 1 (Setup)            | 3     |
| Phase 2 (Foundation)       | 5     |
| Story phases               | 13    |
| Polish / cross-cutting     | 2     |
| Parallel opportunities     | 9     |
| MVP scope (P1 stories)     | 4     |

## Context and Hooks

- Context7 `/golang/go` was consulted for `package main` executable conventions, `internal/` package visibility, and standard `testing`/`os/exec` CLI test planning.
- `.specify/extensions.yml`: not present; no executable `hooks.before_tasks` or `hooks.after_tasks` surfaced.
- Repository setup or agent-context scripts: none found; no setup output captured.
- `/memory/constitution.md`: not present.

## Dependency Graph

```text
T001
  -> T002
      -> T003 [P]
      -> T004
          -> T005 [P]
          -> T006 [P]
          -> T007 [P]
          -> T008 [P]
              -> T009
                  -> T010 [P]
                  -> T011
                  -> T012
                      -> T013
                          -> T014 [P]
                          -> T015
                          -> T016
                              -> T017 [P]
                              -> T018
                              -> T019
                              -> T020
                                  -> T021 [P]
                                      -> T022
                                          -> T023
                                              -> T024
                                                  -> T025
                                                  -> T026
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

- [x] T001 Create Go module metadata for the target CLI in `target/go.mod`
  Depends: none
  Acceptance: `PLAN.md` WP-001 and `ARCHITECTURE.md` Technical Context require a Go module rooted under `target/` with no external CLI framework for F-001.

- [x] T002 Create thin executable entry-point scaffold in `target/cmd/license-checker/main.go`
  Depends: T001
  Acceptance: `SPEC.md` FR-001/FR-002 and `ARCHITECTURE.md` lines 87-94 require command name `license-checker`, `package main`, process wiring only, and no subcommand framework.

- [x] T003 [P] Create CLI subprocess test helper scaffold in `target/internal/testsupport/clitest/clitest.go`
  Depends: T001
  Acceptance: `PLAN.md` WP-006 and `ARCHITECTURE.md` Testing decision require standard `testing` plus `os/exec` integration tests for the built binary.

---

## Phase 2 — Foundation

- [x] T004 Define `app.Invocation` and `Run` orchestration seam in `target/internal/app/app.go`
  Depends: T002
  Acceptance: `PLAN.md` WP-001/WP-002 and `ARCHITECTURE.md` boundaries require `internal/app.Run(ctx, Invocation) int` to receive args, env, cwd, stdout, stderr, and downstream dependencies.

- [x] T005 [P] Define minimal F-001 parsed option fields in `target/internal/cli/options.go`
  Depends: T004
  Acceptance: `SPEC.md` FR-003/FR-011 and `PLAN.md` WP-002 require `Help`, `Version`, `FailOn`, and `OnlyAllow` fields sufficient for preflight decisions.

- [x] T006 [P] Define minimal F-001 parser interface and alias mapping contract in `target/internal/cli/parser.go`
  Depends: T004
  Acceptance: `SPEC.md` FR-011 and assumptions require `-h` to map to `Help`, `-v` to map to `Version`, and all broader `nopt` compatibility to remain deferred to F-002.

- [x] T007 [P] Define downstream pipeline runner interface in `target/internal/app/pipeline.go`
  Depends: T004
  Acceptance: `SPEC.md` FR-008 and `PLAN.md` WP-002/WP-005 require a stub-friendly downstream handoff seam called exactly once only when preflight permits.

- [x] T008 [P] Centralize F-001 exit-code constants in `target/internal/diagnostics/exit.go`
  Depends: T004
  Acceptance: `SPEC.md` FR-004/FR-005/FR-006 require help exit `0`, version exit `1`, conflict exit `1`, and no immediate exit for comma warnings.

---

## Phase 3 — Story: Invoke the license-checker executable (Priority: P1)

- [x] T009 [US1] Wire `main` to construct `app.Invocation` and call `os.Exit(app.Run(...))` in `target/cmd/license-checker/main.go`
  Depends: T004, T007
  Acceptance: `SPEC.md` US1 scenario 1, FR-001, FR-002, FR-008, and `PLAN.md` WP-001 require a single executable command with business logic outside `main.go`.

- [x] T010 [P] [US1] Implement F-001 argument recognition without subcommand dispatch in `target/internal/cli/parser.go`
  Depends: T005, T006
  Acceptance: `SPEC.md` US1 scenario 2, FR-002, FR-011, and Edge Cases require no named subcommand layer and only F-001 fields/aliases at this stage.

- [x] T011 [US1] Implement non-preflight downstream handoff in `target/internal/app/app.go`
  Depends: T007, T009, T010
  Acceptance: `SPEC.md` FR-008 and SC-002 require parsed options to be passed through unchanged and the downstream runner to be invoked exactly once when no preflight exit applies.

- [x] T012 [US1] Add no-preflight and no-subcommand unit coverage in `target/internal/app/app_test.go`
  Depends: T011
  Acceptance: `SPEC.md` US1 Independent Test and SC-005 require stubbed pipeline call-count verification and no subcommand-specific behavior.

---

## Phase 4 — Story: Read help from the CLI (Priority: P2)

- [x] T013 [US2] Implement exact help banner and usage constants in `target/internal/diagnostics/help.go`
  Depends: T008
  Acceptance: `SPEC.md` FR-004/FR-010 and Output Contract lines 116-143 require `license-checker@25.0.1`, ordered help lines, preserved typos, stderr routing, and documented final-spacing limitation handling.

- [x] T014 [P] [US2] Implement help preflight writer behavior in `target/internal/diagnostics/preflight.go`
  Depends: T013
  Acceptance: `SPEC.md` FR-003/FR-004 and Algorithm Fidelity require `Help` to be evaluated first, write only to stderr, skip the pipeline, and return exit `0`.

- [x] T015 [US2] Integrate help preflight short-circuiting in `target/internal/app/app.go`
  Depends: T014
  Acceptance: `SPEC.md` US2 scenarios 1-2 and Edge Cases require `--help`/`-h` to win over version, conflict, warning, and scan behavior.

- [x] T016 [US2] Add help unit and CLI integration tests in `target/internal/diagnostics/preflight_test.go`
  Depends: T015, T003
  Acceptance: `SPEC.md` SC-001/SC-003 require `--help` and `-h` stderr help output, empty stdout, exit `0`, and zero downstream calls.

---

## Phase 5 — Story: Observe version and entry-point policy diagnostics (Priority: P3)

- [x] T017 [P] [US3] Implement exact version diagnostic behavior in `target/internal/diagnostics/preflight.go`
  Depends: T014
  Acceptance: `SPEC.md` FR-005, US3 scenario 1, and Output Contract require exactly `25.0.1\n` on stderr, empty stdout, no pipeline call, and exit `1` when help is absent.

- [x] T018 [US3] Implement failOn-onlyAllow conflict diagnostic behavior in `target/internal/diagnostics/preflight.go`
  Depends: T017
  Acceptance: `SPEC.md` FR-006, US3 scenario 2, and Edge Cases require exact mutual-exclusion stderr, exit `1`, no pipeline call, and conflict precedence over comma warnings.

- [x] T019 [US3] Implement comma warning behavior for policy flags in `target/internal/diagnostics/preflight.go`
  Depends: T018
  Acceptance: `SPEC.md` FR-007, US3 scenario 3, and Output Contract require exact `--failOn`/`--onlyAllow` warning literals preserving `delimeters` and continuing to downstream handoff.

- [x] T020 [US3] Integrate full preflight order in app orchestration in `target/internal/app/app.go`
  Depends: T019
  Acceptance: `SPEC.md` FR-003 and Algorithm Fidelity require order: help, version, conflict, comma warning, then downstream scan/render handoff.

- [x] T021 [P] [US3] Add version, conflict, comma-warning, and priority tests in `target/internal/diagnostics/preflight_test.go`
  Depends: T020
  Acceptance: `SPEC.md` SC-001/SC-003/SC-004 require exact stderr, empty stdout for preflight exits, exit codes, warning-continuation behavior, and branch-priority coverage.

---

## Phase 6 — Validation

- [x] T022 Add end-to-end CLI integration tests for F-001 invocations in `target/cmd/license-checker/main_test.go`
  Depends: T012, T016, T021
  Acceptance: `PLAN.md` WP-006 and `SPEC.md` SC-001 through SC-005 require subprocess coverage for `--help`, `-h`, `--version`, `-v`, conflict, comma warnings, and no-preflight invocation.

- [x] T023 Add golden stderr fixture file for F-001 diagnostics in `target/testdata/f001_cli_preflight.golden`
  Depends: T022
  Acceptance: `ARCHITECTURE.md` Acceptance reference and `SPEC.md` Output Contract require committed golden strings derived from source snippets and SPEC exact literals, without live source execution.

- [x] T024 Verify all F-001 packages with `go test ./...` from module `target/`
  Depends: T023
  Acceptance: `PLAN.md` WP-006 and `ARCHITECTURE.md` Testing decision require the full target module test suite to pass from `target/` before claiming implementation completion.

---

## Phase 7 — Polish & Cross-Cutting

- [x] T025 Ensure F-001 code documents deferred F-002 parser scope in `target/internal/cli/parser.go`
  Depends: T024
  Acceptance: `SPEC.md` assumptions and `PLAN.md` Open Questions require unknown operands, path coercion, defaults, `--no-` forms, and full `nopt` edge behavior to remain explicitly deferred, not guessed.

- [x] T026 Ensure F-001 package boundaries contain no downstream implementation leakage in `target/internal/app/app.go`
  Depends: T024
  Acceptance: `SPEC.md` out-of-scope section and `PLAN.md` WP-005 require no F-003 through F-007 scan/render behavior except injected interfaces/stubs and downstream exit-code propagation.

---

## Implementation Strategy

- **Recommended start**: T001 through T004 to establish the module, executable, and app seam before diagnostics or parser details.
- **Critical path**: T001 → T002 → T004 → T008 → T013 → T014 → T017 → T018 → T019 → T020 → T021 → T022 → T023 → T024.
- **Risk areas**: T013 help final-spacing parity limitation; T010 parser behavior must not drift into F-002; T019 warnings must continue to downstream rather than exiting; T026 must prevent downstream slice scope creep.
- **Handoff checklist**:
  - [x] All P1 story tasks complete
  - [x] Validation phase tasks pass
  - [x] BUILD.md evidence artifact started
  - [x] No unresolved NEEDS CLARIFICATION items in task list

## Decomposition Summary

- **Per-phase breakdown**: Setup 3, Foundation 5, US1 4, US2 4, US3 5, Validation 3, Polish/Cross-cutting 2.
- **Per-story counts**: US1/P1 has 4 tasks (T009-T012), US2/P2 has 4 tasks (T013-T016), US3/P3 has 5 tasks (T017-T021).
- **Parallel opportunities**: 9 tasks are marked `[P]`: T003, T005, T006, T007, T008, T010, T014, T017, T021.
- **MVP scope**: P1 MVP is T009-T012 after setup/foundation prerequisites T001-T008; it delivers the single executable and exactly-once downstream handoff for non-preflight invocation.
- **Format validation**: Every task uses a sequential checklist ID, includes an exact `target/` file path or module path, and has `Depends:` plus `Acceptance:` details.
- **Unresolved items**: none; the help trailing-spacing limitation is documented as a known parity limitation from `SPEC.md`, not a `NEEDS CLARIFICATION` blocker.
