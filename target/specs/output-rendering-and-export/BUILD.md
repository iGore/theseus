# BUILD: Output Rendering and Export (F-006)

## 1. Implementation Scope

- Changed modules:
  - `target/internal/core/model.go`
  - `target/internal/core/render_service.go`
  - `target/internal/core/service_render_mode_test.go`
  - `target/internal/render/{model.go,common.go,tree.go,json.go,csv.go,markdown.go,summary.go,files_export.go}`
  - `target/test/integration/render_export/*.go`
  - `target/test/fixtures/render/f006_modes/*`
  - `target/test/fixtures/render/README.md`
  - `target/test/integration/render_export/README.md`
  - `target/test/integration/render_export/BUILD_EVIDENCE.md`
  - `target/specs/output-rendering-and-export/{PLAN.md,TASKS.md}`
- F-006 requirements covered:
  - FR-001/FR-002: five render modes + strict mode precedence
  - FR-003/FR-004: stdout default and `--out` write with parent dir creation
  - FR-005: `--files` license copy flow
  - FR-006: colorized keys limited to tree + terminal-compatible path
  - FR-008: markdown/files behaviors retained

## 2. Build Steps

1. Added render mode contract (`RenderMode`, `RenderOptions`) and `RenderService` execution path in `internal/core`.
2. Implemented deterministic renderer functions (sorted key iteration) in `internal/render`.
3. Implemented output sink + export helpers (`WriteOutputFile`, `ExportLicenseFiles`) using stdlib file APIs.
4. Added unit precedence tests and fixture-driven integration tests for renderer parity, sink behavior, and export side effects.
5. Updated `PLAN.md` and `TASKS.md` checklists to reflect completed F-006 work.

## 3. Test Evidence

- Tests added/updated:
  - `internal/core/service_render_mode_test.go`
  - `test/integration/render_export/render_modes_golden_test.go`
  - `test/integration/render_export/output_path_creation_test.go`
  - `test/integration/render_export/stdout_default_test.go`
  - `test/integration/render_export/output_content_parity_test.go`
  - `test/integration/render_export/files_export_enabled_test.go`
  - `test/integration/render_export/files_export_omitted_test.go`
- Command run:
  - `go test ./internal/core ./internal/render ./test/integration/render_export`
- Observed result:
  - `ok theseus/target/internal/core (cached)`
  - `? theseus/target/internal/render [no test files]`
  - `ok theseus/target/test/integration/render_export 0.419s`

## 4. Residual Risks

- Open clarification retained from plan: exact TTY detection policy semantics beyond explicit `IsTerminal` option wiring.
- Export collision policy currently follows key-scoped directory layout (`<files>/<name@version>/<basename>`), which is compatibility-safe for this slice but may require tightening if upstream source parity demands a different artifact layout.
- Full repository-wide `go test ./...` was not used as F-006 scope verification; this build used focused F-006 test targets.

## 5. Traceability

- FR-001/FR-002, US1, SC-001
  - Code: `internal/core/render_service.go`, `internal/render/*`
  - Tests: `internal/core/service_render_mode_test.go`, `test/integration/render_export/render_modes_golden_test.go`
- FR-003/FR-004, US2, SC-002/SC-004
  - Code: `internal/core/render_service.go`, `internal/render/files_export.go`
  - Tests: `output_path_creation_test.go`, `stdout_default_test.go`, `output_content_parity_test.go`
- FR-005, US3, SC-003
  - Code: `internal/render/files_export.go`
  - Tests: `files_export_enabled_test.go`, `files_export_omitted_test.go`
- FR-006
  - Code: tree rendering path in `internal/core/render_service.go` + `internal/render/tree.go`
  - Notes: `target/specs/output-rendering-and-export/TASKS.md` builder notes
- FR-008
  - Behavior retained via markdown renderer + files export path; open documentation status preserved in PLAN/TASKS unresolved markers.
