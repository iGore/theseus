# Feature Specification: Filtering, Restriction, and License Policy Decisions

**Created**: 2026-06-22  
**Status**: Draft  
**Input**: User description: "rewrite Sourcebot project `davglass/license-checker` into Go; assigned slice F-005 `filtering-restriction-license-policy`"

## Scope

This specification covers only F-005: the post-scan filtering, package restriction, private-package exclusion, unknown-license filtering, and license policy fail decisions applied after package records and license values have already been produced by F-003/F-004. It is source-grounded in the Analyzer artifacts for pinned source commit `de6e9a42513aa38a58efc6b202ee5281ed61f486`, with Sourcebot-guided reconstruction evidence recorded in `ANALYSIS.md`.

**In scope**:

- `--unknown`, `--onlyunknown`, `--exclude`, `--packages`, `--excludePackages`, `--excludePrivatePackages`, `--failOn`, and `--onlyAllow` behavior. [SA-001]
- Private package policy effect on `licenses` and later exclusion from results. [SA-004]
- Exact exclusion transformation, including escaped comma handling, SPDX/BSD alias handling, guessed-license normalization, and `UNKNOWN` preservation. [SA-005]
- Exact stderr strings and exit code `1` behavior for fail policies and CLI preflight conflicts/warnings. [SA-002]

**Out of scope**:

- Dependency graph discovery and flattened record construction, except as prerequisites consumed by this slice. Those belong to F-003. [SA-006]
- License detection/classification and license-file precedence, except the already-computed license strings consumed by filters. Those belong to F-004. [SA-006]
- Renderer byte formatting for tree/JSON/CSV/Markdown/summary/files. That belongs to F-006; this slice only constrains which records reach renderers and which policy diagnostics are emitted. [SA-007]
- Live source-system execution or newly invented stdout captures. The Analyzer explicitly did not execute the source system. [SA-008]

### Source Artifact References

| ID | Evidence |
|---|---|
| SA-001 | Entry Point flag rows: `--unknown`, `--onlyunknown`, `--exclude`, `--failOn`, `--onlyAllow`, `--packages`, `--excludePackages`, `--excludePrivatePackages` in `ANALYSIS.md:33-40`; F-005 block in `USE-CASES.md:41-48`. |
| SA-002 | Golden Diagnostic Snippets and exact stderr strings in `ANALYSIS.md:121-127`; exit-code summary in `ANALYSIS.md:63-66`; tests `tests/failOn-test.js:7-42`. |
| SA-003 | Function inventory for `exports.init`, `transformBSD`, `invert`, and `spdxIsValid` in `ANALYSIS.md:142-146`; main source `lib/index.js:313-458`. |
| SA-004 | Input-field inventory for package `private` and package `path` in `ANALYSIS.md:165`, `170`; private-module tests in `tests/test.js:291-310` and `tests/packages-test.js:37-44`. |
| SA-005 | Exclusion transformation table in `ANALYSIS.md:206-211`; Source-Dependency Contract for `spdx-correct`/`spdx-satisfies` in `ANALYSIS.md:288`; SPDX risk in `ANALYSIS.md:297`. |
| SA-006 | Flow map and end-to-end pipeline in `ANALYSIS.md:256-266`; dependencies F-003/F-004 in `USE-CASES.md:44`. |
| SA-007 | Output Field Contract for key sort, sentinels, and renderer boundaries in `ANALYSIS.md:213-232`. |
| SA-008 | Analyzer open question forbidding invented live byte fixtures in `ANALYSIS.md:304-309`; shared open questions in `USE-CASES.md:85-89`. |
| SA-009 | Requirements RT-005 and RT-007 in `ANALYSIS.md:348-350`; feature-to-task mapping F-005 -> RT-005/RT-007 in `ANALYSIS.md:360`. |

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Exclude packages by license policy (Priority: P1)

A user scans an npm project and removes packages whose current `licenses` value matches a denied license list, while preserving packages with `UNKNOWN` licenses and source-compatible SPDX semantics.

**Why this priority**: License exclusion is the central compliance filter for this slice and depends on exact SPDX and fallback behavior. [SA-005]

**Independent Test**: Can be tested by supplying a prebuilt sorted package-record map with known license strings and applying `--exclude`; the resulting key set proves policy filtering without invoking renderers.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** package records containing licenses `MIT` and `ISC`, **When** `--exclude "MIT, ISC"` is applied, **Then** records whose `licenses` are exactly `MIT` or `ISC` are absent from the restricted result. [SA-005]
2. **Given** a package record with `licenses` equal to `Apache License, Version 2.0`, **When** `--exclude "Apache License\\, Version 2.0"` is applied, **Then** that record is absent after escaped-comma unescaping and literal invalid-SPDX matching. [SA-005]
3. **Given** a package record with `licenses` equal to `BSD-3-Clause`, **When** `--exclude "BSD"` is applied, **Then** the record is absent because exact exclusion token `BSD` expands to `(0BSD OR BSD-2-Clause OR BSD-3-Clause OR BSD-4-Clause)`. [SA-005]
4. **Given** a package record with `licenses` equal to `UNKNOWN`, **When** any `--exclude` list is applied, **Then** the record remains in the result because `UNKNOWN` is never excluded. [SA-005]
5. **Given** a package record with `licenses` equal to `Custom: MY-LICENSE.md`, **When** `--exclude "MIT"` is applied, **Then** the custom-license record remains because it is neither a literal invalid-SPDX match nor an SPDX satisfaction match. [SA-005]

---

### User Story 2 - Restrict result set by package identity and privacy (Priority: P2)

A user limits scan output to specific `name@version` packages, removes specific packages, or excludes packages marked private.

**Why this priority**: Package allow/deny lists and private exclusion shape the observable output key set before rendering and before fail policies. [SA-004]

**Independent Test**: Can be tested by applying the restrictions to a deterministic filtered record map and checking exact output keys in insertion order.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** filtered records keyed `readable-stream@1.1.14`, `spdx-satisfies@4.0.0`, and `y18n@3.2.1`, **When** `--packages "readable-stream@1.1.14;y18n@3.2.1"` is applied, **Then** the restricted result keys are exactly `readable-stream@1.1.14` followed by `y18n@3.2.1` if those keys appeared in the filtered map in that order. [SA-001]
2. **Given** filtered records including `readable-stream@1.1.14`, `spdx-satisfies@4.0.0`, and `y18n@3.2.1`, **When** `--excludePackages "readable-stream@1.1.14;spdx-satisfies@4.0.0;y18n@3.2.1"` is applied, **Then** none of those exact keys appears in the restricted result. [SA-001]
3. **Given** a package record whose `private` field is true, **When** records are sorted before restrictions, **Then** its `licenses` value is overwritten to `UNLICENSED` regardless of the package's source license metadata. [SA-004]
4. **Given** a package record whose `private` field is true, **When** `--excludePrivatePackages` is applied, **Then** that package key is deleted from the restricted result. [SA-004]

---

### User Story 3 - Focus on unknown or guessed licenses (Priority: P3)

A user requests only ambiguous license records, or requests that guessed classifier results be treated as `UNKNOWN`.

**Why this priority**: Unknown handling affects compliance triage and is explicitly covered by source flags and tests. [SA-001]

**Independent Test**: Can be tested with records containing `MIT`, `MIT*`, `UNKNOWN`, and `UNLICENSED` license strings and checking the resulting records.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** records with `licenses` equal to `MIT*`, `UNKNOWN`, and `MIT`, **When** `--onlyunknown` is applied, **Then** only the records whose `licenses` string contains `*` or contains `UNKNOWN` remain. [SA-001]
2. **Given** a record whose `licenses` value is `MIT*`, **When** `--unknown` is applied during sorting, **Then** its `licenses` value becomes `UNKNOWN` before `--onlyunknown`, `--exclude`, package restrictions, and fail policies run. [SA-001]
3. **Given** a record whose `licenses` value is already `UNKNOWN`, **When** `--unknown` is applied, **Then** the value remains `UNKNOWN`. [SA-001]
4. **Given** a private package whose `licenses` has been overwritten to `UNLICENSED`, **When** `--onlyunknown` is applied, **Then** the package is not included unless its current license string contains `*` or `UNKNOWN`. [SA-004]

---

### User Story 4 - Fail fast on forbidden or non-allowed licenses (Priority: P4)

A user configures policy failure rules so the process exits with code `1` and emits exact stderr diagnostics on the first policy violation.

**Why this priority**: Fail policies are a compliance automation surface, but they run after result-set filtering and are partly shared with F-007 diagnostics. [SA-002]

**Independent Test**: Can be tested by applying fail policies to a deterministic restricted record map while capturing stderr and exit code.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** a restricted record whose `licenses` value is exactly `MIT`, **When** `--failOn "MIT;ISC"` is applied, **Then** stderr receives `Found license defined by the --failOn flag: "MIT". Exiting.` and the process exits with code `1`. [SA-002]
2. **Given** a restricted record whose `licenses` value is `Apache License, Version 2.0`, **When** `--failOn "Apache License, Version 2.0"` is applied, **Then** stderr receives `Found license defined by the --failOn flag: "Apache License, Version 2.0". Exiting.` and the process exits with code `1`. [SA-002]
3. **Given** a restricted record keyed `pkg@1.0.0` whose `licenses` value is `GPL-3.0`, **When** `--onlyAllow "MIT;ISC"` is applied, **Then** stderr receives `Package "pkg@1.0.0" is licensed under "GPL-3.0" which is not permitted by the --onlyAllow flag. Exiting.` and the process exits with code `1`. [SA-002]
4. **Given** CLI arguments containing both `--failOn` and `--onlyAllow`, **When** argument preflight runs, **Then** stderr receives `--failOn and --onlyAllow can not be used at the same time. Choose one or the other.` and the process exits with code `1` before scanning. [SA-002]
5. **Given** CLI arguments `--failOn "MIT,ISC"`, **When** argument preflight runs, **Then** stderr receives `Warning: As of v17 the --failOn argument takes semicolons as delimeters instead of commas (some license names can contain commas)` and scanning continues. [SA-002]

### Edge Cases

- If all records are removed before the sorted result is built, the source sets error text `No packages found in this path..`; downstream error emission belongs to F-007, but this slice MUST preserve the empty-result condition that triggers it. [SA-003]
- If both `--packages` and `--excludePackages` are present, source order means `--excludePackages` rebuilds `restricted` from `filtered`, not from the previous whitelist result. The blacklist therefore overrides the whitelist result construction. [SA-003]
- `--packages` and `--excludePackages` split only on `;` and do not trim tokens. Exact `name@version` string equality is required. [SA-003]
- `--exclude` trims tokens and unescapes `\,`; `--failOn`, `--onlyAllow`, `--packages`, and `--excludePackages` do not share that comma-unescape transform. [SA-003]
- A license containing `UNKNOWN` is assigned into `filtered` immediately during exclusion and MUST NOT be removed by the exclusion rule. [SA-005]
- Guessed license normalization for exclusion removes the final character whenever the license string contains `*`; the source assumes guessed literals use a trailing `*`. [SA-005]
- `--onlyAllow` uses `indexOf` semantics over the current license value: for string license values this is substring matching, so an allowed token contained anywhere in the license string marks the package allowed. [SA-003]
- CLI warns on comma-containing `--failOn` or `--onlyAllow` values but still treats the policy list as semicolon-separated. [SA-002]

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST apply F-005 only after F-003/F-004 have produced a map of package records keyed by exact `name@version` with current `licenses`, optional `private`, and other output fields. [SA-006]
- **FR-002**: During the lexically sorted pass over package keys, the system MUST overwrite `record.licenses` to exact literal `UNLICENSED` whenever `record.private` is truthy, regardless of package license metadata or earlier classification. In colorized tree mode this value may be colorized by F-006; the no-color policy value is exactly `UNLICENSED`. [SA-004]
- **FR-003**: During the same sorted pass, if `record.licenses` is falsy after private handling, the system MUST set it to exact literal `UNKNOWN`. [SA-007]
- **FR-004**: When `--unknown` is true, the system MUST rewrite any non-`UNKNOWN` `record.licenses` value containing `*` to exact literal `UNKNOWN` before `--onlyunknown`, exclusion, package restrictions, private exclusion, and fail policies run. [SA-001]
- **FR-005**: When `--onlyunknown` is true, the system MUST include only records whose current `licenses` value contains `*` or contains `UNKNOWN`; otherwise it MUST include all records from the sorted pass. [SA-001]
- **FR-006**: Top-level package keys MUST be processed in lexical order before filter/restriction decisions, and whitelist/blacklist restriction maps MUST preserve the filtered map's key iteration order. [SA-007]
- **FR-007**: If the sorted result contains no keys, the system MUST set an input error equivalent to `Error('No packages found in this path..')`; downstream stderr formatting for that error is owned by F-007. [SA-003]
- **FR-008**: The `--exclude` option MUST parse its comma list with source regex `/([^\\\][^,]|\\,)+/g`, then for each matched token MUST replace `\,` with `,` and trim leading/trailing whitespace. [SA-005]
- **FR-009**: For exclusion matching, the exact exclusion token `BSD` MUST transform to `(0BSD OR BSD-2-Clause OR BSD-3-Clause OR BSD-4-Clause)` before SPDX validity partitioning. No other token receives this alias transform. [SA-005]
- **FR-010**: Exclusion tokens MUST be partitioned into valid SPDX exclusions where `spdxCorrect(token) === token` and invalid exclusions where that predicate is false. The SPDX excluder expression MUST be built exactly as `( ` + valid exclusions joined by ` OR ` + ` )`. [SA-005]
- **FR-011**: For each candidate package license during exclusion, the system MUST treat the current `licenses` value as an array by concatenating scalar values into a one-item list, then evaluate each license value. [SA-003]
- **FR-012**: During exclusion, if a candidate license contains `UNKNOWN`, the package MUST be kept; `UNKNOWN` MUST NOT be excluded by any `--exclude` list. [SA-005]
- **FR-013**: During exclusion, non-`UNKNOWN` candidate licenses containing `*` MUST be normalized by removing the final character before matching, preserving source behavior for guessed literals such as `MIT*` -> `MIT`. [SA-005]
- **FR-014**: During exclusion, a normalized candidate license exactly equal to `BSD` MUST transform to `(0BSD OR BSD-2-Clause OR BSD-3-Clause OR BSD-4-Clause)` before matching. [SA-005]
- **FR-015**: During exclusion, a package MUST be removed when any normalized candidate license is present in the invalid-exclusion list or when `spdxCorrect(candidate)` returns a truthy corrected license and `spdxSatisfies(spdxCorrect(candidate), spdxExcluder)` is true. [SA-005]
- **FR-016**: During exclusion, a package MUST be kept when no candidate license matches the invalid-exclusion list or SPDX satisfaction expression. [SA-005]
- **FR-017**: `--packages` MUST split its value on `;`, MUST NOT trim tokens, and MUST rebuild `restricted` to contain only exact keys from the filtered map whose `name@version` key is included in that split list. [SA-001]
- **FR-018**: `--excludePackages` MUST split its value on `;`, MUST NOT trim tokens, and MUST rebuild `restricted` from the full filtered map to contain only exact keys not included in that split list, even if `--packages` was previously applied. [SA-001]
- **FR-019**: `--excludePrivatePackages` MUST delete any currently restricted record whose `private` field is truthy. [SA-004]
- **FR-020**: `--failOn` values MUST be split on `;`, each token MUST be trimmed, empty trimmed tokens MUST be ignored, and a restricted record MUST fail only when its current `licenses` value exactly equals one of the resulting policy tokens. [SA-002]
- **FR-021**: On a `--failOn` match, the system MUST emit exactly `Found license defined by the --failOn flag: "<license>". Exiting.` to stderr and exit with code `1`. [SA-002]
- **FR-022**: `--onlyAllow` values MUST be split on `;`, each token MUST be trimmed, empty trimmed tokens MUST be ignored, and a restricted record MUST pass when at least one allowed token is found by source-compatible `indexOf` semantics in the current `licenses` value. [SA-002]
- **FR-023**: On an `--onlyAllow` violation, the system MUST emit exactly `Package "<item>" is licensed under "<license>" which is not permitted by the --onlyAllow flag. Exiting.` to stderr and exit with code `1`. [SA-002]
- **FR-024**: CLI preflight MUST reject simultaneous `--failOn` and `--onlyAllow` before scanning, emit exactly `--failOn and --onlyAllow can not be used at the same time. Choose one or the other.` to stderr, and exit with code `1`. [SA-002]
- **FR-025**: CLI preflight MUST warn, but not fail, when a provided `--failOn` or `--onlyAllow` value contains `,`; the warning string MUST be exactly `Warning: As of v17 the --<argName> argument takes semicolons as delimeters instead of commas (some license names can contain commas)`, where `<argName>` is `failOn` or `onlyAllow`. [SA-002]
- **FR-026**: In the programmatic API path, if both `failOn` and `onlyAllow` are provided without CLI preflight, source-compatible list construction MUST let the later `failOn` branch populate the active fail list and leave the only-allow list empty. [SA-003]
- **FR-027**: The system MUST preserve source dependency semantics for `spdx-correct` and `spdx-satisfies` in exclusion filtering; if the Go target does not use equivalent libraries, it MUST implement a compatibility layer covered by examples for MIT/ISC, escaped-comma Apache text, BSD alias, and Public Domain. [SA-005]

### Output Contract *(mandatory for any slice that produces observable output)*

- **Golden baseline**: This slice has no invented live stdout capture. The committed baselines are: exact stderr strings in `ANALYSIS.md:121-127`; package key-set assertions from `tests/packages-test.js:8-44`; private `UNLICENSED` assertion from `tests/test.js:291-310`; filter/fail behavior assertions from `tests/test.js:137-289`, `555-607`; and fail CLI exit/warning assertions from `tests/failOn-test.js:7-42`. [SA-002] [SA-004] [SA-008]
- **Field order & presence**: This slice MUST NOT reorder fields inside records. Per source record contract, `licenses` is inserted first, `private` second only when truthy, and later fields retain F-003/F-004 insertion order. F-005 mutates `licenses`, deletes records, and rebuilds top-level maps; it does not add output fields. [SA-007]
- **Top-level key order**: Initial `sorted` keys are lexical (`Object.keys(data).sort()`). `filtered`, `restricted`, `--packages`, and `--excludePackages` preserve the iteration order of the source map they copy from, except records removed by filters are absent. [SA-007]
- **Conditionally present records**: Records are conditionally absent when removed by `--onlyunknown`, `--exclude`, `--packages`, `--excludePackages`, or `--excludePrivatePackages`; fail policies do not return a successful output when they exit. [SA-001]
- **Sentinels & literals**: Relevant exact values are `UNKNOWN`, `UNLICENSED`, guessed `*` licenses such as `MIT*`, custom literals such as `Custom: MY-LICENSE.md`, `Public Domain`, and raw SPDX expressions. This slice specifically rewrites private packages to `UNLICENSED` and guessed `*` licenses to `UNKNOWN` under `--unknown`. [SA-007]
- **Sort/aggregation order**: No aggregation is introduced by this slice. It consumes lexically sorted package keys and preserves copied-map order as described above. Summary aggregation belongs to F-006. [SA-007]
- **Whitespace contract**: Stderr diagnostics for this slice are exact single-line strings as specified in FR-021, FR-023, FR-024, and FR-025. Console stderr APIs add their normal line terminator; this spec does not invent a full stderr byte capture beyond the committed exact string snippets. Renderer stdout/file newline behavior belongs to F-006. [SA-002] [SA-008]

### Algorithm Fidelity *(mandatory when the slice transforms values)*

#### Ordered F-005 Pipeline

1. Build `toCheckforOnlyAllow` and `toCheckforFailOn` before scanning. `onlyAllow` initially selects the only-allow list; `failOn`, if present, selects the fail-on list. Split the selected checker string by `;`, trim each token, and push only non-empty tokens. [SA-003]
2. Flatten source package data via F-003/F-004, then iterate `Object.keys(data).sort()`. [SA-006]
3. For each sorted item: private override to `UNLICENSED`; falsy license fallback to `UNKNOWN`; optional `--unknown` rewrite of guessed `*` licenses to `UNKNOWN`; optional `--onlyunknown` inclusion gate; otherwise include in `sorted`. [SA-001] [SA-004]
4. If `--exclude` is present, parse exclusion tokens, perform BSD alias transform, partition valid vs invalid SPDX tokens, and evaluate every sorted record license as described in FR-008 through FR-016. Otherwise set `filtered = sorted`. [SA-005]
5. Set `restricted = filtered`. If `--packages` is present, rebuild from filtered exact included keys. If `--excludePackages` is present, rebuild again from filtered exact non-excluded keys. If `--excludePrivatePackages` is present, delete truthy-private records from the current restricted map. [SA-001] [SA-004]
6. Iterate restricted records in current key order. Apply `failOn` exact equality first. Apply `onlyAllow` containment/index check second. On the first violation, emit exact stderr and exit with code `1`. [SA-002]
7. If no fail policy exits, callback/return the final restricted map to downstream renderers. [SA-003]

#### Complete Exclusion Rule Table

| Order | Input | Transformation / predicate | Result |
|---:|---|---|---|
| 1 | `options.exclude` string | Match tokens with `/([^\\\][^,]|\\,)+/g` | Candidate exclusion tokens. [SA-005] |
| 2 | Each exclusion token | Replace `\,` with `,`; trim leading/trailing whitespace | Normalized exclusion token. [SA-005] |
| 3 | Normalized exclusion token exactly `BSD` | Replace with `(0BSD OR BSD-2-Clause OR BSD-3-Clause OR BSD-4-Clause)` | BSD alias expression. [SA-005] |
| 4 | Transformed exclusion token | If `spdxCorrect(token) === token` | Add to valid SPDX exclusions. [SA-005] |
| 5 | Transformed exclusion token | If `spdxCorrect(token) !== token` | Add to invalid literal exclusions. [SA-005] |
| 6 | Valid SPDX exclusions | Join with ` OR ` inside `( ` and ` )` | SPDX excluder expression. [SA-005] |
| 7 | Missing/falsy record license | Protective branch | Keep record. [SA-003] |
| 8 | Record license value | Convert scalar to one-item array; arrays remain arrays | Candidate license list. [SA-003] |
| 9 | Candidate license containing `UNKNOWN` | Match by substring, including colorized forms | Keep record; `UNKNOWN` is never excluded. [SA-005] |
| 10 | Candidate non-UNKNOWN license containing `*` | Remove final character | Guessed-license normalization, e.g. `MIT*` -> `MIT`. [SA-005] |
| 11 | Candidate normalized license exactly `BSD` | Replace with BSD alias expression | SPDX-compatible BSD matching. [SA-005] |
| 12 | Candidate normalized license in invalid literal exclusions | Literal equality | Mark package as matched for removal. [SA-005] |
| 13 | Candidate normalized license where `spdxCorrect(candidate)` truthy and `spdxSatisfies(spdxCorrect(candidate), spdxExcluder)` true | SPDX satisfaction | Mark package as matched for removal. [SA-005] |
| 14 | No candidate matched removal | Fallback branch | Keep record. [SA-005] |

#### Worked Transform Examples

- `--exclude "MIT, ISC"` with record `licenses: "MIT"` -> record removed. [SA-005]
- `--exclude "Apache License\\, Version 2.0"` with record `licenses: "Apache License, Version 2.0"` -> token unescapes to `Apache License, Version 2.0`; invalid literal match removes record. [SA-005]
- `--exclude "BSD"` with record `licenses: "BSD-3-Clause"` -> `BSD` alias SPDX expression satisfies `BSD-3-Clause`; record removed. [SA-005]
- `--exclude "MIT"` with record `licenses: "Custom: MY-LICENSE.md"` -> no literal or SPDX match; record kept. [SA-005]
- `--unknown` with record `licenses: "Apache*"` -> current license becomes `UNKNOWN`. [SA-001]
- `--onlyunknown` over `MIT`, `MIT*`, `UNKNOWN`, `UNLICENSED` -> keeps `MIT*` and `UNKNOWN`; removes `MIT` and `UNLICENSED`. [SA-001] [SA-004]
- Private record with source package license `MIT` and `private: true` -> current policy license becomes `UNLICENSED`; with `--excludePrivatePackages` the record is deleted. [SA-004]

#### Source Dependency Behavior to Preserve

- `spdx-correct`/`spdx-satisfies` behavior MUST be preserved for exclusion filtering. Required examples: excluding `MIT`/`ISC`, escaped-comma `Apache License, Version 2.0`, BSD alias matching `BSD-3-Clause`, and literal `Public Domain` behavior. [SA-005]
- The target MUST NOT silently replace SPDX satisfaction with naive string equality except for invalid exclusion tokens, because the source delegates valid SPDX expression semantics to these libraries. [SA-005]

### Key Entities *(include if feature involves data)*

- **Package record**: Mutable output record keyed by exact `name@version`; this slice consumes `licenses` and optional `private`, mutates `licenses` for private/unknown policy, and may delete the record from returned maps. [SA-007]
- **Restriction set**: The ordered map after `--exclude`, `--packages`, `--excludePackages`, and `--excludePrivatePackages` are applied. [SA-003]
- **License policy list**: Semicolon-split, trimmed, non-empty tokens used by `--failOn` or `--onlyAllow`. [SA-002]
- **Exclusion token list**: Comma/escaped-comma parsed tokens used by `--exclude`, with BSD alias and SPDX validity partitioning. [SA-005]

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: For fixtures corresponding to `tests/test.js:137-197`, exclusion results match source assertions for MIT/ISC removal, escaped-comma Apache removal, BSD alias removal, Public Domain removal, and custom-license preservation. [SA-005]
- **SC-002**: For fixtures corresponding to `tests/packages-test.js:8-44`, JSON-parsed output key sets match source assertions for `--packages`, `--excludePackages`, and `--excludePrivatePackages`. [SA-001] [SA-004]
- **SC-003**: For fixtures corresponding to `tests/test.js:555-607`, `--onlyunknown` keeps only records whose current license contains `UNKNOWN` or `*`, and does not accidentally behave as if enabled when absent. [SA-001]
- **SC-004**: Policy failure diagnostics are byte-identical for the committed exact stderr substrings: failOn, onlyAllow, mutual exclusion, and comma warning, including misspelling `delimeters`; failing policy paths exit with code `1`. [SA-002]
- **SC-005**: Private package records always receive `UNLICENSED` before filtering, and `--excludePrivatePackages` yields an empty key set for the private fixture covered by `tests/packages-test.js:37-44`. [SA-004]
- **SC-006**: The target's SPDX compatibility layer/library passes the required exclusion examples for valid SPDX, invalid literal, BSD alias, guessed-license normalization, and fallback no-match behavior. [SA-005]
- **SC-007**: For the committed golden fixtures available to this slice, the target's returned restricted package map and policy stderr/exit behavior are parity-equivalent to the source; exact byte parity is required for stderr snippets, while full renderer stdout byte parity is deferred to F-006 because no live stdout captures were generated. [SA-008]

## Assumptions

- F-003 provides source-compatible flattened package records keyed by `name@version` and already sorted or sortable by lexical key. [SA-006]
- F-004 provides source-compatible `licenses` values, including raw SPDX expressions, guessed `*` literals, `Custom: <value>`, `Public Domain`, `UNKNOWN`, `UNLICENSED`, and `Undefined`. [SA-007]
- No-color output is the baseline for exact policy values; colorized sentinel rendering is handled by F-006 and terminal-dependent color support. [SA-007]
- The Go implementation may use a Go SPDX library or a compatibility layer, but its behavior must match the `spdx-correct`/`spdx-satisfies` examples and fallback branches in this spec. [SA-005]
- Full byte-exact stdout fixtures are intentionally not invented here; maintainers may run live differential checks manually out-of-flow as described by repository workflow rules. [SA-008]
- Sourcebot MCP was unavailable during Analyzer reconstruction; this SPEC relies on the Analyzer's Sourcebot-skill-guided artifacts and pinned source references, and should be rechecked if indexed Sourcebot access becomes available. [SA-008]

## Open Questions

- Which exact Go SPDX library or compatibility layer will be selected is an Architecture decision, but it must satisfy the behavior specified here. [SA-005]
- No full live stdout captures exist for combinations of these filters with every renderer; F-006 owns renderer byte contracts and any known gaps. [SA-008]
