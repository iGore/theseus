# Implementation Plan: F-004 License Detection, Classification, and License File Precedence

**Date**: 2026-06-22 | **Spec**: `target/specs/license-detection-classification-files/SPEC.md`  
**Input**: Feature specification from `target/specs/license-detection-classification-files/SPEC.md`; shared architecture from `target/specs/ARCHITECTURE.md`

## Summary

Implement only F-004 in Go: exact license classification, SPDX pass-through via the architecture-selected SPDX adapter, license-file precedence, and package-record license/readme/file-field application for records supplied by F-003. The implementation lives in `target/internal/license` and uses `target/internal/spdxcompat` as the only SPDX boundary, preserving source-compatible sentinels, raw valid SPDX input strings, ordered regex behavior, selected license-file paths/text, NOTICE detection, and downstream record field contracts from the SPEC.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md`: F-004 only; no F-003 scanning, F-005 policy, F-006 rendering, code, `TASKS.md`, or `BUILD.md`.
- [x] Technical context filled with concrete decisions or `NEEDS CLARIFICATION`.
- [x] Work packages derived from the assigned functionality scope.
- [x] Dependencies and execution order documented.
- [x] Validation approach captured.

## Decision Log

- [x] Put classifier, license-file precedence, and record application logic in `target/internal/license`, matching ARCHITECTURE lines 61-63 and 190-191.
- [x] Use `target/internal/spdxcompat` facade backed by `github.com/github/go-spdx/v2/spdxexp v2.7.0` as specified by ARCHITECTURE lines 23 and 169; classifier must return the original raw input unchanged on successful parse.
- [x] Preserve JavaScript source behavior even when it looks suboptimal: remove only the first newline, evaluate regexes in table order, allow BSD source-code rule to be shadowed, and push only the first filename per precedence pattern.
- [x] Treat F-003 as provider of nodes/records/path/readme/package metadata and F-006 as owner of renderer formatting; F-004 only mutates/sets record fields described in SPEC lines 121-128.
- [x] Use standard Go `testing` table-driven unit tests and `go test ./...`; Context7 `/golang/go` confirms internal package visibility and Go test workflow, and Context7 `/git-pkgs/spdx` confirms strict parsing/normalization differences that justify keeping SPDX behavior behind an adapter.

## Open Questions

- [ ] `NEEDS CLARIFICATION`: final verification of `github.com/github/go-spdx/v2/spdxexp v2.7.0` API details is deferred to implementation because the shared architecture mandates it but Context7 returned richer docs for `/git-pkgs/spdx`, not the exact selected module. Do not replace the architecture-selected dependency without Architect approval.
- [ ] `NEEDS CLARIFICATION`: exact Go record/container type names from F-003 are not available in this slice yet. Plan against an interface boundary and wire to concrete F-003 types when that slice exists.

## Technical Context

**Language/Version**: Go, module rooted under `target/`; exact Go toolchain version `NEEDS CLARIFICATION` from future `target/go.mod`.  
**Primary Dependencies**: Go stdlib `regexp`, `strings`, `path/filepath`, `io/fs`, `testing`; `internal/spdxcompat` over `github.com/github/go-spdx/v2/spdxexp v2.7.0`; F-003 record/node interfaces.  
**Storage**: Filesystem reads for selected license/notice files only; no database.  
**Testing**: Standard `go test ./...`, table-driven unit tests for classifier and file precedence, integration-style record-application tests with constructed F-003-compatible fixtures.  
**Target Platform**: POSIX-like CLI on darwin/linux first; path separator behavior explicit where `filepath` can vary.  
**Project Type**: Go CLI internals under `target/internal/license`; not a public Go API.  
**Performance Goals**: Linear over number of candidate filenames plus license/readme text length; no full-tree traversal beyond F-003-provided package context.  
**Constraints**: No live source execution; exact literals/sentinels; exact rule order; raw SPDX pass-through; no renderer formatting in this slice; no policy/filter behavior in this slice.  
**Scale/Scope**: One package record at a time during flattening; supports all package records discovered by F-003.

## Constitution Check

No `/memory/constitution.md` or `.specify/extensions.yml` was found. Applicable repository workflow gates:

- [x] Planner owns only `PLAN.md` in this use-case folder.
- [x] Plan remains spec-first and bounded to F-004.
- [x] Go is the default target language.
- [x] Context7 was used for Go/internal testing conventions and SPDX parser planning.
- [x] No live source-system execution is planned or required.
- [x] No design-side artifacts are needed beyond this plan; no `research.md`, `data-model.md`, `quickstart.md`, `contracts/`, `TASKS.md`, `BUILD.md`, or source code will be written by Planner.

## Project Structure

### Documentation (this feature)

```text
target/specs/license-detection-classification-files/
├── SPEC.md
└── PLAN.md
```

### Source Code (repository root)

```text
target/
├── go.mod
├── cmd/
│   └── license-checker/
│       └── main.go
└── internal/
    ├── app/              # Composes records.Flatten with license.Apply
    ├── license/          # F-004 classifier, file precedence, record field application
    ├── records/          # F-003 ordered record owner consumed/mutated by F-004
    └── spdxcompat/       # SPDX adapter used by license classifier
```

**Structure Decision**: Follow shared architecture. F-004 implementation is concentrated in `internal/license`, with SPDX parse/validate calls routed only through `internal/spdxcompat`. `internal/app` and `internal/records` wire it into the flattening stage, but F-004 must not own scanning, policy filtering, rendering, or CLI flags.

## Algorithm Fidelity Tables to Preserve

### License classifier exact ordered rule table

The classifier MUST evaluate rules in this order and stop at the first match.

| Order | Source rule | Exact result |
|---:|---|---|
| 1 | Try `spdxExpressionParse(str || '')`; on success | Return raw `str` unchanged |
| 2 | If `str` is truthy | Replace the first newline only: `str = str.replace('\n', '')`, then continue |
| 3 | Undefined/falsey input | `Undefined` |
| 4 | `/The ISC License/` | `ISC*` |
| 5 | `/ermission is hereby granted, free of charge, to any/` | `MIT*` |
| 6 | `/edistribution and use in source and binary forms, with or withou/` | `BSD*` |
| 7 | `/edistribution and use of this software in source and binary forms, with or withou/` | `BSD-Source-Code*` |
| 8 | `/DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE/` | `WTFPL*` |
| 9 | `/\bISC\b/` | `ISC*` |
| 10 | `/\bMIT\b/` | `MIT*` |
| 11 | `/\bBSD\b/` | `BSD*` |
| 12 | `/\bWTFPL\b/` | `WTFPL*` |
| 13 | `/\bApache License\b/` | `Apache*` |
| 14 | CC0 deed regex | `CC0-1.0*` |
| 15 | GPL regex `/\bGNU GENERAL PUBLIC LICENSE\s*Version ([^,]*)/i` | Capture version; append `.0` when captured length is `1`; return `GPL-<version>*` |
| 16 | LGPL regex `/(?:LESSER|LIBRARY) GENERAL PUBLIC LICENSE\s*Version ([^,]*)/i` | Capture version; append `.0` when captured length is `1`; return `LGPL-<version>*` |
| 17 | `/[Pp]ublic [Dd]omain/` | `Public Domain` |
| 18 | URL regex or `SEE LICENSE IN (.*)` | `Custom: <capture>` |
| 19 | Fallback/no match | `null` |

**SPDX parser dependency behavior to preserve**: `spdx-expression-parse ^3.0.0` accepts valid SPDX identifiers/expressions and throws on invalid input; this slice uses only success/failure and returns original input unchanged on success. Worked examples include `MIT`, `LGPL-2.0`, `Apache-2.0`, `BSD-2-Clause`, `(GPL-2.0+ WITH Bison-exception-2.2)`, `LGPL-2.0 OR (ISC AND BSD-3-Clause+)`, and `Apache-2.0 OR ISC OR MIT`.

### License-file precedence table

The detector MUST uppercase the basename, ignore the extension, evaluate patterns in this order, and push only the first matching filename for each pattern.

| Order | Pattern |
|---:|---|
| 1 | `^LICENSE$` |
| 2 | `^LICENSE\-\w+$` |
| 3 | `^LICENCE$` |
| 4 | `^LICENCE\-\w+$` |
| 5 | `^COPYING$` |
| 6 | `^README$` |
| 7 | Fallback/no match: select no file for that pattern and return an empty selection when no patterns match |

## Work Packages

### WP-001 — SPDX adapter contract for classifier pass-through

**Depends:** ARCHITECTURE SPDX decision (`internal/spdxcompat`, `github.com/github/go-spdx/v2/spdxexp v2.7.0`); SPEC FR-002 and Algorithm Fidelity rule 1.  
**Scope:** Define the F-004-facing adapter method needed by `internal/license`, e.g. `ParseExpression(raw string) error` or `IsValidExpression(raw string) bool`, without exposing parser internals. Ensure success/failure only; do not return normalized parser output to classifier.  
**Acceptance:** Valid examples from SPEC line 158 succeed and classifier returns exact original raw input. Invalid/non-SPDX text falls through to regex rules.

### WP-002 — Exact license classifier implementation

**Depends:** WP-001; SPEC FR-001, FR-008, Edge Cases, Algorithm Fidelity table.  
**Scope:** Implement ordered classifier in `internal/license` with explicit rules and comments mapping to table order. Represent fallback `null` in a way downstream records can preserve as a null value where needed, while falsey input returns the literal `Undefined`. Remove only the first newline before regex matching.  
**Acceptance:** Table tests cover SPDX pass-through, falsey input, each guessed `*` literal, `Public Domain`, `Custom: <capture>`, shadowed BSD behavior, GPL/LGPL single-digit `.0` expansion, and fallback `null`.

### WP-003 — Package metadata license normalization

**Depends:** WP-002; F-003 record/node metadata shape; SPEC FR-003, FR-004.  
**Scope:** Apply `json.license || json.licenses` source behavior. Classify strings and supported object/string-like values according to source expectations; map arrays through classification. Use README classification only when package license fields do not provide a license value.  
**Acceptance:** Constructed tests prove package license takes precedence over README, README fills missing license, arrays are mapped, and falsey/missing metadata preserves expected sentinel flow.

### WP-004 — License-file precedence detector

**Depends:** SPEC FR-005 and license-file precedence table.  
**Scope:** Given a package directory filename list, compare `strings.ToUpper(basenameWithoutExtension)` against the exact precedence patterns, pushing only the first filename matching each pattern in table order. Do not select arbitrary fallback files.  
**Acceptance:** Tests mirror `tests/license-files-test.js:10-87`: empty/no-match, extension-insensitive, case-insensitive, `LICENSE.md` before `COPYING`/`README.txt`, and only first `LICENSE-*` for that pattern.

### WP-005 — License/notice file extraction and record fields

**Depends:** WP-004; F-003 path and record APIs; SPEC FR-006, FR-009, FR-010, Output Contract field order.  
**Scope:** Read selected license-file content, set `licenseFile`, set `licenseText`, preserve copyright extraction if provided by source-compatible logic, and detect `noticeFile` when NOTICE exists. Respect ownership boundaries: relative path conversion from `--relativeLicensePath` is coordinated with F-003/F-006 contracts, and renderer text transformations stay out of this package.
**Acceptance:** Record-application tests assert conditional field presence and insertion order expectations: `licenses` exists first from F-003, `licenseFile` after `path`, then `licenseText`/`copyright`, then `noticeFile`.

### WP-006 — License-file override flow

**Depends:** WP-002, WP-004, WP-005; SPEC FR-007, SC-003.  
**Scope:** When current `licenses` is absent, contains `UNKNOWN`, or starts with `Custom:`, classify selected license-file text and allow it to overwrite `licenses` without any CLI flag. Preserve values such as `MIT*` exactly.  
**Acceptance:** Constructed fixtures cover missing license, `UNKNOWN`, `Custom: <url>`, and non-overridden valid package license. The license-file-over-custom-url fixture expects exact `MIT*`.

### WP-007 — Integration seam with F-003 flattening

**Depends:** WP-003 through WP-006; ARCHITECTURE ordered stage sequence line 102.  
**Scope:** Provide a small `Detector`/`Apply` API consumed during `records.Flatten`, with injected SPDX adapter and filesystem reader to keep tests deterministic. Do not introduce CLI flags or renderer dependencies.  
**Acceptance:** F-003-compatible fixture records are mutated only in F-004-owned fields: `licenses`, `licenseFile`, `licenseText`, `copyright`, `noticeFile`.

### WP-008 — Validation suite and conformance evidence

**Depends:** WP-001 through WP-007.  
**Scope:** Add table-driven tests under `target/internal/license` and adapter tests under `target/internal/spdxcompat`. Use Analyzer-recorded/source test expectations only; no live source execution. Run with `go test ./...`.  
**Acceptance:** Tests satisfy SPEC SC-001 through SC-005 and prove downstream renderers can receive exact `UNKNOWN`, `MIT*`, and `ISC` values from this slice.

## Dependencies and Execution Order

1. Implement/lock `internal/spdxcompat` classifier-facing validation behavior (WP-001).
2. Implement pure classifier before any record integration (WP-002).
3. Implement package/README metadata flow over the classifier (WP-003).
4. Implement pure filename precedence detector (WP-004).
5. Implement file-content extraction and field population (WP-005).
6. Implement override conditions using selected file text (WP-006).
7. Wire a minimal detector/apply API into F-003 flattening seam (WP-007).
8. Complete tests and run `go test ./...` (WP-008).

## Validation Approach

- **Unit tests — classifier:** one table row per SPEC classifier rule and edge case; expected output is exact string or null-equivalent. Include all worked SPDX examples and ensure the raw input string is returned, not normalized output.
- **Unit tests — license-file detector:** filename-list tests for exact precedence, case-insensitive matching, extension-insensitive matching, first-match-per-pattern, and no fallback selection.
- **Record-flow tests:** constructed package metadata/readme/license-file fixtures for package license, missing package license with README, absent/`UNKNOWN`/`Custom:` overridden by license-file text, and non-overridden valid SPDX package license.
- **Field contract tests:** ordered record assertions for conditional `licenseFile`, `licenseText`, `copyright`, and `noticeFile` placement consumed by F-006.
- **Command:** `go test ./...` from `target/` after implementation exists.
- **Prohibited validation:** do not install or run the live Node source system; do not invent new byte fixtures.

## Traceability Matrix

| Plan item | SPEC dependency | ARCHITECTURE dependency |
|---|---|---|
| WP-001 | FR-002; Algorithm Fidelity rule 1; SA-007 | Technical Context SPDX; Package Decisions `github.com/github/go-spdx/v2/spdxexp v2.7.0`; Implementation Order step 4 |
| WP-002 | FR-001, FR-008; Edge Cases; SC-001 | `internal/license`; stdlib `regexp`, `strings` |
| WP-003 | FR-003, FR-004; User Story 3 | Ordered stage sequence: `records.Flatten` invokes `license.Apply` |
| WP-004 | FR-005; License-file precedence table; SC-002 | `internal/license` exact file precedence |
| WP-005 | FR-006, FR-009, FR-010; Output Contract field order | `records` ordered records; `render` consumes fields later |
| WP-006 | FR-007; SC-003; license-file-over-custom-url `MIT*` | `license.Apply` during flatten |
| WP-007 | Assumptions lines 192-193; F-003/F-006 boundaries | Structure Decision and Composition Root |
| WP-008 | SC-001 through SC-005; SA-006 no live source execution | Testing decision: standard `testing`, golden fixtures from specs |

## Complexity Tracking

No constitution violations were identified. The SPDX adapter is justified by the shared architecture and SPEC FR-002; hand-rolling broad SPDX parsing is rejected because compatibility with `spdx-expression-parse ^3.0.0` is high-risk and dependency-backed parsing is required.
