# Feature Specification: Argument parsing, defaults, and custom format loading

**Created**: 2026-06-22  
**Status**: Draft  
**Input**: User description: "rewrite Sourcebot project `davglass/license-checker` into Go; assigned slice F-002 — Argument parsing, defaults, and custom format loading"

## Scope

This specification covers only functionality slice **F-002** for the pinned source module `davglass/license-checker` version `25.0.1` at commit `de6e9a42513aa38a58efc6b202ee5281ed61f486`.

**In scope**

- Parse the single-command `license-checker [flags]` argument surface declared in `lib/args.js`.
- Preserve short aliases `-v` -> `--version` and `-h` -> `--help`.
- Apply source defaults for `color`, `start`, `relativeLicensePath`, and `direct`.
- Remove `nopt` `argv` metadata from parsed options before defaults are returned.
- Load custom format JSON from `customPath` in `checker.init` and expose the parsed object through `options.customFormat`.
- Preserve custom-format field inclusion/default semantics that affect downstream records and renderers.

**Out of scope**

- CLI process lifecycle, help/version printing, stderr, and exit codes owned by F-001/F-007.
- Dependency graph reconstruction beyond passing `start`, `dev`, and `depth` options to the scanner, owned by F-003.
- License detection and file precedence, owned by F-004.
- Filtering/restriction policy behavior, owned by F-005.
- Final renderer byte output except for the custom-format field order and field-presence contract consumed by F-006.

**Source evidence IDs**

- **SA-001**: `lib/args.js:7-43`, `65-98` — `nopt` known flags, shorts, clean/default/parse exports.
- **SA-002**: `lib/index.js:264-270` — `customPath` assignment to `customFormat` and scanner `depth` option.
- **SA-003**: `lib/index.js:56-58`, `95-105` — custom-format `include` rule, string-only source-field copy, default insertion, and `false` suppression.
- **SA-004**: `lib/index.js:593-608` — `parseJson(jsonPath)` success/error behavior.
- **SA-005**: `tests/test.js:389-476` — defaults, `direct`, and custom-format integration assertions.
- **SA-006**: `tests/test.js:690-716` — JSON parsing success and error assertions.
- **SA-007**: `ANALYSIS.md` Entry Points, Input-Field Inventory, Source-Dependency Contracts, RT-001/RT-006, and Risk/Open Questions for F-002.
- **SA-008**: `USE-CASES.md` F-002 block and shared open questions.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Parse CLI flags with source defaults (Priority: P1)

A CLI user invokes `license-checker` with any supported flag combination and the Go target returns the same normalized option values that the source parser would pass into the scan pipeline.

**Why this priority**: Every downstream feature depends on normalized options. Incorrect `start`, `direct`, or format flags changes the scan root, traversal depth, and output behavior.

**Independent Test**: Can be tested by calling the parser with argument vectors and asserting the normalized option object without running a package scan.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** no parsed options are supplied to defaults, **When** defaults are applied, **Then** `color` equals the runtime equivalent of `chalk.supportsColor`, `start` equals the current working directory, `relativeLicensePath` equals `false`, and `direct` equals `Infinity`. [SA-001][SA-005]
2. **Given** parsed options contain `direct: true`, **When** defaults are applied, **Then** `direct` is normalized to numeric depth `0`. [SA-001][SA-005]
3. **Given** parsed options omit `direct`, **When** defaults are applied, **Then** `direct` is normalized to `Infinity`. [SA-001][SA-005]
4. **Given** an argument vector uses `-v` or `-h`, **When** the parser runs, **Then** the resulting options are equivalent to using `--version` or `--help`, respectively. [SA-001][SA-007]
5. **Given** a parsed option object includes `json`, `markdown`, or `csv` as truthy, **When** defaults are applied, **Then** `color` is forced to `false`. [SA-001]
6. **Given** the cooked argument vector contains either `--color` or `--no-color`, **When** `has('color')` is evaluated, **Then** it returns `true` because source `has(a)` recognizes both positive and `--no-` cooked forms. [SA-001]

---

### User Story 2 - Select scan root and direct dependency depth (Priority: P2)

A user controls the package-tree scan root and whether only direct dependencies are traversed by using parser-level `start` and `direct` values.

**Why this priority**: These values are the F-002 inputs consumed by the F-003 dependency scanner. They must be normalized before scanner invocation.

**Independent Test**: Can be tested with a fake scanner boundary by asserting that normalized options produce scanner inputs `{ start, depth }` equal to source behavior.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** no `--start` value is provided, **When** the scan pipeline receives parsed options, **Then** `options.start` is the current working directory. [SA-001][SA-007]
2. **Given** `--start <path>` is provided, **When** the scan pipeline receives parsed options, **Then** `options.start` is the provided value without additional F-002 normalization. [SA-001][SA-007]
3. **Given** `--direct` is present, **When** `checker.init` builds scanner options, **Then** scanner `depth` equals `0`. [SA-001][SA-002]
4. **Given** `--direct` is absent, **When** `checker.init` builds scanner options, **Then** scanner `depth` equals `Infinity`. [SA-001][SA-002]

---

### User Story 3 - Load and apply custom format fields (Priority: P3)

A CLI or programmatic user supplies a custom format, and package records receive exactly the configured additional fields, default values, and suppressions used by source renderers.

**Why this priority**: Custom format is the bridge between parsing and output identity for CSV/Markdown/custom record fields.

**Independent Test**: Can be tested by passing a small package metadata node and a custom format object into record-building logic, then asserting field presence, order, and values.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** `customPath` points to valid JSON, **When** `checker.init` begins, **Then** `options.customFormat` is set to the parsed JSON object before scanning. [SA-002][SA-004][SA-006]
2. **Given** `customPath` points to a missing file, invalid JSON file, or a non-string path, **When** custom JSON is parsed, **Then** `parseJson` returns an `Error` object rather than throwing. [SA-004][SA-006]
3. **Given** custom format contains `"pewpew": "<<Should Never be set>>` and package metadata does not contain string field `pewpew`, **When** a record is built, **Then** the record field `pewpew` equals `<<Should Never be set>>`. [SA-003][SA-005]
4. **Given** custom format contains a field set to `false`, **When** a record is built, **Then** that field is suppressed from source-driven/custom default insertion. [SA-003]
5. **Given** custom format keys are `name`, `description`, `pewpew`, **When** downstream custom CSV uses the custom format, **Then** the header order is exactly `"module name","name","description","pewpew"`. [SA-003][SA-005][SA-007]

### Edge Cases

- `defaults(undefined)` MUST behave as if parsing the current process arguments and then applying defaults. [SA-001][SA-005]
- `color` MUST remain the runtime color-support default unless explicitly supplied or disabled by `json`, `markdown`, or `csv`; this slice does not add extra summary/out color rules beyond source defaults. [SA-001][SA-007]
- `relativeLicensePath` MUST be boolean-coerced with source truthiness rules: truthy -> `true`, missing/falsey -> `false`. [SA-001]
- `customPath` parse failures return `Error` values; no runtime-captured CLI diagnostics for these failures are available in Analyzer evidence. [SA-004][SA-006][SA-008]
- Exact `nopt` edge behavior for malformed values, unknown flags, and `--no-<flag>` cooked metadata is a parity risk and must not be guessed from live execution. [SA-007][SA-008]

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST provide a single-command argument parser for `license-checker [flags]` with no subcommand dispatch in this slice. [SA-001][SA-007]
- **FR-002**: The parser MUST recognize the following known options with source-equivalent types: `production:Boolean`, `development:Boolean`, `json:Boolean`, `csv:Boolean`, `csvComponentPrefix:String`, `markdown:Boolean`, `out:path`, `unknown:Boolean`, `onlyunknown:Boolean`, `version:Boolean`, `color:Boolean`, `start:String`, `help:Boolean`, `relativeLicensePath:Boolean`, `exclude:String`, `customPath:path`, `customFormat:Object`, `files:path`, `summary:Boolean`, `failOn:String`, `onlyAllow:String`, `direct:Boolean`, `packages:String`, `excludePackages:String`, and `excludePrivatePackages:Boolean`. [SA-001][SA-007]
- **FR-003**: The parser MUST recognize short aliases `-v` as `--version` and `-h` as `--help`. [SA-001][SA-007]
- **FR-004**: The raw parser MUST preserve source-equivalent `nopt ^4.0.1` parsing behavior for declared option types and aliases; if the Go implementation uses a different parser, any observed incompatibility for malformed values, unknown flags, boolean negation, path coercion, or argv metadata MUST be recorded as an explicit conformance decision before release. [SA-001][SA-007][SA-008]
- **FR-005**: The cleaned parser result MUST remove the parser metadata field equivalent to `argv` before returning normalized options. [SA-001]
- **FR-006**: The parser utility `has(a)` MUST return true when source-equivalent cooked arguments contain either `--<a>` or `--no-<a>`. [SA-001]
- **FR-007**: Defaults MUST set `color` to the runtime equivalent of `chalk.supportsColor` only when parsed `color` is `undefined`/unset. [SA-001][SA-005]
- **FR-008**: Defaults MUST force `color` to `false` when `json`, `markdown`, or `csv` is truthy. [SA-001]
- **FR-009**: Defaults MUST set `start` to the current working directory when parsed `start` is falsey; otherwise they MUST preserve the provided `start` string. [SA-001][SA-005][SA-007]
- **FR-010**: Defaults MUST set `relativeLicensePath` to the boolean coercion of its parsed value. [SA-001][SA-007]
- **FR-011**: Defaults MUST convert truthy `direct` to numeric scanner depth `0` and falsey/missing `direct` to `Infinity`. [SA-001][SA-002][SA-005]
- **FR-012**: `checker.init` MUST, when `options.customPath` is truthy, set `options.customFormat` to the return value of `parseJson(options.customPath)` before scanner options are built. [SA-002][SA-004]
- **FR-013**: `checker.init` MUST pass scanner options with `depth: options.direct`; F-002 normalization is the sole source for translating `--direct` to `0` or `Infinity`. [SA-001][SA-002]
- **FR-014**: `parseJson` MUST return `new Error('did not specify a path')` or a target-equivalent Error when `jsonPath` is not a string. [SA-004][SA-006]
- **FR-015**: `parseJson` MUST read string paths as UTF-8 JSON files, return the parsed object on success, and return the caught Error object on file-read or JSON-parse failure without throwing. [SA-004][SA-006]
- **FR-016**: Custom-format inclusion MUST treat a property as included when `customFormat` is absent or `customFormat[property] !== false`; exactly `false` suppresses the property. [SA-003]
- **FR-017**: When `customFormat` is present, record construction MUST iterate `Object.keys(customFormat)` in source-equivalent key order and handle each key before later `path`, `licenseFile`, `licenseText`, `copyright`, and `noticeFile` insertions. [SA-003][SA-007]
- **FR-018**: For each included custom-format key, if package metadata has a truthy string value for that key, the record MUST use the package string value. [SA-003][SA-005]
- **FR-019**: For each included custom-format key, if package metadata lacks a truthy value for that key, the record MUST use the configured custom-format default value. [SA-003][SA-005]
- **FR-020**: For each included custom-format key where package metadata has a truthy non-string value, the implementation MUST preserve source behavior by not assigning the non-string value as a custom field and not substituting the default in that branch. [SA-003]
- **FR-021**: Programmatic callers MUST be able to supply `customFormat` directly without `customPath`, and that object MUST drive the same include/default behavior as a parsed custom JSON file. [SA-003][SA-005]

### Output Contract *(mandatory for any slice that produces observable output)*

This slice does not own final stdout/stderr renderers, but it produces normalized options and custom-format-controlled record fields consumed by F-003/F-006. These contracts are therefore mandatory for parity.

- **Golden baseline snippets**:
  - Defaults: source tests assert `defaults(undefined).color === chalk.supportsColor`, `defaults(undefined).start === path.resolve(path.join(__dirname, '../'))`, missing `direct` -> `Infinity`, and `direct: true` -> `0`. [SA-005]
  - Custom CSV field order/value fixture:
    ```text
    "module name","name","description","pewpew"
    "abbrev@1.0.9","abbrev","Like ruby's abbrev module, but in js","<<Should Never be set>>"
    ```
    [SA-005][SA-007]
  - Custom CSV with component prefix fixture:
    ```text
    "component","module name","name","description","pewpew"
    "main-module","abbrev@1.0.9","abbrev","Like ruby's abbrev module, but in js","<<Should Never be set>>"
    ```
    [SA-005][SA-007]
  - Example custom JSON file fields and defaults:
    ```json
    {
      "name": "",
      "version": "",
      "description": "",
      "licenses": "",
      "copyright": "",
      "licenseFile": "none",
      "licenseText": "none",
      "licenseModified": "no"
    }
    ```
    [SA-004][SA-006]
- **Normalized option field order/presence**: no output schema order is guaranteed for the parser object beyond source-equivalent field names and values. The `argv` metadata field MUST be absent from the cleaned/parsed result. [SA-001]
- **Default literals and sentinel values**:
  - `color`: runtime equivalent of `chalk.supportsColor`; forced literal `false` when `json`, `markdown`, or `csv` is truthy. [SA-001]
  - `start`: `process.cwd()` equivalent when missing/falsey. [SA-001]
  - `relativeLicensePath`: boolean literal `true` or `false` from `!!parsed.relativeLicensePath`. [SA-001]
  - `direct`: numeric `0` when present/truthy; `Infinity` when absent/falsey. [SA-001][SA-002]
  - `parseJson` non-string path error message: `did not specify a path`. [SA-004]
- **Custom-format field order & presence**:
  - Iterate custom-format keys in source-equivalent JSON/object key order. [SA-003]
  - Include a custom-format field unless its configured value is exactly `false`. [SA-003]
  - Use package string value when truthy and string; otherwise use configured default when package value is missing/falsey; suppress exact-`false` fields. [SA-003]
- **Sort/aggregation order**: F-002 does not sort package records. It MUST preserve custom-format key order for downstream custom renderers. [SA-003][SA-007]
- **Whitespace contract**: F-002 `parseJson` consumes UTF-8 JSON and returns structured data; it does not emit text or normalize whitespace. Downstream CSV/Markdown whitespace is owned by F-006, but depends on the field order above. [SA-004][SA-007]

### Algorithm Fidelity *(mandatory when the slice transforms values)*

#### Argument parsing and defaults ordered rule table

1. Call source-equivalent `nopt(known, shorts, args || process.argv)` for raw parsing. [SA-001]
2. Clean parsed options by deleting `argv`. [SA-001]
3. If defaults receive `undefined`, first compute the cleaned parse result from current process arguments. [SA-001][SA-005]
4. If `parsed.color === undefined`, set `parsed.color = chalk.supportsColor` equivalent. [SA-001][SA-005]
5. If `parsed.json || parsed.markdown || parsed.csv`, set `parsed.color = false`. [SA-001]
6. Set `parsed.start = parsed.start || process.cwd()` equivalent. [SA-001][SA-005]
7. Set `parsed.relativeLicensePath = !!parsed.relativeLicensePath`. [SA-001]
8. If `parsed.direct` is truthy, set `parsed.direct = 0`; otherwise set `parsed.direct = Infinity`. [SA-001][SA-002][SA-005]
9. Return the mutated parsed object. [SA-001]

#### Custom JSON loading ordered rule table

1. If `typeof jsonPath !== 'string'`, return `new Error('did not specify a path')`. [SA-004][SA-006]
2. Initialize an empty-result object. [SA-004]
3. Try to read `jsonPath` as UTF-8 text and `JSON.parse` it. [SA-004]
4. On success, return the parsed JSON value. [SA-004][SA-006]
5. On file-read or JSON-parse error, catch the error and return it. [SA-004][SA-006]

#### Custom-format field ordered rule table

1. Define `include(property)` as true when `options.customFormat === undefined || options.customFormat[property] !== false`. [SA-003]
2. If no `options.customFormat` exists, skip custom-format key iteration. [SA-003]
3. If `options.customFormat` exists, iterate `Object.keys(options.customFormat)` in source-equivalent order. [SA-003]
4. For each key, if `include(key)` and package metadata has a truthy value at that key, assign it only when `typeof json[key] === 'string'`. [SA-003]
5. Otherwise, if `include(key)`, assign `options.customFormat[key]` as the record default. [SA-003]
6. If `include(key)` is false, assign nothing for that key. [SA-003]
7. No-match/fallback branch: absent custom format adds no custom fields; included custom keys with truthy non-string package values add no value and do not fall back to the default in source control flow. [SA-003]

**Source dependencies whose behavior must be replicated**

- `nopt ^4.0.1`: preserve known option typing, aliases, boolean negation/cooked option behavior relevant to `has(a)`, path-typed options (`out`, `customPath`, `files`), and metadata deletion after parsing. If not using `nopt`, the Go parser MUST be validated against source-derived vectors and unresolved edge behavior must remain an explicit conformance risk. [SA-001][SA-007][SA-008]
- `chalk ^2.4.1`: this slice uses only `chalk.supportsColor` as the default source for `color`; ANSI rendering is owned by F-006. [SA-001][SA-007]

### Key Entities *(include if feature involves data)*

- **ParsedOptions**: normalized option object passed from CLI parsing into `checker.init`; key attributes include all known flags, no `argv` metadata, and defaults for `color`, `start`, `relativeLicensePath`, and `direct`.
- **CustomFormat**: JSON/object map whose keys define extra output fields, default values, and exact-`false` suppressions.
- **ScannerOptions boundary**: object passed to dependency scanning containing `dev`, `log`, and `depth`, where `depth` comes directly from normalized `direct`.

## Open Questions & Conformance Risks

- **OQ-001**: Exact `nopt ^4.0.1` behavior for unknown flags, malformed typed values, and all `--no-<flag>` cases was not exhaustively captured by committed tests; the Go target must either match verified source-derived vectors or document a non-identical conformance row. [SA-007][SA-008]
- **OQ-002**: Analyzer did not execute the live source CLI, so no new runtime captures may be invented for parser edge cases or custom-path failures. [SA-007][SA-008]
- **OQ-003**: If `parseJson` returns an Error and `checker.init` assigns it to `options.customFormat`, downstream object-key behavior follows source control flow but has no dedicated committed golden output; keep this as a parity-risk fixture rather than adding guessed diagnostics. [SA-002][SA-004][SA-006]

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Parser unit tests cover every known option name, both short aliases, and the absence of `argv` metadata in cleaned results. [SA-001]
- **SC-002**: Defaults tests prove `color`, `start`, `relativeLicensePath`, and `direct` match source behavior for missing, falsey, and truthy inputs. [SA-001][SA-005]
- **SC-003**: Custom JSON tests prove valid JSON returns an object with fields such as `licenseModified: "no"` and invalid/missing/non-string paths return Error objects. [SA-004][SA-006]
- **SC-004**: Custom-format record tests prove package string values override defaults, missing fields receive defaults, exact `false` suppresses fields, and field order matches `Object.keys(customFormat)`. [SA-003][SA-005]
- **SC-005**: For source-derived fixtures in this slice, the target normalized options and custom-format-controlled fields are behavior-identical to the source across default parsing, `--direct`, `--customPath`, direct `customFormat`, and parse-error cases. [SA-001][SA-004][SA-005][SA-006]
- **SC-006**: Any non-identical behavior versus `nopt ^4.0.1` for parser edge cases is documented in the conformance matrix before the rewrite is marked complete; no unverified live source execution is used to fill gaps. [SA-007][SA-008]
- **SC-007**: For the custom-format fixture with keys `name`, `description`, `pewpew`, downstream custom CSV header and first `abbrev@1.0.9` row are byte-identical to the golden baseline snippets in this spec. [SA-005][SA-007]

## Assumptions

- The Go target can represent source `Infinity` scanner depth either as an explicit unbounded-depth sentinel or an equivalent internal value, as long as the scanner boundary preserves absent-`direct` semantics. [SA-001][SA-002]
- The final CLI process behavior for help/version/error exits is provided by F-001/F-007 and consumes this slice's parsed options.
- The final package-record construction and renderers are provided by F-003/F-006 and consume this slice's custom-format contract.
- Because Analyzer did not execute the live source system and Sourcebot indexing was unavailable, exact parser behavior not visible in source/tests remains a conformance-risk item rather than a guessed requirement. [SA-007][SA-008]
