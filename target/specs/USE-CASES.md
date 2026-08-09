# USE-CASES — davglass/license-checker Analyzer Dispatch

## Functionality Checklist

- [ ] **F-001 — CLI entry point and invocation model**
  - **Status:** ready
  - **Slug:** `cli-entrypoint-invocation`
  - **Dependencies:** none
  - **Relevant analysis slice:** Entry Points; Inputs and Outputs; Requirements RT-001/RT-007
  - **Source references:** `package.json:91-94`; `bin/license-checker:1-72`; `README.md:69-90`
  - **Evidence Summary:** Single executable `license-checker`; no subcommands; help exits `0`, version exits `1`, failOn/onlyAllow conflict exits `1`.
  - **Open Questions:** none blocking.

- [ ] **F-002 — Argument parsing, defaults, and custom format loading**
  - **Status:** ready
  - **Slug:** `argument-parsing-defaults-custom-format`
  - **Dependencies:** F-001
  - **Relevant analysis slice:** CLI flag inventory; Input-Field Inventory; Requirements RT-001/RT-006
  - **Source references:** `lib/args.js:7-98`; `lib/index.js:264-270`; `lib/index.js:593-608`; `tests/test.js:389-476`, `690-716`
  - **Evidence Summary:** `nopt` known flags/shorts; defaults for color/start/relative/direct; custom JSON parsing returns parsed object or `Error`.
  - **Open Questions:** Exact nopt edge behavior should be pinned in Spec if Go uses another parser.

- [ ] **F-003 — Dependency scan and metadata flattening**
  - **Status:** ready
  - **Slug:** `dependency-scan-metadata-flattening`
  - **Dependencies:** F-001, F-002
  - **Relevant analysis slice:** Input-Field Inventory; Output Field Contract; Domain Map; Flow Map; Requirements RT-002/RT-003
  - **Source references:** `lib/index.js:27-110`; `lib/index.js:235-309`; `tests/test.js:19-45`, `480-530`
  - **Evidence Summary:** `read-installed` tree is flattened into lexically sorted `name@version` records; repository URLs normalized; package paths and dependency paths populated conditionally.
  - **Open Questions:** `read-installed` graph/dedupe/extraneous/root behavior is high-risk but sufficiently identified for a Spec-Writer to require equivalent behavior.

- [ ] **F-004 — License detection, classification, and license file precedence**
  - **Status:** ready
  - **Slug:** `license-detection-classification-files`
  - **Dependencies:** F-003
  - **Relevant analysis slice:** Value-Transformation Tables; Input-Field Inventory; Requirements RT-004
  - **Source references:** `lib/license.js:1-83`; `lib/license-files.js:3-29`; `lib/index.js:112-230`; `tests/license.js:10-181`; `tests/license-files-test.js:10-87`
  - **Evidence Summary:** Ordered SPDX-pass-through/regex/fallback classifier and ordered license filename detector are fully transcribed, including fallback `null` and sentinel literals.
  - **Open Questions:** Full byte fixture for license-file-over-custom-url path not captured because live execution is prohibited; committed assertion expects `MIT*`.

- [ ] **F-005 — Filtering, restriction, and license policy decisions**
  - **Status:** ready
  - **Slug:** `filtering-restriction-license-policy`
  - **Dependencies:** F-003, F-004
  - **Relevant analysis slice:** Value-Transformation Tables; Risk Map; Requirements RT-005/RT-007
  - **Source references:** `lib/index.js:313-458`; `tests/test.js:137-289`, `555-607`; `tests/failOn-test.js:7-42`; `tests/packages-test.js:8-44`
  - **Evidence Summary:** Covers `--unknown`, `--onlyunknown`, `--exclude`, `--packages`, `--excludePackages`, `--excludePrivatePackages`, `--failOn`, and `--onlyAllow`, including exact stderr snippets.
  - **Open Questions:** SPDX satisfaction should use an equivalent Go SPDX library or an explicitly tested compatibility layer.

- [ ] **F-006 — Output rendering and side-effect outputs**
  - **Status:** ready
  - **Slug:** `output-rendering-side-effects`
  - **Dependencies:** F-003, F-004, F-005
  - **Relevant analysis slice:** Golden Output Snippets; Output Field Contract; Source-Dependency Contracts; Requirements RT-006
  - **Source references:** `bin/license-checker:74-104`; `lib/index.js:471-627`; `README.md:24-67`; `tests/test.js:37-43`, `70-104`, `609-688`
  - **Evidence Summary:** Tree, JSON, CSV, Markdown, summary, `--out`, and `--files` contracts are reconstructed with field order, sentinels, and newline caveats.
  - **Open Questions:** Full byte-exact stdout fixtures for JSON/full tree/summary/color/files are not present in committed tests and were not generated because the source CLI was not executed.

- [ ] **F-007 — Diagnostics, error handling, and exit codes**
  - **Status:** ready
  - **Slug:** `diagnostics-errors-exit-codes`
  - **Dependencies:** F-001, F-005, F-006
  - **Relevant analysis slice:** Inputs and Outputs; Golden Output Snippets; Risk Map; Requirements RT-007
  - **Source references:** `bin/license-checker:17-72`; `lib/index.js:437-467`; `lib/index.js:624`; `README.md:158-170`; `tests/bin-test.js:7-17`; `tests/failOn-test.js:7-42`
  - **Evidence Summary:** Help/version/fail/mutual-exclusion/warning/error output paths and exit codes are explicit in source/tests; debug namespace documented.
  - **Open Questions:** Error object formatting for `console.error(err)` depends on Node runtime and was not byte-captured.

## Status

All functionality items are marked `ready` for Spec-Writer dispatch because core behavior and source references are reconstructed. Parity risks are preserved as open questions rather than hidden assumptions, especially for byte-exact golden fixtures not committed by the source project.

## Slug

Use the listed slug as the intended output folder under `target/specs/<slug>/`.

## Evidence Summary

Primary evidence is from pinned commit `de6e9a42513aa38a58efc6b202ee5281ed61f486`:

- CLI and parsing: `bin/license-checker`, `lib/args.js`, `README.md` options.
- Core pipeline: `lib/index.js`.
- License transforms: `lib/license.js`, `lib/license-files.js`.
- Committed behavior examples: `README.md`, `tests/test.js`, `tests/license.js`, `tests/license-files-test.js`, `tests/failOn-test.js`, `tests/packages-test.js`, `tests/bin-test.js`.

## Open Questions

- Sourcebot service was unavailable for direct indexed queries; re-run evidence check if the MCP becomes available.
- Do not treat missing live byte fixtures as source execution permission. Any final live differential run remains manual/out-of-flow.
- `read-installed`, SPDX packages, and `treeify` are the highest-risk third-party semantics for a Go rewrite.
