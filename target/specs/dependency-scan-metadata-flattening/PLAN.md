# Implementation Plan: F-003 Dependency Scan and Metadata Flattening

**Date**: 2026-06-22 | **Spec**: [`SPEC.md`](./SPEC.md)  
**Input**: Feature specification from `target/specs/dependency-scan-metadata-flattening/SPEC.md` and shared architecture from `target/specs/ARCHITECTURE.md`

**Scope**: Plan for Go target implementation of F-003 only. This plan does not include F-001/F-002 parsing, F-004 license classification/file extraction, F-005 policy filtering, F-006 rendering, or F-007 diagnostics except where F-003 must preserve data contracts consumed by those slices.

## Summary

Implement the Go scan/flatten slice that receives parsed options and optional custom format, calls a `read-installed ~4.0.3`-compatible npm dependency-tree provider, flattens package nodes into ordered records keyed by `name@version`, applies production/development pruning, repository normalization, metadata/path/custom field population, duplicate/incomplete-node handling, lexical key sorting, and F-003 sorted-pass license sentinels. The implementation belongs primarily in `target/internal/npmgraph`, `target/internal/records`, and shared ordered data support under `target/internal/ordered`, wired later through `internal/app` according to the shared architecture.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md` for F-003 only.
- [x] Technical context filled with concrete decisions from `ARCHITECTURE.md`; no blocking clarification remains for planning.
- [x] Work packages derived from F-003 functional requirements FR-001 through FR-024.
- [x] Dependencies and execution order documented with explicit `Depends:` notes.
- [x] Validation approach captured with source-derived fixtures and no live source-system execution.

## Decision Log

- [x] Use Go under `target/` with `internal/npmgraph` for `read-installed` compatibility and `internal/records` for flattening, matching `ARCHITECTURE.md` package boundaries.
- [x] Preserve observable JavaScript insertion order by using explicit ordered object/map types; do not expose plain Go map iteration at record or renderer boundaries.
- [x] Sort top-level package keys with exact lexical string ordering equivalent to `Object.keys(data).sort()`.
- [x] Treat `read-installed` behavior as a compatibility surface, not as a naive filesystem walk; fixtures must cover `extraneous`, `root`, nested dependencies, paths, and dedupe/duplicate key behavior.
- [x] Keep license classifier/file fields reserved for F-004, but leave record ordering stable so F-004 fields can append after F-003-owned fields.
- [x] Keep `--unknown` split explicit: F-003 adds `dependencyPath` and rewrites guessed `*` license strings to `UNKNOWN` during sorted pass; F-005 owns `--onlyunknown` filtering.

## Open Questions

- [x] No planning blocker. The author URL overwrite lacks committed source fixture coverage, but `SPEC.md` requires a synthetic fixture and treats behavior as mandatory.
- [x] No live source-output baselines are available or required in-flow; final live differential execution remains a maintainer-owned manual step.

## Technical Context

**Language/Version**: Go module rooted under `target/`; exact Go version is not pinned in the artifacts, so implementation should use the repository-selected Go toolchain while avoiding version-specific APIs not needed by this slice.  
**Primary Dependencies**: Standard library for `context`, `os`, `io/fs`, `path/filepath`, `strings`, `sort`, and `testing`; no external package is required for F-003 itself.  
**Storage**: Filesystem-backed npm project tree under `Options.Start`; no database.  
**Testing**: Standard `testing` package, table-driven unit tests, source-derived testdata fixtures, and eventual `go test ./...`.  
**Target Platform**: POSIX-like CLI behavior on darwin/linux first; path tests must make separator/root behavior explicit.  
**Project Type**: Go CLI with private internal packages.  
**Performance Goals**: Deterministic traversal and sorted output for installed npm dependency trees; no numeric throughput target is specified.  
**Constraints**: No live source-system execution in-flow; byte/function parity must be based on committed snippets and source-derived contracts; ordered fields and exact literals are observable downstream.  
**Scale/Scope**: Bounded to F-003 package-record map generation for the pinned `davglass/license-checker` 25.0.1 behavior.

### Context7 Planning Note

Context7 `/golang/go` was consulted for Go implementation conventions relevant to this plan: use `package main` for command entry points, keep private implementation in `internal/`, rely on standard `go test` testing workflows, use standard sorting primitives for lexical ordering, and avoid unspecified map traversal where output order is observable. These reinforce the shared architecture decisions rather than changing the F-003 scope.

## Constitution Check

No `/memory/constitution.md` file and no `.specify/extensions.yml` file were found in the repository root. Applicable repository workflow gates from `AGENTS.md` are satisfied: this plan is spec-first, Go-targeted, scoped to one use-case folder, owns only `PLAN.md`, does not write `TASKS.md`/`BUILD.md`/code, and keeps live source-system execution out of the in-flow validation plan.

## Project Structure

### Documentation (this feature)

```text
target/specs/dependency-scan-metadata-flattening/
├── SPEC.md              # Existing F-003 specification
└── PLAN.md              # This Planner-owned artifact
```

No `research.md`, `data-model.md`, `quickstart.md`, or `contracts/` artifacts are required for this slice because the shared architecture and F-003 specification already resolve the planning decisions and no external API contract is introduced by F-003.

### Source Code (repository root)

```text
target/
├── go.mod
├── cmd/
│   └── license-checker/
│       └── main.go
└── internal/
    ├── app/             # Later orchestration wiring; calls scanner then flattener
    ├── cli/             # F-001/F-002 parsed Options consumed by this slice
    ├── debuglog/        # Provider log option source
    ├── npmgraph/        # F-003 read-installed-compatible scanner and Node model
    ├── ordered/         # Ordered object/map primitives for records and sorted maps
    ├── records/         # F-003 flattening, metadata, ordering, sentinel pass
    └── license/         # F-004 integration point; must append fields after F-003 fields
```

**Structure Decision**: Implement F-003 in `internal/npmgraph`, `internal/records`, and `internal/ordered`, with `internal/app` wiring left as integration work after F-001/F-002 option structs exist. This follows `ARCHITECTURE.md` lines 55-67 and the ordered stage sequence lines 95-106.

## Work Packages

### WP-01 — Define F-003 data contracts and option mapping

**Depends:** F-001/F-002 option model shape or a temporary internal test fixture matching `SPEC.md` FR-001.  
**Artifacts:** `SPEC.md` FR-001, FR-002, FR-003, FR-004; `ARCHITECTURE.md` flag-to-stage propagation table.  
**Work:**

- Define/confirm `npmgraph.Options{Dev, Depth, Log}` and `records.Options{Production, Development, Unknown, CustomFormat, BasePath/Start, Include fields}`.
- Ensure `Dev` is `true` by default and `false` when `Production` or `Development` is true.
- Pass `DirectDepth` unchanged to provider `Depth`, including `0` for `--direct` and unbounded sentinel when absent.
- Keep `relativeLicensePath`, `color`, and license-specific options out of F-003 logic except where options must pass through for later slices.

### WP-02 — Implement ordered record/map primitives required by F-003

**Depends:** WP-01 field contract.  
**Artifacts:** `SPEC.md` Output Contract field-order table; `ARCHITECTURE.md` ordered data decision and risk.  
**Work:**

- Provide an ordered record type that supports insert-if-new, overwrite-without-moving, deletion, key iteration, and later custom JSON/tree/CSV traversal.
- Provide an ordered package map type whose top-level package order is explicitly set from sorted keys.
- Make field-order tests possible without relying on Go map iteration.
- Reserve compatibility for F-004 fields appended after `path` without reordering F-003 fields.

### WP-03 — Build `npmgraph` scanner compatibility layer

**Depends:** WP-01; can be developed with fixtures before full app wiring.  
**Artifacts:** `SPEC.md` FR-002, FR-003, FR-004, FR-024 and Algorithm Fidelity dependency-provider contract; `ARCHITECTURE.md` dependency decision for `read-installed`.  
**Work:**

- Model npm dependency nodes with fields required by F-003: `name`, `version`, `private`, `repository`, `url`, `author`, `path`, `dependencies`, `extraneous`, and `root`; keep room for downstream `license`, `licenses`, and `readme` fields.
- Implement filesystem scanning to reproduce logical installed npm tree semantics required by `read-installed ~4.0.3`, including nested dependency shape and absolute `path` propagation under `start`.
- Avoid substituting a simple recursive directory walk unless fixture tests demonstrate compatibility for dedupe, nested dependencies, `extraneous`, `root`, and path propagation.
- Route debug logging through a non-behavioral `Log` option; debug formatting is not owned by F-003.

### WP-04 — Implement flatten recursion and pruning

**Depends:** WP-02 for ordered records; WP-03 node model or fixtures.  
**Artifacts:** `SPEC.md` FR-005, FR-006, FR-007, FR-008, FR-009, FR-010, FR-017, FR-018; pruning table.  
**Work:**

- Construct keys as exact `name + "@" + version` string concatenation.
- Short-circuit immediately when key already exists, including recursion avoidance for duplicate/circular graph cases.
- Apply production pruning before insertion when `Production && node.Extraneous`.
- Apply development pruning before insertion when `Development && !node.Extraneous && !node.Root`.
- Initialize every new record with `licenses = UNKNOWN`, then optional `private = true` immediately after licenses.
- Delete incomplete-node records when `name` or `version` is missing before returning flattened data.
- Recurse through dependencies carrying F-003 option context unchanged.

### WP-05 — Implement metadata population and repository normalization

**Depends:** WP-04 record creation and ordered field insertion.  
**Artifacts:** `SPEC.md` FR-011 through FR-016, repository normalization table, field-order table.  
**Work:**

- Normalize `repository.url` only under the source conditions and in the exact ordered steps: `git+ssh://git@` → `git://`, `git+https://github.com` → `https://github.com`, `git://github.com` → `https://github.com`, `git@github.com:` → `https://github.com/`, trailing `.git` removal, otherwise preserve.
- Insert `url` from `json.url.web` when eligible, then overwrite it from `author.url` when eligible without moving the original insertion position.
- Insert `publisher` and `email` from `author.name` and `author.email` in source order.
- Add `dependencyPath` from node path only when `Unknown` is true.
- Iterate custom format keys in source order; exclude keys whose configured value is `false`; copy string package values; use configured default only when source value is absent/falsey; avoid defaulting non-string truthy values.
- Add/overwrite `path` only when requested and node path is a string, preserving earlier custom-format insertion position if custom format inserted `path` first.

### WP-06 — Implement sorted pass and F-003 license sentinels

**Depends:** WP-04/WP-05 flattened records.  
**Artifacts:** `SPEC.md` FR-019, FR-020, FR-021, FR-022; Sentinels & literals; `ARCHITECTURE.md` ordered stage sequence.  
**Work:**

- Build the returned package map by sorting exact key strings lexically before F-005 filtering.
- During the sorted pass, overwrite private records' `licenses` with `UNLICENSED` without changing field order.
- During the sorted pass, overwrite falsey `licenses` with `UNKNOWN` without changing field order.
- When `Unknown` is true, rewrite any non-`UNKNOWN` license string containing `*` to `UNKNOWN`; do not implement broader F-004 license classification or F-005 filtering here.
- Keep colorized sentinel bytes outside F-003; F-006 owns renderer/color behavior.

### WP-07 — Integrate scanner and flattener through app seams

**Depends:** WP-01 through WP-06 and availability of F-001/F-002 parsed options.  
**Artifacts:** `ARCHITECTURE.md` composition root and ordered stage sequence; `SPEC.md` User Story 1 provider-call acceptance.  
**Work:**

- Wire `app` stage order so F-003 runs after CLI/custom-format loading and before F-004/F-005/F-006.
- Pass `options.Start` to the scanner and pass the scanner root node to `records.Flatten`.
- Preserve returned ordered map for downstream license classification, policy, and render packages.
- Keep errors as values for later diagnostics conversion; do not emit stdout/stderr from F-003.

### WP-08 — Build F-003 validation fixtures and tests

**Depends:** WP-02 through WP-07 as relevant; test fixtures can start in parallel with implementation contracts.  
**Artifacts:** `SPEC.md` Success Criteria SC-001 through SC-006 and Edge Cases; `ARCHITECTURE.md` testing decision and risks.  
**Work:**

- Add unit tests for every repository normalization row, including fallback and `abbrev@1.0.9` expected repository value.
- Add field insertion-order tests for `licenses`, `private`, `repository`, `url`, `publisher`, `email`, author URL overwrite, `dependencyPath`, custom-format fields, and `path`.
- Add flattening tests for duplicate keys/circular guard, missing name/version deletion, private `UNLICENSED`, falsey license `UNKNOWN`, and `--unknown` guessed `*` rewrite.
- Add scan-scope tests for default, production, development, and direct depth provider options and prune outcomes.
- Add compatibility fixtures for `read-installed`-risk behavior: nested dependencies, `extraneous`, `root`, absolute paths under `start`, and dedupe behavior.
- Use committed examples and source-derived fixtures only; do not run the live Node source tool during tests.

## Dependencies and Execution Order

1. WP-01 first establishes option/provider contracts needed by all later work.
2. WP-02 should follow immediately because every observable record package depends on ordered mutation semantics.
3. WP-03 and WP-04 may proceed in parallel after WP-01/WP-02 if WP-04 uses synthetic node fixtures while WP-03 implements filesystem scanning.
4. WP-05 depends on WP-04 because metadata rules attach to records created by flatten recursion.
5. WP-06 depends on completed record population and must run after all F-003 metadata and any F-004 license hook has had a chance to set `licenses` in the eventual full pipeline.
6. WP-07 depends on the scanner and flattener APIs being stable and on parsed options from upstream slices.
7. WP-08 is continuous but cannot be complete until each corresponding work package is implemented.

## Validation Approach

- Run package-level unit tests with `go test ./...` from `target/` once implementation exists.
- Prefer table-driven tests for normalization, pruning, field order, sentinel rewriting, duplicate handling, and custom-format behavior.
- Validate top-level sorted order by comparing explicit key sequences, not map stringification.
- Validate ordered record fields by inspecting ordered key slices and value overwrites separately.
- Validate `read-installed` compatibility through committed synthetic fixture trees and npm-project testdata; include `extraneous`, `root`, nested dependencies, duplicate package identity, and absolute path propagation.
- Validate `abbrev@1.0.9` source-observed compatibility by asserting key presence and repository `https://github.com/isaacs/abbrev-js`; downstream CSV/Markdown byte rendering remains F-006.
- Do not add live source-system differential tests to the automated flow.

## Compatibility Risks to Preserve

- **Ordered records:** Go maps cannot represent source JavaScript insertion-order semantics. All F-003 records and sorted maps must use ordered structures, especially where custom-format keys overwrite earlier fields.
- **`read-installed` semantics:** The provider must model logical npm installation metadata, not only directory layout. `extraneous`, `root`, nested `dependencies`, `path`, and duplicate/dedupe behavior are high-risk parity areas.
- **Author URL overwrite:** `author.url` must overwrite `json.url.web` while preserving insertion position; this is mandatory despite lacking a committed source fixture.
- **Incomplete nodes:** Source creates a temporary key from missing values but deletes the record when name/version is absent; tests must not simplify this into pre-validation that changes side effects visible to recursion/order.
- **Sentinel ownership:** F-003 owns initial `UNKNOWN`, private `UNLICENSED`, falsey-license `UNKNOWN`, and `--unknown` guessed-license rewrite; F-004 owns classifier values and F-006 owns color bytes.

## Traceability Matrix

| Plan item | SPEC coverage | ARCHITECTURE coverage |
|---|---|---|
| WP-01 option/provider mapping | FR-001 to FR-004; provider option table | Stage sequence step 4; flag propagation rows for production/development/direct/unknown/start |
| WP-02 ordered primitives | Output Contract field order; SC-003 | Ordered data decision; risk “Go map ordering would break parity” |
| WP-03 scanner | FR-002, FR-024; Algorithm Fidelity provider contract | `internal/npmgraph`; dependency decision for `read-installed` |
| WP-04 flatten recursion/pruning | FR-005 to FR-010, FR-017, FR-018; pruning table | `internal/records`; implementation order step 6 |
| WP-05 metadata/normalization | FR-011 to FR-016; repository normalization table | `internal/records`; ordered data constraint |
| WP-06 sorted sentinels | FR-019 to FR-022; sentinels table | stage sequence before policy/render; conformance rows for private/unknown |
| WP-07 app integration | User Story 1; whitespace/no-output contract | composition root and ordered stage sequence |
| WP-08 validation | SC-001 to SC-006; Edge Cases | testing decision; risks and conformance matrix |

## Complexity Tracking

No constitution or workflow violations require justification. The ordered-data and scanner abstractions are not optional complexity: they are required by the F-003 output contract and the shared architecture risks.
