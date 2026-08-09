# BUILD — Output rendering and side effects

## 1. Implementation Scope

- Changed: `target/internal/render`, `target/internal/ordered`, app render/emission wiring, golden fixtures under `target/testdata/golden/`.
- Covered: tree, JSON, CSV, Markdown, summary, stdout newline behavior, `--out`, `--files`, missing-license warnings, custom format, component prefix.

## 2. Build Steps

1. Implemented ordered JSON/tree/CSV/Markdown/summary formatters.
2. Implemented side-effect precedence: `--files`, else `--out`, else stdout.
3. Added committed golden fixture files from `ANALYSIS.md` snippets and byte-level tests.

## 3. Test Evidence

- Tests added: `internal/render/render_test.go`.
- Commands run: `gofmt -w ./cmd ./internal`; `go test ./...`; `go vet ./...`; `go build ./...`.
- Observed: all commands passed.

## 3a. Golden-Output Parity Evidence

Golden fixtures are committed under `target/testdata/golden/`. Fixtures for rows without live captures are source-contract-derived from `ANALYSIS.md`, `SPEC.md`, and implemented target fixture data; no live Node.js source system was run.

| Mode / flag row | Fixture file(s) | Test name(s) | Status |
|---|---|---|---|
| Default invocation / tree | `tree_excerpt.golden`, `guessed_tree.golden` | `TestTreeGoldenSnippets`, `TestGoldenRenderSnippets` | identical |
| `--production` | `production_tree.golden` | `TestGoldenConformanceRowsFromSourceContracts/production_prunes_extraneous` | identical |
| `--development` | `development_tree.golden` | `TestGoldenConformanceRowsFromSourceContracts/development_keeps_extraneous_dependency` | identical |
| `--json` | `json_default.golden` | `TestGoldenStructuredModesAndSideEffects/json` | identical |
| `--csv` | `csv_default.golden`, `csv_stdout.golden`, `out_csv_file.golden` | `TestGoldenRenderSnippets/csv_default`, `TestGoldenStructuredModesAndSideEffects/csv`, `.../out_writes_raw_bytes_and_no_stdout` | identical |
| `--csvComponentPrefix` | `csv_component.golden` | `TestGoldenRenderSnippets/csv_component`, `TestGoldenStructuredModesAndSideEffects/csv_component_prefix` | identical |
| `--markdown` | `markdown_default.golden`, `markdown_cli.golden`, `markdown_stdout.golden` | `TestGoldenRenderSnippets/markdown_default`, `TestGoldenStructuredModesAndSideEffects/markdown` | identical |
| `--summary` | `summary_counts.golden` | `TestGoldenStructuredModesAndSideEffects/summary` | identical |
| `--out` | `out_csv_file.golden`, `empty_stdout.golden` | `TestGoldenStructuredModesAndSideEffects/out_writes_raw_bytes_and_no_stdout` | identical |
| `--files` | `files_warning.golden` plus copied `foo-LICENSE.txt` bytes | `TestFilesSideEffectAndWarning` | identical |
| `--start` | `start_path_tree.golden` | `TestGoldenConformanceRowsFromSourceContracts/start_path_propagates_into_output` | identical |
| `--unknown` | `unknown_dependency_path_tree.golden` | `TestGoldenConformanceRowsFromSourceContracts/unknown_rewrites_guessed_and_includes_dependencyPath` | identical |
| `--onlyunknown` | `onlyunknown_tree.golden` | `TestGoldenConformanceRowsFromSourceContracts/policy_filters_render_byte_fixtures/onlyunknown` | identical |
| `--exclude` | `exclude_tree.golden` | `TestGoldenConformanceRowsFromSourceContracts/policy_filters_render_byte_fixtures/exclude` | identical |
| `--failOn` | `failon_warning.golden`, `failon_violation_stderr.golden` | `TestPreflightGolden`, `TestGoldenConformanceRowsFromSourceContracts/failOn_stderr` | identical |
| `--onlyAllow` | `onlyallow_violation_stderr.golden` | `TestGoldenConformanceRowsFromSourceContracts/onlyAllow_stderr` | identical |
| `--packages` | `packages_tree.golden` | `TestGoldenConformanceRowsFromSourceContracts/policy_filters_render_byte_fixtures/packages` | identical |
| `--excludePackages` | `exclude_packages_tree.golden` | `TestGoldenConformanceRowsFromSourceContracts/policy_filters_render_byte_fixtures/excludePackages` | identical |
| `--excludePrivatePackages` | `exclude_private_empty_tree.golden` | `TestGoldenConformanceRowsFromSourceContracts/policy_filters_render_byte_fixtures/excludePrivatePackages` | identical |
| `--relativeLicensePath` | `relative_license_path_tree.golden` | `TestGoldenConformanceRowsFromSourceContracts/relativeLicensePath_output` | identical |
| `--customPath` | `custom_path_csv.golden` | `TestCustomPathCLIGolden` | identical |
| `--direct` | `direct_tree.golden` | `TestGoldenConformanceRowsFromSourceContracts/direct_depth_zero_provider_contract_output` | identical |
| Programmatic `customFormat` | `csv_custom.golden`, `markdown_custom_row.golden` | `TestGoldenRenderSnippets/csv_custom`, `TestGoldenRenderSnippets/markdown_custom` | identical |
| `--help`, `-h` | `help_stderr.golden` | `TestPreflightGolden` | identical |
| `--version`, `-v` | `version_stderr.golden` | `TestPreflightGolden` plus manual binary command | identical |

Normalization allowlist: none for whitespace, trailing newlines, sentinels, field order, stdout, stderr, or file bytes. The only normalization is the volatile `t.TempDir()` absolute root replacement with `/fixture/root` in `TestGoldenConformanceRowsFromSourceContracts/relativeLicensePath_output`; relative `licenseFile` and `noticeFile` bytes are compared exactly after that temp-root substitution.

## 4. Residual Risks

- Full live stdout for JSON/full tree/summary/color was not captured by Analyzer; tests cover source-derived contracts and snippets.
- Color output is compatible-only and not byte-pinned.

## 5. Traceability

- F-006 FR-001–FR-029 map to `internal/render/render.go` and `internal/ordered/ordered.go`.
