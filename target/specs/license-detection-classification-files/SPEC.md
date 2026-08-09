# Feature Specification: License Detection, Classification, and License File Precedence

**Created**: 2026-06-22  
**Status**: Draft  
**Input**: User description: "rewrite Sourcebot project `davglass/license-checker` into Go; assigned slice F-004 — License detection, classification, and license file precedence"

## Scope

This specification covers only F-004 from `USE-CASES.md` and implements Analyzer requirement RT-004: reproducing `davglass/license-checker` license detection, classification, and license-file precedence for package records produced by F-003. It is grounded in pinned source commit `de6e9a42513aa38a58efc6b202ee5281ed61f486`, package version `25.0.1`. [SA-001]

**In scope**:

- Classifying package `license` / `licenses`, README text, and license-file text with the exact ordered SPDX/regex/fallback classifier. [SA-002]
- Detecting candidate license files with exact basename precedence, case-insensitive matching, and extension-insensitive matching. [SA-003]
- Applying the source package-record license flow where README may fill missing package license data and license-file text may override absent, `UNKNOWN`, or `Custom:` license values without any CLI flag. [SA-004]
- Preserving the relevant `licenses`, `licenseFile`, `licenseText`, `copyright`, and `noticeFile` record contracts needed by downstream renderers. [SA-005]

**Out of scope**:

- Dependency graph scanning and base metadata flattening owned by F-003.
- Filtering/restriction/license-policy decisions owned by F-005.
- Tree/JSON/CSV/Markdown/files rendering and stdout/file newline behavior owned by F-006, except where this slice supplies values consumed by those renderers.
- Live source execution or new byte-output captures; only committed snippets/tests from Analyzer evidence may be used. [SA-006]

## Source Evidence

| ID | Evidence |
|---|---|
| SA-001 | `ANALYSIS.md` Scope lines 3-9; Requirements Task List RT-004 lines 342-348; `USE-CASES.md` F-004 lines 32-39. |
| SA-002 | `ANALYSIS.md` Value-Transformation Tables lines 178-200; `lib/license.js:1-83`; `tests/license.js:10-181`. |
| SA-003 | `ANALYSIS.md` Value-Transformation Tables lines 202-204; `lib/license-files.js:3-29`; `tests/license-files-test.js:10-87`. |
| SA-004 | `ANALYSIS.md` Input-Field Inventory lines 168-172; Flow Map lines 256-260; `lib/index.js:112-230`; `tests/test.js:313-323`. |
| SA-005 | `ANALYSIS.md` Output Field Contract lines 215-223; Input-Field Inventory lines 170-173. |
| SA-006 | `ANALYSIS.md` Golden Output Snippets lines 68-92 and Open Questions lines 304-309; `USE-CASES.md` Open Questions lines 85-89. |
| SA-007 | `ANALYSIS.md` Source-Dependency Contracts line 287 and Dependency Notes lines 270-277. |

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Classify license strings and text exactly (Priority: P1)

As a user scanning npm packages, I need license metadata and license/readme text classified into the same values as the source tool so that downstream filters and renderers receive source-compatible license values.

**Why this priority**: Classification is the core behavior of this slice and every later license policy/output feature depends on its exact values.

**Independent Test**: Can be tested by invoking the classifier directly with the committed `tests/license.js` examples and verifying exact returned values.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** classifier input `MIT`, **When** license classification runs, **Then** the result is the raw string `MIT` unchanged because SPDX parsing succeeds. [SA-002]
2. **Given** classifier input `Apache-2.0 OR ISC OR MIT`, **When** license classification runs, **Then** the result is the raw string `Apache-2.0 OR ISC OR MIT` unchanged because SPDX expression parsing succeeds. [SA-002]
3. **Given** classifier text containing `ermission is hereby granted, free of charge, to any`, **When** license classification runs, **Then** the result is `MIT*`. [SA-002]
4. **Given** classifier text containing `GNU GENERAL PUBLIC LICENSE Version 2`, **When** license classification runs, **Then** the result is `GPL-2.0*`. [SA-002]
5. **Given** classifier text containing `SEE LICENSE IN LICENSE.md`, **When** license classification runs, **Then** the result is `Custom: LICENSE.md`. [SA-002]
6. **Given** classifier input that matches no SPDX expression and no ordered source rule, **When** license classification runs, **Then** the result is `null`. [SA-002]

---

### User Story 2 - Detect license files by source precedence (Priority: P2)

As a user scanning installed packages, I need license-like files discovered in the same order as the source tool so that the selected file and guessed license are source-compatible.

**Why this priority**: License files are a source of guessed licenses and file paths; wrong precedence can change visible output and copied-file behavior.

**Independent Test**: Can be tested by feeding directory filename lists into the license-file detector and comparing exact selected filename order to `tests/license-files-test.js` expectations.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** directory filenames containing `LICENSE.md`, `COPYING`, and `README.txt`, **When** license-file detection runs, **Then** `LICENSE.md` is selected before `COPYING` and `README.txt` because `LICENSE` has higher precedence. [SA-003]
2. **Given** directory filenames containing both `LICENSE-MIT` and `LICENSE-APACHE`, **When** license-file detection runs, **Then** only the first matching `LICENSE-*` filename for that pattern is pushed. [SA-003]
3. **Given** directory filenames with mixed case and extensions such as `licence.txt`, **When** license-file detection runs, **Then** matching is case-insensitive and extension-insensitive. [SA-003]
4. **Given** directory filenames with no license-precedence match, **When** license-file detection runs, **Then** no license file is selected for this slice. [SA-003]

---

### User Story 3 - Apply package/readme/file license precedence into records (Priority: P3)

As a programmatic caller or CLI user, I need each package record to receive license fields through the same package/readme/license-file flow so that output values such as `UNKNOWN`, `MIT*`, `Custom: <value>`, and license-file paths match the source behavior.

**Why this priority**: Integration with F-003 records makes the classifier observable in JSON/tree/CSV/Markdown outputs and in later filters.

**Independent Test**: Can be tested by constructing package metadata records with package license fields, README text, directory files, and license-file contents, then verifying the resulting record fields before renderer-specific formatting.

**Acceptance Scenarios** *(Gherkin format — Given / When / Then)*:

1. **Given** a package record has no package `license`/`licenses` but README text classifies to a license, **When** F-004 applies license detection, **Then** the README-derived classification fills the package record license value. [SA-004]
2. **Given** a package record has a missing license, `UNKNOWN`, or a license value beginning with `Custom:`, and a selected license file classifies as MIT text, **When** F-004 applies license-file content detection, **Then** `licenses` becomes `MIT*` without requiring any flag. [SA-004]
3. **Given** a package directory has a selected license file and a notice-like file, **When** F-004 applies file extraction, **Then** the record may include `licenseFile`, `licenseText`, `copyright`, and `noticeFile` in the source insertion order for downstream renderers. [SA-005]

### Edge Cases

- Undefined or falsey classifier input MUST return `Undefined`, not `UNKNOWN` or `null`. [SA-002]
- The classifier MUST remove only the first newline before regex matching. [SA-002]
- The BSD source-code regex appears after the broader BSD regex and can be shadowed; the Go rewrite MUST preserve this ordered behavior rather than optimizing it away. [SA-002]
- A valid SPDX expression MUST pass through unchanged even if later regex rules could also match substrings. [SA-002]
- License-file detection MUST ignore filename extension and use uppercased basename matching. [SA-003]
- License-file detection MUST push only the first filename matching each precedence pattern. [SA-003]
- Full byte fixture for the license-file-over-custom-url path was not captured because live source execution is prohibited; the committed assertion expects the resulting license value `MIT*`. [SA-004] [SA-006]

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST implement the F-004 license classifier as the complete ordered rule table in **Algorithm Fidelity**, returning exact string literals, raw valid SPDX input strings, or `null` for no match. [SA-002]
- **FR-002**: The system MUST use SPDX expression parsing semantics equivalent to source dependency `spdx-expression-parse ^3.0.0`: valid SPDX identifiers/expressions are accepted and returned unchanged; invalid expressions fall through to regex rules. [SA-002] [SA-007]
- **FR-003**: The system MUST classify package `license` / `licenses` metadata according to source behavior: use `json.license || json.licenses`; classify strings/objects; map arrays through classification. [SA-004]
- **FR-004**: The system MUST use README content as a license source when no package license fields provide a license value. [SA-004]
- **FR-005**: The system MUST detect license files with the exact precedence table in **Algorithm Fidelity**, matching on uppercased basename with extension ignored. [SA-003]
- **FR-006**: The system MUST add selected license-file path information to package records as `licenseFile` when a license file is found, preserving downstream relative-path handling ownership for F-003/F-006. [SA-005]
- **FR-007**: The system MUST perform the source license-file override without any CLI flag: when the current license value is absent, contains `UNKNOWN`, or starts with `Custom:`, selected license-file content MUST be classified and allowed to overwrite `licenses`. [SA-004]
- **FR-008**: The system MUST preserve exact license sentinel/literal values emitted by this slice: `UNKNOWN`, `Undefined`, guessed values ending in `*`, `Custom: <capture>`, `Public Domain`, raw valid SPDX expressions, and `null` fallback. [SA-002] [SA-005]
- **FR-009**: The system MUST preserve `licenseText` extraction as the selected license-file content for downstream renderers, with CSV-specific newline/quote transformations owned by F-006. [SA-004] [SA-005]
- **FR-010**: The system MUST preserve `noticeFile` detection as a conditional package-record field when a NOTICE file is found; relative-path conversion remains owned by F-003/F-006. [SA-005]
- **FR-011**: The system MUST NOT invent live-output fixtures or require live source execution for this slice; tests must use Analyzer-recorded snippets and committed source test expectations. [SA-006]

### Output Contract *(mandatory for any slice that produces observable output)*

- **Golden baseline**:
  - README default tree excerpt includes exact license field lines `└─ licenses: UNKNOWN` and `└─ licenses: MIT*`; these establish observable sentinels/guessed marker values for downstream tree output. [SA-006]
  - CSV committed snippet includes exact row `"abbrev@1.0.9","ISC","https://github.com/isaacs/abbrev-js"`, establishing observable SPDX pass-through value `ISC`. [SA-006]
  - Committed assertion for the license-file-over-custom-url fixture expects exact license value `MIT*`; no full byte fixture was captured. [SA-004] [SA-006]
- **Field order & presence**:
  1. `licenses` is inserted first in each output record by F-003 initial record creation and is the primary field this slice overwrites. [SA-005]
  2. `licenseFile` is conditionally present after `path` when a license file is detected. [SA-005]
  3. `licenseText` and `copyright` are conditionally present after `licenseFile` when file extraction populates them. [SA-005]
  4. `noticeFile` is conditionally present after `copyright` when a notice file is found. [SA-005]
- **Sentinels & literals**: exact values are `UNKNOWN`, `UNLICENSED` when supplied by the private-package override in F-003, `Undefined`, guessed literals `ISC*`, `MIT*`, `BSD*`, `BSD-Source-Code*`, `WTFPL*`, `Apache*`, `CC0-1.0*`, `GPL-<version>*`, `LGPL-<version>*`, `Public Domain`, `Custom: <capture>`, raw valid SPDX expressions, and `null`. [SA-002] [SA-005]
- **Sort/aggregation order**: top-level package keys are sorted lexically by F-003 before downstream filtering; within this slice, classifier and license-file rules are ordered exactly as listed in **Algorithm Fidelity**. [SA-003] [SA-005]
- **Whitespace contract**: classifier preprocessing removes only the first newline from truthy input before regex matching. Renderer-level trailing newline/stdout-vs-file behavior is owned by F-006. `licenseText` carries selected license-file content forward; CSV-specific newline and quote changes are owned by F-006. [SA-002] [SA-005]

### Algorithm Fidelity *(mandatory when the slice transforms values)*

#### License classifier exact ordered rule table

The classifier MUST evaluate rules in this order and stop at the first match. [SA-002]

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

**SPDX parser dependency behavior to preserve**: `spdx-expression-parse ^3.0.0` accepts valid SPDX identifiers/expressions and throws on invalid input; this slice uses only success/failure and returns original input unchanged on success. Worked examples include `MIT`, `LGPL-2.0`, `Apache-2.0`, `BSD-2-Clause`, `(GPL-2.0+ WITH Bison-exception-2.2)`, `LGPL-2.0 OR (ISC AND BSD-3-Clause+)`, and `Apache-2.0 OR ISC OR MIT`. [SA-002] [SA-007]

#### License-file precedence table

The detector MUST uppercase the basename, ignore the extension, evaluate patterns in this order, and push only the first matching filename for each pattern. [SA-003]

| Order | Pattern |
|---:|---|
| 1 | `^LICENSE$` |
| 2 | `^LICENSE\-\w+$` |
| 3 | `^LICENCE$` |
| 4 | `^LICENCE\-\w+$` |
| 5 | `^COPYING$` |
| 6 | `^README$` |
| 7 | Fallback/no match: select no file for that pattern and return an empty selection when no patterns match |

### Key Entities *(include if feature involves data)*

- **License value**: A value assigned to the `licenses` field: raw SPDX expression, classifier literal, guessed `*` literal, `Custom: <capture>`, `Public Domain`, `UNKNOWN`, `UNLICENSED`, `Undefined`, or `null`. [SA-002] [SA-005]
- **Package record**: Mutable output record keyed by `name@version` from F-003; this slice updates `licenses` and conditionally adds license/notice file fields. [SA-004] [SA-005]
- **License-file candidate**: Directory file whose basename matches one license-file precedence pattern after uppercasing and extension removal. [SA-003]

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The target classifier returns byte-identical strings or `null` for every committed `tests/license.js:10-181` example, including valid SPDX pass-through, guessed `*` results, `Custom: <capture>`, `Public Domain`, `Undefined`, and fallback `null`. [SA-002]
- **SC-002**: The target license-file detector returns the same selected filename sequence as committed `tests/license-files-test.js:10-87` for empty/no-match, extension-insensitive, case-insensitive, multiple-candidate, and first-`LICENSE-*` cases. [SA-003]
- **SC-003**: For fixtures covering package license, missing package license with README, and license-file-over-custom-url, the target package record has the same `licenses` value as source expectations, including `MIT*` for the committed license-file-over-custom-url assertion. [SA-004]
- **SC-004**: For Analyzer golden snippets, downstream renderers receiving this slice's values can reproduce exact license literals `UNKNOWN`, `MIT*`, and `ISC` in the recorded output contexts. [SA-006]
- **SC-005**: Differential-parity baseline: for committed classifier and license-file fixtures, the target outputs are byte-identical to source expected values across SPDX pass-through, regex classification, fallback/no-match, and license-file precedence edge cases. [SA-002] [SA-003]

## Assumptions

- F-003 provides package records, package metadata fields, package paths, initial `UNKNOWN` license values, and lexical package-key sorting before this slice's results are rendered. [SA-004] [SA-005]
- F-006 owns renderer-specific formatting, including JSON/tree/CSV/Markdown byte layout and stdout-vs-file newline differences. [SA-005] [SA-006]
- The Go implementation will either use an SPDX parser compatible with `spdx-expression-parse ^3.0.0` for the documented examples or faithfully reproduce the same success/failure behavior for this slice. [SA-007]
- Sourcebot MCP was unavailable during Analyzer reconstruction; line evidence should be rechecked against the indexed pinned commit if Sourcebot later becomes available, but the current spec remains grounded in Analyzer artifacts and committed source/test references. [SA-001] [SA-006]
- No final live differential run is part of this agent flow; any live source-vs-target comparison is a later manual maintainer step. [SA-006]

## Risks and Open Questions

- SPDX behavior is delegated in the source to `spdx-expression-parse`; choosing or implementing a compatible Go behavior is high-risk and must be validated against the recorded examples. [SA-007]
- Full byte fixture for the license-file-over-custom-url path was not captured because live source execution is prohibited; only the committed expectation `MIT*` is available for this slice. [SA-004] [SA-006]
