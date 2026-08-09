# Feature Specification: Output Rendering and Side-Effect Outputs

**Created**: 2026-06-22  
**Status**: Draft  
**Input**: User description: "rewrite Sourcebot project `davglass/license-checker` into Go; assigned functionality slice F-006 — Output rendering and side-effect outputs"

## Scope

This specification covers only F-006: rendering already-scanned, already-filtered package records into tree, JSON, CSV, Markdown, summary, output-file, and copied-license-file artifacts for the pinned `davglass/license-checker` source version `de6e9a42513aa38a58efc6b202ee5281ed61f486` / package version `25.0.1`. It depends on F-003, F-004, and F-005 to provide the sorted package map, license values, path fields, custom-format fields, and filtered record set.

In scope:

- CLI renderer selection after `checker.init(args, cb)` returns records.
- Programmatic renderer helpers equivalent to `print`, `asTree`, `asSummary`, `asCSV`, `asMarkDown`, and `asFiles`.
- Output field order, renderer-specific row/header order, sentinel literals, collection ordering, and stdout-vs-file newline behavior.
- Side effects for `--out` and `--files`, including directory creation and copied license-file naming.
- Output-affecting unconditional overrides from the analyzer input-field inventory that must be visible in rendered output.

Out of scope:

- CLI preflight help/version/mutual-exclusion behavior and process exit codes, except where renderer selection reaches stdout/file side effects; those are F-001/F-007.
- Dependency graph discovery, license classification rule internals, filtering policy, and package restrictions, except as renderer input contracts supplied by F-003/F-004/F-005.
- Live differential execution of the Node.js source. No new live-output captures may be invented for this spec.

## Source Evidence Index

- **SA-001**: `ANALYSIS.md` Scope and Golden Output Snippets lines 3-9 and 68-127; source was not installed or executed, and available examples are committed snippets rather than full stdout fixtures for every mode.
- **SA-002**: `ANALYSIS.md` Entry Points, Inputs/Outputs, Flow Map, and RT-006 lines 21-45, 61-65, 254-266, 342-350; `bin/license-checker:74-104` renderer selection, colorization, `--files`, `--out`, and stdout paths.
- **SA-003**: `ANALYSIS.md` Function Inventory lines 147-153 and `lib/index.js:471-627`; renderer helper functions and file-copy side effect.
- **SA-004**: `ANALYSIS.md` Input-Field Inventory lines 165-174 and Output Field Contract lines 213-232; output field order, sentinels, and output-affecting overrides.
- **SA-005**: `ANALYSIS.md` Source-Dependency Contracts lines 289-290 and Dependency Notes/Risk Map lines 278-300; `treeify`, `chalk`, and output byte-risk contracts.
- **SA-006**: `USE-CASES.md` F-006 block lines 50-57 and shared status/open questions lines 68-89; F-006 readiness, dependencies, evidence summary, and missing full-byte fixture limitation.
- **SA-007**: `README.md:24-67`; default tree and guessed-license marker examples.
- **SA-008**: `tests/test.js:37-43`, `70-104`, `609-688`; CSV, Markdown, tree/summary helper, and `asFiles` examples.
- **SA-009**: Source dependency contract for `mkdirp` from `ANALYSIS.md` lines 270-281 plus `bin/license-checker:96-104` and `lib/index.js:610-627`; directory creation before file writes.

## User Scenarios & Testing

### User Story 1 - Render default tree output (Priority: P1)

As a CLI user, I want the default command output to render the sorted package map as the same box-drawing tree produced by the source so that a normal `license-checker` invocation remains visually and structurally compatible.

**Why this priority**: Tree output is the default behavior when no structured-output flag is selected.

**Independent Test**: Provide a sorted package map with repository and license fields, render with no `--json`, `--csv`, `--markdown`, `--summary`, `--files`, or `--out`, and compare the tree glyphs and field order against the committed README snippets.

**Acceptance Scenarios**:

1. **Given** a sorted package map containing `cli@0.4.3` with `repository: http://github.com/chriso/cli` and `licenses: MIT`, **When** the default renderer is selected, **Then** the output includes the tree lines `├─ cli@0.4.3`, `│  ├─ repository: http://github.com/chriso/cli`, and `│  └─ licenses: MIT` in treeify-compatible form. [SA-002][SA-003][SA-007]
2. **Given** a package whose rendered license value is the guessed literal `MIT*`, **When** the default tree renderer emits it, **Then** the output preserves the exact `MIT*` literal rather than replacing or explaining the asterisk. [SA-004][SA-007]

---

### User Story 2 - Emit machine-readable JSON, CSV, and Markdown (Priority: P1)

As a user or automation script, I want JSON, CSV, and Markdown modes to preserve the source field order, row order, literals, and newline behavior so downstream consumers do not observe schema drift.

**Why this priority**: Structured outputs are the primary integration surface and are explicitly documented/tested by the source project.

**Independent Test**: Render a known sorted map in each mode and compare headers, first rows, indentation/newlines, and missing-value handling against the golden snippets and renderer contracts.

**Acceptance Scenarios**:

1. **Given** a sorted map whose first key is `abbrev@1.0.9`, `licenses` is `ISC`, and `repository` is `https://github.com/isaacs/abbrev-js`, **When** CSV mode renders without custom format, **Then** line 1 is exactly `"module name","license","repository"` and line 2 is exactly `"abbrev@1.0.9","ISC","https://github.com/isaacs/abbrev-js"`. [SA-001][SA-003][SA-008]
2. **Given** the same record, **When** Markdown mode renders without custom format, **Then** the row is exactly `[abbrev@1.0.9](https://github.com/isaacs/abbrev-js) - ISC`. [SA-001][SA-003][SA-008]
3. **Given** CLI JSON mode is selected, **When** records are rendered, **Then** the renderer uses two-space JSON indentation over the source field order and appends one newline to the formatted JSON string before any stdout `console.log`-equivalent line ending. [SA-002][SA-004]

---

### User Story 3 - Render summaries and custom-format tables (Priority: P2)

As a user, I want summary, CSV custom format, CSV component prefix, and Markdown custom format to match the source helpers so that reporting options remain compatible.

**Why this priority**: These modes are explicit user-facing report variants but depend on the core record map from P1.

**Independent Test**: Render records with a custom format object containing `name`, `description`, and `pewpew`, and render a summary from records with repeated license values.

**Acceptance Scenarios**:

1. **Given** custom format keys `name`, `description`, `pewpew` in that insertion order, **When** CSV custom-format mode renders, **Then** the header is exactly `"module name","name","description","pewpew"` and the first data row can be `"abbrev@1.0.9","abbrev","Like ruby's abbrev module, but in js","<<Should Never be set>>"`. [SA-003][SA-008]
2. **Given** the same custom format and `csvComponentPrefix` value `main-module`, **When** CSV renders, **Then** the header is exactly `"component","module name","name","description","pewpew"` and each row starts with `"main-module"`. [SA-003][SA-008]
3. **Given** the same custom format, **When** Markdown custom-format mode renders `abbrev@1.0.9`, **Then** the package row is exactly ` - **[abbrev@1.0.9](https://github.com/isaacs/abbrev-js)**` and custom fields follow as four-space-indented `- key: value` lines. [SA-003][SA-008]
4. **Given** multiple records with exact license strings, **When** summary mode renders, **Then** the system counts records by exact `licenses` string, sorts counts descending, and renders the count object with treeify-compatible tree glyphs. [SA-003][SA-004][SA-008]

---

### User Story 4 - Write output files and copy license files (Priority: P2)

As a CLI user, I want `--out` and `--files` side effects to match the source so that existing scripts relying on generated files keep working.

**Why this priority**: File side effects are explicitly supported CLI outputs and differ from stdout in newline behavior.

**Independent Test**: Render to a temporary output file and copy license files from a map with one valid `licenseFile` and one missing `licenseFile`; assert file content, names, and warnings.

**Acceptance Scenarios**:

1. **Given** `--out /tmp/licenses/licenses.csv` and CSV mode, **When** rendering completes, **Then** the parent directory is created as needed and the file receives exactly the CSV formatted string without an additional stdout `console.log` newline. [SA-002][SA-009]
2. **Given** `--files /tmp/lc` and a record key `foo` whose `licenseFile` points to an existing source license file, **When** file-copy mode runs, **Then** `/tmp/lc/foo-LICENSE.txt` is written with the source license-file contents. [SA-003][SA-008][SA-009]
3. **Given** `--files` and a record key with no usable `licenseFile`, **When** file-copy mode runs, **Then** the system emits the exact warning text `no license file found for: <moduleName>` for that key and continues processing other keys. [SA-003][SA-004]

### Edge Cases

- CSV values are wrapped in double quotes but the source does not implement general CSV escaping; the Go implementation MUST preserve this behavior rather than silently switching to RFC-compliant escaping for all fields. [SA-003][SA-004]
- Default CSV missing `licenses` or `repository` values render as empty quoted fields (`""`), demonstrated by the source helper test containing `"foo","",""`. [SA-008]
- Custom CSV reads `module[item]` for each custom key; upstream custom-format population is expected to supply defaults, and absent values must not be guessed by the renderer. [SA-003][SA-004]
- Summary equal-count ordering is not stabilized by source code beyond descending numeric count; the spec does not require deterministic tie ordering unless a downstream fixture pins a specific source-engine result. [SA-004][SA-006]
- Color output is terminal/environment-dependent through `chalk.supportsColor`; no-color output is the exact parity baseline unless a color test explicitly controls the color environment. [SA-005][SA-006]
- Full byte-exact stdout fixtures for JSON, full tree, summary, `--files`, `--out`, color, and error paths are not present in committed tests and were not generated because live source execution is prohibited. [SA-001][SA-006]

## Requirements

### Functional Requirements

- **FR-001**: The system MUST select renderers in this precedence after scan completion: JSON when `json` is truthy; else CSV when `csv` is truthy; else Markdown when `markdown` is truthy; else summary when `summary` is truthy; else default tree. [SA-002]
- **FR-002**: After formatting, the system MUST perform side effects in this precedence: if `files` is truthy, copy license files and do not emit the formatted renderer output; else if `out` is truthy, create the output parent directory and write the raw formatted string as UTF-8; else write the formatted string to stdout using a `console.log`-equivalent that appends its own line ending. [SA-002][SA-009]
- **FR-003**: Default tree and programmatic `asTree` MUST match `treeify.asTree(sorted, true)` behavior, including box-drawing glyphs shown in the README snippets (`├─`, `│`, `└─`) and traversal over object insertion order. [SA-003][SA-005][SA-007]
- **FR-004**: Programmatic `print(sorted)` MUST write exactly the `asTree(sorted)` result through a `console.log`-equivalent line-emitting operation. [SA-003]
- **FR-005**: JSON mode MUST format records with two-space indentation equivalent to `JSON.stringify(json, null, 2)`, preserving top-level key order and per-record field insertion order supplied by F-003/F-004/F-005, and MUST append exactly one newline to the formatted JSON string before stdout or `--out` handling. [SA-002][SA-004]
- **FR-006**: CSV mode without custom format MUST emit header fields in exact order `"module name","license","repository"`, with optional leading `"component"` only when `csvComponentPrefix` is truthy. [SA-003][SA-008]
- **FR-007**: CSV mode without custom format MUST emit each row in sorted-map key order as `"<key>","<licenses-or-empty>","<repository-or-empty>"`, with optional leading `"<csvComponentPrefix>"`; missing `licenses` and missing `repository` MUST render as empty strings. [SA-003][SA-008]
- **FR-008**: CSV mode with custom format MUST emit header fields in exact order: optional `"component"`, then `"module name"`, then `Object.keys(customFormat)` order; each data row MUST emit optional component prefix, the module key, then `module[item]` for each custom-format key in that same order. [SA-003][SA-008]
- **FR-009**: CSV rendering MUST join lines with `\n`, MUST NOT append a trailing newline in `asCSV`, and MUST wrap values in double quotes without general CSV escaping; only upstream `licenseText` CSV preprocessing may replace double quotes with single quotes and line breaks with spaces. [SA-003][SA-004]
- **FR-010**: Markdown mode without custom format MUST emit one row per package in key order as exact template `[<key>](<repository>) - <licenses>`, joined with `\n`, with no trailing newline in `asMarkDown`. [SA-003][SA-008]
- **FR-011**: Markdown mode with custom format MUST emit for each package an exact first line ` - **[<key>](<repository>)**`, followed by one line per custom-format key in `Object.keys(customFormat)` order using exact template `    - <customKey>: <value>`, joined with `\n`, with no trailing newline in `asMarkDown`. [SA-003][SA-008]
- **FR-012**: CLI Markdown mode MUST append exactly one newline to the `asMarkDown` result before stdout or `--out` handling. [SA-002][SA-003]
- **FR-013**: Summary mode MUST count records by exact `licenses` string, create license/count entries, sort them by descending count using the source comparator `b.count - a.count`, and render the resulting object with treeify-compatible `asTree(..., true)` behavior. [SA-003][SA-004][SA-005]
- **FR-014**: Summary equal-count ordering MUST be treated as source-unstable unless a committed fixture pins a specific order; the implementation MUST NOT claim deterministic lexical tie sorting as a source behavior. [SA-004][SA-006]
- **FR-015**: `--out` MUST create `path.dirname(out)` using mkdirp-compatible recursive directory creation and MUST write only the raw formatted string as UTF-8, without the additional stdout line ending. [SA-002][SA-009]
- **FR-016**: `--files`/programmatic `asFiles(json, outDir)` MUST create `outDir`, iterate `Object.keys(json)` order, test the record's `licenseFile` path exactly as supplied, and for each existing file write a file named exactly `<moduleName>-LICENSE.txt` under `outDir` with the source license-file contents. [SA-003][SA-008][SA-009]
- **FR-017**: `--files` MUST create parent directories for the computed output path before writing, preserving source behavior for module names that produce nested paths when joined with `outDir`. [SA-003][SA-009]
- **FR-018**: For each `--files` record lacking an existing `licenseFile`, the system MUST emit exact warning text `no license file found for: <moduleName>`. [SA-003]
- **FR-019**: Colorized CLI tree-key rewriting MUST occur only when `color` is truthy, `out` is falsy, and none of `csv`, `json`, or `markdown` is truthy; it MUST split each top-level key on `@` and rewrite it as `chalk.blue(keyParts[0]) + chalk.dim('@') + chalk.green(keyParts[1])`-compatible ANSI styling before rendering, preserving the source's first-two-parts behavior rather than inventing scoped-package repair logic. [SA-002][SA-005]
- **FR-020**: Sentinel colorization from upstream scan/filter behavior MUST be preserved when `color` is truthy: private-package `UNLICENSED` and fallback `UNKNOWN` may be wrapped with `chalk.bold.red`-compatible styling before rendering. No-color literals remain exactly `UNLICENSED` and `UNKNOWN`. [SA-004][SA-005]
- **FR-021**: The renderer input map MUST preserve JavaScript insertion-order semantics from the source: top-level package keys arrive lexically sorted; per-record fields appear in insertion order beginning with `licenses`, optional `private`, then conditional fields `repository`, `url`, `publisher`, `email`, author `url` assignment, `dependencyPath`, custom-format fields, `path`, `licenseFile`, `licenseText`, `copyright`, and `noticeFile`; if author `url` overwrites an existing `url`, no second `url` field is created. [SA-004]
- **FR-022**: The output must surface the private-package unconditional override: whenever an input record has `private` truthy and it has not been removed by F-005, the rendered `licenses` field MUST be exactly `UNLICENSED` in no-color modes regardless of package license metadata. [SA-004]
- **FR-023**: The output must surface repository URL normalization supplied by F-003: rendered `repository` values MUST reflect source replacements `git+ssh://git@` -> `git://`, `git+https://github.com` -> `https://github.com`, `git://github.com` -> `https://github.com`, `git@github.com:` -> `https://github.com/`, and removal of trailing `.git`. [SA-004][SA-008]
- **FR-024**: The output must surface author URL override semantics: when both `json.url.web` and included `json.author.url` exist, the rendered `url` field MUST be the author URL because source assignment overwrites the earlier `url.web` value. [SA-004]
- **FR-025**: The output must surface license-field and file/readme-derived values exactly as supplied by F-004, including literals `UNKNOWN`, `UNLICENSED`, `Undefined`, guessed `*` values such as `MIT*`, `Custom: <value>`, `Public Domain`, and raw valid SPDX expressions. [SA-004][SA-007]
- **FR-026**: The output must surface path fields exactly as supplied by F-003/F-004: `dependencyPath` is present only for `--unknown`, `path` is included when not suppressed and available, and `licenseFile`/`noticeFile` use relative paths only when `relativeLicensePath` was enabled upstream. [SA-004]
- **FR-027**: The output must surface custom-format inclusion semantics: a custom-format value of `false` suppresses that field; otherwise existing string metadata is used, and missing included fields receive the custom-format default before rendering. [SA-004][SA-008]
- **FR-028**: The output must surface `licenseText` preprocessing semantics: in non-CSV custom-format output it is trimmed file content; in CSV custom-format output it has double quotes replaced with single quotes, all line breaks replaced with spaces, and leading/trailing whitespace trimmed. [SA-004]
- **FR-029**: The output must surface `noticeFile` when F-004 supplies it; no dedicated source test exists, so parity is source-code-derived rather than fixture-backed. [SA-004][SA-006]

### Output Contract

#### Golden baseline snippets

Full live stdout fixtures were not captured. The following committed snippets are the byte baselines available to this slice and MUST be reproduced where the same input conditions are used:

Default tree excerpt:

```text
├─ cli@0.4.3
│  ├─ repository: http://github.com/chriso/cli
│  └─ licenses: MIT
├─ glob@3.1.14
│  ├─ repository: https://github.com/isaacs/node-glob
│  └─ licenses: UNKNOWN
└─ yui-lint@0.1.1
   ├─ licenses: BSD
      └─ repository: http://github.com/yui/yui-lint
```

Guessed-license marker:

```text
└─ debug@2.0.0
   ├─ repository: https://github.com/visionmedia/debug
   └─ licenses: MIT*
```

CSV default:

```text
"module name","license","repository"
"abbrev@1.0.9","ISC","https://github.com/isaacs/abbrev-js"
```

CSV custom format:

```text
"module name","name","description","pewpew"
"abbrev@1.0.9","abbrev","Like ruby's abbrev module, but in js","<<Should Never be set>>"
```

CSV custom format with component prefix:

```text
"component","module name","name","description","pewpew"
"main-module","abbrev@1.0.9","abbrev","Like ruby's abbrev module, but in js","<<Should Never be set>>"
```

Markdown default:

```text
[abbrev@1.0.9](https://github.com/isaacs/abbrev-js) - ISC
```

Markdown custom format package row:

```text
 - **[abbrev@1.0.9](https://github.com/isaacs/abbrev-js)**
```

`--files` copied-file example:

```text
foo-LICENSE.txt
```

Missing `--files` warning template:

```text
no license file found for: <moduleName>
```

#### Field order & presence

- Top-level package key order: lexical sort supplied before rendering; package whitelist/blacklist rebuilds keep filtered key order from prior maps. [SA-004]
- Per-record field order for JSON/tree traversal: `licenses`; optional `private`; then conditional `repository`, `url`, `publisher`, `email`, author `url` assignment, `dependencyPath`, custom-format fields in custom key order, `path`, `licenseFile`, `licenseText`, `copyright`, `noticeFile`; author `url` overwrites an existing `url` value rather than creating a duplicate field. [SA-004]
- CSV default fields: optional `component`, `module name`, `license`, `repository`. [SA-003]
- CSV custom fields: optional `component`, `module name`, then every `customFormat` key in insertion order. [SA-003]
- Markdown default fields: key, repository URL, licenses. [SA-003]
- Markdown custom fields: key/repository package row, then every `customFormat` key in insertion order. [SA-003]
- `licenseFile` and `noticeFile` are conditionally present only when detected and not suppressed by custom format; relative vs absolute path is determined upstream by `relativeLicensePath`. [SA-004]

#### Sentinels & literals

- Initial license sentinel: `UNKNOWN`.
- Private package no-color override: `UNLICENSED`.
- Classifier/record literals that render unchanged when supplied: `Undefined`, guessed `*` values such as `MIT*`, `Apache*`, `BSD*`, `CC0-1.0*`, `Custom: <value>`, `Public Domain`, and raw valid SPDX expressions.
- Default custom-format example value: `<<Should Never be set>>`.
- Missing default CSV `licenses`/`repository`: empty string inside quotes, e.g. `"foo","",""`.
- Missing `--files` license file warning: `no license file found for: <moduleName>`.

#### Sort/aggregation order

- Package records render in `Object.keys(sorted)` order, where upstream sorted maps are lexical by package key after dependency flattening and before filters. [SA-004]
- Summary aggregation counts exact `licenses` strings, converts counts to `{ license, count }` entries, sorts by descending count, inserts into a new object in that sorted order, and renders through treeify. Equal-count order is not specified. [SA-003][SA-004]
- `--files` iterates `Object.keys(json)` order. [SA-003]

#### Whitespace and stdout-vs-file behavior

- `asTree` and `asSummary`: return the treeify string; stdout adds one `console.log`-equivalent line ending; `--out` writes the raw treeify string with no added stdout line ending. [SA-002][SA-003]
- CLI JSON: `JSON.stringify(json, null, 2) + '\n'`; stdout then adds one additional `console.log`-equivalent line ending; `--out` writes the one-newline formatted JSON string. [SA-002]
- `asCSV`: joins lines with `\n` and appends no trailing newline; stdout adds one line ending; `--out` writes no trailing newline beyond the joined lines. [SA-002][SA-003]
- `asMarkDown`: joins rows with `\n` and appends no trailing newline; CLI Markdown appends one newline; stdout adds one additional line ending; `--out` writes the one-newline CLI Markdown string. [SA-002][SA-003]
- `--files`: bypasses stdout/file rendering of the formatted string and writes copied license-file artifacts instead. [SA-002][SA-003]

### Algorithm Fidelity

#### Renderer selection and side-effect rule table

1. Initialize formatted output as empty string. [SA-002]
2. If color eligibility is true (`color && !out && !(csv || json || markdown)`), rewrite each top-level package key by splitting on `@` and applying blue/dim/green-compatible ANSI styling to `keyParts[0]`, `@`, and `keyParts[1]`, respectively. [SA-002][SA-005]
3. If `json` is true, set formatted output to two-space JSON plus `\n`. [SA-002]
4. Else if `csv` is true, set formatted output to `asCSV(json, customFormat, csvComponentPrefix)`. [SA-002][SA-003]
5. Else if `markdown` is true, set formatted output to `asMarkDown(json, customFormat) + '\n'`. [SA-002][SA-003]
6. Else if `summary` is true, set formatted output to `asSummary(json)`. [SA-002][SA-003]
7. Else, set formatted output to `asTree(json)`. [SA-002][SA-003]
8. If `files` is true, run `asFiles(json, files)` and do not write formatted output to stdout or `out`. [SA-002][SA-003]
9. Else if `out` is true, create the parent directory and write formatted output as UTF-8. [SA-002][SA-009]
10. Else, write formatted output through stdout line-emitting behavior. [SA-002]

#### Replicated source dependencies

- **treeify `^1.1.0`**: The target MUST match `treeify.asTree(obj, true)` for default tree and summary tree output, including glyphs, indentation, and object traversal behavior visible in README and tests. If no Go dependency is adopted, the reimplementation must be fixture-tested against the snippets above and any committed golden fixtures. [SA-005][SA-007][SA-008]
- **chalk `^2.4.1`**: The target MUST preserve source color eligibility and ANSI role mapping for tree keys (`blue`, `dim`, `green`) and sentinels (`bold.red`) when color support is explicitly enabled. Because terminal-dependent color support is a known gap, no-color output is the default byte baseline. [SA-002][SA-005][SA-006]
- **mkdirp `^0.5.1`**: The target MUST recursively create output directories before writing `--out` and `--files` artifacts, treating existing directories as successful. [SA-009]

### Key Entities

- **Rendered Package Map**: Ordered mapping from package key `name@version` to an output record. Top-level key order and record field insertion order are part of the observable output contract for JSON/tree renderers.
- **Output Record**: Per-package field set containing `licenses` plus optional metadata and file fields. Values are supplied by F-003/F-004/F-005 but rendered by this slice without schema reshaping.
- **Custom Format**: Ordered object whose keys define additional CSV/Markdown/JSON/tree fields and default values; `false` suppresses a field upstream.
- **Rendered Output String**: The raw string produced by JSON, CSV, Markdown, summary, or tree formatting before stdout or `--out` line-ending differences are applied.
- **Copied License Artifact**: File produced by `asFiles` at `<outDir>/<moduleName>-LICENSE.txt` from the source `licenseFile` contents.

## Success Criteria

### Measurable Outcomes

- **SC-001**: For every committed golden snippet in this spec, the target renderer produces byte-identical output for equivalent in-memory package maps, including exact field order, quotes, glyphs, indentation, and line joins. [SA-001][SA-007][SA-008]
- **SC-002**: For JSON, CSV, Markdown, summary, default tree, `--out`, and stdout modes, tests assert the specified trailing-newline differences between raw formatter output, file output, and stdout output. [SA-002][SA-003]
- **SC-003**: CSV tests cover default fields, missing `licenses`/`repository`, custom-format field order, component prefix, and the source's lack of general CSV escaping. [SA-003][SA-008]
- **SC-004**: `--files` tests cover recursive output directory creation, `foo-LICENSE.txt` naming, copied license-file content, and exact missing-file warning text. [SA-003][SA-008][SA-009]
- **SC-005**: Summary tests cover exact-license-string aggregation and descending count sorting, and they explicitly avoid asserting deterministic equal-count ordering unless a fixture pins it. [SA-003][SA-004]
- **SC-006**: Differential-parity baseline: for the committed source snippets and any analyzer-provided golden fixtures available to this slice, target output is byte-identical across no-color tree, JSON, CSV, Markdown, summary, `--out`, and `--files` behaviors covered by fixtures, including edge inputs with missing fields, private packages, guessed licenses, relative license paths, custom-format defaults, and license-text CSV preprocessing. [SA-001][SA-004][SA-006]
- **SC-007**: Parity limitation is explicit: passing this spec does not imply live full-byte parity for JSON/full tree/summary/color/files modes beyond committed snippets, because no live source execution was performed and those full stdout fixtures are absent. [SA-001][SA-006]

## Assumptions

- The input package map has already been produced, sorted, filtered, and enriched according to F-003, F-004, and F-005; this slice does not rescan dependencies or reclassify licenses.
- No-color output is the default byte-identity baseline. Color behavior is required only when tests explicitly control color eligibility and expected ANSI output.
- The target Go implementation may choose native Go libraries or local reimplementations, but observable behavior must match the source contracts above.
- JSON object and record ordering must be deliberately preserved in Go despite Go map iteration being unordered by default.
- Full live differential comparison against the installed Node.js source remains a manual, out-of-flow release activity and is not a Phase 1 or builder gate.

## Open Questions & Parity Limitations

- Full byte-exact stdout fixtures for JSON, full tree, summary, `--files`, `--out`, color, and error paths are not present in committed tests and were not generated because source execution was prohibited. [SA-001][SA-006]
- There is no dedicated source test for `NOTICE` files or author URL overriding `url.web`; these requirements are source-code-derived and must be covered by target tests using constructed inputs. [SA-004][SA-006]
- Exact terminal/environment color support equivalent to `chalk.supportsColor` is terminal-dependent; no-color output remains the stable parity baseline unless the test matrix pins color behavior. [SA-005][SA-006]
