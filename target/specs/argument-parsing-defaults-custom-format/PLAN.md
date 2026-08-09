# Implementation Plan: Argument parsing, defaults, and custom format loading

**Date**: 2026-06-22 | **Spec**: `target/specs/argument-parsing-defaults-custom-format/SPEC.md`  
**Input**: Feature specification from `target/specs/argument-parsing-defaults-custom-format/SPEC.md` plus shared architecture `target/specs/ARCHITECTURE.md`

**Scope**: F-002 only. Do not implement CLI process lifecycle/help/version exits, dependency scanning, license detection, filtering, final renderers, or side effects except where F-002 exposes normalized options and custom-format contracts to those later stages.

## Summary

Implement the Go target's source-compatible option normalization and custom-format loading for `davglass/license-checker` v25.0.1. The work lives primarily in `target/internal/cli` and `target/internal/customformat`, with a small app-boundary integration point that attaches `CustomFormat` before scanner options are built. The plan preserves the documented `nopt ^4.0.1` parity risks, especially unknown/malformed flags, `--no-` cooked-argument behavior, path coercion, and metadata removal.

Context7 check: `/golang/go` confirms the architecture's Go layout assumptions: `cmd/.../package main`, `internal/` visibility boundaries, stdlib-first testing with `go test`, and avoiding public APIs for private migration internals. Architecture already rejects Go `flag` as the primary parser because F-002 requires `nopt`-compatible semantics.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md`
- [x] Technical context filled with concrete decisions or `NEEDS CLARIFICATION`
- [x] Work packages derived from F-002 only
- [x] Dependencies and execution order documented
- [x] Validation approach captured

## Decision Log

- [x] Implement parser/defaults in `target/internal/cli`; use a source-compatible custom parser rather than `flag.FlagSet` as primary parser. Depends on SPEC FR-001 through FR-011 and ARCH lines 18-20, 57-60, 97, 163, 175.
- [x] Represent unbounded direct depth with an explicit Go sentinel rather than a floating `Infinity`, while preserving scanner-boundary semantics. Depends on SPEC FR-011/FR-013 and ARCH flag propagation row `--direct`.
- [x] Implement custom JSON parsing/loading in `target/internal/customformat`, using `os.ReadFile` plus `encoding/json` while preserving error-as-value behavior. Depends on SPEC FR-012, FR-014, FR-015 and ARCH lines 99, 164-165, 175.
- [x] Preserve custom-format key order with ordered data structures at the record/render boundary; do not rely on Go map iteration for observable field order. Depends on SPEC FR-016 through FR-021 and ARCH lines 22, 48, 61, 164, 240.
- [x] Keep parser edge behavior not evidenced by source/tests as explicit conformance risk, not guessed behavior. Depends on SPEC OQ-001/OQ-002 and ARCH conformance rows for `--customPath`, programmatic `customFormat`, `--direct`, and `--color`.

## Open Questions

- [ ] `NEEDS CLARIFICATION`: Exact `nopt ^4.0.1` behavior for unknown flags, malformed typed values, all path coercions, and full `--no-<flag>` cases is not exhaustively captured. Builder must implement fixture-backed known behavior and record any non-identical rows before release. [SPEC OQ-001, SC-006]
- [ ] `NEEDS CLARIFICATION`: If `parseJson` returns an Error and `checker.init` assigns it to `options.customFormat`, downstream object-key behavior has no dedicated golden output. Keep as parity-risk fixture and do not invent diagnostics. [SPEC OQ-003]

## Technical Context

**Language/Version**: Go; exact Go version inherited from target `go.mod` when created by Builder.  
**Primary Dependencies**: Go standard library for this slice: `os`, `encoding/json`, `errors`, `fmt`, `testing`; no external CLI framework.  
**Storage**: Filesystem only for `customPath` UTF-8 JSON reads; no database.  
**Testing**: Standard `testing` package, table-driven unit tests, fixtures under target testdata as needed, run with `go test ./...`.  
**Target Platform**: POSIX-like CLI on darwin/linux first; path behavior must be explicit if Windows support is tested.  
**Project Type**: Go CLI with private internal packages.  
**Performance Goals**: No special throughput target; parser/default/custom-format work should be linear in argument count and custom-format JSON size.  
**Constraints**: No live source-system execution in-flow; preserve exact source literals, field presence/order, `argv` metadata deletion, and source-compatible error-as-value behavior.  
**Scale/Scope**: One command, all known F-002 option names, direct `customFormat` and `customPath` loading, scanner-boundary depth propagation.

## Constitution Check

GATE status: Pass for planning. This plan follows repository workflow rules by staying spec-first, targeting Go, writing only `PLAN.md` in the assigned use-case folder, and not producing `TASKS.md`, `BUILD.md`, code, or additional artifacts. No live source execution is required or planned.

## Project Structure

### Documentation (this feature)

```text
target/specs/argument-parsing-defaults-custom-format/
├── PLAN.md              # This file
├── SPEC.md              # Planner input
└── TASKS.md             # Produced later by Task Decomposer, not by Planner

target/specs/
└── ARCHITECTURE.md      # Shared Planner input
```

### Source Code (repository root)

```text
target/
├── go.mod
├── cmd/
│   └── license-checker/
│       └── main.go
└── internal/
    ├── app/             # attach customFormat before scanner options; pass depth from normalized DirectDepth
    ├── cli/             # F-002 parser, known options, aliases, defaults, has(), metadata cleanup
    ├── color/           # supportsColor-compatible detector consumed by cli defaults
    ├── customformat/    # parseJson/customPath loader and ordered custom-format representation
    ├── npmgraph/        # scanner boundary consuming Start and DirectDepth; implementation owned by F-003
    ├── ordered/         # ordered object/record primitives for observable key order
    ├── records/         # custom-format field inclusion/default/suppression contract consumed by F-003/F-006
    └── render/          # downstream custom CSV/Markdown field-order consumers owned by F-006
```

**Structure Decision**: Use the shared architecture's standard Go layout under `target/`. F-002 implementation work is concentrated in `internal/cli` and `internal/customformat`, with small integration surfaces in `internal/app`, `internal/ordered`, and `internal/records` only to expose the normalized options and custom-format contract.

## Work Packages

### WP-001 — Define F-002 option model and sentinels

**Depends:** None.  
**Spec trace:** FR-001, FR-002, Key Entities `ParsedOptions`, `ScannerOptions boundary`; Output Contract default literals.  
**Architecture trace:** `internal/cli` boundary, flag-to-stage propagation table, `--direct` row.  
**Work:** Define the Go `Options` shape for every known F-002 option name, including booleans, strings/path strings, `CustomFormat`, explicit color state, and a scanner-depth representation with values `0` and unbounded sentinel. Include a way to preserve whether color was explicitly provided when needed by later color/render rules.  
**Validation:** Compile-time use by parser/default tests; unit tests assert all known option names can be represented without `argv` metadata.

### WP-002 — Implement nopt-compatible parser core for known F-002 flags

**Depends:** WP-001.  
**Spec trace:** FR-001 through FR-006, Algorithm Fidelity steps 1-2, Edge Cases for `--no-<flag>`, OQ-001.  
**Architecture trace:** Package decision rejecting `flag` as primary parser; dependency decision for `nopt ^4.0.1`.  
**Work:** Parse `license-checker [flags]` argv without subcommands; recognize all known F-002 option names, path/string/boolean types, `-v`/`-h` aliases, and a cooked-argument presence tracker for `has(a)` that returns true for both `--<a>` and `--no-<a>`. Ensure parser metadata equivalent to `argv` is not exposed by cleaned results.  
**Validation:** Table tests for every known option, short aliases, `--color`/`--no-color` presence, and absence of `argv`. Add tests only from source-derived behavior; mark malformed/unknown gaps as conformance risks instead of guessing.

### WP-003 — Implement defaults normalization

**Depends:** WP-001, WP-002.  
**Spec trace:** FR-007 through FR-011, Algorithm Fidelity steps 3-9, SC-002.  
**Architecture trace:** Composition stage 1 `cli.Parse(argv, cwd, color.Supports)`.  
**Work:** Apply defaults in source order: parse current process args when input is undefined at the API boundary; set color from `color.Supports` only when unset; force color false for `json`, `markdown`, or `csv`; default `start` to cwd when falsey; boolean-coerce `relativeLicensePath`; normalize truthy `direct` to depth `0` and falsey/missing to unbounded sentinel.  
**Validation:** Table tests for missing/falsey/truthy values, forced color for output modes, cwd default, direct `true -> 0`, direct missing/falsey -> unbounded.

### WP-004 — Implement custom JSON parse/load behavior

**Depends:** WP-001.  
**Spec trace:** FR-012, FR-014, FR-015, Custom JSON loading rule table, SC-003.  
**Architecture trace:** Composition stage 3 `customformat.Load(options.CustomPath)`; stdlib filesystem/json decisions.  
**Work:** Implement `parseJson`-equivalent behavior: non-string path returns an Error value with message `did not specify a path`; string path reads UTF-8 JSON and returns parsed structured data on success; file-read or JSON-parse errors are returned as values, not thrown/panicked. Integrate loader so `options.CustomFormat` is assigned from `CustomPath` before scanner options are built.  
**Validation:** Tests for valid JSON containing `licenseModified: "no"`, missing file, invalid JSON, and non-string path. Assert errors are values and no diagnostics are invented.

### WP-005 — Preserve custom-format ordered field contract

**Depends:** WP-001, WP-004; downstream record construction may be completed with F-003/F-006 but contract belongs here.  
**Spec trace:** FR-016 through FR-021, Custom-format field rule table, Output Contract custom-format field order/presence, SC-004, SC-007.  
**Architecture trace:** `internal/ordered`, `internal/records`, ordered data risk.  
**Work:** Provide an ordered custom-format representation and helper rules: include property when custom format absent or value is not exactly `false`; iterate keys in source-equivalent JSON/object key order; use package truthy string value when present; use configured default when missing/falsey; suppress exact `false`; preserve source branch where truthy non-string package values neither assign the non-string nor fall back to default.  
**Validation:** Unit tests over a small package metadata node for field order `name`, `description`, `pewpew`, string override, missing-field default, exact-false suppression, and truthy non-string no-default branch. Downstream CSV byte snippets remain F-006 validation consumers but must use this field order.

### WP-006 — Wire scanner-boundary inputs for F-003 handoff

**Depends:** WP-003, WP-004.  
**Spec trace:** FR-012, FR-013, User Story 2, ScannerOptions boundary.  
**Architecture trace:** Composition stages 3-4; flag propagation rows `--start`, `--customPath`, `--direct`.  
**Work:** Ensure `internal/app` or equivalent composition attaches `CustomFormat` before scanning and passes `options.Start` plus scanner options where `Depth` comes directly from normalized `DirectDepth`. Do not implement scanner internals in this slice.  
**Validation:** Boundary tests with a fake scanner asserting default/provided `start`, direct-present depth `0`, direct-absent unbounded depth, and `CustomFormat` assignment before scanner invocation.

### WP-007 — Conformance matrix and parity-risk evidence for F-002 rows

**Depends:** WP-002 through WP-006.  
**Spec trace:** OQ-001 through OQ-003, SC-005, SC-006.  
**Architecture trace:** Conformance Matrix rows `--customPath`, programmatic `customFormat`, `--direct`, `--color`; Risks for missing live byte fixtures and Go map ordering.  
**Work:** Record test evidence and any explicit non-identical parser edge behavior before the rewrite is marked complete. Keep custom-path parse-error downstream behavior as a known risk unless maintainers provide a baseline.  
**Validation:** `go test ./...` passes; F-002 conformance notes identify all unresolved `nopt` parity gaps without live source execution.

## Execution Order

1. WP-001 option model and sentinels.
2. WP-002 parser core and metadata cleanup.
3. WP-003 defaults normalization.
4. WP-004 custom JSON parse/load.
5. WP-005 ordered custom-format field contract.
6. WP-006 scanner-boundary wiring.
7. WP-007 conformance evidence and remaining parity-risk documentation.

## Validation Approach

- Unit-test `internal/cli` with table-driven cases for every known option, `-v`, `-h`, `--color`, `--no-color`, `has(a)`, no `argv`, and defaults ordering.
- Unit-test `internal/customformat` for valid/missing/invalid/non-string JSON path behavior and error-as-value semantics.
- Unit-test ordered custom-format application independent of scanner/renderers using a minimal metadata node.
- Boundary-test app/scanner integration with fakes; assert custom format is attached before scan and `Depth` equals normalized direct sentinel.
- Run `go test ./...` from `target/` once implementation exists.
- Do not use live source execution for new parser-edge captures; rely on committed snippets/source-derived contracts and document unverified `nopt` edges.

## Complexity Tracking

No constitution/workflow violations are required. The custom parser is justified by source compatibility: Go `flag`/Cobra/urfave would diverge from `nopt` aliases, `--no-` handling, path coercion, metadata, and malformed-value risks.
