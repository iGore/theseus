# FUNCTIONALITY-INDEX.md

## Functionality Checklist

- [ ] **F-001 — CLI options model and normalization**
  - **Status:** ready
  - **Slug:** `cli-options-model`
  - **Evidence Summary:** `lib/args.js:7-99`; CLI guards/help/version in `bin/license-checker:17-63`.
  - **Open Questions:** Confirm compatibility requirement for `--version` exiting with code `1`.

- [ ] **F-002 — Dependency graph loading and traversal control**
  - **Status:** ready
  - **Slug:** `dependency-graph-loading`
  - **Evidence Summary:** `read-installed` orchestration and depth/dev toggles in `lib/index.js:267-309`; README scan description `177-182`.
  - **Open Questions:** None blocking.

- [ ] **F-003 — Flatten dependency graph into module inventory**
  - **Status:** ready
  - **Slug:** `dependency-flattening`
  - **Evidence Summary:** recursive `flatten` metadata extraction and dedupe/circular guard in `lib/index.js:27-110`, `235-259`.
  - **Open Questions:** None blocking.

- [ ] **F-004 — License detection and normalization pipeline**
  - **Status:** ready
  - **Slug:** `license-detection-pipeline`
  - **Evidence Summary:** license source precedence and heuristics in `lib/index.js:112-172`, `lib/license.js:22-83`, `lib/license-files.js:3-29`; parser tests in `tests/license.js`.
  - **Open Questions:** None blocking.

- [ ] **F-005 — Filtering, package restriction, and policy exits**
  - **Status:** ready
  - **Slug:** `filtering-and-policy-enforcement`
  - **Evidence Summary:** exclude/SPDX handling, package filters, private exclusion, `failOn`/`onlyAllow` exits in `lib/index.js:313-457`; behavior tests in `tests/test.js`, `tests/packages-test.js`, `tests/failOn-test.js`.
  - **Open Questions:** Clarify whether substring-based `onlyAllow` semantics must be preserved exactly.

- [ ] **F-006 — Output rendering and persistence**
  - **Status:** ready
  - **Slug:** `output-rendering-and-export`
  - **Evidence Summary:** mode selection in `bin/license-checker:84-104`; serializers in `lib/index.js:475-591`; file exports in `lib/index.js:610-627`.
  - **Open Questions:** Decide whether under-documented modes (`markdown`, `files`) should stay hidden or be documented.

- [ ] **F-007 — Custom format and JSON config loading**
  - **Status:** ready
  - **Slug:** `custom-format-config`
  - **Evidence Summary:** `customFormat` propagation/defaulting in `lib/index.js:95-106`, `181-221`; `--customPath` + `parseJson` in `264-266`, `593-608`; tests `tests/test.js:429-478`, `690-716`.
  - **Open Questions:** Should invalid `customPath` be a hard runtime error in rewrite?

- [ ] **F-008 — Ancillary Stack utility lifecycle**
  - **Status:** ready
  - **Slug:** `stack-utility-lifecycle`
  - **Evidence Summary:** standalone utility in `lib/stack.js:6-44`; repository search found no references.
  - **Open Questions:** Keep for compatibility/package surface or remove as dead code in rewrite target?

## Status

- Ready: F-001, F-002, F-003, F-004, F-005, F-006, F-007, F-008
- Blocked: none
