# OVERVIEW.md — license-checker Analyzer Overview

## App Tasks

1. Parse CLI/programmatic options and normalize runtime defaults.  
2. Load dependency tree from a start path with depth and dev/prod constraints.  
3. Build a flattened package map with canonical `name@version` keys.  
4. Resolve license information from metadata, SPDX expressions, README/license files.  
5. Apply license/package/private filters and compliance gates (`failOn`/`onlyAllow`).  
6. Emit output in selected format and optionally persist artifacts to disk.

## Use Cases

### UC-01: Audit all dependencies for license inventory
- User runs `license-checker` in a project folder.
- System scans dependencies and prints tree output.
- Evidence: `bin/license-checker:65-104`, `lib/index.js:261-477`, README `13-57`.

### UC-02: Generate machine-readable compliance artifacts
- User runs with `--json`, `--csv`, `--markdown`, or `--summary`, optionally `--out`.
- System formats data and writes to stdout or file.
- Evidence: `bin/license-checker:84-104`, `lib/index.js:508-591`, README `77-85`.

### UC-03: Enforce license policy gates in CI
- User defines forbidden licenses (`--failOn`) or allowlist (`--onlyAllow`).
- System exits `1` on first violating package.
- Evidence: `lib/index.js:437-457`, `tests/failOn-test.js:7-25`, `tests/test.js:233-288`.

### UC-04: Focus scans to relevant scope
- User restricts to production/dev/direct deps, package lists, or excludes.
- System applies filters before returning output.
- Evidence: `lib/args.js:79-83`, `lib/index.js:49-51`, `313-435`, `tests/packages-test.js:8-45`.

### UC-05: Capture license evidence files for legal review
- User passes `--files <dir>`.
- System copies detected license files per package into output folder.
- Evidence: `bin/license-checker:96-98`, `lib/index.js:610-627`, `tests/test.js:669-686`.

### UC-06: Extend output schema with custom format
- User provides inline `customFormat` or `--customPath` JSON.
- System includes selected fields/default values, optional `licenseText`/`copyright`.
- Evidence: `lib/index.js:95-106`, `181-221`, `264-266`, `593-608`, `tests/test.js:429-478`.

## Requirements Task List

- [ ] **R-001 (F-001)** Implement option parsing/default logic compatible with current flag model (`lib/args.js:7-99`).
- [ ] **R-002 (F-002)** Recreate dependency discovery using start path + depth/dev behavior (`lib/index.js:267-309`).
- [ ] **R-003 (F-003)** Rebuild flattening logic with duplicate guard and metadata extraction (`lib/index.js:27-110`, `235-259`).
- [ ] **R-004 (F-004)** Preserve license-resolution precedence: SPDX -> heuristics -> files (`lib/index.js:112-172`; `lib/license.js`; `lib/license-files.js`).
- [ ] **R-005 (F-005)** Reproduce filtering semantics (`unknown`, `onlyunknown`, `exclude`, package include/exclude, private excludes) (`lib/index.js:313-435`).
- [ ] **R-006 (F-005)** Preserve policy-fail process exit behavior for `failOn` and `onlyAllow` (`lib/index.js:437-457`; `bin/license-checker:54-57`).
- [ ] **R-007 (F-006)** Recreate all output serializers and mode precedence (tree/json/csv/markdown/summary) (`bin/license-checker:84-94`; `lib/index.js:475-591`).
- [ ] **R-008 (F-006)** Support output persistence: `--out` file writes and `--files` per-license exports (`bin/license-checker:96-102`; `lib/index.js:610-627`).
- [ ] **R-009 (F-007)** Support custom format ingestion from object or JSON file, including parse-failure behavior (`lib/index.js:264-266`, `593-608`).
- [ ] **R-010 (F-008)** Decide whether unused `Stack` utility is retained, removed, or explicitly deprecated (`lib/stack.js:6-44`; no references found).

## Feature-to-Task Mapping

| Functionality ID | Title | Primary requirements |
|---|---|---|
| F-001 | CLI options model and normalization | R-001 |
| F-002 | Dependency graph loading/traversal control | R-002 |
| F-003 | Flatten graph to package map | R-003 |
| F-004 | License detection + normalization pipeline | R-004 |
| F-005 | Filtering and policy enforcement | R-005, R-006 |
| F-006 | Output rendering and persistence | R-007, R-008 |
| F-007 | Custom format/config model | R-009 |
| F-008 | Ancillary Stack utility handling | R-010 |

## Open Questions

1. Should the rewrite keep `--version` exit code `1` for strict compatibility (`bin/license-checker:49-52`)?
2. Should currently under-documented flags (`--files`, `--markdown`) remain supported as-is?
3. Should invalid `--customPath` become hard error instead of silent `Error`-as-config flow?
4. Should `onlyAllow` matching remain substring-based (current behavior), or become strict SPDX-expression-aware matching?
