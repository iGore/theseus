# Implementation Plan: Filtering, Restriction, and License Policy Decisions

**Date**: 2026-06-22 | **Spec**: `target/specs/filtering-restriction-license-policy/SPEC.md`  
**Input**: Feature specification from `target/specs/filtering-restriction-license-policy/SPEC.md` and shared architecture from `target/specs/ARCHITECTURE.md`

**Note**: This plan is scoped to F-005 only. It does not plan renderer byte formatting, dependency scanning, license classification, TASKS.md, BUILD.md, or Go source edits.

## Summary

Implement `internal/policy` for the Go `license-checker` target so already-flattened package records from F-003/F-004 are transformed by source-compatible filtering, package restriction, private-package policy, unknown-license policy, exclusion policy, and fail policy rules. The policy engine must preserve ordered output behavior, exact sentinel literals (`UNKNOWN`, `UNLICENSED`), `spdx-correct`/`spdx-satisfies` semantics through `internal/spdxcompat`, BSD alias behavior, and exact stderr/exit interactions delegated through diagnostics.

## Progress Checklist

- [x] Planning scope confirmed against `SPEC.md`: F-005 only; consumes F-003/F-004 records and feeds F-006 renderers.
- [x] Technical context filled with concrete decisions from `ARCHITECTURE.md` and F-005 spec.
- [x] Work packages derived from assigned functionality scope.
- [x] Dependencies and execution order documented.
- [x] Validation approach captured.

## Decision Log

- [x] Implement policy behavior in `target/internal/policy`, matching architecture lines 63-65 and implementation order line 193.
- [x] Keep all SPDX validity/correction/satisfaction calls behind `target/internal/spdxcompat`, per architecture lines 23, 63, 169, and 178-179.
- [x] Use ordered record/map abstractions from `target/internal/ordered`/`target/internal/records`; do not expose plain Go map iteration at observable boundaries.
- [x] Preserve diagnostics as exact strings through `target/internal/diagnostics`; policy should return/write source-compatible exit decisions without inventing additional output.
- [x] Treat Context7 `/golang/go` guidance as supporting the internal package boundary and standard `testing`/`go test` workflow; no external test/assertion framework is planned.

## Open Questions

- [ ] Confirm the concrete `spdxcompat` adapter behavior once implemented against `github.com/github/go-spdx/v2/spdxexp v2.7.0`; the architecture selected the facade, but source-compatible correction edge cases must be locked by tests before broad use.
- [ ] Confirm the final ordered record type API from earlier work packages; this plan assumes a mutable ordered collection can iterate keys lexically and rebuild ordered subsets.

## Technical Context

**Language/Version**: Go, module rooted under `target/`; exact Go version follows `target/go.mod` when created.  
**Primary Dependencies**: Standard library (`sort`, `strings`, `regexp`, `io`, `testing`) plus architecture-selected `github.com/github/go-spdx/v2/spdxexp v2.7.0` hidden behind `internal/spdxcompat`.  
**Storage**: N/A for F-005; operates in memory on package records.  
**Testing**: Standard `testing` package, table-driven unit tests for `internal/policy`, adapter contract tests for `internal/spdxcompat`, CLI/process tests for preflight diagnostics where owned by diagnostics/app integration; run with `go test ./...`.  
**Target Platform**: POSIX-like CLI on darwin/linux first; F-005 itself is path-independent.  
**Project Type**: Go CLI with internal packages; `cmd/license-checker/main.go` delegates to `internal/app`, and F-005 lives under `internal/policy`.  
**Performance Goals**: Linear policy passes over already-flattened records except initial lexical key sort (`O(n log n)`); no additional filesystem work.  
**Constraints**: No live source-system execution in-flow; preserve exact source literals, ordered top-level record behavior, stderr strings, exit code `1` for policy failures, and manual/out-of-flow live differential rule.  
**Scale/Scope**: Bounded to post-scan policy decisions over the package record map produced by F-003/F-004.

## Constitution Check

*GATE: Must align with applicable workflow rules before planning proceeds. Re-check after design-oriented planning artifacts are produced.*

- Inputs were limited by the assignment to this folder's `SPEC.md` and shared `ARCHITECTURE.md`; no additional design artifacts are produced.
- Go remains the target language, matching repository workflow and architecture.
- Artifact ownership is respected: Planner writes only this `PLAN.md`.
- F-005 remains bounded and does not implement F-003/F-004/F-006/F-007 behavior except where policy must interact with already-defined diagnostics/preflight contracts.
- Context7 was used for Go package/internal visibility and `testing`/`go test` workflow planning; no extra framework dependency is introduced.

## Project Structure

### Documentation (this feature)

```text
target/specs/filtering-restriction-license-policy/
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
    ├── app/             # invokes policy.Apply after records/license stages
    ├── cli/             # parses F-005 flags into Options
    ├── diagnostics/     # exact preflight and policy stderr/exit literals
    ├── ordered/         # observable ordered map/record primitives
    ├── policy/          # F-005 engine and tests
    ├── records/         # provides flattened mutable records consumed here
    └── spdxcompat/      # SPDX correction/satisfaction facade used by policy
```

**Structure Decision**: Implement this slice in `internal/policy` with narrow dependencies on option structs, ordered records, `diagnostics`, and `spdxcompat`. The CLI/app layers only pass parsed options and wire stderr/exit behavior.

## Complexity Tracking

No constitution/workflow violations are planned.

## Work Packages

### WP-001 — Define F-005 policy input/output contract

**Depends:** F-003/F-004 record contract; `ARCHITECTURE.md` ordered data decision lines 22, 61, 64.  
**Spec trace:** FR-001, Output Contract lines 154-159, Key Entities lines 208-213.  
**Deliverable:** Policy-facing types/options that accept an ordered package record collection keyed by exact `name@version`, with mutable `licenses` and optional truthy `private` fields.  
**Notes:** The contract must preserve field order inside records and top-level key order in rebuilt collections; F-005 may mutate `licenses` and delete/rebuild records, but must not add fields.

### WP-002 — Sorted pass, private override, unknown normalization, only-unknown filtering

**Depends:** WP-001.  
**Spec trace:** FR-002 through FR-007; User Story 2 scenario 3; User Story 3 scenarios 1-4; Algorithm Pipeline steps 2-3.  
**Deliverable:** First policy pass that iterates lexical keys, rewrites truthy-private records to `UNLICENSED`, fills falsy licenses with `UNKNOWN`, optionally rewrites guessed `*` licenses to `UNKNOWN` under `--unknown`, and applies `--onlyunknown`.  
**Notes:** Private override occurs before `--unknown` and `--onlyunknown`; `UNLICENSED` is not unknown unless its current license string contains `*` or `UNKNOWN`. Empty sorted result must expose the `No packages found in this path..` condition for F-007 diagnostics without formatting it here.

### WP-003 — Exclusion parser and SPDX/BSD compatibility hooks

**Depends:** WP-001; `internal/spdxcompat` adapter skeleton from architecture implementation order step 4.  
**Spec trace:** FR-008 through FR-010, FR-027; Exclusion Rule Table rows 1-6; Source Dependency Behavior lines 203-206.  
**Deliverable:** `--exclude` token parser using source-compatible escaped-comma semantics, `BSD` alias expansion to `(0BSD OR BSD-2-Clause OR BSD-3-Clause OR BSD-4-Clause)`, valid-vs-invalid partitioning using `spdxCorrect(token) === token`, and construction of `( <valids joined by OR> )`.  
**Notes:** Only `--exclude` receives comma/escaped-comma parsing and trimming. `--failOn`, `--onlyAllow`, `--packages`, and `--excludePackages` must not reuse this parser.

### WP-004 — Exclusion evaluator over sorted records

**Depends:** WP-002 and WP-003.  
**Spec trace:** FR-011 through FR-016; acceptance scenarios 1.1-1.5; Exclusion Rule Table rows 7-14; SC-001 and SC-006.  
**Deliverable:** Filtering pass that keeps `UNKNOWN`-containing licenses, normalizes guessed `*` candidate licenses by removing the final character, applies candidate `BSD` aliasing, removes invalid-literal matches, and removes valid SPDX satisfaction matches through `spdxcompat`.  
**Notes:** Scalar license values must be evaluated as one-item lists; array-like record license support should follow the source-compatible candidate-list contract. Custom licenses such as `Custom: MY-LICENSE.md` remain unless they match invalid literal exclusions. Do not replace SPDX satisfaction with naive equality.

### WP-005 — Package whitelist, blacklist, and private exclusion restrictions

**Depends:** WP-004 when `--exclude` is present; WP-002 otherwise.  
**Spec trace:** FR-017 through FR-019; User Story 2 scenarios 1, 2, 4; Edge Cases lines 111-114; SC-002 and SC-005.  
**Deliverable:** Restriction pass that starts from `filtered`, applies `--packages` by semicolon exact key inclusion without trimming, applies `--excludePackages` by rebuilding again from the full `filtered` map without trimming, and then deletes truthy-private records for `--excludePrivatePackages`.  
**Notes:** If both `--packages` and `--excludePackages` are present, blacklist rebuilds from `filtered`, not the whitelist result. Preserve filtered iteration order in all rebuilt maps.

### WP-006 — Fail policy list construction and restricted-record policy checks

**Depends:** WP-005; diagnostics literals available from `internal/diagnostics`.  
**Spec trace:** FR-020 through FR-023, FR-026; User Story 4 scenarios 1-3; Algorithm Pipeline steps 1 and 6; SC-004.  
**Deliverable:** Semicolon-split, trimmed, non-empty policy token lists for `failOn` and `onlyAllow`; restricted-record iteration that applies `failOn` exact equality first, then `onlyAllow` source-compatible `indexOf`/substring containment, stopping on the first violation with exit code `1`.  
**Notes:** Programmatic API behavior must match source branch ordering: when both values reach policy without CLI preflight, later `failOn` populates the active fail list and leaves only-allow empty. Diagnostics strings must be exactly:

- `Found license defined by the --failOn flag: "<license>". Exiting.`
- `Package "<item>" is licensed under "<license>" which is not permitted by the --onlyAllow flag. Exiting.`

### WP-007 — CLI preflight diagnostics interaction for F-005 flags

**Depends:** CLI parser option fields from F-001/F-002; diagnostics package.  
**Spec trace:** FR-024, FR-025; User Story 4 scenarios 4-5; Output Contract line 160; architecture stage sequence lines 97-103.  
**Deliverable:** Integration plan for diagnostics preflight to reject simultaneous `--failOn` and `--onlyAllow` before scanning, and to warn but continue when either value contains `,`.  
**Notes:** Exact strings and exit code behavior are diagnostics-owned, but F-005 planning requires tests proving policy is not reached on conflict and is reached after warnings. Preserve misspelling `delimeters`.

### WP-008 — Validation fixtures and parity tests for F-005

**Depends:** WP-001 through WP-007.  
**Spec trace:** Success Criteria SC-001 through SC-007; architecture acceptance reference lines 152-154; conformance matrix lines 215-222.  
**Deliverable:** Table-driven tests over in-memory records plus minimal CLI/preflight tests for diagnostics interactions.  
**Notes:** Tests should derive fixtures from committed SPEC/Analyzer contracts, not from live source execution. Renderer stdout byte parity remains F-006.

## Dependencies and Execution Order

1. Confirm ordered record API from F-003/F-004 implementation (`records`/`ordered`).
2. Implement/test `spdxcompat` examples needed by F-005 before exclusion integration.
3. Implement the sorted normalization pass before exclusion and restrictions, because every later decision uses current `licenses` values.
4. Implement exclusion parser/evaluator before package restrictions, matching pipeline step 4.
5. Implement package restrictions and private deletion before fail policies, matching pipeline step 5.
6. Implement fail policy checks after restriction, matching pipeline step 6.
7. Wire diagnostics preflight tests for mutual exclusion and comma warnings at CLI/app level.
8. Run complete F-005 test set and then `go test ./...`.

## Validation Approach

- Unit-test `--onlyunknown` and `--unknown` with `MIT`, `MIT*`, `UNKNOWN`, and `UNLICENSED` records; verify key retention and license mutation order.
- Unit-test private package behavior: private records become `UNLICENSED` before filtering and are deleted only when `--excludePrivatePackages` is set.
- Unit-test `--exclude` examples for `MIT, ISC`, escaped-comma `Apache License\, Version 2.0`, `BSD` alias removing `BSD-3-Clause`, `Public Domain`, guessed-license normalization, and custom-license no-match.
- Unit-test `--packages` and `--excludePackages` with semicolon splitting, no trimming, exact `name@version` equality, preserved order, and blacklist-overrides-whitelist reconstruction.
- Unit-test `--failOn` exact equality and `--onlyAllow` substring/index semantics against restricted records; assert first violation, exact stderr string, and exit code `1`.
- Integration-test CLI preflight for `--failOn` plus `--onlyAllow` conflict and comma-containing warning values; assert exact stderr strings and scan/policy reachability behavior.
- Run `go test ./...`; no live source-system execution is part of this validation.

## Traceability Matrix

| Plan item | SPEC dependency | ARCHITECTURE dependency |
|---|---|---|
| WP-001 | FR-001; Output Contract; Key Entities | Ordered data line 22; `records`/`policy` boundaries lines 61 and 64 |
| WP-002 | FR-002-FR-007; SC-003, SC-005 | Stage sequence line 103; conformance rows `--unknown`, `--onlyunknown`, `--excludePrivatePackages` |
| WP-003 | FR-008-FR-010, FR-027 | SPDX decision lines 23, 169, 178-179 |
| WP-004 | FR-011-FR-016; SC-001, SC-006 | `policy` + `spdxcompat` boundaries lines 63-64 |
| WP-005 | FR-017-FR-019; SC-002, SC-005 | Flag propagation rows for packages/exclude/private lines 140-142 |
| WP-006 | FR-020-FR-023, FR-026; SC-004 | Diagnostics boundary line 66; stage sequence line 103 |
| WP-007 | FR-024-FR-025; Output whitespace contract | Preflight stage lines 97-99; flag propagation rows 138-139 |
| WP-008 | SC-001-SC-007 | Testing decision line 24; implementation order line 193; conformance rows 215-222 |

## Handoff Notes for Task Decomposition

- Break tasks around the work-package boundaries above; do not combine SPDX adapter work with policy filtering until adapter examples pass.
- Keep diagnostics strings centralized; tests may assert strings through diagnostics helpers but policy must not duplicate divergent literals.
- Preserve source-compatible quirks intentionally: BSD alias expression, `UNKNOWN` never excluded, guessed-license final-character removal, blacklist rebuild from filtered map, no trimming for package lists, `onlyAllow` substring semantics, and `delimeters` misspelling.
- Do not add renderer assertions or full stdout golden generation in this slice beyond key-set/record-map outcomes and stderr snippets.
