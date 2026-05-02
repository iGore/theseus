# BUILD: Custom Format and JSON Config Loading (F-007)

## 1. Implementation Scope

- Files changed (strict F-007 slice):
  - `target/internal/config/custom_format.go`
  - `target/internal/config/custom_format_test.go`
  - `target/internal/flatten/flatten.go`
  - `target/internal/cli/options.go`
  - `target/internal/cli/command.go`
  - `target/internal/core/model.go`
  - `target/test/integration/custom_format_inline_test.go`
  - `target/test/integration/custom_format_exclusion_test.go`
  - `target/test/fixtures/custom-format/valid_custom_path.json`
  - `target/specs/custom-format-config/PLAN.md`
  - `target/specs/custom-format-config/TASKS.md`
- Spec requirements covered:
  - FR-001, FR-002, FR-003 (custom field mapping + default fallback + false exclusion)
  - FR-004, FR-005, FR-006, FR-007 (customPath + parse matrix behavior)
  - FR-008, FR-009 (licenseText CSV normalization + copyright extraction helper)
  - FR-010 resolved per architecture compatibility-mode default for this PoC slice.

## 2. Build Steps

1. Read `SPEC.md`, `PLAN.md`, `TASKS.md`, and shared `ARCHITECTURE` for F-007 constraints.
2. Added `internal/config/custom_format.go` with:
   - `CustomFormatConfig` contract
   - `ParseJSON(path any) any` returning object-or-error value
   - `LoadCustomFormat` loader behavior for inline/customPath
   - deterministic custom-field mapping including default fallback and `false` exclusion
3. Extended F-007 boundary types in `internal/core/model.go` and option plumbing in CLI (`internal/cli/options.go`, `internal/cli/command.go`) for `customPath` ownership boundary.
4. Added enrichment helpers in `internal/flatten/flatten.go` for `licenseText` normalization and copyright extraction rule filtering.
5. Added tests and fixture for parser matrix, inline mapping, fallback, exclusion, and customPath loading.
6. Updated plan/task checklist statuses and FR-010 decision state in use-case docs.

## 3. Test Evidence

- Tests added/updated:
  - `target/internal/config/custom_format_test.go`
  - `target/test/integration/custom_format_inline_test.go`
  - `target/test/integration/custom_format_exclusion_test.go`
- Command run:
  - `go test ./...` (workdir: `target/`)
- Observed result:
  - All packages passed, including:
    - `ok theseus/target/internal/config`
    - `ok theseus/target/internal/flatten`
    - `ok theseus/target/test/integration`
  - No failing tests.

## 4. Residual Risks

- Custom fields are currently produced as a dedicated mapping helper and not yet threaded through all render modes end-to-end.
- FR-010 strict-mode alternative (hard-fail initialization) remains deferred; this slice intentionally keeps compatibility-mode semantics.

## 5. Traceability

- FR-001/FR-002/FR-003 -> `internal/config/custom_format.go` (`ApplyCustomFields`) + tests in `internal/config/custom_format_test.go` and `test/integration/custom_format_exclusion_test.go`.
- FR-004/FR-005/FR-006/FR-007 -> `internal/config/custom_format.go` (`ParseJSON`, `LoadCustomFormat`) + parser matrix tests.
- FR-008/FR-009 -> `internal/flatten/flatten.go` (`EnrichLicenseDerivedFields`, `extractCopyright`) + normalization assertion in config tests.
- FR-010 -> decision recorded in `target/specs/custom-format-config/PLAN.md` and implemented with compatibility-preserving parse behavior.
