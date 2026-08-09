# BUILD — CLI entry point and invocation model

## 1. Implementation Scope

- Changed: `target/go.mod`, `target/cmd/license-checker/main.go`, `target/internal/app`, `target/internal/cli`, `target/internal/diagnostics`.
- Covered: executable `license-checker`, single-command invocation, `--help`/`-h`, `--version`/`-v`, `--failOn` + `--onlyAllow` preflight conflict, comma warnings, stdout/stderr routing, and downstream handoff into scan/render pipeline.

## 2. Build Steps

1. Created Go module rooted at `target/` and thin `cmd/license-checker` entry point.
2. Implemented normalized option parsing and preflight decision order in `internal/cli` and `internal/diagnostics`.
3. Wired `internal/app.Run` to parse, preflight, scan, policy-filter, format, and emit output.

## 3. Test Evidence

- Tests added: `internal/diagnostics/diagnostics_test.go` plus binary verification for help/version.
- Commands run from `target/`: `gofmt -w ./cmd ./internal`; `go test ./...`; `go vet ./...`; `go build ./...`; `go build -o /tmp/license-checker ./cmd/license-checker`; `/tmp/license-checker --help`; `/tmp/license-checker --version`.
- Observed: tests/build/vet passed; help exit `0` with stderr only; version exit `1` with stderr `25.0.1\n` and no stdout.

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

- Full `nopt` malformed/unknown-argument edge behavior remains a documented compatibility risk.
- Help final trailing spacing has no live byte fixture, as declared by SPEC/ARCHITECTURE.

## 5. Traceability

- F-001 FR-001–FR-011 map to `cmd/license-checker/main.go`, `internal/app/app.go`, `internal/cli/cli.go`, and `internal/diagnostics/diagnostics.go`.
- Architecture composition-root decisions are implemented by `app.Run` and a business-logic-free `main.go`.
