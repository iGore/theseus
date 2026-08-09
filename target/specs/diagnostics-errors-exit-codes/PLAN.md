# Implementation Plan: Diagnostics, Error Handling, and Exit Codes

**Date**: 2026-06-22 | **Spec**: `target/specs/diagnostics-errors-exit-codes/SPEC.md`  
**Input**: Feature specification from `target/specs/diagnostics-errors-exit-codes/SPEC.md` and shared architecture from `target/specs/ARCHITECTURE.md`

## Summary

Plan the Go implementation of F-007 only: source-compatible diagnostics, stderr/stdout routing, debug namespace enablement, and process exit-code behavior for the `license-checker` CLI. The implementation belongs in the architecture-selected Go CLI layout under `target/`, with `cmd/license-checker/main.go` delegating process exit to `internal/app.Run`, exact messages centralized in `internal/diagnostics`, preflight checks coupled to `internal/cli`, policy-failure exits in `internal/policy`, file-copy warnings in `internal/render`/`internal/files`, and debug namespace gating in `internal/debuglog`.

All exact strings, misspellings, punctuation, interpolation quotes, channel choices, and exit-code decisions from the F-007 Output Contract must be preserved. Known non-identical/open rows remain explicit: full help trailing whitespace after Node's two-argument `console.error`, debug dependency byte formatting, and Node `console.error(err)` object formatting.

Context7 was used for Go planning validation against `/golang/go`: use a thin `package main`, keep implementation in `internal/`, pass stdout/stderr as writers for testability, call `os.Exit` only from main, and validate CLI behavior with standard `testing` plus `os/exec` capture of stdout/stderr/status.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md` for F-007 only
- [x] Technical context filled with concrete decisions and documented open compatibility decisions
- [x] Work packages derived from the assigned diagnostics/error/exit-code scope
- [x] Dependencies and execution order documented
- [x] Validation approach captured

## Decision Log

- [x] Implement F-007 as CLI-observable behavior only; do not expand parsing, scanning, filtering, or rendering semantics beyond the F-007 trigger/diagnostic points. Depends: `SPEC.md` Scope/Out of scope, `ARCHITECTURE.md` Composition Root.
- [x] Centralize exact diagnostic literals and exit-code decisions in `internal/diagnostics`; expose helpers to `cli`, `policy`, `app`, and renderer side-effect paths instead of duplicating strings. Depends: `SPEC.md` FR-015, `ARCHITECTURE.md` lines 64-66.
- [x] Keep `os.Exit` in `cmd/license-checker/main.go`; `internal/app.Run` returns an int exit code. Depends: `ARCHITECTURE.md` Entry point and Context7 `/golang/go` CLI guidance.
- [x] Preserve compatible, not byte-identical, status for debug formatting and callback error object formatting until maintainers provide approved baselines. Depends: `SPEC.md` FR-011/FR-014 and `ARCHITECTURE.md` Conformance Matrix.
- [x] Use standard library `testing` and `os/exec` for CLI integration tests that assert stderr, stdout, and status. Depends: `ARCHITECTURE.md` Testing decision and Context7 `/golang/go` examples.

## Open Questions

- [ ] Maintainers-approved byte baseline for the second callback-error line equivalent to Node `console.error(err)` is still required before claiming byte-identical parity. Current plan preserves `Found error` exactly and treats the subsequent error formatting as compatible/open. Depends: `SPEC.md` FR-010/FR-011, SA-007.
- [ ] Full byte-exact help trailing whitespace after source `console.error(usage.join('\n'), '\n')` remains uncaptured. Current plan requires exact listed help lines and compatible documented trailing behavior. Depends: `SPEC.md` Output Contract whitespace notes.
- [ ] Debug package timestamp/color/prefix bytes are not captured. Current plan requires `license-checker:log` and `license-checker:error` namespace enablement and no exit-code impact only. Depends: `SPEC.md` FR-014, SA-006/SA-007.

## Technical Context

**Language/Version**: Go, module rooted under `target/`; exact Go version follows `target/go.mod` created by earlier implementation tasks or repository default.  
**Primary Dependencies**: Go standard library for this slice (`os`, `io`, `fmt`, `strings`, `testing`, `os/exec`); no new third-party dependency for F-007 diagnostics.  
**Storage**: N/A for diagnostics; consumes filesystem/render state only for `--files` missing-license warnings owned by F-006 side effects.  
**Testing**: Standard `go test ./...`, table-driven unit tests for diagnostics/preflight/policy helpers, and CLI integration tests using built binary or `go run`/`os/exec` with captured stdout/stderr/status.  
**Target Platform**: POSIX-like Go CLI on darwin/linux first, consistent with shared architecture.  
**Project Type**: Single Go CLI command named `license-checker`.  
**Performance Goals**: No extra performance-sensitive work in this slice; diagnostics must short-circuit before scanning for help/version/mutual-exclusion fatal paths.  
**Constraints**: Preserve exact stderr/stdout routing; do not add stdout copies of diagnostics; preserve exit `0` for help and normal success, exit `1` for version/mutual-exclusion/policy failures; do not invent non-zero exit solely for callback-equivalent init errors; do not invent structured replacement output for Node error object/debug formatting.  
**Scale/Scope**: Bounded to F-007 diagnostic events over the package set produced by F-005/F-006; no new scan/render scope.

## Constitution Check

No `/memory/constitution.md` or `.specify/extensions.yml` was consulted because this Planner invocation explicitly limited inputs to the assigned `SPEC.md` and shared `ARCHITECTURE.md`. Applicable workflow gates from those inputs are satisfied:

- Spec-first gate: plan maps only to F-007 requirements and references the existing architecture.
- Go target gate: plan assumes the shared Go CLI architecture under `target/`.
- Artifact ownership gate: this stage owns only this `PLAN.md`; no `TASKS.md`, `BUILD.md`, code, or auxiliary artifacts are produced.
- Parity gate: exact source-derived diagnostics are planned as byte/channel/status assertions; declared compatible/open rows remain documented instead of being silently treated as identical.

## Project Structure

### Documentation (this feature)

```text
target/specs/diagnostics-errors-exit-codes/
├── SPEC.md
└── PLAN.md              # This file
```

### Source Code (repository root)

```text
target/
├── go.mod
├── cmd/
│   └── license-checker/
│       └── main.go              # calls app.Run and os.Exit(code)
└── internal/
    ├── app/                     # ordered pipeline, callback-error handling, exit-code return
    ├── cli/                     # parsed flags plus preflight trigger ordering
    ├── debuglog/                # DEBUG=license-checker* namespace enablement
    ├── diagnostics/             # exact literals, stderr writers, exit-code results
    ├── policy/                  # failOn/onlyAllow diagnostics after F-005 restrictions
    └── render/                  # --files missing-license-file warning integration
```

**Structure Decision**: Use the shared architecture's standard Go CLI layout. F-007 work should not create public `pkg/` APIs. Diagnostic strings and status decisions should be centralized in `internal/diagnostics`, while trigger ownership remains at the architectural stage that observes the condition.

## Work Packages

### WP-001 — Diagnostic constants and writer contract

**Goal**: Create a single source of truth for exact F-007 diagnostic strings, interpolation functions, and channel/exit metadata.

**Implementation scope**:

- Define exact literals for version, mutual-exclusion, comma warnings, failOn failure, onlyAllow failure, callback error prefix, and missing license-file warning.
- Preserve source misspellings and literals: `delimeters`, `can not`, `Found error`, `Exiting.`, and `no license file found for:`.
- Expose writer helpers that write one line to the intended writer without duplicating to stdout.
- Model exit outcomes as explicit constants or typed values (`0`, `1`, and continue/no-direct-exit).

**Depends:** `SPEC.md` FR-002 through FR-015, Output Contract rows 109-122; `ARCHITECTURE.md` package boundary `internal/diagnostics` and errors-at-boundary rule.

**Validation:** Unit tests assert exact bytes for every exact literal/interpolation function, including final newline behavior chosen for line-oriented `console.error`/`console.warn` equivalents.

### WP-002 — Preflight diagnostic ordering

**Goal**: Implement F-007 pre-scan checks in source order: help, version, mutual exclusion, comma warning, then continue into scan.

**Implementation scope**:

- Help present: emit source help/version lines to stderr, no diagnostic text to stdout, exit `0`, skip later preflight.
- Version present without help: emit `25.0.1\n` to stderr, exit `1`, skip later preflight.
- Both `--failOn` and `--onlyAllow`: emit exactly `--failOn and --onlyAllow can not be used at the same time. Choose one or the other.` to stderr, exit `1`, skip comma warnings/scanning.
- Active policy value containing comma: emit the correct flag-specific warning to stderr and continue without direct exit-code change.

**Depends:** WP-001; F-001/F-002 parsed `Options.Help`, `Options.Version`, `Options.FailOn`, `Options.OnlyAllow`; `SPEC.md` FR-001 through FR-005 and Output Contract field order; `ARCHITECTURE.md` Composition Root steps 1-2 and flag propagation table.

**Validation:** CLI integration tests for `--help`, `--version`, `--failOn` + `--onlyAllow`, `--failOn MIT,ISC`, and `--onlyAllow MIT,ISC` assert stderr, stdout, status, and short-circuit behavior using a scanner spy or fixture that proves fatal preflight does not scan.

### WP-003 — Policy-failure diagnostics after F-005 restriction

**Goal**: Attach source-compatible fatal diagnostics to the policy stage after the restricted package set has been produced.

**Implementation scope**:

- Build non-empty trimmed semicolon token lists for `--failOn` and `--onlyAllow`; ignore empty trimmed tokens without diagnostics.
- For `--failOn`, iterate restricted packages in the order supplied by F-005/F-006; when `record.licenses` exactly equals a token, write `Found license defined by the --failOn flag: "<license>". Exiting.` to stderr and return exit `1` immediately.
- For `--onlyAllow`, when no token is contained as a substring of `record.licenses`, write `Package "<item>" is licensed under "<license>" which is not permitted by the --onlyAllow flag. Exiting.` to stderr and return exit `1` immediately.
- Ensure policy-failure diagnostics occur before renderer callback output for that run.

**Depends:** WP-001; WP-002 for mutual-exclusion precedence; F-005 restricted package records and key ordering; `SPEC.md` FR-006 through FR-009, Edge Cases, Algorithm Fidelity steps 7-9; `ARCHITECTURE.md` Composition Root step 7 and Conformance Matrix rows `--failOn`/`--onlyAllow`.

**Validation:** Unit tests over ordered restricted records for exact-match, semicolon multi-token, empty-token ignore, first-match immediate exit, onlyAllow substring acceptance, onlyAllow violation message interpolation, stdout empty, stderr exact, status `1`.

### WP-004 — Init callback-equivalent error handling

**Goal**: Preserve the CLI callback error prefix and avoid inventing a new exit-code policy for scan/init errors.

**Implementation scope**:

- When the app receives a scanner/init error equivalent to `checker.init` callback error, write `Found error` to stderr before rendering the error object.
- Render the second error-object line using the target's declared Node-compatible decision; do not replace it with JSON/YAML/error-code output.
- Do not return a new explicit non-zero exit solely because of this callback-equivalent error unless another source-compatible fatal path has already set it.
- Debug-log init errors through the error namespace when enabled, without altering functional output or exit code.

**Depends:** WP-001; WP-006 for debug namespace gate; scanner/app error propagation from F-003/F-006; `SPEC.md` FR-010/FR-011/FR-014 and Output Contract callback/debug rows; `ARCHITECTURE.md` risks for Node error-object formatting.

**Validation:** App-level tests inject scanner/init errors and assert `Found error` precedes the compatible error rendering, no stdout diagnostic copy, no newly forced status `1` solely for this condition, and debug enablement does not change status.

### WP-005 — `--files` missing-license-file warning integration

**Goal**: Preserve warning behavior emitted during file-copy side effects without stopping remaining packages.

**Implementation scope**:

- During `--files`, for each package lacking a usable `licenseFile`, write exactly `no license file found for: <moduleName>` to stderr.
- Continue processing remaining package records.
- Preserve warning order from the package key order passed to the file-copy renderer.
- Do not invent friendly filesystem error messages for other renderer side-effect failures in this slice.

**Depends:** WP-001; F-006 `--files` side-effect loop and ordered package records; `SPEC.md` FR-012, Edge Cases, Output Contract missing-file row; `ARCHITECTURE.md` Composition Root step 9.

**Validation:** Renderer/file-side-effect tests with ordered records assert one stderr warning per missing license file, exact interpolation, continued copies for later records, stdout unaffected by the warning, and no direct exit-code change from the warning alone.

### WP-006 — Debug namespace compatibility

**Goal**: Enable source-documented debug namespaces without claiming uncaptured byte formatting parity.

**Implementation scope**:

- Parse `DEBUG=license-checker*` and enable at least `license-checker:log` and `license-checker:error`.
- Wire scan-start/general debug output to the log namespace and init-error debug output to the error namespace.
- Keep debug output from changing policy decisions, stdout/stderr diagnostic contracts, or exit-code decisions.
- Document tests as namespace/enablement assertions only; do not pin timestamp/color/prefix bytes unless maintainers later provide baselines.

**Depends:** `SPEC.md` FR-014, Assumptions, Output Contract debug rows; `ARCHITECTURE.md` `internal/debuglog` decision and Conformance Matrix debug row.

**Validation:** Unit tests for DEBUG pattern matching and app-level tests verifying enabled namespace calls occur while normal F-007 stderr/stdout/status assertions remain unchanged.

### WP-007 — End-to-end CLI conformance tests for F-007

**Goal**: Prove F-007 observable compatibility at the process boundary.

**Implementation scope**:

- Build or run the Go CLI through `os/exec` in tests.
- Capture stdout, stderr, and process status for all F-007 acceptance scenarios.
- Include golden/source-derived assertions for exact strings and declared compatible exceptions for help trailing whitespace, debug formatting, and callback error object formatting.
- Keep live Node source differential execution out of the automated flow.

**Depends:** WP-001 through WP-006; F-001/F-002 parser availability; F-005/F-006 injectable fixtures for policy and files; `SPEC.md` Success Criteria SC-001 through SC-005; `ARCHITECTURE.md` Testing and Acceptance reference.

**Validation:** `go test ./...` passes; tests compare exact stderr/stdout/status for help/version/mutual-exclusion/comma warnings/failOn/onlyAllow/missing-file/callback-prefix and assert no diagnostics are emitted on the wrong channel.

## Execution Order

1. WP-001 — constants/writer/exit primitives.
2. WP-002 — preflight ordering and short-circuit tests.
3. WP-003 — policy fatal diagnostics, depending on restricted record fixtures.
4. WP-006 — debug namespace gate, so WP-004 can use it for init errors.
5. WP-004 — callback-equivalent error handling.
6. WP-005 — `--files` warning integration.
7. WP-007 — process-boundary conformance tests across all F-007 scenarios.

## Validation Approach

- **Unit tests**: exact literal/interpolation helpers, preflight state machine, token splitting/trim behavior, policy failure decisions, debug namespace matching.
- **App-level tests**: injected scanner/policy/render conditions for callback errors, missing-file warnings, and debug-enabled paths.
- **CLI integration tests**: standard-library `os/exec` process tests capturing stdout, stderr, and exit status for acceptance scenarios.
- **Golden checks**: source-derived exact strings from `SPEC.md` Output Contract are byte-compared where captured; compatible/open rows are asserted only to their documented limits.
- **Final command**: `go test ./...` from `target/` or repository-defined Go module root after implementation exists.

## Traceability Matrix

| Plan item | SPEC dependency | ARCHITECTURE dependency |
|---|---|---|
| WP-001 Diagnostic constants | FR-002-FR-015, Output Contract | `internal/diagnostics`, errors-at-boundary rule |
| WP-002 Preflight ordering | FR-001-FR-005, Algorithm Fidelity 1-5 | Composition Root steps 1-2, flag propagation help/version/failOn/onlyAllow |
| WP-003 Policy failures | FR-006-FR-009, Edge Cases, Algorithm Fidelity 7-9 | Composition Root step 7, `internal/policy`, Conformance rows `--failOn`/`--onlyAllow` |
| WP-004 Callback error | FR-010-FR-011, SA-004/SA-007 | `internal/app`, diagnostics boundary, risk: Node error-object formatting |
| WP-005 Missing file warning | FR-012, SA-005 | Composition Root step 9, `internal/render`/`files` |
| WP-006 Debug namespaces | FR-014, SA-006/SA-007 | `internal/debuglog`, Dependency Decisions `debug ^3.1.0` |
| WP-007 CLI tests | SC-001-SC-005 | Testing decision, Acceptance reference, Context7 `/golang/go` testing/os-exec guidance |

## Complexity Tracking

No constitution violations are identified from the allowed inputs. The plan intentionally avoids new frameworks, public APIs, new storage, or live Node execution in the automated flow.
