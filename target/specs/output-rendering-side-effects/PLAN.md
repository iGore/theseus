# Implementation Plan: F-006 Output Rendering and Side-Effect Outputs

**Date**: 2026-06-22 | **Spec**: `target/specs/output-rendering-side-effects/SPEC.md`  
**Input**: Feature specification from `target/specs/output-rendering-side-effects/SPEC.md` and shared architecture from `target/specs/ARCHITECTURE.md`

**Note**: This plan is scoped only to F-006. Do not implement F-003/F-004/F-005 scanning, record construction, license classification, or filtering here; consume their ordered package records exactly as supplied.

## Summary

Implement `internal/render` for the Go `license-checker` target so already-scanned and already-filtered ordered package records render to tree, JSON, CSV, Markdown, summary, `--out`, and `--files` outputs with source-compatible ordering, newline, literal, color-compatibility, and filesystem side-effect behavior. The renderer must preserve JavaScript insertion-order semantics using the architecture's ordered data model, not plain Go map iteration.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md` F-006 only
- [x] Technical context filled with concrete decisions
- [x] Work packages derived from F-006 functional requirements and architecture `internal/render` boundary
- [x] Dependencies and execution order documented with `Depends:` notes
- [x] Validation approach captured for golden snippets, newline rules, ordered output, side effects, and color compatibility

## Decision Log

- [x] Place implementation under `target/internal/render`, consuming `target/internal/ordered`, `target/internal/color`, and filesystem/writer abstractions defined by the architecture.
- [x] Use custom ordered JSON marshaling plus `encoding/json.MarshalIndent`-equivalent two-space formatting; append the CLI JSON newline in F-006 formatting logic.
- [x] Reimplement `treeify.asTree(obj, true)` locally rather than adding a Go tree dependency; lock glyphs and traversal with golden tests.
- [x] Use `os.MkdirAll`/filesystem abstraction as the `mkdirp` replacement for `--out` and `--files` recursive directory creation.
- [x] Preserve source CSV non-compliance: quote values with `"` delimiters but do not introduce general RFC CSV escaping in the renderer.
- [x] Treat no-color output as byte-identical baseline; forced-color tests verify ANSI role placement only, consistent with architecture conformance status `compatible` for color detection.

## Open Questions

- [ ] Full live byte fixtures for JSON, full tree, summary, `--files`, `--out`, and color are absent by design; validation must use committed snippets and source-derived contracts only.
- [ ] Summary equal-count order is source-unstable; do not add lexical tie sorting unless a future fixture explicitly pins it.
- [ ] Exact chalk terminal auto-detection is not byte-pinned; renderer should accept an explicit color eligibility value from upstream.

## Technical Context

**Language/Version**: Go; module rooted under `target/` as defined by `ARCHITECTURE.md`.  
**Primary Dependencies**: Go standard library (`encoding/json`, `strings`, `sort`, `io`, `io/fs`, `os`/filesystem facade, `path/filepath`, `testing`); internal packages `ordered`, `render`, `color`, `files`/filesystem abstraction. No external renderer dependency planned.  
**Storage**: Filesystem only for `--out` and `--files`; renderer input is an in-memory ordered package map.  
**Testing**: Standard `testing` package with table-driven unit tests, byte-level golden snippets from `SPEC.md`, filesystem temp-dir tests, and later CLI integration through `go test ./...`.  
**Target Platform**: POSIX-like CLI on darwin/linux first; path-joining behavior must be isolated where `filepath` can differ.  
**Project Type**: Go CLI internal rendering package plus app wiring hooks.  
**Performance Goals**: Linear traversal over ordered package records and fields; avoid unnecessary resorting except summary count sorting.  
**Constraints**: No live Node source execution; preserve exact field/key order, source literals, newline differences, no-general-CSV-escaping behavior, side-effect precedence, and no-color byte identity.  
**Scale/Scope**: F-006 only: render already-prepared package maps and perform output side effects.

### Context7 Planning Note

Context7 `/golang/go` was consulted for Go conventions relevant to this plan: `encoding/json.MarshalIndent` supports prefix/indent formatting, Go `internal/` packages enforce private implementation boundaries, `os.MkdirAll`/file APIs cover recursive directory creation and writes, and the standard `testing` package supports table-driven tests. These confirm the architecture's stdlib-first plan for F-006.

## Constitution Check

*GATE: Applies repository workflow rules provided to the Planner.*

- [x] Spec-first workflow preserved: this plan derives from `SPEC.md` and `ARCHITECTURE.md` only.
- [x] Artifact ownership preserved: Planner writes only this `PLAN.md`; no `TASKS.md`, `BUILD.md`, source code, or extra design artifacts are produced.
- [x] Go default target preserved.
- [x] F-006 scope remains bounded to rendering and side-effect outputs.
- [x] Parity verification remains in-flow via committed/source-derived golden tests; live differential execution remains manual and out-of-flow.
- [x] No-color exact baseline and declared color compatibility limitation preserved.

## Project Structure

### Documentation (this feature)

```text
target/specs/output-rendering-side-effects/
├── SPEC.md
└── PLAN.md
```

### Source Code (repository root)

```text
target/
├── cmd/
│   └── license-checker/
│       └── main.go                  # app wiring only; no renderer business logic
└── internal/
    ├── app/                         # calls render.Format then render.Emit in architecture stage order
    ├── color/                       # supplies chalk-compatible color eligibility/styling helpers
    ├── files/                       # OS filesystem adapter used by render side effects
    ├── ordered/                     # ordered map/record types required for JSON/tree/CSV traversal
    └── render/
        ├── render.go                # renderer selection and public package facade
        ├── json.go                  # ordered two-space JSON + CLI newline rule
        ├── tree.go                  # treeify-compatible tree renderer
        ├── csv.go                   # default/custom CSV rendering
        ├── markdown.go              # default/custom Markdown rendering
        ├── summary.go               # exact-license count aggregation and tree rendering
        ├── emit.go                  # files/out/stdout precedence and newline behavior
        └── *_test.go                # byte-level unit tests for F-006 contracts
```

**Structure Decision**: Use the architecture's `internal/render` package as the F-006 owner. It consumes ordered records from prior slices and exposes formatter/emit functions to `internal/app`; exported public Go APIs are not introduced unless later tasks explicitly require them.

## Complexity Tracking

No constitution/workflow violations require justification.

## Ordered Work Packages

### WP-001 — Renderer input contracts and options

**Depends:** F-003/F-004/F-005 ordered record model from `internal/ordered`/`internal/records`; `ARCHITECTURE.md` ordered-data decision.  
Define the F-006-facing input types/options for renderer selection (`JSON`, `CSV`, `Markdown`, `Summary`, `Color`, `Out`, `Files`, `CustomFormat`, `CSVComponentPrefix`) without reshaping record content. Ensure package key order and per-record field order remain observable.

Traceability: SPEC FR-001, FR-021 through FR-029; ARCH lines 22, 64-65, 104-105, 180.

### WP-002 — Treeify-compatible tree renderer

**Depends:** WP-001 ordered records.  
Implement `asTree` and `print` behavior matching `treeify.asTree(obj, true)`: box glyphs, indentation, object insertion-order traversal, and raw return string without stdout newline. Add default tree fixtures for README tree and guessed `MIT*` snippets.

Traceability: SPEC FR-003, FR-004, FR-021, FR-025; Output Contract tree snippets; ARCH dependency decision for `treeify`.

### WP-003 — Ordered JSON formatter

**Depends:** WP-001 ordered records and ordered JSON marshaler.  
Format CLI JSON as two-space ordered JSON plus exactly one `\n` before side-effect handling. Preserve top-level package order and field insertion order, including private `UNLICENSED`, repository normalization, author URL override, path/license/notice/custom fields as supplied.

Traceability: SPEC FR-005, FR-021 through FR-029; ARCH package decision `encoding/json` with custom ordered marshaler.

### WP-004 — CSV formatter

**Depends:** WP-001 custom-format order and upstream custom-field population.  
Implement default CSV headers/rows, optional component prefix, custom-format headers/rows, missing default `licenses`/`repository` as empty quoted strings, `\n` joins with no trailing newline, and source-compatible lack of general CSV escaping.

Traceability: SPEC FR-006 through FR-009, FR-027, FR-028; CSV golden snippets and edge cases; ARCH conformance rows `--csv`, `--csvComponentPrefix`, `Programmatic customFormat`.

### WP-005 — Markdown formatter

**Depends:** WP-001 custom-format order.  
Implement default Markdown rows and custom-format package/detail rows in source order. `asMarkDown` must not append a trailing newline; CLI Markdown mode must append exactly one newline before emit handling.

Traceability: SPEC FR-010 through FR-012, FR-027; Markdown golden snippets; ARCH conformance row `--markdown`.

### WP-006 — Summary renderer

**Depends:** WP-002 tree renderer and WP-001 ordered records.  
Count exact `licenses` strings, sort entries by descending count only, build an ordered summary object in that order, and render through the treeify-compatible renderer. Do not stabilize equal-count ties lexically.

Traceability: SPEC FR-013, FR-014; Sort/aggregation order; ARCH risk on summary equal-count order.

### WP-007 — Color compatibility hooks

**Depends:** WP-002 tree renderer, upstream `internal/color` eligibility.  
Apply tree key color rewriting only when `color && !out && !(csv || json || markdown)`. Split keys on `@` and use first-two-parts behavior; preserve upstream bold-red sentinel values when supplied. Keep no-color as byte-identical baseline.

Traceability: SPEC FR-019, FR-020; ARCH conformance row `--color` / implicit color support and risk note.

### WP-008 — Renderer selection and formatted-output newline rules

**Depends:** WP-002 through WP-006 and WP-007.  
Implement precedence: JSON, else CSV, else Markdown, else summary, else tree. Centralize raw formatter newline contracts so stdout and `--out` receive the exact strings required per mode.

Traceability: SPEC FR-001, FR-005, FR-009, FR-012, whitespace contract lines 242-248; ARCH ordered stage sequence step 8.

### WP-009 — Emit side effects: `--files`, `--out`, stdout

**Depends:** WP-008 formatted output and filesystem abstraction.  
Implement side-effect precedence: `--files` bypasses formatted output; else `--out` recursively creates parent directory and writes raw UTF-8 formatted string; else stdout writes using line-emitting behavior that appends its own line ending. `asFiles` must create output dirs, iterate package key order, copy `<moduleName>-LICENSE.txt`, and emit exact missing-file warning text.

Traceability: SPEC FR-002, FR-015 through FR-018; `--files`/`--out` acceptance scenarios; ARCH stage sequence step 9 and `mkdirp` dependency decision.

### WP-010 — Validation fixtures and tests

**Depends:** WP-002 through WP-009.  
Create byte-level tests from the snippets and source-derived contracts: default tree, guessed license marker, JSON indentation/order/newline, CSV default/custom/component/missing fields/no-escaping, Markdown default/custom/newline, summary descending counts/no tie assertion, `--out` raw file contents, `--files` copied content/naming/warning, stdout-vs-file newline matrix, and forced-color role placement.

Traceability: SPEC Success Criteria SC-001 through SC-007; ARCH acceptance reference and conformance matrix rows for F-006.

### WP-011 — App integration handoff

**Depends:** WP-008 and WP-009; earlier pipeline slices provide records/options.  
Wire `internal/app` to call `render.Format(records, options)` then `render.Emit(records, formatted, sideEffects, stdout, stderr, fs)` at the architecture-defined stage order, without adding F-007 exit-code logic or F-001 preflight behavior inside `render`.

Traceability: SPEC FR-001, FR-002; ARCH composition root and ordered stage sequence steps 8-10.

## Dependency Order

```text
WP-001
├── WP-002 ─┬── WP-006 ┐
├── WP-003 ┤          │
├── WP-004 ┤          ├── WP-008 ─── WP-009 ─── WP-011
├── WP-005 ┘          │              │
└── WP-007 ───────────┘              └── WP-010
```

Implementation should begin with ordered input contracts and formatter unit tests, then side effects, then app wiring. Do not start `--files`/`--out` integration until formatter newline contracts are locked.

## Validation Approach

- Golden output: assert byte identity for every snippet embedded in `SPEC.md`, including tree glyphs, guessed `MIT*`, CSV headers/rows, Markdown rows, `foo-LICENSE.txt`, and missing-file warning text.
- Ordered output: tests must fail if plain Go map iteration changes top-level package order or per-record field order in JSON/tree/CSV/Markdown.
- Newline matrix: test raw formatter output, `--out` file content, and stdout line-emitting behavior separately for tree, JSON, CSV, Markdown, and summary.
- Side effects: use temp dirs and fake filesystem/writers where practical; verify recursive directory creation, UTF-8 raw writes, copied license content, key-order iteration, nested module-name path behavior, and warning continuation.
- Color compatibility: default no-color tests are byte-exact; forced-color tests assert only documented blue/dim/green key segmentation and bold-red sentinel preservation.
- Summary: verify descending count sort and exact-license-string aggregation; avoid equal-count deterministic assertions unless future fixtures pin them.
- Final command for Builder phase: `go test ./...` from `target/`, with no live Node source execution.

## Guardrails for Downstream Tasks

- Do not use `map[string]any` at renderer boundaries where order is observable.
- Do not replace CSV behavior with `encoding/csv` if it changes source-compatible no-escaping output.
- Do not add trailing newlines inside `asCSV`, `asMarkDown`, `asTree`, or `asSummary` beyond the CLI-specific rules.
- Do not emit formatted output when `--files` is set.
- Do not invent live-source fixtures or claim full live byte parity beyond committed/source-derived contracts.
