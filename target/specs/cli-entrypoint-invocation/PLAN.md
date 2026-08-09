# Implementation Plan: CLI entry point and invocation model

**Date**: 2026-06-22 | **Spec**: `target/specs/cli-entrypoint-invocation/SPEC.md`  
**Input**: Feature specification from `target/specs/cli-entrypoint-invocation/SPEC.md`; shared architecture from `target/specs/ARCHITECTURE.md`  
**Scope**: F-001 only — top-level `license-checker [flags]` entry point, preflight behavior, stderr/stdout/exit-code contract, and handoff to downstream scan/render pipeline.

## Summary

Implement a Go CLI executable named `license-checker` with a thin `package main` at `target/cmd/license-checker/main.go` delegating to `internal/app.Run`. F-001 work is limited to source-compatible command invocation and ordered preflight handling for `--help`/`-h`, `--version`/`-v`, simultaneous `--failOn` and `--onlyAllow`, comma warnings for the one present policy flag, and exactly-once handoff to the downstream pipeline when no preflight exit applies. The plan follows `ARCHITECTURE.md` decisions: no full CLI framework, `internal/` package boundaries, stderr diagnostics at the diagnostics/app boundary, standard `testing` plus `os/exec` CLI integration tests, and golden strings from `SPEC.md`.

Context7 note: `/golang/go` was consulted for Go executable conventions (`package main` + `main()`), `internal/` boundaries, standard `os/exec` subprocess testing, and stdlib-first CLI/test planning.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md` F-001 only.
- [x] Technical context filled with concrete decisions from `ARCHITECTURE.md`.
- [x] Work packages derived from F-001 functional requirements and success criteria.
- [x] Dependencies and execution order documented with explicit `Depends:` notes.
- [x] Validation approach captured with SPEC/ARCHITECTURE traceability.

## Decision Log

- [x] Use `target/cmd/license-checker/main.go` only as process wiring: args/env/stdout/stderr/getwd and `os.Exit(app.Run(...))`. Depends on `ARCHITECTURE.md` lines 55-59 and 87-94.
- [x] Put entry orchestration in `internal/app` and preflight diagnostics in `internal/diagnostics`; parser-owned compatibility details remain in `internal/cli` and F-002. Depends on `ARCHITECTURE.md` lines 57-66 and `SPEC.md` FR-003/FR-008.
- [x] Do not introduce Cobra/urfave or subcommand dispatch. Depends on `SPEC.md` FR-001/FR-002 and `ARCHITECTURE.md` package decision rejecting full CLI frameworks.
- [x] Treat version as fixed `25.0.1` for this slice. Depends on `SPEC.md` assumptions and FR-005.
- [x] Preserve exact stderr literals, stdout emptiness, and exit codes for preflight paths; final help trailing byte identity remains the documented parity limitation. Depends on `SPEC.md` Output Contract and Success Criteria.

## Open Questions

- [ ] Final live Node `console.error(usage.join('\n'), '\n')` trailing-space byte identity is not available in-flow; preserve source-derived structure and leave live differential confirmation to maintainers as documented in `SPEC.md`.
- [ ] Full unknown operand, path coercion, defaults, and `nopt` edge behavior are intentionally deferred to F-002; F-001 only needs the parsed fields used by preflight.

## Technical Context

**Language/Version**: Go, module under `target/`; exact Go toolchain version not pinned in F-001 artifacts.  
**Primary Dependencies**: Go standard library for this slice (`context`, `os`, `io`, `strings`, `testing`, `os/exec`); no external CLI framework.  
**Storage**: N/A for F-001 except downstream invocation placeholders; filesystem scan/output belongs to later slices.  
**Testing**: Standard `testing`, table-driven unit tests for preflight, and CLI integration tests using built binary/`os/exec`; run with `go test ./...` from `target/`.  
**Target Platform**: POSIX-like CLI on darwin/linux first, matching architecture.  
**Project Type**: Go CLI application.  
**Performance Goals**: No F-001 throughput goal; preflight must short-circuit before scanner setup or execution when applicable.  
**Constraints**: No live source-system execution in-flow; exact stderr/stdout routing and exit codes; no subcommands; no implementation beyond F-001 contracts.  
**Scale/Scope**: One executable and minimal app/cli/diagnostics seams needed to hand off to F-002+ downstream implementation.

## Constitution Check

No additional constitution or extension files were consulted because this run was explicitly constrained to the matching `SPEC.md` plus shared `ARCHITECTURE.md`. Applicable repo workflow rules from the role prompt are reflected here: Go default target, artifact under `target/specs/<slug>/`, Planner owns only `PLAN.md`, no `TASKS.md`, no code.

## Project Structure

### Documentation (this feature)

```text
target/specs/cli-entrypoint-invocation/
├── SPEC.md
└── PLAN.md              # This file; Planner-owned
```

### Source Code (repository root)

```text
target/
├── go.mod
├── cmd/
│   └── license-checker/
│       └── main.go              # package main; process wiring only
└── internal/
    ├── app/                     # Run(ctx, Invocation) int; orchestration seam
    ├── cli/                     # parsed Options fields needed by F-001/F-002
    └── diagnostics/             # preflight decision table and exact stderr literals
```

**Structure Decision**: Use the architecture-selected standard Go CLI layout with a thin command and private `internal/` packages. For F-001, only `cmd/license-checker`, `internal/app`, `internal/cli`, and `internal/diagnostics` are required; downstream scanner/render packages may be represented by interfaces or stubs for validation but their behavior is out of scope.

## Work Packages

### WP-001 — Establish executable and invocation seam

- **Goal**: Create the `license-checker` command surface as `target/cmd/license-checker/main.go` with no subcommands.
- **Work**:
  - Wire `os.Args[1:]`, environment, stdout, stderr, current working directory, and `os.Exit` through an `app.Invocation` value.
  - Keep all business logic out of `main.go`.
  - Ensure binary naming/build path produces executable command `license-checker`.
- **Depends:** none.
- **Traceability:** `SPEC.md` FR-001, FR-002, SC-005; `ARCHITECTURE.md` Entry point and Structure Decision.

### WP-002 — Define F-001 option and invocation contracts

- **Goal**: Provide the minimal internal types that allow preflight and handoff without implementing full F-002 parser semantics.
- **Work**:
  - Define an `Options` shape containing at least `Help`, `Version`, `FailOn`, and `OnlyAllow`.
  - Define an app-level downstream runner/scanner interface that can be stubbed in tests and called exactly once on non-preflight paths.
  - Ensure `-h` maps to `Help` and `-v` maps to `Version` at this surface, while documenting that full `nopt` compatibility belongs to F-002.
- **Depends:** WP-001.
- **Traceability:** `SPEC.md` FR-003, FR-008, FR-011; `ARCHITECTURE.md` boundaries for `internal/app` and `internal/cli`.

### WP-003 — Implement ordered preflight decision table

- **Goal**: Encode the exact F-001 branch order before any downstream scan/render work.
- **Work**:
  - Evaluate in order: help; version; `failOn && onlyAllow`; comma warning for exactly one present policy flag; downstream handoff.
  - Guarantee `--help` wins over all other flags and exits `0`.
  - Guarantee `--version` wins over conflict/warning when help is absent and exits `1`.
  - Guarantee conflict suppresses comma warnings and exits `1`.
  - Guarantee comma warnings do not cause immediate exit and are followed by downstream handoff.
- **Depends:** WP-002.
- **Traceability:** `SPEC.md` FR-003 through FR-008, Algorithm Fidelity table, Edge Cases; `ARCHITECTURE.md` Ordered stage sequence items 1-2.

### WP-004 — Centralize exact diagnostics and help text

- **Goal**: Preserve source-derived stderr strings and stdout emptiness for all F-001 preflight outcomes.
- **Work**:
  - Store fixed version `25.0.1` for F-001 diagnostics.
  - Emit help banner `license-checker@25.0.1\n` and the usage option lines in exact source order, preserving typos (`seperated`, `delimeters`).
  - Emit exact version, mutual-exclusion, `--failOn` comma-warning, and `--onlyAllow` comma-warning strings.
  - Write diagnostics to stderr only; never write preflight output to stdout.
  - Preserve the documented help final-spacing limitation as a test note rather than inventing live byte behavior.
- **Depends:** WP-003.
- **Traceability:** `SPEC.md` FR-004 through FR-010, Output Contract, SC-001, SC-003, SC-004; `ARCHITECTURE.md` diagnostics boundary.

### WP-005 — Wire downstream handoff without implementing downstream slices

- **Goal**: Ensure F-001 reaches the downstream scan/render pipeline exactly once only when preflight permits.
- **Work**:
  - In `internal/app.Run`, call the injected downstream pipeline once after successful preflight.
  - Pass parsed options through unchanged for downstream consumers.
  - Return downstream exit code only after handoff; do not define scan/render behavior in this slice.
  - Ensure help/version/conflict paths never call the downstream stub.
- **Depends:** WP-003, WP-004.
- **Traceability:** `SPEC.md` FR-008, SC-002; `ARCHITECTURE.md` composition root and ordered stage sequence.

### WP-006 — F-001 validation suite

- **Goal**: Lock F-001 behavior before Task Decomposer/Builder move into implementation detail.
- **Work**:
  - Unit-test preflight table for branch priority and exact result metadata.
  - Integration-test the built `license-checker` binary with `--help`, `-h`, `--version`, `-v`, conflict, comma warnings, and no-preflight invocation using a stubbed downstream pipeline.
  - Assert exit codes, stderr, empty stdout for preflight exits, and downstream call count.
  - Run `go test ./...` from `target/`.
- **Depends:** WP-001 through WP-005.
- **Traceability:** `SPEC.md` Success Criteria SC-001 through SC-005; `ARCHITECTURE.md` Testing decision and Conformance Matrix rows for `--help`, `--version`, `--failOn`, and `--onlyAllow`.

## Dependency and Execution Order

```text
WP-001
  -> WP-002
      -> WP-003
          -> WP-004
          -> WP-005
              -> WP-006
```

F-002 parser work may later deepen `internal/cli`, but F-001 must be testable with the minimal parser/option fields above. F-003 through F-007 must not be pulled into this implementation except through injected interfaces/stubs.

## Validation Approach

- **Static design validation**: Confirm `main.go` contains process wiring only and no subcommand framework or business logic.
- **Unit validation**: Table-drive all ordered preflight rows from `SPEC.md` Algorithm Fidelity.
- **CLI integration validation**: Build/run the binary with `os/exec`, capturing stdout/stderr and process exit status.
- **Golden/string validation**: Compare exact stderr literals from `SPEC.md` Output Contract for version, conflict, and comma warnings; compare help body lines/order and documented trailing-spacing behavior.
- **Handoff validation**: Use a stubbed downstream runner to assert exactly-zero calls for preflight exits and exactly-one call for no-preflight and comma-warning cases.
- **Command validation**: Verify no subcommand dispatch path exists; arbitrary token handling remains a parser compatibility concern for F-002.

## Complexity Tracking

No constitution-driven complexity violations are identified for F-001. The chosen seams (`app`, `cli`, `diagnostics`) are required by `ARCHITECTURE.md` to keep process wiring, parsing/preflight, and diagnostics independently testable while avoiding a full CLI framework.
