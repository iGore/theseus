# Feature Specification: Dependency Scan and Metadata Flattening

**Created**: 2026-06-22  
**Status**: Draft  
**Input**: User description: "rewrite Sourcebot project `davglass/license-checker` into Go; assigned slice F-003 Dependency scan and metadata flattening"

## Scope

This specification covers only F-003, `dependency-scan-metadata-flattening`, for the pinned source module `davglass/license-checker` version `25.0.1` at commit `de6e9a42513aa38a58efc6b202ee5281ed61f486`. The slice starts after F-001/F-002 have supplied parsed options and any loaded custom format, and ends with a sorted package-record map keyed by `name@version`.

**In scope**:

- Calling the dependency-tree provider with source-equivalent scan options.
- Flattening the returned package tree into package records.
- Applying production/development pruning during flattening.
- Preserving source insertion order for metadata fields owned by this slice.
- Normalizing repository URLs.
- Populating author/url fields, `path`, and `dependencyPath` under the source conditions.
- Applying F-003 sentinel handling that occurs during the sorted pass: initial `UNKNOWN`, private-package `UNLICENSED`, missing-license `UNKNOWN`, and `--unknown` guessed-license rewrite.

**Out of scope**:

- CLI invocation and option parsing details beyond already-parsed options from F-001/F-002.
- License classification rule fidelity, license-file precedence, and `licenseFile`/`noticeFile` content extraction owned by F-004, except where their later fields affect shared record order.
- Filtering/restriction/compliance decisions owned by F-005.
- Renderer byte formatting and stdout/file newline behavior owned by F-006.
- Diagnostics and process exit behavior owned by F-007.

## Source Grounding

The Analyzer used the Sourcebot reconstruction workflow; the local Sourcebot endpoint was unavailable, so concrete line evidence was recorded from a read-only clone of the pinned GitHub commit. No live source-system execution was performed.

| ID | Evidence |
|---|---|
| SA-001 | Scope and pinned source: `ANALYSIS.md:3-9`; shared ready status: `USE-CASES.md:23-30`. |
| SA-002 | Expected inputs: parsed options, `read-installed` tree nodes, and package metadata fields: `ANALYSIS.md:51-59`. |
| SA-003 | `flatten(options)` and `exports.init(options, callback)` roles: `ANALYSIS.md:129-142`; assigned source refs `lib/index.js:27-110`, `235-309` within the broader reconstructed ranges `lib/index.js:27-258`, `261-309`. |
| SA-004 | Input-field inventory for `start`, `direct`, `name`, `version`, `private`, `repository`, `url`, `author`, `path`, custom format: `ANALYSIS.md:158-174`. |
| SA-005 | Output record field order, sentinels, and lexical sorting: `ANALYSIS.md:213-224`. |
| SA-006 | Golden committed examples for `abbrev@1.0.9` normalized repository output: `ANALYSIS.md:94-119`; `tests/test.js:19-45`. |
| SA-007 | Domain map and flow map: `ANALYSIS.md:246-266`. |
| SA-008 | `read-installed ~4.0.3` dependency contract and risk: `ANALYSIS.md:283-286`, `292-301`. |
| SA-009 | Path-related committed assertions: `tests/test.js:480-530`; summarized in `ANALYSIS.md:170`. |
| SA-010 | Requirements mapping RT-002/RT-003: `ANALYSIS.md:342-358`. |
| SA-011 | Open questions: Sourcebot direct query unavailable; no live byte captures; no dedicated test for author URL overwrite: `ANALYSIS.md:304-309`; `USE-CASES.md:85-89`. |

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Scan an installed npm project into stable package records (Priority: P1)

As a CLI or API user, I need a scan of an installed npm project to return a stable map of discovered packages keyed by `name@version`, so downstream filtering and rendering work on the same package identity as the source tool.

**Why this priority**: This is the core value of the migration slice. F-004 through F-007 depend on this map existing and being source-compatible.

**Independent Test**: Can be tested by invoking the scan/flatten API with a fixture equivalent to the source repository and asserting package keys, record presence, sorted order, and known metadata values without exercising CLI rendering.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** parsed options with `start` pointing at an installed package tree and no `production`, `development`, or `direct` override, **When** the scan starts, **Then** the dependency provider is called with `read(options.start, { dev: true, log: debugLog, depth: options.direct }, cb)`-equivalent semantics. [SA-002][SA-003][SA-008]
2. **Given** a dependency tree containing the committed source test dependency `abbrev@1.0.9`, **When** records are flattened and sorted, **Then** the result contains key `abbrev@1.0.9` and preserves the source-observed values that render as `"abbrev@1.0.9","ISC","https://github.com/isaacs/abbrev-js"` in CSV and `[abbrev@1.0.9](https://github.com/isaacs/abbrev-js) - ISC` in Markdown. [SA-006]
3. **Given** two package nodes that produce the same `name@version` key through dependency recursion, **When** the second node is encountered, **Then** the existing record is retained and recursion does not reprocess that key. [SA-003]
4. **Given** a package node missing `name` or missing `version`, **When** flattening finishes processing that node, **Then** no record for that incomplete node remains in the output map. [SA-004]

---

### User Story 2 - Preserve source metadata and repository normalization (Priority: P2)

As a user inspecting package records, I need repository, URL, author, custom metadata, and path fields to match the source tool so output renderers show the same package facts.

**Why this priority**: Metadata identity is observable in CSV, Markdown, tree, JSON, and programmatic API usage.

**Independent Test**: Can be tested with synthetic `read-installed` package nodes containing repository, URL, author, custom format, and path fields; the flattened record can be inspected directly.

**Acceptance Scenarios** *(Gherkin — Given / When / Then)*:

1. **Given** a package node with `repository.url` equal to `git://github.com/example/pkg.git`, **When** the record is flattened, **Then** `repository` is `https://github.com/example/pkg` after ordered source normalization. [SA-004][SA-005]
2. **Given** a package node with both `url.web` and `author.url`, and the `url` field is included, **When** the record is flattened, **Then** `url` is first eligible for `json.url.web` and is then overwritten by `json.author.url`; if no earlier `url` field was inserted, author URL insertion occurs at the author-field point in source order. [SA-004][SA-005][SA-011]
3. **Given** a custom format containing `name`, `description`, and `pewpew`, where `name` and `description` exist as strings on the package node and `pewpew` is absent, **When** the record is flattened, **Then** the record contains the package `name`, package `description`, and default `pewpew` value exactly as the source custom-format tests expect. [SA-004][SA-006]

---

### User Story 3 - Honor scan-scope and path options (Priority: P3)

As a user narrowing a scan or troubleshooting unknown licenses, I need source-equivalent `production`, `development`, `direct`, `path`, and `dependencyPath` behavior.

**Why this priority**: These controls affect which packages appear and which location fields users can inspect, but they depend on P1/P2 record creation.

**Independent Test**: Can be tested with fixture dependency trees containing `extraneous`, `root`, nested dependencies, package paths, and `unknown` mode, without invoking renderers.

**Acceptance Scenarios** *(Gherkin — Given / When / Then)*:

1. **Given** parsed options with `production: true`, **When** a package node has `extraneous: true`, **Then** flattening omits that node and its record. [SA-003][SA-004][SA-010]
2. **Given** parsed options with `development: true`, **When** a package node is neither `extraneous` nor `root`, **Then** flattening omits that node and its record. [SA-003][SA-004][SA-010]
3. **Given** F-002 has resolved `direct` to `0` for `--direct` or `Infinity` when absent, **When** F-003 invokes the dependency provider, **Then** that value is passed unchanged as the provider `depth` option. [SA-004][SA-008]
4. **Given** a package node with a string `path` and the `path` field is included, **When** the record is flattened, **Then** the record contains `path` as an absolute package path under the scan root, matching the source path assertions. [SA-004][SA-009]
5. **Given** options include `unknown: true`, **When** a package record is created, **Then** `dependencyPath` is populated from the package node `path`; if `unknown` is absent or false, `dependencyPath` is not populated by this slice. [SA-004][SA-005]

### Edge Cases

- A duplicate `name@version` key MUST short-circuit recursion and preserve the first record, preventing circular dependency loops. [SA-003]
- Package nodes with missing `name` or `version` MUST be deleted from the final flattened data, even though the temporary key construction is `json.name + '@' + json.version`. [SA-004]
- `read-installed` graph semantics are high risk and MUST NOT be replaced by a naive filesystem walk without compatibility tests for nested dependencies, `extraneous`, `root`, `path`, `readme`, and duplicate/dedupe behavior. [SA-008]
- Author URL overwrite is source-explicit but has no dedicated committed fixture; tests SHOULD include a synthetic fixture rather than treating the behavior as optional. [SA-011]
- No live byte-output captures were generated; tests for this slice MUST use committed examples and source-derived fixtures, not newly invented live-output baselines. [SA-001][SA-011]

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept parsed F-001/F-002 options including `start`, `direct`, `production`, `development`, `unknown`, `customFormat`, `relativeLicensePath`, and `color` as inputs to the scan/flatten step. [SA-002][SA-004]
- **FR-002**: System MUST invoke a dependency-tree provider with source-equivalent `read-installed ~4.0.3` semantics using `start`, `dev`, `log`, and `depth` inputs; the output MUST provide package node fields used by this slice (`name`, `version`, `private`, `repository`, `url`, `author`, `path`, `dependencies`, `extraneous`, and `root`). [SA-002][SA-008]
- **FR-003**: System MUST set provider option `dev` to `true` by default and to `false` when either `production` or `development` is true. [SA-003][SA-008]
- **FR-004**: System MUST pass `options.direct` through as provider `depth` without changing the F-002-resolved value. [SA-004][SA-008]
- **FR-005**: System MUST create package record keys as exact string concatenation `name + "@" + version`. [SA-004][SA-005]
- **FR-006**: System MUST initialize each new record with `licenses` inserted first and value `UNKNOWN`. [SA-005]
- **FR-007**: System MUST mark `private` as `true` immediately after `licenses` only when the package node has truthy `private`. [SA-004][SA-005]
- **FR-008**: System MUST omit a package node from the flattened data when `production` is true and the node has truthy `extraneous`. [SA-003][SA-010]
- **FR-009**: System MUST omit a package node from the flattened data when `development` is true and the node is not `extraneous` and not `root`. [SA-003][SA-010]
- **FR-010**: System MUST return early without reprocessing or recursing when a `name@version` key already exists in the flattened data. [SA-003]
- **FR-011**: System MUST normalize `repository.url` only when `repository` is included, `json.repository` is an object, and `json.repository.url` is a string; normalization MUST follow the ordered table in the Output Contract. [SA-004][SA-005]
- **FR-012**: System MUST set `url` from `json.url.web` only when `url` is included and `json.url` is an object. [SA-004][SA-005]
- **FR-013**: System MUST set author-derived fields in source order: `publisher` from `author.name` when included, `email` from `author.email` when included, and `url` from `author.url` when included; `author.url` MUST overwrite any earlier `url` value from `json.url.web`. [SA-004][SA-005][SA-011]
- **FR-014**: System MUST add `dependencyPath` from `json.path` only when `options.unknown` is true. [SA-004][SA-005]
- **FR-015**: System MUST apply custom-format field inclusion exactly as source: a field is included unless `customFormat[field] === false`; when included and `json[field]` is a string, copy that string; when included and `json[field]` is absent or falsey, write the configured default; when `json[field]` exists but is non-string, do not write the default from this branch. [SA-004][SA-006]
- **FR-016**: System MUST add `path` from `json.path` only when `path` is included, `json.path` is truthy, and `json.path` is a string. [SA-004][SA-009]
- **FR-017**: System MUST recurse through `json.dependencies` using each child dependency node and carry forward `customFormat`, `production`, `development`, `basePath`, `unknown`, and the original options context needed by downstream slices. [SA-003][SA-007]
- **FR-018**: System MUST delete the record for a node when `json.name` or `json.version` is missing before returning flattened data. [SA-004]
- **FR-019**: System MUST build the returned sorted map by iterating `Object.keys(data).sort()`-equivalent lexical key order before later F-005 filters. [SA-005]
- **FR-020**: During the sorted pass, System MUST overwrite a private record's `licenses` value with `UNLICENSED` while preserving field order. [SA-004][SA-005]
- **FR-021**: During the sorted pass, System MUST overwrite any falsey `licenses` value with `UNKNOWN`. [SA-005]
- **FR-022**: When `options.unknown` is true, System MUST rewrite any non-`UNKNOWN` license string containing `*` to `UNKNOWN`; this requirement only covers the F-003 sentinel override, not F-004 classifier fidelity. [SA-004][SA-005]
- **FR-023**: System MUST preserve the committed source-observed golden package identity for `abbrev@1.0.9`: the flattened/sorted records must contain that key with repository value `https://github.com/isaacs/abbrev-js` in the source fixture context. [SA-006]
- **FR-024**: System MUST treat `read-installed` behavior as a replicated source dependency; if not using an equivalent library, the reimplementation MUST provide compatibility tests for nested dependency tree shape, duplicate key avoidance, `extraneous`/`root` pruning, and path propagation. [SA-008]

### Output Contract *(mandatory for any slice that produces observable output)*

- **Golden baseline**: This slice does not own renderer output and no live output was captured. It owns the package-record map that downstream renderers consume. Committed source examples that must remain reproducible through downstream renderers are:
  - CSV first data row: `"abbrev@1.0.9","ISC","https://github.com/isaacs/abbrev-js"` (`tests/test.js:37-39`; `ANALYSIS.md:94-99`). [SA-006]
  - Markdown first row: `[abbrev@1.0.9](https://github.com/isaacs/abbrev-js) - ISC` (`tests/test.js:41-43`; `ANALYSIS.md:111-119`). [SA-006]
  - Path assertions: every `path` starts with the scan-root prefix in the source fixture; license-file path behavior belongs to F-004/F-006 but shares the same `relativeLicensePath` base-path decision (`tests/test.js:480-530`). [SA-009]

- **Field order & presence**: Records MUST preserve source JavaScript insertion-order semantics. For fields owned or conditionally touched by this slice, exact emission order is:

  | Order | Field | Presence condition | Value contract |
  |---:|---|---|---|
  | 1 | `licenses` | Always inserted for every new record | Initial value `UNKNOWN`; may later be overwritten by F-004 classification and F-003 sorted-pass sentinels. |
  | 2 | `private` | Only when package node `private` is truthy | Boolean `true`; later causes `licenses = UNLICENSED`. |
  | 3 | `repository` | `repository` included and `json.repository.url` is a string | Ordered normalized repository URL. |
  | 4 | `url` | `url` included and `json.url` is an object | `json.url.web`; may be overwritten by author URL without changing order. |
  | 5 | `publisher` | `author` object, `publisher` included, `author.name` truthy | `json.author.name`. |
  | 6 | `email` | `author` object, `email` included, `author.email` truthy | `json.author.email`. |
  | 7 | `url` | `author` object, `url` included, `author.url` truthy | `json.author.url`; overwrites earlier `url`; if no earlier `url` existed, insertion occurs here. |
  | 8 | `dependencyPath` | `options.unknown` true | `json.path`. |
  | 9 | custom-format fields | `options.customFormat` present; iterate `Object.keys(customFormat)` in source order; field not set to `false` | String package value when `json[field]` is a string; otherwise configured default only when source value is absent/falsey. If a custom key names a field already inserted earlier, assignment overwrites the value without moving the field. If a custom key names `path`, it can insert `path` here and the later source `path` assignment overwrites without moving it. |
  | 10 | `path` | `path` included, `json.path` is a string, and `path` was not already inserted by an earlier custom-format assignment | `json.path`. |
  | 11+ | F-004 fields | Added later by license-file/notice logic | `licenseFile`, `licenseText`, `copyright`, `noticeFile`; value rules are owned by F-004 but order after `path` must remain available for the shared record. |

- **Sentinels & literals**:
  - Initial license sentinel: `UNKNOWN`.
  - Private override sentinel: `UNLICENSED`.
  - Missing/falsey license sorted-pass sentinel: `UNKNOWN`.
  - `--unknown` guessed-license rewrite sentinel: `UNKNOWN` for non-`UNKNOWN` license values containing `*`.
  - Exact literals above are the no-color baseline; if `options.color` is true, the source applies `colorizeString` to sorted-pass sentinel rewrites, and byte-level ANSI behavior is owned by F-006 color rendering.
  - Other license literals (`Undefined`, guessed `*` values, `Custom: <value>`, `Public Domain`, and raw valid SPDX expressions) are produced by F-004 and must not be reordered by this slice.

- **Sort/aggregation order**: Top-level output package keys MUST be lexically sorted by exact key string before F-005 filtering. This slice performs no aggregation.

- **Whitespace contract**: F-003 returns an in-memory map/object and MUST NOT emit stdout, write files, or append newlines. Renderer whitespace, stdout `console.log` behavior, and `--out` differences are F-006 responsibilities.

- **Repository normalization exact ordered rule table**:

  | Step | Source operation | Example input after previous step | Example output |
  |---:|---|---|---|
  | 1 | Replace `git+ssh://git@` with `git://` | `git+ssh://git@github.com/example/pkg.git` | `git://github.com/example/pkg.git` |
  | 2 | Replace `git+https://github.com` with `https://github.com` | `git+https://github.com/example/pkg.git` | `https://github.com/example/pkg.git` |
  | 3 | Replace `git://github.com` with `https://github.com` | `git://github.com/example/pkg.git` | `https://github.com/example/pkg.git` |
  | 4 | Replace `git@github.com:` with `https://github.com/` | `git@github.com:example/pkg.git` | `https://github.com/example/pkg.git` |
  | 5 | Remove trailing `.git` matching `/\.git$/` | `https://github.com/example/pkg.git` | `https://github.com/example/pkg` |
  | 6 | Fallback | Any string not matching a prior rule | Preserve the string after prior rule applications. |

  Golden observed normalized value: source fixture record `abbrev@1.0.9` renders repository as `https://github.com/isaacs/abbrev-js`. [SA-006]

### Algorithm Fidelity *(mandatory when the slice transforms values)*

- **Dependency-provider contract**: The implementation MUST replicate the source reliance on `read-installed ~4.0.3`: a logical npm installed dependency tree, not a naive directory walk. Required node fields for this slice are `name`, `version`, `private`, `repository`, `url`, `author`, `path`, `dependencies`, `extraneous`, and `root`; downstream slices also rely on `license`, `licenses`, and `readme`. [SA-002][SA-008]
- **Provider option table**:

  | Condition | `dev` option | `depth` option | Notes |
  |---|---:|---:|---|
  | Default | `true` | `options.direct` | F-002 resolves absent `direct` to `Infinity`. |
  | `production` true | `false` | `options.direct` | Flatten additionally omits `extraneous` nodes. |
  | `development` true | `false` | `options.direct` | Flatten additionally omits non-`extraneous`, non-`root` nodes. |
  | `direct` true from CLI | unchanged by this slice | `0` from F-002 | Provider limits depth. |

- **Production/development pruning table**:

  | Ordered check | Condition | Action |
  |---:|---|---|
  | 1 | Existing output already has `name@version` key | Return current data unchanged. |
  | 2 | `options.production && json.extraneous` | Return current data unchanged. |
  | 3 | `options.development && !json.extraneous && !json.root` | Return current data unchanged. |
  | 4 | Otherwise | Insert record and recurse dependencies. |

- **No-match/fallback behavior**: Repository strings that do not match the normalization rules remain unchanged; custom-format fields with `false` are excluded; package nodes with missing name/version are removed from final data. [SA-004][SA-005]

### Key Entities *(include if feature involves data)*

- **Package node**: A dependency-tree node returned by the source-equivalent provider. Key attributes for this slice include `name`, `version`, `private`, `repository`, `url`, `author`, `path`, `dependencies`, `extraneous`, and `root`. [SA-002][SA-007]
- **Output record**: Mutable package metadata object keyed by `name@version`, initially `{ licenses: 'UNKNOWN' }`, with insertion order preserved for downstream JSON/tree/CSV/Markdown rendering. [SA-005][SA-007]
- **Sorted package map**: Lexically sorted map/object built from flattened records before F-005 filtering and F-006 rendering. [SA-005]
- **Custom format**: Optional object from F-002 that controls field inclusion, exclusion, and default values during record construction. [SA-004]

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: For the committed source fixture represented by `tests/test.js:19-45`, the target scan returns a sorted package map containing `abbrev@1.0.9` with repository value `https://github.com/isaacs/abbrev-js`; when passed to compatible renderers, the CSV and Markdown snippets in the golden baseline are byte-identical. [SA-006]
- **SC-002**: Unit tests cover every repository normalization table row, including the fallback row, and assert exact output strings. [SA-004][SA-005]
- **SC-003**: Unit tests cover record field insertion order for `licenses`, `private`, `repository`, `url`, `publisher`, `email`, author `url`, `dependencyPath`, custom-format fields, and `path`. [SA-005]
- **SC-004**: Scan-scope tests cover default, `production`, `development`, and `direct` cases, including `extraneous` and `root` pruning behavior. [SA-008][SA-010]
- **SC-005**: Path tests match source assertions: every populated `path` is an absolute package path rooted under `start`, and `dependencyPath` appears only with `unknown: true`. [SA-009]
- **SC-006**: For migration parity, the target package-record map is source-compatible across the edge fixture set for this slice: duplicate keys/circular guards, missing name/version deletion, private package `UNLICENSED`, custom-format defaults/exclusions, author URL overwrite, repository normalization, path propagation, and production/development pruning. [SA-004][SA-008][SA-011]

## Assumptions

- F-001/F-002 provide already-parsed options with defaults, including `start`, `direct`, `color`, `relativeLicensePath`, and optional `customFormat`; this slice does not parse CLI arguments.
- F-004 will supply license classification and license-file/notice field value rules while preserving the shared record-order positions reserved here.
- F-005 will apply `exclude`, `packages`, `excludePackages`, `excludePrivatePackages`, `failOn`, and `onlyAllow` after this slice returns the sorted map.
- F-006 will own renderer-specific whitespace, stdout/file behavior, CSV/Markdown/tree/JSON bytes, and color rendering; this slice only preserves the data and field order those renderers consume.
- Because the Analyzer did not execute the live source system, no new live-output captures are assumed. Any final live differential run remains manual/out-of-flow.
- The Go implementation may replace `read-installed` only if it reproduces the documented source dependency contract and passes compatibility fixtures for this slice.
